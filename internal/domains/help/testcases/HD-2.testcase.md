---
criterion_id: HD-2
domain: help
level: A
testability: AUTO
spec_version: "1.0"
type: go
last_reviewed: 2026-09-07
---

# HD-2: Subcommand Help

## What This Tests

Whether every discovered subcommand supports `--help` with exit code 0 and non-empty output.

## How The Test Works

1. Uses the probe's discovered subcommand list (filtered and limited)
2. For each subcommand, runs `<binary> <subcommand> --help`
3. Verifies each exits 0 with non-empty stdout or stderr
4. Reports per-subcommand results

## Pass Criteria

- All tested subcommands exit 0 on `--help`
- All tested subcommands produce non-empty help output
- Partially supports if some subcommands pass but not all

## Fail Criteria

- All subcommands fail `--help` (non-zero exit or empty output)

## Preconditions

- Tool must have subcommands (has_subcommands from probe)

## References

- CLI-ACS v1.0 spec, Section 4.3, HD-2
