---
criterion_id: CV-8
domain: color
level: AA
testability: SEMI
spec_version: "1.0"
type: go
last_reviewed: 2026-09-07
---

# CV-8: 4-Bit ANSI Color Preference

## What This Tests

Whether the tool prefers 4-bit ANSI colors (the 16 user-customizable colors) over fixed 8-bit (256-color) or 24-bit (truecolor) for informational content.

## How The Test Works

1. Runs the binary and captures output with ANSI codes
2. Classifies all color codes as 4-bit (SGR 30-37, 40-47, 90-97, 100-107), 8-bit (38;5;N where N>=16), or 24-bit (38;2;R;G;B)
3. If >50% of color codes are 8-bit or 24-bit, flags for review
4. 4-bit colors are user-customizable in terminal settings; fixed colors bypass user preferences

## Pass Criteria

- Majority of color codes are 4-bit ANSI (user-customizable)
- Any 8-bit/24-bit colors are used only for decorative elements

## Fail Criteria

- >50% of color codes are fixed 8-bit or 24-bit
- Error/warning/status indicators use fixed colors instead of ANSI 16
