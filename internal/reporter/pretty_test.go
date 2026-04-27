package reporter_test

import (
	"bytes"
	"regexp"
	"strings"
	"testing"

	"github.com/foomo/obacht/internal/reporter"
	"github.com/foomo/obacht/pkg/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var ansiRe = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripANSI(s string) string {
	return ansiRe.ReplaceAllString(s, "")
}

func sampleScanResult() *schema.ScanResult {
	results := []schema.CheckResult{
		{
			RuleID:      "SSH001",
			Title:       "SSH private key permissions",
			Severity:    schema.SeverityCritical,
			Category:    "SSH",
			Status:      schema.StatusFail,
			Evidence:    "~/.ssh has mode 0755",
			Remediation: "Run: chmod 700 ~/.ssh",
		},
		{
			RuleID:   "SSH002",
			Title:    "SSH agent running",
			Severity: schema.SeverityInfo,
			Category: "SSH",
			Status:   schema.StatusPass,
		},
		{
			RuleID:   "GIT001",
			Title:    "Git user configured",
			Severity: schema.SeverityWarn,
			Category: "Git",
			Status:   schema.StatusSkip,
		},
		{
			RuleID:      "GIT002",
			Title:       "Git hooks installed",
			Severity:    schema.SeverityHigh,
			Category:    "Git",
			Status:      schema.StatusError,
			Evidence:    "hook directory missing",
			Remediation: "Run: git init",
		},
	}
	sr := schema.NewScanResult(results)

	return &sr
}

func TestPrettyReporter_ContainsPassMarker(t *testing.T) {
	var buf bytes.Buffer

	r := &reporter.PrettyReporter{ShowPassing: true}
	err := r.Report(&buf, sampleScanResult())
	require.NoError(t, err)

	out := buf.String()

	// Pass marker (checkmark) should be present when ShowPassing is enabled.
	assert.Contains(t, out, "\u2713")
}

func TestPrettyReporter_HidesPassByDefault(t *testing.T) {
	var buf bytes.Buffer

	r := reporter.NewPrettyReporter()
	err := r.Report(&buf, sampleScanResult())
	require.NoError(t, err)

	out := buf.String()

	// Pass marker and the passing rule must not appear in the per-check listing.
	assert.NotContains(t, out, "\u2713")
	assert.NotContains(t, out, "SSH002")

	// Failing/skip/error checks still render.
	assert.Contains(t, out, "SSH001")
	assert.Contains(t, out, "GIT001")
	assert.Contains(t, out, "GIT002")

	// Summary line still reflects the full counts including the hidden pass.
	assert.Contains(t, out, "1 passed")
}

func TestPrettyReporter_SkipsEmptyCategoryHeader(t *testing.T) {
	results := []schema.CheckResult{
		{
			RuleID:   "TOOL001",
			Title:    "tool installed",
			Severity: schema.SeverityInfo,
			Category: "Tools",
			Status:   schema.StatusPass,
		},
		{
			RuleID:      "SSH001",
			Title:       "SSH private key permissions",
			Severity:    schema.SeverityCritical,
			Category:    "SSH",
			Status:      schema.StatusFail,
			Evidence:    "~/.ssh has mode 0755",
			Remediation: "Run: chmod 700 ~/.ssh",
		},
	}
	sr := schema.NewScanResult(results)

	var buf bytes.Buffer

	r := reporter.NewPrettyReporter()
	err := r.Report(&buf, &sr)
	require.NoError(t, err)

	out := stripANSI(buf.String())

	// The Tools category is all-pass, so its header must not be rendered.
	assert.NotContains(t, out, "Tools")
	// The SSH category still appears.
	assert.Contains(t, out, "SSH")
	assert.Contains(t, out, "SSH001")
}

func TestPrettyReporter_ContainsFailMarker(t *testing.T) {
	var buf bytes.Buffer

	r := reporter.NewPrettyReporter()
	err := r.Report(&buf, sampleScanResult())
	require.NoError(t, err)

	out := buf.String()

	// Fail marker (cross) should be present.
	assert.Contains(t, out, "\u2717")
}

func TestPrettyReporter_ContainsRuleIDs(t *testing.T) {
	var buf bytes.Buffer

	r := &reporter.PrettyReporter{ShowPassing: true}
	err := r.Report(&buf, sampleScanResult())
	require.NoError(t, err)

	out := buf.String()

	assert.Contains(t, out, "SSH001")
	assert.Contains(t, out, "SSH002")
	assert.Contains(t, out, "GIT001")
	assert.Contains(t, out, "GIT002")
}

func TestPrettyReporter_ContainsSummary(t *testing.T) {
	var buf bytes.Buffer

	r := reporter.NewPrettyReporter()
	err := r.Report(&buf, sampleScanResult())
	require.NoError(t, err)

	out := buf.String()

	assert.Contains(t, out, "Summary:")
	assert.Contains(t, out, "1 failed")
	assert.Contains(t, out, "1 passed")
	assert.Contains(t, out, "1 skipped")
	assert.Contains(t, out, "1 critical")
	assert.Contains(t, out, "1 high")
}

func TestPrettyReporter_ContainsEvidence(t *testing.T) {
	var buf bytes.Buffer

	r := reporter.NewPrettyReporter()
	err := r.Report(&buf, sampleScanResult())
	require.NoError(t, err)

	out := buf.String()

	assert.Contains(t, out, "Evidence: ~/.ssh has mode 0755")
	assert.Contains(t, out, "Fix: Run: chmod 700 ~/.ssh")
}

