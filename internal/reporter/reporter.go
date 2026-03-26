package reporter

import (
	"io"

	"github.com/franklinkim/bouncer/pkg/schema"
)

// Reporter writes scan results to an output.
type Reporter interface {
	Report(w io.Writer, result *schema.ScanResult) error
}

// ForFormat returns the appropriate reporter for the given format.
func ForFormat(format string) Reporter {
	switch format {
	case "json":
		return &JSONReporter{}
	case "pretty":
		return &PrettyReporter{}
	default:
		return &PrettyReporter{}
	}
}
