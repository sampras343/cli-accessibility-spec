package domains

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/sampras343/cli-accessibility-spec/internal/engine"
	"github.com/sampras343/cli-accessibility-spec/internal/probe"
	"gopkg.in/yaml.v3"
)

// YAMLCriterion implements the Criterion interface from YAML files
type YAMLCriterion struct {
	id           string
	domain       string
	level        engine.Level
	testability  engine.Testability
	specVersion  string
	name         string
	requires     []string
	steps        []yamlStep
	passOutcome  engine.Outcome
	failOutcome  engine.Outcome
}

type yamlStep struct {
	Name     string            `yaml:"name"`
	Exec     *yamlExec         `yaml:"exec,omitempty"`
	ExecPTY  *yamlExec         `yaml:"exec_pty,omitempty"`
	Assert   yamlAssert        `yaml:"assert"`
}

type yamlExec struct {
	Args    []string          `yaml:"args"`
	Env     map[string]string `yaml:"env,omitempty"`
	Timeout int               `yaml:"timeout_ms,omitempty"`
}

type yamlAssert struct {
	ExitCode         *int     `yaml:"exit_code,omitempty"`
	ExitCodeNonzero  *bool    `yaml:"exit_code_nonzero,omitempty"`
	StdoutMatches    []string `yaml:"stdout_matches,omitempty"`
	StdoutNotMatches []string `yaml:"stdout_not_matches,omitempty"`
	StderrMatches    []string `yaml:"stderr_matches,omitempty"`
	StderrNotMatches []string `yaml:"stderr_not_matches,omitempty"`
	StdoutContains   []string `yaml:"stdout_contains,omitempty"`
	StdoutNotContains []string `yaml:"stdout_not_contains,omitempty"`
	StderrContains   []string `yaml:"stderr_contains,omitempty"`
	StderrNotContains []string `yaml:"stderr_not_contains,omitempty"`
	StdoutEmpty      *bool    `yaml:"stdout_empty,omitempty"`
	StderrEmpty      *bool    `yaml:"stderr_empty,omitempty"`
	TimingUnderMs    *int     `yaml:"timing_under_ms,omitempty"`
}

type yamlCriterionFile struct {
	ID              string     `yaml:"id"`
	Domain          string     `yaml:"domain"`
	Level           string     `yaml:"level"`
	Testability     string     `yaml:"testability"`
	SpecVersion     string     `yaml:"spec_version"`
	Name            string     `yaml:"name"`
	Requires        []string   `yaml:"requires,omitempty"`
	Steps           []yamlStep `yaml:"steps"`
	ResultOnAllPass string     `yaml:"result_on_all_pass"`
	ResultOnAnyFail string     `yaml:"result_on_any_fail"`
}

// LoadYAMLCriteria loads all YAML criteria from a directory
func LoadYAMLCriteria(dir string) ([]Criterion, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read dir %s: %w", dir, err)
	}

	var criteria []Criterion
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".yaml") && !strings.HasSuffix(entry.Name(), ".yml") {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", path, err)
		}

		var yf yamlCriterionFile
		if err := yaml.Unmarshal(data, &yf); err != nil {
			return nil, fmt.Errorf("parse %s: %w", path, err)
		}

		level, err := parseLevel(yf.Level)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}

		testability, err := parseTestability(yf.Testability)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}

		passOutcome, err := parseOutcome(yf.ResultOnAllPass)
		if err != nil {
			return nil, fmt.Errorf("%s: result_on_all_pass: %w", path, err)
		}

		failOutcome, err := parseOutcome(yf.ResultOnAnyFail)
		if err != nil {
			return nil, fmt.Errorf("%s: result_on_any_fail: %w", path, err)
		}

		yc := &YAMLCriterion{
			id:          yf.ID,
			domain:      yf.Domain,
			level:       level,
			testability: testability,
			specVersion: yf.SpecVersion,
			name:        yf.Name,
			requires:    yf.Requires,
			steps:       yf.Steps,
			passOutcome: passOutcome,
			failOutcome: failOutcome,
		}

		criteria = append(criteria, yc)
	}

	return criteria, nil
}

// ID implements Criterion
func (y *YAMLCriterion) ID() string {
	return y.id
}

// Name implements Criterion
func (y *YAMLCriterion) Name() string {
	return y.name
}

// Domain implements Criterion
func (y *YAMLCriterion) Domain() string {
	return y.domain
}

// Level implements Criterion
func (y *YAMLCriterion) Level() engine.Level {
	return y.level
}

// Testability implements Criterion
func (y *YAMLCriterion) Testability() engine.Testability {
	return y.testability
}

// SpecVersion implements Criterion
func (y *YAMLCriterion) SpecVersion() string {
	return y.specVersion
}

