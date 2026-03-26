package collector

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"slices"

	"github.com/franklinkim/bouncer/pkg/schema"
)

// DockerCollector gathers facts about the Docker installation and socket.
type DockerCollector struct {
	// HomeDir overrides the user's home directory (unused currently but kept
	// for consistency with other collectors). When empty, os.UserHomeDir is used.
	HomeDir string
	// SocketPath overrides the default Docker socket path for testing.
	SocketPath string
	// LookPath overrides exec.LookPath for testing.
	LookPath func(string) (string, error)
	// LookupGroup overrides user.LookupGroup for testing.
	LookupGroup func(string) (*user.Group, error)
	// CurrentUser overrides user.Current for testing.
	CurrentUser func() (*user.User, error)
	// GroupIds overrides u.GroupIds for testing.
	GroupIds func(*user.User) ([]string, error)
}

// NewDockerCollector returns a DockerCollector that uses system defaults.
func NewDockerCollector() *DockerCollector {
	return &DockerCollector{}
}

// Name returns the collector name.
func (c *DockerCollector) Name() string {
	return "docker"
}

// Collect populates facts.Docker with information about Docker installation,
// socket permissions, and group membership.
func (c *DockerCollector) Collect(ctx context.Context, facts *schema.Facts) Result {
	lookPath := c.LookPath
	if lookPath == nil {
		lookPath = exec.LookPath
	}

	// Check if docker is installed.
	_, err := lookPath("docker")
	if err != nil {
		facts.Docker = schema.DockerFacts{
			Installed: false,
		}

		return Result{Name: c.Name(), Status: StatusSkipped}
	}

	facts.Docker.Installed = true

	// Check Docker socket.
	sockPath := c.SocketPath
	if sockPath == "" {
		sockPath = "/var/run/docker.sock"
	}

	info, err := os.Stat(sockPath)
	if err == nil {
		facts.Docker.SocketExists = true
		facts.Docker.SocketMode = fmt.Sprintf("%04o", info.Mode().Perm())
	}

	// Check if current user is in the docker group.
	facts.Docker.UserInGroup = c.isUserInDockerGroup()

	return Result{Name: c.Name(), Status: StatusOK}
}

// isUserInDockerGroup checks whether the current user is a member of the "docker" group.
func (c *DockerCollector) isUserInDockerGroup() bool {
	lookupGroup := c.LookupGroup
	if lookupGroup == nil {
		lookupGroup = user.LookupGroup
	}

	grp, err := lookupGroup("docker")
	if err != nil {
		return false // group doesn't exist
	}

	currentUser := c.CurrentUser
	if currentUser == nil {
		currentUser = user.Current
	}

	u, err := currentUser()
	if err != nil {
		return false
	}

	groupIds := c.GroupIds
	if groupIds == nil {
		groupIds = func(u *user.User) ([]string, error) {
			return u.GroupIds()
		}
	}

	gids, err := groupIds(u)
	if err != nil {
		return false
	}

	return slices.Contains(gids, grp.Gid)
}
