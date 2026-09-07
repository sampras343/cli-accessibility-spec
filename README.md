# CLI Accessibility Conformance Specification (CLI-ACS)

A CLI-native accessibility standard that defines measurable criteria for evaluating how well command-line interface tools serve users with disabilities — plus a conformance suite that evaluates any binary against those criteria and produces a detailed report.

## Why CLI-ACS?

Existing accessibility frameworks (WCAG 2.2, EN 301 549, Section 508) apply to CLI tools in principle but provide little CLI-specific guidance. Current conformance reports (VPATs) for CLI tools are often meaningless — marking criteria like "Pointer Gestures" and "Dragging Movements" as "Supports" for text-based tools, or marking everything "Not Applicable."

CLI-ACS addresses this gap with 98 criteria native to the CLI context, organized by functional domain, mapped to disability categories, and tagged with testability classifications.

## Quick Start

### Install

```bash
go install github.com/sampras343/cli-accessibility-spec/cmd/cli-acs@latest
```

Or build from source:

```bash
git clone https://github.com/sampras343/cli-accessibility-spec.git
cd cli-accessibility-spec
go build -o cli-acs ./cmd/cli-acs/
```

### Check Any Binary

```bash
# Evaluate a binary — that's it. No config, no manifest, no setup.
cli-acs check /usr/bin/gh

# Check a binary on your PATH
cli-acs check $(which kubectl)

# Check your own tool
cli-acs check ./my-cli-tool
```

The suite auto-discovers the binary's capabilities by parsing `--help`, probing flags, and detecting color output. You get a conformance report in your terminal immediately.

## Usage Guide

### Basic Usage

```bash
cli-acs check <binary>
```

This runs all applicable criteria against the binary and prints a terminal report showing pass/fail per criterion, grouped by domain.

Example output:

```
CLI-ACS Conformance Report
==========================
Product: gh v2.87.3
Date:    2026-09-06
Spec:    CLI-ACS v1.0

Overall: Level A (4 criteria evaluated)

Domain: Color & Visual Presentation
  [PASS] CV-2 NO_COLOR Support
  [FAIL] CV-3 --no-color Flag
    Remark: --no-color returns "unknown flag". No per-invocation flag.
  [PASS] CV-4 TTY-Aware Color
  [PASS] CV-5 TERM=dumb Respect
```

### Report Formats

The suite produces reports in four formats:

```bash
# Terminal output (default) — colored, with [PASS]/[FAIL] prefixes
cli-acs check ./my-tool

# JSON — machine-readable, for CI pipelines and automation
cli-acs check --format json ./my-tool

# Markdown — for PRs, wikis, documentation
cli-acs check --format markdown ./my-tool

# HTML — standalone file with embedded CSS, expandable evidence sections
cli-acs check --format html ./my-tool

# Write to a file instead of stdout
cli-acs check --format html --output report.html ./my-tool
cli-acs check --format json --output report.json ./my-tool
cli-acs check --format markdown --output report.md ./my-tool
```

### Filtering What Gets Checked

```bash
# Check only a specific domain
cli-acs check --domain color ./my-tool
cli-acs check --domain help ./my-tool
cli-acs check --domain errors ./my-tool

# Check only Level A criteria (minimum accessibility)
cli-acs check --level A ./my-tool

# Check only Level A + AA (standard conformance target)
cli-acs check --level AA ./my-tool

# Skip specific criteria you've decided to defer
cli-acs check --skip CV-7,HD-12 ./my-tool

# Run only fully automated checks (skip SEMI and MANUAL)
cli-acs check --auto-only ./my-tool
```

Available domains: `output`, `color`, `help`, `errors`, `input`, `environ`, `timing`, `i18n`, `lifecycle`

### CI/CD Integration

Use `cli-acs` as a quality gate in your CI pipeline:

```bash
# Fail the build if Level A criteria are not met
cli-acs check --threshold A --format json --output a11y-report.json ./my-tool

# Fail on Level AA (stricter)
cli-acs check --threshold AA --quiet ./my-tool
```

**Exit codes:**

| Code | Meaning |
|---|---|
| 0 | All criteria pass at the configured threshold level |
| 1 | One or more Level A criteria fail |
| 2 | Level A passes but one or more Level AA criteria fail |
| 3 | Level AA passes but one or more Level AAA criteria fail |
| 4 | Suite internal error (invalid binary, permission denied) |

**GitHub Actions example:**

```yaml
name: Accessibility Check
on: [push, pull_request]

jobs:
  a11y:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'

      - name: Install cli-acs
        run: go install github.com/sampras343/cli-accessibility-spec/cmd/cli-acs@latest

      - name: Build your tool
        run: go build -o my-tool ./cmd/my-tool/

      - name: Check accessibility
        run: cli-acs check --threshold AA --format json --output a11y-report.json ./my-tool

      - name: Upload report
        if: always()
        uses: actions/upload-artifact@v4
        with:
          name: accessibility-report
          path: a11y-report.json
```

### Screen-Reader-Friendly Output

