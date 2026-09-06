---
criterion_id: CV-4
domain: color
level: A
testability: AUTO
spec_version: "1.0"
type: yaml
last_reviewed: 2026-09-06
---

# CV-4: TTY-Aware Color

## What This Tests

Whether the tool automatically suppresses ANSI color codes when stdout is not connected to a terminal (e.g., when piped to another command or redirected to a file).

## How The Test Works

1. Runs the binary with `--help` in a non-TTY context (piped execution)
2. Scans stdout for ANSI SGR color sequences (3-bit, 8-bit, and 24-bit color codes)
3. Verifies exit code is 0 (command succeeds)

## Pass Criteria

- Exit code is 0
- No ANSI color escape sequences in stdout when output is piped (not a TTY)
- Tool detects non-TTY environment and disables color automatically

## Fail Criteria

- ANSI color escape sequences found in piped output
- Tool emits color codes regardless of TTY status
- Command fails (non-zero exit code)

## References

- CLI-ACS v1.0 spec, Section 4.2, CV-4
- Common implementation: `isatty(STDOUT_FILENO)` or equivalent
