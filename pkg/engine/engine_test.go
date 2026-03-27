package engine_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franklinkim/bouncer/pkg/engine"
	"github.com/franklinkim/bouncer/pkg/schema"
)

var testPolicy = `package bouncer.test

import rego.v1

findings contains f if {
    input.directory_mode != "0700"
    f := {"rule_id": "SSH002", "evidence": sprintf("~/.ssh has mode %s", [input.directory_mode])}
}
`

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
	ruleFiles := []schema.RulesFile{
		{
			Input:  `printf '{"directory_mode": "0755"}'`,
			Policy: testPolicy,
			Rules:  testRules,
		},
	}

	result, err := engine.Evaluate(t.Context(), ruleFiles)
	require.NoError(t, err)
	require.Len(t, result.Results, 1)

	cr := result.Results[0]
	assert.Equal(t, "SSH002", cr.RuleID)
	assert.Equal(t, schema.StatusFail, cr.Status)
	assert.Equal(t, "~/.ssh has mode 0755", cr.Evidence)
	assert.Equal(t, schema.SeverityHigh, cr.Severity)
}

func TestEvaluate_Pass(t *testing.T) {
	ruleFiles := []schema.RulesFile{
		{
			Input:  `printf '{"directory_mode": "0700"}'`,
			Policy: testPolicy,
			Rules:  testRules,
		},
	}

	result, err := engine.Evaluate(t.Context(), ruleFiles)
	require.NoError(t, err)
	require.Len(t, result.Results, 1)

	cr := result.Results[0]
	assert.Equal(t, "SSH002", cr.RuleID)
	assert.Equal(t, schema.StatusPass, cr.Status)
	assert.Empty(t, cr.Evidence)
}

func TestEvaluate_NoInput(t *testing.T) {
	ruleFiles := []schema.RulesFile{
		{
			Policy: testPolicy,
			Rules:  testRules,
		},
	}

	result, err := engine.Evaluate(t.Context(), ruleFiles)
	require.NoError(t, err)
	require.Len(t, result.Results, 1)

	cr := result.Results[0]
	assert.Equal(t, "SSH002", cr.RuleID)
	assert.Equal(t, schema.StatusSkip, cr.Status)
	assert.Empty(t, cr.Evidence)
}

func TestEvaluate_InputError(t *testing.T) {
	ruleFiles := []schema.RulesFile{
		{
			Input:  `exit 1`,
			Policy: testPolicy,
			Rules:  testRules,
		},
	}

	result, err := engine.Evaluate(t.Context(), ruleFiles)
	require.NoError(t, err)
	require.Len(t, result.Results, 1)

	cr := result.Results[0]
	assert.Equal(t, "SSH002", cr.RuleID)
	assert.Equal(t, schema.StatusError, cr.Status)
}

func TestEvaluate_RuleLevelOverride(t *testing.T) {
	ruleFiles := []schema.RulesFile{
		{
			Input:  `printf '{"directory_mode": "0700"}'`,
			Policy: testPolicy,
			Rules: []schema.Rule{
				{
					ID:       "SSH002",
					Title:    "SSH directory permissions",
					Severity: schema.SeverityHigh,
					Category: "ssh",
					// Rule-level input overrides file-level.
					Input: `printf '{"directory_mode": "0755"}'`,
				},
			},
		},
	}

	result, err := engine.Evaluate(t.Context(), ruleFiles)
	require.NoError(t, err)
	require.Len(t, result.Results, 1)

	cr := result.Results[0]
	assert.Equal(t, schema.StatusFail, cr.Status)
}
