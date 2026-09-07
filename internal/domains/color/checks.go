package color

import (
	"context"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/sampras343/cli-accessibility-spec/internal/engine"
	"github.com/sampras343/cli-accessibility-spec/internal/probe"
)

// CV-1: Color Not Sole Information Channel [SEMI]
type CV1Check struct{}

func (c *CV1Check) ID() string                  { return "CV-1" }
func (c *CV1Check) Name() string                { return "Color Not Sole Information Channel" }
func (c *CV1Check) Domain() string              { return "color" }
func (c *CV1Check) Level() engine.Level          { return engine.LevelA }
func (c *CV1Check) Testability() engine.Testability { return engine.Semi }
func (c *CV1Check) SpecVersion() string          { return "1.0" }

func (c *CV1Check) Precondition(p *probe.ProbeResult) bool {
	return p.HasColor
}

func (c *CV1Check) Run(ctx context.Context, binary string, p *probe.ProbeResult) *engine.Result {
	start := time.Now()
	result := &engine.Result{
		ID: "CV-1", Name: c.Name(), Domain: "color",
		Level: engine.LevelA, Testability: engine.Semi,
		SpecVersion: "1.0", NeedsReview: true,
	}

	colorOut, err := probe.Run(ctx, binary, probe.ExecOpts{Args: []string{"--help"}})
	if err != nil {
		result.Outcome = engine.OutcomeError
		result.Remarks = fmt.Sprintf("Failed to run with color: %v", err)
		result.Duration = time.Since(start)
		return result
	}

	noColorOut, err := probe.Run(ctx, binary, probe.ExecOpts{
		Args: []string{"--help"},
		Env:  map[string]string{"NO_COLOR": "1"},
	})
	if err != nil {
		result.Outcome = engine.OutcomeError
		result.Remarks = fmt.Sprintf("Failed to run with NO_COLOR: %v", err)
		result.Duration = time.Since(start)
		return result
	}

	result.Evidence = []engine.Evidence{
		{Command: colorOut.Command, Stdout: string(colorOut.Stdout), ExitCode: colorOut.ExitCode, Duration: colorOut.Duration, Note: "Default output (may contain ANSI)"},
		{Command: noColorOut.Command, Stdout: string(noColorOut.Stdout), ExitCode: noColorOut.ExitCode, Duration: noColorOut.Duration, Env: map[string]string{"NO_COLOR": "1"}, Note: "NO_COLOR=1 output"},
	}

	stripped := probe.StripANSI(colorOut.Stdout)
	similarity := computeSimilarity(string(stripped), string(noColorOut.Stdout))

	if similarity > 0.95 {
		result.Outcome = engine.Supports
		result.Remarks = fmt.Sprintf("%.0f%% content similarity — information preserved without color", similarity*100)
	} else if similarity > 0.80 {
		result.Outcome = engine.PartiallySupports
		result.Remarks = fmt.Sprintf("%.0f%% content similarity — some information may depend on color; review needed", similarity*100)
	} else {
		result.Outcome = engine.DoesNotSupport
		result.Remarks = fmt.Sprintf("%.0f%% content similarity — significant information likely conveyed by color alone", similarity*100)
	}

	result.Duration = time.Since(start)
	return result
}

// CV-8: 4-Bit ANSI Color Preference [SEMI]
type CV8Check struct{}

func (c *CV8Check) ID() string                  { return "CV-8" }
func (c *CV8Check) Name() string                { return "4-Bit ANSI Color Preference" }
func (c *CV8Check) Domain() string              { return "color" }
func (c *CV8Check) Level() engine.Level          { return engine.LevelAA }
func (c *CV8Check) Testability() engine.Testability { return engine.Semi }
func (c *CV8Check) SpecVersion() string          { return "1.0" }

func (c *CV8Check) Precondition(p *probe.ProbeResult) bool {
	return p.HasColor
}

