// internal/report/json.go
package report

import (
	"encoding/json"
	"io"

	"github.com/sampras343/cli-accessibility-spec/internal/engine"
)

// JSONReporter renders conformance reports as JSON.
type JSONReporter struct{}

// Name returns the reporter identifier.
func (r *JSONReporter) Name() string {
	return "json"
}

// FileExtension returns the recommended file extension.
func (r *JSONReporter) FileExtension() string {
	return ".json"
}

// Render writes the report as formatted JSON.
func (r *JSONReporter) Render(report *engine.ConformanceReport, w io.Writer) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(report)
}
