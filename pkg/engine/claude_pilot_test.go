package engine_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/foomo/obacht/internal/cli"
	"github.com/foomo/obacht/pkg/engine"
	"github.com/foomo/obacht/pkg/schema"
	"github.com/foomo/obacht/rules"
)

func TestPilot_ClaudeCategory_LoadsAndEvaluates(t *testing.T) {
	files, err := cli.LoadEmbeddedRuleFiles()
	require.NoError(t, err)

	var claudeFile *schema.RulesFile

	for i := range files {
		for _, r := range files[i].Rules {
			if strings.HasPrefix(r.ID, "CLD") {
				claudeFile = &files[i]
				break
			}
		}

		if claudeFile != nil {
			break
		}
	}

	require.NotNil(t, claudeFile, "claude category not found")
	require.Len(t, claudeFile.Rules, 41)

	for _, r := range claudeFile.Rules {
		assert.NotEmpty(t, r.Input, "rule %s has no input", r.ID)
		assert.NotEmpty(t, r.Policy, "rule %s has no policy", r.ID)
		assert.Equal(t, "claude", r.Category, "rule %s has wrong category", r.ID)
	}

	result, err := engine.EvaluateWithFS(t.Context(), rules.Embedded, []schema.RulesFile{*claudeFile})
	require.NoError(t, err)
	require.Len(t, result.Results, 41)
}
