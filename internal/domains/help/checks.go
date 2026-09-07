package help

import (
	"context"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/sampras343/cli-accessibility-spec/internal/engine"
	"github.com/sampras343/cli-accessibility-spec/internal/probe"
)

// ---------------------------------------------------------------------------
// HD-1: --help and -h Flags [A, AUTO]
// Tries --help, -h, and bare "help" subcommand. Passes if any form works.
// ---------------------------------------------------------------------------

type HD1Check struct{}

func (c *HD1Check) ID() string                         { return "HD-1" }
func (c *HD1Check) Name() string                       { return "--help and -h Flags" }
func (c *HD1Check) Domain() string                     { return "help" }
func (c *HD1Check) Level() engine.Level                { return engine.LevelA }
func (c *HD1Check) Testability() engine.Testability    { return engine.Auto }
func (c *HD1Check) SpecVersion() string                { return "1.0" }
func (c *HD1Check) Precondition(_ *probe.ProbeResult) bool { return true }

func (c *HD1Check) Run(ctx context.Context, binary string, p *probe.ProbeResult) *engine.Result {
	start := time.Now()
	result := &engine.Result{
		ID: "HD-1", Name: c.Name(), Domain: "help",
		Level: engine.LevelA, Testability: engine.Auto, SpecVersion: "1.0",
	}

	forms := []struct {
		name string
		args []string
	}{
		{"--help", []string{"--help"}},
		{"-h", []string{"-h"}},
		{"help (bare subcommand)", []string{"help"}},
	}

	var passedForms []string
	var evidence []engine.Evidence

	for _, form := range forms {
		out, err := probe.Run(ctx, binary, probe.ExecOpts{
			Args:    form.args,
			Timeout: 5 * time.Second,
		})
		if err != nil {
			continue
		}
		ev := engine.Evidence{
			Command:  out.Command,
			Stdout:   truncate(string(out.Stdout), 500),
			Stderr:   truncate(string(out.Stderr), 500),
			ExitCode: out.ExitCode,
			Duration: out.Duration,
		}
		if out.ExitCode == 0 && len(out.Stdout) > 0 {
			ev.Note = fmt.Sprintf("%s: exit 0, %d bytes on stdout — PASS", form.name, len(out.Stdout))
			passedForms = append(passedForms, form.name)
		} else {
			ev.Note = fmt.Sprintf("%s: exit %d, %d bytes stdout — does not satisfy criterion", form.name, out.ExitCode, len(out.Stdout))
		}
		evidence = append(evidence, ev)
	}

	result.Evidence = evidence
	result.Duration = time.Since(start)

	if len(passedForms) >= 2 {
		result.Outcome = engine.Supports
		result.Remarks = fmt.Sprintf("Help available via: %s", strings.Join(passedForms, ", "))
	} else if len(passedForms) == 1 {
		result.Outcome = engine.PartiallySupports
		result.Remarks = fmt.Sprintf("Help available only via %s — tools should support at least --help and one alternative (e.g., -h or bare 'help' subcommand)", passedForms[0])
	} else {
		result.Outcome = engine.DoesNotSupport
		result.Remarks = "No help form works (tried --help, -h, and bare 'help' subcommand)"
	}
	return result
}

// ---------------------------------------------------------------------------
// HD-3: --version Flag [A, AUTO]
// Tries --version, -V, and bare "version" subcommand. Passes if any works.
// ---------------------------------------------------------------------------

type HD3Check struct{}

func (c *HD3Check) ID() string                         { return "HD-3" }
func (c *HD3Check) Name() string                       { return "--version Flag" }
func (c *HD3Check) Domain() string                     { return "help" }
func (c *HD3Check) Level() engine.Level                { return engine.LevelA }
func (c *HD3Check) Testability() engine.Testability    { return engine.Auto }
func (c *HD3Check) SpecVersion() string                { return "1.0" }
func (c *HD3Check) Precondition(_ *probe.ProbeResult) bool { return true }

func (c *HD3Check) Run(ctx context.Context, binary string, p *probe.ProbeResult) *engine.Result {
	start := time.Now()
	result := &engine.Result{
		ID: "HD-3", Name: c.Name(), Domain: "help",
		Level: engine.LevelA, Testability: engine.Auto, SpecVersion: "1.0",
	}

	forms := []struct {
		name string
		args []string
	}{
		{"--version", []string{"--version"}},
		{"-V", []string{"-V"}},
		{"version (bare subcommand)", []string{"version"}},
	}

	var passedForms []string
	var evidence []engine.Evidence

	for _, form := range forms {
		out, err := probe.Run(ctx, binary, probe.ExecOpts{
			Args:    form.args,
			Timeout: 5 * time.Second,
		})
		if err != nil {
			continue
		}
		ev := engine.Evidence{
			Command:  out.Command,
			Stdout:   truncate(string(out.Stdout), 500),
			Stderr:   truncate(string(out.Stderr), 500),
			ExitCode: out.ExitCode,
			Duration: out.Duration,
		}
		if out.ExitCode == 0 && len(out.Stdout) > 0 {
			ev.Note = fmt.Sprintf("%s: exit 0, output present — PASS", form.name)
			passedForms = append(passedForms, form.name)
		} else {
			ev.Note = fmt.Sprintf("%s: exit %d — does not satisfy criterion", form.name, out.ExitCode)
		}
		evidence = append(evidence, ev)
	}

	result.Evidence = evidence
	result.Duration = time.Since(start)

	if len(passedForms) >= 1 {
		result.Outcome = engine.Supports
		result.Remarks = fmt.Sprintf("Version info available via: %s", strings.Join(passedForms, ", "))
	} else {
		result.Outcome = engine.DoesNotSupport
		result.Remarks = "No version form works (tried --version, -V, and bare 'version' subcommand)"
	}
	return result
}

