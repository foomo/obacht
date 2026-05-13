package runner

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"os/exec"
	"time"
)

// DefaultTimeout is the maximum time a single input script may run.
const DefaultTimeout = 30 * time.Second

// InputResult holds the outcome of running a single input script.
type InputResult struct {
	// Data is the parsed JSON data from the script's envelope (or bare JSON for legacy).
	Data any
	// Status indicates whether the script succeeded, was skipped, or errored.
	Status Status
	// SkipReason is set when Status is StatusSkipped due to an envelope skip.
	SkipReason string
	// Error is set when Status is StatusError.
	Error error
}

// Status represents the outcome of running an input script.
type Status string

const (
	StatusOK      Status = "ok"
	StatusSkipped Status = "skipped"
	StatusError   Status = "error"
)

// RunInput executes a shell script and parses its stdout as JSON.
// If script is empty, the result is skipped.
// This is the legacy entrypoint that accepts only bare-JSON output and
// has no preprocessing. New code should use RunInputForRule.
func RunInput(ctx context.Context, script string) InputResult {
	if script == "" {
		return InputResult{Status: StatusSkipped}
	}

	return runScript(ctx, script, "")
}

// RunInputForRule preprocesses a rule script (resolving `# include:` directives
// against fsys) and runs it. Output is parsed as an Envelope; if expectedID
// is non-empty, the envelope's rule_id must match.
func RunInputForRule(ctx context.Context, fsys fs.FS, scriptPath, script, expectedID string) InputResult {
	if script == "" {
		return InputResult{Status: StatusSkipped}
	}

	resolved, err := Preprocess(fsys, scriptPath, script)
	if err != nil {
		return InputResult{Status: StatusError, Error: err}
	}

	return runScript(ctx, resolved, expectedID)
}

func runScript(ctx context.Context, script, expectedID string) InputResult {
	ctx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sh", "-c", script)

	var stdout, stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return InputResult{
			Status: StatusError,
			Error:  fmt.Errorf("input script failed: %w: %s", err, stderr.String()),
		}
	}

	env, err := ParseEnvelope(stdout.Bytes(), expectedID)
	if err != nil {
		return InputResult{
			Status: StatusError,
			Error:  fmt.Errorf("input script output: %w\noutput: %s", err, stdout.String()),
		}
	}

	switch env.Status {
	case "skip":
		return InputResult{Status: StatusSkipped, SkipReason: env.SkipReason}
	case "ok":
		return InputResult{Status: StatusOK, Data: env.Data}
	default:
		return InputResult{Status: StatusError, Error: fmt.Errorf("unhandled envelope status: %q", env.Status)}
	}
}
