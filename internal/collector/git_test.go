package collector

import (
	"context"
	"fmt"
	"testing"

	"github.com/franklinkim/bouncer/pkg/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeRunner is a test double for CommandRunner.
type fakeRunner struct {
	// results maps "name arg1 arg2 ..." to the output or error.
	results map[string]fakeResult
}

type fakeResult struct {
	output []byte
	err    error
}

func (f *fakeRunner) Run(_ context.Context, name string, args ...string) ([]byte, error) {
	key := name
	if len(args) > 0 {
		key += " " + fmt.Sprintf("%s", args)
	}
	if r, ok := f.results[key]; ok {
		return r.output, r.err
	}
	return nil, fmt.Errorf("unexpected command: %s", key)
}

func TestParseGitConfig_FullSigning(t *testing.T) {
	input := `user.name=John Doe
user.email=john@example.com
credential.helper=store
commit.gpgsign=true
gpg.format=ssh
`
	var git schema.GitFacts
	parseGitConfig(input, &git)

	assert.Equal(t, "store", git.CredentialHelper)
	assert.True(t, git.SigningEnabled)
	assert.Equal(t, "ssh", git.SigningFormat)
}

func TestParseGitConfig_NoSigning(t *testing.T) {
	input := `user.name=Jane Doe
user.email=jane@example.com
`
	var git schema.GitFacts
	parseGitConfig(input, &git)

	assert.Equal(t, "", git.CredentialHelper)
	assert.False(t, git.SigningEnabled)
	assert.Equal(t, "", git.SigningFormat)
}

func TestGitCollector_Collect(t *testing.T) {
	runner := &fakeRunner{
		results: map[string]fakeResult{
			"git [--version]": {
				output: []byte("git version 2.45.0\n"),
			},
			"git [config --global --list]": {
				output: []byte("credential.helper=store\ncommit.gpgsign=true\ngpg.format=ssh\n"),
			},
		},
	}

	collector := &GitCollector{Runner: runner}
	facts := schema.NewFacts()

	// Note: this test only works when git is actually on PATH.
	// If git is not installed, the collector returns StatusSkipped
	// before it even calls the runner.
	result := collector.Collect(context.Background(), &facts)

	if result.Status == StatusSkipped {
		t.Skip("git not found on PATH, skipping integration-style test")
	}

	require.NoError(t, result.Error)
	assert.Equal(t, StatusOK, result.Status)
	assert.Equal(t, "git", result.Name)
	assert.True(t, facts.Git.Installed)
	assert.Equal(t, "git version 2.45.0", facts.Git.Version)
	assert.Equal(t, "store", facts.Git.CredentialHelper)
	assert.True(t, facts.Git.SigningEnabled)
	assert.Equal(t, "ssh", facts.Git.SigningFormat)
}

func TestGitCollector_VersionError(t *testing.T) {
	runner := &fakeRunner{
		results: map[string]fakeResult{
			"git [--version]": {
				err: fmt.Errorf("command failed"),
			},
		},
	}

	collector := &GitCollector{Runner: runner}
	facts := schema.NewFacts()

	result := collector.Collect(context.Background(), &facts)

	if result.Status == StatusSkipped {
		t.Skip("git not found on PATH, skipping test")
	}

	assert.Equal(t, StatusError, result.Status)
	assert.Error(t, result.Error)
}
