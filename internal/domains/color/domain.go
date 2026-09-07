package color

import (
	"embed"
	"fmt"
	"io/fs"
	"strings"

	"github.com/sampras343/cli-accessibility-spec/internal/domains"
	"github.com/sampras343/cli-accessibility-spec/internal/engine"
	"gopkg.in/yaml.v3"
)

//go:embed testcases/*.yaml
var testcaseFS embed.FS

func init() {
	d := &ColorDomain{}
	// Register Go-defined criteria (complex checks requiring multi-step logic)
	d.criteria = append(d.criteria, &CV1Check{})
	d.criteria = append(d.criteria, &CV8Check{})
	d.criteria = append(d.criteria, &CV9Check{})
	d.criteria = append(d.criteria, &CV12Check{})
	// Load YAML-defined criteria (simple env/flag/output checks)
	if err := d.loadCriteria(); err != nil {
		panic(fmt.Sprintf("failed to load color domain criteria: %v", err))
	}
	domains.Register("color", d)
}

type ColorDomain struct {
	criteria []domains.Criterion
}

func (d *ColorDomain) Name() string {
	return "color"
}

func (d *ColorDomain) Criteria() []domains.Criterion {
	return d.criteria
}

// loadCriteria loads YAML criteria from the embedded filesystem
func (d *ColorDomain) loadCriteria() error {
	entries, err := fs.ReadDir(testcaseFS, "testcases")
	if err != nil {
		return fmt.Errorf("read embedded testcases dir: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".yaml") && !strings.HasSuffix(entry.Name(), ".yml") {
			continue
		}

		path := "testcases/" + entry.Name()
		data, err := testcaseFS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}

		criterion, err := parseYAMLCriterion(data, entry.Name())
		if err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}

		d.criteria = append(d.criteria, criterion)
	}

	return nil
}

// parseYAMLCriterion parses a YAML criterion from bytes
func parseYAMLCriterion(data []byte, filename string) (domains.Criterion, error) {
	var yf yamlCriterionFile
	if err := yaml.Unmarshal(data, &yf); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	level, err := parseLevel(yf.Level)
	if err != nil {
		return nil, err
	}

	testability, err := parseTestability(yf.Testability)
	if err != nil {
		return nil, err
	}

	passOutcome, err := parseOutcome(yf.ResultOnAllPass)
	if err != nil {
		return nil, fmt.Errorf("result_on_all_pass: %w", err)
	}

	failOutcome, err := parseOutcome(yf.ResultOnAnyFail)
	if err != nil {
		return nil, fmt.Errorf("result_on_any_fail: %w", err)
	}

	yc := &yamlCriterion{
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

	return yc, nil
}

// yamlCriterionFile matches the structure in domains.YAMLCriterion
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
