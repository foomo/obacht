//go:build safe

package cli_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franklinkim/bouncer/internal/cli"
	"github.com/franklinkim/bouncer/pkg/schema"
)

var testRuleFiles = []schema.RulesFile{
	{
		Input:  "echo ssh",
		Policy: "package bouncer.ssh",
		Rules: []schema.Rule{
			{ID: "SSH001", Category: "ssh", Severity: schema.SeverityHigh},
			{ID: "SSH002", Category: "ssh", Severity: schema.SeverityCritical},
		},
	},
	{
		Input:  "echo git",
		Policy: "package bouncer.git",
		Rules: []schema.Rule{
			{ID: "GIT001", Category: "git", Severity: schema.SeverityWarn},
			{ID: "GIT002", Category: "git", Severity: schema.SeverityInfo},
		},
	},
}

func TestParseRuleIDs(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  map[string]bool
	}{
		{"empty", "", nil},
		{"single", "GIT001", map[string]bool{"GIT001": true}},
		{"multiple", "GIT001,SSH002", map[string]bool{"GIT001": true, "SSH002": true}},
		{"whitespace", " GIT001 , SSH002 ", map[string]bool{"GIT001": true, "SSH002": true}},
		{"trailing comma", "GIT001,", map[string]bool{"GIT001": true}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cli.ParseRuleIDs(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCollectRuleIDs(t *testing.T) {
	got := cli.CollectRuleIDs(testRuleFiles)
	assert.Equal(t, map[string]bool{
		"SSH001": true,
		"SSH002": true,
		"GIT001": true,
		"GIT002": true,
	}, got)
}

func TestCollectRuleIDs_Empty(t *testing.T) {
	got := cli.CollectRuleIDs(nil)
	assert.Empty(t, got)
}

func TestValidateRuleIDs(t *testing.T) {
	known := map[string]bool{"SSH001": true, "GIT001": true}

	t.Run("all known", func(t *testing.T) {
		err := cli.ValidateRuleIDs(map[string]bool{"SSH001": true}, known)
		require.NoError(t, err)
	})

	t.Run("nil requested", func(t *testing.T) {
		err := cli.ValidateRuleIDs(nil, known)
		require.NoError(t, err)
	})

	t.Run("one unknown", func(t *testing.T) {
		err := cli.ValidateRuleIDs(map[string]bool{"FAKE001": true}, known)
		require.EqualError(t, err, "unknown rule IDs: FAKE001")
	})

	t.Run("multiple unknown sorted", func(t *testing.T) {
		err := cli.ValidateRuleIDs(map[string]bool{"ZZZ001": true, "AAA001": true}, known)
		require.EqualError(t, err, "unknown rule IDs: AAA001, ZZZ001")
	})
}

func TestFilterRuleFilesByID(t *testing.T) {
	t.Run("nil returns input unchanged", func(t *testing.T) {
		got := cli.FilterRuleFilesByID(testRuleFiles, nil)
		assert.Equal(t, testRuleFiles, got)
	})

	t.Run("single ID", func(t *testing.T) {
		got := cli.FilterRuleFilesByID(testRuleFiles, map[string]bool{"SSH001": true})
		require.Len(t, got, 1)
		assert.Equal(t, "echo ssh", got[0].Input)
		require.Len(t, got[0].Rules, 1)
		assert.Equal(t, "SSH001", got[0].Rules[0].ID)
	})

	t.Run("IDs across rule files", func(t *testing.T) {
		got := cli.FilterRuleFilesByID(testRuleFiles, map[string]bool{"SSH001": true, "GIT002": true})
		require.Len(t, got, 2)
		assert.Len(t, got[0].Rules, 1)
		assert.Equal(t, "SSH001", got[0].Rules[0].ID)
		assert.Len(t, got[1].Rules, 1)
		assert.Equal(t, "GIT002", got[1].Rules[0].ID)
	})

	t.Run("no match drops all", func(t *testing.T) {
		got := cli.FilterRuleFilesByID(testRuleFiles, map[string]bool{"NONE001": true})
		assert.Empty(t, got)
	})

	t.Run("preserves input and policy", func(t *testing.T) {
		got := cli.FilterRuleFilesByID(testRuleFiles, map[string]bool{"GIT001": true})
		require.Len(t, got, 1)
		assert.Equal(t, "echo git", got[0].Input)
		assert.Equal(t, "package bouncer.git", got[0].Policy)
	})
}

func TestExcludeRuleFilesByID(t *testing.T) {
	t.Run("nil returns input unchanged", func(t *testing.T) {
		got := cli.ExcludeRuleFilesByID(testRuleFiles, nil)
		assert.Equal(t, testRuleFiles, got)
	})

	t.Run("exclude one rule", func(t *testing.T) {
		got := cli.ExcludeRuleFilesByID(testRuleFiles, map[string]bool{"SSH001": true})
		require.Len(t, got, 2)
		require.Len(t, got[0].Rules, 1)
		assert.Equal(t, "SSH002", got[0].Rules[0].ID)
		assert.Len(t, got[1].Rules, 2)
	})

	t.Run("exclude all in a rule file drops it", func(t *testing.T) {
		got := cli.ExcludeRuleFilesByID(testRuleFiles, map[string]bool{"SSH001": true, "SSH002": true})
		require.Len(t, got, 1)
		assert.Equal(t, "echo git", got[0].Input)
	})

	t.Run("exclude all rules", func(t *testing.T) {
		got := cli.ExcludeRuleFilesByID(testRuleFiles, map[string]bool{
			"SSH001": true, "SSH002": true, "GIT001": true, "GIT002": true,
		})
		assert.Empty(t, got)
	})
}
