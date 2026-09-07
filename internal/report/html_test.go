// internal/report/html_test.go
package report

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/sampras343/cli-accessibility-spec/internal/engine"
)

func TestHTMLReporterName(t *testing.T) {
	r := &HTMLReporter{}
	if r.Name() != "html" {
		t.Errorf("expected name=html, got %s", r.Name())
	}
}

func TestHTMLReporterExtension(t *testing.T) {
	r := &HTMLReporter{}
	if r.FileExtension() != ".html" {
		t.Errorf("expected extension=.html, got %s", r.FileExtension())
	}
}

func TestHTMLReporterOutput(t *testing.T) {
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
				Evidence: []engine.Evidence{
					{
						Command:  "test-tool input",
						Stdout:   "test output",
						ExitCode: 1,
						Duration: 100 * time.Millisecond,
					},
				},
				SpecVersion: "1.0",
			},
		},
		OverallLevel: "A",
		Threshold:    "AA",
	}
	report.ComputeSummaries()

	var buf bytes.Buffer
	r := &HTMLReporter{}

	if err := r.Render(report, &buf); err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	output := buf.String()

	// Check for valid HTML structure
	expectedTags := []string{
		"<!DOCTYPE html>",
		"<html",
		"<head>",
		"<body>",
		"</body>",
		"</html>",
		"<style>",
		"</style>",
	}

	for _, tag := range expectedTags {
		if !strings.Contains(output, tag) {
			t.Errorf("expected HTML tag %q not found in output", tag)
		}
	}

	// Check for product information
	if !strings.Contains(output, "test-tool") {
		t.Error("product name not found in output")
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

	// Check for details tags (expandable sections)
	if !strings.Contains(output, "<details>") {
		t.Error("expandable sections (details) not found in output")
	}

	// Check for embedded CSS (no external dependencies)
	if strings.Contains(output, "href=") && strings.Contains(output, ".css") {
		t.Error("external CSS stylesheet found; expected embedded CSS only")
	}
	if strings.Contains(output, "cdn.") || strings.Contains(output, "unpkg.") {
		t.Error("CDN reference found; expected no external dependencies")
	}
}

func TestHTMLReporterEmptyResults(t *testing.T) {
	report := &engine.ConformanceReport{
		CLIACSVersion: "1.0",
		SuiteVersion:  "0.1.0",
		Product:       engine.ProductInfo{Name: "empty", Version: "1.0"},
		ReportDate:    time.Now(),
		Results:       []engine.Result{},
	}
	report.ComputeSummaries()

	var buf bytes.Buffer
	r := &HTMLReporter{}

	if err := r.Render(report, &buf); err != nil {
		t.Fatalf("Render failed for empty results: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "<!DOCTYPE html>") {
		t.Error("DOCTYPE not found for empty report")
	}
	if !strings.Contains(output, "CLI Accessibility Conformance Report") {
		t.Error("title not found for empty report")
	}
}

func TestHTMLReporterWithEvidence(t *testing.T) {
	report := &engine.ConformanceReport{
		CLIACSVersion: "1.0",
		SuiteVersion:  "0.1.0",
		Product:       engine.ProductInfo{Name: "test-tool", Version: "1.0"},
		ReportDate:    time.Now(),
		Results: []engine.Result{
			{
				ID:          "OV-1",
				Name:        "Test with evidence",
				Domain:      "output",
				Level:       engine.LevelA,
				Testability: engine.Auto,
				Outcome:     engine.Supports,
				Evidence: []engine.Evidence{
					{
						Command:  "test-tool --help",
						Stdout:   "Usage: test-tool",
						Stderr:   "",
						ExitCode: 0,
						Duration: 100 * time.Millisecond,
					},
				},
				SpecVersion: "1.0",
			},
		},
	}
	report.ComputeSummaries()

	var buf bytes.Buffer
	r := &HTMLReporter{}

	if err := r.Render(report, &buf); err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	output := buf.String()

	// Check for evidence section
	if !strings.Contains(output, "test-tool --help") {
		t.Error("evidence command not found in output")
	}
	if !strings.Contains(output, "Usage: test-tool") {
		t.Error("evidence stdout not found in output")
	}
}
