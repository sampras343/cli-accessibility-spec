// internal/report/example_test.go
package report

import (
	"os"
	"time"

	"github.com/sampras343/cli-accessibility-spec/internal/engine"
)

// This example demonstrates the terminal reporter output.
func ExampleTerminalReporter() {
	report := &engine.ConformanceReport{
		CLIACSVersion: "1.0",
		SuiteVersion:  "0.1.0",
		Product: engine.ProductInfo{
			Name:    "mytool",
			Version: "2.0",
			Path:    "/usr/local/bin/mytool",
		},
		ReportDate: time.Date(2026, 9, 6, 12, 0, 0, 0, time.UTC),
		Environment: engine.TestEnvironment{
			OS:       "linux",
			Arch:     "amd64",
			Terminal: "xterm-256color",
			Shell:    "bash",
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
				Remarks:     "Failed to accept keyboard input in non-canonical mode",
				SpecVersion: "1.0",
			},
			{
				ID:          "OV-3",
				Name:        "Progress indicators are accessible",
				Domain:      "output",
				Level:       engine.LevelAA,
				Testability: engine.Semi,
				Outcome:     engine.PartiallySupports,
				Remarks:     "Progress bar visible but lacks text alternative",
				SpecVersion: "1.0",
			},
		},
		OverallLevel: "A",
		Threshold:    "AA",
	}
	report.ComputeSummaries()

	// Use plain mode for consistent test output
	reporter := &TerminalReporter{Plain: true}
	reporter.Render(report, os.Stdout)

	// Output is tested in terminal_test.go
}

// This example demonstrates the JSON reporter output.
func ExampleJSONReporter() {
	report := &engine.ConformanceReport{
		CLIACSVersion: "1.0",
		SuiteVersion:  "0.1.0",
		Product: engine.ProductInfo{
			Name:    "mytool",
			Version: "2.0",
		},
		ReportDate: time.Date(2026, 9, 6, 12, 0, 0, 0, time.UTC),
		Results: []engine.Result{
			{
				ID:          "CV-2",
				Name:        "Color is not the only indicator",
				Domain:      "color",
				Level:       engine.LevelA,
				Outcome:     engine.Supports,
				SpecVersion: "1.0",
			},
		},
	}
	report.ComputeSummaries()

	reporter := &JSONReporter{}
	reporter.Render(report, os.Stdout)

	// Output is tested in json_test.go
}
