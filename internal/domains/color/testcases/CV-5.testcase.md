---
criterion_id: CV-5
domain: color
level: A
testability: AUTO
spec_version: "1.0"
type: yaml
last_reviewed: 2026-09-07
---

# CV-5: TERM=dumb Respect

## What This Tests

Whether the tool suppresses ALL ANSI escape sequences — not just color
codes, but also cursor movement, screen clearing, text styling, erase
in line, and OSC sequences — when `TERM=dumb` is set.

`TERM=dumb` signals a terminal with zero capability for escape sequence
interpretation (e.g., Emacs shell buffers, CI log viewers, serial
consoles). Any escape byte (`\x1b`) in output indicates a violation.

## How The Test Works

1. Runs the binary with `TERM=dumb` and `--help`
2. Scans stdout for ANY occurrence of the ESC byte (`\x1b`, 0x1B)
3. This catches all escape types: CSI color codes (`\x1b[31m`), cursor
   movement (`\x1b[H`), erase in line (`\x1b[K`), OSC sequences
   (`\x1b]8;`), and any other ANSI escapes

## Pass Criteria

- No ESC byte (`\x1b`) present anywhere in stdout when `TERM=dumb`
- Exit code is 0

## Fail Criteria

- Any ESC byte found in stdout — indicates the tool emits escape
  sequences in a context that cannot interpret them
- Common violations: GCC emitting `\x1b[K` (Erase in Line), tools
  emitting bold (`\x1b[1m`) despite `TERM=dumb`

## References

- CLI-ACS v1.0 spec, Section 4.2, CV-5
- POSIX terminal handling: `TERM=dumb` indicates no terminfo capabilities
