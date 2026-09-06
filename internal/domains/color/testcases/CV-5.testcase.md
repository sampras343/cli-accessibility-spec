---
criterion_id: CV-5
domain: color
level: A
testability: AUTO
spec_version: "1.0"
type: yaml
last_reviewed: 2026-09-06
---

# CV-5: TERM=dumb Respect

## What This Tests

Whether the tool respects the `TERM=dumb` environment variable and suppresses ANSI color codes and other terminal formatting when running in a basic terminal environment.

## How The Test Works

1. Runs the binary with `TERM=dumb` set and `--help`
2. Scans stdout for ANSI SGR color sequences (3-bit, 8-bit, and 24-bit color codes)
3. Verifies exit code is 0 (command succeeds)

## Pass Criteria

- Exit code is 0
- No ANSI color escape sequences in stdout when `TERM=dumb` is set
- Tool detects dumb terminal and disables color/formatting

## Fail Criteria

- ANSI color escape sequences found in output despite `TERM=dumb`
- Tool ignores TERM environment variable
- Command fails (non-zero exit code)

## References

- CLI-ACS v1.0 spec, Section 4.2, CV-5
- TERM environment variable specification
- Common practice in CLI tools to respect TERM=dumb
