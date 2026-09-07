# CLI Accessibility Conformance Specification (CLI-ACS) v1.0

## Design Document

**Date:** 2026-08-25
**Status:** Draft
**Authors:** Sachin Sampras M

---

## Table of Contents

1. [Introduction & Scope](#1-introduction--scope)
2. [Terminology](#2-terminology)
3. [Conformance Model](#3-conformance-model)
4. [Functional Domain Criteria](#4-functional-domain-criteria)
   - 4.1 [Output Structure & Content](#41-output-structure--content)
   - 4.2 [Color & Visual Presentation](#42-color--visual-presentation)
   - 4.3 [Help & Documentation](#43-help--documentation)
   - 4.4 [Error Handling & Feedback](#44-error-handling--feedback)
   - 4.5 [Interactivity & Input](#45-interactivity--input)
   - 4.6 [Environment Awareness & Configuration](#46-environment-awareness--configuration)
   - 4.7 [Timing & Motion](#47-timing--motion)
   - 4.8 [Internationalization & Localization](#48-internationalization--localization)
   - 4.9 [Installation & Lifecycle](#49-installation--lifecycle)
5. [TUI Extension Module](#5-tui-extension-module)
6. [Disability Impact Matrix](#6-disability-impact-matrix)
7. [Conformance Report Template](#7-conformance-report-template)
8. [Normative References & Standards Crosswalk](#8-normative-references--standards-crosswalk)

---

## 1. Introduction & Scope

### 1.1 Purpose

The CLI Accessibility Conformance Specification (CLI-ACS) defines measurable accessibility criteria for command-line interface (CLI) tools. It provides a framework for evaluating how well CLI tools serve users with disabilities, and a standard conformance report format for communicating the results.

Unlike web accessibility, which is well-served by WCAG 2.2 and its extensive techniques, CLI tools operate in a fundamentally different environment — terminal emulators — where there is no DOM, no accessibility tree, no ARIA roles, and no standardized semantic markup. Existing accessibility frameworks (WCAG, EN 301 549, Section 508) apply to CLI tools in principle but provide little CLI-specific guidance, resulting in conformance reports (VPATs/ACRs) that are either meaningless (marking all 55 WCAG criteria as "Supports" with no remarks) or inapplicable (evaluating criteria like "Pointer Gestures" and "Dragging Movements" for a text-based tool).

CLI-ACS addresses this gap by defining criteria native to the CLI context — organized by the functional domains CLI developers work in, mapped to the disability categories they serve, and tagged with testability classifications that directly guide conformance suite implementation.

### 1.2 Scope

**In scope (Core CLI):**

- Non-interactive command-line tools that accept input via arguments, flags, and stdin, and produce output to stdout and stderr (e.g., `grep`, `curl`, `jq`, `git`)
- Lightly interactive CLI tools that use prompts, simple menus, progress indicators, or confirmation dialogs (e.g., `npm init`, `gh pr create`, `aws configure`)
- CLI tools distributed as standalone binaries, scripts, or interpreted programs

**In scope (TUI Extension Module):**

- Full-screen Terminal User Interface applications that render content into a character cell matrix and manage their own input handling (e.g., `vim`, `htop`, `tmux`, `ranger`)

**Out of scope:**

- GUI applications launched from the terminal
- Web-based terminal emulators (evaluate via WCAG 2.2)
- Terminal emulators themselves (these are the "user agent" in the CLI accessibility model)
- Shell environments (bash, zsh, fish) — these are platforms, not tools under test

### 1.3 Relationship to Existing Standards

CLI-ACS is a **CLI-native accessibility standard** that references and cross-maps to existing standards without duplicating them:

- **WCAG 2.2** — CLI-ACS criteria map to applicable WCAG Success Criteria where relevant. The Standards Crosswalk (Section 8.3) provides a complete mapping.
- **WCAG2ICT** — CLI-ACS builds on the W3C's guidance for applying WCAG to non-web ICT, particularly its treatment of terminal emulators as "user agents" and its interpretation of "programmatically determined" for text interfaces.
- **EN 301 549 v3.2.1** — CLI-ACS criteria map to Chapter 11 (Non-web Software) requirements. The crosswalk enables EN 301 549 conformance claims.
- **Section 508 (Revised)** — CLI-ACS criteria map to Chapter 3 (Functional Performance Criteria) and Chapter 5 (Software) requirements.
- **VPAT 2.5** — The CLI-ACS Conformance Report Template (Section 7) is designed to complement, not replace, VPATs. The crosswalk enables generation of VPAT-format reports from CLI-ACS evaluations.

### 1.4 The Terminal Accessibility Model

In the terminal accessibility model, the roles map as follows:

| Web Model | Terminal Model |
|---|---|
| Web page / web application | CLI tool / TUI application |
| Web browser | Terminal emulator |
| HTML DOM / accessibility tree | Character cell matrix / text buffer |
| ARIA roles and properties | (No equivalent — text is unstructured) |
| CSS styling | ANSI escape sequences |
| Screen reader (JAWS, NVDA, VoiceOver) | Terminal screen reader (NVDA terminal mode, VoiceOver Terminal, Fenrir, tdsr, Emacspeak) |

The terminal emulator renders characters on screen and exposes the character cell matrix to assistive technology. Screen readers for terminal applications obtain text from this matrix and apply heuristics to detect structure (headings via capitalization/spacing, input fields via inverse video, columns via alignment). Unlike web screen readers, terminal screen readers have no semantic markup to rely on — all structure must be inferred from visual layout or explicitly communicated through text formatting conventions.

This means CLI tools bear a greater responsibility for producing well-structured, predictable, linear text output than web applications do, because there is no accessibility tree to compensate for poor visual structure.

---

## 2. Terminology

**Assistive Technology (AT):** Software or hardware that provides functionality to meet the requirements of users with disabilities. In the CLI context, this includes screen readers (NVDA, JAWS, VoiceOver, Orca, Fenrir, tdsr, Emacspeak), screen magnifiers, braille displays, switch access devices, and voice input systems.

**ANSI Escape Sequence:** A sequence of characters beginning with the ASCII Escape character (0x1B) followed by control codes, used to control cursor position, text color, text styling, and other terminal display features. Defined by ECMA-48 and commonly called "ANSI codes."

**CLI (Command-Line Interface):** A text-based interface where users interact with software by typing commands, flags, and arguments, and the software responds with text output. Distinguished from TUIs by not managing the full terminal screen.

**Conformance Suite:** An automated tool that evaluates a CLI binary against CLI-ACS criteria, producing a Conformance Report.

**Exit Code:** A numeric value (0–255) returned by a process to its parent upon termination. By POSIX convention, 0 indicates success and non-zero indicates failure.

**Pager:** A program that displays text one screenful at a time (e.g., `less`, `more`). Users configure their preferred pager via the `PAGER` environment variable.

**Programmatically Determined:** Information that can be obtained by assistive technology from the software's output. In the terminal context, this means text that appears in the terminal's character cell matrix and can be read by a screen reader, or structured data available via machine-readable output formats.

**stderr (Standard Error):** File descriptor 2. The conventional output stream for diagnostic messages, errors, warnings, and progress information.

**stdout (Standard Output):** File descriptor 1. The conventional output stream for a program's primary data output.

**Terminal Emulator:** Software that emulates a video terminal within a graphical environment. Acts as the "user agent" for CLI and TUI applications, rendering characters and exposing them to assistive technology. Examples: GNOME Terminal, iTerm2, Windows Terminal, xterm.

**TTY (Teletypewriter):** In modern usage, refers to a terminal device or pseudoterminal. The `isatty()` function tests whether a file descriptor is connected to a terminal, as opposed to a pipe or file redirect.

**TUI (Terminal User Interface):** A full-screen text-based application that manages the entire terminal display, placing characters at specific positions in the character cell matrix. Distinguished from CLIs by direct screen management. Examples: vim, htop, tmux.

---

## 3. Conformance Model

### 3.1 Conformance Levels

CLI-ACS defines three conformance levels, modeled after WCAG 2.2:

**Level A — Minimum Accessibility**

Criteria whose failure blocks access to core functionality for users with disabilities. Every CLI tool MUST meet all applicable Level A criteria to claim any CLI-ACS conformance.

Level A criteria address the most fundamental barriers: output that can be read by screen readers, exit codes that indicate success/failure, color that doesn't carry sole meaning, help text that exists and is findable, and errors that go to stderr.

**Level AA — Standard Conformance Target**

Criteria that significantly improve the experience for users with disabilities beyond the minimum. Level AA is the recommended conformance target for production software.

Level AA criteria address quality and usability: structured output formats, actionable error messages, environment variable respect, configurable behavior, and consistent documentation.

**Level AAA — Enhanced Accessibility**

Aspirational criteria that represent best-in-class CLI accessibility. Level AAA is not required for conformance claims but demonstrates leadership in accessibility.

Level AAA criteria address polish and completeness: shell completions, pager support, plain-text modes, bug report facilitation, and advanced screen reader support.

### 3.2 Testability Classifications

Each criterion is tagged with a testability classification that guides conformance suite implementation:

**`[AUTO]` — Automated**

The conformance suite can fully evaluate this criterion by running the binary with specific inputs, environment variables, and flags, then examining outputs, exit codes, and stream behavior. No human judgment is required.

Example: Testing `NO_COLOR` support — run the tool with `NO_COLOR=1` set and verify output contains no ANSI escape sequences.

**`[SEMI]` — Semi-Automated**

The conformance suite can detect relevant signals (presence of ANSI codes, output length, error message format) but a human evaluator must judge quality or completeness.

Example: Testing "human-readable error messages" — the suite can verify errors go to stderr and contain text (not just a code), but a human must judge whether the text is actually understandable.

**`[MANUAL]` — Manual**

The criterion requires testing with assistive technology (screen reader, magnifier) or human judgment about interaction quality that cannot be automated.

Example: Testing "keyboard-navigable selections" — requires a screen reader user to interact with the tool's interactive menus and verify they are announced correctly.

### 3.3 Conformance Claims

A tool may make the following conformance claims:

- **"CLI-ACS Level A Conformant"** — Meets all applicable Level A criteria across all applicable domains.
- **"CLI-ACS Level AA Conformant"** — Meets all applicable Level A and Level AA criteria.
- **"CLI-ACS Level AAA Conformant"** — Meets all applicable Level A, Level AA, and Level AAA criteria.
- **Partial conformance** — Tools MAY report per-domain conformance (e.g., "Level AA in Output Structure, Level A in Help & Documentation"). Partial conformance MUST NOT be described as "CLI-ACS Level AA Conformant" without qualification.
- **"CLI-ACS Level [A/AA/AAA] + TUI Extension Conformant"** — For TUI applications, conformance includes both core criteria and TUI Extension Module criteria.

### 3.4 Applicability

Not all criteria apply to all tools. A criterion is **Not Applicable** when:

- The tool does not have the feature the criterion addresses (e.g., "Subcommand help" for a tool with no subcommands)
- The tool's execution model makes the criterion irrelevant (e.g., "Configurable timeouts" for a tool that performs no network or time-bounded operations)

A criterion MUST NOT be marked Not Applicable merely because implementing it would be difficult. The evaluator MUST document the rationale for each Not Applicable determination in the conformance report.

### 3.5 Evaluation Methodology

A complete CLI-ACS evaluation involves:

1. **Automated scan** — Run the conformance suite against the binary, which tests all `[AUTO]` criteria.
2. **Semi-automated review** — Review the suite's flagged `[SEMI]` criteria, applying human judgment to each.
3. **Manual testing** — Test all `[MANUAL]` criteria using assistive technology on at least one platform/screen reader combination.
4. **Report generation** — Compile results into the CLI-ACS Conformance Report Template (Section 7).

The testing environment (OS, terminal emulator, shell, screen reader, suite version) MUST be documented in the report.

---

## 4. Functional Domain Criteria

### 4.1 Output Structure & Content

This domain covers how CLI tools structure and present their output to stdout/stderr.

#### OS-1: Standard Stream Separation

| | |
|---|---|
| **Level** | A |
| **Testability** | `[AUTO]` |
| **Requirement** | Primary output MUST go to stdout. Diagnostic messages, progress indicators, warnings, and errors MUST go to stderr. Mixing primary output with diagnostics on stdout prevents piping and confuses screen readers that read stdout. |
| **Test method** | Run the tool with stdout and stderr redirected to separate files. Verify primary output appears only in stdout and diagnostics only in stderr. |
| **Disability impact** | Visual |
| **WCAG mapping** | — |

#### OS-2: Meaningful Exit Codes

| | |
|---|---|
| **Level** | A |
| **Testability** | `[AUTO]` |
| **Requirement** | The tool MUST return exit code 0 on success and a non-zero exit code on failure. The tool MUST NOT return 0 when the requested operation has failed. |
| **Test method** | Run the tool with valid inputs and verify exit code 0. Run with invalid inputs (missing required args, nonexistent files, invalid flags) and verify non-zero exit code. |
| **Disability impact** | Motor, Cognitive |
| **WCAG mapping** | — |

#### OS-3: Clean Piped Output

| | |
|---|---|
| **Level** | A |
| **Testability** | `[AUTO]` |
| **Requirement** | When stdout is not a TTY, output MUST be free of ANSI escape sequences, progress animations, spinner characters, and decorative elements. The tool MUST detect non-TTY stdout (via `isatty()` or equivalent) and strip formatting automatically. |
| **Test method** | Pipe the tool's output through `cat` (making stdout a pipe, not a TTY). Scan output for ANSI escape sequences (regex: `\x1b\[[\d;]*[a-zA-Z]`). Verify none are present. |
| **Disability impact** | Visual |
| **WCAG mapping** | — |

#### OS-4: Linear Reading Order

| | |
|---|---|
| **Level** | A |
| **Testability** | `[SEMI]` |
| **Requirement** | Output MUST be meaningful when read top-to-bottom, left-to-right. No reliance on spatial positioning (multi-column layouts, side-by-side comparisons) as the sole way to convey information. |
| **Test method** | Capture output and verify it does not use whitespace-based multi-column layout as the only representation. If columns are used, verify a `--plain` or `--json` alternative exists. |
| **Disability impact** | Visual, Cognitive |
| **WCAG mapping** | 1.3.2 Meaningful Sequence |

#### OS-5: No ASCII Art as Sole Information

| | |
|---|---|
| **Level** | A |
| **Testability** | `[SEMI]` |
| **Requirement** | ASCII art, box-drawing diagrams, sparklines, bar charts, and decorative borders MUST NOT be the only means of conveying information. A text alternative or `--plain` mode MUST exist if the tool produces such output. |
| **Test method** | Run the tool and scan for box-drawing characters (U+2500–U+257F), repeated decorative characters, and common ASCII art patterns. If found, verify a `--plain` or `--json` alternative produces equivalent information without these elements. |
| **Disability impact** | Visual, Cognitive |
| **WCAG mapping** | 1.1.1 Non-text Content |

#### OS-6: Machine-Readable Output Mode

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[AUTO]` |
| **Requirement** | The tool MUST provide at least one structured output format (e.g., `--json`, `--csv`, `--tsv`) for all commands that produce tabular or structured data. The structured output MUST contain the same information as the human-readable output. |
| **Test method** | Test for `--json` flag acceptance (exit code 0, valid JSON output). Test for `--csv` or `--tsv` if applicable. Compare information content between structured and human-readable output. |
| **Disability impact** | Visual, Motor |
| **WCAG mapping** | — |

#### OS-7: Quiet Mode

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[AUTO]` |
| **Requirement** | The tool MUST support a `-q` or `--quiet` flag that suppresses all non-essential output (decorations, tips, banners, progress), emitting only the primary result or nothing on success. |
| **Test method** | Run the tool with `-q` and `--quiet`. Verify accepted (exit code 0). Verify output is reduced compared to default invocation. |
| **Disability impact** | Cognitive |
| **WCAG mapping** | — |

#### OS-8: Output Width Awareness

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[AUTO]` |
| **Requirement** | The tool MUST respect terminal width when formatting output. It SHOULD check `COLUMNS` environment variable and/or terminal ioctl. When stdout is not a TTY, the tool MUST NOT assume a fixed width — it SHOULD either omit width-dependent formatting or default to a reasonable width (e.g., 80 columns). |
| **Test method** | Set `COLUMNS=40` and run the tool. Verify output lines do not exceed 40 characters (where the tool formats output). Pipe through `cat` and verify output does not assume a specific width. |
| **Disability impact** | Visual |
| **WCAG mapping** | — |

#### OS-9: Pager Support for Long Output

| | |
|---|---|
| **Level** | AAA |
| **Testability** | `[AUTO]` |
| **Requirement** | For output exceeding terminal height, the tool SHOULD support piping through a pager when stdout is a TTY. The tool MUST respect the `PAGER` environment variable. The tool MUST NOT invoke a pager when stdout is not a TTY. |
| **Test method** | Set `PAGER=cat` and run a command producing long output. Verify the tool invokes the pager when stdout is a TTY. Verify no pager is invoked when output is piped. |
| **Disability impact** | Visual, Motor |
| **WCAG mapping** | — |

#### OS-10: Plain Text Alternative for All Visual Elements

| | |
|---|---|
| **Level** | AAA |
| **Testability** | `[AUTO]` |
| **Requirement** | The tool SHOULD provide a `--plain` or `--accessible` flag that strips all decorative elements (borders, box-drawing, emoji, Unicode symbols, ASCII art) and produces screen-reader-friendly linear text output. |
| **Test method** | Run the tool with `--plain` or `--accessible`. Verify accepted. Scan output for decorative elements and verify they are absent. Compare information content with default output. |
| **Disability impact** | Visual, Cognitive |
| **WCAG mapping** | 1.1.1 Non-text Content |

---

### 4.2 Color & Visual Presentation

This domain covers how CLI tools use color, contrast, and visual styling in terminal output.

#### CV-1: Color Not Sole Information Channel

| | |
|---|---|
| **Level** | A |
| **Testability** | `[SEMI]` |
| **Requirement** | Color MUST NOT be the only visual means of conveying information, indicating an action, prompting a response, or distinguishing a visual element. The tool MUST use additional indicators alongside color: whitespace, indentation, text prefixes (e.g., `[ERROR]`, `[WARN]`, `[OK]`), symbols, or formatting changes. |
| **Test method** | Run the tool with `NO_COLOR=1`. Compare output with default colored output. Verify that all information distinguishable by color in the default output is also distinguishable in the no-color output via other means. |
| **Disability impact** | Visual, Cognitive |
| **WCAG mapping** | 1.4.1 Use of Color |

#### CV-2: `NO_COLOR` Support

| | |
|---|---|
| **Level** | A |
| **Testability** | `[AUTO]` |
| **Requirement** | When the `NO_COLOR` environment variable is set and non-empty (regardless of value), the tool MUST suppress all ANSI color codes from output. Other styling (bold, underline, inverse) MAY be retained. Per the no-color.org specification. |
| **Test method** | Run the tool with `NO_COLOR=1` and capture output. Scan for ANSI color escape sequences (SGR codes for foreground/background color: `\x1b\[[\d;]*[34][0-9]m`). Verify none are present. Run without `NO_COLOR` on a TTY and verify color codes ARE present (confirming the tool uses color by default). |
| **Disability impact** | Visual |
| **WCAG mapping** | 1.4.1 Use of Color |

#### CV-3: `--no-color` Flag

| | |
|---|---|
| **Level** | A |
| **Testability** | `[AUTO]` |
| **Requirement** | The tool MUST provide a `--no-color` command-line flag that suppresses ANSI color codes, giving users per-invocation control independent of environment variables. |
| **Test method** | Run the tool with `--no-color`. Verify accepted (exit 0). Scan output for ANSI color escape sequences. Verify none are present. |
| **Disability impact** | Visual |
| **WCAG mapping** | 1.4.1 Use of Color |

#### CV-4: TTY-Aware Color

| | |
|---|---|
| **Level** | A |
| **Testability** | `[AUTO]` |
| **Requirement** | When stdout is not a TTY, ANSI color and styling escape sequences MUST be suppressed by default. When stderr is not a TTY, the same MUST apply independently to stderr. This ensures piped output and file redirections are clean. |
| **Test method** | Pipe stdout through `cat` and scan for ANSI escape sequences. Redirect stderr to a file and scan for ANSI escape sequences. Verify none are present in either case. |
| **Disability impact** | Visual |
| **WCAG mapping** | — |

#### CV-5: `TERM=dumb` Respect

| | |
|---|---|
| **Level** | A |
| **Testability** | `[AUTO]` |
| **Requirement** | When the `TERM` environment variable is set to `dumb`, the tool MUST suppress ALL escape sequences — this includes SGR color and styling codes (`\x1b[...m`), CSI cursor movement (`\x1b[A`–`\x1b[D`, `\x1b[H`), screen/line clearing (`\x1b[2J`, `\x1b[K`), OSC sequences including terminal hyperlinks (`\x1b]8;;...`), and any other ECMA-48 control sequences. `TERM=dumb` signals a terminal with zero capability for escape sequence interpretation. When `TERM` is empty or unset, the tool SHOULD behave as if `TERM=dumb`. |
| **Test method** | Run the tool with `TERM=dumb` and capture raw output bytes. Scan for ALL escape sequence types: CSI sequences (`\x1b\[`), OSC sequences (`\x1b\]`), and other escape types. Verify none are present. Also run with `TERM=` (empty) and verify the same behavior. |
| **Disability impact** | Visual |
| **WCAG mapping** | — |

#### CV-6: `--color` Flag with Modes

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[AUTO]` |
| **Requirement** | The tool SHOULD support a `--color=WHEN` flag accepting at least three values: `always` (force color regardless of TTY/environment), `never` (suppress color), and `auto` (default behavior — enable color only when output is a TTY and `NO_COLOR`/`TERM=dumb` are not set). Tools MAY use non-standard flag names for the same concept (e.g., GCC uses `-fdiagnostics-color=auto/always/never`). Such tools satisfy the intent of this criterion if the flag provides equivalent auto/always/never modes. |
| **Test method** | Test `--color=never` (verify no ANSI color codes), `--color=always` with piped output (verify ANSI color codes present), and `--color=auto` with both TTY and non-TTY stdout. |
| **Disability impact** | Visual, Motor |
| **WCAG mapping** | — |

#### CV-7: `FORCE_COLOR` Support

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[AUTO]` |
| **Requirement** | When the `FORCE_COLOR` environment variable is set and non-empty, the tool SHOULD force color output even when stdout is not a TTY. `FORCE_COLOR` overrides `NO_COLOR` in the precedence chain. Per the force-color.org specification. |
| **Test method** | Pipe stdout through `cat` with `FORCE_COLOR=1` set. Verify ANSI color codes ARE present in piped output. Set both `NO_COLOR=1` and `FORCE_COLOR=1` and verify color is present (`FORCE_COLOR` wins). |
| **Disability impact** | Visual |
| **WCAG mapping** | — |

#### CV-8: 4-Bit ANSI Color Preference

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[SEMI]` |
| **Requirement** | The tool SHOULD prefer 4-bit ANSI colors (the 16 standard colors: codes 30–37, 40–47, 90–97, 100–107) over 8-bit (256-color) or 24-bit (truecolor) palettes. 4-bit colors are user-customizable in terminal emulator settings, giving users control over contrast and palette. When 8-bit or 24-bit colors are used, they SHOULD be limited to non-essential decorative elements. |
| **Test method** | Capture ANSI color codes from output. Classify as 4-bit (`\x1b\[\d+m`), 8-bit (`\x1b\[38;5;\d+m`), or 24-bit (`\x1b\[38;2;\d+;\d+;\d+m`). Flag predominant use of 8-bit or 24-bit for informational content. |
| **Disability impact** | Visual |
| **WCAG mapping** | 1.4.3 Contrast (Minimum) |

#### CV-9: No Background Color Assumption

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[SEMI]` |
| **Requirement** | The tool MUST NOT assume a specific terminal background color (light or dark). Color choices MUST be readable against both light and dark backgrounds, or the tool MUST detect and adapt to the terminal theme. The safest approach is to use ANSI colors that contrast well against any background, or to rely solely on the terminal's default foreground color. |
| **Test method** | Semi-automated: extract foreground colors used and evaluate whether any would produce low contrast against common light (#FFFFFF) or dark (#000000) backgrounds. Flag colors that fail 4.5:1 contrast ratio against either. |
| **Disability impact** | Visual |
| **WCAG mapping** | 1.4.3 Contrast (Minimum) |

#### CV-10: Configuration Precedence

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[AUTO]` |
| **Requirement** | Color configuration MUST follow this precedence order (highest to lowest): command-line flag (`--color`) > `FORCE_COLOR` environment variable > `NO_COLOR` environment variable > `TERM=dumb` > application-specific config > TTY auto-detection default. |
| **Test method** | Test conflict scenarios: set `NO_COLOR=1` and run with `--color=always` (flag should win). Set `NO_COLOR=1` and `FORCE_COLOR=1` (`FORCE_COLOR` should win). Set `TERM=dumb` with `NO_COLOR` unset (should suppress escapes). |
| **Disability impact** | Visual |
| **WCAG mapping** | — |

#### CV-11: High Contrast Mode Support

| | |
|---|---|
| **Level** | AAA |
| **Testability** | `[MANUAL]` |
| **Requirement** | The tool SHOULD produce readable output when the operating system's high-contrast mode is enabled. No information should be lost when terminal colors are overridden by the user or the OS. |
| **Test method** | Manual: enable OS high-contrast mode, run the tool, verify all output remains readable and no information is conveyed solely through colors that have been overridden. |
| **Disability impact** | Visual |
| **WCAG mapping** | 1.4.3 Contrast (Minimum) |

#### CV-12: Bold/Underline as Structural Cues

| | |
|---|---|
| **Level** | AAA |
| **Testability** | `[SEMI]` |
| **Requirement** | When using text styling (bold, underline, inverse) to convey structure (headings, emphasis, selection), the same structure MUST also be discernible without styling — via indentation, prefixes, blank lines, or whitespace patterns. |
| **Test method** | Run the tool with `TERM=dumb` (stripping all styling). Compare structural readability with styled output. Verify that sections, headings, and emphasis remain distinguishable through text formatting alone. |
| **Disability impact** | Visual, Cognitive |
| **WCAG mapping** | 1.3.1 Info and Relationships |

#### CV-13: Terminal Hyperlink Accessibility

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[SEMI]` |
| **Requirement** | When the tool emits OSC 8 terminal hyperlinks (`\x1b]8;;URI\a text \x1b]8;;\a`), the linked text MUST be meaningful on its own without the hyperlink — the URL MUST NOT be the only way to discover the link target. OSC 8 sequences MUST be suppressed when `TERM=dumb` is set and SHOULD be suppressed when `NO_COLOR` is set. Tools MUST NOT rely on hyperlink functionality as the only means of providing a reference, since screen readers and many terminal emulators do not support OSC 8. |
| **Test method** | Run with `TERM=dumb` and scan output for OSC 8 sequences (`\x1b]8;`). Run with `NO_COLOR=1` and scan. If hyperlinks are present in default output, verify the linked text is descriptive (not just "click here" or "link"). Verify the information is accessible without the hyperlink. |
| **Disability impact** | Visual |
| **WCAG mapping** | 2.4.4 Link Purpose (In Context) |

#### CV-14: Foreground-Background Pair Contrast

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[SEMI]` |
| **Requirement** | When the tool sets both foreground AND background colors in its output (e.g., colored status badges, highlighted sections, LS_COLORS-style file type indicators), the contrast ratio between the foreground and background colors MUST be at least 4.5:1 for normal text. This applies to 4-bit ANSI color pairs (using xterm default RGB mappings), fixed 8-bit color pairs, and 24-bit color pairs. The tool MUST NOT produce foreground-background combinations that fail this ratio in its default configuration. |
| **Test method** | Capture output and extract all SGR sequences that set both foreground (30–37, 90–97, 38;5;N, 38;2;R;G;B) and background (40–47, 100–107, 48;5;N, 48;2;R;G;B) colors. Map 4-bit ANSI codes to their xterm default RGB values. Calculate contrast ratio using the WCAG relative luminance formula. Flag any pair below 4.5:1. |
| **Disability impact** | Visual |
| **WCAG mapping** | 1.4.3 Contrast (Minimum) |

#### CV-15: Unicode/Emoji Symbol Accessibility

| | |
|---|---|
| **Level** | A |
| **Testability** | `[SEMI]` |
| **Requirement** | Unicode symbols and emoji used as status indicators, progress markers, or informational icons (e.g., ✓, ✗, ●, ▶, ⚠, 🔴) MUST be accompanied by a text alternative that conveys the same meaning. Screen readers announce the Unicode character name (e.g., "HEAVY CHECK MARK"), which may not convey the intended semantic meaning. Symbols that render as boxes or tofu (□) in terminals without appropriate fonts further degrade the experience. A `--plain` or `--ascii` mode SHOULD be available to replace Unicode symbols with ASCII text alternatives (e.g., `[OK]` instead of `✓`, `[FAIL]` instead of `✗`). |
| **Test method** | Scan output for common Unicode status symbols: checkmarks (U+2713–U+2717), circles (U+25CF, U+25CB), arrows (U+25B6, U+25C0, U+2192), warning signs (U+26A0), and emoji (U+1F300–U+1F9FF). For each symbol found, verify adjacent text provides equivalent meaning. Flag bare symbols without text context. |
| **Disability impact** | Visual, Cognitive |
| **WCAG mapping** | 1.1.1 Non-text Content |

---

### 4.3 Help & Documentation

This domain covers the discoverability, structure, and accessibility of help text, documentation, and usage guidance.

#### HD-1: `--help` and `-h` Flags

| | |
|---|---|
| **Level** | A |
| **Testability** | `[AUTO]` |
| **Requirement** | The tool MUST respond to both `--help` and `-h` with usage information sent to stdout, and MUST exit with code 0. |
| **Test method** | Run `tool --help` and `tool -h`. Verify exit code 0 for both. Verify stdout is non-empty and contains recognizable help content (usage patterns, flag descriptions, or command lists). |
| **Disability impact** | Visual, Cognitive |
| **WCAG mapping** | 3.3.2 Labels or Instructions |

#### HD-2: Subcommand Help

| | |
|---|---|
| **Level** | A |
| **Testability** | `[AUTO]` |
| **Requirement** | Every subcommand MUST support `--help` and `-h`, producing help specific to that subcommand. For git-style tools (tools with `help` as a subcommand), `tool help subcommand` MUST also work. |
| **Test method** | Enumerate subcommands from top-level help. For each, run `tool subcommand --help` and verify exit code 0 and non-empty stdout. For git-style tools, also test `tool help subcommand`. |
| **Disability impact** | Visual, Cognitive |
| **WCAG mapping** | 3.3.2 Labels or Instructions |

#### HD-3: `--version` Flag

| | |
|---|---|
| **Level** | A |
| **Testability** | `[AUTO]` |
| **Requirement** | The tool MUST respond to `--version` (and optionally `-V`) by printing a version string to stdout and exiting with code 0. |
| **Test method** | Run `tool --version`. Verify exit code 0 and stdout contains a version-like string (semver pattern, date, or version identifier). |
| **Disability impact** | Cognitive |
| **WCAG mapping** | — |

#### HD-4: Missing-Argument Guidance

| | |
|---|---|
| **Level** | A |
| **Testability** | `[AUTO]` |
| **Requirement** | When invoked with missing required arguments, the tool MUST print a concise error message to stderr identifying the missing argument(s) and a pointer to `--help`, and MUST exit with a non-zero code. The tool MUST NOT print the full help text to stderr on simple usage errors. |
| **Test method** | Identify a command requiring arguments (from help text parsing). Run without arguments. Verify: non-zero exit code, stderr contains error text, stderr references `--help` or `-h`. |
| **Disability impact** | Cognitive |
| **WCAG mapping** | 3.3.1 Error Identification |

#### HD-5: Help Text Structure

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[SEMI]` |
| **Requirement** | Help text MUST include: (1) a one-line description of the tool/command, (2) a usage synopsis showing the invocation pattern, (3) a list of available subcommands (if applicable), (4) a list of flags with descriptions, and (5) version or contact/support information. Sections MUST be separated by blank lines for screen reader navigation. |
| **Test method** | Parse `--help` output. Check for presence of: description line, usage/synopsis pattern, flag list, section separators (blank lines). Flag missing sections for human review. |
| **Disability impact** | Visual, Cognitive |
| **WCAG mapping** | 2.4.6 Headings and Labels, 1.3.1 Info and Relationships |

#### HD-6: Examples in Help

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[SEMI]` |
| **Requirement** | Help text SHOULD include at least one usage example demonstrating a common invocation. Examples SHOULD appear early in the help output or in a clearly labeled "Examples" section. |
| **Test method** | Scan `--help` output for example patterns (lines starting with `$`, `>`, or containing the tool name with arguments in a context suggesting example usage). Flag absence for human review. |
| **Disability impact** | Cognitive |
| **WCAG mapping** | — |

#### HD-7: Consistent Flag Descriptions

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[SEMI]` |
| **Requirement** | Flag descriptions MUST follow a consistent format: short flag, long flag, description, and default value (if any). Alignment SHOULD use spaces, not tabs. All flags across all subcommands SHOULD follow the same formatting convention. |
| **Test method** | Parse flag listings from `--help` output. Verify consistent alignment and format across entries. Flag inconsistencies for human review. |
| **Disability impact** | Visual, Cognitive |
| **WCAG mapping** | 3.2.4 Consistent Identification |

#### HD-8: No Hang on Empty Stdin

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[AUTO]` |
| **Requirement** | When the tool expects piped stdin input and stdin is a TTY (no data piped), the tool MUST either display help/usage and exit, or clearly prompt the user — not hang silently waiting for input indefinitely. |
| **Test method** | Run the tool in a scenario where it reads stdin, with stdin connected to a TTY (or `/dev/null`), with a timeout. Verify the tool either exits, displays help, or prompts within 5 seconds. |
| **Disability impact** | Motor, Cognitive |
| **WCAG mapping** | — |

#### HD-9: Man Page Availability

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[AUTO]` |
| **Requirement** | The tool SHOULD provide a man page installable via the system package manager. Man pages have a standardized, predictable, searchable format that screen reader users have specialized tooling for (w3mman, Emacspeak man-page scripts). The man page MUST match the current version of the tool. |
| **Test method** | Run `man tool-name` and check exit code. Verify man page exists and contains content. Cross-reference version in man page with `--version` output. |
| **Disability impact** | Visual |
| **WCAG mapping** | — |

#### HD-10: Typo Suggestion

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[SEMI]` |
| **Requirement** | When the user provides an unrecognized command, subcommand, or flag, the tool SHOULD suggest the closest valid alternative (e.g., "Did you mean `--format`?"). |
| **Test method** | Run the tool with a slightly misspelled subcommand or flag. Check if stderr contains a suggestion. Flag absence for human review. |
| **Disability impact** | Cognitive |
| **WCAG mapping** | 3.3.3 Error Suggestion |

#### HD-11: Web Documentation Link

| | |
|---|---|
| **Level** | AAA |
| **Testability** | `[SEMI]` |
| **Requirement** | Help text SHOULD include a URL to web-based documentation, which can be accessed in a browser with full web accessibility support (resizable text, screen reader compatibility, searchability). The URL SHOULD deep-link to the relevant subcommand page where possible. |
| **Test method** | Scan `--help` output for URLs (regex for http/https). Verify at least one documentation URL is present. |
| **Disability impact** | Visual, Cognitive |
| **WCAG mapping** | — |

#### HD-12: Shell Completion Support

| | |
|---|---|
| **Level** | AAA |
| **Testability** | `[AUTO]` |
| **Requirement** | The tool SHOULD provide shell completion scripts for at least bash and zsh, enabling command and flag discoverability without reading help text. Completions reduce the need to memorize commands and flags — critical for users with cognitive disabilities and efficient for screen reader users. |
| **Test method** | Check for a `completion` or `completions` subcommand, or `--generate-completions`, `--completion`, or `--bash-completion` flags. Verify at least one produces output. |
| **Disability impact** | Visual, Motor |
| **WCAG mapping** | — |

#### HD-13: `whatis`/`apropos` Compatibility

| | |
|---|---|
| **Level** | AAA |
| **Testability** | `[AUTO]` |
| **Requirement** | The man page (if provided) MUST include a proper NAME section following the `name - description` format so that `whatis` and `apropos` can index the tool for system-wide discoverability. |
| **Test method** | Run `whatis tool-name`. Verify it returns a result (exit code 0) with a description. |
| **Disability impact** | Visual |
| **WCAG mapping** | — |

---

### 4.4 Error Handling & Feedback

This domain covers how CLI tools communicate errors, warnings, and operational status to users.

#### EF-1: Errors to Stderr

| | |
|---|---|
| **Level** | A |
| **Testability** | `[AUTO]` |
| **Requirement** | All error messages MUST be sent to stderr, not stdout. This ensures errors are visible even when stdout is piped, and prevents error text from corrupting downstream tool input. |
| **Test method** | Trigger an error (invalid flag, missing file, etc.) with stdout redirected to a file. Verify error message appears on stderr (captured separately) and NOT in the stdout file. |
| **Disability impact** | Visual |
| **WCAG mapping** | — |

#### EF-2: Non-Zero Exit on Error

| | |
|---|---|
| **Level** | A |
| **Testability** | `[AUTO]` |
| **Requirement** | Any error condition MUST result in a non-zero exit code. The tool MUST NOT exit 0 when it has failed to perform the requested operation. |
| **Test method** | Trigger known error conditions (invalid input, nonexistent resource, permission error). Verify non-zero exit code for each. |
| **Disability impact** | Motor, Cognitive |
| **WCAG mapping** | — |

#### EF-3: Human-Readable Error Messages

| | |
|---|---|
| **Level** | A |
| **Testability** | `[SEMI]` |
| **Requirement** | Error messages MUST be written in plain language, identifying: (1) what went wrong, and (2) what the user can do about it. The tool MUST NOT present raw exception traces, regex patterns, memory addresses, or internal identifiers as the primary error message. Technical details MAY be available via `--debug` or `--verbose`. |
| **Test method** | Trigger errors and capture stderr. Verify messages contain natural language text (not just codes or traces). Flag messages shorter than 10 characters or containing stack trace patterns for human review. |
| **Disability impact** | Cognitive |
| **WCAG mapping** | 3.3.1 Error Identification |

#### EF-4: Distinct Error Codes

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[AUTO]` |
| **Requirement** | Different categories of failure (e.g., invalid input, network error, permission denied, resource not found) SHOULD produce distinct non-zero exit codes. Exit code meanings SHOULD be documented in help text or man page. |
| **Test method** | Trigger different error categories and collect exit codes. Verify at least two distinct non-zero codes are used for different error types. Check if help/man page documents exit codes. |
| **Disability impact** | Motor |
| **WCAG mapping** | — |

#### EF-5: Actionable Error Suggestions

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[SEMI]` |
| **Requirement** | Error messages SHOULD include a suggested corrective action or next step where determinable (e.g., "Run `tool auth login` to authenticate" rather than just "Authentication failed"). |
| **Test method** | Trigger common errors (auth failure, missing config, invalid format). Check if stderr contains actionable text (commands to run, flags to add, files to check). Flag generic-only messages for human review. |
| **Disability impact** | Motor, Cognitive |
| **WCAG mapping** | 3.3.3 Error Suggestion |

#### EF-6: Warning Distinction

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[SEMI]` |
| **Requirement** | Warnings (non-fatal issues) MUST be visually and programmatically distinguishable from errors (fatal issues). The tool MUST use consistent prefixes (e.g., `Warning:` vs `Error:`) and distinct exit behavior (0 for warnings-only, non-zero for errors). |
| **Test method** | Trigger warning conditions and error conditions separately. Verify distinct prefixes and exit codes. Verify warnings go to stderr. |
| **Disability impact** | Visual, Cognitive |
| **WCAG mapping** | — |

#### EF-7: Error Aggregation

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[SEMI]` |
| **Requirement** | When multiple errors occur in a single invocation, the tool SHOULD group or summarize them rather than printing repetitive per-item error lines. A summary count with representative examples is preferable to hundreds of identical or near-identical lines. |
| **Test method** | Trigger batch operations with multiple failures. Count error lines. Flag if the tool produces more than 20 identical-format error lines without a summary. |
| **Disability impact** | Cognitive |
| **WCAG mapping** | — |

#### EF-8: No Jargon in Default Errors

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[SEMI]` |
| **Requirement** | Default error messages MUST avoid internal jargon, stack traces, memory addresses, hex values, or developer-facing debug information. The tool SHOULD provide a `--debug` or `--verbose` flag that reveals technical details when needed. |
| **Test method** | Trigger errors and scan stderr for patterns indicating jargon: stack trace frames, hex addresses (0x...), internal function names, unformatted JSON error objects. Flag occurrences for human review. |
| **Disability impact** | Cognitive |
| **WCAG mapping** | — |

#### EF-9: Structured Error Output

| | |
|---|---|
| **Level** | AAA |
| **Testability** | `[AUTO]` |
| **Requirement** | When `--json` mode is active, errors SHOULD also be emitted in structured JSON format with fields for error code, message, and suggested action. This enables assistive technology tools and scripts to parse and re-present errors accessibly. |
| **Test method** | Trigger an error with `--json` flag. Verify the output (stdout or stderr) contains a JSON object with error information fields. |
| **Disability impact** | Visual |
| **WCAG mapping** | — |

#### EF-10: Bug Report Facilitation

| | |
|---|---|
| **Level** | AAA |
| **Testability** | `[SEMI]` |
| **Requirement** | For unexpected or internal errors, the tool SHOULD provide instructions for reporting the bug, ideally with a URL pre-populated with diagnostic information (OS, tool version, error details). |
| **Test method** | Check if error output contains a URL or reporting instructions when an unexpected error occurs. |
| **Disability impact** | Motor, Cognitive |
| **WCAG mapping** | — |

---

### 4.5 Interactivity & Input

This domain covers how CLI tools handle user input, prompts, interactive elements, and confirmations.

#### II-1: Full Non-Interactive Mode

| | |
|---|---|
| **Level** | A |
| **Testability** | `[AUTO]` |
| **Requirement** | Every operation achievable through interactive prompts MUST also be achievable via flags and arguments alone. No functionality may be locked behind interactive-only prompts. |
| **Test method** | Identify interactive prompts from help text or documentation. For each, verify that equivalent flags exist. Run the tool with `--no-input` or equivalent flag and verify operations complete without prompting. |
| **Disability impact** | Motor |
| **WCAG mapping** | 2.1.1 Keyboard |

#### II-2: TTY-Aware Prompting

| | |
|---|---|
| **Level** | A |
| **Testability** | `[AUTO]` |
| **Requirement** | When stdin is not a TTY, the tool MUST NOT prompt for interactive input. It MUST either use default values, fail with a clear error identifying the required flag, or read from stdin as data input. |
| **Test method** | Run the tool with stdin redirected from `/dev/null` in a scenario requiring interactive input. Verify the tool does not hang waiting for input. Verify it either succeeds with defaults or fails with a helpful error. |
| **Disability impact** | Motor |
| **WCAG mapping** | — |

#### II-3: `--no-input` / `--non-interactive` Flag

| | |
|---|---|
| **Level** | A |
| **Testability** | `[AUTO]` |
| **Requirement** | The tool MUST provide a flag to explicitly disable all interactive prompts (e.g., `--no-input`, `--non-interactive`, `--batch`, or `--yes`). When set, the tool MUST use defaults or fail with guidance on which flags to provide. |
| **Test method** | Test `--no-input`, `--non-interactive`, `--batch`, and `--yes` flags. Verify at least one is accepted (exit 0). Run an operation that normally prompts with this flag and verify no prompt occurs. |
| **Disability impact** | Motor |
| **WCAG mapping** | — |

#### II-4: Ctrl-C Responsiveness

| | |
|---|---|
| **Level** | A |
| **Testability** | `[AUTO]` |
| **Requirement** | The tool MUST respond to SIGINT (Ctrl-C) by exiting promptly. The tool MUST NOT trap SIGINT and silently ignore it. A brief cleanup phase (up to 5 seconds) is acceptable, but a second SIGINT MUST force immediate exit. |
| **Test method** | Start the tool in a long-running operation. Send SIGINT. Verify the tool exits within 5 seconds. If it does not exit, send a second SIGINT and verify immediate exit. |
| **Disability impact** | Motor |
| **WCAG mapping** | 2.1.1 Keyboard |

#### II-5: Password/Secret Masking

| | |
|---|---|
| **Level** | A |
| **Testability** | `[AUTO]` |
| **Requirement** | When prompting for passwords or secrets interactively, the tool MUST disable terminal echo. Input MUST NOT be displayed as the user types. |
| **Test method** | If the tool has an authentication flow, trigger a password prompt. Verify that terminal echo is disabled during input (check that the terminal's ECHO flag is cleared via `stty` or equivalent). |
| **Disability impact** | Visual |
| **WCAG mapping** | — |

#### II-6: Confirmation for Destructive Actions

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[SEMI]` |
| **Requirement** | Irreversible or destructive operations (delete, overwrite, format, purge) MUST require explicit confirmation in interactive mode, or a `--force` / `-f` flag in non-interactive mode. The confirmation prompt MUST describe what will be destroyed. |
| **Test method** | Identify destructive commands from help text. Run without `--force` and verify a confirmation prompt appears (or the tool requires `--force`). Verify the prompt describes the destructive action. |
| **Disability impact** | Cognitive |
| **WCAG mapping** | 3.3.4 Error Prevention |

#### II-7: Prompt Labels and Context

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[SEMI]` |
| **Requirement** | Interactive prompts MUST clearly state: (1) what is being asked, (2) what the valid options are, and (3) what the default value is (if any). A screen reader user hearing only the prompt text must understand what to input without seeing surrounding visual context. |
| **Test method** | Trigger interactive prompts. Verify each prompt contains a question/label, lists valid options (for choice prompts), and indicates the default value (e.g., `[Y/n]`, `(default: 8080)`). |
| **Disability impact** | Visual, Cognitive |
| **WCAG mapping** | 3.3.2 Labels or Instructions |

#### II-8: `--dry-run` Support

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[AUTO]` |
| **Requirement** | For commands that modify state (files, remote resources, configuration), the tool SHOULD provide a `--dry-run`, `--whatif`, or `--simulate` flag that shows what would happen without performing the action. |
| **Test method** | Test `--dry-run`, `--whatif`, and `--simulate` flags. Verify at least one is accepted for state-modifying commands. Verify the tool produces descriptive output without actually modifying state. |
| **Disability impact** | Cognitive |
| **WCAG mapping** | 3.3.4 Error Prevention |

#### II-9: Keyboard-Navigable Selections

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[MANUAL]` |
| **Requirement** | Interactive menus and selection lists MUST be navigable using arrow keys, and MUST announce the currently selected item in a way accessible to screen readers. Selection indicators MUST NOT rely solely on redraw-based visual changes that screen readers cannot detect. |
| **Test method** | Manual: use a screen reader to navigate an interactive selection list. Verify each item is announced when focused. Verify the selection can be confirmed with Enter. |
| **Disability impact** | Visual, Motor |
| **WCAG mapping** | 2.1.1 Keyboard |

#### II-10: Stdin as File with `-`

| | |
|---|---|
| **Level** | AAA |
| **Testability** | `[AUTO]` |
| **Requirement** | When the tool accepts file input, it SHOULD support `-` as a filename alias for stdin, enabling piped input as an alternative to interactive file selection or argument passing. |
| **Test method** | If the tool accepts a filename argument, test with `-` as the filename while piping data on stdin. Verify the tool processes the piped data. |
| **Disability impact** | Motor |
| **WCAG mapping** | — |

#### II-11: Selection State Announcement

| | |
|---|---|
| **Level** | AAA |
| **Testability** | `[MANUAL]` |
| **Requirement** | Interactive multi-select prompts MUST announce checked/unchecked state to assistive technology. Selection state MUST NOT rely solely on visual checkmarks, color changes, or character substitution that screen readers cannot interpret. |
| **Test method** | Manual: use a screen reader to interact with a multi-select prompt. Verify checked/unchecked state is announced for each item. Verify toggling selection is announced. |
| **Disability impact** | Visual |
| **WCAG mapping** | 4.1.2 Name, Role, Value |

---

### 4.6 Environment Awareness & Configuration

This domain covers how CLI tools detect, respect, and interact with the user's environment, configuration, and assistive technology signals.

#### EC-1: `NO_COLOR` Respect

| | |
|---|---|
| **Level** | A |
| **Testability** | `[AUTO]` |
| **Requirement** | Cross-reference CV-2. When `NO_COLOR` is set and non-empty, suppress ANSI color codes. |
| **Note** | This is a cross-reference to CV-2, listed here because environment awareness is a distinct concern from color handling. Evaluation counts once. |

#### EC-2: `TERM=dumb` Respect

| | |
|---|---|
| **Level** | A |
| **Testability** | `[AUTO]` |
| **Requirement** | Cross-reference CV-5. When `TERM=dumb`, suppress all ANSI escape sequences. |
| **Note** | Cross-reference to CV-5. Evaluation counts once. |

#### EC-3: Configuration Precedence

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[AUTO]` |
| **Requirement** | Configuration MUST follow a documented, predictable precedence order: CLI flags (highest) > environment variables > project-level config (e.g., `.toolrc` in project directory) > user-level config (e.g., `~/.config/tool/config`) > system-level config > built-in defaults (lowest). The precedence order MUST be documented in help text or man page. |
| **Test method** | Set conflicting configuration at multiple levels (flag + env var, env var + config file). Verify the higher-precedence source wins. Check help/man page for documentation of precedence. |
| **Disability impact** | Cognitive |
| **WCAG mapping** | — |

#### EC-4: XDG Base Directory Compliance

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[AUTO]` |
| **Requirement** | On systems supporting the XDG Base Directory Specification, configuration files SHOULD respect `XDG_CONFIG_HOME` (default `~/.config`), data files SHOULD respect `XDG_DATA_HOME` (default `~/.local/share`), and cache files SHOULD respect `XDG_CACHE_HOME` (default `~/.cache`), rather than creating dotfiles/dotdirs directly in `$HOME`. |
| **Test method** | Set `XDG_CONFIG_HOME` to a custom directory. Run the tool in a way that creates config. Verify config is created under the custom directory, not in `$HOME`. |
| **Disability impact** | Cognitive |
| **WCAG mapping** | — |

#### EC-5: `PAGER` Respect

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[AUTO]` |
| **Requirement** | When invoking a pager for long output, the tool MUST respect the `PAGER` environment variable. Users may set this to a pager compatible with their assistive technology (e.g., `most`, a custom script, or `cat` to disable paging). |
| **Test method** | Set `PAGER=cat` and run a command that would normally page. Verify all output goes to stdout without invoking the system default pager. |
| **Disability impact** | Visual |
| **WCAG mapping** | — |

#### EC-6: `EDITOR`/`VISUAL` Respect

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[AUTO]` |
| **Requirement** | When launching an external editor (for commit messages, config editing, interactive input), the tool MUST respect the `VISUAL` environment variable (preferred for terminal editors with full-screen capability) and fall back to `EDITOR`. Users with disabilities may need a specific accessible editor (e.g., Emacs with Emacspeak). |
| **Test method** | Set `VISUAL` or `EDITOR` to a custom script that logs its invocation. Trigger an editor launch. Verify the custom script was invoked. |
| **Disability impact** | Visual, Motor |
| **WCAG mapping** | — |

#### EC-7: Locale and `LANG` Respect

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[AUTO]` |
| **Requirement** | The tool SHOULD respect `LANG`, `LC_ALL`, and `LC_MESSAGES` for output language where translations are available. Error messages and help text SHOULD follow the user's locale when supported. At minimum, the tool MUST correctly handle UTF-8 locales without producing mojibake. |
| **Test method** | Set `LANG=C` and run the tool. Verify output is in English (C locale). If the tool supports translations, set a supported locale and verify translated output. Set a UTF-8 locale and verify output contains valid UTF-8. |
| **Disability impact** | Cognitive |
| **WCAG mapping** | 3.1.1 Language of Page |

#### EC-8: No Secrets in Environment Variables

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[SEMI]` |
| **Requirement** | The tool SHOULD NOT require secrets (passwords, API tokens, private keys) to be passed via environment variables as the primary or only mechanism. Environment variables are visible via `/proc/<pid>/environ`, `ps eww`, and crash dumps. The tool SHOULD prefer credential files (with appropriate permissions), stdin pipes, OS keychain integration, or dedicated secret management tools. |
| **Test method** | Review help text and documentation for authentication mechanisms. Flag if the only documented method for providing secrets is an environment variable. |
| **Disability impact** | Cross-cutting (security) |
| **WCAG mapping** | — |

#### EC-9: No Telemetry Without Consent

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[SEMI]` |
| **Requirement** | The tool MUST NOT transmit usage data, crash reports, or analytics without explicit user opt-in. First-run disclosure is acceptable only if the user is prompted and can decline before any data is sent. The opt-out mechanism MUST be documented and accessible. |
| **Test method** | Run the tool for the first time. Monitor network traffic. Verify no outbound connections are made for telemetry without user consent. Check help/docs for telemetry opt-out instructions. |
| **Disability impact** | Cognitive |
| **WCAG mapping** | — |

#### EC-10: Proxy Environment Respect

| | |
|---|---|
| **Level** | AAA |
| **Testability** | `[AUTO]` |
| **Requirement** | The tool SHOULD respect `HTTP_PROXY`, `HTTPS_PROXY`, `ALL_PROXY`, and `NO_PROXY` environment variables for network operations. Users behind corporate proxies or using accessibility-related network routing depend on these standard variables. |
| **Test method** | Set `HTTPS_PROXY` to a non-functional proxy. Run a command requiring network access. Verify the tool attempts to use the proxy (connection error mentioning proxy) rather than connecting directly. |
| **Disability impact** | Cross-cutting |
| **WCAG mapping** | — |

#### EC-11: `COLUMNS`/`LINES` Respect

| | |
|---|---|
| **Level** | AAA |
| **Testability** | `[AUTO]` |
| **Requirement** | The tool SHOULD respect `COLUMNS` and `LINES` environment variables for output formatting, in addition to terminal ioctl detection. Users with screen magnification may set smaller values to ensure output fits their visible area. |
| **Test method** | Set `COLUMNS=40` and `LINES=10`. Run the tool and verify output respects these dimensions (line length, pager behavior). |
| **Disability impact** | Visual |
| **WCAG mapping** | — |

---

### 4.7 Timing & Motion

This domain covers animations, progress indicators, timeouts, and content that changes over time.

#### TM-1: No Animation When Not TTY

| | |
|---|---|
| **Level** | A |
| **Testability** | `[AUTO]` |
| **Requirement** | When stdout is not a TTY, the tool MUST NOT emit cursor-movement escape sequences, screen-clearing codes, spinner animations, carriage return (`\r`) based overwrites, or any content that relies on overwriting previous output. |
| **Test method** | Pipe stdout through `cat` during a long-running operation. Scan for carriage return characters (`\r`) not followed by newline, cursor movement sequences (`\x1b\[\d*[ABCDHJ]`), and screen clearing sequences. Verify none are present. |
| **Disability impact** | Visual |
| **WCAG mapping** | — |

#### TM-2: Static Progress Alternative

| | |
|---|---|
| **Level** | A |
| **Testability** | `[AUTO]` |
| **Requirement** | If the tool displays animated progress indicators (spinners, progress bars using carriage returns, braille-character animations), it MUST provide a static alternative. This may be activated via `--no-animation`, `--progress=plain`, an environment variable, or automatic detection when stdout is not a TTY. The static alternative SHOULD emit periodic line-based progress messages (e.g., `Processing: 45% complete`). |
| **Test method** | Run a long operation with stdout piped. Verify progress is reported as line-based messages (newline-terminated) rather than in-place overwrites. Test `--no-animation` flag if present. |
| **Disability impact** | Visual |
| **WCAG mapping** | 2.2.2 Pause, Stop, Hide |

#### TM-3: No Flashing/Strobing Content

| | |
|---|---|
| **Level** | A |
| **Testability** | `[SEMI]` |
| **Requirement** | Output MUST NOT flash or strobe more than 3 times per second. Rapidly alternating inverse video, color cycling, blinking text attributes (`\x1b\[5m`), or fast cursor movement creating visual strobing can trigger photosensitive seizures. |
| **Test method** | Scan output for ANSI blink attribute (`\x1b\[5m` or `\x1b\[6m`). Monitor update frequency of carriage-return-based overwrites. Flag if overwrites exceed 3 per second. |
| **Disability impact** | Vestibular/Photosensitive |
| **WCAG mapping** | 2.3.1 Three Flashes or Below Threshold |

#### TM-4: Configurable Timeouts

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[AUTO]` |
| **Requirement** | Network timeouts, operation timeouts, and any time-limited interactions MUST be configurable via flags or environment variables. Default timeouts MUST be reasonable (not less than 30 seconds for network operations). Users with motor disabilities or who use assistive technology may need more time. |
| **Test method** | Check help text for timeout-related flags (e.g., `--timeout`, `--connect-timeout`). Verify timeout flags are accepted and modify behavior. |
| **Disability impact** | Motor |
| **WCAG mapping** | 2.2.1 Timing Adjustable |

#### TM-5: Progress Indication for Long Operations

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[SEMI]` |
| **Requirement** | Operations taking longer than 2 seconds SHOULD provide some form of progress feedback on stderr — at minimum a static "Working..." message, ideally a percentage, item count, or elapsed time. Silence during long operations is inaccessible: screen reader users cannot distinguish "working" from "frozen," and all users benefit from feedback. |
| **Test method** | Run a known long operation and monitor stderr within the first 5 seconds. Verify some progress output appears. Flag complete silence for human review. |
| **Disability impact** | Visual, Cognitive |
| **WCAG mapping** | — |

#### TM-6: Responsive Startup

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[AUTO]` |
| **Requirement** | The tool SHOULD produce its first output (even if just the beginning of a progress message) within 500ms of invocation. The `--help` flag MUST respond within 500ms. Slow startup disproportionately affects users relying on assistive technology feedback loops — a screen reader user waiting for `--help` output cannot tell if the tool is loading or has crashed. |
| **Test method** | Time `tool --help` execution. Verify total time is under 500ms. Time initial output of a data command. Verify first byte appears within 500ms. |
| **Disability impact** | Visual, Cognitive |
| **WCAG mapping** | — |

#### TM-7: Pause/Resume for Scrolling Output

| | |
|---|---|
| **Level** | AAA |
| **Testability** | `[MANUAL]` |
| **Requirement** | For streaming output (log tailing, watch modes, continuous monitoring), the tool SHOULD support pausing output so users with screen readers can catch up. This may be implemented via Ctrl-S/Ctrl-Q flow control, a pause keystroke, or a `--page` flag that buffers output. |
| **Test method** | Manual: during streaming output, attempt Ctrl-S to pause and Ctrl-Q to resume. Verify output pauses and resumes correctly without data loss. |
| **Disability impact** | Visual, Motor |
| **WCAG mapping** | 2.2.1 Timing Adjustable |

---

### 4.8 Internationalization & Localization

This domain covers language support, character encoding, and bidirectional text handling.

#### IL-1: UTF-8 Output Support

| | |
|---|---|
| **Level** | A |
| **Testability** | `[AUTO]` |
| **Requirement** | The tool MUST produce valid UTF-8 output when the locale indicates UTF-8 encoding (via `LANG`, `LC_ALL`, or `LC_CTYPE`). Mojibake (garbled characters from encoding mismatches) blocks access for users in non-Latin-script locales and corrupts screen reader output. |
| **Test method** | Set `LANG=en_US.UTF-8`. Run the tool and validate output is valid UTF-8 (no invalid byte sequences). If the tool produces non-ASCII characters, verify they are correctly encoded. |
| **Disability impact** | Cognitive |
| **WCAG mapping** | — |

#### IL-2: No Hardcoded English-Only Errors

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[SEMI]` |
| **Requirement** | Error messages and user-facing strings SHOULD be externalizable for translation via standard i18n mechanisms (gettext, message catalogs, resource bundles, i18n libraries). Even if translations are not yet available, the architecture should not prevent future localization. |
| **Test method** | Semi-automated: check if the tool uses gettext or equivalent (look for `.po`/`.mo` files, i18n library dependencies, or message catalog references in the binary). Flag hardcoded strings for human review. |
| **Disability impact** | Cognitive |
| **WCAG mapping** | 3.1.1 Language of Page |

#### IL-3: Bidirectional Text Safety

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[SEMI]` |
| **Requirement** | When displaying user-provided content that may contain right-to-left (RTL) text (Arabic, Hebrew, Persian, Urdu), the tool SHOULD NOT corrupt the reading order. At minimum, the tool MUST NOT strip or mangle Unicode bidirectional control characters (U+200E–U+200F, U+202A–U+202E, U+2066–U+2069). |
| **Test method** | Pipe RTL text (Arabic or Hebrew) through the tool. Verify the RTL characters appear intact in output (byte-level comparison). Verify BiDi control characters are preserved. |
| **Disability impact** | Cognitive |
| **WCAG mapping** | 1.3.2 Meaningful Sequence |

#### IL-4: Unicode Box-Drawing Correctness

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[AUTO]` |
| **Requirement** | If the tool uses Unicode box-drawing characters (U+2500–U+257F), they MUST be used correctly — forming complete, unbroken borders. The tool MUST NOT mix ASCII pseudo-borders (`+`, `-`, `|`) with Unicode box-drawing characters, and MUST NOT break borders with inline descriptive text. Accessible terminal emulators and screen readers detect box structures from these characters; inconsistent use produces garbled output. |
| **Test method** | Scan output for box-drawing characters. Verify they form consistent, complete structures. Check for mixed ASCII/Unicode border characters. |
| **Disability impact** | Visual |
| **WCAG mapping** | — |

#### IL-5: Locale-Aware Formatting

| | |
|---|---|
| **Level** | AAA |
| **Testability** | `[AUTO]` |
| **Requirement** | Numeric, date, and time output in human-readable mode SHOULD respect locale formatting conventions (e.g., decimal separators, date order). Machine-readable modes (`--json`, `--csv`) SHOULD use ISO 8601 dates and unformatted numbers for interoperability. |
| **Test method** | Set locale to one with distinct formatting (e.g., `de_DE.UTF-8` for comma decimal separator). Run the tool and check if numeric/date formatting follows locale conventions. Verify `--json` output uses ISO 8601 dates. |
| **Disability impact** | Cognitive |
| **WCAG mapping** | — |

---

### 4.9 Installation & Lifecycle

This domain covers accessibility of the installation, update, and uninstallation process.

#### IN-1: Accessible Installation Process

| | |
|---|---|
| **Level** | A |
| **Testability** | `[SEMI]` |
| **Requirement** | The installation process (whether via package manager, install script, or binary download) MUST be completable using only a keyboard and screen reader. Piped-to-shell install scripts (`curl | sh`) MUST NOT require interactive input without providing flag alternatives. Installation instructions MUST be available as accessible text (not only in images or videos). |
| **Test method** | Review the tool's documented installation methods. For script-based installers, run with `--help` to check for non-interactive flags. Verify at least one installation method does not require mouse interaction. |
| **Disability impact** | Visual |
| **WCAG mapping** | — |

#### IN-2: Version Check Mechanism

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[AUTO]` |
| **Requirement** | The tool SHOULD provide a way to check for available updates without modifying the system (e.g., `tool version --check`, `tool update --dry-run`, or a non-destructive update notification). |
| **Test method** | Check help text for update-check related commands or flags. Test if a version check command exists and runs without modifying the installation. |
| **Disability impact** | Cognitive |
| **WCAG mapping** | — |

#### IN-3: Deprecation Warnings in Output

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[SEMI]` |
| **Requirement** | When features, flags, or behaviors are deprecated, the tool MUST emit clear warnings to stderr that: (1) name the deprecated element, (2) state what replaces it, and (3) indicate the removal timeline. Deprecation warnings MUST NOT be silently suppressed — users who cannot read changelogs or release notes depend on in-tool warnings. |
| **Test method** | If deprecated features are known, invoke them and verify a deprecation warning appears on stderr. Verify the warning names a replacement and a timeline. |
| **Disability impact** | Cognitive |
| **WCAG mapping** | — |

#### IN-4: Non-Destructive Updates

| | |
|---|---|
| **Level** | AAA |
| **Testability** | `[SEMI]` |
| **Requirement** | Update mechanisms SHOULD preserve user configuration, shell completions, and customizations. The tool MUST NOT silently overwrite user config files during updates. If a config migration is needed, the tool SHOULD notify the user and provide a migration path. |
| **Test method** | If the tool has a self-update mechanism, create a custom config file, run the update, and verify the config file is preserved. |
| **Disability impact** | Cognitive |
| **WCAG mapping** | — |

---

## 5. TUI Extension Module

This module defines additional criteria for full-screen Terminal User Interface applications. A TUI tool claims conformance by meeting applicable core CLI-ACS criteria PLUS these TUI-specific criteria.

### Applicability

The TUI Extension Module applies when the tool:

- Manages the full terminal screen (clears screen, positions cursor at arbitrary locations)
- Uses a curses/ncurses library or equivalent for screen management
- Provides interactive elements (menus, forms, panels, dialogs) within the terminal

Tools that only use simple line-based prompts, progress bars, or single-line interactive elements do NOT need the TUI Extension — they are covered by the core CLI criteria.

### TUI Criteria

#### TU-1: Keyboard-Only Full Operability

| | |
|---|---|
| **Level** | A |
| **Testability** | `[MANUAL]` |
| **Requirement** | Every function of the TUI MUST be operable via keyboard alone. No operation may require a mouse click as the only input method. |
| **Disability impact** | Motor |
| **WCAG mapping** | 2.1.1 Keyboard |

#### TU-2: Visible Focus Indicator

| | |
|---|---|
| **Level** | A |
| **Testability** | `[MANUAL]` |
| **Requirement** | The currently focused element (cursor position, selected row, active pane, highlighted menu item) MUST be visually distinguishable at all times. The indicator MUST NOT rely solely on color — it MUST use inverse video, brackets, underline, cursor positioning, or a combination. |
| **Disability impact** | Visual |
| **WCAG mapping** | 2.4.7 Focus Visible |

#### TU-3: Logical Focus Order

| | |
|---|---|
| **Level** | A |
| **Testability** | `[MANUAL]` |
| **Requirement** | Tab/focus traversal between UI regions (panes, panels, dialogs, form fields) MUST follow a logical, predictable order consistent with the visual layout. Focus MUST NOT jump to unexpected regions or skip interactive elements. |
| **Disability impact** | Visual, Cognitive |
| **WCAG mapping** | 2.4.3 Focus Order |

#### TU-4: No Keyboard Trap

| | |
|---|---|
| **Level** | A |
| **Testability** | `[MANUAL]` |
| **Requirement** | The user MUST be able to navigate away from any component using standard keyboard mechanisms. Modal dialogs MUST provide a clear, documented escape mechanism (e.g., Esc, q, Ctrl-C). No component may trap keyboard focus indefinitely. |
| **Disability impact** | Motor |
| **WCAG mapping** | 2.1.2 No Keyboard Trap |

#### TU-5: Screen Reader Region Announcements

| | |
|---|---|
| **Level** | A |
| **Testability** | `[MANUAL]` |
| **Requirement** | When the TUI updates a region of the screen (status bar change, new content in a pane, notification), the update SHOULD be detectable by screen readers. The tool SHOULD avoid full-screen redraws when only a portion has changed — use targeted cursor movement and writes to minimize the area that screen readers must re-scan. |
| **Disability impact** | Visual |
| **WCAG mapping** | 4.1.3 Status Messages |

#### TU-6: Discoverable Keybindings

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[SEMI]` |
| **Requirement** | All keyboard shortcuts MUST be documented and discoverable within the application — via a help screen (typically `?` or `F1`), a persistent status bar hint showing common keys, or an in-app key legend. |
| **Disability impact** | Cognitive |
| **WCAG mapping** | — |

#### TU-7: Rebindable Keys

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[SEMI]` |
| **Requirement** | Keyboard shortcuts SHOULD be user-configurable via a configuration file. Users with motor disabilities or alternative input devices may need to remap keys to avoid chords (Ctrl+Shift+X), hard-to-reach keys, or keys that conflict with their assistive technology. |
| **Disability impact** | Motor |
| **WCAG mapping** | — |

#### TU-8: No Single-Character-Only Shortcuts for Critical Actions

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[SEMI]` |
| **Requirement** | Destructive or irreversible actions (delete, quit without save, format) MUST NOT be triggered by a single unmodified keypress alone (e.g., just `d` or `q`). Require a modifier key (Ctrl, Alt), a confirmation prompt, or a multi-key sequence (e.g., `dd` in vim). Single-character shortcuts for non-destructive actions (navigation, scrolling) are acceptable. |
| **Disability impact** | Motor, Cognitive |
| **WCAG mapping** | 2.1.4 Character Key Shortcuts |

#### TU-9: Minimum Contrast in Custom Themes

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[MANUAL]` |
| **Requirement** | If the TUI defines its own color theme (using specific color values rather than terminal default colors), the contrast ratio between foreground text and background MUST be at least 4.5:1 for normal text and 3:1 for bold or heading text. Where possible, the TUI SHOULD defer to terminal colors to give users maximum control. |
| **Disability impact** | Visual |
| **WCAG mapping** | 1.4.3 Contrast (Minimum) |

#### TU-10: Resize/Reflow Handling

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[MANUAL]` |
| **Requirement** | The TUI MUST handle terminal resize events (SIGWINCH) gracefully — reflowing content, adjusting layout proportions, and maintaining focus position. Content MUST NOT be clipped, lost, or rendered unreadable after a resize. The TUI SHOULD announce or visually indicate the resize. |
| **Disability impact** | Visual, Cognitive |
| **WCAG mapping** | 1.4.10 Reflow |

#### TU-11: Mouse-Optional Enhancement

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[MANUAL]` |
| **Requirement** | If mouse support is provided, it MUST be a supplement to keyboard operation, not a replacement. Every mouse-accessible function MUST have a keyboard equivalent. Mouse support SHOULD be disableable (some terminal screen readers conflict with mouse capture mode). |
| **Disability impact** | Motor |
| **WCAG mapping** | 2.1.1 Keyboard |

#### TU-12: Alternative View Modes

| | |
|---|---|
| **Level** | AAA |
| **Testability** | `[MANUAL]` |
| **Requirement** | The TUI SHOULD provide alternative display modes for complex visualizations (graphs, dashboards, tree views). For example: a list/table view as an alternative to a graphical chart, a flat list alternative to a tree hierarchy, or a summary view alternative to a detailed dashboard. |
| **Disability impact** | Visual |
| **WCAG mapping** | 1.1.1 Non-text Content |

#### TU-13: Spatial Navigation Cues

| | |
|---|---|
| **Level** | AAA |
| **Testability** | `[MANUAL]` |
| **Requirement** | The TUI SHOULD provide non-visual cues for spatial layout — announcing pane names, positions (e.g., "left pane: file list, right pane: preview"), boundaries, and the number of panes when the user navigates between regions. This enables screen reader users to build a mental model of the layout. |
| **Disability impact** | Visual |
| **WCAG mapping** | 1.3.1 Info and Relationships |

#### TU-14: Configurable Refresh Rate

| | |
|---|---|
| **Level** | AAA |
| **Testability** | `[AUTO]` |
| **Requirement** | For TUIs that periodically refresh their display (monitoring dashboards, log viewers, system monitors), the refresh interval SHOULD be user-configurable via flag or config. Rapid refreshes overwhelm screen readers (which must re-scan the changed screen) and can trigger vestibular discomfort. |
| **Disability impact** | Vestibular/Photosensitive |
| **WCAG mapping** | 2.2.2 Pause, Stop, Hide |

---

## 6. Disability Impact Matrix

### 6.1 Coverage by Disability Category

This matrix shows which criteria serve each disability category, enabling evaluators to assess a tool's accessibility for specific user groups.

#### Visual (Blind, Low Vision, Color Blind)

Users who rely on screen readers, screen magnifiers, braille displays, high-contrast modes, or custom color configurations.

**Level A:** OS-1, OS-3, OS-4, OS-5, CV-1, CV-2, CV-3, CV-4, CV-5, HD-1, HD-2, EF-1, II-5, TM-1, TM-2, IN-1 (16 criteria)

**Level AA:** OS-6, OS-8, CV-6, CV-7, CV-8, CV-9, CV-10, HD-5, HD-7, HD-9, EF-6, II-7, II-9, EC-5, EC-6, TM-5, TM-6, IL-4 (18 criteria)

**Level AAA:** OS-9, OS-10, CV-11, CV-12, HD-11, HD-12, HD-13, EF-9, II-11, EC-11, TM-7 (11 criteria)

**TUI Extension:** TU-2, TU-3, TU-5, TU-9, TU-10, TU-12, TU-13 (7 criteria)

#### Motor (Limited Dexterity, RSI, Switch Access, Voice Input)

Users who rely on keyboard-only navigation, switch access, voice input, sticky keys, or who experience fatigue from repetitive input.

**Level A:** OS-2, EF-2, II-1, II-2, II-3, II-4 (6 criteria)

**Level AA:** OS-6, CV-6, HD-8, EF-4, EF-5, II-9, EC-6, TM-4 (8 criteria)

**Level AAA:** OS-9, HD-12, EF-10, II-10, TM-7 (5 criteria)

**TUI Extension:** TU-1, TU-4, TU-7, TU-8, TU-11 (5 criteria)

#### Cognitive (Learning Disabilities, Attention Disorders, Neurodivergent)

Users who benefit from clear language, predictable behavior, reduced complexity, and consistent patterns.

**Level A:** OS-2, OS-4, OS-5, CV-1, HD-1, HD-2, HD-3, HD-4, EF-2, EF-3, IL-1 (11 criteria)

**Level AA:** OS-7, HD-5, HD-6, HD-7, HD-8, HD-10, EF-5, EF-6, EF-7, EF-8, II-6, II-7, II-8, EC-3, EC-4, EC-7, EC-9, TM-5, TM-6, IL-2, IL-3, IN-2, IN-3 (23 criteria)

**Level AAA:** OS-10, CV-12, HD-11, EF-10, IL-5, IN-4 (6 criteria)

**TUI Extension:** TU-3, TU-6, TU-8, TU-10 (4 criteria)

#### Auditory (Deaf, Hard of Hearing)

CLI tools are inherently non-auditory — they produce no sound by default. Auditory criteria score 0 not because of a gap but because the medium does not create auditory barriers. If a CLI tool produces terminal bell alerts (`\a`) or system sounds, criterion CV-1 (information not solely via one channel) provides coverage.

#### Vestibular / Photosensitive

Users susceptible to motion sickness, dizziness, or photosensitive seizures from flashing or rapidly changing content.

**Level A:** TM-3 (1 criterion)

**Level AAA:** TU-14 (1 criterion)

### 6.2 Summary Counts

| Category | Level A | Level AA | Level AAA | TUI Ext | Total |
|---|---|---|---|---|---|
| Visual | 16 | 18 | 11 | 7 | 52 |
| Motor | 6 | 8 | 5 | 5 | 24 |
| Cognitive | 11 | 23 | 6 | 4 | 44 |
| Auditory | 0 | 0 | 0 | 0 | 0 |
| Vestibular | 1 | 0 | 0 | 1 | 2 |

---

## 7. Conformance Report Template

### 7.1 Report Structure

The CLI-ACS Conformance Report (CLI-ACR) follows this structure:

#### Header Section

| Field | Description |
|---|---|
| **Product Name** | Full name and version of the CLI tool under evaluation |
| **Report Date** | Date the evaluation was completed (ISO 8601: YYYY-MM-DD) |
| **Evaluator** | Organization or individual who performed the evaluation |
| **Spec Version** | CLI-ACS version used (e.g., "CLI-ACS v1.0") |
| **Scope** | "Core CLI" or "Core CLI + TUI Extension" |
| **Methodology** | Description of testing methods: automated suite version, manual testing approaches, assistive technology used |

#### Conformance Summary

| Field | Description |
|---|---|
| **Overall Level** | Highest level fully met: A, AA, AAA, or Partial |
| **Per-Domain Breakdown** | Table showing pass/fail per domain per level |
| **Disability Impact Summary** | Percentage of criteria met per disability category |

#### Detailed Results

For each criterion:

| Field | Description |
|---|---|
| **ID** | Criterion identifier (e.g., OS-1, CV-2) |
| **Criterion** | Short name |
| **Level** | A, AA, or AAA |
| **Test Method** | `[AUTO]`, `[SEMI]`, or `[MANUAL]` |
| **Result** | Supports / Partially Supports / Does Not Support / Not Applicable / Not Evaluated |
| **Evidence** | Automated test output, manual observation notes, or screen reader test results |
| **Remarks** | Free-text explanation, required if result is not "Supports" |

#### Testing Environment

| Field | Description |
|---|---|
| **Operating System** | Name and version (e.g., Ubuntu 24.04 LTS, macOS 15.1, Windows 11 24H2) |
| **Terminal Emulator** | Name and version (e.g., GNOME Terminal 3.50, iTerm2 3.5, Windows Terminal 1.20) |
| **Shell** | Name and version (e.g., bash 5.2, zsh 5.9, PowerShell 7.4) |
| **Screen Reader** | Name and version (e.g., NVDA 2025.1, VoiceOver macOS 15, Orca 46) |
| **Conformance Suite** | CLI-ACS Suite version used for automated tests |

#### Standards Crosswalk (Optional)

If the report needs to serve as input for a VPAT or EN 301 549 conformance declaration, include the crosswalk table mapping CLI-ACS criteria to WCAG 2.2, EN 301 549, and Section 508.

### 7.2 Conformance Level Definitions

| Term | Definition |
|---|---|
| **Supports** | The tool meets the criterion in all tested scenarios with no known defects. |
| **Partially Supports** | Some functionality meets the criterion but at least one tested scenario fails or has known limitations. The evaluator MUST describe the gap in Remarks. |
| **Does Not Support** | The majority of tested scenarios fail to meet the criterion. The evaluator MUST describe the failure in Remarks. |
| **Not Applicable** | The criterion does not apply to this tool. The evaluator MUST document the rationale (e.g., "Tool has no subcommands" for HD-2). |
| **Not Evaluated** | The tool has not been tested against this criterion. The evaluator SHOULD explain why (e.g., "Manual screen reader testing not performed in this evaluation cycle"). |

### 7.3 Machine-Readable Report Format

The conformance suite SHOULD produce reports in both human-readable (Markdown or HTML) and machine-readable (JSON) formats. The JSON format enables aggregation, trending, and automated compliance checks.

```json
{
  "cli_acs_version": "1.0",
  "product": {
    "name": "example-tool",
    "version": "2.3.1"
  },
  "report_date": "2026-08-25",
  "evaluator": "Example Corp Accessibility Team",
  "scope": "core_cli",
  "overall_level": "AA",
  "environment": {
    "os": "Ubuntu 24.04 LTS",
    "terminal": "GNOME Terminal 3.50",
    "shell": "bash 5.2",
    "screen_reader": "NVDA 2025.1 (via RDP)",
    "suite_version": "cli-acs-suite 1.0.0"
  },
  "results": [
    {
      "id": "OS-1",
      "criterion": "Standard Stream Separation",
      "level": "A",
      "test_method": "AUTO",
      "result": "Supports",
      "evidence": "stdout and stderr correctly separated in all 12 tested commands",
      "remarks": null
    },
    {
      "id": "CV-2",
      "criterion": "NO_COLOR Support",
      "level": "A",
      "test_method": "AUTO",
      "result": "Does Not Support",
      "evidence": "ANSI color codes present in output with NO_COLOR=1",
      "remarks": "Tool does not check NO_COLOR environment variable. Color codes emitted regardless."
    }
  ],
  "summary": {
    "by_domain": {
      "output_structure": {"A": "pass", "AA": "pass", "AAA": "partial"},
      "color_visual": {"A": "fail", "AA": "not_evaluated", "AAA": "not_evaluated"}
    },
    "by_disability": {
      "visual": {"met": 35, "total": 42, "percentage": 83.3},
      "motor": {"met": 17, "total": 20, "percentage": 85.0},
      "cognitive": {"met": 28, "total": 36, "percentage": 77.8},
      "vestibular": {"met": 2, "total": 2, "percentage": 100.0}
    }
  }
}
```

---

## 8. Normative References & Standards Crosswalk

### 8.1 Normative References

| Standard | Version | Publisher | Relevance to CLI-ACS |
|---|---|---|---|
| [WCAG 2.2](https://www.w3.org/TR/WCAG22/) | W3C Recommendation, Oct 2023 | W3C/WAI | Primary web accessibility standard; CLI-ACS criteria map to applicable Success Criteria |
| [WCAG2ICT](https://www.w3.org/TR/wcag2ict-22/) | W3C Group Note, Oct 2024 | W3C/WAI | Guidance on applying WCAG 2.2 to non-web ICT including terminal applications |
| [EN 301 549](https://www.etsi.org/deliver/etsi_en/301500_301599/301549/03.02.01_60/en_301549v030201p.pdf) | v3.2.1, Mar 2021 | ETSI/CEN/CENELEC | European ICT accessibility standard; Chapter 11 covers non-web software |
| [Section 508 (Revised)](https://www.access-board.gov/ict/) | 36 CFR Part 1194, Jan 2018 | U.S. Access Board | U.S. federal ICT accessibility requirements |
| [VPAT 2.5](https://www.itic.org/policy/accessibility/vpat) | Nov 2023 | ITI | Accessibility Conformance Report template |
| [NO_COLOR](https://no-color.org/) | Informal, 2017 | no-color.org | Environment variable standard for disabling ANSI color |
| [FORCE_COLOR](https://force-color.org/) | Informal, 2023 | force-color.org | Environment variable standard for forcing ANSI color |
| [XDG Base Directory Spec](https://specifications.freedesktop.org/basedir-spec/latest/) | v0.8, May 2021 | freedesktop.org | Standard config/data/cache directory locations |
| [POSIX.1-2017](https://pubs.opengroup.org/onlinepubs/9699919799/) | IEEE Std 1003.1-2017 | IEEE/Open Group | Terminal, locale, signal, and environment variable standards |

### 8.2 Informative References

| Source | Description |
|---|---|
| Sampath, H., Merrick, A., & Macvean, A. (2021). [Accessibility of Command Line Interfaces](https://dl.acm.org/doi/abs/10.1145/3411764.3445544). CHI '21, ACM. DOI: 10.1145/3411764.3445544 | Seminal study on CLI accessibility with 12 screen reader users, identifying key barriers in unstructured text output |
| [Command Line Interface Guidelines](https://clig.dev/) (clig.dev) | Community guidelines for CLI design including accessibility, output formatting, and error handling |
| Seirdy. (2022). [Best practices for inclusive CLIs](https://seirdy.one/posts/2022/06/10/cli-best-practices/). | Comprehensive accessibility recommendations including espeak-ng testing methodology |
| GitHub Blog. (2025). [Building a more accessible GitHub CLI](https://github.blog/engineering/user-experience/building-a-more-accessible-github-cli/). | Practical implementation lessons from the `gh` CLI accessibility initiative |
| W3C. (2024). [Background on Text / Command-Line / Terminal Applications and Interfaces](https://github.com/w3c/wcag2ict/blob/main/background-on-text-command-line-terminal-applications-and-interfaces.md). | W3C working document explaining how terminal emulators serve as "user agents" for text applications |
| [GitHub CLI Accessibility Conformance Report](https://accessibility.github.com/conformance/cli/) | Example CLI-specific VPAT/ACR from GitHub |
| [npm CLI Accessibility Conformance Report](https://accessibility.github.com/conformance/npm-cli/) | Example CLI-specific VPAT/ACR from npm |
| [ImageMagick VPAT](https://imagemagick.org/script/vpat.php) | Example CLI tool VPAT |
| [Google Cloud CLI VPAT](https://cloud.google.com/security/compliance/vpat) | Example cloud CLI VPAT |

### 8.3 Standards Crosswalk

This table maps CLI-ACS criteria to their corresponding requirements in WCAG 2.2, EN 301 549, and Section 508. This enables organizations to generate VPAT-format reports from CLI-ACS evaluations and to demonstrate compliance with existing regulatory frameworks.

| CLI-ACS ID | CLI-ACS Criterion | WCAG 2.2 SC | EN 301 549 | Section 508 |
|---|---|---|---|---|
| OS-4 | Linear Reading Order | 1.3.2 Meaningful Sequence | 11.1.3.2 | Ch. 5, 502.3.1 |
| OS-5 | No ASCII Art Sole Info | 1.1.1 Non-text Content | 11.1.1.1 | Ch. 5, 502.3.1 |
| CV-1 | Color Not Sole Channel | 1.4.1 Use of Color | 11.1.4.1 | Ch. 5, 502.3.1 |
| CV-8 | 4-Bit ANSI Preference | 1.4.3 Contrast (Minimum) | 11.1.4.3 | Ch. 5, 502.3.1 |
| CV-9 | No Background Assumption | 1.4.3 Contrast (Minimum) | 11.1.4.3 | Ch. 5, 502.3.1 |
| CV-12 | Bold/Underline Structural | 1.3.1 Info and Relationships | 11.1.3.1 | Ch. 5, 502.3.1 |
| CV-13 | Terminal Hyperlink A11y | 2.4.4 Link Purpose (In Context) | 11.2.4.4 | Ch. 5, 502.3.1 |
| CV-14 | FG-BG Pair Contrast | 1.4.3 Contrast (Minimum) | 11.1.4.3 | Ch. 5, 502.3.1 |
| CV-15 | Unicode/Emoji Symbol A11y | 1.1.1 Non-text Content | 11.1.1.1 | Ch. 5, 502.3.1 |
| HD-1 | --help and -h | 3.3.2 Labels or Instructions | 11.3.3.2 | Ch. 5, 502.3.1 |
| HD-4 | Missing-Arg Guidance | 3.3.1 Error Identification | 11.3.3.1 | Ch. 5, 502.3.1 |
| HD-5 | Help Text Structure | 1.3.1, 2.4.6 | 11.1.3.1, 11.2.4.6 | Ch. 5, 502.3.1 |
| HD-7 | Consistent Flag Format | 3.2.4 Consistent Identification | 11.3.2.4 | Ch. 5, 502.3.1 |
| HD-10 | Typo Suggestion | 3.3.3 Error Suggestion | 11.3.3.3 | Ch. 5, 502.3.1 |
| EF-3 | Human-Readable Errors | 3.3.1 Error Identification | 11.3.3.1 | Ch. 5, 502.3.1 |
| EF-5 | Actionable Suggestions | 3.3.3 Error Suggestion | 11.3.3.3 | Ch. 5, 502.3.1 |
| II-1 | Non-Interactive Mode | 2.1.1 Keyboard | 11.2.1.1 | Ch. 5, 502.3.1 |
| II-4 | Ctrl-C Responsiveness | 2.1.1 Keyboard | 11.2.1.1 | Ch. 5, 502.3.1 |
| II-6 | Destructive Confirmation | 3.3.4 Error Prevention | 11.3.3.4 | Ch. 5, 502.3.1 |
| II-7 | Prompt Labels | 3.3.2 Labels or Instructions | 11.3.3.2 | Ch. 5, 502.3.1 |
| II-9 | Keyboard-Nav Selections | 2.1.1 Keyboard | 11.2.1.1 | Ch. 5, 502.3.1 |
| II-11 | Selection State Announce | 4.1.2 Name, Role, Value | 11.4.1.2 | Ch. 5, 502.3.1 |
| EC-7 | Locale Respect | 3.1.1 Language of Page | 11.3.1.1 | Ch. 5, 502.3.1 |
| TM-2 | Static Progress Alt | 2.2.2 Pause, Stop, Hide | 11.2.2.2 | Ch. 5, 502.3.1 |
| TM-3 | No Flashing Content | 2.3.1 Three Flashes | 11.2.3.1 | Ch. 5, 502.3.1 |
| TM-4 | Configurable Timeouts | 2.2.1 Timing Adjustable | 11.2.2.1 | Ch. 5, 502.3.1 |
| IL-3 | BiDi Text Safety | 1.3.2 Meaningful Sequence | 11.1.3.2 | Ch. 5, 502.3.1 |
| TU-1 | Keyboard-Only Operability | 2.1.1 Keyboard | 11.2.1.1 | Ch. 5, 502.3.1 |
| TU-2 | Visible Focus Indicator | 2.4.7 Focus Visible | 11.2.4.7 | Ch. 5, 502.3.1 |
| TU-3 | Logical Focus Order | 2.4.3 Focus Order | 11.2.4.3 | Ch. 5, 502.3.1 |
| TU-4 | No Keyboard Trap | 2.1.2 No Keyboard Trap | 11.2.1.2 | Ch. 5, 502.3.1 |
| TU-5 | Region Announcements | 4.1.3 Status Messages | 11.4.1.3 | Ch. 5, 502.3.1 |
| TU-8 | No Single-Char Critical | 2.1.4 Character Key Shortcuts | 11.2.1.4 | — |
| TU-9 | Minimum Contrast | 1.4.3 Contrast (Minimum) | 11.1.4.3 | Ch. 5, 502.3.1 |
| TU-10 | Resize/Reflow | 1.4.10 Reflow | 11.1.4.10 | Ch. 5, 502.3.1 |
| TU-12 | Alternative View Modes | 1.1.1 Non-text Content | 11.1.1.1 | Ch. 5, 502.3.1 |
| TU-13 | Spatial Navigation Cues | 1.3.1 Info and Relationships | 11.1.3.1 | Ch. 5, 502.3.1 |
| TU-14 | Configurable Refresh Rate | 2.2.2 Pause, Stop, Hide | 11.2.2.2 | Ch. 5, 502.3.1 |

### 8.4 Criteria Without Direct Standards Mapping

The following CLI-ACS criteria address CLI-specific accessibility concerns that are not directly covered by WCAG 2.2, EN 301 549, or Section 508. These represent the unique value of a CLI-native accessibility standard:

| CLI-ACS ID | Criterion | Rationale |
|---|---|---|
| OS-1 | Stream Separation | CLI-specific: stdout/stderr distinction has no web equivalent |
| OS-2 | Meaningful Exit Codes | CLI-specific: programmatic success/failure signaling |
| OS-3 | Clean Piped Output | CLI-specific: pipe and redirect behavior |
| OS-6 | Machine-Readable Output | CLI-specific: `--json` / `--csv` for AT tool consumption |
| OS-7 | Quiet Mode | CLI-specific: output verbosity control |
| OS-8 | Width Awareness | CLI-specific: terminal column awareness |
| OS-9 | Pager Support | CLI-specific: long output handling |
| CV-2 | NO_COLOR Support | CLI-specific: de facto color control standard |
| CV-3 | --no-color Flag | CLI-specific: per-invocation color control |
| CV-4 | TTY-Aware Color | CLI-specific: TTY detection for formatting |
| CV-5 | TERM=dumb Respect | CLI-specific: terminal capability detection |
| CV-6 | --color Flag | CLI-specific: color mode control |
| CV-7 | FORCE_COLOR Support | CLI-specific: de facto color forcing standard |
| CV-10 | Config Precedence | CLI-specific: flag/env/config hierarchy |
| HD-2 | Subcommand Help | CLI-specific: subcommand discoverability |
| HD-3 | --version Flag | CLI-specific: version identification |
| HD-8 | No Hang on Empty Stdin | CLI-specific: TTY/pipe stdin behavior |
| HD-9 | Man Page Availability | CLI-specific: system documentation format |
| HD-12 | Shell Completion | CLI-specific: shell integration |
| HD-13 | whatis/apropos Compat | CLI-specific: system-wide discoverability |
| EF-1 | Errors to Stderr | CLI-specific: stream discipline |
| EF-4 | Distinct Error Codes | CLI-specific: exit code semantics |
| EF-7 | Error Aggregation | CLI-specific: batch error handling |
| EF-9 | Structured Error Output | CLI-specific: JSON error format |
| II-1 | Non-Interactive Mode | CLI-specific: flag-based alternatives to prompts |
| II-2 | TTY-Aware Prompting | CLI-specific: stdin TTY detection |
| II-3 | --no-input Flag | CLI-specific: interactive mode control |
| II-5 | Password Masking | CLI-specific: terminal echo control |
| II-8 | --dry-run Support | CLI-specific: safe preview mode |
| II-10 | Stdin as File | CLI-specific: `-` convention |
| EC-3 | Config Precedence | CLI-specific: configuration layering |
| EC-4 | XDG Compliance | CLI-specific: filesystem conventions |
| EC-5 | PAGER Respect | CLI-specific: pager integration |
| EC-6 | EDITOR/VISUAL Respect | CLI-specific: editor integration |
| EC-8 | No Secrets in Env | CLI-specific: environment security |
| EC-9 | No Telemetry w/o Consent | CLI-specific: privacy control |
| EC-10 | Proxy Respect | CLI-specific: network configuration |
| EC-11 | COLUMNS/LINES Respect | CLI-specific: terminal dimension awareness |
| TM-1 | No Animation When !TTY | CLI-specific: TTY-aware animation |
| TM-5 | Progress Indication | CLI-specific: long operation feedback |
| TM-6 | Responsive Startup | CLI-specific: help response time |
| TM-7 | Pause/Resume Scrolling | CLI-specific: flow control |
| IL-1 | UTF-8 Output | CLI-specific: encoding correctness |
| IL-2 | Externalizable Strings | CLI-specific: i18n architecture |
| IL-4 | Box-Drawing Correctness | CLI-specific: Unicode border integrity |
| IL-5 | Locale-Aware Formatting | CLI-specific: number/date formatting |
| IN-1 | Accessible Install | CLI-specific: install script accessibility |
| IN-2 | Version Check | CLI-specific: update discoverability |
| IN-3 | Deprecation Warnings | CLI-specific: in-tool deprecation |
| IN-4 | Non-Destructive Updates | CLI-specific: config preservation |
| TU-6 | Discoverable Keybindings | TUI-specific: in-app key help |
| TU-7 | Rebindable Keys | TUI-specific: key remapping |
| TU-11 | Mouse-Optional | TUI-specific: mouse/keyboard parity |

---

## Appendix A: Criteria Summary Table

| ID | Criterion | Level | Test | Visual | Motor | Cognitive | Vestibular |
|---|---|---|---|---|---|---|---|
| **Output Structure & Content** |||||||
| OS-1 | Standard Stream Separation | A | AUTO | ● | | | |
| OS-2 | Meaningful Exit Codes | A | AUTO | | ● | ● | |
| OS-3 | Clean Piped Output | A | AUTO | ● | | | |
| OS-4 | Linear Reading Order | A | SEMI | ● | | ● | |
| OS-5 | No ASCII Art Sole Info | A | SEMI | ● | | ● | |
| OS-6 | Machine-Readable Output | AA | AUTO | ● | ● | | |
| OS-7 | Quiet Mode | AA | AUTO | | | ● | |
| OS-8 | Output Width Awareness | AA | AUTO | ● | | | |
| OS-9 | Pager Support | AAA | AUTO | ● | ● | | |
| OS-10 | Plain Text Alternative | AAA | AUTO | ● | | ● | |
| **Color & Visual Presentation** |||||||
| CV-1 | Color Not Sole Channel | A | SEMI | ● | | ● | |
| CV-2 | NO_COLOR Support | A | AUTO | ● | | | |
| CV-3 | --no-color Flag | A | AUTO | ● | | | |
| CV-4 | TTY-Aware Color | A | AUTO | ● | | | |
| CV-5 | TERM=dumb Respect | A | AUTO | ● | | | |
| CV-6 | --color Flag with Modes | AA | AUTO | ● | ● | | |
| CV-7 | FORCE_COLOR Support | AA | AUTO | ● | | | |
| CV-8 | 4-Bit ANSI Preference | AA | SEMI | ● | | | |
| CV-9 | No Background Assumption | AA | SEMI | ● | | | |
| CV-10 | Configuration Precedence | AA | AUTO | ● | | | |
| CV-11 | High Contrast Support | AAA | MANUAL | ● | | | |
| CV-12 | Bold/Underline Structural | AAA | SEMI | ● | | ● | |
| CV-13 | Terminal Hyperlink A11y | AA | SEMI | ● | | | |
| CV-14 | FG-BG Pair Contrast | AA | SEMI | ● | | | |
| CV-15 | Unicode/Emoji Symbol A11y | A | SEMI | ● | | ● | |
| **Help & Documentation** |||||||
| HD-1 | --help and -h Flags | A | AUTO | ● | | ● | |
| HD-2 | Subcommand Help | A | AUTO | ● | | ● | |
| HD-3 | --version Flag | A | AUTO | | | ● | |
| HD-4 | Missing-Argument Guidance | A | AUTO | | | ● | |
| HD-5 | Help Text Structure | AA | SEMI | ● | | ● | |
| HD-6 | Examples in Help | AA | SEMI | | | ● | |
| HD-7 | Consistent Flag Format | AA | SEMI | ● | | ● | |
| HD-8 | No Hang on Empty Stdin | AA | AUTO | | ● | ● | |
| HD-9 | Man Page Availability | AA | AUTO | ● | | | |
| HD-10 | Typo Suggestion | AA | SEMI | | | ● | |
| HD-11 | Web Documentation Link | AAA | SEMI | ● | | ● | |
| HD-12 | Shell Completion Support | AAA | AUTO | ● | ● | | |
| HD-13 | whatis/apropos Compat | AAA | AUTO | ● | | | |
| **Error Handling & Feedback** |||||||
| EF-1 | Errors to Stderr | A | AUTO | ● | | | |
| EF-2 | Non-Zero Exit on Error | A | AUTO | | ● | ● | |
| EF-3 | Human-Readable Errors | A | SEMI | | | ● | |
| EF-4 | Distinct Error Codes | AA | AUTO | | ● | | |
| EF-5 | Actionable Suggestions | AA | SEMI | | ● | ● | |
| EF-6 | Warning Distinction | AA | SEMI | ● | | ● | |
| EF-7 | Error Aggregation | AA | SEMI | | | ● | |
| EF-8 | No Jargon in Errors | AA | SEMI | | | ● | |
| EF-9 | Structured Error Output | AAA | AUTO | ● | | | |
| EF-10 | Bug Report Facilitation | AAA | SEMI | | ● | ● | |
| **Interactivity & Input** |||||||
| II-1 | Full Non-Interactive Mode | A | AUTO | | ● | | |
| II-2 | TTY-Aware Prompting | A | AUTO | | ● | | |
| II-3 | --no-input Flag | A | AUTO | | ● | | |
| II-4 | Ctrl-C Responsiveness | A | AUTO | | ● | | |
| II-5 | Password/Secret Masking | A | AUTO | ● | | | |
| II-6 | Destructive Confirmation | AA | SEMI | | | ● | |
| II-7 | Prompt Labels and Context | AA | SEMI | ● | | ● | |
| II-8 | --dry-run Support | AA | AUTO | | | ● | |
| II-9 | Keyboard-Nav Selections | AA | MANUAL | ● | ● | | |
| II-10 | Stdin as File with - | AAA | AUTO | | ● | | |
| II-11 | Selection State Announce | AAA | MANUAL | ● | | | |
| **Environment Awareness & Config** |||||||
| EC-1 | NO_COLOR Respect | A | AUTO | ● | | | |
| EC-2 | TERM=dumb Respect | A | AUTO | ● | | | |
| EC-3 | Configuration Precedence | AA | AUTO | | | ● | |
| EC-4 | XDG Compliance | AA | AUTO | | | ● | |
| EC-5 | PAGER Respect | AA | AUTO | ● | | | |
| EC-6 | EDITOR/VISUAL Respect | AA | AUTO | ● | ● | | |
| EC-7 | Locale Respect | AA | AUTO | | | ● | |
| EC-8 | No Secrets in Env | AA | SEMI | | | | |
| EC-9 | No Telemetry w/o Consent | AA | SEMI | | | ● | |
| EC-10 | Proxy Respect | AAA | AUTO | | | | |
| EC-11 | COLUMNS/LINES Respect | AAA | AUTO | ● | | | |
| **Timing & Motion** |||||||
| TM-1 | No Animation When !TTY | A | AUTO | ● | | | |
| TM-2 | Static Progress Alt | A | AUTO | ● | | | |
| TM-3 | No Flashing Content | A | SEMI | | | | ● |
| TM-4 | Configurable Timeouts | AA | AUTO | | ● | | |
| TM-5 | Progress Indication | AA | SEMI | ● | | ● | |
| TM-6 | Responsive Startup | AA | AUTO | ● | | ● | |
| TM-7 | Pause/Resume Scrolling | AAA | MANUAL | ● | ● | | |
| **Internationalization** |||||||
| IL-1 | UTF-8 Output Support | A | AUTO | | | ● | |
| IL-2 | No Hardcoded English | AA | SEMI | | | ● | |
| IL-3 | BiDi Text Safety | AA | SEMI | | | ● | |
| IL-4 | Unicode Box-Drawing | AA | AUTO | ● | | | |
| IL-5 | Locale-Aware Formatting | AAA | AUTO | | | ● | |
| **Installation & Lifecycle** |||||||
| IN-1 | Accessible Install | A | SEMI | ● | | | |
| IN-2 | Version Check Mechanism | AA | AUTO | | | ● | |
| IN-3 | Deprecation Warnings | AA | SEMI | | | ● | |
| IN-4 | Non-Destructive Updates | AAA | SEMI | | | ● | |
| **TUI Extension Module** |||||||
| TU-1 | Keyboard-Only Operability | A | MANUAL | | ● | | |
| TU-2 | Visible Focus Indicator | A | MANUAL | ● | | | |
| TU-3 | Logical Focus Order | A | MANUAL | ● | | ● | |
| TU-4 | No Keyboard Trap | A | MANUAL | | ● | | |
| TU-5 | Region Announcements | A | MANUAL | ● | | | |
| TU-6 | Discoverable Keybindings | AA | SEMI | | | ● | |
| TU-7 | Rebindable Keys | AA | SEMI | | ● | | |
| TU-8 | No Single-Char Critical | AA | SEMI | | ● | ● | |
| TU-9 | Minimum Contrast | AA | MANUAL | ● | | | |
| TU-10 | Resize/Reflow Handling | AA | MANUAL | ● | | ● | |
| TU-11 | Mouse-Optional Enhancement | AA | MANUAL | | ● | | |
| TU-12 | Alternative View Modes | AAA | MANUAL | ● | | | |
| TU-13 | Spatial Navigation Cues | AAA | MANUAL | ● | | | |
| TU-14 | Configurable Refresh Rate | AAA | AUTO | | | | ● |

---

## Appendix B: Testability Summary

| Classification | Core CLI | TUI Extension | Total |
|---|---|---|---|
| `[AUTO]` — Fully Automated | 49 | 1 | 50 |
| `[SEMI]` — Semi-Automated | 31 | 3 | 34 |
| `[MANUAL]` — Manual AT Testing | 4 | 10 | 14 |
| **Total** | **84** | **14** | **98** |

Note: EC-1 and EC-2 are cross-references to CV-2 and CV-5 respectively. They are counted once (under Color & Visual Presentation, as CV-2 and CV-5). The unique criteria count is:

- **Core CLI:** 84 unique criteria (28 Level A, 40 Level AA, 16 Level AAA)
- **TUI Extension:** 14 unique criteria (5 Level A, 6 Level AA, 3 Level AAA)
- **Grand Total:** 98 unique testable criteria

Of these, 50 (51%) are fully automatable by a conformance suite, enabling meaningful accessibility evaluation of any CLI binary without requiring manual testing for over half the criteria.
