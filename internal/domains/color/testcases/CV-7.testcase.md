---
criterion_id: CV-7
domain: color
level: AA
testability: AUTO
spec_version: "1.0"
type: yaml
last_reviewed: 2026-09-07
---

# CV-7: FORCE_COLOR Support

## What This Tests

Whether the tool forces color output when the `FORCE_COLOR` environment variable is set, and whether `FORCE_COLOR` overrides `NO_COLOR`.

## How The Test Works

1. Runs with `FORCE_COLOR=1` in piped mode and verifies ANSI codes are present
2. Runs with both `NO_COLOR=1` and `FORCE_COLOR=1` and verifies color wins (FORCE_COLOR overrides NO_COLOR)

## Pass Criteria

- ANSI color codes present in piped output when `FORCE_COLOR=1`
- Color present when both `NO_COLOR=1` and `FORCE_COLOR=1` are set

## Fail Criteria

- No ANSI codes in output despite `FORCE_COLOR=1`
- `NO_COLOR` wins over `FORCE_COLOR` (wrong precedence)
