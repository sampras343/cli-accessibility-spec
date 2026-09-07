---
criterion_id: HD-6
domain: help
level: AA
testability: SEMI
spec_version: "1.0"
type: go
last_reviewed: 2026-09-07
---

# HD-6: Examples in Help

## What This Tests

Whether the help text includes usage examples that demonstrate how to invoke the tool.

## How The Test Works

1. Scans `--help` output for example-related patterns
2. Looks for an explicit "Examples" section header
3. Looks for lines starting with shell prompts (`$`, `>`, `#`)
4. Looks for lines containing the binary name with arguments

## Pass Criteria

- Found an "Examples" section header AND example-like lines
- Partially supports if example-like lines exist without a header

## Fail Criteria

- No examples found in help output

## References

- CLI-ACS v1.0 spec, Section 4.3, HD-6
