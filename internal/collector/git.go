package collector

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/franklinkim/bouncer/pkg/schema"
)

// CommandRunner abstracts command execution so it can be replaced in tests.
type CommandRunner interface {
	// Run executes the named program with the given arguments and returns
	// combined stdout output or an error.
	Run(ctx context.Context, name string, args ...string) ([]byte, error)
}

// execRunner is the default CommandRunner that delegates to os/exec.
type execRunner struct{}

func (e *execRunner) Run(ctx context.Context, name string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, name, args...)

	var out bytes.Buffer

	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

// GitCollector gathers Git-related facts.
type GitCollector struct {
	Runner CommandRunner
}

// NewGitCollector returns a GitCollector that uses real command execution.
func NewGitCollector() *GitCollector {
	return &GitCollector{Runner: &execRunner{}}
}

// Name returns the collector name.
func (g *GitCollector) Name() string {
	return "git"
}

// Collect populates Git facts in the provided Facts struct.
func (g *GitCollector) Collect(ctx context.Context, facts *schema.Facts) Result {
	result := Result{Name: g.Name()}

	// Check if git is installed.
	if _, err := exec.LookPath("git"); err != nil {
		result.Status = StatusSkipped
		return result
	}

	facts.Git.Installed = true

	// Get git version.
	versionOut, err := g.Runner.Run(ctx, "git", "--version")
	if err != nil {
		result.Status = StatusError
		result.Error = fmt.Errorf("git --version: %w", err)

		return result
	}

	facts.Git.Version = strings.TrimSpace(string(versionOut))

	// Get global git config.
	configOut, err := g.Runner.Run(ctx, "git", "config", "--global", "--list")
	if err != nil {
		// git config --global --list exits non-zero when there is no global
		// config file; treat that as empty config rather than an error.
		configOut = nil
	}

	parseGitConfig(string(configOut), &facts.Git)

	result.Status = StatusOK

	return result
}

// parseGitConfig extracts relevant fields from git config key=value output.
func parseGitConfig(output string, git *schema.GitFacts) {
	for line := range strings.SplitSeq(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Format: key=value (first '=' separates key from value)
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.ToLower(parts[0])
		value := parts[1]

		switch key {
		case "credential.helper":
			git.CredentialHelper = value
		case "commit.gpgsign":
			git.SigningEnabled = strings.EqualFold(value, "true")
		case "gpg.format":
			git.SigningFormat = value
		}
	}
}
