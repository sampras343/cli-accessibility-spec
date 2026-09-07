---
criterion_id: CV-6
domain: color
level: AA
testability: AUTO
spec_version: "1.0"
type: yaml
last_reviewed: 2026-09-07
---

# CV-6: --color Flag with Modes

## What This Tests

Whether the tool supports a `--color=WHEN` flag with at least `never` and `always` modes for per-invocation color control.

## How The Test Works

1. Runs with `--color=never --help` and verifies no ANSI color codes in stdout
2. Runs with `--color=always --help` in a pipe and verifies ANSI codes ARE present

## Pass Criteria

- `--color=never` suppresses all color escape sequences
- `--color=always` forces color even when stdout is piped (non-TTY)

## Fail Criteria

- `--color=never` or `--color=always` flags are not recognized
- Color codes present despite `--color=never`
- No color codes despite `--color=always` in pipe

## Non-Standard Flag Names

Some tools implement equivalent functionality with different flag names:

- GCC: `-fdiagnostics-color=auto/always/never`
- git: `-c color.ui=always/false` (config-based, not a flag)
- clang: `-fcolor-diagnostics` / `-fno-color-diagnostics`

Tools using non-standard names satisfy the **intent** of this criterion
(user-controllable color per invocation) even if they fail the specific
`--color=WHEN` test. Evaluators should note non-standard mechanisms in
the remarks and consider marking as "Partially Supports" rather than
"Does Not Support" when equivalent control exists.

## References

- CLI-ACS v1.0 spec, Section 4.2, CV-6
- GNU coreutils `--color=auto/always/never` convention
