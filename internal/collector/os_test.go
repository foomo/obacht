package collector_test

import (
	"runtime"
	"testing"

	"github.com/franklinkim/bouncer/internal/collector"
	"github.com/franklinkim/bouncer/pkg/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOSCollector(t *testing.T) {
	c := collector.NewOSCollector()
	facts := schema.NewFacts()
	result := c.Collect(t.Context(), &facts)

	assert.Equal(t, "os", result.Name)
	assert.Equal(t, collector.StatusOK, result.Status)
	require.NoError(t, result.Error)

	assert.Equal(t, runtime.GOOS, facts.OS.OS)
	assert.Equal(t, runtime.GOARCH, facts.OS.Arch)
	assert.NotEmpty(t, facts.OS.Hostname)
}
