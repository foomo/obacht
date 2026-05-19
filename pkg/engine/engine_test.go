package engine_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/foomo/obacht/pkg/engine"
	"github.com/foomo/obacht/pkg/schema"
)

var testPolicy = `package obacht.test

import rego.v1

findings contains f if {
    input.directory_mode != "0700"
    f := {"rule_id": "SSH002", "evidence": sprintf("~/.ssh has mode %s", [input.directory_mode])}
}
`

// testPolicyAutoPrefix has no package/import — the engine should auto-prefix them.
var testPolicyAutoPrefix = `findings contains f if {
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

func TestEvaluate_EnvelopeSkip_PropagatesReason(t *testing.T) {
	ruleFiles := []schema.RulesFile{
		{
			Input:  `printf '{"rule_id":"SSH002","status":"skip","skip_reason":"no .ssh dir"}'`,
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
	assert.Equal(t, "no .ssh dir", cr.Evidence)
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

func TestEvaluate_AutoPrefix(t *testing.T) {
	ruleFiles := []schema.RulesFile{
		{
			Input: `printf '{"directory_mode": "0755"}'`,
			Rules: []schema.Rule{
				{
					ID:       "SSH002",
					Title:    "SSH directory permissions",
					Severity: schema.SeverityHigh,
					Category: "ssh",
					Policy:   testPolicyAutoPrefix,
				},
			},
		},
	}

	result, err := engine.Evaluate(t.Context(), ruleFiles)
	require.NoError(t, err)
	require.Len(t, result.Results, 1)

	cr := result.Results[0]
	assert.Equal(t, "SSH002", cr.RuleID)
	assert.Equal(t, schema.StatusFail, cr.Status)
	assert.Equal(t, "~/.ssh has mode 0755", cr.Evidence)
}

func TestEvaluate_AutoPrefix_FileLevelPolicy(t *testing.T) {
	ruleFiles := []schema.RulesFile{
		{
			Input:  `printf '{"directory_mode": "0755"}'`,
			Policy: testPolicyAutoPrefix,
			Rules: []schema.Rule{
				{
					ID:       "SSH002",
					Title:    "SSH directory permissions",
					Severity: schema.SeverityHigh,
					Category: "ssh",
				},
			},
		},
	}

	result, err := engine.Evaluate(t.Context(), ruleFiles)
	require.NoError(t, err)
	require.Len(t, result.Results, 1)

	cr := result.Results[0]
	assert.Equal(t, "SSH002", cr.RuleID)
	assert.Equal(t, schema.StatusFail, cr.Status)
	assert.Equal(t, "~/.ssh has mode 0755", cr.Evidence)
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

var testPolicySkip = `package obacht.test

import rego.v1

skips contains s if {
    input.disk_connected == false
    s := {"rule_id": "SSH002", "evidence": "disk not connected"}
}
`

func TestEvaluate_RegoSkip_MapsToStatusSkip(t *testing.T) {
	ruleFiles := []schema.RulesFile{
		{
			Input:  `printf '{"disk_connected": false}'`,
			Policy: testPolicySkip,
			Rules:  testRules,
		},
	}

	result, err := engine.Evaluate(t.Context(), ruleFiles)
	require.NoError(t, err)
	require.Len(t, result.Results, 1)

	cr := result.Results[0]
	assert.Equal(t, "SSH002", cr.RuleID)
	assert.Equal(t, schema.StatusSkip, cr.Status)
	assert.Equal(t, "disk not connected", cr.Evidence)
}

var testPolicyFailAndSkip = `package obacht.test

import rego.v1

findings contains f if {
    input.directory_mode != "0700"
    f := {"rule_id": "SSH002", "evidence": sprintf("~/.ssh has mode %s", [input.directory_mode])}
}

