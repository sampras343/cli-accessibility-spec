// internal/report/reporter_test.go
package report

import (
	"io"
	"testing"

	"github.com/sampras343/cli-accessibility-spec/internal/engine"
)

// customReporter is a test reporter for testing registration
type customReporter struct{}

func (r *customReporter) Name() string {
	return "custom"
}

func (r *customReporter) FileExtension() string {
	return ".custom"
}

func (r *customReporter) Render(report *engine.ConformanceReport, w io.Writer) error {
	return nil
}

func TestGetReporter(t *testing.T) {
	tests := []struct {
		name      string
		wantError bool
	}{
		{"json", false},
		{"terminal", false},
		{"unknown", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := GetReporter(tt.name)
			if tt.wantError {
				if err == nil {
					t.Errorf("GetReporter(%s) expected error, got nil", tt.name)
				}
			} else {
				if err != nil {
					t.Errorf("GetReporter(%s) unexpected error: %v", tt.name, err)
				}
				if r == nil {
					t.Errorf("GetReporter(%s) returned nil reporter", tt.name)
				}
				if r.Name() != tt.name {
					t.Errorf("GetReporter(%s) returned reporter with name %s", tt.name, r.Name())
				}
			}
		})
	}
}

func TestListReporters(t *testing.T) {
	reporters := ListReporters()
	if len(reporters) < 2 {
		t.Errorf("expected at least 2 reporters, got %d", len(reporters))
	}

	// Check that json and terminal are in the list
	hasJSON := false
	hasTerminal := false
	for _, name := range reporters {
		if name == "json" {
			hasJSON = true
		}
		if name == "terminal" {
			hasTerminal = true
		}
	}

	if !hasJSON {
		t.Error("ListReporters() missing 'json' reporter")
	}
	if !hasTerminal {
		t.Error("ListReporters() missing 'terminal' reporter")
	}
}

func TestRegister(t *testing.T) {
	// Register a custom reporter
	cr := &customReporter{}
	Register(cr)

	// Verify it can be retrieved
	r, err := GetReporter("custom")
	if err != nil {
		t.Fatalf("GetReporter(custom) failed: %v", err)
	}
	if r.Name() != "custom" {
		t.Errorf("expected custom reporter, got %s", r.Name())
	}
}
