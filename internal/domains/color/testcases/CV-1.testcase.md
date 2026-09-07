---
criterion_id: CV-1
domain: color
level: A
testability: SEMI
spec_version: "1.0"
type: go
last_reviewed: 2026-09-07
---

# CV-1: Color Not Sole Information Channel

## What This Tests

Whether the tool uses additional indicators (text prefixes, symbols, indentation) alongside color to convey information, so that information is not lost when color is removed.

## How The Test Works

1. Runs the binary with color enabled (default) and captures output
2. Runs the binary with `NO_COLOR=1` and captures output
3. Strips ANSI escape sequences from the colored output
4. Compares the stripped colored output against the no-color output
5. If >95% similar, information is preserved without color — passes
6. Flags for human review since automated comparison is heuristic

## Pass Criteria

- Stripped colored output is >95% similar to no-color output
- Information structure (sections, items, status) is preserved without color

## Fail Criteria

- Significant content differs between colored and no-color output
- Status indicators, categories, or groupings disappear without color
