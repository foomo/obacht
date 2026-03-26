package collector

import (
	"context"
	"runtime"
	"testing"

	"github.com/franklinkim/bouncer/pkg/schema"
	"github.com/stretchr/testify/assert"
)

func TestOSCollector(t *testing.T) {
	c := NewOSCollector()
	facts := schema.NewFacts()
	result := c.Collect(context.Background(), &facts)

	assert.Equal(t, "os", result.Name)
	assert.Equal(t, StatusOK, result.Status)
	assert.Nil(t, result.Error)

	assert.Equal(t, runtime.GOOS, facts.OS.OS)
	assert.Equal(t, runtime.GOARCH, facts.OS.Arch)
	assert.NotEmpty(t, facts.OS.Hostname)
}
