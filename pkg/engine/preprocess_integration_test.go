package engine_test

import (
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/foomo/obacht/pkg/engine"
	"github.com/foomo/obacht/pkg/schema"
)

func TestEvaluate_WithFS_ResolvesInclude(t *testing.T) {
	fsys := fstest.MapFS{
		"inputs/_lib/json.sh": &fstest.MapFile{Data: []byte(`emit_ok() {
  printf '{"rule_id":"%s","status":"ok","data":%s}' "$1" "$2"
}
`)},
	}

	script := `#!/bin/sh
# include: _lib/json.sh
emit_ok SSH002 '{"directory_mode":"0755"}'
`

	ruleFiles := []schema.RulesFile{
		{
			Rules: []schema.Rule{
				{
					ID:       "SSH002",
					Title:    "SSH directory permissions",
					Severity: schema.SeverityHigh,
					Category: "ssh",
					Input:    script,
					Policy:   testPolicy,
				},
			},
		},
	}

	result, err := engine.EvaluateWithFS(t.Context(), fsys, ruleFiles)
	require.NoError(t, err)
	require.Len(t, result.Results, 1)

	cr := result.Results[0]
	assert.Equal(t, "SSH002", cr.RuleID)
	assert.Equal(t, schema.StatusFail, cr.Status)
}
