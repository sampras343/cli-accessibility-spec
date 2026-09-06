# CLI-ACS Conformance Suite Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a Go binary (`cli-acs`) that accepts any CLI binary, auto-discovers its capabilities, evaluates it against CLI-ACS v1.0 criteria, and produces accessibility conformance reports in JSON, Markdown, Terminal, and HTML formats.

**Architecture:** Hybrid criterion system — declarative YAML for simple checks, Go for complex checks. Pure auto-discovery via --help parsing. Domain self-registration via init(). Reporter interface for output formats. Process-level isolation with timeout and safety blocklist per criterion execution.

**Tech Stack:** Go 1.22+, cobra (CLI), creack/pty (PTY), gopkg.in/yaml.v3 (YAML), santhosh-tekuri/jsonschema (validation), Go stdlib for everything else.

**Spec:** `spec/conformance-suite-design.md` (design), `spec/CLI_ACS_v1.0.md` (criteria)

## Global Constraints

- Go 1.22+ (for `slices`, `maps` packages)
- Single `go.mod` at repo root
- Domain packages under `internal/domains/` — not public API
- All subprocess execution via `internal/probe/exec.go` — no direct `os/exec` elsewhere
- Per-criterion timeout: 10s default, 30s for timing criteria
- Suite itself must respect `NO_COLOR`, `TERM=dumb`, `--no-color` (self-dogfooding)
- No external test frameworks — use Go `testing` package only
- Commits use `-s` flag for signoff, no Co-Authored-By

---

### Task 1: Go Module Init + Core Types

**Files:**
- Create: `go.mod`
- Create: `internal/engine/types.go`
- Create: `internal/engine/types_test.go`

**Interfaces:**
- Consumes: nothing (foundation)
- Produces: `Level` (type + constants `LevelA`, `LevelAA`, `LevelAAA`), `Testability` (type + constants `Auto`, `Semi`, `Manual`), `Outcome` (type + constants `Supports`, `PartiallySupports`, `DoesNotSupport`, `NotApplicable`, `NotEvaluated`, `Error`, `Timeout`), `Evidence` struct, `Result` struct, `ConformanceReport` struct, `ProductInfo` struct, `TestEnvironment` struct, `DomainSummary` struct. Also `Level.String()`, `Testability.String()`, `Outcome.String()` methods.

- [ ] **Step 1: Initialize Go module**

```bash
cd /home/sacm/Documents/Study/my_projects/cli-accessibility-spec
go mod init github.com/sampras343/cli-accessibility-spec
```

- [ ] **Step 2: Write tests for core type String() methods**

```go
// internal/engine/types_test.go
package engine

import "testing"

func TestLevelString(t *testing.T) {
	tests := []struct {
		level Level
		want  string
	}{
		{LevelA, "A"},
		{LevelAA, "AA"},
		{LevelAAA, "AAA"},
	}
	for _, tt := range tests {
		if got := tt.level.String(); got != tt.want {
			t.Errorf("Level(%d).String() = %q, want %q", tt.level, got, tt.want)
		}
	}
}

func TestOutcomeString(t *testing.T) {
	tests := []struct {
		outcome Outcome
		want    string
	}{
		{Supports, "Supports"},
		{PartiallySupports, "Partially Supports"},
		{DoesNotSupport, "Does Not Support"},
		{NotApplicable, "Not Applicable"},
		{NotEvaluated, "Not Evaluated"},
		{OutcomeError, "Error"},
		{OutcomeTimeout, "Timeout"},
	}
	for _, tt := range tests {
		if got := tt.outcome.String(); got != tt.want {
			t.Errorf("Outcome(%d).String() = %q, want %q", tt.outcome, got, tt.want)
		}
	}
}

func TestTestabilityString(t *testing.T) {
	tests := []struct {
		t    Testability
		want string
	}{
		{Auto, "AUTO"},
		{Semi, "SEMI"},
		{Manual, "MANUAL"},
	}
	for _, tt := range tests {
		if got := tt.t.String(); got != tt.want {
			t.Errorf("Testability(%d).String() = %q, want %q", tt.t, got, tt.want)
		}
	}
}

func TestResultOutcomeIsFailure(t *testing.T) {
	passing := []Outcome{Supports, NotApplicable, NotEvaluated}
	for _, o := range passing {
		r := Result{Outcome: o}
		if r.IsFailure() {
			t.Errorf("Outcome %v should not be a failure", o)
		}
	}
	failing := []Outcome{DoesNotSupport, PartiallySupports, OutcomeError, OutcomeTimeout}
	for _, o := range failing {
		r := Result{Outcome: o}
		if !r.IsFailure() {
			t.Errorf("Outcome %v should be a failure", o)
		}
	}
}
```

- [ ] **Step 3: Run tests — verify they fail**

