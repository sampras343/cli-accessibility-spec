package probe

import (
	"bufio"
	"regexp"
	"strings"
)

type HelpParser interface {
	Name() string
	CanParse(text string) bool
	Parse(text string) (*ParsedHelp, error)
}

type ParsedHelp struct {
	Subcommands []ParsedSubcommand
	Flags       []ParsedFlag
}

type ParsedSubcommand struct {
	Name     string
	HelpText string
}

type ParsedFlag struct {
	Short       string
	Long        string
	Description string
	TakesValue  bool
}

// CobraParser parses help text from Cobra-based CLIs (Go)
type CobraParser struct{}

func (p *CobraParser) Name() string {
	return "cobra"
}

func (p *CobraParser) CanParse(text string) bool {
	lower := strings.ToLower(text)
	// Look for cobra-specific patterns
	hasCobraCommands := strings.Contains(lower, "available commands:") ||
		strings.Contains(text, "CORE COMMANDS") ||
		strings.Contains(text, "GITHUB ACTIONS COMMANDS")
	hasCobraUsage := strings.Contains(lower, "use \"") && strings.Contains(lower, "--help")
	hasFlags := strings.Contains(lower, "flags:")

	return hasCobraCommands || (hasCobraUsage && hasFlags)
}

func (p *CobraParser) Parse(text string) (*ParsedHelp, error) {
	result := &ParsedHelp{
		Subcommands: []ParsedSubcommand{},
		Flags:       []ParsedFlag{},
	}

	scanner := bufio.NewScanner(strings.NewReader(text))
	inCommands := false
	inFlags := false

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		lower := strings.ToLower(trimmed)

		// Check for command sections (case-insensitive)
		// Matches: "Available Commands:", "CORE COMMANDS", "Commands:"
		isCommandSection := strings.Contains(lower, "commands") && !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t")
		if isCommandSection {
			inCommands = true
			inFlags = false
			continue
		}

		// Check for flags section (case-insensitive)
		// Matches: "Flags:", "FLAGS", "Options:"
		isFlagSection := (strings.HasPrefix(lower, "flags") || strings.HasPrefix(lower, "options")) &&
			!strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t")
		if isFlagSection {
			inFlags = true
			inCommands = false
			continue
		}

		// Check for end markers
		if strings.HasPrefix(line, "Use ") || strings.HasPrefix(line, "EXAMPLES") ||
		   strings.HasPrefix(line, "LEARN MORE") || strings.HasPrefix(line, "ENVIRONMENT") {
			inCommands = false
			inFlags = false
			continue
		}

		if len(trimmed) == 0 {
			// Don't reset on empty lines - sections might have blank lines
			continue
		}

		if inCommands {
			// Parse subcommand line. Two formats:
			// Indented: "  list        List all items" (most cobra tools)
			// Flush:    "list                  List all items" (cosign-style)
			// Both have a command name followed by 2+ spaces then description.
			isIndented := strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t")
			hasMultiSpaceGap := regexp.MustCompile(`^\S+\s{2,}`).MatchString(trimmed)
			if isIndented || hasMultiSpaceGap {
				parts := strings.Fields(trimmed)
				if len(parts) >= 1 {
					name := strings.TrimSuffix(parts[0], ":")
					desc := ""
					if len(parts) > 1 {
						desc = strings.Join(parts[1:], " ")
					}
					result.Subcommands = append(result.Subcommands, ParsedSubcommand{
						Name:     name,
						HelpText: desc,
					})
				}
			}
		}

		if inFlags {
			// Parse flag line: "  -h, --help      help for mytool"
			// or: "      --json      Output as JSON"
			if strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") {
				result.Flags = append(result.Flags, parseFlag(line)...)
			}
		}
	}

	return result, nil
}

// ClapParser parses help text from Clap-based CLIs (Rust)
type ClapParser struct{}

func (p *ClapParser) Name() string {
	return "clap"
}

func (p *ClapParser) CanParse(text string) bool {
	hasUsage := strings.Contains(text, "USAGE:") || strings.Contains(text, "Usage:")
	hasOptions := strings.Contains(text, "OPTIONS:") || strings.Contains(text, "Args:")
	return hasUsage && hasOptions
}

func (p *ClapParser) Parse(text string) (*ParsedHelp, error) {
	result := &ParsedHelp{
		Subcommands: []ParsedSubcommand{},
		Flags:       []ParsedFlag{},
	}

	scanner := bufio.NewScanner(strings.NewReader(text))
	inSubcommands := false
	inOptions := false

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(strings.TrimSpace(line), "SUBCOMMANDS:") {
			inSubcommands = true
			inOptions = false
			continue
		}
		if strings.HasPrefix(strings.TrimSpace(line), "OPTIONS:") ||
			strings.HasPrefix(strings.TrimSpace(line), "Args:") {
			inOptions = true
			inSubcommands = false
			continue
		}
		if len(strings.TrimSpace(line)) == 0 {
			inSubcommands = false
			inOptions = false
			continue
		}

		if inSubcommands {
			// Parse subcommand line: "    list    List all items"
			trimmed := strings.TrimSpace(line)
			parts := strings.Fields(trimmed)
			if len(parts) >= 1 {
				name := parts[0]
				desc := ""
				if len(parts) > 1 {
					desc = strings.Join(parts[1:], " ")
				}
				result.Subcommands = append(result.Subcommands, ParsedSubcommand{
					Name:     name,
					HelpText: desc,
				})
			}
		}

		if inOptions {
			// Parse option line: "    -h, --help       Print help information"
			result.Flags = append(result.Flags, parseFlag(line)...)
		}
	}

	return result, nil
}

