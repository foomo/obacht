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

func TestPilot_PathCategory_LoadsAndEvaluates(t *testing.T) {
	files, err := cli.LoadEmbeddedRuleFiles()
	require.NoError(t, err)

	var pathFile *schema.RulesFile

	for i := range files {
		for _, r := range files[i].Rules {
			if strings.HasPrefix(r.ID, "PTH") {
				pathFile = &files[i]
				break
			}
		}

		if pathFile != nil {
			break
		}
	}

	require.NotNil(t, pathFile, "path category not found")
	require.Len(t, pathFile.Rules, 2)

	for _, r := range pathFile.Rules {
		assert.NotEmpty(t, r.Input, "rule %s has no input", r.ID)
		assert.NotEmpty(t, r.Policy, "rule %s has no policy", r.ID)
	}

	result, err := engine.EvaluateWithFS(t.Context(), rules.Embedded, []schema.RulesFile{*pathFile})
	require.NoError(t, err)
	require.Len(t, result.Results, 2)
}
