package collector

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/franklinkim/bouncer/pkg/schema"
)

// PathCollector gathers facts about directories in $PATH.
type PathCollector struct{}

// NewPathCollector returns a new PathCollector.
func NewPathCollector() *PathCollector {
	return &PathCollector{}
}

// Name returns the collector name.
func (p *PathCollector) Name() string {
	return "path"
}

// Collect populates Path facts in the provided Facts struct.
func (p *PathCollector) Collect(_ context.Context, facts *schema.Facts) Result {
	result := Result{Name: p.Name()}

	pathEnv := os.Getenv("PATH")
	entries := strings.Split(pathEnv, string(os.PathListSeparator))

	dirs := make([]schema.PathDir, 0, len(entries))
	for _, entry := range entries {
		if entry == "" {
			continue
		}

		dir := schema.PathDir{
			Path:       entry,
			IsRelative: !filepath.IsAbs(entry),
		}

		info, err := os.Stat(filepath.Clean(entry))
		if err == nil && info.IsDir() {
			dir.Exists = true
			dir.Writable = isDirWritable(entry)
		}

		dirs = append(dirs, dir)
	}

	facts.Path.Dirs = dirs
	result.Status = StatusOK

	return result
}

// isDirWritable checks if the directory is writable by the current user
// by attempting to create a temporary file in it.
func isDirWritable(path string) bool {
	f, err := os.CreateTemp(path, ".bouncer-check-*")
	if err != nil {
		return false
	}

	name := f.Name()
	f.Close()
	os.Remove(filepath.Clean(name))

	return true
}
