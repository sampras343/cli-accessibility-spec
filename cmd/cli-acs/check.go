// cmd/cli-acs/check.go
package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sampras343/cli-accessibility-spec/internal/domains"
	_ "github.com/sampras343/cli-accessibility-spec/internal/domains/color"
	_ "github.com/sampras343/cli-accessibility-spec/internal/domains/help"
	"github.com/sampras343/cli-accessibility-spec/internal/engine"
	"github.com/sampras343/cli-accessibility-spec/internal/report"
	"github.com/spf13/cobra"
)

type checkFlags struct {
	domain         string
	level          string
	autoOnly       bool
	threshold      string
	format         string
	output         string
	crosswalk      bool
	plain          bool
	quiet          bool
	skip           string
	timeout        string
	maxSubcommands int
	criteriaDir    string
}

// domainAdapter adapts domains.Domain to engine.Domain
type domainAdapter struct {
	domains.Domain
}

func (d *domainAdapter) Criteria() []engine.Criterion {
	domainCriteria := d.Domain.Criteria()
	result := make([]engine.Criterion, len(domainCriteria))
	for i, c := range domainCriteria {
		result[i] = c
	}
	return result
}

func newCheckCommand() *cobra.Command {
	flags := &checkFlags{}

	cmd := &cobra.Command{
		Use:   "check <binary>",
		Short: "Evaluate a CLI binary for conformance",
		Long: `Evaluate a CLI binary against CLI-ACS v1.0 criteria.

The check command probes the binary, discovers its capabilities,
and runs conformance tests across all applicable domains.

By default, all domains, levels, and testability types are included.
Use filters to focus on specific criteria or domains.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			binaryPath := args[0]
			return runCheck(binaryPath, flags)
		},
	}

	// Add all flags from the design doc Section 2
	cmd.Flags().StringVar(&flags.domain, "domain", "", "Comma-separated domain filter (output,color,help,errors,input,environ,timing,i18n,lifecycle)")
	cmd.Flags().StringVar(&flags.level, "level", "", "Filter criteria by level (A, AA, AAA)")
	cmd.Flags().BoolVar(&flags.autoOnly, "auto-only", false, "Run only AUTO criteria, skip SEMI and MANUAL")
	cmd.Flags().StringVar(&flags.threshold, "threshold", "A", "Conformance level required for exit 0 (A, AA, AAA)")
	cmd.Flags().StringVar(&flags.format, "format", "terminal", "Output format: terminal, json, markdown, html")
	cmd.Flags().StringVar(&flags.output, "output", "", "Write report to file instead of stdout")
	cmd.Flags().BoolVar(&flags.crosswalk, "crosswalk", false, "Include WCAG/508/EN 301 549 mapping in report")
	cmd.Flags().BoolVar(&flags.plain, "plain", false, "Screen-reader-friendly output (no tables, no box-drawing, no color)")
	cmd.Flags().BoolVar(&flags.quiet, "quiet", false, "Suppress all output, communicate only via exit code")
	cmd.Flags().StringVar(&flags.skip, "skip", "", "Comma-separated criterion IDs to skip")
	cmd.Flags().StringVar(&flags.timeout, "timeout", "10s", "Per-criterion timeout")
	cmd.Flags().IntVar(&flags.maxSubcommands, "max-subcommands", 10, "Maximum subcommands to test during discovery")
	cmd.Flags().StringVar(&flags.criteriaDir, "criteria-dir", "", "Additional directory of YAML criteria to load")

	return cmd
}

func runCheck(binaryPath string, flags *checkFlags) error {
	// Parse timeout
	timeout, err := time.ParseDuration(flags.timeout)
	if err != nil {
		return fmt.Errorf("invalid timeout: %w", err)
	}

	// Parse threshold level
	threshold, err := parseLevel(flags.threshold)
	if err != nil {
		return fmt.Errorf("invalid threshold: %w", err)
	}

	// Parse level filter
	var levelFilter []engine.Level
	if flags.level != "" {
		level, err := parseLevel(flags.level)
		if err != nil {
			return fmt.Errorf("invalid level filter: %w", err)
		}
		levelFilter = []engine.Level{level}
	}

	// Parse domain filter
	var domainFilter []string
	if flags.domain != "" {
		domainFilter = splitAndTrim(flags.domain)
	}

	// Parse skip criteria
	var skipCriteria []string
	if flags.skip != "" {
		skipCriteria = splitAndTrim(flags.skip)
	}

	// Apply noColor or plain flags to config
	noColorFlag := noColor || flags.plain

	// Build engine config
	cfg := engine.Config{
		DomainFilter:   domainFilter,
		LevelFilter:    levelFilter,
		AutoOnly:       flags.autoOnly,
		Threshold:      threshold,
		SkipCriteria:   skipCriteria,
		Timeout:        timeout,
		MaxSubcommands: flags.maxSubcommands,
		CriteriaDir:    flags.criteriaDir,
		Format:         flags.format,
		OutputPath:     flags.output,
		Quiet:          flags.quiet,
		Plain:          flags.plain,
		Crosswalk:      flags.crosswalk,
		NoColor:        noColorFlag,
	}

	// Create domain provider that wraps domains.AllDomains()
	domainProvider := func() []engine.Domain {
		allDomains := domains.AllDomains()
		// Adapt []domains.Domain to []engine.Domain
		result := make([]engine.Domain, len(allDomains))
		for i, d := range allDomains {
			result[i] = &domainAdapter{d}
		}
		return result
	}

	// Create runner
	runner := engine.NewRunner(cfg, domainProvider)

	// Run conformance check
	ctx := context.Background()
	conformanceReport, err := runner.Run(ctx, binaryPath)
	if err != nil {
		return fmt.Errorf("conformance check failed: %w", err)
	}

	// Select reporter
	reporter, err := report.GetReporter(flags.format)
	if err != nil {
		return fmt.Errorf("failed to get reporter: %w", err)
	}

	// Determine output destination
	var output *os.File
	if flags.output != "" {
		output, err = os.Create(flags.output)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer output.Close()
	} else {
		output = os.Stdout
	}

	// Render report (unless quiet mode)
	if !flags.quiet {
		if err := reporter.Render(conformanceReport, output); err != nil {
			return fmt.Errorf("failed to render report: %w", err)
		}
	}

	// Compute and exit with appropriate code
	exitCode := engine.ComputeExitCode(conformanceReport, threshold)
	os.Exit(exitCode)

	return nil
}

// parseLevel converts a string to an engine.Level
func parseLevel(s string) (engine.Level, error) {
	switch strings.ToUpper(s) {
	case "A":
		return engine.LevelA, nil
	case "AA":
		return engine.LevelAA, nil
	case "AAA":
		return engine.LevelAAA, nil
	default:
		return engine.LevelA, fmt.Errorf("invalid level %q (must be A, AA, or AAA)", s)
	}
}

// splitAndTrim splits a comma-separated string and trims whitespace
func splitAndTrim(s string) []string {
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
