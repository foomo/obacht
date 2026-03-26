package engine

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/open-policy-agent/opa/v1/rego"

	"github.com/franklinkim/bouncer/internal/collector"
	"github.com/franklinkim/bouncer/pkg/schema"
)

// opaFinding represents a single finding from an OPA evaluation.
type opaFinding struct {
	RuleID   string `json:"rule_id"`
	Evidence string `json:"evidence"`
}

// Engine evaluates Rego policies against facts via the embedded OPA library.
type Engine struct {
	Policies [][]byte
	Rules    []schema.Rule
}

// NewEngine creates an Engine with the given policies and rules.
func NewEngine(policies [][]byte, rules []schema.Rule) (*Engine, error) {
	return &Engine{
		Policies: policies,
		Rules:    rules,
	}, nil
}

// Evaluate runs the embedded OPA engine against the provided facts and returns a ScanResult.
func (e *Engine) Evaluate(ctx context.Context, facts *schema.Facts, collectorResults []collector.Result) (*schema.ScanResult, error) {
	// Marshal facts to generic interface{} for OPA input.
	factsJSON, err := json.Marshal(facts)
	if err != nil {
		return nil, fmt.Errorf("marshaling facts: %w", err)
	}

	var input any
	if err := json.Unmarshal(factsJSON, &input); err != nil {
		return nil, fmt.Errorf("unmarshaling facts to interface: %w", err)
	}

	// Build rego options.
	opts := []func(*rego.Rego){
		rego.Query("data.bouncer"),
		rego.Input(input),
	}

	for i, p := range e.Policies {
		opts = append(opts, rego.Module(fmt.Sprintf("policy_%d.rego", i), string(p)))
	}

	// Evaluate policies.
	rs, err := rego.New(opts...).Eval(ctx)
	if err != nil {
		return nil, fmt.Errorf("evaluating rego: %w", err)
	}

	// Extract findings from the result set.
	findings := extractFindings(rs)

	// Build category -> collector status map.
	categoryStatus := buildCategoryStatusMap(collectorResults)

	// Build finding set keyed by rule_id for quick lookup.
	findingMap := make(map[string]opaFinding)
	for _, f := range findings {
		findingMap[f.RuleID] = f
	}

	// Build check results.
	results := make([]schema.CheckResult, 0, len(e.Rules))
	for _, rule := range e.Rules {
		cr := schema.CheckResult{
			RuleID:      rule.ID,
			Title:       rule.Title,
			Severity:    rule.Severity,
			Category:    rule.Category,
			Remediation: rule.Remediation,
		}

		if status, ok := categoryStatus[rule.Category]; ok {
			switch status {
			case collector.StatusSkipped:
				cr.Status = schema.StatusSkip
				results = append(results, cr)

				continue
			case collector.StatusError:
				cr.Status = schema.StatusError
				results = append(results, cr)

				continue
			case collector.StatusOK:
				// Collector succeeded; fall through to OPA evaluation below.
			}
		}

		if f, found := findingMap[rule.ID]; found {
			cr.Status = schema.StatusFail
			cr.Evidence = f.Evidence
		} else {
			cr.Status = schema.StatusPass
		}

		results = append(results, cr)
	}

	scanResult := schema.NewScanResult(results)

	return &scanResult, nil
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

// buildCategoryStatusMap creates a map from collector name (category) to its status.
func buildCategoryStatusMap(results []collector.Result) map[string]collector.CollectorStatus {
	m := make(map[string]collector.CollectorStatus, len(results))
	for _, r := range results {
		m[r.Name] = r.Status
	}

	return m
}
