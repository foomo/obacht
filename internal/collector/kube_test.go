package collector_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/franklinkim/bouncer/internal/collector"
	"github.com/franklinkim/bouncer/pkg/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testKubeconfig = `apiVersion: v1
kind: Config
contexts:
- name: dev-context
  context:
    cluster: dev-cluster
    user: dev-user
- name: prod-context
  context:
    cluster: prod-cluster
    user: prod-user
clusters:
- name: dev-cluster
  cluster:
    server: https://dev.example.com
- name: prod-cluster
  cluster:
    server: https://prod.example.com
`

func TestKubeCollector_NoConfig(t *testing.T) {
	home := t.TempDir() // no .kube inside

	c := &collector.KubeCollector{HomeDir: home}
	facts := schema.NewFacts()
	result := c.Collect(context.Background(), &facts)

	assert.Equal(t, "kube", result.Name)
	assert.Equal(t, collector.StatusSkipped, result.Status)
	assert.False(t, facts.Kube.ConfigExists)
}

func TestKubeCollector_WithConfig(t *testing.T) {
	home := t.TempDir()
	kubeDir := filepath.Join(home, ".kube")
	require.NoError(t, os.MkdirAll(kubeDir, 0700))

	configPath := filepath.Join(kubeDir, "config")
	require.NoError(t, os.WriteFile(configPath, []byte(testKubeconfig), 0600))

	c := &collector.KubeCollector{HomeDir: home}
	facts := schema.NewFacts()
	result := c.Collect(context.Background(), &facts)

	assert.Equal(t, collector.StatusOK, result.Status)
	require.NoError(t, result.Error)
	assert.True(t, facts.Kube.ConfigExists)
	assert.Equal(t, "0600", facts.Kube.ConfigMode)

	require.Len(t, facts.Kube.Contexts, 2)

	// Build a map for order-independent assertions.
	ctxByName := map[string]schema.KubeContext{}
	for _, c := range facts.Kube.Contexts {
		ctxByName[c.Name] = c
	}

	assert.Equal(t, "dev-cluster", ctxByName["dev-context"].Cluster)
	assert.Equal(t, "prod-cluster", ctxByName["prod-context"].Cluster)
}

func TestKubeCollector_WeakPermissions(t *testing.T) {
	home := t.TempDir()
	kubeDir := filepath.Join(home, ".kube")
	require.NoError(t, os.MkdirAll(kubeDir, 0700))

	configPath := filepath.Join(kubeDir, "config")
	require.NoError(t, os.WriteFile(configPath, []byte(testKubeconfig), 0644))

	c := &collector.KubeCollector{HomeDir: home}
	facts := schema.NewFacts()
	result := c.Collect(context.Background(), &facts)

	assert.Equal(t, collector.StatusOK, result.Status)
	assert.True(t, facts.Kube.ConfigExists)
	assert.Equal(t, "0644", facts.Kube.ConfigMode)
}

func TestKubeCollector_EmptyConfig(t *testing.T) {
	home := t.TempDir()
	kubeDir := filepath.Join(home, ".kube")
	require.NoError(t, os.MkdirAll(kubeDir, 0700))

	configPath := filepath.Join(kubeDir, "config")
	require.NoError(t, os.WriteFile(configPath, []byte("apiVersion: v1\nkind: Config\n"), 0600))

	c := &collector.KubeCollector{HomeDir: home}
	facts := schema.NewFacts()
	result := c.Collect(context.Background(), &facts)

	assert.Equal(t, collector.StatusOK, result.Status)
	assert.True(t, facts.Kube.ConfigExists)
	assert.Empty(t, facts.Kube.Contexts)
}
