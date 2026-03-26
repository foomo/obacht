package reporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/franklinkim/bouncer/internal/reporter"
	"github.com/franklinkim/bouncer/pkg/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

	r := reporter.NewPrettyReporter()
	err := r.Report(&buf, sampleScanResult())
	require.NoError(t, err)

	out := buf.String()

	// Pass marker (checkmark) should be present.
	assert.Contains(t, out, "\u2713")
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

	r := reporter.NewPrettyReporter()
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
	assert.Contains(t, out, "- GIT001")
}

func TestPrettyReporter_ErrorMarker(t *testing.T) {
	var buf bytes.Buffer

	r := reporter.NewPrettyReporter()
	err := r.Report(&buf, sampleScanResult())
	require.NoError(t, err)

	out := buf.String()

	assert.Contains(t, out, "! GIT002")
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
