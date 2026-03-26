package engine_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franklinkim/bouncer/internal/collector"
	"github.com/franklinkim/bouncer/pkg/engine"
	"github.com/franklinkim/bouncer/pkg/schema"
)

var testPolicy = []byte(`package bouncer.test

import rego.v1

findings contains f if {
    input.ssh.directory_mode != "0700"
    f := {"rule_id": "SSH002", "evidence": sprintf("~/.ssh has mode %s", [input.ssh.directory_mode])}
}
`)

var testRules = []schema.Rule{
	{
		ID:          "SSH002",
		Title:       "SSH directory permissions",
		Severity:    schema.SeverityHigh,
		Category:    "ssh",
		Description: "SSH directory should have mode 0700",
		Remediation: "chmod 0700 ~/.ssh",
	},
}

func TestEvaluate_Fail(t *testing.T) {
	eng, err := engine.NewEngine([][]byte{testPolicy}, testRules)
	require.NoError(t, err)

	facts := schema.NewFacts()
	facts.SSH = schema.SSHFacts{
		DirectoryExists: true,
		DirectoryMode:   "0755",
	}

	collectorResults := []collector.Result{
		{Name: "ssh", Status: collector.StatusOK},
	}

	result, err := eng.Evaluate(t.Context(), &facts, collectorResults)
	require.NoError(t, err)
	require.Len(t, result.Results, 1)

	cr := result.Results[0]
	assert.Equal(t, "SSH002", cr.RuleID)
	assert.Equal(t, schema.StatusFail, cr.Status)
	assert.Equal(t, "~/.ssh has mode 0755", cr.Evidence)
	assert.Equal(t, schema.SeverityHigh, cr.Severity)
}

func TestEvaluate_Pass(t *testing.T) {
	eng, err := engine.NewEngine([][]byte{testPolicy}, testRules)
	require.NoError(t, err)

	facts := schema.NewFacts()
	facts.SSH = schema.SSHFacts{
		DirectoryExists: true,
		DirectoryMode:   "0700",
	}

	collectorResults := []collector.Result{
		{Name: "ssh", Status: collector.StatusOK},
	}

	result, err := eng.Evaluate(t.Context(), &facts, collectorResults)
	require.NoError(t, err)
	require.Len(t, result.Results, 1)

	cr := result.Results[0]
	assert.Equal(t, "SSH002", cr.RuleID)
	assert.Equal(t, schema.StatusPass, cr.Status)
	assert.Empty(t, cr.Evidence)
}

func TestEvaluate_CollectorSkipped(t *testing.T) {
	eng, err := engine.NewEngine([][]byte{testPolicy}, testRules)
	require.NoError(t, err)

	facts := schema.NewFacts()
	facts.SSH = schema.SSHFacts{
		DirectoryExists: false,
		DirectoryMode:   "0755",
	}

	collectorResults := []collector.Result{
		{Name: "ssh", Status: collector.StatusSkipped},
	}

	result, err := eng.Evaluate(t.Context(), &facts, collectorResults)
	require.NoError(t, err)
	require.Len(t, result.Results, 1)

	cr := result.Results[0]
	assert.Equal(t, "SSH002", cr.RuleID)
	assert.Equal(t, schema.StatusSkip, cr.Status)
	assert.Empty(t, cr.Evidence)
}
