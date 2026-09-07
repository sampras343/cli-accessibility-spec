---
criterion_id: HD-9
domain: help
level: AA
testability: AUTO
spec_version: "1.0"
type: go
last_reviewed: 2026-09-07
---

# HD-9: Man Page Availability

## What This Tests

Whether a man page is installed for the binary.

## How The Test Works

1. Extracts the binary name from its path
2. Runs `man <binary-name>` with a 5-second timeout
3. Checks whether man exits 0 with non-empty output

## Pass Criteria

- `man <binary-name>` exits 0 with non-empty output

## Fail Criteria

- `man` exits non-zero or produces no output (no man page installed)

## Notes

This checks the system's man page database. Tools installed via package managers typically include man pages; tools installed via `go install`, `cargo install`, or `npm install -g` often do not.

## References

- CLI-ACS v1.0 spec, Section 4.3, HD-9
- man-pages(7)
