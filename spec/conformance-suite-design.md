# CLI-ACS Conformance Suite — Design Document

**Date:** 2026-09-06
**Status:** Draft
**Spec:** CLI-ACS v1.0

---

## Table of Contents

1. [Purpose](#1-purpose)
2. [User Interface](#2-user-interface)
3. [Architecture](#3-architecture)
4. [Probe & Auto-Discovery](#4-probe--auto-discovery)
5. [Criterion System (Hybrid YAML + Go)](#5-criterion-system-hybrid-yaml--go)
6. [YAML DSL Specification](#6-yaml-dsl-specification)
7. [Go Criterion Interface](#7-go-criterion-interface)
8. [Testcase Documentation (testcase.md)](#8-testcase-documentation-testcasemd)
9. [Safety & Isolation](#9-safety--isolation)
10. [Report Generation](#10-report-generation)
11. [Self-Testing & Fixtures](#11-self-testing--fixtures)
12. [Coverage Enforcement](#12-coverage-enforcement)
13. [Dependencies](#13-dependencies)
14. [Directory Structure](#14-directory-structure)

---

## 1. Purpose

The CLI-ACS Conformance Suite is a Go binary (`cli-acs`) that accepts any CLI binary path, evaluates it against the CLI-ACS v1.0 criteria (95 criteria across 9 domains + TUI extension), and produces a comprehensive accessibility conformance report.

### Design Principles

1. **Zero config for end users.** Run `cli-acs check <binary>` and get a report. No manifest, no YAML, no setup.
2. **Hybrid criterion system.** Simple checks in declarative YAML (easy for non-developers to add/maintain). Complex checks in Go (for multi-step logic, parsing, heuristics).
3. **Every criterion documented.** Each criterion — YAML or Go — has a `testcase.md` with frontmatter linking it to the spec.
4. **Self-dogfooding.** The suite itself respects `NO_COLOR`, `TERM=dumb`, `--no-color`, `--plain`. It passes its own criteria.
5. **Safe by default.** Never run destructive subcommands. Process-level isolation with timeouts. Hardcoded safety blocklist.
6. **Extensible.** External YAML criteria via `--criteria-dir`. Domain self-registration via `init()`. Reporter interface for new output formats.

---

## 2. User Interface

### Primary Commands

```
cli-acs check <binary>                       # full evaluation, terminal report
cli-acs check --domain color <binary>        # single domain
cli-acs check --level A <binary>             # only Level A criteria
cli-acs check --auto-only <binary>           # skip SEMI/MANUAL (for CI)
cli-acs check --threshold AA <binary>        # exit 1 if AA not met
cli-acs check --format json <binary>         # JSON to stdout
cli-acs check --format markdown <binary>     # Markdown to stdout
cli-acs check --format html <binary>         # HTML to stdout
cli-acs check --output report.json <binary>  # write to file
cli-acs check --crosswalk <binary>           # include WCAG/508/EN 301 549 mapping
cli-acs check --plain <binary>               # screen-reader-friendly linear output
cli-acs check --quiet <binary>               # exit code only, no output
cli-acs check --skip CV-7,HD-12 <binary>     # skip specific criteria
cli-acs check --timeout 30s <binary>         # per-criterion timeout
cli-acs check --max-subcommands 20 <binary>  # cap subcommand discovery
cli-acs check --criteria-dir ./custom/ <binary>  # load external YAML criteria
cli-acs check --no-color <binary>            # suppress suite's own color

cli-acs validate [path]                      # lint YAML + testcase.md files
cli-acs new-criterion <ID>                   # scaffold new YAML + testcase.md
cli-acs coverage                             # verify spec ↔ tests ↔ docs sync
cli-acs version                              # suite version + spec version
```

### Exit Codes

| Code | Meaning |
|---|---|
| 0 | All criteria pass at the configured threshold level |
| 1 | One or more Level A criteria fail |
| 2 | Level A passes but one or more Level AA criteria fail |
| 3 | Level AA passes but one or more Level AAA criteria fail |
| 4 | Suite internal error (invalid binary path, permission denied, etc.) |

### Flags Summary

| Flag | Default | Description |
|---|---|---|
| `--domain` | all | Comma-separated domain filter (output,color,help,errors,input,environ,timing,i18n,lifecycle) |
| `--level` | all | Filter criteria by level (A, AA, AAA) |
| `--auto-only` | false | Run only `[AUTO]` criteria, skip `[SEMI]` and `[MANUAL]` |
| `--threshold` | A | Conformance level required for exit 0 |
| `--format` | terminal | Output format: terminal, json, markdown, html |
| `--output` | stdout | Write report to file instead of stdout |
| `--crosswalk` | false | Include WCAG/508/EN 301 549 mapping columns in report |
| `--plain` | false | Screen-reader-friendly output (no tables, no box-drawing, no color) |
| `--quiet` | false | Suppress all output, communicate only via exit code |
| `--skip` | none | Comma-separated criterion IDs to skip |
| `--timeout` | 10s | Per-criterion timeout |
| `--max-subcommands` | 10 | Maximum subcommands to test during discovery |
| `--criteria-dir` | none | Additional directory of YAML criteria to load |
| `--no-color` | false | Disable color in suite output |

---

## 3. Architecture

```
                                ┌─────────────────────┐
                                │   cmd/cli-acs/main   │
                                │   (cobra CLI)        │
                                └─────────┬───────────┘
                                          │
                              ┌───────────▼───────────┐
                              │   internal/engine/     │
                              │   runner.go            │
                              │   - loads config       │
                              │   - runs probe         │
                              │   - iterates domains   │
                              │   - collects results   │
                              │   - invokes reporter   │
                              └───┬───────────┬───────┘
                                  │           │
                    ┌─────────────▼──┐   ┌────▼────────────┐
                    │ internal/probe/ │   │ internal/report/ │
                    │ discovery.go    │   │ reporter.go      │
                    │ parsers.go      │   │ json.go          │
                    │ tty.go          │   │ markdown.go      │
                    │ exec.go         │   │ terminal.go      │
                    │ ansi.go         │   │ html.go          │
                    │ safety.go       │   └──────────────────┘
                    └─────────────────┘
                          │
            ┌─────────────▼──────────────┐
            │   internal/domains/        │
            │   registry.go              │
            │   yaml_loader.go           │
            │   criterion.go (interface) │
            ├────────────────────────────┤
            │   output/                  │
            │     domain.go  (init reg)  │
            │     checks.go  (Go checks) │
            │     testcases/             │
            │       OS-1.yaml            │
            │       OS-1.testcase.md     │
            ├────────────────────────────┤
            │   color/                   │
            │   help/                    │
            │   errors/                  │
            │   input/                   │
            │   environ/                 │
            │   timing/                  │
            │   i18n/                    │
            │   lifecycle/               │
            └────────────────────────────┘
```

### Flow

1. **Parse CLI flags** (cobra) → build `Config`
2. **Probe** the binary: run `<binary> --help`, `<binary> --version`, detect subcommands, flags, capabilities
3. **Load criteria** from all registered domains (YAML + Go), filtered by `--domain`, `--level`, `--auto-only`, `--skip`
4. **Execute criteria** sequentially within each domain, each in an isolated subprocess with timeout
5. **Collect results** into a `ConformanceReport` struct
6. **Render report** via the selected `Reporter` (terminal/JSON/Markdown/HTML)
7. **Exit** with appropriate code based on `--threshold`

---

## 4. Probe & Auto-Discovery

The probe package analyzes the binary under test before any criteria run, producing a `ProbeResult` that criteria use for preconditions and test setup.

### ProbeResult Structure

```go
type ProbeResult struct {
    BinaryPath     string
    BinaryName     string
    Version        string            // from --version output
    HelpText       string            // raw --help output
    Subcommands    []Subcommand      // discovered from help
    GlobalFlags    []Flag            // discovered from help
    HasColor       bool              // emits ANSI codes on TTY
    HasHelp        bool              // responds to --help
    HasVersion     bool              // responds to --version
    HasSubcommands bool              // git-style subcommand structure
    HasJSONFlag    bool              // supports --json
    HasQuietFlag   bool              // supports -q/--quiet
    HasColorFlag   bool              // supports --color/--no-color
    HasDryRunFlag  bool              // supports --dry-run
    HasNoInputFlag bool              // supports --no-input/--non-interactive
    HelpFormat     string            // detected format: cobra, argparse, clap, custom
    SampleErrors   []ErrorSample     // outputs from deliberately triggered errors
}

type Subcommand struct {
    Name     string
    HelpText string
    Flags    []Flag
}

type Flag struct {
    Short       string   // e.g., "-h"
    Long        string   // e.g., "--help"
    Description string
    TakesValue  bool
}

type ErrorSample struct {
    Trigger  string   // what was run
    Stdout   []byte
    Stderr   []byte
    ExitCode int
}
```

### Discovery Steps

1. **Run `<binary> --help`** — capture stdout, stderr, exit code. If exit 0 and stdout non-empty, parse for subcommands and flags.
2. **Run `<binary> --version`** — capture version string.
3. **Parse help text** using a cascade of parsers:
   - Cobra parser (looks for `Available Commands:`, `Flags:`, `Use "... --help"`)
   - Clap/Rust parser (looks for `USAGE:`, `OPTIONS:`, `SUBCOMMANDS:`)
   - Argparse/Python parser (looks for `positional arguments:`, `options:`)
   - Generic fallback (lines starting with `--` or `-` are flags; indented words after command-like headers are subcommands)
4. **Detect color** — run in a PTY and check for ANSI escape sequences in output.
5. **Detect common flags** — scan parsed flags for `--json`, `--quiet`, `--color`, `--no-color`, `--dry-run`, `--no-input`, `--non-interactive`.
6. **Discover subcommands** — up to `--max-subcommands` (default 10). For each, run `<binary> <subcommand> --help` to get subcommand-specific flags.
7. **Trigger errors** — run `<binary> --nonexistent-flag-a11y-probe` to capture error output format and exit code. Run `<binary> /tmp/cli-acs-nonexistent-<random>` as a nonexistent-file error trigger.

### Help Parser Cascade

```go
type HelpParser interface {
    Name() string
    CanParse(helpText string) bool
    Parse(helpText string) (*ParsedHelp, error)
}

// Tried in order; first CanParse() == true wins.
var parsers = []HelpParser{
    &CobraParser{},
    &ClapParser{},
    &ArgparseParser{},
    &GenericParser{},  // always returns true for CanParse
}
```

---

## 5. Criterion System (Hybrid YAML + Go)

Every criterion in the suite implements the same `Criterion` interface, regardless of whether it's defined in YAML or Go. The domain's `init()` function registers both:

```go
// domains/color/domain.go
func init() {
    d := &ColorDomain{}
    // Register Go-defined criteria
    d.AddCheck(&CV1Check{})
    d.AddCheck(&CV9Check{})
    // YAML criteria loaded automatically from testcases/*.yaml
    domains.Register("color", d)
}
```

### When to Use YAML vs Go

| Use YAML when... | Use Go when... |
|---|---|
| Check is env var + flag + output scan | Check requires parsing structured output |
| Single exec → assert pattern | Check needs multi-step logic (run A, then run B, compare) |
| Non-developer should be able to add it | Check needs PTY allocation |
| Test method is straightforward regex | Check needs timing measurement |
| Examples: CV-2, CV-3, CV-5, HD-1, HD-3, EF-1, EF-2 | Examples: CV-1, CV-9, OS-4, HD-5, TM-6, II-4 |

---

## 6. YAML DSL Specification

### Schema

```yaml
# Required fields
id: "CV-2"                      # criterion ID, must match spec
domain: "color"                  # domain name
level: "A"                       # A, AA, or AAA
testability: "AUTO"              # AUTO, SEMI, or MANUAL
spec_version: "1.0"              # CLI-ACS spec version
name: "NO_COLOR Support"         # human-readable name

# Optional: skip this check if probe doesn't detect capability
requires:                        # preconditions from probe result
  - uses_color                   # tool emits ANSI on TTY
  # Available: has_help, has_version, has_subcommands, uses_color,
  #            has_json_flag, has_quiet_flag, has_color_flag,
  #            has_dry_run_flag, has_no_input_flag

# Test steps — executed in order
steps:
  - name: "Color suppressed with NO_COLOR=1"
    exec:
      args: ["--help"]           # arguments to pass
      env:                       # environment overrides
        NO_COLOR: "1"
      timeout: "10s"             # per-step timeout (optional)
      # Execution modes (pick one):
      # exec:      → stdout piped (non-TTY) — default
      # exec_pty:  → stdout in PTY (TTY)
      # exec_pipe: → stdout piped, explicit (same as exec)
    assert:
      exit_code: 0               # exact exit code
      # exit_code_nonzero: true  # any non-zero (alternative)
      stdout_not_matches:        # regex patterns that must NOT appear
        - '\x1b\[(3[0-7]|9[0-7]|38;[25];)'   # foreground color
        - '\x1b\[(4[0-7]|10[0-7]|48;[25];)'   # background color
      # stdout_matches:          # regex patterns that MUST appear
      # stderr_matches:          # same for stderr
      # stderr_not_matches:
      # stderr_empty: true       # stderr must be empty
      # stdout_empty: true       # stdout must be empty
      # stdout_contains:         # substring match (simpler than regex)
      # stderr_contains:
      # stdout_not_contains:
      # timing_under_ms: 500     # execution time threshold

  - name: "Color present when NO_COLOR unset (confirms tool uses color)"
    exec_pty:
      args: ["--help"]
      env: {}                    # clean env, no NO_COLOR
    assert:
      stdout_matches:
        - '\x1b\['              # any ANSI escape sequence

  - name: "Empty NO_COLOR treated as unset"
    exec_pty:
      args: ["--help"]
      env:
        NO_COLOR: ""             # empty = unset per spec
    assert:
      stdout_matches:
        - '\x1b\['              # color should still be present

# How to interpret results
result_on_all_pass: "Supports"
result_on_any_fail: "Does Not Support"
# For SEMI criteria:
# result_on_all_pass: "Supports (needs review)"
# human_review_prompt: "Verify that the error message is understandable"
```

### Execution Modes

| Mode | Syntax | Behavior |
|---|---|---|
| `exec` | `exec:` or `exec_pipe:` | Run binary with stdout/stderr as pipes (non-TTY). Default. |
| `exec_pty` | `exec_pty:` | Run binary with stdout connected to a PTY (TTY). Required for testing TTY-aware behavior. |

### Assert Operators

| Operator | Type | Description |
|---|---|---|
| `exit_code` | int | Exact exit code match |
| `exit_code_nonzero` | bool | Any non-zero exit code |
| `stdout_matches` | []regex | All patterns must appear in stdout |
| `stdout_not_matches` | []regex | No pattern may appear in stdout |
| `stderr_matches` | []regex | All patterns must appear in stderr |
| `stderr_not_matches` | []regex | No pattern may appear in stderr |
| `stdout_contains` | []string | Substring presence check |
| `stdout_not_contains` | []string | Substring absence check |
| `stderr_contains` | []string | Same for stderr |
| `stderr_not_contains` | []string | Same for stderr |
| `stdout_empty` | bool | stdout must be empty |
| `stderr_empty` | bool | stderr must be empty |
| `timing_under_ms` | int | Total execution time must be under N ms |

### Requires (Preconditions)

| Precondition | Probe field | Meaning |
|---|---|---|
| `has_help` | `ProbeResult.HasHelp` | Binary responds to --help |
| `has_version` | `ProbeResult.HasVersion` | Binary responds to --version |
| `has_subcommands` | `ProbeResult.HasSubcommands` | Binary has git-style subcommands |
| `uses_color` | `ProbeResult.HasColor` | Binary emits ANSI codes on TTY |
| `has_json_flag` | `ProbeResult.HasJSONFlag` | Binary supports --json |
| `has_quiet_flag` | `ProbeResult.HasQuietFlag` | Binary supports -q/--quiet |
| `has_color_flag` | `ProbeResult.HasColorFlag` | Binary supports --color/--no-color |
| `has_dry_run_flag` | `ProbeResult.HasDryRunFlag` | Binary supports --dry-run |
| `has_no_input_flag` | `ProbeResult.HasNoInputFlag` | Binary supports --no-input |

When a precondition is not met, the criterion result is `Not Applicable` (not a failure).

### JSON Schema

A `schemas/criterion.schema.json` file validates all YAML criteria in CI and enables editor autocomplete. The `cli-acs validate` command uses this schema.

---

## 7. Go Criterion Interface

```go
package domains

import (
    "context"
    "time"
)

// Level represents a CLI-ACS conformance level.
type Level int

const (
    LevelA   Level = iota
    LevelAA
    LevelAAA
)

// Testability represents how a criterion can be tested.
type Testability int

const (
    Auto   Testability = iota
    Semi
    Manual
)

// Outcome represents the result of evaluating a criterion.
type Outcome int

const (
    Supports          Outcome = iota
    PartiallySupports
    DoesNotSupport
    NotApplicable
    NotEvaluated
    Error
    Timeout
)

// Evidence captures proof of a test result.
type Evidence struct {
    Command  string // the exact command run
    Env      map[string]string
    Stdout   string // captured stdout (truncated to 4KB)
    Stderr   string // captured stderr (truncated to 4KB)
    ExitCode int
    Duration time.Duration
    Note     string // human-readable observation
}

// Result is the outcome of a single criterion evaluation.
type Result struct {
    ID          string
    Name        string
    Domain      string
    Level       Level
    Testability Testability
    Outcome     Outcome
    Evidence    []Evidence
    Remarks     string      // required if Outcome != Supports
    Duration    time.Duration
    NeedsReview bool        // true for SEMI results flagged for human review
    SpecVersion string
}

// Criterion is the interface every check implements (YAML or Go).
type Criterion interface {
    ID() string
    Name() string
    Domain() string
    Level() Level
    Testability() Testability
    SpecVersion() string
    Precondition(probe *ProbeResult) bool
    Run(ctx context.Context, binary string, probe *ProbeResult) *Result
}

// Domain groups related criteria and self-registers via init().
type Domain interface {
    Name() string
    Criteria() []Criterion
}
```

### Example Go Criterion (complex check)

```go
// domains/color/cv1_check.go
package color

// CV1Check tests that color is not the sole information channel.
// This requires running the binary twice (with and without NO_COLOR)
// and comparing whether information is lost.
type CV1Check struct{}

func (c *CV1Check) ID() string          { return "CV-1" }
func (c *CV1Check) Name() string        { return "Color Not Sole Information Channel" }
func (c *CV1Check) Domain() string      { return "color" }
func (c *CV1Check) Level() Level        { return LevelA }
func (c *CV1Check) Testability() Testability { return Semi }
func (c *CV1Check) SpecVersion() string { return "1.0" }

func (c *CV1Check) Precondition(probe *ProbeResult) bool {
    return probe.HasColor // only test if the tool uses color
}

func (c *CV1Check) Run(ctx context.Context, binary string, probe *ProbeResult) *Result {
    // 1. Run with color (PTY)
    colorOut, colorEvidence := execPTY(ctx, binary, []string{"--help"}, nil)

    // 2. Run without color
    noColorOut, noColorEvidence := execPipe(ctx, binary, []string{"--help"},
        map[string]string{"NO_COLOR": "1"})

    // 3. Strip ANSI from colored output
    stripped := stripANSI(colorOut.Stdout)

    // 4. Compare: if stripped colored output == no-color output,
    //    then no information is lost when color is removed.
    //    This is a heuristic — flag for human review.
    similarity := computeSimilarity(stripped, noColorOut.Stdout)

    result := &Result{
        ID: "CV-1", Name: c.Name(), Domain: "color",
        Level: LevelA, Testability: Semi, SpecVersion: "1.0",
        Evidence:    []Evidence{colorEvidence, noColorEvidence},
        NeedsReview: true,
    }

    if similarity > 0.95 {
        result.Outcome = Supports
        result.Remarks = "Output structure preserved when color is removed"
    } else {
        result.Outcome = PartiallySupports
        result.Remarks = fmt.Sprintf(
            "%.0f%% content similarity between colored and no-color output — "+
            "review for information loss", similarity*100)
    }
    return result
}
```

---

## 8. Testcase Documentation (testcase.md)

Every criterion — YAML or Go — has a companion `testcase.md` file in its domain's `testcases/` directory.

### Template

```markdown
---
criterion_id: CV-2
domain: color
level: A
testability: AUTO
spec_version: "1.0"
type: yaml           # yaml or go
last_reviewed: 2026-09-06
---

# CV-2: NO_COLOR Support

## What This Tests

Whether the tool suppresses ANSI color codes when the NO_COLOR
environment variable is set to a non-empty value, per the
no-color.org specification.

## How The Test Works

1. Runs the binary with `NO_COLOR=1` and `--help`, scans stdout for
   ANSI SGR color sequences (foreground and background).
2. Runs without `NO_COLOR` in a PTY to confirm the tool uses color
   by default (if it never uses color, this criterion is trivially met).
3. Runs with `NO_COLOR=""` (empty) to confirm empty is treated as unset.

## Pass Criteria

- No ANSI color escape sequences in stdout when `NO_COLOR=1`
- Color sequences present when `NO_COLOR` is unset (on TTY)
- Color sequences present when `NO_COLOR=""` (empty)

## Fail Criteria

- ANSI color escape sequences found in output despite `NO_COLOR=1`
- Tool treats empty `NO_COLOR=` as set (suppresses color on empty)

## References

- [no-color.org](https://no-color.org/)
- CLI-ACS v1.0 spec, Section 4.2, CV-2
```

### Enforcement

The `cli-acs validate` command checks:

1. Frontmatter fields are present and valid
2. `criterion_id` matches the filename
3. `type` matches whether a `.yaml` or Go check exists
4. All required sections are present (What, How, Pass, Fail)
5. `spec_version` matches the suite's target spec version

The `cli-acs coverage` command checks:

1. Every criterion ID in the spec has a YAML or Go check
2. Every criterion ID has a `testcase.md`
3. No orphan tests exist for removed criteria
4. No orphan `testcase.md` files exist without a corresponding check

---

## 9. Safety & Isolation

### Hard Safety Guarantees

1. **Isolated filesystem.** All tests run with `HOME` and `PWD` set to a fresh temp directory. The binary under test cannot access the user's real home directory or current working directory.

2. **Clean environment.** Each criterion gets a fresh `exec.Cmd` with `Env` set explicitly — not inheriting the suite's process environment. The base env includes only:
   ```
   PATH=<system paths>
   HOME=/tmp/cli-acs-<random>/home
   TMPDIR=/tmp/cli-acs-<random>/tmp
   LANG=C.UTF-8
   TERM=xterm-256color
   ```
   Each test adds its specific overrides (e.g., `NO_COLOR=1`).

3. **Per-criterion timeout.** Default 10s, configurable via `--timeout`. Timing criteria get 30s. On timeout, the process is killed and the result is `Timeout`.

4. **Subcommand blocklist.** The suite never invokes subcommands matching these patterns (case-insensitive):
   ```
   delete, rm, remove, destroy, purge, drop, format, init, reset,
   clean, wipe, nuke, truncate, uninstall, erase
   ```
   If the only way to test a criterion is through a blocked subcommand, the result is `Not Evaluated` with a note.

5. **Read-only invocations only.** The suite only invokes: `--help`, `--version`, `-h`, `-V`, discovered flags that appear read-only (no value taking), and error-trigger patterns (invalid flags, nonexistent file paths).

6. **Process groups.** Subprocesses run in their own process group (`Setpgid: true`). On timeout, the entire group is signaled, preventing orphaned child processes.

7. **No network by default.** Tests that would require network access (EC-10 proxy) are skipped unless `--allow-network` is explicitly passed.

### What the Suite Does NOT Guarantee

- Container-level isolation. Users testing untrusted binaries should run the suite inside a container or VM.
- Prevention of all side effects. A sufficiently adversarial binary can still write to `/tmp`, make network calls, or consume system resources within the timeout window.

---

## 10. Report Generation

### Reporter Interface

```go
type Reporter interface {
    Name() string
    FileExtension() string
    Render(report *ConformanceReport, w io.Writer) error
}
```

New reporters are added by implementing this interface and registering via `init()`. The `--format` flag auto-discovers registered reporters.

### ConformanceReport Structure

```go
type ConformanceReport struct {
    // Header
    CLIACSVersion  string
    SuiteVersion   string
    Product        ProductInfo
    ReportDate     time.Time
    Environment    TestEnvironment

    // Results
    Results        []Result
    DomainSummary  map[string]DomainSummary
    LevelSummary   map[Level]LevelSummary
    DisabilitySummary map[string]DisabilitySummary // keyed by category name
    OverallLevel   Level  // highest level fully passed

    // Config
    Threshold      Level
    SkippedCriteria []string
    AutoOnly       bool
    IncludeCrosswalk bool
}

type ProductInfo struct {
    Name    string
    Path    string
    Version string  // from --version probe
}

type TestEnvironment struct {
    OS       string
    Terminal string
    Shell    string
    Arch     string
    SuiteVer string
}

type DomainSummary struct {
    Domain      string
    Total       int
    Supports    int
    Partial     int
    Fails       int
    NA          int
    NotEval     int
    LevelA      string  // "pass" or "fail"
    LevelAA     string
    LevelAAA    string
}
```

### Report Formats

**Terminal** — colored summary table with per-domain breakdown. Uses `[PASS]`, `[PARTIAL]`, `[FAIL]`, `[N/A]`, `[SKIP]` prefixes for screen reader compatibility. Respects `NO_COLOR`, `TERM=dumb`, `--no-color`, `--plain`.

**JSON** — machine-readable, matches the schema from CLI-ACS v1.0 Section 7.3. Includes full evidence for each criterion (commands run, stdout/stderr captured, exit codes, timing).

**Markdown** — matches the CLI-ACR template from CLI-ACS v1.0 Section 7. Suitable for embedding in PRs, wikis, or documentation.

**HTML** — standalone HTML file with embedded CSS. Includes expandable evidence sections, a summary dashboard, and per-domain drill-down. No external dependencies (single file, no CDN links).

### Terminal Output Example (--plain mode)

```
CLI-ACS Conformance Report
Product: gh v2.87.3
Date: 2026-09-06
Spec: CLI-ACS v1.0
Suite: cli-acs v0.1.0

Overall: Partial Level A (4 Level A failures)

Domain: Output Structure & Content
  [PASS] OS-1 Standard Stream Separation
  [PASS] OS-2 Meaningful Exit Codes
  [PASS] OS-3 Clean Piped Output
  [PASS] OS-4 Linear Reading Order
  [PASS] OS-5 No ASCII Art Sole Info
  [PASS] OS-6 Machine-Readable Output
  [PARTIAL] OS-7 Quiet Mode
    Remark: No -q/--quiet flag. -q is aliased to --jq.
  [PARTIAL] OS-8 Output Width Awareness
    Remark: Help output does not respect COLUMNS=40.
  [PASS] OS-9 Pager Support
  [FAIL] OS-10 Plain Text Alternative
    Remark: No --plain or --accessible flag.

Domain: Color & Visual Presentation
  [PASS] CV-1 Color Not Sole Information Channel
  [PASS] CV-2 NO_COLOR Support
  [FAIL] CV-3 --no-color Flag
    Remark: --no-color returns "unknown flag". No per-invocation flag.
  ...

Summary by disability category:
  Visual: 45/52 criteria met (86.5%)
  Motor: 22/24 criteria met (91.7%)
  Cognitive: 38/44 criteria met (86.4%)
  Vestibular: 2/2 criteria met (100.0%)
```

---

## 11. Self-Testing & Fixtures

### Fixture Binaries

The `testdata/fixtures/` directory contains small Go programs built via `go generate` that exhibit specific accessibility behaviors:

| Fixture | Purpose | Behavior |
|---|---|---|
| `fixture-good` | Passes all AUTO criteria | Respects NO_COLOR, TERM=dumb, has --help/-h/--version, uses exit codes correctly, structured output, --json, --quiet, --no-color, --color=WHEN |
| `fixture-no-nocolor` | Fails CV-2 | Ignores NO_COLOR, always emits ANSI |
| `fixture-no-help` | Fails HD-1 | No --help or -h flag |
| `fixture-color-only` | Fails CV-1 | Uses red/green for pass/fail with no text indicator |
| `fixture-bad-exit` | Fails EF-2 | Exits 0 on error |
| `fixture-hangs` | Tests timeout handling | Sleeps 60s on invocation |
| `fixture-stderr-mix` | Fails OS-1 | Sends errors to stdout |
| `fixture-ansi-pipe` | Fails OS-3 | Emits ANSI codes even when piped |
| `fixture-jargon-errors` | Fails EF-8 | Prints stack traces on errors |
| `fixture-interactive` | Tests II-2, II-3 | Prompts on stdin even when not TTY |
| `fixture-slow-help` | Fails TM-6 | --help takes 2 seconds |
| `fixture-subcommands` | Tests subcommand discovery | Has 5 subcommands with distinct behaviors |

Each fixture is a single `main.go` with clear documentation of what it tests.

### Suite Self-Tests

```go
func TestCV2_GoodFixture(t *testing.T) {
    result := runCriterion("CV-2", fixtureGoodPath)
    assert(t, result.Outcome == Supports)
}

func TestCV2_NoNocolorFixture(t *testing.T) {
    result := runCriterion("CV-2", fixtureNoNocolorPath)
    assert(t, result.Outcome == DoesNotSupport)
}
```

Every YAML and Go criterion has at least two fixture tests: one that should pass and one that should fail. This ensures criteria actually detect violations.

### Timing Criteria Handling

Timing-sensitive criteria (TM-6) run 5 times and report the median. The test environment is documented in the report. CI environments may override timing thresholds via `--timeout`.

---

## 12. Coverage Enforcement

The `cli-acs coverage` command runs in CI and enforces:

1. **Spec → Test:** Every criterion ID found in the spec (regex `[A-Z]{2,3}-\d+`) has a corresponding YAML or Go check.
2. **Spec → Docs:** Every criterion ID has a `testcase.md` file.
3. **Test → Spec:** No orphan YAML or Go checks exist for criterion IDs not in the spec.
4. **Docs → Spec:** No orphan `testcase.md` files exist.
5. **Frontmatter consistency:** `criterion_id` in testcase.md matches filename; `spec_version` matches suite target.
6. **YAML schema validation:** All YAML files pass JSON Schema validation.

This runs as a CI check on every PR to prevent drift.

---

## 13. Dependencies

Minimal dependency footprint:

| Dependency | Purpose | Justification |
|---|---|---|
| `github.com/spf13/cobra` | CLI framework | Industry standard for Go CLIs |
| `github.com/creack/pty` | PTY allocation | Required for TTY simulation (CV-4, OS-3, TM-1) |
| `gopkg.in/yaml.v3` | YAML parsing | Criterion DSL loading |
| `github.com/santhosh-tekuri/jsonschema` | JSON Schema validation | YAML criterion validation |
| Go stdlib | Everything else | `os/exec`, `regexp`, `encoding/json`, `html/template`, `time`, `testing` |

No test frameworks beyond `testing`. No assertion libraries. No color libraries (the suite implements its own ANSI output respecting NO_COLOR).

---

## 14. Directory Structure

```
cli-acs/
├── cmd/
│   └── cli-acs/
│       └── main.go                    # cobra CLI entrypoint
├── internal/
│   ├── engine/
│   │   ├── runner.go                  # orchestrates probe → criteria → report
│   │   ├── config.go                  # CLI config, thresholds, skip lists
│   │   └── result.go                  # Result, Evidence, Outcome types
│   ├── probe/
│   │   ├── discovery.go               # top-level probe orchestration
│   │   ├── parsers.go                 # help text parsers (cobra, clap, argparse, generic)
│   │   ├── tty.go                     # PTY allocation and TTY simulation
│   │   ├── exec.go                    # isolated subprocess execution
│   │   ├── ansi.go                    # ANSI detection, classification, stripping
│   │   └── safety.go                  # subcommand blocklist, env isolation
│   ├── domains/
│   │   ├── registry.go                # domain registration + discovery
│   │   ├── criterion.go               # Criterion interface definition
│   │   ├── yaml_criterion.go          # YAMLCriterion implements Criterion
│   │   ├── yaml_loader.go             # loads and validates YAML files
│   │   ├── output/
│   │   │   ├── domain.go              # init() self-registration
│   │   │   ├── checks.go              # Go-defined complex checks
│   │   │   └── testcases/
│   │   │       ├── OS-1.yaml
│   │   │       ├── OS-1.testcase.md
│   │   │       ├── OS-2.yaml
│   │   │       ├── OS-2.testcase.md
│   │   │       └── ...
│   │   ├── color/
│   │   │   ├── domain.go
│   │   │   ├── checks.go
│   │   │   └── testcases/
│   │   │       ├── CV-1.testcase.md   # Go check, no YAML
│   │   │       ├── CV-2.yaml
│   │   │       ├── CV-2.testcase.md
│   │   │       └── ...
│   │   ├── help/
│   │   │   ├── domain.go
│   │   │   ├── checks.go
│   │   │   └── testcases/
│   │   ├── errors/
│   │   │   └── ...
│   │   ├── input/
│   │   │   └── ...
│   │   ├── environ/
│   │   │   └── ...
│   │   ├── timing/
│   │   │   └── ...
│   │   ├── i18n/
│   │   │   └── ...
│   │   └── lifecycle/
│   │       └── ...
│   └── report/
│       ├── reporter.go                # Reporter interface
│       ├── json.go                    # JSON renderer
│       ├── markdown.go                # Markdown renderer
│       ├── terminal.go                # Terminal renderer (colored + plain)
│       └── html.go                    # Standalone HTML renderer
├── schemas/
│   └── criterion.schema.json          # JSON Schema for YAML validation
├── testdata/
│   └── fixtures/
│       ├── generate.go                # go:generate build script
│       ├── good/main.go               # passes all criteria
│       ├── no-nocolor/main.go         # ignores NO_COLOR
│       ├── no-help/main.go            # no --help flag
│       ├── color-only/main.go         # color as sole info channel
│       ├── bad-exit/main.go           # exits 0 on error
│       ├── hangs/main.go              # sleeps 60s
│       ├── stderr-mix/main.go         # errors to stdout
│       ├── ansi-pipe/main.go          # ANSI when piped
│       ├── jargon-errors/main.go      # stack traces on error
│       ├── interactive/main.go        # prompts when not TTY
│       ├── slow-help/main.go          # 2s --help
│       └── subcommands/main.go        # 5 subcommands
├── spec/                              # specification documents (existing)
│   ├── CLI_ACS_v1.0.md
│   ├── color-and-visual-presentation.md
│   └── conformance-suite-design.md    # this document
├── go.mod
├── go.sum
├── Makefile
└── README.md
```
