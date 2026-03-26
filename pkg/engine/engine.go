package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/franklinkim/bouncer/internal/collector"
	"github.com/franklinkim/bouncer/pkg/schema"
)

// opaFinding represents a single finding from an OPA evaluation.
type opaFinding struct {
	RuleID   string `json:"rule_id"`
	Evidence string `json:"evidence"`
}

// opaOutput represents the top-level OPA JSON output structure.
type opaOutput struct {
	Result []struct {
		Expressions []struct {
			Value map[string]struct {
				Findings []opaFinding `json:"findings"`
			} `json:"value"`
		} `json:"expressions"`
	} `json:"result"`
}

// CommandRunner abstracts command execution for testability.
type CommandRunner func(ctx context.Context, name string, args ...string) ([]byte, error)

// defaultRunner executes a command using os/exec.
func defaultRunner(ctx context.Context, name string, args ...string) ([]byte, error) {
	return exec.CommandContext(ctx, name, args...).CombinedOutput()
}

// Engine evaluates Rego policies against facts via the OPA binary.
type Engine struct {
	Policies [][]byte
	Rules    []schema.Rule
	OPAPath  string
	RunCmd   CommandRunner
}

// NewEngine creates an Engine, locating the OPA binary in PATH.
func NewEngine(policies [][]byte, rules []schema.Rule) (*Engine, error) {
	opaPath, err := exec.LookPath("opa")
	if err != nil {
		return nil, fmt.Errorf("opa not found in PATH: %w", err)
	}

	return &Engine{
		Policies: policies,
		Rules:    rules,
		OPAPath:  opaPath,
		RunCmd:   defaultRunner,
	}, nil
}

// Evaluate runs OPA against the provided facts and returns a ScanResult.
func (e *Engine) Evaluate(ctx context.Context, facts *schema.Facts, collectorResults []collector.Result) (*schema.ScanResult, error) {
	tmpDir, err := os.MkdirTemp("", "bouncer-opa-*")
	if err != nil {
		return nil, fmt.Errorf("creating temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Write policy files.
	for i, policy := range e.Policies {
		policyPath := filepath.Join(tmpDir, fmt.Sprintf("policy_%d.rego", i))
		if err := os.WriteFile(policyPath, policy, 0600); err != nil {
			return nil, fmt.Errorf("writing policy file: %w", err)
		}
	}

	// Write facts as JSON input.
	factsJSON, err := json.Marshal(facts)
	if err != nil {
		return nil, fmt.Errorf("marshaling facts: %w", err)
	}
	factsPath := filepath.Join(tmpDir, "facts.json")
	if err := os.WriteFile(factsPath, factsJSON, 0600); err != nil {
		return nil, fmt.Errorf("writing facts file: %w", err)
	}

	// Run OPA eval.
	out, err := e.RunCmd(ctx, e.OPAPath,
		"eval",
		"-d", tmpDir,
		"-i", factsPath,
		"--format", "json",
		"data.bouncer",
	)
	if err != nil {
		return nil, fmt.Errorf("running opa eval: %w\noutput: %s", err, string(out))
	}

	// Parse OPA output.
	findings, err := parseOPAFindings(out)
	if err != nil {
		return nil, fmt.Errorf("parsing opa output: %w", err)
	}

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

// parseOPAFindings extracts all findings from OPA JSON output.
func parseOPAFindings(data []byte) ([]opaFinding, error) {
	var out opaOutput
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("unmarshaling opa output: %w", err)
	}

	var findings []opaFinding
	for _, r := range out.Result {
		for _, expr := range r.Expressions {
			for _, category := range expr.Value {
				findings = append(findings, category.Findings...)
			}
		}
	}
	return findings, nil
}

// buildCategoryStatusMap creates a map from collector name (category) to its status.
func buildCategoryStatusMap(results []collector.Result) map[string]collector.CollectorStatus {
	m := make(map[string]collector.CollectorStatus, len(results))
	for _, r := range results {
		m[r.Name] = r.Status
	}
	return m
}
