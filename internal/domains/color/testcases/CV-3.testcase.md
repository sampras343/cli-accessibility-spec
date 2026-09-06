---
criterion_id: CV-3
domain: color
level: A
testability: AUTO
spec_version: "1.0"
type: yaml
last_reviewed: 2026-09-06
---

# CV-3: --no-color Flag Support

## What This Tests

Whether the tool provides a `--no-color` flag (or equivalent) that suppresses ANSI color codes in output.

## How The Test Works

1. Runs the binary with `--no-color --help`
2. Scans stdout for ANSI SGR color sequences (3-bit, 8-bit, and 24-bit color codes)
3. Verifies exit code is 0 (command succeeds)

## Pass Criteria

- Exit code is 0
- No ANSI color escape sequences in stdout when `--no-color` flag is used
- The flag is accepted without error

## Fail Criteria

- ANSI color escape sequences found in output despite `--no-color` flag
- Flag is not recognized (command fails)
- Command fails (non-zero exit code)

## References

- CLI-ACS v1.0 spec, Section 4.2, CV-3
- Common practice: `--no-color`, `--color=never`, `--plain`
