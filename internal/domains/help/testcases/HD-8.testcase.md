---
criterion_id: HD-8
domain: help
level: AA
testability: AUTO
spec_version: "1.0"
type: go
last_reviewed: 2026-09-07
---

# HD-8: No Hang on Empty Stdin

## What This Tests

Whether the tool exits within a reasonable time when invoked with no arguments and stdin connected to /dev/null (empty input).

## How The Test Works

1. Runs the binary with no arguments and a 5-second timeout
2. In Go's os/exec, stdin defaults to os.DevNull, so the process receives immediate EOF
3. Checks whether the process exits before the timeout

## Pass Criteria

- Process exits within 5 seconds (any exit code is acceptable)

## Fail Criteria

- Process times out (hangs waiting for interactive input despite empty stdin)

## Notes

Tools that are designed to read from stdin (like `cat`, `sed`, `awk`) will pass this test because they receive EOF from /dev/null and exit immediately. The criterion catches tools that block waiting for interactive input without checking whether stdin is a terminal.

## References

- CLI-ACS v1.0 spec, Section 4.3, HD-8