The suite itself is accessible — it practices what it preaches:

```bash
# Plain mode: no color, no tables, no box-drawing — linear text with prefixes
cli-acs check --plain ./my-tool

# Respects NO_COLOR (just like it tests for)
NO_COLOR=1 cli-acs check ./my-tool

# Respects TERM=dumb
TERM=dumb cli-acs check ./my-tool

# Explicit --no-color flag
cli-acs check --no-color ./my-tool
```

### Crosswalk to WCAG / Section 508 / EN 301 549

For procurement or regulatory compliance, include the standards crosswalk in your report:

```bash
cli-acs check --crosswalk --format markdown --output report.md ./my-tool
```

This adds mapping columns showing how each CLI-ACS criterion maps to WCAG 2.2 Success Criteria, EN 301 549 clauses, and Section 508 requirements.

### Advanced Options

```bash
# Set per-criterion timeout (default 10s)
cli-acs check --timeout 30s ./my-tool

# Limit subcommand discovery (default 10, useful for large tools like kubectl)
cli-acs check --max-subcommands 5 ./my-tool

# Load additional custom criteria from a directory
cli-acs check --criteria-dir ./my-custom-criteria/ ./my-tool
```

### All Flags Reference

```
cli-acs check <binary> [flags]

Flags:
      --auto-only               Run only AUTO criteria, skip SEMI and MANUAL
      --criteria-dir string     Additional directory of YAML criteria to load
      --crosswalk               Include WCAG/508/EN 301 549 mapping in report
      --domain string           Comma-separated domain filter
      --format string           Output format: terminal, json, markdown, html (default "terminal")
  -h, --help                    help for check
      --level string            Filter by level: A, AA, AAA
      --max-subcommands int     Max subcommands to test (default 10)
      --no-color                Disable color in suite output
      --output string           Write report to file
      --plain                   Screen-reader-friendly linear output
      --quiet                   Exit code only, no output
      --skip string             Comma-separated criterion IDs to skip
      --threshold string        Conformance level required for exit 0 (default "A")
      --timeout duration        Per-criterion timeout (default 10s)
```

## Other Commands

### Validate Criteria Files

Lint YAML criterion definitions and testcase.md documentation for correctness:

```bash
# Validate all criteria in a directory
cli-acs validate internal/domains/color/testcases/

# Checks: YAML parsing, required fields, field value enums,
# testcase.md frontmatter, required sections, filename conventions
```

### Coverage Check

Verify that the spec, test implementations, and documentation are in sync:

```bash
cli-acs coverage --spec spec/CLI_ACS_v1.0.md

# Reports: missing tests, missing docs, orphan tests, orphan docs
```

### Version

```bash
cli-acs version
# cli-acs suite version: 1.0.0
# CLI-ACS spec version: 1.0.0
```

## How Auto-Discovery Works

When you run `cli-acs check <binary>`, the suite probes the binary before testing:

1. **Runs `<binary> --help`** — parses help text to extract subcommands and flags
2. **Runs `<binary> --version`** — captures version string for the report
3. **Detects help format** — recognizes Cobra (Go), Clap (Rust), argparse (Python), and generic formats
4. **Detects color usage** — checks if the binary emits ANSI escape codes
5. **Discovers common flags** — scans for `--json`, `--quiet`, `--color`, `--no-color`, `--dry-run`, `--no-input`
6. **Probes subcommands** — runs `--help` on up to 10 discovered subcommands (skips destructive ones like `delete`, `rm`)
7. **Triggers errors** — runs with an invalid flag to capture error output format

This auto-discovery feeds into criterion preconditions. For example, CV-2 (NO_COLOR Support) only runs if the binary actually uses color. If a criterion's precondition isn't met, it reports "Not Applicable" — not a failure.

### Safety Guarantees

The suite never runs destructive commands. Specifically:

- Only invokes `--help`, `--version`, `-h`, and error-trigger patterns
- Hardcoded blocklist: never runs subcommands containing `delete`, `rm`, `remove`, `destroy`, `purge`, `drop`, `format`, `init`, `reset`, `clean`, `wipe`, `nuke`, `uninstall`, `erase`
- Each criterion runs in an isolated environment (temp HOME, clean env vars, process group with timeout)
- Default 10-second timeout per criterion — a hanging binary won't block the suite

## Adding Custom Criteria

You can extend the suite with your own criteria using YAML:

```yaml
# my-criteria/ORG-1.yaml
id: "ORG-1"
domain: "custom"
level: "A"
testability: "AUTO"
spec_version: "1.0"
name: "Must support --org flag"
steps:
  - name: "Check --org flag exists"
    exec:
      args: ["--help"]
    assert:
      stdout_contains:
        - "--org"
result_on_all_pass: "Supports"
result_on_any_fail: "Does Not Support"
```

Then run with:

```bash
cli-acs check --criteria-dir ./my-criteria/ ./my-tool
```

### YAML DSL Reference

Each YAML criterion supports these execution modes and assertions:

