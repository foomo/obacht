package collector_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/franklinkim/bouncer/internal/collector"
	"github.com/franklinkim/bouncer/pkg/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShellCollector_Zsh(t *testing.T) {
	home := t.TempDir()

	// Create a history file with known permissions.
	histPath := filepath.Join(home, ".zsh_history")
	writeFile(t, histPath, 0600)

	t.Setenv("SHELL", "/bin/zsh")
	t.Setenv("HISTCONTROL", "ignorespace")

	c := &collector.ShellCollector{HomeDir: home}
	facts := schema.NewFacts()
	result := c.Collect(context.Background(), &facts)

	assert.Equal(t, collector.StatusOK, result.Status)
	assert.Equal(t, "/bin/zsh", facts.Shell.Shell)
	assert.Equal(t, histPath, facts.Shell.HistoryFile)
	assert.Equal(t, "0600", facts.Shell.HistoryFileMode)
	assert.Equal(t, "ignorespace", facts.Shell.HistControl)
}

func TestShellCollector_Bash(t *testing.T) {
	home := t.TempDir()

	histPath := filepath.Join(home, ".bash_history")
	writeFile(t, histPath, 0644)

	t.Setenv("SHELL", "/bin/bash")
	t.Setenv("HISTCONTROL", "")

	c := &collector.ShellCollector{HomeDir: home}
	facts := schema.NewFacts()
	result := c.Collect(context.Background(), &facts)

	assert.Equal(t, collector.StatusOK, result.Status)
	assert.Equal(t, "/bin/bash", facts.Shell.Shell)
	assert.Equal(t, histPath, facts.Shell.HistoryFile)
	assert.Equal(t, "0644", facts.Shell.HistoryFileMode)
}

func TestShellCollector_Fish(t *testing.T) {
	home := t.TempDir()

	fishDir := filepath.Join(home, ".local", "share", "fish")
	require.NoError(t, os.MkdirAll(fishDir, 0755))
	histPath := filepath.Join(fishDir, "fish_history")
	writeFile(t, histPath, 0600)

	t.Setenv("SHELL", "/usr/bin/fish")
	t.Setenv("HISTCONTROL", "")

	c := &collector.ShellCollector{HomeDir: home}
	facts := schema.NewFacts()
	result := c.Collect(context.Background(), &facts)

	assert.Equal(t, collector.StatusOK, result.Status)
	assert.Equal(t, histPath, facts.Shell.HistoryFile)
}

func TestShellCollector_NoHistoryFile(t *testing.T) {
	home := t.TempDir()

	t.Setenv("SHELL", "/bin/zsh")

	c := &collector.ShellCollector{HomeDir: home}
	facts := schema.NewFacts()
	result := c.Collect(context.Background(), &facts)

	assert.Equal(t, collector.StatusOK, result.Status)
	assert.Empty(t, facts.Shell.HistoryFileMode)
}

func TestHistoryFilePath(t *testing.T) {
	tests := []struct {
		shell    string
		expected string
	}{
		{"/bin/bash", ".bash_history"},
		{"/bin/zsh", ".zsh_history"},
		{"/usr/bin/fish", filepath.Join(".local", "share", "fish", "fish_history")},
		{"/bin/csh", ""},
	}
	for _, tt := range tests {
		t.Run(filepath.Base(tt.shell), func(t *testing.T) {
			result := collector.ExportHistoryFilePath(tt.shell, "/home/user")
			if tt.expected == "" {
				assert.Empty(t, result)
			} else {
				assert.Equal(t, filepath.Join("/home/user", tt.expected), result)
			}
		})
	}
}
