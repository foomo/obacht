package engine

import (
	"context"
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/open-policy-agent/opa/v1/rego"

	"github.com/foomo/obacht/internal/runner"
	"github.com/foomo/obacht/pkg/schema"
)

// opaFinding represents a single finding from an OPA evaluation.
type opaFinding struct {
	RuleID   string `json:"rule_id"`
	Evidence string `json:"evidence"`
}

// ruleGroup is a set of rules that share the same input script and policy.
type ruleGroup struct {
	Input  string
	Policy string
	Rules  []schema.Rule
}

// Evaluate runs input scripts and evaluates rego policies for the given rule files.
// An optional ProgressFunc callback receives events as each rule group starts and completes.
func Evaluate(ctx context.Context, ruleFiles []schema.RulesFile, onProgress ...ProgressFunc) (*schema.ScanResult, error) {
	groups := buildRuleGroups(ruleFiles)

	var notify ProgressFunc
	if len(onProgress) > 0 && onProgress[0] != nil {
		notify = onProgress[0]
	}

	var results []schema.CheckResult

	for i, g := range groups {
		cat := groupCategory(g)

		if notify != nil {
			notify(ProgressEvent{
				Kind:       EventGroupStart,
				Category:   cat,
				RuleCount:  len(g.Rules),
				GroupIndex: i,
				GroupTotal: len(groups),
			})
		}

		groupResults, err := evaluateGroup(ctx, g)
		if err != nil {
			return nil, err
		}

		results = append(results, groupResults...)

		if notify != nil {
			notify(ProgressEvent{
				Kind:       EventGroupDone,
				Category:   cat,
				RuleCount:  len(g.Rules),
				Results:    groupResults,
				GroupIndex: i,
				GroupTotal: len(groups),
			})
		}
	}

	scanResult := schema.NewScanResult(results)

	return &scanResult, nil
}

// groupCategory returns the category from the first rule in the group.
func groupCategory(g ruleGroup) string {
	if len(g.Rules) > 0 {
		return g.Rules[0].Category
	}

	return "unknown"
}

// buildRuleGroups organizes rules into groups that share the same input.
// Rules within a RulesFile that resolve to the same input script are merged
// into a single group — their effective policies are concatenated into one
// rego module so the input runs once and OPA evaluates once per group.
//
// Rules whose effective policy declares its own `package` are isolated into
// their own groups, since concatenating multiple package declarations would
// produce an invalid module.
func buildRuleGroups(ruleFiles []schema.RulesFile) []ruleGroup {
	var groups []ruleGroup

	for _, rf := range ruleFiles {
		type bucket struct {
			input    string
			rules    []schema.Rule
			policies []string
			category string
		}

		var order []string

		buckets := map[string]*bucket{}

		for _, rule := range rf.Rules {
			input := resolveField(rule.Input, rf.Input)
			policy := resolveField(rule.Policy, rf.Policy)

			if strings.HasPrefix(strings.TrimSpace(policy), "package ") {
				groups = append(groups, ruleGroup{
					Input:  input,
					Policy: preparePolicy(policy, rule.Category),
					Rules:  []schema.Rule{rule},
				})

				continue
			}

			b, ok := buckets[input]
			if !ok {
				b = &bucket{input: input, category: rule.Category}
				buckets[input] = b
				order = append(order, input)
			}

			b.rules = append(b.rules, rule)

			if policy != "" && !slices.Contains(b.policies, policy) {
				b.policies = append(b.policies, policy)
			}
		}

		for _, key := range order {
			b := buckets[key]
			groups = append(groups, ruleGroup{
				Input:  b.input,
				Policy: preparePolicy(strings.Join(b.policies, "\n\n"), b.category),
				Rules:  b.rules,
			})
		}
	}

	return groups
}

// preparePolicy ensures the policy string has a package declaration and rego.v1 import.
// If the policy already starts with "package", it is returned unchanged.
// Otherwise, the package name is derived from category.
func preparePolicy(policy, category string) string {
	if policy == "" {
		return ""
	}

	trimmed := strings.TrimSpace(policy)
	if strings.HasPrefix(trimmed, "package ") {
		return policy
	}

	pkg := category
	if pkg == "" {
		pkg = "default"
	}

	return fmt.Sprintf("package obacht.%s\nimport rego.v1\n\n%s", pkg, policy)
}

