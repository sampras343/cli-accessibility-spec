---
criterion_id: HD-3
domain: help
level: A
testability: AUTO
spec_version: "1.0"
type: yaml
last_reviewed: 2026-09-07
---

# HD-3: --version Flag

## What This Tests

Whether the tool provides a `--version` flag that outputs version information and exits with code 0.

## How The Test Works

1. Runs the binary with `--version`
2. Checks for exit code 0
3. Checks for non-empty stdout (version string)

## Pass Criteria

- `--version` exits with code 0
- Stdout contains version information (non-empty)

## Fail Criteria

- Non-zero exit code
- Empty stdout (no version information)

## References

- CLI-ACS v1.0 spec, Section 4.3, HD-3
- GNU Coding Standards, 4.7.1 (--version)
