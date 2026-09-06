---
criterion_id: CV-2
domain: color
level: A
testability: AUTO
spec_version: "1.0"
type: yaml
last_reviewed: 2026-09-06
---

# CV-2: NO_COLOR Support

## What This Tests

Whether the tool suppresses ANSI color codes when the NO_COLOR environment variable is set to a non-empty value, per the no-color.org specification.

## How The Test Works

1. Runs the binary with `NO_COLOR=1` and `--help`
2. Scans stdout for ANSI SGR color sequences (3-bit, 8-bit, and 24-bit color codes)
3. Verifies exit code is 0 (command succeeds)

## Pass Criteria

- Exit code is 0
- No ANSI color escape sequences in stdout when `NO_COLOR=1` is set
- Specifically checks for patterns: `\x1b[(30-37|90-97|38;5;N|38;2;R;G;B)m`

## Fail Criteria

- ANSI color escape sequences found in output despite `NO_COLOR=1`
- Command fails (non-zero exit code)

## References

- [no-color.org](https://no-color.org/)
- CLI-ACS v1.0 spec, Section 4.2, CV-2