func (c *CV8Check) Run(ctx context.Context, binary string, p *probe.ProbeResult) *engine.Result {
	start := time.Now()
	result := &engine.Result{
		ID: "CV-8", Name: c.Name(), Domain: "color",
		Level: engine.LevelAA, Testability: engine.Semi,
		SpecVersion: "1.0", NeedsReview: true,
	}

	out, err := probe.Run(ctx, binary, probe.ExecOpts{Args: []string{"--help"}})
	if err != nil {
		result.Outcome = engine.OutcomeError
		result.Remarks = fmt.Sprintf("Failed to run: %v", err)
		result.Duration = time.Since(start)
		return result
	}

	result.Evidence = []engine.Evidence{
		{Command: out.Command, Stdout: string(out.Stdout), ExitCode: out.ExitCode, Duration: out.Duration},
	}

	classification := probe.ClassifyColors(out.Stdout)
	total := classification.FourBit + classification.EightBit + classification.TwentyFourBit

	if total == 0 {
		result.Outcome = engine.NotApplicable
		result.Remarks = "No color codes detected in output"
		result.Duration = time.Since(start)
		return result
	}

	fourBitPct := float64(classification.FourBit) / float64(total) * 100
	fixedCount := classification.EightBit + classification.TwentyFourBit

	result.Remarks = fmt.Sprintf("Color usage: %d four-bit (%.0f%%), %d eight-bit, %d twenty-four-bit",
		classification.FourBit, fourBitPct, classification.EightBit, classification.TwentyFourBit)

	if fixedCount == 0 {
		result.Outcome = engine.Supports
	} else if fourBitPct >= 50 {
		result.Outcome = engine.PartiallySupports
		result.Remarks += " — majority four-bit but some fixed colors detected"
	} else {
		result.Outcome = engine.DoesNotSupport
		result.Remarks += " — majority of colors are fixed (non-customizable by user)"
	}

	result.Duration = time.Since(start)
	return result
}

// CV-9: No Background Color Assumption [SEMI]
type CV9Check struct{}

func (c *CV9Check) ID() string                  { return "CV-9" }
func (c *CV9Check) Name() string                { return "No Background Color Assumption" }
func (c *CV9Check) Domain() string              { return "color" }
func (c *CV9Check) Level() engine.Level          { return engine.LevelAA }
func (c *CV9Check) Testability() engine.Testability { return engine.Semi }
func (c *CV9Check) SpecVersion() string          { return "1.0" }

func (c *CV9Check) Precondition(p *probe.ProbeResult) bool {
	return p.HasColor
}

