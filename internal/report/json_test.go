// internal/report/json_test.go
package report

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/sampras343/cli-accessibility-spec/internal/engine"
)

func TestJSONReporterOutput(t *testing.T) {
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
		},
		OverallLevel: "A",
		Threshold:    "AA",
	}
	report.ComputeSummaries()

	var buf bytes.Buffer
	r := &JSONReporter{}

	if err := r.Render(report, &buf); err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	// Verify it's valid JSON
	var parsed map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	// Verify key fields are present
	if parsed["cli_acs_version"] != "1.0" {
		t.Errorf("expected cli_acs_version=1.0, got %v", parsed["cli_acs_version"])
	}

	if parsed["suite_version"] != "0.1.0" {
		t.Errorf("expected suite_version=0.1.0, got %v", parsed["suite_version"])
	}

	product, ok := parsed["product"].(map[string]interface{})
	if !ok {
		t.Fatal("product field missing or not an object")
	}
	if product["name"] != "test-tool" {
		t.Errorf("expected product.name=test-tool, got %v", product["name"])
	}

	results, ok := parsed["results"].([]interface{})
	if !ok {
		t.Fatal("results field missing or not an array")
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}

	summaries, ok := parsed["domain_summaries"].([]interface{})
	if !ok {
		t.Fatal("domain_summaries field missing or not an array")
	}
	if len(summaries) != 1 {
		t.Errorf("expected 1 domain summary, got %d", len(summaries))
	}
}

func TestJSONReporterName(t *testing.T) {
	r := &JSONReporter{}
	if r.Name() != "json" {
		t.Errorf("expected name=json, got %s", r.Name())
	}
}

func TestJSONReporterExtension(t *testing.T) {
	r := &JSONReporter{}
	if r.FileExtension() != ".json" {
		t.Errorf("expected extension=.json, got %s", r.FileExtension())
	}
}

func TestJSONReporterEmptyResults(t *testing.T) {
	report := &engine.ConformanceReport{
		CLIACSVersion: "1.0",
		SuiteVersion:  "0.1.0",
		Product:       engine.ProductInfo{Name: "empty", Version: "1.0"},
		ReportDate:    time.Now(),
		Results:       []engine.Result{},
	}
	report.ComputeSummaries()

	var buf bytes.Buffer
	r := &JSONReporter{}

	if err := r.Render(report, &buf); err != nil {
		t.Fatalf("Render failed for empty results: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}
