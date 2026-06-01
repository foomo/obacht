//go:build safe

package cli_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/foomo/obacht/internal/cli"
)

func TestLoadExternal_PerRuleLookup(t *testing.T) {
	dir := t.TempDir()

	require.NoError(t, os.MkdirAll(filepath.Join(dir, "policies"), 0o755))
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "inputs", "demo"), 0o755))
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "policy", "demo"), 0o755))

	require.NoError(t, os.WriteFile(filepath.Join(dir, "policies", "demo.yaml"), []byte(`
rules:
  - id: DEMO001
    title: demo rule
    severity: warn
    category: demo
    description: test
    remediation: none
`), 0o644))

	require.NoError(t, os.WriteFile(filepath.Join(dir, "inputs", "demo", "DEMO001.sh"),
		[]byte(`#!/bin/sh
printf '{"rule_id":"DEMO001","status":"ok","data":{"flag":true}}'
`), 0o644))

	require.NoError(t, os.WriteFile(filepath.Join(dir, "policy", "demo", "DEMO001.rego"),
		[]byte(`package obacht.demo
import rego.v1

findings contains f if {
  input.flag == true
  f := {"rule_id": "DEMO001", "evidence": "flag is true"}
}
`), 0o644))

	files, err := cli.LoadExternalRuleFiles(dir)
	require.NoError(t, err)
	require.Len(t, files, 1)
	require.Len(t, files[0].Rules, 1)

	rule := files[0].Rules[0]
	assert.NotEmpty(t, rule.Input)
	assert.NotEmpty(t, rule.Policy)
	assert.Contains(t, rule.Policy, "package obacht.demo")
}

func TestLoadExternal_SharedInputBundlesPolicy(t *testing.T) {
	dir := t.TempDir()

	require.NoError(t, os.MkdirAll(filepath.Join(dir, "policies"), 0o755))
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "inputs"), 0o755))
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "policy", "demo"), 0o755))

	require.NoError(t, os.WriteFile(filepath.Join(dir, "policies", "demo.yaml"), []byte(`
rules:
  - id: DEMO001
    title: a
    severity: warn
    category: demo
    description: x
    remediation: y
  - id: DEMO002
    title: b
    severity: warn
    category: demo
    description: x
    remediation: y
`), 0o644))

	require.NoError(t, os.WriteFile(filepath.Join(dir, "inputs", "demo.sh"),
		[]byte("#!/bin/sh\nprintf '{\"rule_id\":\"DEMO000\",\"status\":\"ok\",\"data\":{\"items\":[]}}'\n"), 0o644))

	require.NoError(t, os.WriteFile(filepath.Join(dir, "policy", "demo", "_helpers.rego"),
		[]byte("package obacht.demo\n\nimport rego.v1\n\n_ev(x) := sprintf(\"%v\", [x])\n"), 0o644))

	require.NoError(t, os.WriteFile(filepath.Join(dir, "policy", "demo", "DEMO001.rego"),
		[]byte("package obacht.demo\n\nimport rego.v1\n\nfindings contains f if {\n  some x in input.items\n  x == 1\n  f := {\"rule_id\": \"DEMO001\", \"evidence\": _ev(x)}\n}\n"), 0o644))

	require.NoError(t, os.WriteFile(filepath.Join(dir, "policy", "demo", "DEMO002.rego"),
		[]byte("package obacht.demo\n\nimport rego.v1\n\nfindings contains f if {\n  some x in input.items\n  x == 2\n  f := {\"rule_id\": \"DEMO002\", \"evidence\": _ev(x)}\n}\n"), 0o644))

	files, err := cli.LoadExternalRuleFiles(dir)
	require.NoError(t, err)
	require.Len(t, files, 1)

	rf := files[0]
	assert.NotEmpty(t, rf.Input, "shared input should populate file-level Input")
	assert.Contains(t, rf.Policy, "package obacht.demo")
	assert.Contains(t, rf.Policy, `"rule_id": "DEMO001"`)
	assert.Contains(t, rf.Policy, `"rule_id": "DEMO002"`)
	assert.Contains(t, rf.Policy, "_ev(x)")

	// Bundle mode must keep exactly one package declaration.
	assert.Equal(t, 1, strings.Count(rf.Policy, "package obacht.demo"))

	// Per-rule rule.Policy must stay empty so engine groups all rules together.
	for _, r := range rf.Rules {
		assert.Empty(t, r.Policy, "rule %s should not have per-rule policy when bundled", r.ID)
	}
}

func TestLoadExternal_InvalidRuleID_Rejected(t *testing.T) {
	dir := t.TempDir()

	require.NoError(t, os.MkdirAll(filepath.Join(dir, "policies"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "policies", "demo.yaml"), []byte(`
rules:
  - id: "../../etc/passwd"
    title: malicious
    severity: warn
    category: demo
    description: test
    remediation: none
`), 0o644))

	_, err := cli.LoadExternalRuleFiles(dir)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid rule ID")
}

func TestLoadExternal_LowercaseID_Rejected(t *testing.T) {
	dir := t.TempDir()

	require.NoError(t, os.MkdirAll(filepath.Join(dir, "policies"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "policies", "demo.yaml"), []byte(`
rules:
  - id: "ssh001"
    title: lowercase id
    severity: warn
    category: demo
    description: test
    remediation: none
`), 0o644))

	_, err := cli.LoadExternalRuleFiles(dir)
	require.Error(t, err)
}