**Execution modes:**
- `exec:` — run binary with stdout piped (non-TTY)
- `exec_pty:` — run binary in a PTY (TTY mode)

**Assertions:**

| Assertion | Type | Description |
|---|---|---|
| `exit_code` | int | Exact exit code match |
| `exit_code_nonzero` | bool | Any non-zero exit code |
| `stdout_matches` | list of regex | All patterns must appear in stdout |
| `stdout_not_matches` | list of regex | No pattern may appear in stdout |
| `stderr_matches` | list of regex | All patterns must appear in stderr |
| `stderr_not_matches` | list of regex | No pattern may appear in stderr |
| `stdout_contains` | list of strings | Substring presence check |
| `stdout_not_contains` | list of strings | Substring absence check |
| `stderr_contains` | list of strings | Substring presence in stderr |
| `stderr_not_contains` | list of strings | Substring absence in stderr |
| `stdout_empty` | bool | stdout must be empty |
| `stderr_empty` | bool | stderr must be empty |
| `timing_under_ms` | int | Execution must complete under N milliseconds |

**Preconditions** (skip criterion if binary lacks a capability):

```yaml
requires:
  - uses_color        # binary emits ANSI codes
  - has_help          # binary responds to --help
  - has_version       # binary responds to --version
  - has_subcommands   # binary has git-style subcommands
  - has_json_flag     # binary supports --json
  - has_quiet_flag    # binary supports -q/--quiet
  - has_color_flag    # binary supports --color/--no-color
  - has_dry_run_flag  # binary supports --dry-run
  - has_no_input_flag # binary supports --no-input
```

Each criterion should also have a companion `testcase.md` documenting what it tests and why. Run `cli-acs validate` to check your criteria are well-formed.

## Specification

The conformance criteria are defined in the CLI-ACS specification:

- [CLI-ACS v1.0 Specification](spec/CLI_ACS_v1.0.md) — 98 criteria across 9 domains + TUI extension
- [Color & Visual Presentation (detailed)](spec/color-and-visual-presentation.md) — deep-dive on terminal color tiers, SGR attributes, color vision deficiency, and contrast
- [Conformance Suite Design](spec/conformance-suite-design.md) — architecture, YAML DSL, probe system, safety model

### Criteria Summary

| Level | Core CLI | TUI Extension | Total |
|---|---|---|---|
| A (Minimum) | 28 | 5 | 33 |
| AA (Standard) | 40 | 6 | 46 |
| AAA (Enhanced) | 16 | 3 | 19 |
| **Total** | **84** | **14** | **98** |

51% of criteria (50 of 98) are fully automatable by the conformance suite.

### Domains

| Domain | Criteria | What it covers |
|---|---|---|
| Output Structure | OS-1 to OS-10 | Stream separation, exit codes, clean piped output, machine-readable formats |
| Color & Visual | CV-1 to CV-15 | NO_COLOR, TERM=dumb, TTY-aware color, 4-bit ANSI preference, contrast, OSC 8 hyperlinks, fg/bg pair contrast, Unicode symbol accessibility |
| Help & Documentation | HD-1 to HD-13 | --help/-h, --version, subcommand help, man pages, shell completions |
| Error Handling | EF-1 to EF-10 | Errors to stderr, exit codes, human-readable messages, actionable suggestions |
| Interactivity | II-1 to II-11 | Non-interactive mode, Ctrl-C, password masking, dry-run, keyboard navigation |
| Environment | EC-1 to EC-11 | Config precedence, XDG compliance, PAGER/EDITOR respect, no telemetry |
| Timing & Motion | TM-1 to TM-7 | No animation when piped, static progress, responsive startup, no flashing |
| Internationalization | IL-1 to IL-5 | UTF-8, externalizable strings, BiDi safety, locale-aware formatting |
| Lifecycle | IN-1 to IN-4 | Accessible installation, deprecation warnings, non-destructive updates |

## Research Foundation

This specification is informed by:

- [Sampath et al., "Accessibility of Command Line Interfaces," CHI 2021](https://dl.acm.org/doi/abs/10.1145/3411764.3445544) — seminal academic study
- [W3C WCAG2ICT](https://www.w3.org/TR/wcag2ict-22/) — guidance on applying WCAG to non-web ICT
- [W3C Background on Text/Terminal Applications](https://github.com/w3c/wcag2ict/blob/main/background-on-text-command-line-terminal-applications-and-interfaces.md)
- [Command Line Interface Guidelines](https://clig.dev/)
- [GitHub CLI Accessibility Work](https://github.blog/engineering/user-experience/building-a-more-accessible-github-cli/)
- [Seirdy's Inclusive CLI Best Practices](https://seirdy.one/posts/2022/06/10/cli-best-practices/)
- [NO_COLOR](https://no-color.org/) and [FORCE_COLOR](https://force-color.org/) standards
- Existing CLI VPATs from GitHub CLI, npm CLI, Google Cloud CLI, and ImageMagick

## License

See [LICENSE](LICENSE) for details.
