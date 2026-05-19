package runner_test

import (
	"strings"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/foomo/obacht/internal/runner"
)

func TestPreprocess_NoIncludes(t *testing.T) {
	fsys := fstest.MapFS{}
	script := "#!/bin/sh\necho hello\n"

	out, err := runner.Preprocess(fsys, "inputs/ssh/SSH001.sh", script)
	require.NoError(t, err)
	assert.Equal(t, script, out)
}

func TestPreprocess_SingleInclude(t *testing.T) {
	fsys := fstest.MapFS{
		"inputs/_lib/json.sh": &fstest.MapFile{Data: []byte("emit_ok() { echo ok; }\n")},
	}
	script := "#!/bin/sh\n# include: _lib/json.sh\necho hello\n"

	out, err := runner.Preprocess(fsys, "inputs/ssh/SSH001.sh", script)
	require.NoError(t, err)
	assert.Contains(t, out, "emit_ok() { echo ok; }")
	assert.Contains(t, out, "echo hello")
	assert.True(t, strings.HasPrefix(out, "#!/bin/sh"))
	assert.NotContains(t, out, "# include:")
}

func TestPreprocess_MultipleIncludes_PreservesOrder(t *testing.T) {
	fsys := fstest.MapFS{
		"inputs/_lib/a.sh": &fstest.MapFile{Data: []byte("# LIB_A\n")},
		"inputs/_lib/b.sh": &fstest.MapFile{Data: []byte("# LIB_B\n")},
	}
	script := "#!/bin/sh\n# include: _lib/a.sh\n# include: _lib/b.sh\nbody\n"

	out, err := runner.Preprocess(fsys, "inputs/x/X001.sh", script)
	require.NoError(t, err)

	idxA := strings.Index(out, "# LIB_A")
	idxB := strings.Index(out, "# LIB_B")

	assert.Greater(t, idxA, -1)
	assert.Greater(t, idxB, idxA)
}

func TestPreprocess_DeduplicatesIncludes(t *testing.T) {
	fsys := fstest.MapFS{
		"inputs/_lib/a.sh": &fstest.MapFile{Data: []byte("# LIB_A\n")},
	}
	script := "#!/bin/sh\n# include: _lib/a.sh\n# include: _lib/a.sh\nbody\n"

	out, err := runner.Preprocess(fsys, "inputs/x/X001.sh", script)
	require.NoError(t, err)
	assert.Equal(t, 1, strings.Count(out, "# LIB_A"))
}

func TestPreprocess_TransitiveIncludes(t *testing.T) {
	fsys := fstest.MapFS{
		"inputs/_lib/a.sh": &fstest.MapFile{Data: []byte("# include: _lib/b.sh\n# LIB_A\n")},
		"inputs/_lib/b.sh": &fstest.MapFile{Data: []byte("# LIB_B\n")},
	}
	script := "#!/bin/sh\n# include: _lib/a.sh\nbody\n"

	out, err := runner.Preprocess(fsys, "inputs/x/X001.sh", script)
	require.NoError(t, err)

	idxB := strings.Index(out, "# LIB_B")
	idxA := strings.Index(out, "# LIB_A")

	assert.Greater(t, idxB, -1)
	assert.Greater(t, idxA, idxB) // b included before a's body
}

func TestPreprocess_MissingInclude_Errors(t *testing.T) {
	fsys := fstest.MapFS{}
	script := "#!/bin/sh\n# include: _lib/missing.sh\n"

	_, err := runner.Preprocess(fsys, "inputs/x/X001.sh", script)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing.sh")
}

func TestPreprocess_CyclicInclude_Errors(t *testing.T) {
	fsys := fstest.MapFS{
		"inputs/_lib/a.sh": &fstest.MapFile{Data: []byte("# include: _lib/b.sh\n")},
		"inputs/_lib/b.sh": &fstest.MapFile{Data: []byte("# include: _lib/a.sh\n")},
	}
	script := "#!/bin/sh\n# include: _lib/a.sh\n"

	_, err := runner.Preprocess(fsys, "inputs/x/X001.sh", script)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cycle")
}

func TestPreprocess_PathTraversal_Errors(t *testing.T) {
	fsys := fstest.MapFS{}
	script := "#!/bin/sh\n# include: ../../../etc/passwd\n"

	_, err := runner.Preprocess(fsys, "inputs/x/X001.sh", script)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "_lib/")
}

func TestPreprocess_IncludesOnlyAtTop(t *testing.T) {
	fsys := fstest.MapFS{
		"inputs/_lib/a.sh": &fstest.MapFile{Data: []byte("# LIB_A\n")},
	}
	// include directive AFTER a body line — should be ignored as a regular comment
	script := "#!/bin/sh\necho first\n# include: _lib/a.sh\n"

	out, err := runner.Preprocess(fsys, "inputs/x/X001.sh", script)
	require.NoError(t, err)
	assert.NotContains(t, out, "# LIB_A")
}

func TestPreprocess_IndentedIncludeIsComment(t *testing.T) {
	fsys := fstest.MapFS{
		"inputs/_lib/a.sh": &fstest.MapFile{Data: []byte("# LIB_A\n")},
	}
	// Leading whitespace disqualifies the directive; treated as a plain comment.
	script := "#!/bin/sh\n  # include: _lib/a.sh\nbody\n"

	out, err := runner.Preprocess(fsys, "inputs/x/X001.sh", script)
	require.NoError(t, err)
	assert.NotContains(t, out, "# LIB_A")
	assert.Contains(t, out, "# include: _lib/a.sh") // preserved as comment
}
