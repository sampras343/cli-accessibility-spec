# CLI Accessibility Conformance Specification (CLI-ACS)

A CLI-native accessibility standard that defines measurable criteria for evaluating how well command-line interface tools serve users with disabilities.

## Why CLI-ACS?

Existing accessibility frameworks (WCAG 2.2, EN 301 549, Section 508) apply to CLI tools in principle but provide little CLI-specific guidance. Current conformance reports (VPATs) for CLI tools are often meaningless — marking criteria like "Pointer Gestures" and "Dragging Movements" as "Supports" for text-based tools, or marking everything "Not Applicable."

CLI-ACS addresses this gap with 95 criteria native to the CLI context, organized by functional domain, mapped to disability categories, and tagged with testability classifications for conformance suite implementation.

## Specification

- [CLI-ACS v1.0 Design Document](docs/superpowers/specs/2026-08-25-cli-acs-design.md)

## Key Features

- **CLI-native criteria** organized by 9 functional domains (Output, Color, Help, Errors, Input, Environment, Timing, i18n, Lifecycle)
- **TUI Extension Module** for full-screen terminal applications
- **WCAG-style conformance levels** (A / AA / AAA)
- **Testability tags** ([AUTO] / [SEMI] / [MANUAL]) guiding conformance suite implementation
- **Disability Impact Matrix** showing coverage for Visual, Motor, Cognitive, and Vestibular categories
- **Standards Crosswalk** mapping to WCAG 2.2, EN 301 549, and Section 508
- **Conformance Report Template** (CLI-ACR) with both human-readable and machine-readable (JSON) formats

## Criteria Summary

| Level | Core CLI | TUI Extension | Total |
|---|---|---|---|
| A (Minimum) | 27 | 5 | 32 |
| AA (Standard) | 38 | 6 | 44 |
| AAA (Enhanced) | 16 | 3 | 19 |
| **Total** | **81** | **14** | **95** |

53% of criteria (50 of 95) are fully automatable by a conformance suite.

## Research Foundation

This specification is informed by:

- [Sampath et al., "Accessibility of Command Line Interfaces," CHI 2021](https://dl.acm.org/doi/abs/10.1145/3411764.3445544) — seminal academic study
- [W3C WCAG2ICT](https://www.w3.org/TR/wcag2ict-22/) — guidance on applying WCAG to non-web ICT
- [W3C Background on Text/Terminal Applications](https://github.com/w3c/wcag2ict/blob/main/background-on-text-command-line-terminal-applications-and-interfaces.md)
- [Command Line Interface Guidelines](https://clig.dev/)
- [GitHub CLI Accessibility Work](https://github.blog/engineering/user-experience/building-a-more-accessible-github-cli/)
- [Seirdy's Inclusive CLI Best Practices](https://seirdy.one/posts/2022/06/10/cli-best-practices/)
- [NO_COLOR](https://no-color.org/) and [FORCE_COLOR](https://force-color.org/) standards
- Existing CLI VPATs from GitHub CLI, npm CLI, Google Cloud CLI, and ImageMagick

## License

See [LICENSE](LICENSE) for details.
