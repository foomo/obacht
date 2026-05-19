//go:build safe

package cli_test

import (
	"os"
	"path/filepath"
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
