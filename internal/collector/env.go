package collector

import (
	"context"
	"os"
	"strings"

	"github.com/franklinkim/bouncer/pkg/schema"
)

// DefaultExactMatches is the curated allowlist of dangerous env var names.
var DefaultExactMatches = []string{
	"AWS_SECRET_ACCESS_KEY",
	"AWS_ACCESS_KEY_ID",
	"GITHUB_TOKEN",
	"GITLAB_TOKEN",
	"NPM_TOKEN",
	"DOCKER_PASSWORD",
	"SLACK_TOKEN",
	"SLACK_WEBHOOK_URL",
	"DATABASE_URL",
	"MYSQL_PASSWORD",
	"POSTGRES_PASSWORD",
	"REDIS_URL",
}

// DefaultSuffixPatterns is the curated list of dangerous env var suffix patterns.
var DefaultSuffixPatterns = []string{
	"_PASSWORD",
	"_SECRET_KEY",
	"_API_KEY",
	"_PRIVATE_KEY",
}

// EnvCollector scans environment variables against known suspicious patterns.
// It NEVER stores the value — only the variable name and the pattern that matched.
type EnvCollector struct {
	// ExactMatches lists env var names that are flagged by exact match.
	ExactMatches []string
	// SuffixPatterns lists suffixes that flag an env var by suffix match.
	SuffixPatterns []string
	// EnvironFunc overrides os.Environ for testing.
	EnvironFunc func() []string
}

// NewEnvCollector returns an EnvCollector with default patterns.
func NewEnvCollector() *EnvCollector {
	return &EnvCollector{
		ExactMatches:   DefaultExactMatches,
		SuffixPatterns: DefaultSuffixPatterns,
	}
}

// Name returns the collector name.
func (c *EnvCollector) Name() string {
	return "env"
}

// Collect populates facts.Env with any environment variables that match
// the configured suspicious patterns.
func (c *EnvCollector) Collect(_ context.Context, facts *schema.Facts) Result {
	environ := os.Environ
	if c.EnvironFunc != nil {
		environ = c.EnvironFunc
	}

	exactSet := make(map[string]bool, len(c.ExactMatches))
	for _, name := range c.ExactMatches {
		exactSet[name] = true
	}

	var suspicious []schema.SuspiciousVar

	for _, entry := range environ() {
		name, _, ok := strings.Cut(entry, "=")
		if !ok {
			continue
		}

		if exactSet[name] {
			suspicious = append(suspicious, schema.SuspiciousVar{
				Name:    name,
				Pattern: "exact:" + name,
			})

			continue
		}

		for _, suffix := range c.SuffixPatterns {
			if strings.HasSuffix(name, suffix) {
				suspicious = append(suspicious, schema.SuspiciousVar{
					Name:    name,
					Pattern: "*" + suffix,
				})

				break
			}
		}
	}

	facts.Env = schema.EnvFacts{
		SuspiciousVars: suspicious,
	}

	return Result{Name: c.Name(), Status: StatusOK}
}