skips contains s if {
    input.directory_mode != "0700"
    s := {"rule_id": "SSH002", "evidence": "ignored when fail also fires"}
}
`

func TestEvaluate_RegoFailWinsOverSkip(t *testing.T) {
	ruleFiles := []schema.RulesFile{
		{
			Input:  `printf '{"directory_mode": "0755"}'`,
			Policy: testPolicyFailAndSkip,
			Rules:  testRules,
		},
	}

	result, err := engine.Evaluate(t.Context(), ruleFiles)
	require.NoError(t, err)
	require.Len(t, result.Results, 1)

	cr := result.Results[0]
	assert.Equal(t, schema.StatusFail, cr.Status)
	assert.Equal(t, "~/.ssh has mode 0755", cr.Evidence)
}

var testPolicyMultiSkip = `package obacht.test

import rego.v1

skips contains s if {
    input.foo == true
    s := {"rule_id": "SSH002", "evidence": "alpha"}
}

skips contains s if {
    input.bar == true
    s := {"rule_id": "SSH002", "evidence": "bravo"}
}
`

func TestEvaluate_MultipleSkips_AggregateEvidence(t *testing.T) {
	ruleFiles := []schema.RulesFile{
		{
			Input:  `printf '{"foo": true, "bar": true}'`,
			Policy: testPolicyMultiSkip,
			Rules:  testRules,
		},
	}

	result, err := engine.Evaluate(t.Context(), ruleFiles)
	require.NoError(t, err)
	require.Len(t, result.Results, 1)

	cr := result.Results[0]
	assert.Equal(t, schema.StatusSkip, cr.Status)
	assert.Equal(t, "alpha; bravo", cr.Evidence)
}

// OS031 — skip when password_lock_delay_seconds == -1 (cannot determine).
var os031Policy = `package obacht.os

import rego.v1

findings contains f if {
    input.os == "darwin"
    not input.password_required_on_lock
    input.password_lock_delay_seconds > 0
    f := {
        "rule_id": "OS031",
        "evidence": sprintf("Password not required immediately after screen lock (delay: %ds)", [input.password_lock_delay_seconds]),
    }
}

findings contains f if {
    input.os == "darwin"
    input.password_lock_delay_seconds == -2
    f := {
        "rule_id": "OS031",
        "evidence": "Screen lock is disabled — no password required after screensaver",
    }
}

skips contains s if {
    input.os == "darwin"
    input.password_lock_delay_seconds == -1
    s := {
        "rule_id": "OS031",
        "evidence": "Could not determine screen-lock password delay (sysadminctl unavailable and legacy askForPasswordDelay unreadable)",
    }
}
`

var os031Rules = []schema.Rule{
	{
		ID:       "OS031",
		Title:    "Password not required immediately after screen lock",
		Severity: schema.SeverityHigh,
		Category: "os",
	},
}

func TestEvaluate_OS031_Skip_WhenDelayUnknown(t *testing.T) {
	ruleFiles := []schema.RulesFile{
		{
			Input:  `printf '{"os": "darwin", "password_required_on_lock": false, "password_lock_delay_seconds": -1}'`,
			Policy: os031Policy,
			Rules:  os031Rules,
		},
	}

	result, err := engine.Evaluate(t.Context(), ruleFiles)
	require.NoError(t, err)
	require.Len(t, result.Results, 1)

	cr := result.Results[0]
	assert.Equal(t, "OS031", cr.RuleID)
	assert.Equal(t, schema.StatusSkip, cr.Status)
	assert.Contains(t, cr.Evidence, "Could not determine screen-lock password delay")
}

func TestEvaluate_OS031_Fail_WhenScreenLockDisabled(t *testing.T) {
	ruleFiles := []schema.RulesFile{
		{
			Input:  `printf '{"os": "darwin", "password_required_on_lock": false, "password_lock_delay_seconds": -2}'`,
			Policy: os031Policy,
			Rules:  os031Rules,
		},
	}

	result, err := engine.Evaluate(t.Context(), ruleFiles)
	require.NoError(t, err)
	require.Len(t, result.Results, 1)

	cr := result.Results[0]
	assert.Equal(t, schema.StatusFail, cr.Status)
	assert.Contains(t, cr.Evidence, "Screen lock is disabled")
}

func TestEvaluate_OS031_Pass_WhenImmediate(t *testing.T) {
	ruleFiles := []schema.RulesFile{
		{
			Input:  `printf '{"os": "darwin", "password_required_on_lock": true, "password_lock_delay_seconds": 0}'`,
			Policy: os031Policy,
			Rules:  os031Rules,
		},
	}

	result, err := engine.Evaluate(t.Context(), ruleFiles)
	require.NoError(t, err)
	require.Len(t, result.Results, 1)
	assert.Equal(t, schema.StatusPass, result.Results[0].Status)
}

// SSH005 — skip when key.bits == 0 (ssh-keygen could not parse pub key).
var ssh005Policy = `package obacht.ssh

