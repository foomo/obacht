package collector

import (
	"context"
	"os/exec"
	"strings"

	"github.com/franklinkim/bouncer/pkg/schema"
)

// securityTools is the curated list of security-relevant tools to check.
var securityTools = []struct {
	name       string
	versionCmd string
	versionArg string
}{
	{"git", "git", "--version"},
	{"opa", "opa", "version"},
	{"gpg", "gpg", "--version"},
	{"ssh-agent", "ssh-agent", "-V"}, // ssh-agent prints version to stderr, but we try anyway
}

// ToolsCollector gathers facts about installed security-relevant tools.
type ToolsCollector struct {
	Runner CommandRunner
}

// NewToolsCollector returns a ToolsCollector that uses real command execution.
func NewToolsCollector() *ToolsCollector {
	return &ToolsCollector{Runner: &execRunner{}}
}

// Name returns the collector name.
func (t *ToolsCollector) Name() string {
	return "tools"
}

// Collect populates Tools facts in the provided Facts struct.
func (t *ToolsCollector) Collect(ctx context.Context, facts *schema.Facts) Result {
	result := Result{Name: t.Name()}

	tools := make([]schema.ToolInfo, 0, len(securityTools))

	for _, tool := range securityTools {
		info := schema.ToolInfo{Name: tool.name}

		path, err := exec.LookPath(tool.name)
		if err != nil {
			// Tool not installed.
			tools = append(tools, info)
			continue
		}

		info.Installed = true
		info.Path = path

		// Try to get version output.
		out, err := t.Runner.Run(ctx, tool.versionCmd, tool.versionArg)
		if err == nil {
			// Take only the first line of output.
			version := strings.TrimSpace(string(out))
			if idx := strings.IndexByte(version, '\n'); idx >= 0 {
				version = version[:idx]
			}

			info.Version = version
		}

		tools = append(tools, info)
	}

	facts.Tools.Tools = tools
	result.Status = StatusOK

	return result
}
