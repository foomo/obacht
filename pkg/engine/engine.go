package engine

import (
	"context"
	"fmt"
	"io/fs"
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

// EvaluateWithFS runs input scripts (with include preprocessing against fsys)
// and evaluates rego policies for the given rule files.
func EvaluateWithFS(ctx context.Context, fsys fs.FS, ruleFiles []schema.RulesFile, onProgress ...ProgressFunc) (*schema.ScanResult, error) {
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

		groupResults, err := evaluateGroup(ctx, fsys, g)
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

// Evaluate is the legacy entrypoint, using a nil fs.FS. Include directives
// will fail to resolve under this entrypoint. Prefer EvaluateWithFS.
func Evaluate(ctx context.Context, ruleFiles []schema.RulesFile, onProgress ...ProgressFunc) (*schema.ScanResult, error) {
	return EvaluateWithFS(ctx, nil, ruleFiles, onProgress...)
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
func evaluateGroup(ctx context.Context, fsys fs.FS, g ruleGroup) ([]schema.CheckResult, error) {
	// Use the first rule's ID as the expected envelope rule_id when the group
	// contains exactly one rule. For multi-rule groups (shared per-category
	// input), expectedID stays empty and the legacy bare-JSON path applies.
	expectedID := ""
	if len(g.Rules) == 1 {
		expectedID = g.Rules[0].ID
	}

	scriptPath := ""
	if len(g.Rules) > 0 {
		scriptPath = fmt.Sprintf("inputs/%s/%s.sh", g.Rules[0].Category, g.Rules[0].ID)
	}

	var inputResult runner.InputResult
	if fsys != nil {
		inputResult = runner.RunInputForRule(ctx, fsys, scriptPath, g.Input, expectedID)
	} else {
		inputResult = runner.RunInput(ctx, g.Input)
	}

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
				cr.Evidence = inputResult.SkipReason
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

	findings, skips, err := evalRego(ctx, g.Policy, inputResult.Data)
	if err != nil {
		return nil, err
	}

	findingMap := make(map[string][]opaFinding)
	for _, f := range findings {
		findingMap[f.RuleID] = append(findingMap[f.RuleID], f)
	}

	skipMap := make(map[string][]opaFinding)
	for _, s := range skips {
		skipMap[s.RuleID] = append(skipMap[s.RuleID], s)
	}

	for _, rule := range g.Rules {
		cr := schema.CheckResult{
			RuleID:      rule.ID,
			Title:       rule.Title,
			Severity:    rule.Severity,
			Category:    rule.Category,
			Remediation: rule.Remediation,
		}

		switch {
		case len(findingMap[rule.ID]) > 0:
			cr.Status = schema.StatusFail
			cr.Evidence = joinEvidence(findingMap[rule.ID])
		case len(skipMap[rule.ID]) > 0:
			cr.Status = schema.StatusSkip
			cr.Evidence = joinEvidence(skipMap[rule.ID])
		default:
			cr.Status = schema.StatusPass
		}

		results = append(results, cr)
	}

	return results, nil
}

// joinEvidence sorts and joins evidence strings from multiple records for the
// same rule_id with "; ". Sorting keeps output stable across runs.
func joinEvidence(items []opaFinding) string {
	parts := make([]string, 0, len(items))
	for _, it := range items {
		parts = append(parts, it.Evidence)
	}

	sort.Strings(parts)

	return strings.Join(parts, "; ")
}

// evalRego evaluates a rego policy string against the given input data and
// returns both findings and skips collections.
func evalRego(ctx context.Context, policy string, input any) ([]opaFinding, []opaFinding, error) {
	opts := []func(*rego.Rego){
		rego.Query("data.obacht"),
		rego.Input(input),
		rego.Module("policy.rego", policy),
	}

	rs, err := rego.New(opts...).Eval(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("evaluating rego: %w", err)
	}

	return extractCollection(rs, "findings"), extractCollection(rs, "skips"), nil
}

// extractCollection walks the OPA result set and collects all entries from the
// named collection (e.g. "findings" or "skips"). Each entry must be a map with
// a "rule_id" string field; "evidence" is optional.
func extractCollection(rs rego.ResultSet, name string) []opaFinding {
	var out []opaFinding

	if len(rs) == 0 {
		return out
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

				raw, ok := catMap[name]
				if !ok {
					continue
				}

				slice, ok := raw.([]any)
				if !ok {
					continue
				}

				for _, rf := range slice {
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
						out = append(out, f)
					}
				}
			}
		}
	}

	return out
}
