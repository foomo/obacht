package runner_test

import (
	"context"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/foomo/obacht/internal/runner"
)

func TestParseEnvelope_OK(t *testing.T) {
	raw := `{"rule_id":"SSH001","status":"ok","data":{"foo":"bar"}}`

	env, err := runner.ParseEnvelope([]byte(raw), "SSH001")
	require.NoError(t, err)
	assert.Equal(t, "SSH001", env.RuleID)
	assert.Equal(t, "ok", env.Status)

	data, ok := env.Data.(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "bar", data["foo"])
}

func TestParseEnvelope_Skip(t *testing.T) {
	raw := `{"rule_id":"SSH001","status":"skip","skip_reason":"no .ssh directory"}`

	env, err := runner.ParseEnvelope([]byte(raw), "SSH001")
	require.NoError(t, err)
	assert.Equal(t, "skip", env.Status)
	assert.Equal(t, "no .ssh directory", env.SkipReason)
}

func TestParseEnvelope_RuleIDMismatch_Errors(t *testing.T) {
	raw := `{"rule_id":"WRONG","status":"ok","data":{}}`

	_, err := runner.ParseEnvelope([]byte(raw), "SSH001")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "rule_id mismatch")
}

func TestParseEnvelope_MissingDataOnOK_Errors(t *testing.T) {
	raw := `{"rule_id":"SSH001","status":"ok"}`

	_, err := runner.ParseEnvelope([]byte(raw), "SSH001")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "data")
}

func TestParseEnvelope_MissingReasonOnSkip_Errors(t *testing.T) {
	raw := `{"rule_id":"SSH001","status":"skip"}`

	_, err := runner.ParseEnvelope([]byte(raw), "SSH001")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "skip_reason")
}

func TestParseEnvelope_LegacyBareJSON(t *testing.T) {
	// Old-style script output: no envelope, just data. Should be treated
	// as status=ok with the whole payload as `data` to preserve back-compat
	// during migration.
	raw := `{"directory_mode":"0755"}`

	env, err := runner.ParseEnvelope([]byte(raw), "SSH002")
	require.NoError(t, err)
	assert.Equal(t, "SSH002", env.RuleID)
	assert.Equal(t, "ok", env.Status)

	data, ok := env.Data.(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "0755", data["directory_mode"])
}

func TestParseEnvelope_LegacyBareJSON_NoExpectedID(t *testing.T) {
	// When expectedID is empty (e.g. shared category script), the bare-JSON
	// fallback still works and RuleID stays empty.
	raw := `{"directory_mode":"0755"}`

	env, err := runner.ParseEnvelope([]byte(raw), "")
	require.NoError(t, err)
	assert.Empty(t, env.RuleID)
	assert.Equal(t, "ok", env.Status)
}

func TestParseEnvelope_InvalidJSON_Errors(t *testing.T) {
	raw := `not json`

	_, err := runner.ParseEnvelope([]byte(raw), "SSH001")
	require.Error(t, err)
}

func TestParseEnvelope_NonStringStatus_TreatedAsLegacy(t *testing.T) {
	// A legacy script that emits a numeric or null status field is treated
	// as legacy bare JSON, not as a malformed envelope.
	raw := `{"status":1,"foo":"bar"}`

	env, err := runner.ParseEnvelope([]byte(raw), "X001")
	require.NoError(t, err)
	assert.Equal(t, "X001", env.RuleID)
	assert.Equal(t, "ok", env.Status)

	data, ok := env.Data.(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "bar", data["foo"])
}

func TestParseEnvelope_UnknownStringStatus_TreatedAsLegacy(t *testing.T) {
	// A legacy script that emits {"status":"online", ...} (e.g. a daemon
	// state collector) should NOT be misclassified as a malformed envelope.
	raw := `{"status":"online","port":22}`

	env, err := runner.ParseEnvelope([]byte(raw), "X001")
	require.NoError(t, err)
	assert.Equal(t, "ok", env.Status)
}

func TestParseEnvelope_EmptyBody_Errors(t *testing.T) {
	_, err := runner.ParseEnvelope([]byte(""), "X001")
	require.Error(t, err)
}

func TestParseEnvelope_TopLevelArray_Errors(t *testing.T) {
	raw := `[1,2,3]`

	_, err := runner.ParseEnvelope([]byte(raw), "X001")
	require.Error(t, err)
}

func TestParseEnvelope_OK_WithExplicitNullData_Errors(t *testing.T) {
	raw := `{"rule_id":"X001","status":"ok","data":null}`

	_, err := runner.ParseEnvelope([]byte(raw), "X001")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "data")
}

func TestRunInput_EnvelopeSkip(t *testing.T) {
	script := `#!/bin/sh
printf '{"rule_id":"SSH001","status":"skip","skip_reason":"no .ssh dir"}'
`

	res := runner.RunInputForRule(context.Background(), fstest.MapFS{}, "inputs/ssh/SSH001.sh", script, "SSH001")
	require.Equal(t, runner.StatusSkipped, res.Status)
	assert.Equal(t, "no .ssh dir", res.SkipReason)
}

func TestRunInput_EnvelopeOK(t *testing.T) {
	script := `#!/bin/sh
printf '{"rule_id":"SSH001","status":"ok","data":{"foo":"bar"}}'
`

	res := runner.RunInputForRule(context.Background(), fstest.MapFS{}, "inputs/ssh/SSH001.sh", script, "SSH001")
	require.Equal(t, runner.StatusOK, res.Status)

	data, ok := res.Data.(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "bar", data["foo"])
}

func TestRunInput_LegacyBareJSON(t *testing.T) {
	script := `printf '{"directory_mode":"0755"}'`

	res := runner.RunInputForRule(context.Background(), fstest.MapFS{}, "inputs/ssh/SSH001.sh", script, "")
	require.Equal(t, runner.StatusOK, res.Status)

	data, ok := res.Data.(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "0755", data["directory_mode"])
}

func TestRunInput_PreprocessIncludeFromFS(t *testing.T) {
	fsys := fstest.MapFS{
		"inputs/_lib/json.sh": &fstest.MapFile{Data: []byte(`
emit_ok() {
  printf '{"rule_id":"%s","status":"ok","data":%s}' "$1" "$2"
}
`)},
	}
	script := `#!/bin/sh
# include: _lib/json.sh
emit_ok SSH001 '{"foo":"bar"}'
`

	res := runner.RunInputForRule(context.Background(), fsys, "inputs/ssh/SSH001.sh", script, "SSH001")
	require.Equal(t, runner.StatusOK, res.Status)

	data, ok := res.Data.(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "bar", data["foo"])
}
