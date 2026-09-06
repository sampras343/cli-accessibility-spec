// internal/report/terminal_test.go
package report

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/sampras343/cli-accessibility-spec/internal/engine"
)

func TestTerminalReporterBasic(t *testing.T) {
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
	r := &TerminalReporter{}

	if err := r.Render(report, &buf); err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	output := buf.String()

	// Check header
	if !strings.Contains(output, "test-tool") {
		t.Error("output missing product name")
	}
	if !strings.Contains(output, "1.0") {
		t.Error("output missing version")
	}

	// Check results with prefixes
	if !strings.Contains(output, "[PASS]") {
		t.Error("output missing [PASS] prefix")
	}
	if !strings.Contains(output, "[FAIL]") {
		t.Error("output missing [FAIL] prefix")
	}

	// Check remarks are shown for failures
	if !strings.Contains(output, "Failed to accept keyboard input") {
		t.Error("output missing failure remarks")
	}

	// Check domain summaries
	if !strings.Contains(output, "color") || !strings.Contains(output, "input") {
		t.Error("output missing domain summaries")
	}
}

func TestTerminalReporterPlainMode(t *testing.T) {
	report := &engine.ConformanceReport{
		CLIACSVersion: "1.0",
		SuiteVersion:  "0.1.0",
		Product:       engine.ProductInfo{Name: "test", Version: "1.0"},
		ReportDate:    time.Now(),
		Results: []engine.Result{
			{
				ID:          "CV-2",
				Name:        "Test",
				Domain:      "color",
				Level:       engine.LevelA,
				Outcome:     engine.Supports,
				SpecVersion: "1.0",
			},
		},
	}
	report.ComputeSummaries()

	var buf bytes.Buffer
	r := &TerminalReporter{Plain: true}

	if err := r.Render(report, &buf); err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	output := buf.String()

	// In plain mode, should not have ANSI codes
	if strings.Contains(output, "\x1b[") {
		t.Error("plain mode output contains ANSI escape codes")
	}

	// Should still have prefixes
	if !strings.Contains(output, "[PASS]") {
		t.Error("plain mode output missing [PASS] prefix")
	}
}

func TestTerminalReporterNOCOLOR(t *testing.T) {
	// Save original and restore after test
	originalNoColor := os.Getenv("NO_COLOR")
	defer func() {
		if originalNoColor != "" {
			os.Setenv("NO_COLOR", originalNoColor)
		} else {
			os.Unsetenv("NO_COLOR")
		}
	}()

	os.Setenv("NO_COLOR", "1")

	report := &engine.ConformanceReport{
		CLIACSVersion: "1.0",
		SuiteVersion:  "0.1.0",
		Product:       engine.ProductInfo{Name: "test", Version: "1.0"},
		ReportDate:    time.Now(),
		Results: []engine.Result{
			{
				ID:          "CV-2",
				Name:        "Test",
				Domain:      "color",
				Level:       engine.LevelA,
				Outcome:     engine.Supports,
				SpecVersion: "1.0",
			},
		},
	}
	report.ComputeSummaries()

	var buf bytes.Buffer
	r := &TerminalReporter{}

	if err := r.Render(report, &buf); err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	output := buf.String()

	// Should not have ANSI codes when NO_COLOR is set
	if strings.Contains(output, "\x1b[") {
		t.Error("NO_COLOR set but output contains ANSI escape codes")
	}
}

func TestTerminalReporterDumbTerminal(t *testing.T) {
	// Save original and restore after test
	originalTerm := os.Getenv("TERM")
	defer func() {
		if originalTerm != "" {
			os.Setenv("TERM", originalTerm)
		} else {
			os.Unsetenv("TERM")
		}
	}()

	os.Setenv("TERM", "dumb")

	report := &engine.ConformanceReport{
		CLIACSVersion: "1.0",
		SuiteVersion:  "0.1.0",
		Product:       engine.ProductInfo{Name: "test", Version: "1.0"},
		ReportDate:    time.Now(),
		Results: []engine.Result{
			{
				ID:          "CV-2",
				Name:        "Test",
				Domain:      "color",
				Level:       engine.LevelA,
				Outcome:     engine.Supports,
				SpecVersion: "1.0",
			},
		},
	}
	report.ComputeSummaries()

	var buf bytes.Buffer
	r := &TerminalReporter{}

	if err := r.Render(report, &buf); err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	output := buf.String()

	// Should not have ANSI codes when TERM=dumb
	if strings.Contains(output, "\x1b[") {
		t.Error("TERM=dumb but output contains ANSI escape codes")
	}
}

func TestTerminalReporterAllOutcomes(t *testing.T) {
	report := &engine.ConformanceReport{
		CLIACSVersion: "1.0",
		SuiteVersion:  "0.1.0",
		Product:       engine.ProductInfo{Name: "test", Version: "1.0"},
		ReportDate:    time.Now(),
		Results: []engine.Result{
			{ID: "1", Name: "Test 1", Domain: "d", Level: engine.LevelA, Outcome: engine.Supports, SpecVersion: "1.0"},
			{ID: "2", Name: "Test 2", Domain: "d", Level: engine.LevelA, Outcome: engine.PartiallySupports, SpecVersion: "1.0"},
			{ID: "3", Name: "Test 3", Domain: "d", Level: engine.LevelA, Outcome: engine.DoesNotSupport, SpecVersion: "1.0"},
			{ID: "4", Name: "Test 4", Domain: "d", Level: engine.LevelA, Outcome: engine.NotApplicable, SpecVersion: "1.0"},
			{ID: "5", Name: "Test 5", Domain: "d", Level: engine.LevelA, Outcome: engine.NotEvaluated, SpecVersion: "1.0"},
		},
	}
	report.ComputeSummaries()

	var buf bytes.Buffer
	r := &TerminalReporter{Plain: true}

	if err := r.Render(report, &buf); err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	output := buf.String()

	// Check all outcome prefixes
	if !strings.Contains(output, "[PASS]") {
		t.Error("missing [PASS] for Supports outcome")
	}
	if !strings.Contains(output, "[PARTIAL]") {
		t.Error("missing [PARTIAL] for PartiallySupports outcome")
	}
	if !strings.Contains(output, "[FAIL]") {
		t.Error("missing [FAIL] for DoesNotSupport outcome")
	}
	if !strings.Contains(output, "[N/A]") {
		t.Error("missing [N/A] for NotApplicable outcome")
	}
	if !strings.Contains(output, "[SKIP]") {
		t.Error("missing [SKIP] for NotEvaluated outcome")
	}
}

func TestTerminalReporterName(t *testing.T) {
	r := &TerminalReporter{}
	if r.Name() != "terminal" {
		t.Errorf("expected name=terminal, got %s", r.Name())
	}
}

func TestTerminalReporterExtension(t *testing.T) {
	r := &TerminalReporter{}
	if r.FileExtension() != ".txt" {
		t.Errorf("expected extension=.txt, got %s", r.FileExtension())
	}
}