// ArgparseParser parses help text from argparse-based CLIs (Python)
type ArgparseParser struct{}

func (p *ArgparseParser) Name() string {
	return "argparse"
}

func (p *ArgparseParser) CanParse(text string) bool {
	hasUsage := strings.Contains(text, "usage:")
	hasOptions := strings.Contains(text, "options:") || strings.Contains(text, "optional arguments:")
	hasPositional := strings.Contains(text, "positional arguments:")
	return hasUsage && (hasOptions || hasPositional)
}

func (p *ArgparseParser) Parse(text string) (*ParsedHelp, error) {
	result := &ParsedHelp{
		Subcommands: []ParsedSubcommand{},
		Flags:       []ParsedFlag{},
	}

	scanner := bufio.NewScanner(strings.NewReader(text))
	inPositional := false
	inOptions := false

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "positional arguments:") {
			inPositional = true
			inOptions = false
			continue
		}
		if strings.HasPrefix(line, "options:") || strings.HasPrefix(line, "optional arguments:") {
			inOptions = true
			inPositional = false
			continue
		}
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}
		if !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t") {
			inPositional = false
			inOptions = false
			continue
		}

		if inPositional {
			// Parse positional line: "  {list,get}"
			// followed by: "    list        List all items"
			trimmed := strings.TrimSpace(line)
			if !strings.HasPrefix(trimmed, "{") && !strings.HasPrefix(trimmed, "[") {
				parts := strings.Fields(trimmed)
				if len(parts) >= 1 && !strings.HasPrefix(parts[0], "-") {
					name := parts[0]
					desc := ""
					if len(parts) > 1 {
						desc = strings.Join(parts[1:], " ")
					}
					result.Subcommands = append(result.Subcommands, ParsedSubcommand{
						Name:     name,
						HelpText: desc,
					})
				}
			}
		}

		if inOptions {
			// Parse option line: "  -h, --help    show this help message and exit"
			result.Flags = append(result.Flags, parseFlag(line)...)
		}
	}

	return result, nil
}

// GenericParser is a fallback parser that uses simple heuristics
type GenericParser struct{}

func (p *GenericParser) Name() string {
	return "generic"
}

func (p *GenericParser) CanParse(text string) bool {
	// GenericParser always returns true as it's the fallback
	return true
}

func (p *GenericParser) Parse(text string) (*ParsedHelp, error) {
	result := &ParsedHelp{
		Subcommands: []ParsedSubcommand{},
		Flags:       []ParsedFlag{},
	}

	scanner := bufio.NewScanner(strings.NewReader(text))
	inCommandSection := false

	for scanner.Scan() {
		line := scanner.Text()

		// Detect command sections
		lower := strings.ToLower(line)
		if strings.Contains(lower, "command") && strings.HasSuffix(lower, ":") {
			inCommandSection = true
			continue
		}

		// Reset section if we hit an empty line or new section
		if len(strings.TrimSpace(line)) == 0 {
			inCommandSection = false
			continue
		}

		// Parse flags: lines starting with - or --
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "-") {
			result.Flags = append(result.Flags, parseFlag(line)...)
		}

		// Parse subcommands: indented lines in command section
		if inCommandSection && (strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t")) {
			parts := strings.Fields(trimmed)
			if len(parts) >= 1 && !strings.HasPrefix(parts[0], "-") {
				name := parts[0]
				desc := ""
				if len(parts) > 1 {
					desc = strings.Join(parts[1:], " ")
				}
				result.Subcommands = append(result.Subcommands, ParsedSubcommand{
					Name:     name,
					HelpText: desc,
				})
			}
		}
	}

	return result, nil
}

// parseFlag extracts flags from a line like "  -h, --help      help text"
func parseFlag(line string) []ParsedFlag {
	var flags []ParsedFlag

	trimmed := strings.TrimSpace(line)

	// Pattern 1: "-h, --help description" (short and long)
	// Pattern 2: "--help description" (long only)
	// Pattern 3: "-h description" (short only)

	// First, extract all flag tokens
	parts := strings.Fields(trimmed)
	if len(parts) == 0 {
		return flags
	}

	short := ""
	long := ""
	desc := ""
	descStartIdx := -1

	for i, part := range parts {
		// Remove comma if present
		part = strings.TrimSuffix(part, ",")

		if strings.HasPrefix(part, "--") {
			// Long flag
			long = part
		} else if strings.HasPrefix(part, "-") && len(part) == 2 {
			// Short flag (single char after -)
			short = part
		} else if !strings.HasPrefix(part, "-") {
			// This is the start of the description
			descStartIdx = i
			break
		}
	}

	if descStartIdx > 0 {
		desc = strings.Join(parts[descStartIdx:], " ")
	}

	if short != "" || long != "" {
		flags = append(flags, ParsedFlag{
			Short:       short,
			Long:        long,
			Description: desc,
			TakesValue:  false, // TODO: detect if flag takes value
		})
	}

	return flags
}