// ---------------------------------------------------------------------------
// HD-2: Subcommand Help [A, AUTO]
// ---------------------------------------------------------------------------

type HD2Check struct{}

func (c *HD2Check) ID() string                         { return "HD-2" }
func (c *HD2Check) Name() string                       { return "Subcommand Help" }
func (c *HD2Check) Domain() string                     { return "help" }
func (c *HD2Check) Level() engine.Level                { return engine.LevelA }
func (c *HD2Check) Testability() engine.Testability    { return engine.Auto }
func (c *HD2Check) SpecVersion() string                { return "1.0" }
func (c *HD2Check) Precondition(p *probe.ProbeResult) bool { return p.HasSubcommands }

func (c *HD2Check) Run(ctx context.Context, binary string, p *probe.ProbeResult) *engine.Result {
	start := time.Now()
	result := &engine.Result{
		ID: "HD-2", Name: c.Name(), Domain: "help",
		Level: engine.LevelA, Testability: engine.Auto,
		SpecVersion: "1.0",
	}

	if len(p.Subcommands) == 0 {
		result.Outcome = engine.NotApplicable
		result.Remarks = "No subcommands discovered"
		result.Duration = time.Since(start)
		return result
	}

	var failing []string
	for _, sub := range p.Subcommands {
		out, err := probe.Run(ctx, binary, probe.ExecOpts{
			Args:    []string{sub.Name, "--help"},
			Timeout: 5 * time.Second,
		})
		if err != nil {
			failing = append(failing, fmt.Sprintf("%s (exec error: %v)", sub.Name, err))
			continue
		}

		result.Evidence = append(result.Evidence, engine.Evidence{
			Command:  out.Command,
			Stdout:   truncate(string(out.Stdout), 500),
			Stderr:   truncate(string(out.Stderr), 200),
			ExitCode: out.ExitCode,
			Duration: out.Duration,
			Note:     fmt.Sprintf("Subcommand %q help", sub.Name),
		})

		if out.ExitCode != 0 || (len(out.Stdout) == 0 && len(out.Stderr) == 0) {
			failing = append(failing, fmt.Sprintf("%s (exit=%d, stdout=%d bytes)", sub.Name, out.ExitCode, len(out.Stdout)))
		}
	}

	if len(failing) == 0 {
		result.Outcome = engine.Supports
		result.Remarks = fmt.Sprintf("All %d subcommands provide --help with exit 0 and output", len(p.Subcommands))
	} else if len(failing) < len(p.Subcommands) {
		result.Outcome = engine.PartiallySupports
		result.Remarks = fmt.Sprintf("%d/%d subcommands failed: %s", len(failing), len(p.Subcommands), strings.Join(failing, "; "))
	} else {
		result.Outcome = engine.DoesNotSupport
		result.Remarks = fmt.Sprintf("All subcommands failed --help: %s", strings.Join(failing, "; "))
	}

	result.Duration = time.Since(start)
	return result
}

// ---------------------------------------------------------------------------
// HD-5: Help Text Structure [AA, SEMI]
// ---------------------------------------------------------------------------

type HD5Check struct{}

func (c *HD5Check) ID() string                         { return "HD-5" }
func (c *HD5Check) Name() string                       { return "Help Text Structure" }
func (c *HD5Check) Domain() string                     { return "help" }
func (c *HD5Check) Level() engine.Level                { return engine.LevelAA }
func (c *HD5Check) Testability() engine.Testability    { return engine.Semi }
func (c *HD5Check) SpecVersion() string                { return "1.0" }
func (c *HD5Check) Precondition(p *probe.ProbeResult) bool { return p.HasHelp }

func (c *HD5Check) Run(ctx context.Context, binary string, p *probe.ProbeResult) *engine.Result {
	start := time.Now()
	result := &engine.Result{
		ID: "HD-5", Name: c.Name(), Domain: "help",
		Level: engine.LevelAA, Testability: engine.Semi,
		SpecVersion: "1.0", NeedsReview: true,
	}

	helpText := p.HelpText
	if helpText == "" {
		result.Outcome = engine.OutcomeError
		result.Remarks = "No help text available from probe"
		result.Duration = time.Since(start)
		return result
	}

	result.Evidence = []engine.Evidence{
		{Command: binary + " --help", Stdout: truncate(helpText, 1000), Note: "Help text from probe"},
	}

	lines := strings.Split(helpText, "\n")

	// Detect sections
	sections := detectHelpSections(lines)

	// Check for blank-line separation
	hasBlankSep := false
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "" && i > 0 && i < len(lines)-1 {
			hasBlankSep = true
			break
		}
	}

	sectionCount := len(sections)

	if sectionCount >= 3 && hasBlankSep {
		result.Outcome = engine.Supports
		result.Remarks = fmt.Sprintf("Found %d sections (%s) with blank-line separation", sectionCount, strings.Join(sections, ", "))
	} else if sectionCount >= 2 {
		result.Outcome = engine.PartiallySupports
		result.Remarks = fmt.Sprintf("Found %d sections (%s); recommend at least 3 (description, usage, flags) with blank-line separation", sectionCount, strings.Join(sections, ", "))
	} else {
		result.Outcome = engine.DoesNotSupport
		result.Remarks = fmt.Sprintf("Found only %d recognizable section(s); structured help should have description, usage, and flag sections", sectionCount)
	}

	result.Duration = time.Since(start)
	return result
}

