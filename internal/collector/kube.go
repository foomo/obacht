package collector

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/franklinkim/bouncer/pkg/schema"
	"gopkg.in/yaml.v3"
)

// KubeCollector gathers facts about the user's Kubernetes configuration.
type KubeCollector struct {
	// HomeDir overrides the user's home directory. When empty, os.UserHomeDir
	// is used.
	HomeDir string
}

// NewKubeCollector returns a KubeCollector that uses the real home directory.
func NewKubeCollector() *KubeCollector {
	return &KubeCollector{}
}

// Name returns the collector name.
func (c *KubeCollector) Name() string {
	return "kube"
}

// kubeconfig is a minimal representation of a kubeconfig file, containing
// only the fields we need to extract context and cluster names.
type kubeconfig struct {
	Contexts []kubeconfigContext `yaml:"contexts"`
}

type kubeconfigContext struct {
	Name    string               `yaml:"name"`
	Context kubeconfigContextRef `yaml:"context"`
}

type kubeconfigContextRef struct {
	Cluster string `yaml:"cluster"`
}

// Collect populates facts.Kube with information about the ~/.kube/config file,
// including permissions and configured contexts.
func (c *KubeCollector) Collect(ctx context.Context, facts *schema.Facts) Result {
	configPath, err := c.configPath()
	if err != nil {
		return Result{Name: c.Name(), Status: StatusError, Error: fmt.Errorf("determine home dir: %w", err)}
	}

	configPath = resolvePath(configPath)

	info, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		facts.Kube = schema.KubeFacts{
			ConfigExists: false,
		}

		return Result{Name: c.Name(), Status: StatusSkipped}
	}

	if err != nil {
		return Result{Name: c.Name(), Status: StatusError, Error: fmt.Errorf("stat %s: %w", configPath, err)}
	}

	facts.Kube.ConfigExists = true
	facts.Kube.ConfigMode = fmt.Sprintf("%04o", info.Mode().Perm())

	// Parse kubeconfig to extract contexts.
	data, err := os.ReadFile(configPath)
	if err != nil {
		return Result{Name: c.Name(), Status: StatusError, Error: fmt.Errorf("read %s: %w", configPath, err)}
	}

	var kc kubeconfig
	if err := yaml.Unmarshal(data, &kc); err != nil {
		return Result{Name: c.Name(), Status: StatusError, Error: fmt.Errorf("parse kubeconfig: %w", err)}
	}

	for _, ctxEntry := range kc.Contexts {
		facts.Kube.Contexts = append(facts.Kube.Contexts, schema.KubeContext{
			Name:    ctxEntry.Name,
			Cluster: ctxEntry.Context.Cluster,
		})
	}

	return Result{Name: c.Name(), Status: StatusOK}
}

// configPath returns the path to the kubeconfig file.
func (c *KubeCollector) configPath() (string, error) {
	home := c.HomeDir
	if home == "" {
		var err error

		home, err = os.UserHomeDir()
		if err != nil {
			return "", err
		}
	}

	return filepath.Join(home, ".kube", "config"), nil
}
