---
criterion_id: HD-4
domain: help
level: A
testability: AUTO
spec_version: "1.0"
type: yaml
last_reviewed: 2026-09-07
---

# HD-4: Missing-Argument Guidance

## What This Tests

Whether the tool provides helpful guidance when invoked with an unrecognized flag, specifically whether error output mentions how to get help.

## How The Test Works

1. Runs the binary with `--nonexistent-flag-a11y-probe` to trigger an error
2. Asserts non-zero exit code (the tool recognized the flag is invalid)
3. Asserts stderr contains "help", "--help", or "usage" (pointer to documentation)

## Pass Criteria

- Non-zero exit code (tool recognizes the error)
- Stderr contains a pointer to help (mentions "help", "--help", or "usage")

## Fail Criteria

- Exit code 0 (tool silently ignores unknown flags)
- Stderr does not mention help or usage

## References

- CLI-ACS v1.0 spec, Section 4.3, HD-4
