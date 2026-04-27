package engine

import "github.com/foomo/obacht/pkg/schema"

// ProgressKind identifies the type of progress event.
type ProgressKind int

const (
	// EventGroupStart is fired before a rule group is evaluated.
	EventGroupStart ProgressKind = iota
	// EventGroupDone is fired after a rule group has been evaluated.
	EventGroupDone
)

// ProgressEvent describes a lifecycle event during evaluation.
type ProgressEvent struct {
	Kind       ProgressKind
	Category   string
	RuleCount  int
	Results    []schema.CheckResult // populated only for EventGroupDone
	GroupIndex int                  // 0-based index of the current group
	GroupTotal int                  // total number of groups
}

// ProgressFunc is a callback invoked during evaluation to report progress.
type ProgressFunc func(ProgressEvent)
