package collector_test

import (
	"os"
	"testing"

	"github.com/franklinkim/bouncer/internal/collector"
	"github.com/franklinkim/bouncer/pkg/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPathCollector_Collect(t *testing.T) {
	// Create a temp writable directory.
	writableDir := t.TempDir()

	// Build a controlled PATH with:
	// 1. A writable directory that exists
	// 2. A non-existent directory
	// 3. A relative path
	nonExistent := "/tmp/bouncer-does-not-exist-xyz"
	relativePath := "relative/bin"

	controlledPATH := writableDir + ":" + nonExistent + ":" + relativePath
	t.Setenv("PATH", controlledPATH)

	c := collector.NewPathCollector()
	facts := schema.NewFacts()

	result := c.Collect(t.Context(), &facts)

	require.NoError(t, result.Error)
	assert.Equal(t, collector.StatusOK, result.Status)
	assert.Equal(t, "path", result.Name)

	require.Len(t, facts.Path.Dirs, 3)

	// First entry: writable, exists, absolute.
	assert.Equal(t, writableDir, facts.Path.Dirs[0].Path)
	assert.True(t, facts.Path.Dirs[0].Exists, "temp dir should exist")
	assert.True(t, facts.Path.Dirs[0].Writable, "temp dir should be writable")
	assert.False(t, facts.Path.Dirs[0].IsRelative, "temp dir should be absolute")

	// Second entry: non-existent.
	assert.Equal(t, nonExistent, facts.Path.Dirs[1].Path)
	assert.False(t, facts.Path.Dirs[1].Exists, "non-existent dir should not exist")
	assert.False(t, facts.Path.Dirs[1].Writable, "non-existent dir should not be writable")
	assert.False(t, facts.Path.Dirs[1].IsRelative, "absolute path should not be relative")

	// Third entry: relative path.
	assert.Equal(t, relativePath, facts.Path.Dirs[2].Path)
	assert.True(t, facts.Path.Dirs[2].IsRelative, "relative path should be detected")
}

func TestPathCollector_EmptyPATH(t *testing.T) {
	t.Setenv("PATH", "")

	c := collector.NewPathCollector()
	facts := schema.NewFacts()

	result := c.Collect(t.Context(), &facts)

	assert.Equal(t, collector.StatusOK, result.Status)
	assert.Empty(t, facts.Path.Dirs)
}

func TestIsDirWritable(t *testing.T) {
	// Writable temp dir.
	dir := t.TempDir()
	assert.True(t, collector.ExportIsDirWritable(dir))

	// Non-existent dir.
	assert.False(t, collector.ExportIsDirWritable("/tmp/bouncer-nonexistent-dir-check"))

	// Read-only dir.
	roDir := t.TempDir()
	require.NoError(t, os.Chmod(roDir, 0o555))
	t.Cleanup(func() { _ = os.Chmod(roDir, 0o755) })
	assert.False(t, collector.ExportIsDirWritable(roDir))
}