func (c *CV9Check) Run(ctx context.Context, binary string, p *probe.ProbeResult) *engine.Result {
	start := time.Now()
	result := &engine.Result{
		ID: "CV-9", Name: c.Name(), Domain: "color",
		Level: engine.LevelAA, Testability: engine.Semi,
		SpecVersion: "1.0", NeedsReview: true,
	}

	out, err := probe.Run(ctx, binary, probe.ExecOpts{Args: []string{"--help"}})
	if err != nil {
		result.Outcome = engine.OutcomeError
		result.Remarks = fmt.Sprintf("Failed to run: %v", err)
		result.Duration = time.Since(start)
		return result
	}

	result.Evidence = []engine.Evidence{
		{Command: out.Command, Stdout: string(out.Stdout), ExitCode: out.ExitCode, Duration: out.Duration},
	}

	// Check for background color codes (SGR 40-47, 100-107, 48;5;N, 48;2;R;G;B)
	bgPattern := regexp.MustCompile(`\x1b\[(4[0-7]|10[0-7]|48;5;\d+|48;2;\d+;\d+;\d+)m`)
	bgMatches := bgPattern.FindAll(out.Stdout, -1)

	// Check for fixed foreground colors (8-bit index 16+ and 24-bit)
	fg8bit := regexp.MustCompile(`\x1b\[38;5;(\d+)m`)
	fg24bit := regexp.MustCompile(`\x1b\[38;2;(\d+);(\d+);(\d+)m`)

	fg8matches := fg8bit.FindAllSubmatch(out.Stdout, -1)
	fg24matches := fg24bit.FindAllSubmatch(out.Stdout, -1)

	var failingColors []string

	// Check 8-bit colors (only indices 16+ are fixed)
	for _, match := range fg8matches {
		idx, _ := strconv.Atoi(string(match[1]))
		if idx >= 16 {
			failingColors = append(failingColors, fmt.Sprintf("8-bit index %d", idx))
		}
	}

	// Check 24-bit colors against light and dark backgrounds
	lightBg := [3]float64{1.0, 1.0, 1.0}   // #FFFFFF
	darkBg := [3]float64{0.118, 0.118, 0.118} // #1E1E1E
	for _, match := range fg24matches {
		r, _ := strconv.Atoi(string(match[1]))
		g, _ := strconv.Atoi(string(match[2]))
		b, _ := strconv.Atoi(string(match[3]))
		fgLum := relativeLuminance(float64(r)/255.0, float64(g)/255.0, float64(b)/255.0)
		lightContrast := contrastRatio(fgLum, relativeLuminance(lightBg[0], lightBg[1], lightBg[2]))
		darkContrast := contrastRatio(fgLum, relativeLuminance(darkBg[0], darkBg[1], darkBg[2]))
		if lightContrast < 4.5 || darkContrast < 4.5 {
			failingColors = append(failingColors, fmt.Sprintf("rgb(%d,%d,%d) light=%.1f:1 dark=%.1f:1", r, g, b, lightContrast, darkContrast))
		}
	}

	if len(bgMatches) > 0 {
		result.Remarks = fmt.Sprintf("Tool sets background colors (%d instances) — may conflict with user's terminal theme", len(bgMatches))
		result.Outcome = engine.PartiallySupports
	} else if len(failingColors) > 0 {
		result.Remarks = fmt.Sprintf("Fixed colors with potential contrast issues: %v", failingColors)
		result.Outcome = engine.PartiallySupports
	} else if len(fg8matches) == 0 && len(fg24matches) == 0 {
		result.Outcome = engine.Supports
		result.Remarks = "Uses only ANSI 16 named colors or no fixed colors — adapts to any terminal theme"
	} else {
		result.Outcome = engine.Supports
		result.Remarks = "Fixed colors pass contrast checks against both light and dark backgrounds"
	}

	result.Duration = time.Since(start)
	return result
}

// CV-12: Bold/Underline as Structural Cues [SEMI]
type CV12Check struct{}

func (c *CV12Check) ID() string                  { return "CV-12" }
func (c *CV12Check) Name() string                { return "Bold/Underline as Structural Cues" }
func (c *CV12Check) Domain() string              { return "color" }
func (c *CV12Check) Level() engine.Level          { return engine.LevelAAA }
func (c *CV12Check) Testability() engine.Testability { return engine.Semi }
func (c *CV12Check) SpecVersion() string          { return "1.0" }

func (c *CV12Check) Precondition(p *probe.ProbeResult) bool {
	return p.HasHelp
}

func (c *CV12Check) Run(ctx context.Context, binary string, p *probe.ProbeResult) *engine.Result {
	start := time.Now()
	result := &engine.Result{
		ID: "CV-12", Name: c.Name(), Domain: "color",
		Level: engine.LevelAAA, Testability: engine.Semi,
		SpecVersion: "1.0", NeedsReview: true,
	}

	styledOut, err := probe.Run(ctx, binary, probe.ExecOpts{Args: []string{"--help"}})
	if err != nil {
		result.Outcome = engine.OutcomeError
		result.Remarks = fmt.Sprintf("Failed to run styled: %v", err)
		result.Duration = time.Since(start)
		return result
	}

	dumbOut, err := probe.Run(ctx, binary, probe.ExecOpts{
		Args: []string{"--help"},
		Env:  map[string]string{"TERM": "dumb"},
	})
	if err != nil {
		result.Outcome = engine.OutcomeError
		result.Remarks = fmt.Sprintf("Failed to run with TERM=dumb: %v", err)
		result.Duration = time.Since(start)
		return result
	}

	result.Evidence = []engine.Evidence{
		{Command: styledOut.Command, Stdout: string(styledOut.Stdout), ExitCode: styledOut.ExitCode, Duration: styledOut.Duration, Note: "Default styled output"},
		{Command: dumbOut.Command, Stdout: string(dumbOut.Stdout), ExitCode: dumbOut.ExitCode, Duration: dumbOut.Duration, Env: map[string]string{"TERM": "dumb"}, Note: "TERM=dumb output"},
	}

	stripped := probe.StripANSI(styledOut.Stdout)
	similarity := computeSimilarity(string(stripped), string(dumbOut.Stdout))

	if similarity > 0.90 {
		result.Outcome = engine.Supports
		result.Remarks = fmt.Sprintf("%.0f%% structural similarity — content structure preserved without styling", similarity*100)
	} else if similarity > 0.70 {
		result.Outcome = engine.PartiallySupports
		result.Remarks = fmt.Sprintf("%.0f%% structural similarity — some structure may depend on styling; review needed", similarity*100)
	} else {
		result.Outcome = engine.DoesNotSupport
		result.Remarks = fmt.Sprintf("%.0f%% structural similarity — significant structure lost without styling", similarity*100)
	}

	result.Duration = time.Since(start)
	return result
}

