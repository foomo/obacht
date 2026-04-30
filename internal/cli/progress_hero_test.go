//go:build safe

package cli_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/foomo/obacht/internal/cli"
)

func TestHeroRendersVersion(t *testing.T) {
	prev := cli.SetVersionForTest("1.2.3")

	t.Cleanup(func() { cli.SetVersionForTest(prev) })

	out := cli.RenderHeroForTest()

	assert.Contains(t, out, "developer environment scanner")
	assert.Contains(t, out, "v1.2.3")
}

func TestHeroRendersDevWithoutPrefix(t *testing.T) {
	prev := cli.SetVersionForTest("dev")

	t.Cleanup(func() { cli.SetVersionForTest(prev) })

	out := cli.RenderHeroForTest()

	assert.Contains(t, out, "dev")
	assert.NotContains(t, out, "vdev")
}