import rego.v1

findings contains f if {
    key := input.keys[_]
    key.algorithm == "DSA"
    f := {
        "rule_id": "SSH005",
        "evidence": sprintf("SSH key %s uses weak algorithm: DSA/%d bits", [key.path, key.bits]),
    }
}

findings contains f if {
    key := input.keys[_]
    key.algorithm == "RSA"
    key.bits > 0
    key.bits < 3072
    f := {
        "rule_id": "SSH005",
        "evidence": sprintf("SSH key %s uses weak algorithm: RSA/%d bits (minimum 3072)", [key.path, key.bits]),
    }
}

skips contains s if {
    key := input.keys[_]
    key.bits == 0
    s := {
        "rule_id": "SSH005",
        "evidence": sprintf("Could not read SSH key bits: %s (ssh-keygen failed)", [key.path]),
    }
}
`

var ssh005Rules = []schema.Rule{
	{
		ID:       "SSH005",
		Title:    "SSH key uses weak algorithm",
		Severity: schema.SeverityHigh,
		Category: "ssh",
	},
}

func TestEvaluate_SSH005_Skip_WhenBitsUnreadable(t *testing.T) {
	ruleFiles := []schema.RulesFile{
		{
			Input:  `printf '{"keys": [{"path": "/home/u/.ssh/id_rsa.pub", "mode": "0644", "type": "rsa", "bits": 0, "algorithm": ""}]}'`,
			Policy: ssh005Policy,
			Rules:  ssh005Rules,
		},
	}

	result, err := engine.Evaluate(t.Context(), ruleFiles)
	require.NoError(t, err)
	require.Len(t, result.Results, 1)

	cr := result.Results[0]
	assert.Equal(t, "SSH005", cr.RuleID)
	assert.Equal(t, schema.StatusSkip, cr.Status)
	assert.Contains(t, cr.Evidence, "Could not read SSH key bits")
	assert.Contains(t, cr.Evidence, "id_rsa.pub")
}

func TestEvaluate_SSH005_Pass_WhenStrongKey(t *testing.T) {
	ruleFiles := []schema.RulesFile{
		{
			Input:  `printf '{"keys": [{"path": "/home/u/.ssh/id_ed25519.pub", "mode": "0644", "type": "ed25519", "bits": 256, "algorithm": "ED25519"}]}'`,
			Policy: ssh005Policy,
			Rules:  ssh005Rules,
		},
	}

	result, err := engine.Evaluate(t.Context(), ruleFiles)
	require.NoError(t, err)
	require.Len(t, result.Results, 1)
	assert.Equal(t, schema.StatusPass, result.Results[0].Status)
}

func TestEvaluate_SSH005_Fail_WhenWeakRSA(t *testing.T) {
	ruleFiles := []schema.RulesFile{
		{
			Input:  `printf '{"keys": [{"path": "/home/u/.ssh/id_rsa.pub", "mode": "0644", "type": "rsa", "bits": 2048, "algorithm": "RSA"}]}'`,
			Policy: ssh005Policy,
			Rules:  ssh005Rules,
		},
	}

	result, err := engine.Evaluate(t.Context(), ruleFiles)
	require.NoError(t, err)
	require.Len(t, result.Results, 1)
	assert.Equal(t, schema.StatusFail, result.Results[0].Status)
}
