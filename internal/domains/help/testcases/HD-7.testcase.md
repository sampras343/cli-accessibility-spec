---
criterion_id: HD-7
domain: help
level: AA
testability: SEMI
spec_version: "1.0"
type: go
last_reviewed: 2026-09-07
---

# HD-7: Consistent Flag Format

## What This Tests

Whether flags in the help text use a consistent presentation format: consistent indentation, and a uniform style (all short+long, or all long-only).

## How The Test Works

1. Extracts flag-like lines from `--help` output (lines with leading whitespace and dashes)
2. Checks whether flags consistently use short+long format (e.g., `-v, --verbose`) or long-only format (e.g., `--verbose`)
3. Measures indentation consistency across flag lines

## Pass Criteria

- Consistent format (all short+long or all long-only)
- Consistent indentation across flag lines

## Fail Criteria

- Mixed flag formats (some short+long, some long-only)
- Inconsistent indentation

## References

- CLI-ACS v1.0 spec, Section 4.3, HD-7
