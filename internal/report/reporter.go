// internal/report/reporter.go
package report

import (
	"fmt"
	"io"

	"github.com/sampras343/cli-accessibility-spec/internal/engine"
)

// Reporter defines the interface for rendering conformance reports.
type Reporter interface {
	// Name returns the reporter's identifier (e.g., "json", "terminal").
	Name() string

	// FileExtension returns the recommended file extension (e.g., ".json", ".txt").
	FileExtension() string

	// Render writes the conformance report to the provided writer.
	Render(report *engine.ConformanceReport, w io.Writer) error
}

// registry holds all registered reporters.
var registry = make(map[string]Reporter)

// Register adds a reporter to the registry.
func Register(r Reporter) {
	registry[r.Name()] = r
}

// GetReporter retrieves a reporter by name.
func GetReporter(name string) (Reporter, error) {
	r, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("unknown reporter: %s", name)
	}
	return r, nil
}

// ListReporters returns the names of all registered reporters.
func ListReporters() []string {
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	return names
}

func init() {
	// Register built-in reporters
	Register(&JSONReporter{})
	Register(&TerminalReporter{})
}
