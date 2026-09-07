---
criterion_id: CV-9
domain: color
level: AA
testability: SEMI
spec_version: "1.0"
type: go
last_reviewed: 2026-09-07
---

# CV-9: No Background Color Assumption

## What This Tests

Whether the tool avoids assuming a specific terminal background color, ensuring readability on both light and dark backgrounds.

## How The Test Works

1. Runs the binary and captures all ANSI color codes
2. Extracts fixed foreground colors (8-bit indices 16+ and 24-bit RGB)
3. Calculates contrast ratio against white (#FFFFFF) and dark (#1E1E1E) backgrounds
4. Flags any fixed color that fails 4.5:1 contrast against either background
5. ANSI 16 colors (indices 0-15) are exempt since they are user-customizable

## Pass Criteria

- No fixed foreground colors fail 4.5:1 contrast against either light or dark backgrounds
- Or tool uses only ANSI 16 named colors (user-customizable, always exempt)

## Fail Criteria

- Fixed 8-bit or 24-bit foreground colors that fail contrast on light or dark backgrounds
