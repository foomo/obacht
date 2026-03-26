package schema

// Severity represents the severity level of a finding.
type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityHigh     Severity = "high"
	SeverityWarn     Severity = "warn"
	SeverityInfo     Severity = "info"
)

// Status represents the result status of a check.
type Status string

const (
	StatusPass  Status = "pass"
	StatusFail  Status = "fail"
	StatusSkip  Status = "skip"
	StatusError Status = "error"
)

// Rule defines a single policy rule.
type Rule struct {
	ID          string   `json:"id" yaml:"id"`
	Title       string   `json:"title" yaml:"title"`
	Severity    Severity `json:"severity" yaml:"severity"`
	Category    string   `json:"category" yaml:"category"`
	Description string   `json:"description" yaml:"description"`
	Remediation string   `json:"remediation" yaml:"remediation"`
}

// RulesFile is the top-level structure for a rules YAML/JSON file.
type RulesFile struct {
	Rules []Rule `json:"rules" yaml:"rules"`
}

// CheckResult is the outcome of evaluating a single rule against collected facts.
type CheckResult struct {
	RuleID      string   `json:"rule_id"`
	Title       string   `json:"title"`
	Severity    Severity `json:"severity"`
	Category    string   `json:"category"`
	Status      Status   `json:"status"`
	Evidence    string   `json:"evidence,omitempty"`
	Remediation string   `json:"remediation,omitempty"`
}

// Summary aggregates counts from a set of check results.
type Summary struct {
	Total    int `json:"total"`
	Passed   int `json:"passed"`
	Failed   int `json:"failed"`
	Skipped  int `json:"skipped"`
	Errors   int `json:"errors"`
	Critical int `json:"critical"`
	High     int `json:"high"`
	Warn     int `json:"warn"`
	Info     int `json:"info"`
}

// ScanResult is the complete output of a bouncer scan.
type ScanResult struct {
	SchemaVersion string        `json:"schema_version"`
	Results       []CheckResult `json:"results"`
	Summary       Summary       `json:"summary"`
}

// NewScanResult creates a ScanResult from check results, computing the Summary.
func NewScanResult(results []CheckResult) ScanResult {
	s := Summary{
		Total: len(results),
	}
	for _, r := range results {
		switch r.Status {
		case StatusPass:
			s.Passed++
		case StatusFail:
			s.Failed++
		case StatusSkip:
			s.Skipped++
		case StatusError:
			s.Errors++
		}

		switch r.Severity {
		case SeverityCritical:
			s.Critical++
		case SeverityHigh:
			s.High++
		case SeverityWarn:
			s.Warn++
		case SeverityInfo:
			s.Info++
		}
	}

	return ScanResult{
		SchemaVersion: "1.0",
		Results:       results,
		Summary:       s,
	}
}
