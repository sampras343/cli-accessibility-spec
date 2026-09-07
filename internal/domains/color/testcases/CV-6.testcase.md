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