// detectHelpSections identifies structural sections in help text.
func detectHelpSections(lines []string) []string {
	var found []string
	seen := make(map[string]bool)

	hasDescription := false
	hasUsage := false
	hasFlags := false
	hasCommands := false
	hasExamples := false

	usageRe := regexp.MustCompile(`(?i)^(usage|synopsis)\s*:`)
	flagsRe := regexp.MustCompile(`(?i)^(options|flags|global flags|global options)\s*:?\s*$`)
	commandsRe := regexp.MustCompile(`(?i)^(commands|available commands|subcommands)\s*:?\s*$`)
	examplesRe := regexp.MustCompile(`(?i)^(examples?)\s*:?\s*$`)
	flagLineRe := regexp.MustCompile(`^\s+--?\w`)

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		if usageRe.MatchString(trimmed) && !hasUsage {
			hasUsage = true
		}
		if flagsRe.MatchString(trimmed) && !hasFlags {
			hasFlags = true
		}
		if commandsRe.MatchString(trimmed) && !hasCommands {
			hasCommands = true
		}
		if examplesRe.MatchString(trimmed) && !hasExamples {
			hasExamples = true
		}

		// First non-empty line that isn't a section header is likely the description
		if i < 3 && !hasDescription && !usageRe.MatchString(trimmed) && !flagsRe.MatchString(trimmed) {
			hasDescription = true
		}

		// Detect flags section by presence of flag-like lines
		if !hasFlags && flagLineRe.MatchString(line) {
			hasFlags = true
		}
	}

	if hasDescription && !seen["description"] {
		found = append(found, "description")
		seen["description"] = true
	}
	if hasUsage && !seen["usage"] {
		found = append(found, "usage")
		seen["usage"] = true
	}
	if hasFlags && !seen["flags"] {
		found = append(found, "flags")
		seen["flags"] = true
	}
	if hasCommands && !seen["commands"] {
		found = append(found, "commands")
		seen["commands"] = true
	}
	if hasExamples && !seen["examples"] {
		found = append(found, "examples")
		seen["examples"] = true
	}

	return found
}

// ---------------------------------------------------------------------------
// HD-6: Examples in Help [AA, SEMI]
// ---------------------------------------------------------------------------

type HD6Check struct{}

func (c *HD6Check) ID() string                         { return "HD-6" }
func (c *HD6Check) Name() string                       { return "Examples in Help" }
func (c *HD6Check) Domain() string                     { return "help" }
func (c *HD6Check) Level() engine.Level                { return engine.LevelAA }
func (c *HD6Check) Testability() engine.Testability    { return engine.Semi }
func (c *HD6Check) SpecVersion() string                { return "1.0" }
func (c *HD6Check) Precondition(p *probe.ProbeResult) bool { return p.HasHelp }

func (c *HD6Check) Run(ctx context.Context, binary string, p *probe.ProbeResult) *engine.Result {
	start := time.Now()
	result := &engine.Result{
		ID: "HD-6", Name: c.Name(), Domain: "help",
		Level: engine.LevelAA, Testability: engine.Semi,
		SpecVersion: "1.0", NeedsReview: true,
	}

	helpText := p.HelpText
	if helpText == "" {
		result.Outcome = engine.OutcomeError
		result.Remarks = "No help text available from probe"
		result.Duration = time.Since(start)
		return result
	}

	result.Evidence = []engine.Evidence{
		{Command: binary + " --help", Stdout: truncate(helpText, 1000), Note: "Help text from probe"},
	}

	binaryName := filepath.Base(binary)
	lines := strings.Split(helpText, "\n")

	exampleHeaderRe := regexp.MustCompile(`(?i)^(examples?)\s*:?\s*$`)
	promptRe := regexp.MustCompile(`^\s*[$>#]`)

	hasExampleHeader := false
	exampleLineCount := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if exampleHeaderRe.MatchString(trimmed) {
			hasExampleHeader = true
		}

		// Lines starting with shell prompts
		if promptRe.MatchString(line) {
			exampleLineCount++
		}

		// Lines containing the binary name with arguments (likely usage examples)
		// Exclude: usage synopsis, cross-references, help topic listings,
		// and lines with placeholder patterns like [command] or <arg>.
		if strings.Contains(trimmed, binaryName+" ") &&
			!strings.HasPrefix(strings.ToLower(trimmed), "usage") &&
			!strings.Contains(strings.ToLower(trimmed), "use \"") &&
			!strings.Contains(strings.ToLower(trimmed), "use '") &&
			!strings.Contains(strings.ToLower(trimmed), "see ") &&
			!strings.Contains(trimmed, "[command]") &&
			!strings.Contains(trimmed, "<command>") &&
			!strings.Contains(trimmed, "not built with") &&
			(strings.HasPrefix(trimmed, "$") || strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, ">")) {
			exampleLineCount++
		}
	}

	if hasExampleHeader && exampleLineCount > 0 {
		result.Outcome = engine.Supports
		result.Remarks = fmt.Sprintf("Found examples section header and %d example-like lines", exampleLineCount)
	} else if exampleLineCount > 0 {
		result.Outcome = engine.PartiallySupports
		result.Remarks = fmt.Sprintf("Found %d example-like lines but no explicit examples section header", exampleLineCount)
	} else {
		result.Outcome = engine.DoesNotSupport
		result.Remarks = "No examples found in help output; consider adding usage examples"
	}

	result.Duration = time.Since(start)
	return result
}

// ---------------------------------------------------------------------------
// HD-7: Consistent Flag Format [AA, SEMI]
// ---------------------------------------------------------------------------

type HD7Check struct{}

func (c *HD7Check) ID() string                         { return "HD-7" }
func (c *HD7Check) Name() string                       { return "Consistent Flag Format" }
func (c *HD7Check) Domain() string                     { return "help" }
func (c *HD7Check) Level() engine.Level                { return engine.LevelAA }
func (c *HD7Check) Testability() engine.Testability    { return engine.Semi }
func (c *HD7Check) SpecVersion() string                { return "1.0" }
func (c *HD7Check) Precondition(p *probe.ProbeResult) bool { return p.HasHelp }

