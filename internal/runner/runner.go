package runner

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"
)

// DefaultTimeout is the maximum time a single input script may run.
const DefaultTimeout = 30 * time.Second

// InputResult holds the outcome of running a single input script.
type InputResult struct {
	// Data is the parsed JSON output from the script.
	Data any
	// Status indicates whether the script succeeded, was skipped, or errored.
	Status Status
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
func RunInput(ctx context.Context, script string) InputResult {
	if script == "" {
		return InputResult{Status: StatusSkipped}
	}

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

	var data any
	if err := json.Unmarshal(stdout.Bytes(), &data); err != nil {
		return InputResult{
			Status: StatusError,
			Error:  fmt.Errorf("input script output is not valid JSON: %w\noutput: %s", err, stdout.String()),
		}
	}

	return InputResult{
		Data:   data,
		Status: StatusOK,
	}
}
