// internal/report/terminal.go
package report

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sampras343/cli-accessibility-spec/internal/engine"
)

// TerminalReporter renders conformance reports for terminal display.
type TerminalReporter struct {
	// Plain disables colors and fancy formatting for screen readers.
	Plain bool
}

// Name returns the reporter identifier.
func (r *TerminalReporter) Name() string {
	return "terminal"
}

// FileExtension returns the recommended file extension.
func (r *TerminalReporter) FileExtension() string {
	return ".txt"
}

// ANSI color codes
const (
	colorReset  = "\x1b[0m"
	colorRed    = "\x1b[31m"
	colorGreen  = "\x1b[32m"
	colorYellow = "\x1b[33m"
	colorBold   = "\x1b[1m"
)

// shouldUseColor determines if color output should be used.
func (r *TerminalReporter) shouldUseColor() bool {
	if r.Plain {
		return false
	}
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	term := os.Getenv("TERM")
	if term == "dumb" || term == "" {
		return false
	}
	return true
}

// colorize wraps text in ANSI color codes if colors are enabled.
func (r *TerminalReporter) colorize(text, color string) string {
	if !r.shouldUseColor() {
		return text
	}
	return color + text + colorReset
}

// outcomePrefix returns the text prefix for an outcome.
func (r *TerminalReporter) outcomePrefix(outcome engine.Outcome) string {
	switch outcome {
	case engine.Supports:
		return r.colorize("[PASS]", colorGreen)
	case engine.PartiallySupports:
		return r.colorize("[PARTIAL]", colorYellow)
	case engine.DoesNotSupport, engine.OutcomeError, engine.OutcomeTimeout:
		return r.colorize("[FAIL]", colorRed)
	case engine.NotApplicable:
		return "[N/A]"
	case engine.NotEvaluated:
		return "[SKIP]"
	default:
		return "[UNKNOWN]"
	}
}

// Render writes the report to the terminal.
func (r *TerminalReporter) Render(report *engine.ConformanceReport, w io.Writer) error {
	// Header
	fmt.Fprintf(w, "%s\n", r.colorize("CLI Accessibility Conformance Suite Report", colorBold))
	fmt.Fprintf(w, "%s\n\n", strings.Repeat("=", 60))

	fmt.Fprintf(w, "Product:        %s %s\n", report.Product.Name, report.Product.Version)
	if report.Product.Path != "" {
		fmt.Fprintf(w, "Path:           %s\n", report.Product.Path)
	}
	fmt.Fprintf(w, "CLI-ACS:        %s\n", report.CLIACSVersion)
	fmt.Fprintf(w, "Suite Version:  %s\n", report.SuiteVersion)
	fmt.Fprintf(w, "Report Date:    %s\n", report.ReportDate.Format("2006-01-02 15:04:05 MST"))
	fmt.Fprintf(w, "Overall Level:  %s\n", report.OverallLevel)
	fmt.Fprintf(w, "Threshold:      %s\n", report.Threshold)

	if report.Environment.OS != "" {
		fmt.Fprintf(w, "\nEnvironment:\n")
		fmt.Fprintf(w, "  OS:           %s\n", report.Environment.OS)
		fmt.Fprintf(w, "  Architecture: %s\n", report.Environment.Arch)
		if report.Environment.Terminal != "" {
			fmt.Fprintf(w, "  Terminal:     %s\n", report.Environment.Terminal)
		}
		if report.Environment.Shell != "" {
			fmt.Fprintf(w, "  Shell:        %s\n", report.Environment.Shell)
		}
	}

	// Domain summaries
	if len(report.DomainSummaries) > 0 {
		fmt.Fprintf(w, "\n%s\n", r.colorize("Domain Summaries", colorBold))
		fmt.Fprintf(w, "%s\n\n", strings.Repeat("-", 60))

		if r.Plain {
			// Plain text table
			for _, ds := range report.DomainSummaries {
				fmt.Fprintf(w, "Domain: %s\n", ds.Domain)
				fmt.Fprintf(w, "  Total:        %d\n", ds.Total)
				fmt.Fprintf(w, "  Pass:         %d\n", ds.Pass)
				fmt.Fprintf(w, "  Partial:      %d\n", ds.Partial)
				fmt.Fprintf(w, "  Fail:         %d\n", ds.Fail)
				fmt.Fprintf(w, "  N/A:          %d\n", ds.NA)
				fmt.Fprintf(w, "  Not Eval:     %d\n", ds.NotEval)
				fmt.Fprintf(w, "\n")
			}
		} else {
			// Formatted table
			fmt.Fprintf(w, "%-15s %6s %6s %8s %6s %6s %9s\n",
				"Domain", "Total", "Pass", "Partial", "Fail", "N/A", "Not Eval")
			fmt.Fprintf(w, "%s\n", strings.Repeat("-", 60))

			for _, ds := range report.DomainSummaries {
				fmt.Fprintf(w, "%-15s %6d %6d %8d %6d %6d %9d\n",
					ds.Domain, ds.Total, ds.Pass, ds.Partial, ds.Fail, ds.NA, ds.NotEval)
			}
			fmt.Fprintf(w, "\n")
		}
	}

	// Results
	if len(report.Results) > 0 {
		fmt.Fprintf(w, "%s\n", r.colorize("Test Results", colorBold))
		fmt.Fprintf(w, "%s\n\n", strings.Repeat("-", 60))

		// Group by domain for better readability
		domainResults := make(map[string][]engine.Result)
		for _, result := range report.Results {
			domainResults[result.Domain] = append(domainResults[result.Domain], result)
		}

		for domain, results := range domainResults {
			fmt.Fprintf(w, "%s\n", r.colorize(fmt.Sprintf("Domain: %s", domain), colorBold))

			for _, result := range results {
				prefix := r.outcomePrefix(result.Outcome)
				fmt.Fprintf(w, "  %s %s (%s) - %s\n",
					prefix, result.ID, result.Level.String(), result.Name)

				// Show remarks for non-passing results
				if result.Remarks != "" && result.IsFailure() {
					fmt.Fprintf(w, "      Remarks: %s\n", result.Remarks)
				}

				// Show needs review flag
				if result.NeedsReview {
					fmt.Fprintf(w, "      %s\n", r.colorize("(Needs Manual Review)", colorYellow))
				}
			}
			fmt.Fprintf(w, "\n")
		}
	}

	// Summary footer
	totalTests := len(report.Results)
	passCount := 0
	failCount := 0
	partialCount := 0

	for _, result := range report.Results {
		switch result.Outcome {
		case engine.Supports:
			passCount++
		case engine.PartiallySupports:
			partialCount++
		case engine.DoesNotSupport, engine.OutcomeError, engine.OutcomeTimeout:
			failCount++
		}
	}

	fmt.Fprintf(w, "%s\n", strings.Repeat("=", 60))
	fmt.Fprintf(w, "Total Tests: %d | ", totalTests)
	fmt.Fprintf(w, "%s: %d | ", r.colorize("Pass", colorGreen), passCount)
	fmt.Fprintf(w, "%s: %d | ", r.colorize("Partial", colorYellow), partialCount)
	fmt.Fprintf(w, "%s: %d\n", r.colorize("Fail", colorRed), failCount)

	return nil
}
