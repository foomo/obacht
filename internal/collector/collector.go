package collector

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"

	"github.com/franklinkim/bouncer/pkg/schema"
)

// CollectorStatus represents the three-state status of a collector.
type CollectorStatus string

const (
	StatusOK      CollectorStatus = "ok"
	StatusSkipped CollectorStatus = "skipped"
	StatusError   CollectorStatus = "error"
)

// Result holds the outcome of a collector run.
type Result struct {
	Name   string
	Status CollectorStatus
	Error  error
}

// Collector interface for all fact collectors.
type Collector interface {
	Name() string
	Collect(ctx context.Context, facts *schema.Facts) Result
}

// CollectAll runs all collectors concurrently using errgroup.
//
// Safety: each collector MUST only write to its own dedicated sub-struct field
// within Facts (e.g., the SSH collector writes only to facts.SSH). This avoids
// data races since no two goroutines write to the same memory.
func CollectAll(ctx context.Context, collectors []Collector) (*schema.Facts, []Result, error) {
	facts := schema.NewFacts()
	results := make([]Result, len(collectors))

	g, ctx := errgroup.WithContext(ctx)

	for i, c := range collectors {
		g.Go(func() error {
			results[i] = c.Collect(ctx, &facts)
			return nil // collectors don't fail the group, they return status
		})
	}

	if err := g.Wait(); err != nil {
		return nil, nil, fmt.Errorf("collector group failed: %w", err)
	}

	return &facts, results, nil
}
