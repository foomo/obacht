package collector

import (
	"context"
	"testing"

	"github.com/franklinkim/bouncer/pkg/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToolsCollector_Collect(t *testing.T) {
	collector := NewToolsCollector()
	facts := schema.NewFacts()

	result := collector.Collect(context.Background(), &facts)

	require.NoError(t, result.Error)
	assert.Equal(t, StatusOK, result.Status)
	assert.Equal(t, "tools", result.Name)

	// We should have entries for all security tools.
	assert.Len(t, facts.Tools.Tools, len(securityTools))

	// git should be installed in any dev environment.
	var gitTool *schema.ToolInfo
	for i := range facts.Tools.Tools {
		if facts.Tools.Tools[i].Name == "git" {
			gitTool = &facts.Tools.Tools[i]
			break
		}
	}
	require.NotNil(t, gitTool, "git should be present in tool results")
	assert.True(t, gitTool.Installed, "git should be installed")
	assert.NotEmpty(t, gitTool.Path, "git path should be set")
	// Version may be empty in sandboxed environments where exec is restricted.
	// When it is populated, it should contain "git".
	if gitTool.Version != "" {
		assert.Contains(t, gitTool.Version, "git")
	}
}

func TestToolsCollector_Structure(t *testing.T) {
	runner := &fakeRunner{
		results: map[string]fakeResult{
			"git [--version]": {output: []byte("git version 2.45.0\n")},
			"opa [version]":   {output: []byte("Version: 0.62.0\n")},
			"gpg [--version]": {output: []byte("gpg (GnuPG) 2.4.5\nmore lines\n")},
			"ssh-agent [-V]":  {output: []byte("OpenSSH_9.7p1\n")},
		},
	}

	collector := &ToolsCollector{Runner: runner}
	facts := schema.NewFacts()

	result := collector.Collect(context.Background(), &facts)

	assert.Equal(t, StatusOK, result.Status)

	// All tools should be present.
	assert.Len(t, facts.Tools.Tools, len(securityTools))

	// Each tool that is on PATH should have its info populated.
	for _, tool := range facts.Tools.Tools {
		assert.NotEmpty(t, tool.Name, "every tool should have a name")
		if tool.Installed {
			assert.NotEmpty(t, tool.Path, "installed tools should have a path")
		}
	}
}