// Precondition implements Criterion
func (y *YAMLCriterion) Precondition(probeResult *probe.ProbeResult) bool {
	for _, req := range y.requires {
		switch req {
		case "uses_color":
			if !probeResult.HasColor {
				return false
			}
		case "has_help":
			if !probeResult.HasHelp {
				return false
			}
		case "has_version":
			if !probeResult.HasVersion {
				return false
			}
		case "has_subcommands":
			if !probeResult.HasSubcommands {
				return false
			}
		case "has_json_flag":
			if !probeResult.HasJSONFlag {
				return false
			}
		case "has_quiet_flag":
			if !probeResult.HasQuietFlag {
				return false
			}
		case "has_color_flag":
			if !probeResult.HasColorFlag {
				return false
			}
		case "has_dry_run_flag":
			if !probeResult.HasDryRunFlag {
				return false
			}
		case "has_no_input_flag":
			if !probeResult.HasNoInputFlag {
				return false
			}
		default:
			// Unknown requirement - skip this criterion
			return false
		}
	}
	return true
}

// Run implements Criterion
func (y *YAMLCriterion) Run(ctx context.Context, binary string, probeResult *probe.ProbeResult) *engine.Result {
	start := time.Now()

	result := &engine.Result{
		ID:          y.id,
		Name:        y.name,
		Domain:      y.domain,
		Level:       y.level,
		Testability: y.testability,
		SpecVersion: y.specVersion,
		Evidence:    []engine.Evidence{},
	}

	// Check preconditions
	if !y.Precondition(probeResult) {
		result.Outcome = engine.NotApplicable
		result.Remarks = "Preconditions not met"
		result.Duration = time.Since(start)
		return result
	}

	// Execute steps
	allPassed := true
	for i, step := range y.steps {
		stepResult, passed := y.executeStep(ctx, binary, step, i)
		result.Evidence = append(result.Evidence, stepResult)

		if !passed {
			allPassed = false
			result.Outcome = y.failOutcome
			result.Remarks = fmt.Sprintf("Step %d (%s) failed", i+1, step.Name)
			result.Duration = time.Since(start)
			return result
		}
	}

	if allPassed {
		result.Outcome = y.passOutcome
		result.Remarks = "All steps passed"
	}

	result.Duration = time.Since(start)
	return result
}

// executeStep runs a single step and returns evidence and pass/fail status
func (y *YAMLCriterion) executeStep(ctx context.Context, binary string, step yamlStep, stepNum int) (engine.Evidence, bool) {
	var execOpts probe.ExecOpts

	// Determine exec vs exec_pty
	if step.ExecPTY != nil {
		execOpts.Args = step.ExecPTY.Args
		execOpts.Env = step.ExecPTY.Env
		execOpts.UsePTY = true
		if step.ExecPTY.Timeout > 0 {
			execOpts.Timeout = time.Duration(step.ExecPTY.Timeout) * time.Millisecond
		}
	} else if step.Exec != nil {
		execOpts.Args = step.Exec.Args
		execOpts.Env = step.Exec.Env
		execOpts.UsePTY = false
		if step.Exec.Timeout > 0 {
			execOpts.Timeout = time.Duration(step.Exec.Timeout) * time.Millisecond
		}
	} else {
		// No exec specified
		return engine.Evidence{
			Note: fmt.Sprintf("Step %d (%s): no exec/exec_pty specified", stepNum+1, step.Name),
		}, false
	}

	// Execute the command
	execResult, err := probe.Run(ctx, binary, execOpts)
	if err != nil {
		return engine.Evidence{
			Command:  binary + " " + strings.Join(execOpts.Args, " "),
			Note:     fmt.Sprintf("Exec error: %v", err),
			Duration: 0,
		}, false
	}

	evidence := engine.Evidence{
		Command:  execResult.Command,
		Env:      execOpts.Env,
		Stdout:   string(execResult.Stdout),
		Stderr:   string(execResult.Stderr),
		ExitCode: execResult.ExitCode,
		Duration: execResult.Duration,
	}

	// Evaluate assertions
	passed, note := y.evaluateAssertions(step.Assert, execResult)
	evidence.Note = note

	return evidence, passed
}