func (c *HD7Check) Run(ctx context.Context, binary string, p *probe.ProbeResult) *engine.Result {
	start := time.Now()
	result := &engine.Result{
		ID: "HD-7", Name: c.Name(), Domain: "help",
		Level: engine.LevelAA, Testability: engine.Semi,
		SpecVersion: "1.0", NeedsReview: true,
	}

	helpText := p.HelpText
	if helpText == "" {
		result.Outcome = engine.OutcomeError
		result.Remarks = "No help text available from probe"
		result.Duration = time.Since(start)
		return result
	}

	result.Evidence = []engine.Evidence{
		{Command: binary + " --help", Stdout: truncate(helpText, 1000), Note: "Help text from probe"},
	}

	lines := strings.Split(helpText, "\n")

	// Extract flag lines: lines with leading whitespace and a dash
	flagLineRe := regexp.MustCompile(`^\s+(-\w|--\w)`)
	shortLongRe := regexp.MustCompile(`-\w.*--\w`)
	longOnlyRe := regexp.MustCompile(`^\s+--\w`)

	var flagLines []string
	shortLongCount := 0
	longOnlyCount := 0
	var indents []int

	for _, line := range lines {
		if !flagLineRe.MatchString(line) {
			continue
		}
		flagLines = append(flagLines, line)

		// Measure indentation
		indent := len(line) - len(strings.TrimLeft(line, " \t"))
		indents = append(indents, indent)

		if shortLongRe.MatchString(line) {
			shortLongCount++
		} else if longOnlyRe.MatchString(line) {
			longOnlyCount++
		}
	}

	if len(flagLines) == 0 {
		result.Outcome = engine.NotApplicable
		result.Remarks = "No flag-like lines detected in help output"
		result.Duration = time.Since(start)
		return result
	}

	// Check indentation consistency
	indentConsistent := true
	if len(indents) > 1 {
		for _, indent := range indents[1:] {
			if indent != indents[0] {
				indentConsistent = false
				break
			}
		}
	}

	// Check format consistency (all short+long or all long-only)
	formatConsistent := shortLongCount == 0 || longOnlyCount == 0

	var issues []string
	if !formatConsistent {
		issues = append(issues, fmt.Sprintf("mixed flag formats: %d short+long, %d long-only", shortLongCount, longOnlyCount))
	}
	if !indentConsistent {
		issues = append(issues, "inconsistent indentation across flag lines")
	}

	if len(issues) == 0 {
		result.Outcome = engine.Supports
		result.Remarks = fmt.Sprintf("%d flags with consistent format and indentation", len(flagLines))
	} else if len(issues) == 1 {
		result.Outcome = engine.PartiallySupports
		result.Remarks = fmt.Sprintf("%d flags; issue: %s", len(flagLines), issues[0])
	} else {
		result.Outcome = engine.DoesNotSupport
		result.Remarks = fmt.Sprintf("%d flags; issues: %s", len(flagLines), strings.Join(issues, "; "))
	}

	result.Duration = time.Since(start)
	return result
}

// ---------------------------------------------------------------------------
// HD-8: No Hang on Empty Stdin [AA, AUTO]
// ---------------------------------------------------------------------------

type HD8Check struct{}

func (c *HD8Check) ID() string                         { return "HD-8" }
func (c *HD8Check) Name() string                       { return "No Hang on Empty Stdin" }
func (c *HD8Check) Domain() string                     { return "help" }
func (c *HD8Check) Level() engine.Level                { return engine.LevelAA }
func (c *HD8Check) Testability() engine.Testability    { return engine.Auto }
func (c *HD8Check) SpecVersion() string                { return "1.0" }
func (c *HD8Check) Precondition(_ *probe.ProbeResult) bool { return true }

func (c *HD8Check) Run(ctx context.Context, binary string, p *probe.ProbeResult) *engine.Result {
	start := time.Now()
	result := &engine.Result{
		ID: "HD-8", Name: c.Name(), Domain: "help",
		Level: engine.LevelAA, Testability: engine.Auto,
		SpecVersion: "1.0",
	}

	// Run the binary with no arguments. In Go's os/exec, stdin defaults to
	// os.DevNull so the process receives immediate EOF — equivalent to
	// running: binary < /dev/null
	out, err := probe.Run(ctx, binary, probe.ExecOpts{
		Timeout: 5 * time.Second,
	})
	if err != nil {
		result.Outcome = engine.OutcomeError
		result.Remarks = fmt.Sprintf("Exec error: %v", err)
		result.Duration = time.Since(start)
		return result
	}

	result.Evidence = []engine.Evidence{
		{
			Command:  out.Command,
			Stdout:   truncate(string(out.Stdout), 200),
			Stderr:   truncate(string(out.Stderr), 200),
			ExitCode: out.ExitCode,
			Duration: out.Duration,
			Note:     "Run with no args, stdin=/dev/null, 5s timeout",
		},
	}

	if out.TimedOut {
		result.Outcome = engine.DoesNotSupport
		result.Remarks = "Binary hangs when run with no arguments and empty stdin (timed out after 5s)"
	} else {
		result.Outcome = engine.Supports
		result.Remarks = fmt.Sprintf("Binary exited within timeout (exit=%d, %v)", out.ExitCode, out.Duration.Round(time.Millisecond))
	}

	result.Duration = time.Since(start)
	return result
}

// ---------------------------------------------------------------------------
// HD-9: Man Page Availability [AA, AUTO]
// ---------------------------------------------------------------------------

type HD9Check struct{}

