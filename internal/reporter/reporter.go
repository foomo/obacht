package reporter

import (
	"io"

	"github.com/foomo/obacht/pkg/schema"
)

// Reporter writes scan results to an output.
type Reporter interface {
	Report(w io.Writer, result *schema.ScanResult) error
}

// ForFormat returns the appropriate reporter for the given format.
// showPassing controls whether passing checks are included in pretty output;
// it has no effect on JSON output.
func ForFormat(format string, showPassing bool) Reporter {
	switch format {
	case "json":
		return &JSONReporter{}
	case "pretty":
		return &PrettyReporter{ShowPassing: showPassing}
	default:
		return &PrettyReporter{ShowPassing: showPassing}
	}
}