```bash
go test ./internal/engine/ -v
```
Expected: compilation errors (types don't exist yet)

- [ ] **Step 4: Implement core types**

```go
// internal/engine/types.go
package engine

import (
	"time"
)

type Level int

const (
	LevelA   Level = iota
	LevelAA
	LevelAAA
)

func (l Level) String() string {
	switch l {
	case LevelA:
		return "A"
	case LevelAA:
		return "AA"
	case LevelAAA:
		return "AAA"
	default:
		return "Unknown"
	}
}

type Testability int

const (
	Auto   Testability = iota
	Semi
	Manual
)

func (t Testability) String() string {
	switch t {
	case Auto:
		return "AUTO"
	case Semi:
		return "SEMI"
	case Manual:
		return "MANUAL"
	default:
		return "Unknown"
	}
}

type Outcome int

const (
	Supports          Outcome = iota
	PartiallySupports
	DoesNotSupport
	NotApplicable
	NotEvaluated
	OutcomeError
	OutcomeTimeout
)

func (o Outcome) String() string {
	switch o {
	case Supports:
		return "Supports"
	case PartiallySupports:
		return "Partially Supports"
	case DoesNotSupport:
		return "Does Not Support"
	case NotApplicable:
		return "Not Applicable"
	case NotEvaluated:
		return "Not Evaluated"
	case OutcomeError:
		return "Error"
	case OutcomeTimeout:
		return "Timeout"
	default:
		return "Unknown"
	}
}

type Evidence struct {
	Command  string            `json:"command"`
	Env      map[string]string `json:"env,omitempty"`
	Stdout   string            `json:"stdout"`
	Stderr   string            `json:"stderr"`
	ExitCode int               `json:"exit_code"`
	Duration time.Duration     `json:"duration_ms"`
	Note     string            `json:"note,omitempty"`
}

type Result struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Domain      string        `json:"domain"`
	Level       Level         `json:"level"`
	Testability Testability   `json:"testability"`
	Outcome     Outcome       `json:"outcome"`
	Evidence    []Evidence    `json:"evidence,omitempty"`
	Remarks     string        `json:"remarks,omitempty"`
	Duration    time.Duration `json:"duration_ms"`
	NeedsReview bool          `json:"needs_review,omitempty"`
	SpecVersion string        `json:"spec_version"`
}

func (r *Result) IsFailure() bool {
	switch r.Outcome {
	case DoesNotSupport, PartiallySupports, OutcomeError, OutcomeTimeout:
		return true
	default:
		return false
	}
}

type ProductInfo struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Version string `json:"version"`
}

type TestEnvironment struct {
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Terminal string `json:"terminal,omitempty"`
	Shell    string `json:"shell,omitempty"`
	SuiteVer string `json:"suite_version"`
}

type DomainSummary struct {
	Domain   string `json:"domain"`
	Total    int    `json:"total"`
	Pass     int    `json:"pass"`
	Partial  int    `json:"partial"`
	Fail     int    `json:"fail"`
	NA       int    `json:"na"`
	NotEval  int    `json:"not_evaluated"`
}

type ConformanceReport struct {
	CLIACSVersion   string            `json:"cli_acs_version"`
	SuiteVersion    string            `json:"suite_version"`
	Product         ProductInfo       `json:"product"`
	ReportDate      time.Time         `json:"report_date"`
	Environment     TestEnvironment   `json:"environment"`
	Results         []Result          `json:"results"`
	DomainSummaries []DomainSummary   `json:"domain_summaries"`
	OverallLevel    string            `json:"overall_level"`
	Threshold       string            `json:"threshold"`
}

func (r *ConformanceReport) ComputeSummaries() {
	domainMap := make(map[string]*DomainSummary)
	for _, res := range r.Results {
		s, ok := domainMap[res.Domain]
		if !ok {
			s = &DomainSummary{Domain: res.Domain}
			domainMap[res.Domain] = s
		}
		s.Total++
		switch res.Outcome {
		case Supports:
			s.Pass++
		case PartiallySupports:
			s.Partial++
		case DoesNotSupport, OutcomeError, OutcomeTimeout:
			s.Fail++
		case NotApplicable:
			s.NA++
		case NotEvaluated:
			s.NotEval++
		}
	}
	r.DomainSummaries = make([]DomainSummary, 0, len(domainMap))
	for _, s := range domainMap {
		r.DomainSummaries = append(r.DomainSummaries, *s)
	}
}
```

- [ ] **Step 5: Run tests — verify they pass**

```bash
go test ./internal/engine/ -v
```
Expected: all PASS

- [ ] **Step 6: Commit**

```bash
git add go.mod internal/engine/types.go internal/engine/types_test.go
git commit -s -m "feat: add core types for conformance engine (Level, Outcome, Result, Evidence)"
```

---

### Task 2: Subprocess Execution with Isolation

**Files:**
- Create: `internal/probe/exec.go`
- Create: `internal/probe/exec_test.go`
- Create: `internal/probe/safety.go`
- Create: `internal/probe/safety_test.go`

**Interfaces:**
- Consumes: `engine.Evidence`
- Produces: `ExecResult` struct (Stdout, Stderr []byte, ExitCode int, Duration time.Duration), `ExecOpts` struct (Args, Env, Timeout, UsePTY), `Run(ctx, binary, opts) (*ExecResult, error)`, `BaseEnv(tmpDir) []string`, `IsBlockedSubcommand(name) bool`

- [ ] **Step 1: Write tests for safe execution**

```go
// internal/probe/exec_test.go
package probe

import (
	"context"
	"testing"
	"time"
)

func TestRunCapturesStdoutStderr(t *testing.T) {
	result, err := Run(context.Background(), "sh", ExecOpts{
		Args: []string{"-c", "echo hello; echo err >&2"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if string(result.Stdout) != "hello\n" {
		t.Errorf("stdout = %q, want %q", result.Stdout, "hello\n")
	}
	if string(result.Stderr) != "err\n" {
		t.Errorf("stderr = %q, want %q", result.Stderr, "err\n")
	}
	if result.ExitCode != 0 {
		t.Errorf("exit code = %d, want 0", result.ExitCode)
	}
}

func TestRunCapturesNonZeroExit(t *testing.T) {
	result, err := Run(context.Background(), "sh", ExecOpts{
		Args: []string{"-c", "exit 42"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.ExitCode != 42 {
		t.Errorf("exit code = %d, want 42", result.ExitCode)
	}
}

func TestRunTimesOut(t *testing.T) {
	result, err := Run(context.Background(), "sleep", ExecOpts{
		Args:    []string{"60"},
		Timeout: 500 * time.Millisecond,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.ExitCode == 0 {
		t.Error("expected non-zero exit on timeout")
	}
	if result.TimedOut != true {
		t.Error("expected TimedOut = true")
	}
}

func TestRunIsolatesEnv(t *testing.T) {
	result, err := Run(context.Background(), "sh", ExecOpts{
		Args: []string{"-c", "echo $HOME"},
		Env:  map[string]string{"HOME": "/tmp/fake"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if string(result.Stdout) != "/tmp/fake\n" {
		t.Errorf("HOME not isolated: got %q", result.Stdout)
	}
}
```

```go
// internal/probe/safety_test.go
package probe

import "testing"

func TestIsBlockedSubcommand(t *testing.T) {
	blocked := []string{"delete", "rm", "remove", "destroy", "purge", "DROP", "Format", "init", "reset", "clean", "wipe", "nuke", "uninstall"}
	for _, s := range blocked {
		if !IsBlockedSubcommand(s) {
			t.Errorf("%q should be blocked", s)
		}
	}

	allowed := []string{"list", "get", "show", "status", "help", "version", "config", "describe"}
	for _, s := range allowed {
		if IsBlockedSubcommand(s) {
			t.Errorf("%q should be allowed", s)
		}
	}
}
```

- [ ] **Step 2: Run tests — verify they fail**

```bash
go test ./internal/probe/ -v
```

- [ ] **Step 3: Implement exec.go**

```go
// internal/probe/exec.go
package probe

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

type ExecOpts struct {
	Args    []string
	Env     map[string]string
	Timeout time.Duration
	UsePTY  bool
}

type ExecResult struct {
	Stdout   []byte
	Stderr   []byte
	ExitCode int
	Duration time.Duration
	TimedOut bool
	Command  string
}

func Run(ctx context.Context, binary string, opts ExecOpts) (*ExecResult, error) {
	timeout := opts.Timeout
	if timeout == 0 {
		timeout = 10 * time.Second
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, binary, opts.Args...)

	cmd.Env = buildEnv(opts.Env)

	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	start := time.Now()
	err := cmd.Run()
	duration := time.Since(start)

	exitCode := 0
	timedOut := false

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			timedOut = true
			exitCode = -1
		} else if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			return nil, fmt.Errorf("exec %s: %w", binary, err)
		}
	}

	cmdStr := binary
	if len(opts.Args) > 0 {
		cmdStr += " " + strings.Join(opts.Args, " ")
	}

	return &ExecResult{
		Stdout:   stdout.Bytes(),
		Stderr:   stderr.Bytes(),
		ExitCode: exitCode,
		Duration: duration,
		TimedOut: timedOut,
		Command:  cmdStr,
	}, nil
}

func buildEnv(overrides map[string]string) []string {
	tmpDir, _ := os.MkdirTemp("", "cli-acs-*")
	homeDir := filepath.Join(tmpDir, "home")
	os.MkdirAll(homeDir, 0755)

	base := map[string]string{
		"PATH":   os.Getenv("PATH"),
		"HOME":   homeDir,
		"TMPDIR": tmpDir,
		"LANG":   "C.UTF-8",
		"TERM":   "xterm-256color",
	}

	for k, v := range overrides {
		base[k] = v
	}

	env := make([]string, 0, len(base))
	for k, v := range base {
		env = append(env, k+"="+v)
	}
	return env
}
```

- [ ] **Step 4: Implement safety.go**

```go
// internal/probe/safety.go
package probe

import "strings"

var blockedSubcommands = []string{
	"delete", "rm", "remove", "destroy", "purge", "drop",
	"format", "init", "reset", "clean", "wipe", "nuke",
	"truncate", "uninstall", "erase",
}

func IsBlockedSubcommand(name string) bool {
	lower := strings.ToLower(name)
	for _, blocked := range blockedSubcommands {
		if lower == blocked {
			return true
		}
	}
	return false
}
```

- [ ] **Step 5: Run tests — verify they pass**

```bash
go test ./internal/probe/ -v -timeout 30s
```

- [ ] **Step 6: Commit**

```bash
git add internal/probe/exec.go internal/probe/exec_test.go internal/probe/safety.go internal/probe/safety_test.go
git commit -s -m "feat: add isolated subprocess execution with timeout and safety blocklist"
```

---

### Task 3: ANSI Detection & Stripping

**Files:**
- Create: `internal/probe/ansi.go`
- Create: `internal/probe/ansi_test.go`

**Interfaces:**
- Consumes: nothing
- Produces: `HasANSI(data []byte) bool`, `StripANSI(data []byte) []byte`, `HasColorCodes(data []byte) bool`, `ClassifyColors(data []byte) ColorClassification`, `ColorClassification` struct (FourBit, EightBit, TwentyFourBit int)

- [ ] **Step 1: Write tests**

```go
// internal/probe/ansi_test.go
package probe

import "testing"

func TestHasANSI(t *testing.T) {
	if HasANSI([]byte("plain text")) {
		t.Error("plain text should not have ANSI")
	}
	if !HasANSI([]byte("\x1b[31mred\x1b[0m")) {
		t.Error("colored text should have ANSI")
	}
	if !HasANSI([]byte("\x1b[1mbold\x1b[0m")) {
		t.Error("bold text should have ANSI")
	}
}

func TestStripANSI(t *testing.T) {
	input := []byte("\x1b[1;31mError:\x1b[0m file not found")
	got := string(StripANSI(input))
	want := "Error: file not found"
	if got != want {
		t.Errorf("StripANSI = %q, want %q", got, want)
	}
}

func TestHasColorCodes(t *testing.T) {
	if HasColorCodes([]byte("\x1b[1mbold only\x1b[0m")) {
		t.Error("bold-only should not count as color")
	}
	if !HasColorCodes([]byte("\x1b[31mred\x1b[0m")) {
		t.Error("SGR 31 is a color code")
	}
	if !HasColorCodes([]byte("\x1b[38;5;196mred\x1b[0m")) {
		t.Error("8-bit color should be detected")
	}
	if !HasColorCodes([]byte("\x1b[38;2;255;0;0mred\x1b[0m")) {
		t.Error("24-bit color should be detected")
	}
}

func TestClassifyColors(t *testing.T) {
	input := []byte("\x1b[31mred\x1b[0m \x1b[38;5;196m256red\x1b[0m \x1b[38;2;0;255;0mtrue\x1b[0m")
	c := ClassifyColors(input)
	if c.FourBit < 1 {
		t.Error("expected at least 1 four-bit color")
	}
	if c.EightBit < 1 {
		t.Error("expected at least 1 eight-bit color")
	}
	if c.TwentyFourBit < 1 {
		t.Error("expected at least 1 twenty-four-bit color")
	}
}
```

- [ ] **Step 2: Run tests — verify they fail**

```bash
go test ./internal/probe/ -run TestHasANSI -v
```

- [ ] **Step 3: Implement ansi.go**

```go
// internal/probe/ansi.go
package probe

import "regexp"

var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
var colorFGPattern = regexp.MustCompile(`\x1b\[(3[0-7]|9[0-7]|38;5;\d+|38;2;\d+;\d+;\d+)m`)
var colorBGPattern = regexp.MustCompile(`\x1b\[(4[0-7]|10[0-7]|48;5;\d+|48;2;\d+;\d+;\d+)m`)
var eightBitPattern = regexp.MustCompile(`\x1b\[(38|48);5;(\d+)m`)
var twentyFourBitPattern = regexp.MustCompile(`\x1b\[(38|48);2;\d+;\d+;\d+m`)
var fourBitFGPattern = regexp.MustCompile(`\x1b\[(3[0-7]|9[0-7])m`)

type ColorClassification struct {
	FourBit        int
	EightBit       int
	TwentyFourBit  int
}

func HasANSI(data []byte) bool {
	return ansiPattern.Match(data)
}

func StripANSI(data []byte) []byte {
	return ansiPattern.ReplaceAll(data, nil)
}

func HasColorCodes(data []byte) bool {
	return colorFGPattern.Match(data) || colorBGPattern.Match(data)
}

func ClassifyColors(data []byte) ColorClassification {
	var c ColorClassification
	c.TwentyFourBit = len(twentyFourBitPattern.FindAll(data, -1))
	c.EightBit = len(eightBitPattern.FindAll(data, -1))
	c.FourBit = len(fourBitFGPattern.FindAll(data, -1))
	return c
}
```

- [ ] **Step 4: Run tests — verify they pass**

```bash
go test ./internal/probe/ -v
```

- [ ] **Step 5: Commit**

```bash
git add internal/probe/ansi.go internal/probe/ansi_test.go
git commit -s -m "feat: add ANSI escape sequence detection, stripping, and color classification"
```

---

### Task 4: Domain Registry & Criterion Interface

**Files:**
- Create: `internal/domains/criterion.go`
- Create: `internal/domains/registry.go`
- Create: `internal/domains/registry_test.go`

**Interfaces:**
- Consumes: `engine.Level`, `engine.Testability`, `engine.Outcome`, `engine.Result`, `probe.ExecResult`
- Produces: `Criterion` interface (ID, Name, Domain, Level, Testability, SpecVersion, Precondition, Run methods), `Domain` interface (Name, Criteria), `Register(name, Domain)`, `AllDomains() []Domain`, `ProbeResult` struct

- [ ] **Step 1: Write tests**

```go
// internal/domains/registry_test.go
package domains

import (
	"context"
	"testing"

	"github.com/sampras343/cli-accessibility-spec/internal/engine"
)

type mockDomain struct {
	name     string
	criteria []Criterion
}

func (d *mockDomain) Name() string       { return d.name }
func (d *mockDomain) Criteria() []Criterion { return d.criteria }

type mockCriterion struct {
	id     string
	domain string
}

func (c *mockCriterion) ID() string                { return c.id }
func (c *mockCriterion) Name() string              { return c.id + " test" }
func (c *mockCriterion) Domain() string            { return c.domain }
func (c *mockCriterion) Level() engine.Level       { return engine.LevelA }
func (c *mockCriterion) Testability() engine.Testability { return engine.Auto }
func (c *mockCriterion) SpecVersion() string       { return "1.0" }
func (c *mockCriterion) Precondition(_ *ProbeResult) bool { return true }
func (c *mockCriterion) Run(_ context.Context, _ string, _ *ProbeResult) *engine.Result {
	return &engine.Result{ID: c.id, Outcome: engine.Supports}
}

func TestRegisterAndRetrieve(t *testing.T) {
	resetRegistry()
	Register("test", &mockDomain{
		name:     "test",
		criteria: []Criterion{&mockCriterion{id: "T-1", domain: "test"}},
	})
	domains := AllDomains()
	if len(domains) != 1 {
		t.Fatalf("expected 1 domain, got %d", len(domains))
	}
	if domains[0].Name() != "test" {
		t.Errorf("domain name = %q, want %q", domains[0].Name(), "test")
	}
	if len(domains[0].Criteria()) != 1 {
		t.Fatalf("expected 1 criterion, got %d", len(domains[0].Criteria()))
	}
	if domains[0].Criteria()[0].ID() != "T-1" {
		t.Errorf("criterion ID = %q, want %q", domains[0].Criteria()[0].ID(), "T-1")
	}
}

func TestRegisterPanicsOnDuplicate(t *testing.T) {
	resetRegistry()
	Register("dup", &mockDomain{name: "dup"})
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic on duplicate registration")
		}
	}()
	Register("dup", &mockDomain{name: "dup"})
}
```

- [ ] **Step 2: Run tests — verify they fail**

```bash
go test ./internal/domains/ -v
```

- [ ] **Step 3: Implement criterion.go and registry.go**

```go
// internal/domains/criterion.go
package domains

import (
	"context"

	"github.com/sampras343/cli-accessibility-spec/internal/engine"
)

type ProbeResult struct {
	BinaryPath     string
	BinaryName     string
	Version        string
	HelpText       string
	Subcommands    []Subcommand
	GlobalFlags    []Flag
	HasColor       bool
	HasHelp        bool
	HasVersion     bool
	HasSubcommands bool
	HasJSONFlag    bool
	HasQuietFlag   bool
	HasColorFlag   bool
	HasDryRunFlag  bool
	HasNoInputFlag bool
	HelpFormat     string
	SampleErrors   []ErrorSample
}

type Subcommand struct {
	Name     string
	HelpText string
	Flags    []Flag
}

type Flag struct {
	Short       string
	Long        string
	Description string
	TakesValue  bool
}

type ErrorSample struct {
	Trigger  string
	Stdout   []byte
	Stderr   []byte
	ExitCode int
}

type Criterion interface {
	ID() string
	Name() string
	Domain() string
	Level() engine.Level
	Testability() engine.Testability
	SpecVersion() string
	Precondition(probe *ProbeResult) bool
	Run(ctx context.Context, binary string, probe *ProbeResult) *engine.Result
}

type Domain interface {
	Name() string
	Criteria() []Criterion
}
```

```go
// internal/domains/registry.go
package domains

import (
	"fmt"
	"sync"
)

var (
	mu       sync.Mutex
	registry = make(map[string]Domain)
)

func Register(name string, d Domain) {
	mu.Lock()
	defer mu.Unlock()
	if _, exists := registry[name]; exists {
		panic(fmt.Sprintf("domain %q already registered", name))
	}
	registry[name] = d
}

func AllDomains() []Domain {
	mu.Lock()
	defer mu.Unlock()
	domains := make([]Domain, 0, len(registry))
	for _, d := range registry {
		domains = append(domains, d)
	}
	return domains
}

func GetDomain(name string) (Domain, bool) {
	mu.Lock()
	defer mu.Unlock()
	d, ok := registry[name]
	return d, ok
}

func resetRegistry() {
	mu.Lock()
	defer mu.Unlock()
	registry = make(map[string]Domain)
}
```

- [ ] **Step 4: Run tests — verify they pass**

```bash
go test ./internal/domains/ -v
```

- [ ] **Step 5: Commit**

```bash
git add internal/domains/criterion.go internal/domains/registry.go internal/domains/registry_test.go
git commit -s -m "feat: add domain registry with self-registration and criterion interface"
```

---

### Task 5: YAML Criterion Loader

**Files:**
- Create: `internal/domains/yaml_criterion.go`
- Create: `internal/domains/yaml_criterion_test.go`
- Create: `schemas/criterion.schema.json`

**Interfaces:**
- Consumes: `Criterion` interface, `ProbeResult`, `probe.Run()`, `probe.HasANSI()`, `probe.HasColorCodes()`
- Produces: `LoadYAMLCriteria(dir string) ([]Criterion, error)`, `YAMLCriterion` struct (implements Criterion)

- [ ] **Step 1: Write a sample YAML criterion for testing**

Create `internal/domains/testdata/sample.yaml`:

```yaml
id: "TEST-1"
domain: "test"
level: "A"
testability: "AUTO"
spec_version: "1.0"
name: "Sample Test Criterion"
steps:
  - name: "Check exit code"
    exec:
      args: ["--help"]
      env: {}
    assert:
      exit_code: 0
result_on_all_pass: "Supports"
result_on_any_fail: "Does Not Support"
```

- [ ] **Step 2: Write tests**

```go
// internal/domains/yaml_criterion_test.go
package domains

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadYAMLCriteria(t *testing.T) {
	dir := filepath.Join("testdata")
	criteria, err := LoadYAMLCriteria(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(criteria) == 0 {
		t.Fatal("expected at least one criterion")
	}
	c := criteria[0]
	if c.ID() != "TEST-1" {
		t.Errorf("ID = %q, want %q", c.ID(), "TEST-1")
	}
	if c.Name() != "Sample Test Criterion" {
		t.Errorf("Name = %q, want %q", c.Name(), "Sample Test Criterion")
	}
	if c.Domain() != "test" {
		t.Errorf("Domain = %q, want %q", c.Domain(), "test")
	}
}

func TestYAMLCriterionRunPassesWithEcho(t *testing.T) {
	// Create a temp YAML that tests "echo" (which always exits 0)
	dir := t.TempDir()
	yaml := `id: "ECHO-1"
domain: "test"
level: "A"
testability: "AUTO"
spec_version: "1.0"
name: "Echo test"
steps:
  - name: "echo exits 0"
    exec:
      args: ["-c", "echo hello"]
      env: {}
    assert:
      exit_code: 0
      stdout_contains:
        - "hello"
result_on_all_pass: "Supports"
result_on_any_fail: "Does Not Support"
`
	os.WriteFile(filepath.Join(dir, "ECHO-1.yaml"), []byte(yaml), 0644)

	criteria, err := LoadYAMLCriteria(dir)
	if err != nil {
		t.Fatal(err)
	}

	probe := &ProbeResult{HasHelp: true}
	result := criteria[0].Run(context.Background(), "sh", probe)
	if result.Outcome != Supports {
		t.Errorf("outcome = %v, want Supports. Remarks: %s", result.Outcome, result.Remarks)
	}
}
```

- [ ] **Step 3: Run tests — verify they fail**

```bash
go test ./internal/domains/ -run TestLoadYAML -v
```

- [ ] **Step 4: Implement yaml_criterion.go**

This is the YAML DSL engine — it parses YAML criterion files, builds step sequences, executes them against the binary under test, and evaluates assertions. See the design doc Section 6 for the full schema. The implementation parses `steps[].exec/exec_pty`, runs each via `probe.Run()`, evaluates `assert` conditions (exit_code, stdout_matches/not_matches, stderr_matches/not_matches, stdout_contains/not_contains, timing_under_ms), and produces a Result.

Key implementation details:
- Parse `level` string → `engine.Level` constant
- Parse `testability` string → `engine.Testability` constant
- Parse `requires` list → check against `ProbeResult` fields
- Each step: build `probe.ExecOpts` from YAML, call `probe.Run()`, evaluate assertions
- All steps pass → `result_on_all_pass` outcome; any step fails → `result_on_any_fail` outcome
- Collect `Evidence` from each step's execution

- [ ] **Step 5: Run tests — verify they pass**

```bash
go test ./internal/domains/ -v
```

- [ ] **Step 6: Create JSON Schema**

Create `schemas/criterion.schema.json` with the full schema for YAML validation. This enables `cli-acs validate` and editor autocomplete.

- [ ] **Step 7: Commit**

```bash
git add internal/domains/yaml_criterion.go internal/domains/yaml_criterion_test.go internal/domains/testdata/ schemas/
git commit -s -m "feat: add YAML criterion loader with DSL execution engine"
```

---

### Task 6: Probe / Auto-Discovery

**Files:**
- Create: `internal/probe/discovery.go`
- Create: `internal/probe/parsers.go`
- Create: `internal/probe/discovery_test.go`
- Create: `internal/probe/parsers_test.go`

**Interfaces:**
- Consumes: `probe.Run()`, `probe.HasANSI()`, `domains.ProbeResult`, `domains.Subcommand`, `domains.Flag`
- Produces: `Probe(ctx, binary, maxSubcmds) (*domains.ProbeResult, error)`, `HelpParser` interface, `CobraParser`, `ClapParser`, `ArgparseParser`, `GenericParser`

- [ ] **Step 1: Write parser tests with sample help outputs**

```go
// internal/probe/parsers_test.go
package probe

import "testing"

const cobraHelp = `A CLI tool for things

Usage:
  mytool [command]

Available Commands:
  list        List all items
  get         Get a specific item
  completion  Generate shell completions
  help        Help about any command

Flags:
  -h, --help      help for mytool
  -v, --version   version for mytool
      --json      Output as JSON
  -q, --quiet     Suppress output

Use "mytool [command] --help" for more information about a command.
`

func TestCobraParserCanParse(t *testing.T) {
	p := &CobraParser{}
	if !p.CanParse(cobraHelp) {
		t.Error("CobraParser should recognize cobra-style help")
	}
	if p.CanParse("usage: tool [options] file") {
		t.Error("CobraParser should not match non-cobra help")
	}
}

func TestCobraParserExtractsSubcommands(t *testing.T) {
	p := &CobraParser{}
	result, err := p.Parse(cobraHelp)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Subcommands) < 3 {
		t.Errorf("expected >=3 subcommands, got %d", len(result.Subcommands))
	}
	found := false
	for _, s := range result.Subcommands {
		if s.Name == "list" {
			found = true
		}
	}
	if !found {
		t.Error("expected 'list' subcommand")
	}
}

func TestCobraParserExtractsFlags(t *testing.T) {
	p := &CobraParser{}
	result, err := p.Parse(cobraHelp)
	if err != nil {
		t.Fatal(err)
	}
	foundJSON := false
	foundQuiet := false
	for _, f := range result.Flags {
		if f.Long == "--json" {
			foundJSON = true
		}
		if f.Long == "--quiet" {
			foundQuiet = true
		}
	}
	if !foundJSON {
		t.Error("expected --json flag")
	}
	if !foundQuiet {
		t.Error("expected --quiet flag")
	}
}
```

- [ ] **Step 2: Run tests — verify they fail**

```bash
go test ./internal/probe/ -run TestCobra -v
```

- [ ] **Step 3: Implement parsers.go with cobra, clap, argparse, and generic parsers**

Each parser implements `HelpParser` interface: `Name() string`, `CanParse(text string) bool`, `Parse(text string) (*ParsedHelp, error)`. Parsers are tried in cascade order. `ParsedHelp` contains `Subcommands []domains.Subcommand` and `Flags []domains.Flag`.

- [ ] **Step 4: Implement discovery.go**

`Probe()` function: runs `--help`, `--version`, parses help text via cascade, detects color via PTY, detects common flags, discovers subcommands (up to max), triggers errors. Returns `*domains.ProbeResult`.

- [ ] **Step 5: Write and run discovery integration test using `gh` binary**

```go
func TestProbeGH(t *testing.T) {
	if _, err := exec.LookPath("gh"); err != nil {
		t.Skip("gh not installed")
	}
	result, err := Probe(context.Background(), "gh", 5)
	if err != nil {
		t.Fatal(err)
	}
	if !result.HasHelp {
		t.Error("gh should have --help")
	}
	if !result.HasVersion {
		t.Error("gh should have --version")
	}
	if !result.HasSubcommands {
		t.Error("gh should have subcommands")
	}
	if result.HelpFormat != "cobra" {
		t.Errorf("gh help format = %q, want cobra", result.HelpFormat)
	}
}
```

- [ ] **Step 6: Run all probe tests**

```bash
go test ./internal/probe/ -v -timeout 60s
```

- [ ] **Step 7: Commit**

```bash
git add internal/probe/discovery.go internal/probe/parsers.go internal/probe/discovery_test.go internal/probe/parsers_test.go
git commit -s -m "feat: add auto-discovery probe with help parser cascade"
```

---

### Task 7: Engine Runner

**Files:**
- Create: `internal/engine/runner.go`
- Create: `internal/engine/config.go`
- Create: `internal/engine/runner_test.go`

**Interfaces:**
- Consumes: `domains.AllDomains()`, `domains.Criterion`, `domains.ProbeResult`, `probe.Probe()`, `engine.Result`, `engine.ConformanceReport`
- Produces: `Config` struct, `Runner` struct, `NewRunner(config Config) *Runner`, `Runner.Run(ctx, binary) (*ConformanceReport, error)`, `Runner.ComputeExitCode(report) int`

- [ ] **Step 1: Write tests**

```go
// internal/engine/runner_test.go
package engine

import "testing"

func TestComputeExitCode(t *testing.T) {
	tests := []struct {
		name      string
		results   []Result
		threshold Level
		want      int
	}{
		{
			name:      "all pass",
			results:   []Result{{Outcome: Supports, Level: LevelA}},
			threshold: LevelA,
			want:      0,
		},
		{
			name:      "level A failure",
			results:   []Result{{Outcome: DoesNotSupport, Level: LevelA}},
			threshold: LevelA,
			want:      1,
		},
		{
			name:      "level AA failure with AA threshold",
			results:   []Result{{Outcome: Supports, Level: LevelA}, {Outcome: DoesNotSupport, Level: LevelAA}},
			threshold: LevelAA,
			want:      2,
		},
		{
			name:      "level AA failure with A threshold passes",
			results:   []Result{{Outcome: Supports, Level: LevelA}, {Outcome: DoesNotSupport, Level: LevelAA}},
			threshold: LevelA,
			want:      0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report := &ConformanceReport{Results: tt.results, Threshold: tt.threshold.String()}
			got := ComputeExitCode(report, tt.threshold)
			if got != tt.want {
				t.Errorf("ComputeExitCode = %d, want %d", got, tt.want)
			}
		})
	}
}
```

- [ ] **Step 2: Implement config.go and runner.go**

`Config` holds all CLI flags (domain filter, level filter, auto-only, threshold, skip list, timeout, max-subcommands, criteria-dir). `Runner.Run()` orchestrates: probe → filter criteria → execute each → build report.

- [ ] **Step 3: Run tests**

```bash
go test ./internal/engine/ -v
```

- [ ] **Step 4: Commit**

```bash
git add internal/engine/runner.go internal/engine/config.go internal/engine/runner_test.go
git commit -s -m "feat: add engine runner with config, criterion filtering, and exit code computation"
```

---

### Task 8: Terminal + JSON Reporters

**Files:**
- Create: `internal/report/reporter.go`
- Create: `internal/report/json.go`
- Create: `internal/report/terminal.go`
- Create: `internal/report/json_test.go`
- Create: `internal/report/terminal_test.go`

**Interfaces:**
- Consumes: `engine.ConformanceReport`, `engine.Result`, `engine.DomainSummary`
- Produces: `Reporter` interface (Name, FileExtension, Render), `JSONReporter`, `TerminalReporter`, `GetReporter(name) Reporter`

- [ ] **Step 1: Write tests**

```go
// internal/report/json_test.go
package report

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/sampras343/cli-accessibility-spec/internal/engine"
)

func TestJSONReporterOutput(t *testing.T) {
	report := &engine.ConformanceReport{
		CLIACSVersion: "1.0",
		Product:       engine.ProductInfo{Name: "test-tool", Version: "1.0"},
		Results: []engine.Result{
			{ID: "CV-2", Outcome: engine.Supports, Domain: "color"},
		},
	}
	report.ComputeSummaries()

	var buf bytes.Buffer
	r := &JSONReporter{}
	if err := r.Render(report, &buf); err != nil {
		t.Fatal(err)
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if parsed["cli_acs_version"] != "1.0" {
		t.Error("missing cli_acs_version")
	}
}
```

- [ ] **Step 2: Implement reporter.go, json.go, terminal.go**

`terminal.go` renders a colored summary respecting `NO_COLOR` and `--plain`. Uses `[PASS]`, `[FAIL]`, `[PARTIAL]`, `[N/A]`, `[SKIP]` prefixes for screen reader accessibility.

- [ ] **Step 3: Run tests**

```bash
go test ./internal/report/ -v
```

- [ ] **Step 4: Commit**

```bash
git add internal/report/
git commit -s -m "feat: add JSON and terminal report renderers"
```

---

### Task 9: Cobra CLI Entrypoint

**Files:**
- Create: `cmd/cli-acs/main.go`
- Create: `cmd/cli-acs/check.go`
- Create: `cmd/cli-acs/validate.go`
- Create: `cmd/cli-acs/version.go`

**Interfaces:**
- Consumes: `engine.Runner`, `engine.Config`, `report.GetReporter()`, `domains.AllDomains()`
- Produces: `cli-acs` binary with `check`, `validate`, `version` subcommands and all flags from design doc Section 2

- [ ] **Step 1: Install cobra dependency**

```bash
go get github.com/spf13/cobra
```

- [ ] **Step 2: Implement main.go with root command**

Root command with `--no-color` global flag. Self-dogfooding: check `NO_COLOR`, `TERM=dumb` at startup.

- [ ] **Step 3: Implement check.go**

All flags from design doc Section 2. Wires: parse flags → build Config → create Runner → Runner.Run() → select Reporter → Render → exit with computed code.

- [ ] **Step 4: Implement validate.go and version.go**

`validate`: loads YAML from path, validates against JSON Schema, checks testcase.md presence. `version`: prints suite version + spec version.

- [ ] **Step 5: Build and smoke test**

```bash
go build -o cli-acs ./cmd/cli-acs/
./cli-acs version
./cli-acs check --help
./cli-acs check /usr/bin/echo
```

- [ ] **Step 6: Commit**

```bash
git add cmd/ go.sum
git commit -s -m "feat: add cobra CLI with check, validate, and version commands"
```

---

### Task 10: First Domain — Color (CV-2, CV-3, CV-4, CV-5)

**Files:**
- Create: `internal/domains/color/domain.go`
- Create: `internal/domains/color/checks.go`
- Create: `internal/domains/color/testcases/CV-2.yaml`
- Create: `internal/domains/color/testcases/CV-2.testcase.md`
- Create: `internal/domains/color/testcases/CV-3.yaml`
- Create: `internal/domains/color/testcases/CV-3.testcase.md`
- Create: `internal/domains/color/testcases/CV-4.yaml`
- Create: `internal/domains/color/testcases/CV-4.testcase.md`
- Create: `internal/domains/color/testcases/CV-5.yaml`
- Create: `internal/domains/color/testcases/CV-5.testcase.md`

**Interfaces:**
- Consumes: `domains.Register()`, `domains.Criterion`, `domains.ProbeResult`
- Produces: `ColorDomain` (registered as "color"), 4 YAML criteria (CV-2, CV-3, CV-4, CV-5) as the reference implementation for all future criteria

- [ ] **Step 1: Create domain.go with init() self-registration**

```go
// internal/domains/color/domain.go
package color

import "github.com/sampras343/cli-accessibility-spec/internal/domains"

func init() {
	domains.Register("color", &ColorDomain{})
}

type ColorDomain struct {
	criteria []domains.Criterion
}

func (d *ColorDomain) Name() string             { return "color" }
func (d *ColorDomain) Criteria() []domains.Criterion { return d.criteria }
```

- [ ] **Step 2: Create CV-2.yaml and CV-2.testcase.md**

YAML criterion per design doc Section 6 example. testcase.md per design doc Section 8 template.

- [ ] **Step 3: Create CV-3.yaml, CV-4.yaml, CV-5.yaml with testcase.md files**

Each follows the same YAML DSL pattern. CV-3 tests `--no-color` flag. CV-4 tests piped output. CV-5 tests `TERM=dumb`.

- [ ] **Step 4: Verify domain loads and criteria run against `/usr/bin/ls`**

```bash
go build -o cli-acs ./cmd/cli-acs/
./cli-acs check --domain color /usr/bin/ls
```

- [ ] **Step 5: Commit**

```bash
git add internal/domains/color/
git commit -s -m "feat: add color domain with CV-2, CV-3, CV-4, CV-5 YAML criteria"
```

---

### Task 11: Markdown + HTML Reporters

**Files:**
- Create: `internal/report/markdown.go`
- Create: `internal/report/html.go`
- Create: `internal/report/markdown_test.go`
- Create: `internal/report/html_test.go`

**Interfaces:**
- Consumes: `Reporter` interface, `engine.ConformanceReport`
- Produces: `MarkdownReporter`, `HTMLReporter`

- [ ] **Step 1: Implement markdown.go**

Renders CLI-ACR template from spec Section 7: header, conformance summary table, per-criterion results with evidence, disability impact summary, testing environment.

- [ ] **Step 2: Implement html.go**

Standalone HTML with embedded CSS. Uses `html/template`. Includes expandable evidence sections. No external dependencies (single file, no CDN).

- [ ] **Step 3: Write tests and run**

```bash
go test ./internal/report/ -v
```

- [ ] **Step 4: Commit**

```bash
git add internal/report/markdown.go internal/report/html.go internal/report/markdown_test.go internal/report/html_test.go
git commit -s -m "feat: add markdown and HTML report renderers"
```

---

### Task 12: Fixture Binaries + End-to-End Test

**Files:**
- Create: `testdata/fixtures/good/main.go`
- Create: `testdata/fixtures/no-nocolor/main.go`
- Create: `testdata/fixtures/no-help/main.go`
- Create: `testdata/fixtures/generate.go`
- Create: `e2e_test.go`

**Interfaces:**
- Consumes: entire suite
- Produces: fixture binaries, end-to-end integration tests

- [ ] **Step 1: Create fixture-good — a binary that passes all criteria**

```go
// testdata/fixtures/good/main.go
package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

func main() {
	noColor := os.Getenv("NO_COLOR") != ""
	termDumb := os.Getenv("TERM") == "dumb"
	isTTY := term.IsTerminal(int(os.Stdout.Fd()))

	for _, arg := range os.Args[1:] {
		switch arg {
		case "--help", "-h":
			printHelp(noColor || termDumb || !isTTY)
			os.Exit(0)
		case "--version":
			fmt.Println("fixture-good v1.0.0")
			os.Exit(0)
		case "--no-color":
			noColor = true
		case "--json":
			fmt.Println(`{"status":"ok"}`)
			os.Exit(0)
		case "--quiet", "-q":
			os.Exit(0)
		}
	}

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Error: missing command. Run with --help for usage.")
		os.Exit(1)
	}
	fmt.Fprintln(os.Stderr, "Error: unknown command '"+os.Args[1]+"'. Run with --help for usage.")
	os.Exit(1)
}

func printHelp(plain bool) {
	if plain {
		fmt.Println("fixture-good - a test fixture for CLI-ACS")
		fmt.Println()
		fmt.Println("USAGE")
		fmt.Println("  fixture-good [flags]")
	} else {
		fmt.Println("\x1b[1mfixture-good\x1b[0m - a test fixture for CLI-ACS")
		fmt.Println()
		fmt.Println("\x1b[1mUSAGE\x1b[0m")
		fmt.Println("  fixture-good [flags]")
	}
	fmt.Println()
	fmt.Println("FLAGS")
	fmt.Println("  -h, --help       Show help")
	fmt.Println("      --version    Show version")
	fmt.Println("      --json       Output as JSON")
	fmt.Println("  -q, --quiet      Suppress output")
	fmt.Println("      --no-color   Disable color")
}
```

- [ ] **Step 2: Create fixture-no-nocolor — ignores NO_COLOR**

A binary that always emits ANSI color codes regardless of `NO_COLOR`.

- [ ] **Step 3: Create fixture-no-help — no --help flag**

A binary that exits 1 with "unknown flag" for `--help`.

- [ ] **Step 4: Create generate.go with go:generate directives**

```go
// testdata/fixtures/generate.go
package fixtures

//go:generate go build -o ../bin/fixture-good ./good/
//go:generate go build -o ../bin/fixture-no-nocolor ./no-nocolor/
//go:generate go build -o ../bin/fixture-no-help ./no-help/
```

- [ ] **Step 5: Write end-to-end test**

```go
// e2e_test.go
package main

import (
	"os/exec"
	"testing"
)

func TestE2E_GoodFixture(t *testing.T) {
	cmd := exec.Command("./cli-acs", "check", "--domain", "color", "--format", "json", "testdata/bin/fixture-good")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("cli-acs failed: %v\n%s", err, out)
	}
	// exit 0 = all pass
}

func TestE2E_NoNocolorFixture(t *testing.T) {
	cmd := exec.Command("./cli-acs", "check", "--domain", "color", "--format", "json", "testdata/bin/fixture-no-nocolor")
	_ = cmd.Run()
	if cmd.ProcessState.ExitCode() == 0 {
		t.Error("expected non-zero exit for fixture that fails NO_COLOR")
	}
}
```

- [ ] **Step 6: Build fixtures, build suite, run e2e tests**

```bash
go generate ./testdata/fixtures/
go build -o cli-acs ./cmd/cli-acs/
go test -v -tags e2e ./...
```

- [ ] **Step 7: Commit**

```bash
git add testdata/ e2e_test.go
git commit -s -m "feat: add fixture binaries and end-to-end tests"
```

---

### Task 13: Coverage Command + Validate Command

**Files:**
- Create: `cmd/cli-acs/coverage.go`
- Modify: `cmd/cli-acs/validate.go` (enhance from Task 9)

**Interfaces:**
- Consumes: spec file parser, domain registry, YAML schema
- Produces: `cli-acs coverage` command, enhanced `cli-acs validate` command

- [ ] **Step 1: Implement coverage.go**

Parses `spec/CLI_ACS_v1.0.md` for all criterion IDs via regex `[A-Z]{2,3}-\d+`. Cross-references against registered YAML + Go criteria and testcase.md files. Reports: missing tests, missing docs, orphan tests, orphan docs.

- [ ] **Step 2: Enhance validate.go**

Validates YAML against JSON Schema. Checks testcase.md frontmatter fields. Verifies `criterion_id` matches filename. Reports all errors.

- [ ] **Step 3: Run against current state**

```bash
./cli-acs coverage
./cli-acs validate internal/domains/color/testcases/
```

- [ ] **Step 4: Commit**

```bash
git add cmd/cli-acs/coverage.go cmd/cli-acs/validate.go
git commit -s -m "feat: add coverage and validate commands for spec-test-docs sync"
```