func (c *HD9Check) ID() string                         { return "HD-9" }
func (c *HD9Check) Name() string                       { return "Man Page Availability" }
func (c *HD9Check) Domain() string                     { return "help" }
func (c *HD9Check) Level() engine.Level                { return engine.LevelAA }
func (c *HD9Check) Testability() engine.Testability    { return engine.Auto }
func (c *HD9Check) SpecVersion() string                { return "1.0" }
func (c *HD9Check) Precondition(_ *probe.ProbeResult) bool { return true }

func (c *HD9Check) Run(ctx context.Context, binary string, p *probe.ProbeResult) *engine.Result {
	start := time.Now()
	result := &engine.Result{
		ID: "HD-9", Name: c.Name(), Domain: "help",
		Level: engine.LevelAA, Testability: engine.Auto,
		SpecVersion: "1.0",
	}

	binaryName := filepath.Base(binary)

	out, err := probe.Run(ctx, "man", probe.ExecOpts{
		Args:    []string{binaryName},
		Timeout: 5 * time.Second,
	})
	if err != nil {
		result.Outcome = engine.OutcomeError
		result.Remarks = fmt.Sprintf("Could not run man: %v", err)
		result.Duration = time.Since(start)
		return result
	}

	result.Evidence = []engine.Evidence{
		{
			Command:  out.Command,
			Stdout:   truncate(string(out.Stdout), 300),
			Stderr:   truncate(string(out.Stderr), 200),
			ExitCode: out.ExitCode,
			Duration: out.Duration,
			Note:     fmt.Sprintf("man %s", binaryName),
		},
	}

	if out.ExitCode == 0 && len(out.Stdout) > 0 {
		result.Outcome = engine.Supports
		result.Remarks = fmt.Sprintf("Man page found for %s", binaryName)
	} else {
		result.Outcome = engine.DoesNotSupport
		result.Remarks = fmt.Sprintf("No man page found for %s (exit=%d)", binaryName, out.ExitCode)
	}

	result.Duration = time.Since(start)
	return result
}

// ---------------------------------------------------------------------------
// HD-10: Typo Suggestion [AA, SEMI]
// ---------------------------------------------------------------------------

type HD10Check struct{}

func (c *HD10Check) ID() string                         { return "HD-10" }
func (c *HD10Check) Name() string                       { return "Typo Suggestion" }
func (c *HD10Check) Domain() string                     { return "help" }
func (c *HD10Check) Level() engine.Level                { return engine.LevelAA }
func (c *HD10Check) Testability() engine.Testability    { return engine.Semi }
func (c *HD10Check) SpecVersion() string                { return "1.0" }
func (c *HD10Check) Precondition(p *probe.ProbeResult) bool { return p.HasSubcommands && len(p.Subcommands) > 0 }

func (c *HD10Check) Run(ctx context.Context, binary string, p *probe.ProbeResult) *engine.Result {
	start := time.Now()
	result := &engine.Result{
		ID: "HD-10", Name: c.Name(), Domain: "help",
		Level: engine.LevelAA, Testability: engine.Semi,
		SpecVersion: "1.0", NeedsReview: true,
	}

	// Pick the first subcommand and drop the last character to simulate a typo
	subcmd := p.Subcommands[0].Name
	if len(subcmd) < 2 {
		result.Outcome = engine.NotApplicable
		result.Remarks = "First subcommand name too short to simulate typo"
		result.Duration = time.Since(start)
		return result
	}
	typo := subcmd[:len(subcmd)-1]

	out, err := probe.Run(ctx, binary, probe.ExecOpts{
		Args:    []string{typo},
		Timeout: 5 * time.Second,
	})
	if err != nil {
		result.Outcome = engine.OutcomeError
		result.Remarks = fmt.Sprintf("Exec error: %v", err)
		result.Duration = time.Since(start)
		return result
	}

	result.Evidence = []engine.Evidence{
		{
			Command:  out.Command,
			Stdout:   truncate(string(out.Stdout), 500),
			Stderr:   truncate(string(out.Stderr), 500),
			ExitCode: out.ExitCode,
			Duration: out.Duration,
			Note:     fmt.Sprintf("Typo probe: %q (original: %q)", typo, subcmd),
		},
	}

	combined := strings.ToLower(string(out.Stdout) + string(out.Stderr))
	suggestionRe := regexp.MustCompile(`(?i)(did you mean|similar|perhaps you meant|most similar|suggest)`)

	if suggestionRe.MatchString(combined) || strings.Contains(combined, subcmd) {
		result.Outcome = engine.Supports
		result.Remarks = fmt.Sprintf("Tool suggests correction for typo %q (original: %q)", typo, subcmd)
	} else {
		result.Outcome = engine.DoesNotSupport
		result.Remarks = fmt.Sprintf("No typo suggestion found for %q; stderr: %s", typo, truncate(string(out.Stderr), 200))
	}

	result.Duration = time.Since(start)
	return result
}

// ---------------------------------------------------------------------------
// HD-11: Web Documentation Link [AAA, SEMI]
// ---------------------------------------------------------------------------

type HD11Check struct{}

func (c *HD11Check) ID() string                         { return "HD-11" }
func (c *HD11Check) Name() string                       { return "Web Documentation Link" }
func (c *HD11Check) Domain() string                     { return "help" }
func (c *HD11Check) Level() engine.Level                { return engine.LevelAAA }
func (c *HD11Check) Testability() engine.Testability    { return engine.Semi }
func (c *HD11Check) SpecVersion() string                { return "1.0" }
func (c *HD11Check) Precondition(p *probe.ProbeResult) bool { return p.HasHelp }

