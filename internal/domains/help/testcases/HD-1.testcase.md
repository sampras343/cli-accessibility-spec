---
criterion_id: HD-1
domain: help
level: A
testability: AUTO
spec_version: "1.0"
type: yaml
last_reviewed: 2026-09-07
---

# HD-1: --help and -h Flags

## What This Tests

Whether the tool provides both `--help` and `-h` flags that produce help output and exit with code 0.

## How The Test Works

1. Runs the binary with `--help` and checks for exit code 0 and non-empty stdout
2. Runs the binary with `-h` and checks for exit code 0 and non-empty stdout

## Pass Criteria

- Both `--help` and `-h` exit with code 0
- Both produce non-empty output on stdout

## Fail Criteria

- Either flag exits with non-zero code (e.g., `go --help` exits 2, `npm --help` exits 1)
- Either flag produces no output on stdout
- The `-h` flag is used for a different purpose (see note below)

## Notes

Some POSIX tools use `-h` for other purposes (e.g., `ls -h` for human-readable sizes, `grep -h` to suppress filenames). These tools will correctly fail this criterion because they genuinely do not provide `-h` as a help shortcut. This is expected behavior -- the criterion tests for the `-h` convention, not POSIX compatibility.

Some tools exit with non-zero codes for `--help` (e.g., `go` exits 2, `npm` exits 1). The criterion correctly catches this as non-conforming -- help requests should succeed with exit 0.

## References

- CLI-ACS v1.0 spec, Section 4.3, HD-1
- GNU Coding Standards, 4.7.2 (--help)
