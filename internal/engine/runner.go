// internal/engine/runner.go
package engine

import (
	"context"
	"os"
	"runtime"
	"time"

	"github.com/sampras343/cli-accessibility-spec/internal/probe"
)

// Domain represents a conformance domain with criteria
type Domain interface {
	Name() string
	Criteria() []Criterion
}

// Criterion represents a single conformance criterion
type Criterion interface {
	ID() string
	Name() string
	Domain() string
	Level() Level
	Testability() Testability
	SpecVersion() string
	Precondition(probe *probe.ProbeResult) bool
	Run(ctx context.Context, binary string, probe *probe.ProbeResult) *Result
}

// DomainProvider is a function that returns all available domains
type DomainProvider func() []Domain

// Runner orchestrates the full conformance check flow
type Runner struct {
	config         Config
	domainProvider DomainProvider
}

// NewRunner creates a new conformance runner with the given configuration
func NewRunner(cfg Config, domainProvider DomainProvider) *Runner {
	return &Runner{
		config:         cfg,
		domainProvider: domainProvider,
	}
}

// Run executes the conformance test suite against the given binary
func (r *Runner) Run(ctx context.Context, binary string) (*ConformanceReport, error) {
	// Step 1: Probe the binary to discover its capabilities
	probeResult, err := probe.Probe(ctx, binary, r.config.MaxSubcommands)
	if err != nil {
		return nil, err
	}

	// Step 2: Initialize report
	report := &ConformanceReport{
		CLIACSVersion: "1.0.0",
		SuiteVersion:  "1.0.0",
		Product: ProductInfo{
			Name:    probeResult.BinaryName,
			Path:    probeResult.BinaryPath,
			Version: probeResult.Version,
		},
		ReportDate: time.Now(),
		Environment: TestEnvironment{
			OS:       runtime.GOOS,
			Arch:     runtime.GOARCH,
			Terminal: os.Getenv("TERM"),
			Shell:    os.Getenv("SHELL"),
			SuiteVer: "1.0.0",
		},
		Results:         []Result{},
		DomainSummaries: []DomainSummary{},
		Threshold:       r.config.Threshold.String(),
	}

	// Step 3: Load all domains and filter by domain filter
	allDomains := r.domainProvider()
	selectedDomains := r.filterDomains(allDomains)

	// Step 4: For each domain, execute criteria
	for _, domain := range selectedDomains {
		criteria := domain.Criteria()
		for _, criterion := range criteria {
			// Filter by level
			if !r.shouldRunCriterion(criterion) {
				continue
			}

			// Skip if in skip list
			if r.isSkipped(criterion.ID()) {
				continue
			}

			// Check precondition
			if !criterion.Precondition(probeResult) {
				// Record as NotApplicable
				result := &Result{
					ID:          criterion.ID(),
					Name:        criterion.Name(),
					Domain:      criterion.Domain(),
					Level:       criterion.Level(),
					Testability: criterion.Testability(),
					Outcome:     NotApplicable,
					Evidence:    []Evidence{},
					Remarks:     "Precondition not met",
					Duration:    0,
					SpecVersion: criterion.SpecVersion(),
				}
				report.Results = append(report.Results, *result)
				continue
			}

			// Execute criterion with timeout
			result := r.executeCriterion(ctx, criterion, binary, probeResult)
			report.Results = append(report.Results, *result)
		}
	}

	// Step 5: Compute summaries
	report.ComputeSummaries()

	// Step 6: Determine overall level (highest level with all criteria passing)
	report.OverallLevel = r.computeOverallLevel(report.Results)

	return report, nil
}

// filterDomains filters domains based on the domain filter configuration
func (r *Runner) filterDomains(allDomains []Domain) []Domain {
	if len(r.config.DomainFilter) == 0 {
		return allDomains
	}

	// Create a set of allowed domains
	allowedDomains := make(map[string]bool)
	for _, d := range r.config.DomainFilter {
		allowedDomains[d] = true
	}

	filtered := []Domain{}
	for _, domain := range allDomains {
		if allowedDomains[domain.Name()] {
			filtered = append(filtered, domain)
		}
	}
	return filtered
}

// shouldRunCriterion determines if a criterion should be executed based on filters
func (r *Runner) shouldRunCriterion(criterion Criterion) bool {
	// Filter by level
	if len(r.config.LevelFilter) > 0 {
		levelMatch := false
		for _, level := range r.config.LevelFilter {
			if criterion.Level() == level {
				levelMatch = true
				break
			}
		}
		if !levelMatch {
			return false
		}
	}

	// Filter by testability (AutoOnly)
	if r.config.AutoOnly && criterion.Testability() != Auto {
		return false
	}

	return true
}

// isSkipped checks if a criterion is in the skip list
func (r *Runner) isSkipped(criterionID string) bool {
	for _, skipID := range r.config.SkipCriteria {
		if skipID == criterionID {
			return true
		}
	}
	return false
}

// executeCriterion runs a single criterion with timeout
func (r *Runner) executeCriterion(ctx context.Context, criterion Criterion, binary string, probeResult *probe.ProbeResult) *Result {
	// Create a timeout context if timeout is configured
	execCtx := ctx
	var cancel context.CancelFunc
	if r.config.Timeout > 0 {
		execCtx, cancel = context.WithTimeout(ctx, r.config.Timeout)
		defer cancel()
	}

	// Execute the criterion
	start := time.Now()
	result := criterion.Run(execCtx, binary, probeResult)
	result.Duration = time.Since(start)

	// Check if context was cancelled due to timeout
	if execCtx.Err() == context.DeadlineExceeded && result.Outcome != OutcomeTimeout {
		result.Outcome = OutcomeTimeout
		result.Remarks = "Execution timed out"
	}

	return result
}

// computeOverallLevel determines the highest level with all criteria passing
func (r *Runner) computeOverallLevel(results []Result) string {
	// Count failures by level
	levelAFail := false
	levelAAFail := false
	levelAAAFail := false

	for _, result := range results {
		if result.IsFailure() {
			switch result.Level {
			case LevelA:
				levelAFail = true
			case LevelAA:
				levelAAFail = true
			case LevelAAA:
				levelAAAFail = true
			}
		}
	}

	// Determine overall level
	if levelAFail {
		return "None"
	}
	if levelAAFail {
		return "A"
	}
	if levelAAAFail {
		return "AA"
	}
	return "AAA"
}

// ComputeExitCode computes the exit code based on the conformance report and threshold
// Returns:
//   - 0: All criteria at or below threshold level pass
//   - 1: Any Level A criterion fails
//   - 2: Any Level AA criterion fails (when threshold >= AA)
//   - 3: Any Level AAA criterion fails (when threshold >= AAA)
func ComputeExitCode(report *ConformanceReport, threshold Level) int {
	levelAFail := false
	levelAAFail := false
	levelAAAFail := false

	for _, result := range report.Results {
		// Only consider failures at or below the threshold
		if result.IsFailure() && result.Level <= threshold {
			switch result.Level {
			case LevelA:
				levelAFail = true
			case LevelAA:
				levelAAFail = true
			case LevelAAA:
				levelAAAFail = true
			}
		}
	}

	// Return the highest priority failure
	if levelAFail {
		return 1
	}
	if levelAAFail {
		return 2
	}
	if levelAAAFail {
		return 3
	}
	return 0
}
