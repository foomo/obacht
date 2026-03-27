package collector

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/franklinkim/bouncer/pkg/schema"
)

// SSHCollector gathers facts about the user's SSH configuration and keys.
type SSHCollector struct {
	// HomeDir overrides the user's home directory. When empty, os.UserHomeDir
	// is used.
	HomeDir string
}

// NewSSHCollector returns an SSHCollector that uses the real home directory.
func NewSSHCollector() *SSHCollector {
	return &SSHCollector{}
}

// Name returns the collector name.
func (c *SSHCollector) Name() string {
	return "ssh"
}

// Collect populates facts.SSH with information about the ~/.ssh directory,
// private keys, and config file.
func (c *SSHCollector) Collect(ctx context.Context, facts *schema.Facts) Result {
	sshDir, err := c.sshDir()
	if err != nil {
		return Result{Name: c.Name(), Status: StatusError, Error: fmt.Errorf("determine home dir: %w", err)}
	}

	// Resolve symlinks and check if ~/.ssh exists.
	sshDir = resolvePath(sshDir)

	info, err := os.Stat(sshDir)
	if os.IsNotExist(err) {
		facts.SSH = schema.SSHFacts{
			DirectoryExists: false,
		}

		return Result{Name: c.Name(), Status: StatusSkipped}
	}

	if err != nil {
		return Result{Name: c.Name(), Status: StatusError, Error: fmt.Errorf("stat %s: %w", sshDir, err)}
	}

	facts.SSH.DirectoryExists = true
	facts.SSH.DirectoryMode = fmt.Sprintf("%04o", info.Mode().Perm())

	// Discover private keys by globbing id_* and excluding .pub files.
	matches, err := filepath.Glob(filepath.Join(sshDir, "id_*"))
	if err != nil {
		return Result{Name: c.Name(), Status: StatusError, Error: fmt.Errorf("glob keys: %w", err)}
	}

	for _, m := range matches {
		if strings.HasSuffix(m, ".pub") {
			continue
		}

		ki, err := os.Stat(m)
		if err != nil {
			// Skip keys we cannot stat rather than failing the whole collector.
			continue
		}

		facts.SSH.Keys = append(facts.SSH.Keys, schema.SSHKey{
			Path: m,
			Mode: fmt.Sprintf("%04o", ki.Mode().Perm()),
			Type: keyTypeFromFilename(filepath.Base(m)),
		})
	}

	// Check if ~/.ssh/config exists.
	configPath := filepath.Join(sshDir, "config")
	if _, err := os.Stat(configPath); err == nil {
		facts.SSH.ConfigExists = true
	}

	return Result{Name: c.Name(), Status: StatusOK}
}

// sshDir returns the path to the SSH directory.
func (c *SSHCollector) sshDir() (string, error) {
	home := c.HomeDir
	if home == "" {
		var err error

		home, err = os.UserHomeDir()
		if err != nil {
			return "", err
		}
	}

	return filepath.Join(home, ".ssh"), nil
}

// keyTypeFromFilename derives the key type from a private key filename.
// For example "id_ed25519" -> "ed25519", "id_rsa" -> "rsa".
// If the name does not follow the id_<type> convention, the full filename is returned.
func keyTypeFromFilename(name string) string {
	if after, ok := strings.CutPrefix(name, "id_"); ok {
		return after
	}

	return name
}