// CV-14: Foreground-Background Pair Contrast [SEMI]
type CV14Check struct{}

func (c *CV14Check) ID() string                     { return "CV-14" }
func (c *CV14Check) Name() string                   { return "Foreground-Background Pair Contrast" }
func (c *CV14Check) Domain() string                 { return "color" }
func (c *CV14Check) Level() engine.Level             { return engine.LevelAA }
func (c *CV14Check) Testability() engine.Testability { return engine.Semi }
func (c *CV14Check) SpecVersion() string             { return "1.0" }

func (c *CV14Check) Precondition(p *probe.ProbeResult) bool {
	return p.HasColor
}

// Standard xterm default RGB values for the 16 ANSI colors.
var ansi16RGB = [16][3]float64{
	{0, 0, 0},          // 0  black
	{0.502, 0, 0},      // 1  red
	{0, 0.502, 0},      // 2  green
	{0.502, 0.502, 0},  // 3  yellow
	{0, 0, 0.502},      // 4  blue
	{0.502, 0, 0.502},  // 5  magenta
	{0, 0.502, 0.502},  // 6  cyan
	{0.753, 0.753, 0.753}, // 7  white
	{0.502, 0.502, 0.502}, // 8  bright black
	{1, 0, 0},          // 9  bright red
	{0, 1, 0},          // 10 bright green
	{1, 1, 0},          // 11 bright yellow
	{0, 0, 1},          // 12 bright blue
	{1, 0, 1},          // 13 bright magenta
	{0, 1, 1},          // 14 bright cyan
	{1, 1, 1},          // 15 bright white
}

