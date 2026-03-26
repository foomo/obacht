package reporter

import (
	"encoding/json"
	"io"

	"github.com/franklinkim/bouncer/pkg/schema"
)

// JSONReporter outputs scan results as formatted JSON.
type JSONReporter struct{}

// Report writes the ScanResult as indented JSON to w.
func (r *JSONReporter) Report(w io.Writer, result *schema.ScanResult) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(result)
}
