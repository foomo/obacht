package collector

import (
	"context"
	"fmt"

	"github.com/franklinkim/bouncer/pkg/schema"
)

// Exported aliases for internal functions used by collector_test package.
var (
	ExportParseGitConfig      = parseGitConfig
	ExportIsDirWritable       = isDirWritable
	ExportKeyTypeFromFilename = keyTypeFromFilename
	ExportHistoryFilePath     = historyFilePath
	ExportSecurityToolsCount  = len(securityTools)
)

// NewFakeRunner creates a test CommandRunner from a results map.
func NewFakeRunner(results map[string]FakeResult) CommandRunner {
	return &fakeRunnerExport{results: results}
}

// FakeResult holds test output for a faked command.
type FakeResult struct {
	Output []byte
	Err    error
}

type fakeRunnerExport struct {
	results map[string]FakeResult
}

func (f *fakeRunnerExport) Run(_ context.Context, name string, args ...string) ([]byte, error) {
	key := name
	if len(args) > 0 {
		key += " " + fmt.Sprintf("%s", args)
	}

	if r, ok := f.results[key]; ok {
		return r.Output, r.Err
	}

	return nil, fmt.Errorf("unexpected command: %s", key)
}

// ExportParseGitConfigFunc is a type alias for reference.
type ExportParseGitConfigFunc = func(string, *schema.GitFacts)