func (c *CV14Check) Run(ctx context.Context, binary string, p *probe.ProbeResult) *engine.Result {
	start := time.Now()
	result := &engine.Result{
		ID: "CV-14", Name: c.Name(), Domain: "color",
		Level: engine.LevelAA, Testability: engine.Semi,
		SpecVersion: "1.0", NeedsReview: true,
	}

	out, err := probe.Run(ctx, binary, probe.ExecOpts{
		Args: []string{"--help"},
		Env:  map[string]string{"CLICOLOR_FORCE": "1", "FORCE_COLOR": "1"},
	})
	if err != nil {
		result.Outcome = engine.OutcomeError
		result.Remarks = fmt.Sprintf("Failed to run: %v", err)
		result.Duration = time.Since(start)
		return result
	}

	result.Evidence = []engine.Evidence{
		{Command: out.Command, Stdout: string(out.Stdout), ExitCode: out.ExitCode, Duration: out.Duration},
	}

	// Find compound SGR sequences and track active fg/bg state
	sgrPattern := regexp.MustCompile(`\x1b\[([\d;]*)m`)
	matches := sgrPattern.FindAllSubmatch(out.Stdout, -1)

	var failingPairs []string
	activeFG := -1  // -1 = default/unset
	activeBG := -1

	for _, match := range matches {
		params := string(match[1])
		codes := splitSGRCodes(params)

		for _, code := range codes {
			switch {
			case code == 0:
				activeFG = -1
				activeBG = -1
			case code >= 30 && code <= 37:
				activeFG = code - 30
			case code >= 40 && code <= 47:
				activeBG = code - 40
			case code >= 90 && code <= 97:
				activeFG = code - 90 + 8
			case code >= 100 && code <= 107:
				activeBG = code - 100 + 8
			case code == 39:
				activeFG = -1
			case code == 49:
				activeBG = -1
			}
		}

		// If both fg and bg are set to ANSI 16 colors, check contrast
		if activeFG >= 0 && activeFG < 16 && activeBG >= 0 && activeBG < 16 {
			fgRGB := ansi16RGB[activeFG]
			bgRGB := ansi16RGB[activeBG]
			fgLum := relativeLuminance(fgRGB[0], fgRGB[1], fgRGB[2])
			bgLum := relativeLuminance(bgRGB[0], bgRGB[1], bgRGB[2])
			cr := contrastRatio(fgLum, bgLum)
			if cr < 4.5 {
				pair := fmt.Sprintf("fg=%d bg=%d contrast=%.1f:1", activeFG, activeBG, cr)
				// Deduplicate
				found := false
				for _, p := range failingPairs {
					if p == pair {
						found = true
						break
					}
				}
				if !found {
					failingPairs = append(failingPairs, pair)
				}
			}
		}
	}

	if len(matches) == 0 {
		result.Outcome = engine.NotApplicable
		result.Remarks = "No SGR sequences found in output"
	} else if len(failingPairs) > 0 {
		result.Outcome = engine.DoesNotSupport
		result.Remarks = fmt.Sprintf("Foreground-background pairs with contrast < 4.5:1: %v", failingPairs)
	} else {
		result.Outcome = engine.Supports
		result.Remarks = "All foreground-background color pairs meet 4.5:1 contrast ratio, or tool does not set both fg and bg"
	}

	result.Duration = time.Since(start)
	return result
}

func splitSGRCodes(params string) []int {
	if params == "" {
		return []int{0}
	}
	parts := regexp.MustCompile(`;`).Split(params, -1)
	codes := make([]int, 0, len(parts))
	for _, p := range parts {
		if p == "" {
			codes = append(codes, 0)
			continue
		}
		n, err := strconv.Atoi(p)
		if err == nil {
			codes = append(codes, n)
		}
	}
	return codes
}

// CV-15: Unicode/Emoji Symbol Accessibility [SEMI]
type CV15Check struct{}

func (c *CV15Check) ID() string                     { return "CV-15" }
func (c *CV15Check) Name() string                   { return "Unicode/Emoji Symbol Accessibility" }
func (c *CV15Check) Domain() string                 { return "color" }
func (c *CV15Check) Level() engine.Level             { return engine.LevelA }
func (c *CV15Check) Testability() engine.Testability { return engine.Semi }
func (c *CV15Check) SpecVersion() string             { return "1.0" }

func (c *CV15Check) Precondition(p *probe.ProbeResult) bool {
	return p.HasHelp
}

var indicatorSymbols = []rune{
	'✓', '✗', '✔', '✘', '☑', '☒', '⊘',
	'●', '○', '◉', '◆', '◇', '⬤',
	'▶', '▷', '►', '▸', '→', '←', '↑', '↓',
	'⚠', '⚡', 'ℹ', 'ⓘ',
}

var textAlternatives = []string{
	"OK", "PASS", "FAIL", "ERROR", "WARN", "INFO", "SUCCESS",
	"YES", "NO", "DONE", "SKIP", "NOTE",
	"[ok]", "[pass]", "[fail]", "[error]", "[warn]",
	"ok", "pass", "fail", "error", "warn", "info",
}

