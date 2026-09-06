# CLI-ACS Domain 4.2: Color & Visual Presentation

**Specification:** CLI Accessibility Conformance Specification (CLI-ACS) v1.0
**Domain:** 4.2 — Color & Visual Presentation
**Status:** Draft
**Last updated:** 2026-09-06

This document provides the detailed specification for the Color & Visual Presentation domain of CLI-ACS, covering how CLI tools use color, contrast, and visual styling in terminal output. It includes technical background, normative criteria, test methods, and references.

---

## Table of Contents

1. [Background: How Terminal Color Works](#1-background-how-terminal-color-works)
   - 1.1 [The Three Color Tiers](#11-the-three-color-tiers)
   - 1.2 [SGR Text Attributes (Non-Color Styling)](#12-sgr-text-attributes-non-color-styling)
   - 1.3 [Color Capability Detection](#13-color-capability-detection)
   - 1.4 [Color Vision Deficiency (Color Blindness)](#14-color-vision-deficiency-color-blindness)
   - 1.5 [Contrast in Terminal Contexts](#15-contrast-in-terminal-contexts)
   - 1.6 [Screen Readers and Terminal Styling](#16-screen-readers-and-terminal-styling)
2. [Criteria](#2-criteria)
   - [CV-1: Color Not Sole Information Channel](#cv-1-color-not-sole-information-channel)
   - [CV-2: `NO_COLOR` Support](#cv-2-no_color-support)
   - [CV-3: `--no-color` Flag](#cv-3---no-color-flag)
   - [CV-4: TTY-Aware Color](#cv-4-tty-aware-color)
   - [CV-5: `TERM=dumb` Respect](#cv-5-termdumb-respect)
   - [CV-6: `--color` Flag with Modes](#cv-6---color-flag-with-modes)
   - [CV-7: `FORCE_COLOR` Support](#cv-7-force_color-support)
   - [CV-8: 4-Bit ANSI Color Preference](#cv-8-4-bit-ansi-color-preference)
   - [CV-9: No Background Color Assumption](#cv-9-no-background-color-assumption)
   - [CV-10: Color Configuration Precedence](#cv-10-color-configuration-precedence)
   - [CV-11: High Contrast Mode Support](#cv-11-high-contrast-mode-support)
   - [CV-12: Bold/Underline as Structural Cues](#cv-12-boldunderline-as-structural-cues)
3. [References](#3-references)

---

## 1. Background: How Terminal Color Works

Understanding terminal color architecture is essential for evaluating CLI accessibility, because the accessibility properties of color in terminals differ fundamentally from color on the web.

### 1.1 The Three Color Tiers

Terminal color operates in three tiers, each with different accessibility implications:

#### 4-Bit ANSI Colors (16 Colors)

SGR codes 30-37 (foreground), 40-47 (background), 90-97 (bright foreground), 100-107 (bright background). These provide 8 base colors (black, red, green, yellow, blue, magenta, cyan, white) plus 8 bright variants.

Critically, these are **named slots, not fixed RGB values** — the terminal emulator decides what "red" actually looks like. A user's theme maps the 16 names to actual hex values, meaning the same ANSI "red" renders as different hues in different color schemes. This indirection is the foundation of terminal color accessibility:

- Users with color vision deficiency can remap every named color to values they can perceive.
- Users requiring high contrast can ensure sufficient contrast ratios between all colors.
- Users with photosensitivity can choose muted palettes.
- Accessibility-focused terminal color schemes like Tempus (WCAG AA-compliant, 4.5:1 minimum contrast) work automatically with any tool that uses ANSI 16 colors.

All of this happens without any application cooperation — the tool writes "red" and the user's theme decides what that means.

#### 8-Bit Colors (256 Colors)

Accessed via `\x1b[38;5;Nm` (foreground) and `\x1b[48;5;Nm` (background). The 256-color palette comprises three regions:

| Index Range | Count | Description | User-Customizable? |
|---|---|---|---|
| 0-15 | 16 | Mirror the ANSI 16 slots | Yes — same as 4-bit |
| 16-231 | 216 | Fixed 6x6x6 RGB color cube | **No** — theme-independent |
| 232-255 | 24 | Fixed grayscale ramp | **No** — theme-independent |

The RGB cube uses the formula: `index = 16 + 36*r + 6*g + b` where r, g, b are in {0, 1, 2, 3, 4, 5}. The six channel levels map to intensity values: 0, 95, 135, 175, 215, 255.

Indices 16-255 are **fixed hex values that bypass the user's color scheme entirely**. A hardcoded `\x1b[38;5;196m` renders as `#FF0000` regardless of the user's accessibility preferences, cannot be remapped, and may fail against the user's background color.

#### 24-Bit Truecolor (16.7 Million Colors)

Accessed via `\x1b[38;2;R;G;Bm` (foreground) and `\x1b[48;2;R;G;Bm` (background), where R, G, B are decimal 0-255. These specify exact RGB values and are completely theme-independent. Supported by most terminal emulators released after 2018 (GNOME Terminal, iTerm2, Windows Terminal, Konsole, kitty, Alacritty, WezTerm).

The 24-bit color syntax originates from a reading of ISO/IEC 8613-6 (ITU T.416), referenced by ECMA-48 SGR code 38 ("reserved for future standardization; intended for setting character foreground colour as specified in ISO 8613-6"). The semicolon-delimited syntax (`38;2;R;G;B`) is a de facto standard despite technical arguments that colons should be used as sub-parameter separators per ECMA-48 Section 5.4.2.

#### Accessibility Comparison

| Property | ANSI 16 (4-bit) | 256-color (8-bit) | Truecolor (24-bit) |
|---|---|---|---|
| User-remappable | All 16 slots | Only indices 0-15 | None |
| Theme-adaptive | Yes | Partially | No |
| Fixed RGB output | No | Indices 16-255 | Always |
| Accessibility profile | **Best** | Mixed | Requires careful design |
| Universal terminal support | All terminals | Since ~1999 | Modern terminals (post-2018) |
| Detection signal | Default baseline | `TERM` ends with `-256color` | `COLORTERM=truecolor` or `24bit` |

**The accessibility hierarchy is clear:** ANSI 16 gives users maximum control, 8-bit gives partial control, and 24-bit gives zero user control. This is why criterion CV-8 recommends preferring 4-bit ANSI colors for informational content.

### 1.2 SGR Text Attributes (Non-Color Styling)

ANSI Select Graphic Rendition (SGR) provides non-color styling attributes defined by ECMA-48 (ISO/IEC 6429, ANSI X3.64). These attributes are independent of color and can be combined with it.

| SGR Code | Reset Code | Attribute | Support Level | Accessibility Notes |
|---|---|---|---|---|
| 0 | — | Reset all attributes | Universal | Returns to terminal defaults |
| 1 | 22 | **Bold** / increased intensity | Universal | Some terminals render as brighter color instead of heavier weight. Screen readers do NOT announce bold. |
| 2 | 22 | Dim / faint | Uneven | Alpha-blend on some terminals, intensity reduction on others. MUST NOT be used as sole differentiator — may be invisible on some terminals. |
| 3 | 23 | *Italic* | Moderate | Not supported in legacy Windows consoles, some minimal terminals. |
| 4 | 24 | Underline | Universal | Screen readers do NOT announce underline. |
| 5 | 25 | Blink (slow) | Low | Most modern terminals ignore, render as bold, or render without animation. **MUST be avoided** — can trigger photosensitive seizures. See criterion TM-3. |
| 6 | 25 | Blink (rapid) | Very low | Even less supported than slow blink. Same seizure risk. |
| 7 | 27 | Inverse / reverse video | Universal | Swaps foreground and background colors. Useful for focus indicators and selection highlighting. |
| 8 | 28 | Conceal / hidden | Low | Hides text visually but screen readers still read it. Rarely used. |
| 9 | 29 | ~~Strikethrough~~ | Moderate | Not reliable across all terminals. |

Multiple SGR codes can be combined in a single sequence using semicolons: `\x1b[1;4;31m` applies bold + underline + red simultaneously. Each attribute has a specific reset code: SGR 22 resets both bold and dim (they are mutually exclusive), SGR 24 resets underline, SGR 27 resets inverse.

**Bold and color interaction:** SGR 1 (bold) historically served double duty in many terminals — brightening the foreground color in addition to (or instead of) increasing font weight. The bright-color-on-bold behavior is why the "bright" ANSI color range (SGR 90-97) exists as explicit alternatives. Modern terminals generally support true bold weight, but the legacy behavior means bold can provide redundant visual encoding — both weight AND color shift — which is an accessibility benefit.

**Dim is unreliable:** SGR 2 (dim/faint) has the most inconsistent implementation across terminals. Some terminals reduce opacity via alpha-blending, others reduce intensity, and some ignore it entirely. Because of this, dim MUST NOT be used as the sole means of conveying information. It is acceptable as a supplementary cue alongside text-based indicators.

### 1.3 Color Capability Detection

Before emitting colored output, a CLI tool must determine the terminal's color capability tier. The recommended detection flow, in priority order:

```
1. Check app-specific flags (--no-color, --color=WHEN)  ← Highest priority
2. If NO_COLOR is set and non-empty  → MONOCHROME
3. If TERM == "dumb" or TERM is empty → MONOCHROME
4. If !isatty(stdout)                → MONOCHROME (pipe-safe)
5. If COLORTERM == "truecolor" or "24bit" → TRUECOLOR
6. If TERM matches *-direct          → TRUECOLOR
7. If TERM matches *-256color        → 256-COLOR
8. If TERM contains xterm/screen/tmux/rxvt → ANSI 16
9. Fallback                          → ANSI 16 or MONOCHROME
```

#### Environment Variables for Color Detection

| Variable | Values | Meaning | Standard |
|---|---|---|---|
| `NO_COLOR` | Any non-empty value | Disable all color output. Accessibility and user-preference signal, not a capability question. | [no-color.org](https://no-color.org/) (2017) |
| `FORCE_COLOR` | Any non-empty value | Force color even in non-TTY contexts. Overrides `NO_COLOR`. | [force-color.org](https://force-color.org/) (2023) |
| `COLORTERM` | `truecolor` or `24bit` | Terminal supports 24-bit RGB color. Set automatically by VTE-based terminals, Konsole, iTerm2. | De facto standard |
| `TERM` | `dumb` | Monochrome — no escape sequence support at all. | POSIX |
| `TERM` | `*-256color` (e.g., `xterm-256color`) | 256-color palette support. Most modern terminals default to this. | terminfo |
| `TERM` | `xterm`, `screen`, `tmux`, `rxvt`, etc. | ANSI 16-color baseline. | terminfo |
| `CLICOLOR` | `1` | Enable color output (deprecated; prefer `NO_COLOR` / `FORCE_COLOR`). | De facto (BSD origin) |
| `CLICOLOR_FORCE` | `1` | Force color output (deprecated; prefer `FORCE_COLOR`). | De facto (BSD origin) |

**Important caveats:**

- **`TERM` is a claim, not truth.** It tells you how the terminal *wants to be treated*, not what it actually is. Most modern terminals ship `xterm-256color` even when they support truecolor. `COLORTERM` is the truecolor signal; `TERM` is the baseline floor.
- **`NO_COLOR=` (empty) is NOT set.** Only a non-empty value disables color. This is specified by no-color.org and critical for correct implementation.
- **SSH propagates `TERM` but not `COLORTERM`.** Tools should handle the case where `COLORTERM` is lost during SSH sessions gracefully.
- **`CLICOLOR` / `CLICOLOR_FORCE` are deprecated.** New tools should use `NO_COLOR` and `FORCE_COLOR`. Tools that support both should let `NO_COLOR` / `FORCE_COLOR` take precedence.

#### OSC Probes for Runtime Detection

For definitive runtime color information, tools can query the terminal directly using Operating System Command (OSC) sequences. These are especially useful for detecting the terminal's background color (for contrast calculations) and confirming truecolor support.

| Query | Purpose | Syntax | Response Format |
|---|---|---|---|
| OSC 10 | Foreground color | `\x1b]10;?\a` | `\x1b]10;rgb:RRRR/GGGG/BBBB\a` |
| OSC 11 | Background color | `\x1b]11;?\a` | `\x1b]11;rgb:RRRR/GGGG/BBBB\a` |
| OSC 4 | ANSI palette slot N | `\x1b]4;N;?\a` | `\x1b]4;N;rgb:RRRR/GGGG/BBBB\a` |
| OSC 12 | Cursor color | `\x1b]12;?\a` | `\x1b]12;rgb:RRRR/GGGG/BBBB\a` |

Response values are 16-bit hex per channel (e.g., `rgb:ffff/0000/0000` for pure red). To convert to 8-bit, take the first two hex digits of each channel.

**Guidelines for OSC probing:**

- Use a short timeout (100-200ms) — terminals that don't support a query won't respond.
- Switch to raw terminal mode (`stty raw`) before probing; restore afterward.
- Fire queries in parallel and collect responses.
- Accept partial results — missing data should fall back to safe defaults.
- Never block indefinitely waiting for a response.
- OSC probing is OPTIONAL — environment variable detection is sufficient for most tools.

#### Graceful Degradation Strategy

After detecting the color tier, tools should degrade gracefully:

| Tier | Rendering Strategy |
|---|---|
| **Truecolor** | Emit full RGB values for all colors. |
| **256-color** | Quantize RGB colors to the nearest index in the 256-color cube. Indices 0-15 use ANSI named slots. |
| **ANSI 16** | Map to named color slots. Let the user's theme define what those colors look like. |
| **Monochrome** | Strip all color. Use SGR attributes (bold, inverse, underline) and text-based indicators (prefixes, symbols) to convey hierarchy and meaning. |

The recommended architecture: use semantic color tokens (like `$error`, `$warning`, `$success`, `$info`, `$muted`) that resolve to different concrete escape sequences per tier. This keeps component code tier-agnostic and ensures consistent degradation.

### 1.4 Color Vision Deficiency (Color Blindness)

Color vision deficiency (CVD) affects approximately 8% of males and 0.5% of females worldwide — roughly 1 in 12 male users of any CLI tool. Understanding CVD types is critical for choosing accessible color combinations.

#### Types of Color Vision Deficiency

| Type | Category | Prevalence | What Happens | CLI Impact |
|---|---|---|---|---|
| **Deuteranopia** | Red-green (green-blind) | ~6% of males | Green cones absent or shifted. Red and green appear as similar brownish-yellow shades. | Red/green success/failure indicators become indistinguishable. The most impactful CVD type for CLI tools due to high prevalence and the ubiquity of red/green status patterns. |
| **Protanopia** | Red-green (red-blind) | ~1% of males | Red cones absent or shifted. Red appears darker than normal — red text on dark backgrounds can become invisible. | Red error text may disappear against dark terminal backgrounds. Red and green are confused similarly to deuteranopia. |
| **Tritanopia** | Blue-yellow (blue-blind) | ~0.01% | Blue cones absent or shifted. Blue and purple appear as similar shades of blue-green. Yellow and pink become confused. | Blue informational markers and purple link/reference indicators look identical. Rare enough to be a secondary concern. |
| **Achromatopsia** | Total color blindness | ~0.003% | No functioning cones. Vision is entirely in grayscale. | All color-only information is lost. Only luminance (brightness) differences are perceived. Bold, inverse, and underline become the only available emphasis channels. |

#### Dangerous Color Pairings for CLI Tools

| Pairing | Risk | Affected Population | Alternative |
|---|---|---|---|
| **Red + Green** (e.g., pass/fail, add/delete) | Both appear as brownish-yellow | ~7% of males (deutan + protan) | Use blue + orange, or add text prefixes (`✓`/`✗`, `+`/`-`, `PASS`/`FAIL`) |
| **Green + Brown/Orange** | Nearly identical under red-green CVD | ~7% of males | Use distinct luminance levels; add text indicators |
| **Blue + Purple** | Confused by tritanopia | ~0.01% | Use blue + orange, or add text indicators |
| **Red on dark background** | Red appears very dark to protanopes — invisible on dark terminals | ~1% of males | Use bright red (SGR 91) instead of standard red (SGR 31); always pair with text prefix |
| **Any pastel-on-pastel** | Insufficient luminance difference for all users | Everyone with reduced contrast sensitivity | Ensure minimum 4.5:1 contrast ratio |

#### Safe Color Strategies

1. **Blue + Orange:** The most reliable combination — distinguishable across all common CVD types. Blue is perceived by all CVD types; orange sits in a spectrum region where red-green and blue-yellow confusion lines don't overlap.

2. **Color Brewer Safe Palette:** Originally developed for cartography, this 7-color set remains distinguishable under deuteranopia, protanopia, and tritanopia: orange, sky blue, bluish green, yellow, blue, vermilion, reddish purple.

3. **Vary luminance, not just hue:** Any color pairing that varies in lightness will remain distinguishable even in full grayscale (achromatopsia), when printed in black and white, or when viewed through any type of CVD.

4. **Always pair with text:** Use text prefixes (`[ERROR]`, `[WARN]`, `[OK]`), Unicode symbols (`✓`, `✗`, `●`, `▶`, `!`), or formatting alongside color. This eliminates dependence on color perception entirely and is the single most effective accessibility practice.

### 1.5 Contrast in Terminal Contexts

WCAG 2.2 Success Criterion 1.4.3 requires a contrast ratio of at least 4.5:1 for normal text and 3:1 for large text (Level AA). SC 1.4.6 raises the requirement to 7:1 for normal text (Level AAA).

#### The Contrast Ratio Formula

```
contrast_ratio = (L1 + 0.05) / (L2 + 0.05)
```

where L1 is the relative luminance of the lighter color and L2 is the relative luminance of the darker color. The ratio ranges from 1:1 (identical colors) to 21:1 (black on white).

**Relative luminance calculation:**

1. Convert sRGB values (0-255) to linear light:
   ```
   C_srgb = C / 255
   C_linear = C_srgb <= 0.04045
              ? C_srgb / 12.92
              : ((C_srgb + 0.055) / 1.055) ^ 2.4
   ```

2. Calculate relative luminance:
   ```
   L = 0.2126 * R_linear + 0.7152 * G_linear + 0.0722 * B_linear
   ```

   (Green is weighted most heavily because human vision is most sensitive to green wavelengths.)

#### The Terminal Background Problem

In web development, the application controls both foreground and background colors via CSS. In terminals, the application controls only the foreground — the background is set by the user's terminal emulator theme. This creates a unique challenge:

- A foreground color designed for a dark background (`#1E1E1E`) may fail contrast against a light background (`#FFFFFF`).
- A foreground color designed for a light background may vanish on a dark background.
- The application has no direct way to know the background color (though OSC 11 can probe it at runtime).

**Mitigation strategies (in order of preference):**

1. **Use ANSI 16 named colors.** These adapt to the user's theme. If the user has a well-designed theme (Solarized, Dracula, Gruvbox, Tempus), the theme author has already ensured contrast ratios are sufficient for both their light and dark variants. The application gets accessible colors for free.

2. **Don't set colors at all.** WCAG 2.2 Technique G148: "Not specifying background color, not specifying text color, and not using technology features that change those defaults." The terminal's default foreground color on the user's chosen background is almost always readable — the user chose it that way.

3. **If using fixed colors, test against both backgrounds.** Any 8-bit (index 16+) or 24-bit color must be verified to maintain 4.5:1 contrast against both `#FFFFFF` (common light background) and `#1E1E1E` (common dark background). Colors that fail against either should be avoided.

4. **Detect background via OSC 11.** If the tool must use fixed colors and needs to ensure contrast, it can query the terminal's background color at runtime and select an appropriate palette. This is the most robust approach but adds complexity.

#### Dark Mode Contrast Considerations

WCAG contrast math is symmetric — the formula produces the same ratio regardless of which side is lighter. However, human perception is not symmetric:

- A 4.5:1 ratio that reads crisply on a white background can appear slightly blurry on a dark background. This is due to the "halation" effect — light text on dark backgrounds causes slight optical blurring.
- Accessibility-conscious tools targeting dark terminals should aim for 7:1 (WCAG AAA) rather than 4.5:1 (WCAG AA) for body text.
- On a typical dark background like `#121212`, text at `#B0B0B0` produces ~8:1 contrast — readable but soft.

#### Contrast and ANSI 16 Colors

Because ANSI 16 colors are user-defined, the *application* cannot guarantee a specific contrast ratio — the user's theme determines the actual RGB values. This is by design and is an accessibility feature, not a limitation. The responsibility for contrast in ANSI 16 mode lies with the theme author and the user, not the application.

For tools that use fixed 8-bit or 24-bit colors, contrast responsibility shifts to the application, and the conformance suite can verify ratios. This is another reason to prefer ANSI 16 colors.

### 1.6 Screen Readers and Terminal Styling

Screen readers interact with terminals at the text buffer level — they read the matrix of characters in the terminal's display buffer. This has critical implications for color and styling:

**Screen readers do NOT announce ANSI styling.** Bold, underline, italic, color, inverse, dim, and strikethrough are purely visual attributes. A screen reader user hears only the text content. If a heading is made bold with no other distinguishing feature, the screen reader user experiences it as a continuous run of text. If an error is shown in red with no `Error:` prefix, the screen reader user has no way to know it's an error.

**The terminal as a "user agent" with no accessibility tree.** The W3C WCAG2ICT document explains that "unlike the semantic objects of graphical user interfaces and web pages, the output of text-based applications consists of plain text." A terminal emulator "might render some content such as escape codes as semantic elements, but otherwise exposes only lines of text to assistive technology." There is no DOM, no ARIA, no accessibility tree — just rows and columns of characters.

**Screen reader strategies for terminal text.** Terminal screen readers (NVDA's terminal support, VoiceOver Terminal, Fenrir, tdsr, Emacspeak) use heuristic analysis of the text buffer to detect structure:

- Scanning for screen updates to determine what changed
- Analyzing indentation and spacing to infer hierarchy
- Detecting inverse video to identify input fields and selections
- Tracking cursor position to follow user input

These heuristics work best with clean, predictable, linear text output. Complex formatting, rapid redraws, and color-only semantics defeat them.

**Practical implication for all criteria in this domain:** Every piece of information conveyed by color or styling MUST also be conveyed by text. This is not a general best practice — it is a hard requirement for screen reader accessibility. The text-based alternative is not a fallback for edge cases; it is the primary information channel for screen reader users.

---

## 2. Criteria

### CV-1: Color Not Sole Information Channel

| | |
|---|---|
| **Level** | A |
| **Testability** | `[SEMI]` |
| **Requirement** | Color MUST NOT be the only visual means of conveying information, indicating an action, prompting a response, or distinguishing a visual element. The tool MUST use additional indicators alongside color: whitespace, indentation, text prefixes (e.g., `[ERROR]`, `[WARN]`, `[OK]`), Unicode symbols (e.g., `✓`, `✗`, `●`), or formatting changes. This applies to all output including status indicators, diff output, tables, logs, and interactive prompts. |
| **Rationale** | Approximately 8% of males have some form of color vision deficiency. Screen readers cannot announce color changes — they read only text content. Even for sighted users, color meaning is lost when output is piped, redirected to a file, printed on paper, or viewed on a monochrome display. Red/green success/failure indicators — the single most common color-only pattern in CLIs — are indistinguishable to the ~7% of males with deuteranopia or protanopia. |
| **Test method** | Run the tool with `NO_COLOR=1`. Compare output side-by-side with default colored output. For every piece of information distinguishable by color in the default output, verify it is also distinguishable in the no-color output via text, symbols, prefixes, indentation, or whitespace. Common violations to check: diff output using only red/green with no `+`/`-` prefixes; status output using only green/red with no `[OK]`/`[FAIL]` text; table columns using only color to indicate categories; log levels using only color with no `INFO:`/`WARN:`/`ERROR:` prefixes. |
| **Positive examples** | `git diff` uses `+`/`-` prefixes AND color. `rustc` errors use `error:` prefix AND red. `pytest` uses `PASSED`/`FAILED` text AND color. `eslint` uses `✓`/`✗` symbols AND color for pass/fail. |
| **Negative examples** | A tool that shows passing tests in green and failing tests in red with no textual distinction. A tool that highlights search matches using only background color with no other indicator. A log viewer that distinguishes log levels solely by color. |
| **Disability impact** | Visual (color blindness, screen readers), Cognitive (color meaning varies by culture and context) |
| **WCAG mapping** | 1.4.1 Use of Color |

### CV-2: `NO_COLOR` Support

| | |
|---|---|
| **Level** | A |
| **Testability** | `[AUTO]` |
| **Requirement** | When the `NO_COLOR` environment variable is set and non-empty (regardless of its specific value), the tool MUST suppress all ANSI color codes from output. Non-color styling (bold via SGR 1, underline via SGR 4, inverse via SGR 7) MAY be retained, as the `NO_COLOR` standard specifies only color suppression. An empty `NO_COLOR=` value MUST be treated as unset (color allowed). |
| **Rationale** | `NO_COLOR` is the de facto universal signal for "I don't want color in my terminal." Users set it for many accessibility reasons: screen reader compatibility (raw escape sequences produce garbled speech output), color vision deficiency, photosensitivity, cognitive preference for simpler output, or piping output through tools that don't handle escape codes. It is an accessibility and user-preference signal, not a capability question — the terminal may well support color, but the user has chosen not to see it. As of 2026, the standard has 180+ adopters including Python 3.13+, npm, ripgrep, GitHub CLI, Ansible, jq, bat, fzf, and Homebrew. Applications that override or ignore `NO_COLOR` break user trust and accessibility workflows. |
| **Test method** | (1) Run the tool with `NO_COLOR=1` and capture raw output bytes. Scan for ANSI SGR color sequences using these patterns: foreground color `\x1b\[(3[0-7]|9[0-7]|38;5;\d+|38;2;\d+;\d+;\d+)m`, background color `\x1b\[(4[0-7]|10[0-7]|48;5;\d+|48;2;\d+;\d+;\d+)m`. Verify none are present. (2) Run without `NO_COLOR` on a TTY and verify color codes ARE present (confirming the tool uses color by default — a tool that never uses color trivially satisfies this criterion). (3) Run with `NO_COLOR=` (empty value) and verify color codes ARE present (empty is treated as unset per spec). |
| **Disability impact** | Visual |
| **WCAG mapping** | 1.4.1 Use of Color |
| **References** | [no-color.org specification](https://no-color.org/) |

### CV-3: `--no-color` Flag

| | |
|---|---|
| **Level** | A |
| **Testability** | `[AUTO]` |
| **Requirement** | The tool MUST provide a `--no-color` command-line flag that suppresses ANSI color codes, giving users per-invocation control independent of environment variables. The flag MUST override any application-specific color configuration. The tool SHOULD also accept `--no-colour` (British spelling) as an alias where the flag framework supports it. |
| **Rationale** | While `NO_COLOR` provides session-wide control, `--no-color` gives per-invocation control. A user who generally wants color may need to disable it for a specific command — when piping output to a screen reader script, when capturing output for a bug report, or when sharing terminal output in a text-only medium (email, chat, issue tracker). Per-invocation flags also work in contexts where environment variables are cumbersome (aliases, shell functions, one-off commands). |
| **Test method** | Run the tool with `--no-color`. Verify: exit code 0 (flag accepted), output contains no ANSI color escape sequences. Verify the flag works even when `FORCE_COLOR=1` is set (flag takes highest precedence per CV-10). |
| **Disability impact** | Visual |
| **WCAG mapping** | 1.4.1 Use of Color |

### CV-4: TTY-Aware Color

| | |
|---|---|
| **Level** | A |
| **Testability** | `[AUTO]` |
| **Requirement** | When stdout is not a TTY (detected via `isatty(1)` / `isatty(stdout)` or language equivalent), ANSI color and styling escape sequences MUST be suppressed from stdout by default. When stderr is not a TTY, the same MUST apply independently to stderr. Each stream's TTY status MUST be checked independently — a non-TTY stdout does not imply non-TTY stderr, and vice versa. |
| **Rationale** | When CLI output is piped to another program (`tool | grep pattern`), redirected to a file (`tool > output.txt`), or captured in a variable (`result=$(tool)`), ANSI escape sequences become raw bytes in the data stream. These appear as garbled characters in text files (e.g., `[31m`), corrupt downstream tool parsing (breaking `awk`, `cut`, `wc`), produce gibberish in screen reader output, and break machine-readable processing. Independent stream checking is critical because a tool may pipe stdout while displaying progress on stderr — stderr should retain color in that scenario if stderr is still a TTY. |
| **Test method** | (1) Pipe stdout through `cat` (making stdout a pipe) and scan stdout for any ANSI escape sequences matching `\x1b[\[\]()]`. Verify none are present. (2) Redirect stderr to a file (`2>/tmp/stderr.txt`) and scan the file for ANSI escape sequences. Verify none are present. (3) Verify that piping stdout does NOT strip color from stderr when stderr is still a TTY (independent stream behavior). |
| **Disability impact** | Visual |
| **WCAG mapping** | — |

### CV-5: `TERM=dumb` Respect

| | |
|---|---|
| **Level** | A |
| **Testability** | `[AUTO]` |
| **Requirement** | When the `TERM` environment variable is set to `dumb`, the tool MUST suppress ALL ANSI escape sequences — not just color, but also cursor movement (`\x1b[H`, `\x1b[A`-`\x1b[D`), screen clearing (`\x1b[2J`, `\x1b[K`), text styling (bold, underline, inverse), and any other CSI or OSC sequences. `TERM=dumb` is the POSIX signal for a terminal with zero capability for escape sequence interpretation. When `TERM` is empty or unset, the tool SHOULD behave as if `TERM=dumb`. |
| **Rationale** | `TERM=dumb` is set by terminal environments that cannot interpret escape sequences: Emacs shell buffers (a common environment for Emacspeak users — the primary screen reader for blind Emacs users), Emacs `M-x compile` output buffers, CI/CD log viewers, some IDE integrated terminals, serial console connections, and custom assistive technology terminals. Unlike `NO_COLOR` (which suppresses only color), `TERM=dumb` signals that the rendering environment cannot process ANY escape sequences. Raw escape sequences in these environments appear as visible garbage characters like `[31m`, `[0m`, `[2J`, `[?25l`, which are not just ugly but actively confuse screen readers and Braille display users who encounter them as literal text. |
| **Test method** | Run the tool with `TERM=dumb` and capture raw output bytes. Scan for ANY ANSI escape sequences using a broad regex: `\x1b[\[\]()#;][^\x1b]*` (covers CSI, OSC, and other escape types). Verify none are present. Also run with `TERM=` (empty) and verify the same behavior. |
| **Disability impact** | Visual |
| **WCAG mapping** | — |

### CV-6: `--color` Flag with Modes

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[AUTO]` |
| **Requirement** | The tool SHOULD support a `--color=WHEN` flag accepting at least three values: `always` (force color regardless of TTY status, `NO_COLOR`, or `TERM=dumb`), `never` (suppress all color — equivalent to `--no-color`), and `auto` (default behavior — enable color only when output is a TTY and no disable signals are set). The flag MAY also be spelled `--colour=WHEN`. The `--color` flag takes the highest precedence in the color configuration stack — it overrides all environment variables and config files. |
| **Rationale** | The `auto`/`always`/`never` tri-state gives users complete control over color behavior for any context. `auto` is the safe default. `always` is needed for piping colored output through tools like `less -R` that can interpret escape sequences, or for capturing colored output in tools that render ANSI (e.g., `bat`, `delta`). `never` provides explicit opt-out. This flag pattern is established by coreutils (`ls --color=auto`), `grep --color=auto`, `git config color.ui auto`, `ripgrep --color=auto`, and many other widely-used tools. |
| **Test method** | (1) `--color=never`: verify no ANSI color codes in output. (2) `--color=always` with piped output (stdout is not a TTY): verify ANSI color codes ARE present despite non-TTY. (3) `--color=auto` on a TTY: verify color present. (4) `--color=auto` piped through `cat`: verify no color. |
| **Disability impact** | Visual, Motor (reduces number of steps needed to control color per-invocation) |
| **WCAG mapping** | — |

### CV-7: `FORCE_COLOR` Support

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[AUTO]` |
| **Requirement** | When the `FORCE_COLOR` environment variable is set and non-empty, the tool SHOULD force color output even when stdout is not a TTY. `FORCE_COLOR` overrides `NO_COLOR` in the environment variable precedence chain (but is itself overridden by explicit CLI flags like `--color=never`). An empty `FORCE_COLOR=` value MUST be treated as unset. |
| **Rationale** | `FORCE_COLOR` addresses a legitimate accessibility and usability need: users who have configured their pager or log viewer to handle ANSI escape codes need color preserved through pipes (`less -R`, `bat`). CI environments where colored output aids readability in web-based log viewers but TTY detection fails. Users piping through tools that re-render ANSI sequences accessibly. The environment variable approach works in contexts where CLI flags are inaccessible — internal piping between programs, tool chains where the user cannot add flags to intermediate commands. As of 2026, adopted by Node.js, Python 3.13+, Deno, Jest, pytest, npm, Sphinx, uv, and many others. |
| **Test method** | (1) Pipe stdout through `cat` with `FORCE_COLOR=1` set. Verify ANSI color codes ARE present in piped output despite non-TTY stdout. (2) Set both `NO_COLOR=1` and `FORCE_COLOR=1`. Verify color IS present (`FORCE_COLOR` overrides `NO_COLOR`). (3) Set `FORCE_COLOR=1` and run with `--color=never`. Verify no color (CLI flag overrides `FORCE_COLOR` per CV-10). |
| **Disability impact** | Visual |
| **WCAG mapping** | — |
| **References** | [force-color.org specification](https://force-color.org/) |

### CV-8: 4-Bit ANSI Color Preference

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[SEMI]` |
| **Requirement** | For all informational and semantic color use (errors, warnings, success indicators, status markers, diff output, log levels, prompts), the tool SHOULD prefer 4-bit ANSI colors (the 16 standard colors: SGR codes 30-37, 40-47, 90-97, 100-107) over 8-bit (256-color, indices 16-255) or 24-bit (truecolor) palettes. 8-bit and 24-bit colors MAY be used for decorative or aesthetic elements (syntax highlighting gradients, chart fills, decorative borders) where the specific shade is not semantically required. When the tool uses 8-bit or 24-bit colors, it MUST provide a fallback to ANSI 16 when `COLORTERM` is absent and `TERM` does not indicate extended color support. |
| **Rationale** | ANSI 16 colors are named slots, not fixed RGB values. The terminal emulator maps each slot to a user-chosen hex color (see Section 1.1). This means users with color vision deficiency can select schemes where all 16 slots are perceptually distinguishable; users requiring high contrast can ensure sufficient ratios; users with photosensitivity can choose muted palettes. Accessibility-focused schemes like Tempus provide WCAG AA-compliant contrast (4.5:1 minimum) across all 16 colors. Fixed 8-bit (indices 16-255) and 24-bit colors bypass the user's theme entirely — a hardcoded color renders identically regardless of accessibility preferences, cannot be remapped, and may fail against the user's chosen background. GitHub's CLI team learned this lesson practically: they aligned their color palette to 4-bit ANSI colors so users could "completely customize their experience using their terminal preferences." |
| **Test method** | Capture raw output bytes from a representative set of tool commands. Extract all ANSI color SGR sequences. Classify each as: 4-bit (SGR 30-37, 40-47, 90-97, 100-107, and default reset 39/49), 8-bit (`38;5;N` or `48;5;N` where N >= 16), or 24-bit (`38;2;R;G;B` or `48;2;R;G;B`). If >50% of informational color output uses 8-bit or 24-bit, flag for human review. If any error, warning, or status indicator uses fixed 8-bit or 24-bit color, flag as a violation. |
| **Disability impact** | Visual |
| **WCAG mapping** | 1.4.3 Contrast (Minimum) |
| **References** | [Terminal Color Fundamentals (terminfo.dev)](https://terminfo.dev/fundamentals/color-fundamentals), [Tempus Themes](https://protesilaos.com/codelog/tempus-themes-intro/), [GitHub Blog: Building a more accessible GitHub CLI](https://github.blog/engineering/user-experience/building-a-more-accessible-github-cli/) |

### CV-9: No Background Color Assumption

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[SEMI]` |
| **Requirement** | The tool MUST NOT assume a specific terminal background color (light or dark). Color choices MUST be readable against both light and dark backgrounds, or the tool MUST detect the terminal's background color (e.g., via OSC 11 query) and adapt. The safest approaches, in order of preference: (1) use only the terminal's default foreground color (SGR 39) and rely on structural formatting; (2) use ANSI 16 named colors, which adapt to the user's theme; (3) if using fixed colors, detect the background via OSC 11 and select an appropriate palette. The tool MUST NOT set a background color for normal text output (SGR 40-47, 100-107, or `48;...`) except for specific UI elements like selection highlighting or progress bars. |
| **Rationale** | Unlike web browsers where the application controls both foreground and background via CSS, terminal applications control only the foreground — the background is set by the user's terminal emulator theme. This is the single most common color accessibility failure in CLI tools. GitHub's CLI team documented this exact problem: "a terminal's background color is not set by the application" — their legacy Markdown palette "did not take the terminal's background color into account, leading to low contrast in some cases." Bright yellow text designed for a dark background produces invisible output on a light theme. Dark blue text designed for a light background vanishes on a dark theme. |
| **Test method** | Semi-automated: (1) Extract all fixed foreground colors (8-bit indices 16-255 and 24-bit RGB values) from the tool's output. (2) Calculate the contrast ratio against a standard light background (#FFFFFF) and a standard dark background (#1E1E1E) using the WCAG relative luminance formula (see Section 1.5). (3) Flag any fixed color that fails 4.5:1 contrast ratio against either background. (4) ANSI 16 colors (indices 0-15 / SGR 30-37, 90-97) are exempt from this check since they are user-customizable. |
| **Disability impact** | Visual |
| **WCAG mapping** | 1.4.3 Contrast (Minimum) |
| **References** | [Understanding SC 1.4.3 Contrast (Minimum)](https://www.w3.org/WAI/WCAG21/Understanding/contrast-minimum.html), [GitHub Blog: Building a more accessible GitHub CLI](https://github.blog/engineering/user-experience/building-a-more-accessible-github-cli/) |

### CV-10: Color Configuration Precedence

| | |
|---|---|
| **Level** | AA |
| **Testability** | `[AUTO]` |
| **Requirement** | Color configuration MUST follow a documented, predictable precedence order (highest to lowest): |

```
1. CLI flag (--color=never/always/auto, --no-color)     ← highest priority
2. FORCE_COLOR environment variable
3. NO_COLOR environment variable
4. TERM=dumb (or TERM empty/unset)
5. Application-specific config file
6. TTY auto-detection (isatty)                           ← lowest priority
```

| | |
|---|---|
| **Rationale** | Predictable precedence prevents conflicts and ensures users can always override color behavior at the level they need. A user may set `NO_COLOR=1` globally in `.bashrc` for most tools but use `FORCE_COLOR=1` in a specific CI pipeline. The `--color` flag provides the ultimate override for one-off invocations. Without documented precedence, users cannot predict what will happen when multiple signals conflict, and accessibility configurations become unreliable. The precedence order reflects signal specificity: CLI flags are the most specific (this invocation), environment variables are session-level, config files are persistent, and TTY detection is automatic. |
| **Test method** | Test all precedence conflict scenarios: (1) `NO_COLOR=1` + `--color=always` → color IS present (flag wins). (2) `NO_COLOR=1` + `FORCE_COLOR=1` → color IS present (`FORCE_COLOR` wins). (3) `TERM=dumb` with `NO_COLOR` unset → no escape sequences (`TERM=dumb` wins over TTY detection). (4) `FORCE_COLOR=1` + `--color=never` → no color (flag wins). (5) Verify the tool's help text or man page documents the precedence order. |
| **Disability impact** | Visual |
| **WCAG mapping** | — |

### CV-11: High Contrast Mode Support

| | |
|---|---|
| **Level** | AAA |
| **Testability** | `[MANUAL]` |
| **Requirement** | The tool SHOULD produce readable output when the operating system's high-contrast mode is enabled (Windows High Contrast, macOS Increase Contrast, GNOME High Contrast theme). No information should be lost when terminal colors are overridden by the user, the OS, or a forced-colors stylesheet. The tool SHOULD NOT override or fight the high-contrast color scheme — if the tool uses custom 8-bit or 24-bit colors, it SHOULD detect high-contrast mode and switch to ANSI 16 colors or monochrome output (using bold, inverse, and underline for emphasis). |
| **Rationale** | High-contrast mode is a critical accessibility feature for users with low vision, cataracts, or age-related vision loss. In high-contrast mode, the OS overrides application colors with a small, high-contrast palette — typically white on black or black on white with very few accent colors. Tools that hardcode specific foreground or background colors may produce unreadable output in these modes. GitHub's Primer design system guidance: forced colors mode "should not be overridden to satisfy aesthetic desires." |
| **Test method** | Manual: (1) Enable OS high-contrast mode. On Windows: Settings > Accessibility > Contrast Themes > select a high-contrast theme. On macOS: System Settings > Accessibility > Display > Increase Contrast. On GNOME: Settings > Accessibility > High Contrast. (2) Open a terminal emulator. (3) Run the tool with several representative commands. (4) Verify all output is readable — text is not invisible, no essential information is lost. (5) Verify the tool does not produce visual artifacts from color codes conflicting with the high-contrast palette. |
| **Disability impact** | Visual |
| **WCAG mapping** | 1.4.3 Contrast (Minimum) |
| **References** | [Primer: Color considerations](https://primer.style/accessibility/design-guidance/color-considerations/) |

### CV-12: Bold/Underline as Structural Cues

| | |
|---|---|
| **Level** | AAA |
| **Testability** | `[SEMI]` |
| **Requirement** | When using SGR text styling (bold, underline, inverse, italic) to convey structure (headings, emphasis, selection states, active items), the same structure MUST also be discernible without styling — via indentation, prefixes, blank lines, whitespace patterns, or textual markers. |
| **Rationale** | SGR text attributes are purely visual — screen readers read only the text content and do not announce "this text is bold" or "this text is underlined" (see Section 1.6). The W3C WCAG2ICT background document confirms that terminal screen readers "otherwise exposes only lines of text to assistive technology." If a heading is conveyed solely by bold formatting with no other distinguishing feature, a screen reader user experiences it as a continuous run of text with no structural delineation. Additionally, SGR support varies by terminal — dim (SGR 2) has uneven support, italic (SGR 3) is ignored by some terminals, and blink (SGR 5) is widely suppressed. Structural meaning that depends solely on these attributes may silently degrade. |
| **Test method** | (1) Run the tool with `TERM=dumb` (stripping all styling). (2) Compare the dumb output against styled output. (3) For every structural element identifiable in the styled output (headings, sections, emphasized items, selected states), verify it remains identifiable in the dumb output through text-only means. Common violations: headings that are only bold with no blank-line separation or prefix; emphasized items that are only underlined with no other indicator; selected items in a list shown only in inverse video; dim text used as the sole indicator for "disabled" or "inactive" items. |
| **Disability impact** | Visual, Cognitive |
| **WCAG mapping** | 1.3.1 Info and Relationships |
| **References** | [W3C: Background on Text / Command-Line / Terminal Applications](https://github.com/w3c/wcag2ict/blob/main/background-on-text-command-line-terminal-applications-and-interfaces.md), [ECMA-48 SGR Codes](https://strasis.com/documentation/limelight-xe/reference/ecma-48-sgr-codes) |

---

## 3. References

### Normative

| Source | Description |
|---|---|
| [ECMA-48](https://ecma-international.org/publications-and-standards/standards/ecma-48/) | Control Functions for Coded Character Sets. Defines CSI, SGR, and escape sequence standards. 5th edition, 1991. |
| [ISO/IEC 8613-6](https://www.itu.int/rec/T-REC-T.416) (ITU T.416) | Character Content Architectures. Referenced by ECMA-48 SGR 38/48 for extended color specification. |
| [WCAG 2.2 SC 1.4.1](https://www.w3.org/TR/WCAG22/#use-of-color) | Use of Color: "Color is not used as the only visual means of conveying information." |
| [WCAG 2.2 SC 1.4.3](https://www.w3.org/TR/WCAG22/#contrast-minimum) | Contrast (Minimum): 4.5:1 for normal text, 3:1 for large text. |
| [WCAG 2.2 SC 1.3.1](https://www.w3.org/TR/WCAG22/#info-and-relationships) | Info and Relationships: structure conveyed through presentation must also be programmatically determinable or available in text. |
| [WCAG2ICT](https://www.w3.org/TR/wcag2ict-22/) | Guidance on Applying WCAG 2.2 to Non-Web ICT. W3C Group Note, Oct 2024. |
| [NO_COLOR specification](https://no-color.org/) | Informal standard (2017). Command-line software which adds ANSI color by default should check for `NO_COLOR`. 180+ adopters. |
| [FORCE_COLOR specification](https://force-color.org/) | Informal standard (2023). Environment variable to force ANSI color output. Adopted by Node.js, Python 3.13+, Deno, pytest. |

### Informative

| Source | Description |
|---|---|
| [Terminal Color Fundamentals (terminfo.dev)](https://terminfo.dev/fundamentals/color-fundamentals) | Comprehensive reference on ANSI 16, 256-color, and truecolor terminal rendering, including accessibility implications of each tier. |
| [Terminal Color Detection (terminfo.dev)](https://terminfo.dev/fundamentals/color-detection) | Complete color capability detection stack: environment variables, OSC probes, degradation strategies. |
| [Tempus Themes](https://protesilaos.com/codelog/tempus-themes-intro/) | WCAG AA-compliant (4.5:1 minimum contrast) terminal color schemes for the ANSI 16 palette. |
| [GitHub Blog: Building a more accessible GitHub CLI](https://github.blog/engineering/user-experience/building-a-more-accessible-github-cli/) | Practical implementation lessons: replacing 8-bit colors with 4-bit ANSI for user customizability, replacing animated spinners with static text, building on Primer's accessibility foundations. |
| [Primer: Color considerations](https://primer.style/accessibility/design-guidance/color-considerations/) | GitHub's design system guidance on color contrast, light/dark mode, high contrast mode, and functional vs. decorative color. |
| [W3C: Background on Text / Terminal Applications](https://github.com/w3c/wcag2ict/blob/main/background-on-text-command-line-terminal-applications-and-interfaces.md) | W3C working document on how terminal emulators act as "user agents" and how screen readers interact with terminal text buffers. |
| [Understanding SC 1.4.3 Contrast (Minimum)](https://www.w3.org/WAI/WCAG21/Understanding/contrast-minimum.html) | W3C detailed explanation of the contrast ratio formula, relative luminance calculation, and large text exception. |
| [Coloring for Colorblindness (David Nichols)](https://davidmathlogic.com/colorblind/) | Interactive tool for evaluating color palette accessibility across all CVD types. |
| [Color Brewer 2.0](https://colorbrewer2.org/) | Colorblind-safe palettes originally designed for cartography, widely adopted for data visualization. |
| [Julia Evans: Standards for ANSI escape codes](https://jvns.ca/blog/2025/03/07/escape-code-standards/) | Practical overview of ANSI escape code standards, terminal support quirks, and the semicolon vs. colon debate. |
| Sampath, H., Merrick, A., & Macvean, A. (2021). [Accessibility of Command Line Interfaces](https://dl.acm.org/doi/abs/10.1145/3411764.3445544). CHI '21, ACM. | Seminal academic study on CLI accessibility. Documents how screen readers struggle with ANSI-formatted CLI output and how users resort to workarounds like `--json` flags and output redirection. |
| [Seirdy: Best practices for inclusive CLIs](https://seirdy.one/posts/2022/06/10/cli-best-practices/) | Comprehensive accessibility recommendations including espeak-ng testing methodology, `NO_COLOR` guidance, and WCAG plain-text techniques. |
