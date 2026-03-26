package collector

import (
	"context"
	"os"
	"runtime"

	"github.com/franklinkim/bouncer/pkg/schema"
)

// OSCollector gathers facts about the operating system.
type OSCollector struct{}

// NewOSCollector returns a new OSCollector.
func NewOSCollector() *OSCollector {
	return &OSCollector{}
}

// Name returns the collector name.
func (c *OSCollector) Name() string {
	return "os"
}

// Collect populates facts.OS with operating system information.
func (c *OSCollector) Collect(_ context.Context, facts *schema.Facts) Result {
	facts.OS.OS = runtime.GOOS
	facts.OS.Arch = runtime.GOARCH

	hostname, _ := os.Hostname()
	facts.OS.Hostname = hostname

	return Result{Name: c.Name(), Status: StatusOK}
}