func TestPrettyReporter_SkipMarker(t *testing.T) {
	var buf bytes.Buffer

	r := reporter.NewPrettyReporter()
	err := r.Report(&buf, sampleScanResult())
	require.NoError(t, err)

	out := buf.String()

	// Skip uses a dash.
	assert.Contains(t, stripANSI(out), "- GIT001")
}

func TestPrettyReporter_ErrorMarker(t *testing.T) {
	var buf bytes.Buffer

	r := reporter.NewPrettyReporter()
	err := r.Report(&buf, sampleScanResult())
	require.NoError(t, err)

	out := buf.String()

	assert.Contains(t, stripANSI(out), "! GIT002 [high]")
}

func TestPrettyReporter_GroupsByCategory(t *testing.T) {
	var buf bytes.Buffer

	r := reporter.NewPrettyReporter()
	err := r.Report(&buf, sampleScanResult())
	require.NoError(t, err)

	out := buf.String()

	// Both category headers should appear.
	assert.Contains(t, out, "SSH")
	assert.Contains(t, out, "Git")
}

func TestPrettyReporter_SeveritySortWithinCategory(t *testing.T) {
	var buf bytes.Buffer

	r := reporter.NewPrettyReporter()
	err := r.Report(&buf, sampleScanResult())
	require.NoError(t, err)

	out := buf.String()

	// In the Git category, GIT002 (high) should appear before GIT001 (warn).
	git002Pos := strings.Index(out, "GIT002")
	git001Pos := strings.Index(out, "GIT001")
	assert.Greater(t, git001Pos, git002Pos, "GIT002 (high) should appear before GIT001 (warn)")
}

func TestPrettyReporter_BulletsMultipleEvidence(t *testing.T) {
	results := []schema.CheckResult{
		{
			RuleID:      "ENV001",
			Title:       "Sensitive credentials found in environment variables",
			Severity:    schema.SeverityHigh,
			Category:    "Env",
			Status:      schema.StatusFail,
			Evidence:    "Suspicious env var: GITHUB_TOKEN (matched pattern: exact:GITHUB_TOKEN); Suspicious env var: JIRA_API_TOKEN (matched pattern: *_TOKEN)",
			Remediation: "Use a secrets manager",
		},
	}
	sr := schema.NewScanResult(results)

	var buf bytes.Buffer

	r := reporter.NewPrettyReporter()
	err := r.Report(&buf, &sr)
	require.NoError(t, err)

	out := stripANSI(buf.String())

	// Header line on its own (no inline value).
	assert.Contains(t, out, "      Evidence:\n")
	// Each finding rendered as its own bullet, indented 8 spaces.
	assert.Contains(t, out, "        - Suspicious env var: GITHUB_TOKEN (matched pattern: exact:GITHUB_TOKEN)\n")
	assert.Contains(t, out, "        - Suspicious env var: JIRA_API_TOKEN (matched pattern: *_TOKEN)\n")
	// The original joined wall-of-text must NOT appear on a single line.
	assert.NotContains(t, out, "exact:GITHUB_TOKEN); Suspicious env var: JIRA_API_TOKEN")
}

func TestPrettyReporter_InlineSingleEvidence(t *testing.T) {
	var buf bytes.Buffer

	r := reporter.NewPrettyReporter()
	err := r.Report(&buf, sampleScanResult())
	require.NoError(t, err)

	out := stripANSI(buf.String())

	// Single-part evidence stays inline.
	assert.Contains(t, out, "      Evidence: ~/.ssh has mode 0755\n")
	// Must not be promoted to a bullet.
	assert.NotContains(t, out, "        - ~/.ssh has mode 0755")
}

func TestPrettyReporter_EmptyEvidence(t *testing.T) {
	results := []schema.CheckResult{
		{
			RuleID:   "ENV999",
			Title:    "Test rule",
			Severity: schema.SeverityHigh,
			Category: "Env",
			Status:   schema.StatusFail,
			Evidence: "",
		},
	}
	sr := schema.NewScanResult(results)

	var buf bytes.Buffer

	r := reporter.NewPrettyReporter()
	err := r.Report(&buf, &sr)
	require.NoError(t, err)

	out := stripANSI(buf.String())

	// Rule still appears.
	assert.Contains(t, out, "ENV999")
	// No Evidence line at all.
	assert.NotContains(t, out, "Evidence")
}

func TestPrettyReporter_TrailingSeparator(t *testing.T) {
	results := []schema.CheckResult{
		{
			RuleID:   "ENV998",
			Title:    "Test rule",
			Severity: schema.SeverityHigh,
			Category: "Env",
			Status:   schema.StatusFail,
			Evidence: "first finding; second finding; ",
		},
	}
	sr := schema.NewScanResult(results)

	var buf bytes.Buffer

	r := reporter.NewPrettyReporter()
	err := r.Report(&buf, &sr)
	require.NoError(t, err)

	out := stripANSI(buf.String())

	assert.Contains(t, out, "        - first finding\n")
	assert.Contains(t, out, "        - second finding\n")
	// Trailing empty part must not produce an empty bullet line.
	assert.NotContains(t, out, "        - \n")
}