func (c *HD11Check) Run(ctx context.Context, binary string, p *probe.ProbeResult) *engine.Result {
	start := time.Now()
	result := &engine.Result{
		ID: "HD-11", Name: c.Name(), Domain: "help",
		Level: engine.LevelAAA, Testability: engine.Semi,
		SpecVersion: "1.0", NeedsReview: true,
	}

	helpText := p.HelpText
	if helpText == "" {
		result.Outcome = engine.OutcomeError
		result.Remarks = "No help text available from probe"
		result.Duration = time.Since(start)
		return result
	}

	result.Evidence = []engine.Evidence{
		{Command: binary + " --help", Stdout: truncate(helpText, 1000), Note: "Help text from probe"},
	}

	urlRe := regexp.MustCompile(`https?://[^\s)>\]]+`)
	urls := urlRe.FindAllString(helpText, -1)

	if len(urls) > 0 {
		result.Outcome = engine.Supports
		result.Remarks = fmt.Sprintf("Found %d URL(s) in help: %s", len(urls), strings.Join(urls, ", "))
	} else {
		result.Outcome = engine.DoesNotSupport
		result.Remarks = "No URLs (http:// or https://) found in help output"
	}

	result.Duration = time.Since(start)
	return result
}

// ---------------------------------------------------------------------------
// HD-12: Shell Completion Support [AAA, AUTO]
// ---------------------------------------------------------------------------

type HD12Check struct{}

func (c *HD12Check) ID() string                         { return "HD-12" }
func (c *HD12Check) Name() string                       { return "Shell Completion Support" }
func (c *HD12Check) Domain() string                     { return "help" }
func (c *HD12Check) Level() engine.Level                { return engine.LevelAAA }
func (c *HD12Check) Testability() engine.Testability    { return engine.Auto }
func (c *HD12Check) SpecVersion() string                { return "1.0" }
func (c *HD12Check) Precondition(_ *probe.ProbeResult) bool { return true }

func (c *HD12Check) Run(ctx context.Context, binary string, p *probe.ProbeResult) *engine.Result {
	start := time.Now()
	result := &engine.Result{
		ID: "HD-12", Name: c.Name(), Domain: "help",
		Level: engine.LevelAAA, Testability: engine.Auto,
		SpecVersion: "1.0",
	}

	// Try various completion-related subcommands and flags
	completionProbes := [][]string{
		{"completion", "--help"},
		{"completions", "--help"},
		{"completion", "bash"},
		{"completions", "bash"},
		{"--generate-completions"},
		{"--completion"},
		{"--bash-completion"},
	}

	for _, args := range completionProbes {
		out, err := probe.Run(ctx, binary, probe.ExecOpts{
			Args:    args,
			Timeout: 5 * time.Second,
		})
		if err != nil {
			continue
		}

		result.Evidence = append(result.Evidence, engine.Evidence{
			Command:  out.Command,
			Stdout:   truncate(string(out.Stdout), 300),
			Stderr:   truncate(string(out.Stderr), 200),
			ExitCode: out.ExitCode,
			Duration: out.Duration,
			Note:     fmt.Sprintf("Completion probe: %s", strings.Join(args, " ")),
		})

		if out.ExitCode == 0 && (len(out.Stdout) > 0 || len(out.Stderr) > 0) {
			result.Outcome = engine.Supports
			result.Remarks = fmt.Sprintf("Shell completion supported via: %s %s", binary, strings.Join(args, " "))
			result.Duration = time.Since(start)
			return result
		}
	}

	result.Outcome = engine.DoesNotSupport
	result.Remarks = "No shell completion mechanism detected (tried completion, completions, --generate-completions, --completion, --bash-completion)"

	result.Duration = time.Since(start)
	return result
}

// ---------------------------------------------------------------------------
// HD-13: whatis/apropos Compatibility [AAA, AUTO]
// ---------------------------------------------------------------------------

type HD13Check struct{}

func (c *HD13Check) ID() string                         { return "HD-13" }
func (c *HD13Check) Name() string                       { return "whatis/apropos Compatibility" }
func (c *HD13Check) Domain() string                     { return "help" }
func (c *HD13Check) Level() engine.Level                { return engine.LevelAAA }
func (c *HD13Check) Testability() engine.Testability    { return engine.Auto }
func (c *HD13Check) SpecVersion() string                { return "1.0" }
func (c *HD13Check) Precondition(_ *probe.ProbeResult) bool { return true }

func (c *HD13Check) Run(ctx context.Context, binary string, p *probe.ProbeResult) *engine.Result {
	start := time.Now()
	result := &engine.Result{
		ID: "HD-13", Name: c.Name(), Domain: "help",
		Level: engine.LevelAAA, Testability: engine.Auto,
		SpecVersion: "1.0",
	}

	binaryName := filepath.Base(binary)

	out, err := probe.Run(ctx, "whatis", probe.ExecOpts{
		Args:    []string{binaryName},
		Timeout: 5 * time.Second,
	})
	if err != nil {
		result.Outcome = engine.OutcomeError
		result.Remarks = fmt.Sprintf("Could not run whatis: %v", err)
		result.Duration = time.Since(start)
		return result
	}

	result.Evidence = []engine.Evidence{
		{
			Command:  out.Command,
			Stdout:   truncate(string(out.Stdout), 300),
			Stderr:   truncate(string(out.Stderr), 200),
			ExitCode: out.ExitCode,
			Duration: out.Duration,
			Note:     fmt.Sprintf("whatis %s", binaryName),
		},
	}

	if out.ExitCode == 0 && len(out.Stdout) > 0 {
		result.Outcome = engine.Supports
		result.Remarks = fmt.Sprintf("whatis entry found: %s", strings.TrimSpace(string(out.Stdout)))
	} else {
		result.Outcome = engine.DoesNotSupport
		result.Remarks = fmt.Sprintf("No whatis entry for %s (exit=%d)", binaryName, out.ExitCode)
	}

	result.Duration = time.Since(start)
	return result
}