func (c *CV15Check) Run(ctx context.Context, binary string, p *probe.ProbeResult) *engine.Result {
	start := time.Now()
	result := &engine.Result{
		ID: "CV-15", Name: c.Name(), Domain: "color",
		Level: engine.LevelA, Testability: engine.Semi,
		SpecVersion: "1.0", NeedsReview: true,
	}

	out, err := probe.Run(ctx, binary, probe.ExecOpts{Args: []string{"--help"}})
	if err != nil {
		result.Outcome = engine.OutcomeError
		result.Remarks = fmt.Sprintf("Failed to run: %v", err)
		result.Duration = time.Since(start)
		return result
	}

	result.Evidence = []engine.Evidence{
		{Command: out.Command, Stdout: string(out.Stdout), ExitCode: out.ExitCode, Duration: out.Duration},
	}

	content := string(probe.StripANSI(out.Stdout))
	runes := []rune(content)

	var bareSymbols []string
	var pairedSymbols []string
	symbolSet := make(map[rune]bool)
	for _, s := range indicatorSymbols {
		symbolSet[s] = true
	}

	for i, r := range runes {
		if !symbolSet[r] && !isEmoji(r) {
			continue
		}
		// Check surrounding context (20 chars before and after) for text alternatives
		contextStart := i - 20
		if contextStart < 0 {
			contextStart = 0
		}
		contextEnd := i + 20
		if contextEnd > len(runes) {
			contextEnd = len(runes)
		}
		surrounding := strings.ToLower(string(runes[contextStart:contextEnd]))

		hasText := false
		for _, alt := range textAlternatives {
			if strings.Contains(surrounding, strings.ToLower(alt)) {
				hasText = true
				break
			}
		}

		sym := fmt.Sprintf("%c (U+%04X)", r, r)
		if hasText {
			pairedSymbols = append(pairedSymbols, sym)
		} else {
			bareSymbols = append(bareSymbols, sym)
		}
	}

	totalSymbols := len(bareSymbols) + len(pairedSymbols)

	if totalSymbols == 0 {
		result.Outcome = engine.Supports
		result.Remarks = "No Unicode indicator symbols or emoji found in output"
	} else if len(bareSymbols) == 0 {
		result.Outcome = engine.Supports
		result.Remarks = fmt.Sprintf("%d symbols found, all have adjacent text alternatives: %v", totalSymbols, pairedSymbols)
	} else if len(pairedSymbols) > 0 {
		result.Outcome = engine.PartiallySupports
		result.Remarks = fmt.Sprintf("%d symbols without text context: %v; %d with text: %v", len(bareSymbols), bareSymbols, len(pairedSymbols), pairedSymbols)
	} else {
		result.Outcome = engine.DoesNotSupport
		result.Remarks = fmt.Sprintf("%d symbols used without text alternatives: %v", len(bareSymbols), bareSymbols)
	}

	result.Duration = time.Since(start)
	return result
}

func isEmoji(r rune) bool {
	return (r >= 0x1F300 && r <= 0x1F9FF) ||
		(r >= 0x2600 && r <= 0x26FF) ||
		(r >= 0x2700 && r <= 0x27BF) ||
		(r >= 0x1FA00 && r <= 0x1FA6F) ||
		(r >= 0x1FA70 && r <= 0x1FAFF)
}

// computeSimilarity returns a 0.0-1.0 score of how similar two strings are.
func computeSimilarity(a, b string) float64 {
	if a == b {
		return 1.0
	}
	if len(a) == 0 && len(b) == 0 {
		return 1.0
	}
	if len(a) == 0 || len(b) == 0 {
		return 0.0
	}
	longer := len(a)
	if len(b) > longer {
		longer = len(b)
	}
	common := 0
	bBytes := []byte(b)
	for i, c := range []byte(a) {
		if i < len(bBytes) && c == bBytes[i] {
			common++
		}
	}
	return float64(common) / float64(longer)
}

// relativeLuminance calculates WCAG relative luminance from linear sRGB values (0-1).
func relativeLuminance(r, g, b float64) float64 {
	r = linearize(r)
	g = linearize(g)
	b = linearize(b)
	return 0.2126*r + 0.7152*g + 0.0722*b
}

func linearize(c float64) float64 {
	if c <= 0.04045 {
		return c / 12.92
	}
	return math.Pow((c+0.055)/1.055, 2.4)
}

func contrastRatio(l1, l2 float64) float64 {
	if l1 < l2 {
		l1, l2 = l2, l1
	}
	return (l1 + 0.05) / (l2 + 0.05)
}
