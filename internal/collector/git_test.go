package collector_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/franklinkim/bouncer/internal/collector"
	"github.com/franklinkim/bouncer/pkg/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseGitConfig_FullSigning(t *testing.T) {
	input := `user.name=John Doe
user.email=john@example.com
credential.helper=store
commit.gpgsign=true
gpg.format=ssh
`

	var git schema.GitFacts
	collector.ExportParseGitConfig(input, &git)

	assert.Equal(t, "store", git.CredentialHelper)
	assert.True(t, git.SigningEnabled)
	assert.Equal(t, "ssh", git.SigningFormat)
}

func TestParseGitConfig_NoSigning(t *testing.T) {
	input := `user.name=Jane Doe
user.email=jane@example.com
`

	var git schema.GitFacts
	collector.ExportParseGitConfig(input, &git)

	assert.Empty(t, git.CredentialHelper)
	assert.False(t, git.SigningEnabled)
	assert.Empty(t, git.SigningFormat)
}

func TestGitCollector_Collect(t *testing.T) {
	runner := collector.NewFakeRunner(map[string]collector.FakeResult{
		"git [--version]": {
			Output: []byte("git version 2.45.0\n"),
		},
		"git [config --global --list]": {
			Output: []byte("credential.helper=store\ncommit.gpgsign=true\ngpg.format=ssh\n"),
		},
	})

	c := &collector.GitCollector{Runner: runner}
	facts := schema.NewFacts()

	// Note: this test only works when git is actually on PATH.
	// If git is not installed, the collector returns StatusSkipped
	// before it even calls the runner.
	result := c.Collect(context.Background(), &facts)

	if result.Status == collector.StatusSkipped {
		t.Skip("git not found on PATH, skipping integration-style test")
	}

	require.NoError(t, result.Error)
	assert.Equal(t, collector.StatusOK, result.Status)
	assert.Equal(t, "git", result.Name)
	assert.True(t, facts.Git.Installed)
	assert.Equal(t, "git version 2.45.0", facts.Git.Version)
	assert.Equal(t, "store", facts.Git.CredentialHelper)
	assert.True(t, facts.Git.SigningEnabled)
	assert.Equal(t, "ssh", facts.Git.SigningFormat)
}

func TestGitCollector_VersionError(t *testing.T) {
	runner := collector.NewFakeRunner(map[string]collector.FakeResult{
		"git [--version]": {
			Err: fmt.Errorf("command failed"),
		},
	})

	c := &collector.GitCollector{Runner: runner}
	facts := schema.NewFacts()

	result := c.Collect(context.Background(), &facts)

	if result.Status == collector.StatusSkipped {
		t.Skip("git not found on PATH, skipping test")
	}

	assert.Equal(t, collector.StatusError, result.Status)
	assert.Error(t, result.Error)
}