// ---------------------------------------------------------------------------
// HD-14: Help Cross-References [AA, SEMI]
// ---------------------------------------------------------------------------

type HD14Check struct{}

func (c *HD14Check) ID() string                         { return "HD-14" }
func (c *HD14Check) Name() string                       { return "Help Cross-References" }
func (c *HD14Check) Domain() string                     { return "help" }
func (c *HD14Check) Level() engine.Level                { return engine.LevelAA }
func (c *HD14Check) Testability() engine.Testability    { return engine.Semi }
func (c *HD14Check) SpecVersion() string                { return "1.0" }
func (c *HD14Check) Precondition(p *probe.ProbeResult) bool { return p.HasHelp }

func (c *HD14Check) Run(ctx context.Context, binary string, p *probe.ProbeResult) *engine.Result {
	start := time.Now()
	result := &engine.Result{
		ID: "HD-14", Name: c.Name(), Domain: "help",
		Level: engine.LevelAA, Testability: engine.Semi,
		SpecVersion: "1.0", NeedsReview: true,
	}

	helpText := p.HelpText
	if helpText == "" {
		result.Outcome = engine.OutcomeError
		result.Remarks = "No help text available from probe"
		result.Duration = time.Since(start)
		return result
	}

	result.Evidence = []engine.Evidence{
		{Command: binary + " --help", Stdout: truncate(helpText, 1000), Note: "Help text from probe"},
	}

	lower := strings.ToLower(helpText)
	crossRefPatterns := []string{
		"see also",
		"see '",
		`see "`,
		"for more information",
		"for more details",
		"use .* --help",
		"run .* for help",
		"run .* --help",
		"learn more",
		"see the documentation",
		"refer to",
	}

	var found []string
	for _, pattern := range crossRefPatterns {
		re, err := regexp.Compile(`(?i)` + pattern)
		if err != nil {
			continue
		}
		if re.MatchString(lower) {
			matches := re.FindAllString(helpText, 3)
			for _, m := range matches {
				found = append(found, strings.TrimSpace(m))
			}
		}
	}

	if len(found) > 0 {
		result.Outcome = engine.Supports
		result.Remarks = fmt.Sprintf("Found %d cross-reference(s): %s", len(found), strings.Join(found, "; "))
	} else if p.HasSubcommands {
		result.Outcome = engine.DoesNotSupport
		result.Remarks = "No cross-references found; tools with subcommands SHOULD include navigational pointers"
	} else {
		result.Outcome = engine.PartiallySupports
		result.Remarks = "No cross-references found; simple tools may not need them, but consider adding 'See also' or help pointers"
	}

	result.Duration = time.Since(start)
	return result
}

// ---------------------------------------------------------------------------
// HD-15: Grouped/Categorized Help [AA, SEMI]
// ---------------------------------------------------------------------------

type HD15Check struct{}

func (c *HD15Check) ID() string                         { return "HD-15" }
func (c *HD15Check) Name() string                       { return "Grouped/Categorized Help" }
func (c *HD15Check) Domain() string                     { return "help" }
func (c *HD15Check) Level() engine.Level                { return engine.LevelAA }
func (c *HD15Check) Testability() engine.Testability    { return engine.Semi }
func (c *HD15Check) SpecVersion() string                { return "1.0" }
func (c *HD15Check) Precondition(p *probe.ProbeResult) bool { return p.HasHelp }

func (c *HD15Check) Run(ctx context.Context, binary string, p *probe.ProbeResult) *engine.Result {
	start := time.Now()
	result := &engine.Result{
		ID: "HD-15", Name: c.Name(), Domain: "help",
		Level: engine.LevelAA, Testability: engine.Semi,
		SpecVersion: "1.0", NeedsReview: true,
	}

	helpText := p.HelpText
	if helpText == "" {
		result.Outcome = engine.OutcomeError
		result.Remarks = "No help text available from probe"
		result.Duration = time.Since(start)
		return result
	}

	result.Evidence = []engine.Evidence{
		{Command: binary + " --help", Stdout: truncate(helpText, 1000), Note: "Help text from probe"},
	}

	// Count flags in top-level help
	lines := strings.Split(helpText, "\n")
	flagLineRe := regexp.MustCompile(`^\s+--?\w`)
	flagCount := 0
	for _, line := range lines {
		if flagLineRe.MatchString(line) {
			flagCount++
		}
	}

	// If top-level has few flags but has subcommands, sample a few
	// subcommands' help to find the one with the most flags — tools like
	// cosign have 4 global flags but 50+ per subcommand.
	maxSubcmdFlags := 0
	maxSubcmdName := ""
	if flagCount < 20 && p.HasSubcommands && len(p.Subcommands) > 0 {
		sampled := p.Subcommands
		if len(sampled) > 5 {
			sampled = sampled[:5]
		}
		for _, sub := range sampled {
			subFlags := 0
			subOut, err := probe.Run(ctx, binary, probe.ExecOpts{
				Args:    []string{sub.Name, "--help"},
				Timeout: 5 * time.Second,
			})
			if err == nil && subOut.ExitCode == 0 {
				for _, l := range strings.Split(string(subOut.Stdout), "\n") {
					if flagLineRe.MatchString(l) {
						subFlags++
					}
				}
			}
			if subFlags > maxSubcmdFlags {
				maxSubcmdFlags = subFlags
				maxSubcmdName = sub.Name
			}
		}
		_ = maxSubcmdName
	}

	effectiveFlags := flagCount
	if maxSubcmdFlags > effectiveFlags {
		effectiveFlags = maxSubcmdFlags
	}

	if effectiveFlags < 20 {
		result.Outcome = engine.NotApplicable
		result.Remarks = fmt.Sprintf("Only %d flags detected (global: %d, largest subcommand: %d); grouping recommended for tools with 20+ flags", effectiveFlags, flagCount, maxSubcmdFlags)
		result.Duration = time.Since(start)
		return result
	}

	// Check for section headings (lines that look like categories)
	// Section headings are typically: uppercase or title-case words followed
	// by a colon, or all-caps lines, or lines that are not indented and are
	// followed by indented content
	sectionHeadingRe := regexp.MustCompile(`(?i)^[A-Z][\w\s]*:\s*$|^[A-Z][A-Z\s]+$`)
	sectionCount := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if sectionHeadingRe.MatchString(trimmed) {
			sectionCount++
		}
	}

	if sectionCount >= 3 {
		result.Outcome = engine.Supports
		result.Remarks = fmt.Sprintf("%d flags organized under %d section headings", flagCount, sectionCount)
	} else if sectionCount >= 1 {
		result.Outcome = engine.PartiallySupports
		result.Remarks = fmt.Sprintf("%d flags with only %d section heading(s); consider more grouping", flagCount, sectionCount)
	} else {
		result.Outcome = engine.DoesNotSupport
		result.Remarks = fmt.Sprintf("%d flags presented as a flat list without section headings or categories", flagCount)
	}

	result.Duration = time.Since(start)
	return result
}

