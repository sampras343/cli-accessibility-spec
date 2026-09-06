// internal/report/markdown.go
package report

import (
	"fmt"
	"io"
	"strings"

	"github.com/sampras343/cli-accessibility-spec/internal/engine"
)

// MarkdownReporter renders conformance reports as Markdown.
type MarkdownReporter struct{}

// Name returns the reporter identifier.
func (r *MarkdownReporter) Name() string {
	return "markdown"
}

// FileExtension returns the recommended file extension.
func (r *MarkdownReporter) FileExtension() string {
	return ".md"
}

// Render writes the report as Markdown following the CLI-ACR template.
func (r *MarkdownReporter) Render(report *engine.ConformanceReport, w io.Writer) error {
	// Header
	fmt.Fprintf(w, "# CLI Accessibility Conformance Report\n\n")

	// Product Information section
	fmt.Fprintf(w, "## Product Information\n\n")
	fmt.Fprintf(w, "| Field | Value |\n")
	fmt.Fprintf(w, "|-------|-------|\n")
	fmt.Fprintf(w, "| Product Name | %s |\n", report.Product.Name)
	fmt.Fprintf(w, "| Product Version | %s |\n", report.Product.Version)
	if report.Product.Path != "" {
		fmt.Fprintf(w, "| Product Path | %s |\n", report.Product.Path)
	}
	fmt.Fprintf(w, "| CLI-ACS Version | %s |\n", report.CLIACSVersion)
	fmt.Fprintf(w, "| Suite Version | %s |\n", report.SuiteVersion)
	fmt.Fprintf(w, "| Report Date | %s |\n", report.ReportDate.Format("2006-01-02 15:04:05 MST"))
	fmt.Fprintf(w, "| Overall Level | %s |\n", report.OverallLevel)
	fmt.Fprintf(w, "| Threshold | %s |\n\n", report.Threshold)

	// Conformance Summary
	fmt.Fprintf(w, "## Conformance Summary\n\n")
	if len(report.DomainSummaries) > 0 {
		fmt.Fprintf(w, "| Domain | Total | Pass | Partial | Fail | N/A | Not Evaluated |\n")
		fmt.Fprintf(w, "|--------|-------|------|---------|------|-----|---------------|\n")
		for _, ds := range report.DomainSummaries {
			fmt.Fprintf(w, "| %s | %d | %d | %d | %d | %d | %d |\n",
				ds.Domain, ds.Total, ds.Pass, ds.Partial, ds.Fail, ds.NA, ds.NotEval)
		}
		fmt.Fprintf(w, "\n")
	} else {
		fmt.Fprintf(w, "No test results available.\n\n")
	}

	// Domain Results
	if len(report.Results) > 0 {
		fmt.Fprintf(w, "## Domain Results\n\n")

		// Group results by domain
		domainResults := make(map[string][]engine.Result)
		for _, result := range report.Results {
			domainResults[result.Domain] = append(domainResults[result.Domain], result)
		}

		for domain, results := range domainResults {
			fmt.Fprintf(w, "### Domain: %s\n\n", domain)

			for _, result := range results {
				// Criterion header with outcome
				outcomeSymbol := r.outcomeSymbol(result.Outcome)
				fmt.Fprintf(w, "#### %s %s - %s (Level %s)\n\n",
					outcomeSymbol, result.ID, result.Name, result.Level.String())

				// Details table
				fmt.Fprintf(w, "| Field | Value |\n")
				fmt.Fprintf(w, "|-------|-------|\n")
				fmt.Fprintf(w, "| Outcome | %s |\n", result.Outcome.String())
				fmt.Fprintf(w, "| Testability | %s |\n", result.Testability.String())
				fmt.Fprintf(w, "| Spec Version | %s |\n", result.SpecVersion)

				if result.Remarks != "" {
					fmt.Fprintf(w, "| Remarks | %s |\n", result.Remarks)
				}

				if result.NeedsReview {
					fmt.Fprintf(w, "| Manual Review Required | Yes |\n")
				}

				fmt.Fprintf(w, "\n")

				// Evidence (if present)
				if len(result.Evidence) > 0 {
					fmt.Fprintf(w, "**Evidence:**\n\n")
					for i, ev := range result.Evidence {
						fmt.Fprintf(w, "Evidence %d:\n", i+1)
						fmt.Fprintf(w, "- Command: `%s`\n", ev.Command)
						if ev.Note != "" {
							fmt.Fprintf(w, "- Note: %s\n", ev.Note)
						}
						fmt.Fprintf(w, "- Exit Code: %d\n", ev.ExitCode)
						fmt.Fprintf(w, "- Duration: %s\n", ev.Duration)

						if ev.Stdout != "" {
							fmt.Fprintf(w, "- Stdout:\n  ```\n  %s\n  ```\n", strings.TrimSpace(ev.Stdout))
						}
						if ev.Stderr != "" {
							fmt.Fprintf(w, "- Stderr:\n  ```\n  %s\n  ```\n", strings.TrimSpace(ev.Stderr))
						}
						fmt.Fprintf(w, "\n")
					}
				}
			}
		}
	}

	// Testing Environment
	fmt.Fprintf(w, "## Testing Environment\n\n")
	fmt.Fprintf(w, "| Field | Value |\n")
	fmt.Fprintf(w, "|-------|-------|\n")
	fmt.Fprintf(w, "| Operating System | %s |\n", report.Environment.OS)
	fmt.Fprintf(w, "| Architecture | %s |\n", report.Environment.Arch)
	if report.Environment.Terminal != "" {
		fmt.Fprintf(w, "| Terminal | %s |\n", report.Environment.Terminal)
	}
	if report.Environment.Shell != "" {
		fmt.Fprintf(w, "| Shell | %s |\n", report.Environment.Shell)
	}
	fmt.Fprintf(w, "\n")

	return nil
}

// outcomeSymbol returns a visual symbol for the outcome.
func (r *MarkdownReporter) outcomeSymbol(outcome engine.Outcome) string {
	switch outcome {
	case engine.Supports:
		return "✓"
	case engine.PartiallySupports:
		return "⚠"
	case engine.DoesNotSupport, engine.OutcomeError, engine.OutcomeTimeout:
		return "✗"
	case engine.NotApplicable:
		return "N/A"
	case engine.NotEvaluated:
		return "⊘"
	default:
		return "?"
	}
}
