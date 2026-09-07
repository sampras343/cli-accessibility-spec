// internal/report/markdown_test.go
package report

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/sampras343/cli-accessibility-spec/internal/engine"
)

func TestMarkdownReporterName(t *testing.T) {
	r := &MarkdownReporter{}
	if r.Name() != "markdown" {
		t.Errorf("expected name=markdown, got %s", r.Name())
	}
}

func TestMarkdownReporterExtension(t *testing.T) {
	r := &MarkdownReporter{}
	if r.FileExtension() != ".md" {
		t.Errorf("expected extension=.md, got %s", r.FileExtension())
	}
}

func TestMarkdownReporterOutput(t *testing.T) {
	report := &engine.ConformanceReport{
		CLIACSVersion: "1.0",
		SuiteVersion:  "0.1.0",
		Product: engine.ProductInfo{
			Name:    "test-tool",
			Version: "1.0",
			Path:    "/usr/bin/test-tool",
		},
		ReportDate: time.Date(2026, 9, 6, 12, 0, 0, 0, time.UTC),
		Environment: engine.TestEnvironment{
			OS:       "linux",
			Arch:     "amd64",
			Terminal: "xterm-256color",
			Shell:    "bash",
			SuiteVer: "0.1.0",
		},
		Results: []engine.Result{
			{
				ID:          "CV-2",
				Name:        "Color is not the only indicator",
				Domain:      "color",
				Level:       engine.LevelA,
				Testability: engine.Manual,
				Outcome:     engine.Supports,
				SpecVersion: "1.0",
			},
			{
				ID:          "IV-1",
				Name:        "Keyboard input works",
				Domain:      "input",
				Level:       engine.LevelA,
				Testability: engine.Auto,
				Outcome:     engine.DoesNotSupport,
				Remarks:     "Failed to accept keyboard input",
				SpecVersion: "1.0",
			},
		},
		OverallLevel: "A",
		Threshold:    "AA",
	}
	report.ComputeSummaries()

	var buf bytes.Buffer
	r := &MarkdownReporter{}

	if err := r.Render(report, &buf); err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	output := buf.String()

	// Check for expected headers
	expectedHeaders := []string{
		"# CLI Accessibility Conformance Report",
		"## Product Information",
		"## Conformance Summary",
		"## Domain Results",
		"### Domain: color",
		"### Domain: input",
		"## Testing Environment",
	}

	for _, header := range expectedHeaders {
		if !strings.Contains(output, header) {
			t.Errorf("expected header %q not found in output", header)
		}
	}

	// Check for product information
	if !strings.Contains(output, "test-tool") {
		t.Error("product name not found in output")
	}
	if !strings.Contains(output, "1.0") {
		t.Error("product version not found in output")
	}

	// Check for results
	if !strings.Contains(output, "CV-2") {
		t.Error("criterion ID CV-2 not found in output")
	}
	if !strings.Contains(output, "IV-1") {
		t.Error("criterion ID IV-1 not found in output")
	}

	// Check for remarks
	if !strings.Contains(output, "Failed to accept keyboard input") {
		t.Error("remarks not found in output")
	}

	// Check for environment
	if !strings.Contains(output, "linux") {
		t.Error("OS not found in output")
	}
	if !strings.Contains(output, "amd64") {
		t.Error("architecture not found in output")
	}
}

func TestMarkdownReporterEmptyResults(t *testing.T) {
	report := &engine.ConformanceReport{
		CLIACSVersion: "1.0",
		SuiteVersion:  "0.1.0",
		Product:       engine.ProductInfo{Name: "empty", Version: "1.0"},
		ReportDate:    time.Now(),
		Results:       []engine.Result{},
	}
	report.ComputeSummaries()

	var buf bytes.Buffer
	r := &MarkdownReporter{}

	if err := r.Render(report, &buf); err != nil {
		t.Fatalf("Render failed for empty results: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "# CLI Accessibility Conformance Report") {
		t.Error("main header not found for empty report")
	}
}