// evaluateAssertions checks all assertions in a step
func (y *YAMLCriterion) evaluateAssertions(assert yamlAssert, result *probe.ExecResult) (bool, string) {
	// Exit code
	if assert.ExitCode != nil {
		if result.ExitCode != *assert.ExitCode {
			return false, fmt.Sprintf("exit_code: got %d, want %d", result.ExitCode, *assert.ExitCode)
		}
	}

	// Exit code nonzero
	if assert.ExitCodeNonzero != nil {
		if *assert.ExitCodeNonzero && result.ExitCode == 0 {
			return false, "exit_code_nonzero: expected nonzero, got 0"
		}
		if !*assert.ExitCodeNonzero && result.ExitCode != 0 {
			return false, fmt.Sprintf("exit_code_nonzero: expected 0, got %d", result.ExitCode)
		}
	}

	// Stdout matches
	for _, pattern := range assert.StdoutMatches {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return false, fmt.Sprintf("stdout_matches: invalid regex %q: %v", pattern, err)
		}
		if !re.Match(result.Stdout) {
			return false, fmt.Sprintf("stdout_matches: %q not found", pattern)
		}
	}

	// Stdout not matches
	for _, pattern := range assert.StdoutNotMatches {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return false, fmt.Sprintf("stdout_not_matches: invalid regex %q: %v", pattern, err)
		}
		if re.Match(result.Stdout) {
			return false, fmt.Sprintf("stdout_not_matches: %q found but should not match", pattern)
		}
	}

	// Stderr matches
	for _, pattern := range assert.StderrMatches {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return false, fmt.Sprintf("stderr_matches: invalid regex %q: %v", pattern, err)
		}
		if !re.Match(result.Stderr) {
			return false, fmt.Sprintf("stderr_matches: %q not found", pattern)
		}
	}

	// Stderr not matches
	for _, pattern := range assert.StderrNotMatches {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return false, fmt.Sprintf("stderr_not_matches: invalid regex %q: %v", pattern, err)
		}
		if re.Match(result.Stderr) {
			return false, fmt.Sprintf("stderr_not_matches: %q found but should not match", pattern)
		}
	}

	// Stdout contains
	stdout := string(result.Stdout)
	for _, substr := range assert.StdoutContains {
		if !strings.Contains(stdout, substr) {
			return false, fmt.Sprintf("stdout_contains: %q not found", substr)
		}
	}

	// Stdout not contains
	for _, substr := range assert.StdoutNotContains {
		if strings.Contains(stdout, substr) {
			return false, fmt.Sprintf("stdout_not_contains: %q found but should not be present", substr)
		}
	}

	// Stderr contains
	stderr := string(result.Stderr)
	for _, substr := range assert.StderrContains {
		if !strings.Contains(stderr, substr) {
			return false, fmt.Sprintf("stderr_contains: %q not found", substr)
		}
	}

	// Stderr not contains
	for _, substr := range assert.StderrNotContains {
		if strings.Contains(stderr, substr) {
			return false, fmt.Sprintf("stderr_not_contains: %q found but should not be present", substr)
		}
	}

	// Stdout empty
	if assert.StdoutEmpty != nil {
		isEmpty := len(result.Stdout) == 0
		if *assert.StdoutEmpty != isEmpty {
			if *assert.StdoutEmpty {
				return false, "stdout_empty: expected empty, but has content"
			} else {
				return false, "stdout_empty: expected content, but is empty"
			}
		}
	}

	// Stderr empty
	if assert.StderrEmpty != nil {
		isEmpty := len(result.Stderr) == 0
		if *assert.StderrEmpty != isEmpty {
			if *assert.StderrEmpty {
				return false, "stderr_empty: expected empty, but has content"
			} else {
				return false, "stderr_empty: expected content, but is empty"
			}
		}
	}

	// Timing under ms
	if assert.TimingUnderMs != nil {
		maxDuration := time.Duration(*assert.TimingUnderMs) * time.Millisecond
		if result.Duration > maxDuration {
			return false, fmt.Sprintf("timing_under_ms: took %v, want under %v", result.Duration, maxDuration)
		}
	}

	return true, "All assertions passed"
}

// parseLevel converts string to engine.Level
func parseLevel(s string) (engine.Level, error) {
	switch strings.ToUpper(s) {
	case "A":
		return engine.LevelA, nil
	case "AA":
		return engine.LevelAA, nil
	case "AAA":
		return engine.LevelAAA, nil
	default:
		return 0, fmt.Errorf("invalid level %q (must be A, AA, or AAA)", s)
	}
}

// parseTestability converts string to engine.Testability
func parseTestability(s string) (engine.Testability, error) {
	switch strings.ToUpper(s) {
	case "AUTO":
		return engine.Auto, nil
	case "SEMI":
		return engine.Semi, nil
	case "MANUAL":
		return engine.Manual, nil
	default:
		return 0, fmt.Errorf("invalid testability %q (must be AUTO, SEMI, or MANUAL)", s)
	}
}

// parseOutcome converts string to engine.Outcome
func parseOutcome(s string) (engine.Outcome, error) {
	switch s {
	case "Supports":
		return engine.Supports, nil
	case "Partially Supports":
		return engine.PartiallySupports, nil
	case "Does Not Support":
		return engine.DoesNotSupport, nil
	case "Not Applicable":
		return engine.NotApplicable, nil
	case "Not Evaluated":
		return engine.NotEvaluated, nil
	case "Error":
		return engine.OutcomeError, nil
	case "Timeout":
		return engine.OutcomeTimeout, nil
	default:
		return 0, fmt.Errorf("invalid outcome %q", s)
	}
}
