package collector_test

import (
	"errors"
	"os"
	"os/user"
	"path/filepath"
	"testing"

	"github.com/franklinkim/bouncer/internal/collector"
	"github.com/franklinkim/bouncer/pkg/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testDockerPath = "/usr/bin/docker"

func TestDockerCollector_NotInstalled(t *testing.T) {
	c := &collector.DockerCollector{
		LookPath: func(name string) (string, error) {
			return "", errors.New("not found")
		},
	}
	facts := schema.NewFacts()
	result := c.Collect(t.Context(), &facts)

	assert.Equal(t, "docker", result.Name)
	assert.Equal(t, collector.StatusSkipped, result.Status)
	assert.False(t, facts.Docker.Installed)
}

func TestDockerCollector_WithSocket(t *testing.T) {
	tmp := t.TempDir()
	sockPath := filepath.Join(tmp, "docker.sock")
	require.NoError(t, os.WriteFile(sockPath, []byte(""), 0600))
	require.NoError(t, os.Chmod(sockPath, 0660))

	c := &collector.DockerCollector{
		SocketPath: sockPath,
		LookPath: func(name string) (string, error) {
			return testDockerPath, nil
		},
		LookupGroup: func(name string) (*user.Group, error) {
			return nil, errors.New("group not found")
		},
	}
	facts := schema.NewFacts()
	result := c.Collect(t.Context(), &facts)

	assert.Equal(t, collector.StatusOK, result.Status)
	assert.True(t, facts.Docker.Installed)
	assert.True(t, facts.Docker.SocketExists)
	assert.Equal(t, "0660", facts.Docker.SocketMode)
	assert.False(t, facts.Docker.UserInGroup)
}

func TestDockerCollector_NoSocket(t *testing.T) {
	c := &collector.DockerCollector{
		SocketPath: "/nonexistent/docker.sock",
		LookPath: func(name string) (string, error) {
			return testDockerPath, nil
		},
		LookupGroup: func(name string) (*user.Group, error) {
			return nil, errors.New("group not found")
		},
	}
	facts := schema.NewFacts()
	result := c.Collect(t.Context(), &facts)

	assert.Equal(t, collector.StatusOK, result.Status)
	assert.True(t, facts.Docker.Installed)
	assert.False(t, facts.Docker.SocketExists)
}

func TestDockerCollector_UserInGroup(t *testing.T) {
	tmp := t.TempDir()
	sockPath := filepath.Join(tmp, "docker.sock")
	require.NoError(t, os.WriteFile(sockPath, []byte(""), 0600))
	require.NoError(t, os.Chmod(sockPath, 0660))

	c := &collector.DockerCollector{
		SocketPath: sockPath,
		LookPath: func(name string) (string, error) {
			return testDockerPath, nil
		},
		LookupGroup: func(name string) (*user.Group, error) {
			return &user.Group{Gid: "999", Name: "docker"}, nil
		},
		CurrentUser: func() (*user.User, error) {
			return &user.User{Uid: "1000", Username: "testuser"}, nil
		},
		GroupIds: func(u *user.User) ([]string, error) {
			return []string{"1000", "999", "100"}, nil
		},
	}
	facts := schema.NewFacts()
	result := c.Collect(t.Context(), &facts)

	assert.Equal(t, collector.StatusOK, result.Status)
	assert.True(t, facts.Docker.UserInGroup)
}

func TestDockerCollector_UserNotInGroup(t *testing.T) {
	c := &collector.DockerCollector{
		SocketPath: "/nonexistent/docker.sock",
		LookPath: func(name string) (string, error) {
			return testDockerPath, nil
		},
		LookupGroup: func(name string) (*user.Group, error) {
			return &user.Group{Gid: "999", Name: "docker"}, nil
		},
		CurrentUser: func() (*user.User, error) {
			return &user.User{Uid: "1000", Username: "testuser"}, nil
		},
		GroupIds: func(u *user.User) ([]string, error) {
			return []string{"1000", "100"}, nil
		},
	}
	facts := schema.NewFacts()
	result := c.Collect(t.Context(), &facts)

	assert.Equal(t, collector.StatusOK, result.Status)
	assert.False(t, facts.Docker.UserInGroup)
}
