---
criterion_id: HD-5
domain: help
level: AA
testability: SEMI
spec_version: "1.0"
type: go
last_reviewed: 2026-09-07
---

# HD-5: Help Text Structure

## What This Tests

Whether the help text follows a structured format with distinct sections: description, usage/synopsis, flags/options, and blank-line separation between sections.

## How The Test Works

1. Parses the probe's captured `--help` output
2. Scans for recognizable section indicators: usage/synopsis lines, flag/option headers, command listings, example sections
3. Checks for a description in the first few lines
4. Checks for blank-line separation between sections
5. Counts total recognizable sections

## Pass Criteria

- At least 3 sections identified (e.g., description, usage, flags)
- Blank-line separation between sections

## Fail Criteria

- Fewer than 2 recognizable sections
- Unstructured wall of text without clear organization

## References

- CLI-ACS v1.0 spec, Section 4.3, HD-5
