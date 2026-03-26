package collector_test

import (
	"context"
	"testing"

	"github.com/franklinkim/bouncer/internal/collector"
	"github.com/franklinkim/bouncer/pkg/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnvCollector_ExactMatch(t *testing.T) {
	c := &collector.EnvCollector{
		ExactMatches:   []string{"GITHUB_TOKEN", "AWS_SECRET_ACCESS_KEY"},
		SuffixPatterns: collector.DefaultSuffixPatterns,
		EnvironFunc: func() []string {
			return []string{
				"GITHUB_TOKEN=ghp_secret123",
				"HOME=/home/user",
				"PATH=/usr/bin",
			}
		},
	}

	facts := schema.NewFacts()
	result := c.Collect(context.Background(), &facts)

	assert.Equal(t, collector.StatusOK, result.Status)
	require.Len(t, facts.Env.SuspiciousVars, 1)
	assert.Equal(t, "GITHUB_TOKEN", facts.Env.SuspiciousVars[0].Name)
	assert.Equal(t, "exact:GITHUB_TOKEN", facts.Env.SuspiciousVars[0].Pattern)
}

func TestEnvCollector_SuffixMatch(t *testing.T) {
	c := &collector.EnvCollector{
		ExactMatches:   collector.DefaultExactMatches,
		SuffixPatterns: collector.DefaultSuffixPatterns,
		EnvironFunc: func() []string {
			return []string{
				"MY_APP_API_KEY=abc123",
				"STRIPE_PRIVATE_KEY=sk_test_xxx",
				"HOME=/home/user",
			}
		},
	}

	facts := schema.NewFacts()
	result := c.Collect(context.Background(), &facts)

	assert.Equal(t, collector.StatusOK, result.Status)
	require.Len(t, facts.Env.SuspiciousVars, 2)

	byName := map[string]schema.SuspiciousVar{}
	for _, v := range facts.Env.SuspiciousVars {
		byName[v.Name] = v
	}

	assert.Equal(t, "*_API_KEY", byName["MY_APP_API_KEY"].Pattern)
	assert.Equal(t, "*_PRIVATE_KEY", byName["STRIPE_PRIVATE_KEY"].Pattern)
}

func TestEnvCollector_NoMatches(t *testing.T) {
	c := &collector.EnvCollector{
		ExactMatches:   collector.DefaultExactMatches,
		SuffixPatterns: collector.DefaultSuffixPatterns,
		EnvironFunc: func() []string {
			return []string{
				"HOME=/home/user",
				"EDITOR=vim",
				"TERM=xterm-256color",
			}
		},
	}

	facts := schema.NewFacts()
	result := c.Collect(context.Background(), &facts)

	assert.Equal(t, collector.StatusOK, result.Status)
	assert.Empty(t, facts.Env.SuspiciousVars)
}

func TestEnvCollector_NeverStoresValues(t *testing.T) {
	c := &collector.EnvCollector{
		ExactMatches:   []string{"GITHUB_TOKEN"},
		SuffixPatterns: []string{"_PASSWORD"},
		EnvironFunc: func() []string {
			return []string{
				"GITHUB_TOKEN=super_secret_value",
				"DB_PASSWORD=hunter2",
			}
		},
	}

	facts := schema.NewFacts()
	c.Collect(context.Background(), &facts)

	for _, v := range facts.Env.SuspiciousVars {
		assert.NotContains(t, v.Name, "super_secret_value")
		assert.NotContains(t, v.Pattern, "super_secret_value")
		assert.NotContains(t, v.Name, "hunter2")
		assert.NotContains(t, v.Pattern, "hunter2")
	}
}