// ---------------------------------------------------------------------------
// HD-16: Flag Value Enumeration [AAA, SEMI]
// ---------------------------------------------------------------------------

type HD16Check struct{}

func (c *HD16Check) ID() string                         { return "HD-16" }
func (c *HD16Check) Name() string                       { return "Flag Value Enumeration" }
func (c *HD16Check) Domain() string                     { return "help" }
func (c *HD16Check) Level() engine.Level                { return engine.LevelAAA }
func (c *HD16Check) Testability() engine.Testability    { return engine.Semi }
func (c *HD16Check) SpecVersion() string                { return "1.0" }
func (c *HD16Check) Precondition(p *probe.ProbeResult) bool { return p.HasHelp }

func (c *HD16Check) Run(ctx context.Context, binary string, p *probe.ProbeResult) *engine.Result {
	start := time.Now()
	result := &engine.Result{
		ID: "HD-16", Name: c.Name(), Domain: "help",
		Level: engine.LevelAAA, Testability: engine.Semi,
		SpecVersion: "1.0", NeedsReview: true,
	}

	helpText := p.HelpText
	if helpText == "" {
		result.Outcome = engine.OutcomeError
		result.Remarks = "No help text available from probe"
		result.Duration = time.Since(start)
		return result
	}

	result.Evidence = []engine.Evidence{
		{Command: binary + " --help", Stdout: truncate(helpText, 1000), Note: "Help text from probe"},
	}

	// Look for enum-value patterns in flag descriptions:
	// - "one of: ..." or "one of ..."
	// - "[value1|value2|value3]"
	// - "{value1,value2,value3}"
	// - "(value1|value2|value3)"
	// - "format: json, table, yaml"
	// - "choices: ..."
	// - "valid values: ..."
	enumPatterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)\bone of[:\s]+\w+[,|/]\s*\w+`),
		regexp.MustCompile(`\[(\w+\|)+\w+\]`),
		regexp.MustCompile(`\{(\w+,\s*)+\w+\}`),
		regexp.MustCompile(`\((\w+\|)+\w+\)`),
		regexp.MustCompile(`(?i)\b(choices|valid values|allowed values|possible values)\s*[:\s]+`),
		regexp.MustCompile(`(?i):\s+\w+,\s+\w+,\s+\w+`),         // "format: json, table, yaml"
		regexp.MustCompile(`(?i)\b(always|never|auto)\b.*\b(always|never|auto)\b`), // common tristate
	}

	var foundEnums []string
	for _, re := range enumPatterns {
		matches := re.FindAllString(helpText, 5)
		for _, m := range matches {
			foundEnums = append(foundEnums, strings.TrimSpace(m))
		}
	}

	// Count flags that take values — check both probe data AND help text.
	// Some tools (cosign) use --flag=default format that the probe may not
	// detect as TakesValue, so also scan help text for = patterns in flags.
	valueFlagCount := 0
	for _, f := range p.GlobalFlags {
		if f.TakesValue {
			valueFlagCount++
		}
	}
	// Also detect flags with = in help text (e.g., --timeout=3m0s, --format=json)
	flagWithValueRe := regexp.MustCompile(`--\w[\w-]+=\S+`)
	helpFlagValues := flagWithValueRe.FindAllString(helpText, -1)
	if len(helpFlagValues) > valueFlagCount {
		valueFlagCount = len(helpFlagValues)
	}

	if len(foundEnums) > 0 {
		result.Outcome = engine.Supports
		result.Remarks = fmt.Sprintf("Found %d flag value enumeration(s): %s", len(foundEnums), strings.Join(dedup(foundEnums), "; "))
	} else if valueFlagCount == 0 {
		result.Outcome = engine.NotApplicable
		result.Remarks = "No flags that take values detected"
	} else {
		result.Outcome = engine.DoesNotSupport
		result.Remarks = fmt.Sprintf("%d flags take values but none enumerate valid options in help text", valueFlagCount)
	}

	result.Duration = time.Since(start)
	return result
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// truncate returns at most maxLen bytes from s, appending "..." if truncated.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// dedup returns unique strings preserving order.
func dedup(ss []string) []string {
	seen := make(map[string]bool)
	var out []string
	for _, s := range ss {
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	return out
}
