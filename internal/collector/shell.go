package collector

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/franklinkim/bouncer/pkg/schema"
)

// ShellCollector gathers facts about the user's shell configuration.
type ShellCollector struct {
	// HomeDir overrides the user's home directory for testing.
	HomeDir string
}

// NewShellCollector returns a ShellCollector that uses the real home directory.
func NewShellCollector() *ShellCollector {
	return &ShellCollector{}
}

// Name returns the collector name.
func (c *ShellCollector) Name() string {
	return "shell"
}

// Collect populates facts.Shell with the current shell, history file info,
// and HISTCONTROL setting.
func (c *ShellCollector) Collect(_ context.Context, facts *schema.Facts) Result {
	home := c.HomeDir
	if home == "" {
		var err error

		home, err = os.UserHomeDir()
		if err != nil {
			return Result{Name: c.Name(), Status: StatusError, Error: fmt.Errorf("determine home dir: %w", err)}
		}
	}

	shell := os.Getenv("SHELL")
	facts.Shell.Shell = shell

	histFile := historyFilePath(shell, home)
	facts.Shell.HistoryFile = histFile

	if histFile != "" {
		if info, err := os.Stat(histFile); err == nil {
			facts.Shell.HistoryFileMode = fmt.Sprintf("%04o", info.Mode().Perm())
		}
	}

	facts.Shell.HistControl = os.Getenv("HISTCONTROL")

	return Result{Name: c.Name(), Status: StatusOK}
}

// historyFilePath returns the conventional history file path for the given shell.
func historyFilePath(shell, home string) string {
	switch filepath.Base(shell) {
	case "bash":
		return filepath.Join(home, ".bash_history")
	case "zsh":
		return filepath.Join(home, ".zsh_history")
	case "fish":
		return filepath.Join(home, ".local", "share", "fish", "fish_history")
	default:
		return ""
	}
}
