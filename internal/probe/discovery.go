package probe

import (
	"context"
	"path/filepath"
	"strings"
	"time"
)

// Probe auto-discovers a binary's capabilities by running help, version,
// parsing help text, detecting color, and triggering errors.
func Probe(ctx context.Context, binary string, maxSubcmds int) (*ProbeResult, error) {
	result := &ProbeResult{
		BinaryPath:     binary,
		BinaryName:     filepath.Base(binary),
		Subcommands:    []Subcommand{},
		GlobalFlags:    []Flag{},
		SampleErrors:   []ErrorSample{},
		HasHelp:        false,
		HasVersion:     false,
		HasSubcommands: false,
		HasColor:       false,
		HasJSONFlag:    false,
		HasQuietFlag:   false,
		HasColorFlag:   false,
		HasDryRunFlag:  false,
		HasNoInputFlag: false,
		HelpFormat:     "unknown",
	}

	// Step 1: Try --help
	helpResult, helpErr := Run(ctx, binary, ExecOpts{
		Args:    []string{"--help"},
		Timeout: 5 * time.Second,
	})

	// If we get an error (not just a non-zero exit), the binary likely doesn't exist
	if helpErr != nil {
		return nil, helpErr
	}

	if helpResult.ExitCode == 0 {
		result.HasHelp = true
		// Prefer stdout, but fall back to stderr (some tools output help to stderr)
		helpText := string(helpResult.Stdout)
		if len(helpText) == 0 {
			helpText = string(helpResult.Stderr)
		}
		result.HelpText = helpText

		// Detect ANSI usage in help output (piped mode).
		// Check for any ANSI sequences (color, bold, underline) — not just
		// color codes. Tools that use bold/underline still need NO_COLOR,
		// TERM=dumb, and TTY-awareness testing.
		if HasANSI(helpResult.Stdout) || HasANSI(helpResult.Stderr) {
			result.HasColor = true
		}

		// If no ANSI detected in pipe mode, try forcing color output.
		// TTY-aware tools suppress all ANSI when piped. Probe with
		// CLICOLOR_FORCE=1 (BSD convention, used by gh) and FORCE_COLOR=1
		// (modern standard) to detect tools that use ANSI on a real terminal.
		if !result.HasColor {
			forceColorResult, forceErr := Run(ctx, binary, ExecOpts{
				Args: []string{"--help"},
				Env: map[string]string{
					"CLICOLOR_FORCE": "1",
					"FORCE_COLOR":    "1",
				},
				Timeout: 5 * time.Second,
			})
			if forceErr == nil && forceColorResult.ExitCode == 0 {
				if HasANSI(forceColorResult.Stdout) || HasANSI(forceColorResult.Stderr) {
					result.HasColor = true
				}
			}
		}

		// If still nothing, try --color=always flag (common convention)
		if !result.HasColor {
			colorFlagResult, colorErr := Run(ctx, binary, ExecOpts{
				Args:    []string{"--color=always", "--help"},
				Timeout: 5 * time.Second,
			})
			if colorErr == nil && colorFlagResult.ExitCode == 0 {
				if HasANSI(colorFlagResult.Stdout) || HasANSI(colorFlagResult.Stderr) {
					result.HasColor = true
					result.HasColorFlag = true
				}
			}
		}

		// Parse help text using parser cascade
		parsed := parseHelpText(helpText)
		if parsed != nil {
			result.HelpFormat = parsed.format

			// Convert parsed subcommands to Subcommand
			for _, s := range parsed.help.Subcommands {
				result.Subcommands = append(result.Subcommands, Subcommand{
					Name:     s.Name,
					HelpText: s.HelpText,
					Flags:    []Flag{},
				})
			}

			// Convert parsed flags to Flag
			for _, f := range parsed.help.Flags {
				result.GlobalFlags = append(result.GlobalFlags, Flag{
					Short:       f.Short,
					Long:        f.Long,
					Description: f.Description,
					TakesValue:  f.TakesValue,
				})
			}
		}
	}

	// Step 2: Try --version
	versionResult, versionErr := Run(ctx, binary, ExecOpts{
		Args:    []string{"--version"},
		Timeout: 5 * time.Second,
	})

	if versionErr == nil && versionResult.ExitCode == 0 {
		result.HasVersion = true
		// Extract version from stdout or stderr
		versionText := strings.TrimSpace(string(versionResult.Stdout))
		if len(versionText) == 0 {
			versionText = strings.TrimSpace(string(versionResult.Stderr))
		}
		// Take first line only
		if idx := strings.IndexByte(versionText, '\n'); idx > 0 {
			versionText = versionText[:idx]
		}
		result.Version = versionText
	}

	// Step 3: Detect common flags
	result.HasJSONFlag = hasFlag(result.GlobalFlags, "--json")
	result.HasQuietFlag = hasFlag(result.GlobalFlags, "--quiet", "-q")
	result.HasColorFlag = hasFlag(result.GlobalFlags, "--color")
	result.HasDryRunFlag = hasFlag(result.GlobalFlags, "--dry-run")
	result.HasNoInputFlag = hasFlag(result.GlobalFlags, "--no-input", "--non-interactive", "-n")

	// Step 4: Filter and limit subcommands
	if len(result.Subcommands) > 0 {
		result.HasSubcommands = true
		filtered := []Subcommand{}
		for _, sc := range result.Subcommands {
			// Skip blocked subcommands
			if IsBlockedSubcommand(sc.Name) {
				continue
			}
			// Skip common help/completion commands
			if sc.Name == "help" || sc.Name == "completion" {
				continue
			}
			filtered = append(filtered, sc)
			if len(filtered) >= maxSubcmds {
				break
			}
		}
		result.Subcommands = filtered
	}

	// Step 5: Trigger error with nonexistent flag to get error output sample
	errorResult, _ := Run(ctx, binary, ExecOpts{
		Args:    []string{"--nonexistent-flag-a11y-probe"},
		Timeout: 5 * time.Second,
	})

	if errorResult != nil && errorResult.ExitCode != 0 {
		result.SampleErrors = append(result.SampleErrors, ErrorSample{
			Trigger:  "--nonexistent-flag-a11y-probe",
			Stdout:   errorResult.Stdout,
			Stderr:   errorResult.Stderr,
			ExitCode: errorResult.ExitCode,
		})
	}

	return result, nil
}

type parsedResult struct {
	format string
	help   *ParsedHelp
}

// parseHelpText runs the parser cascade on help text
func parseHelpText(text string) *parsedResult {
	if len(text) == 0 {
		return nil
	}

	parsers := []HelpParser{
		&CobraParser{},
		&ClapParser{},
		&ArgparseParser{},
		&GenericParser{},
	}

	for _, p := range parsers {
		if p.CanParse(text) {
			parsed, err := p.Parse(text)
			if err == nil && parsed != nil {
				return &parsedResult{
					format: p.Name(),
					help:   parsed,
				}
			}
		}
	}

	return nil
}

// hasFlag checks if any flag matches the given names
func hasFlag(flags []Flag, names ...string) bool {
	for _, f := range flags {
		for _, name := range names {
			if f.Long == name || f.Short == name {
				return true
			}
		}
	}
	return false
}
