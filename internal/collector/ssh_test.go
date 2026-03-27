package collector_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/franklinkim/bouncer/internal/collector"
	"github.com/franklinkim/bouncer/pkg/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSSHCollector_NoSSHDir(t *testing.T) {
	home := t.TempDir() // no .ssh inside

	c := &collector.SSHCollector{HomeDir: home}
	facts := schema.NewFacts()
	result := c.Collect(t.Context(), &facts)

	assert.Equal(t, "ssh", result.Name)
	assert.Equal(t, collector.StatusSkipped, result.Status)
	assert.False(t, facts.SSH.DirectoryExists)
}

func TestSSHCollector_WithKeys(t *testing.T) {
	home := t.TempDir()
	sshDir := filepath.Join(home, ".ssh")
	require.NoError(t, os.MkdirAll(sshDir, 0700))

	// Create private keys and a public key that should be excluded.
	writeFile(t, filepath.Join(sshDir, "id_rsa"), 0600)
	writeFile(t, filepath.Join(sshDir, "id_rsa.pub"), 0644)
	writeFile(t, filepath.Join(sshDir, "id_ed25519"), 0600)
	writeFile(t, filepath.Join(sshDir, "id_ed25519.pub"), 0644)

	// Create config file.
	writeFile(t, filepath.Join(sshDir, "config"), 0644)

	c := &collector.SSHCollector{HomeDir: home}
	facts := schema.NewFacts()
	result := c.Collect(t.Context(), &facts)

	assert.Equal(t, collector.StatusOK, result.Status)
	require.NoError(t, result.Error)

	assert.True(t, facts.SSH.DirectoryExists)
	assert.Equal(t, "0700", facts.SSH.DirectoryMode)
	assert.True(t, facts.SSH.ConfigExists)

	require.Len(t, facts.SSH.Keys, 2)

	// Build a map for order-independent assertions.
	keysByType := map[string]schema.SSHKey{}
	for _, k := range facts.SSH.Keys {
		keysByType[k.Type] = k
	}

	// resolvePath accounts for macOS /var -> /private/var symlinks in temp dirs.
	resolvedSSHDir, err := filepath.EvalSymlinks(sshDir)
	require.NoError(t, err)

	rsaKey := keysByType["rsa"]
	assert.Equal(t, filepath.Join(resolvedSSHDir, "id_rsa"), rsaKey.Path)
	assert.Equal(t, "0600", rsaKey.Mode)

	edKey := keysByType["ed25519"]
	assert.Equal(t, filepath.Join(resolvedSSHDir, "id_ed25519"), edKey.Path)
	assert.Equal(t, "0600", edKey.Mode)
}

func TestSSHCollector_NoConfig(t *testing.T) {
	home := t.TempDir()
	sshDir := filepath.Join(home, ".ssh")
	require.NoError(t, os.MkdirAll(sshDir, 0700))

	c := &collector.SSHCollector{HomeDir: home}
	facts := schema.NewFacts()
	result := c.Collect(t.Context(), &facts)

	assert.Equal(t, collector.StatusOK, result.Status)
	assert.True(t, facts.SSH.DirectoryExists)
	assert.False(t, facts.SSH.ConfigExists)
	assert.Empty(t, facts.SSH.Keys)
}

func TestSSHCollector_SymlinkedDir(t *testing.T) {
	// Create a real .ssh directory with a key.
	realHome := t.TempDir()
	realSSH := filepath.Join(realHome, ".ssh")
	require.NoError(t, os.MkdirAll(realSSH, 0700))
	writeFile(t, filepath.Join(realSSH, "id_ed25519"), 0600)

	// Create a fake home where .ssh is a symlink to the real one.
	fakeHome := t.TempDir()
	require.NoError(t, os.Symlink(realSSH, filepath.Join(fakeHome, ".ssh")))

	c := &collector.SSHCollector{HomeDir: fakeHome}
	facts := schema.NewFacts()
	result := c.Collect(t.Context(), &facts)

	assert.Equal(t, collector.StatusOK, result.Status)
	assert.True(t, facts.SSH.DirectoryExists)
	assert.Equal(t, "0700", facts.SSH.DirectoryMode)
	require.Len(t, facts.SSH.Keys, 1)
	assert.Equal(t, "ed25519", facts.SSH.Keys[0].Type)
}

func TestKeyTypeFromFilename(t *testing.T) {
	tests := []struct {
		filename string
		expected string
	}{
		{"id_rsa", "rsa"},
		{"id_ed25519", "ed25519"},
		{"id_ecdsa", "ecdsa"},
		{"id_dsa", "dsa"},
		{"custom_key", "custom_key"},
	}
	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			assert.Equal(t, tt.expected, collector.ExportKeyTypeFromFilename(tt.filename))
		})
	}
}

// writeFile creates a file with the given permissions and dummy content.
func writeFile(t *testing.T, path string, mode os.FileMode) {
	t.Helper()
	require.NoError(t, os.WriteFile(path, []byte("test"), mode))
}