// resolveField returns the rule-level value if set, otherwise the file-level fallback.
func resolveField(ruleLevel, fileLevel string) string {
	if ruleLevel != "" {
		return ruleLevel
	}

	return fileLevel
}

// evaluateGroup runs the input script and rego policy for a group of rules.
func evaluateGroup(ctx context.Context, g ruleGroup) ([]schema.CheckResult, error) {
	// Run input script.
	inputResult := runner.RunInput(ctx, g.Input)

	results := make([]schema.CheckResult, 0, len(g.Rules))

	// If input was skipped or errored, mark all rules accordingly.
	if inputResult.Status != runner.StatusOK {
		for _, rule := range g.Rules {
			cr := schema.CheckResult{
				RuleID:      rule.ID,
				Title:       rule.Title,
				Severity:    rule.Severity,
				Category:    rule.Category,
				Remediation: rule.Remediation,
			}

			switch inputResult.Status {
			case runner.StatusSkipped:
				cr.Status = schema.StatusSkip
			case runner.StatusError:
				cr.Status = schema.StatusError
				if inputResult.Error != nil {
					cr.Evidence = inputResult.Error.Error()
				}
			}

			results = append(results, cr)
		}

		return results, nil
	}

	// Evaluate rego policy.
	if g.Policy == "" {
		// No policy — all rules pass by default.
		for _, rule := range g.Rules {
			results = append(results, schema.CheckResult{
				RuleID:      rule.ID,
				Title:       rule.Title,
				Severity:    rule.Severity,
				Category:    rule.Category,
				Remediation: rule.Remediation,
				Status:      schema.StatusPass,
			})
		}

		return results, nil
	}

	findings, err := evalRego(ctx, g.Policy, inputResult.Data)
	if err != nil {
		return nil, err
	}

	// Group findings by rule_id. A single rule may produce multiple findings
	// (e.g., one per offending PATH dir); aggregate their evidence strings so
	// callers see every match rather than an arbitrary one.
	findingMap := make(map[string][]opaFinding)
	for _, f := range findings {
		findingMap[f.RuleID] = append(findingMap[f.RuleID], f)
	}

	// Build check results.
	for _, rule := range g.Rules {
		cr := schema.CheckResult{
			RuleID:      rule.ID,
			Title:       rule.Title,
			Severity:    rule.Severity,
			Category:    rule.Category,
			Remediation: rule.Remediation,
		}

		if fs, found := findingMap[rule.ID]; found {
			cr.Status = schema.StatusFail

			parts := make([]string, 0, len(fs))
			for _, f := range fs {
				parts = append(parts, f.Evidence)
			}

			sort.Strings(parts)
			cr.Evidence = strings.Join(parts, "; ")
		} else {
			cr.Status = schema.StatusPass
		}

		results = append(results, cr)
	}

	return results, nil
}

// evalRego evaluates a rego policy string against the given input data.
func evalRego(ctx context.Context, policy string, input any) ([]opaFinding, error) {
	opts := []func(*rego.Rego){
		rego.Query("data.obacht"),
		rego.Input(input),
		rego.Module("policy.rego", policy),
	}

	rs, err := rego.New(opts...).Eval(ctx)
	if err != nil {
		return nil, fmt.Errorf("evaluating rego: %w", err)
	}

	return extractFindings(rs), nil
}

// extractFindings walks the OPA result set and collects all findings.
func extractFindings(rs rego.ResultSet) []opaFinding {
	var findings []opaFinding

	if len(rs) == 0 {
		return findings
	}

	for _, result := range rs {
		for _, expr := range result.Expressions {
			categories, ok := expr.Value.(map[string]any)
			if !ok {
				continue
			}

			for _, catVal := range categories {
				catMap, ok := catVal.(map[string]any)
				if !ok {
					continue
				}

				rawFindings, ok := catMap["findings"]
				if !ok {
					continue
				}

				findingsSlice, ok := rawFindings.([]any)
				if !ok {
					continue
				}

				for _, rf := range findingsSlice {
					fm, ok := rf.(map[string]any)
					if !ok {
						continue
					}

					f := opaFinding{}
					if v, ok := fm["rule_id"].(string); ok {
						f.RuleID = v
					}

					if v, ok := fm["evidence"].(string); ok {
						f.Evidence = v
					}

					if f.RuleID != "" {
						findings = append(findings, f)
					}
				}
			}
		}
	}

	return findings
}
