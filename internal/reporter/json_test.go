package reporter_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/franklinkim/bouncer/internal/reporter"
	"github.com/franklinkim/bouncer/pkg/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONReporter_Report(t *testing.T) {
	result := schema.NewScanResult([]schema.CheckResult{
		{
			RuleID:   "R001",
			Title:    "Test rule",
			Severity: schema.SeverityHigh,
			Category: "security",
			Status:   schema.StatusFail,
			Evidence: "found issue",
		},
		{
			RuleID:   "R002",
			Title:    "Another rule",
			Severity: schema.SeverityInfo,
			Category: "best-practice",
			Status:   schema.StatusPass,
		},
	})

	var buf bytes.Buffer

	r := &reporter.JSONReporter{}
	err := r.Report(&buf, &result)
	require.NoError(t, err)

	// Verify it is valid JSON
	var decoded schema.ScanResult

	err = json.Unmarshal(buf.Bytes(), &decoded)
	require.NoError(t, err, "output must be valid JSON")

	// Verify expected fields
	assert.Equal(t, "1.0", decoded.SchemaVersion)
	assert.Len(t, decoded.Results, 2)
	assert.Equal(t, 2, decoded.Summary.Total)
	assert.Equal(t, 1, decoded.Summary.Passed)
	assert.Equal(t, 1, decoded.Summary.Failed)

	// Verify first result fields
	assert.Equal(t, "R001", decoded.Results[0].RuleID)
	assert.Equal(t, schema.SeverityHigh, decoded.Results[0].Severity)
	assert.Equal(t, schema.StatusFail, decoded.Results[0].Status)
	assert.Equal(t, "found issue", decoded.Results[0].Evidence)

	// Verify second result fields
	assert.Equal(t, "R002", decoded.Results[1].RuleID)
	assert.Equal(t, schema.StatusPass, decoded.Results[1].Status)
	assert.Empty(t, decoded.Results[1].Evidence, "omitempty should exclude empty evidence")
}

func TestJSONReporter_OutputIsIndented(t *testing.T) {
	result := schema.NewScanResult([]schema.CheckResult{
		{
			RuleID:   "R001",
			Title:    "Test rule",
			Severity: schema.SeverityWarn,
			Category: "test",
			Status:   schema.StatusPass,
		},
	})

	var buf bytes.Buffer

	r := &reporter.JSONReporter{}
	err := r.Report(&buf, &result)
	require.NoError(t, err)

	// Verify output uses 2-space indentation
	assert.Contains(t, buf.String(), "  \"schema_version\"")
}
