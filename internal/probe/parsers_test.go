package probe

import "testing"

const cobraHelp = `A CLI tool for things

Usage:
  mytool [command]

Available Commands:
  list        List all items
  get         Get a specific item
  completion  Generate shell completions
  help        Help about any command

Flags:
  -h, --help      help for mytool
  -v, --version   version for mytool
      --json      Output as JSON
  -q, --quiet     Suppress output

Use "mytool [command] --help" for more information about a command.
`

const clapHelp = `mytool 1.0
A tool for doing things

USAGE:
    mytool [OPTIONS] [SUBCOMMANDS]

OPTIONS:
    -h, --help       Print help information
    -V, --version    Print version information
        --json       Output as JSON
    -q, --quiet      Suppress output

SUBCOMMANDS:
    list    List all items
    get     Get a specific item
    help    Print this message or the help of the given subcommand(s)
`

const argparseHelp = `usage: mytool [-h] [--json] [-q] {list,get} ...

A tool for doing things

positional arguments:
  {list,get}
    list        List all items
    get         Get a specific item

options:
  -h, --help    show this help message and exit
  --json        Output as JSON
  -q, --quiet   Suppress output
`

const genericHelp = `mytool: A simple tool

Usage: mytool [options] <command>

Options:
  -h, --help      Show this help
  -v, --version   Show version
  --json          Output as JSON
  -q, --quiet     Suppress output

Commands:
  list    List all items
  get     Get a specific item
`

func TestCobraParserCanParse(t *testing.T) {
	p := &CobraParser{}
	if !p.CanParse(cobraHelp) {
		t.Error("CobraParser should recognize cobra-style help")
	}
	if p.CanParse("usage: tool [options] file") {
		t.Error("CobraParser should not match non-cobra help")
	}
}

func TestCobraParserExtractsSubcommands(t *testing.T) {
	p := &CobraParser{}
	result, err := p.Parse(cobraHelp)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Subcommands) < 3 {
		t.Errorf("expected >=3 subcommands, got %d", len(result.Subcommands))
	}
	found := false
	for _, s := range result.Subcommands {
		if s.Name == "list" {
			found = true
		}
	}
	if !found {
		t.Error("expected 'list' subcommand")
	}
}

func TestCobraParserExtractsFlags(t *testing.T) {
	p := &CobraParser{}
	result, err := p.Parse(cobraHelp)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Found %d flags", len(result.Flags))
	for _, f := range result.Flags {
		t.Logf("Flag: short=%q long=%q desc=%q", f.Short, f.Long, f.Description)
	}
	foundJSON := false
	foundQuiet := false
	for _, f := range result.Flags {
		if f.Long == "--json" {
			foundJSON = true
		}
		if f.Long == "--quiet" {
			foundQuiet = true
		}
	}
	if !foundJSON {
		t.Error("expected --json flag")
	}
	if !foundQuiet {
		t.Error("expected --quiet flag")
	}
}

func TestClapParserCanParse(t *testing.T) {
	p := &ClapParser{}
	if !p.CanParse(clapHelp) {
		t.Error("ClapParser should recognize clap-style help")
	}
	if p.CanParse(cobraHelp) {
		t.Error("ClapParser should not match cobra help")
	}
}

func TestClapParserExtractsSubcommands(t *testing.T) {
	p := &ClapParser{}
	result, err := p.Parse(clapHelp)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Subcommands) < 2 {
		t.Errorf("expected >=2 subcommands, got %d", len(result.Subcommands))
	}
	found := false
	for _, s := range result.Subcommands {
		if s.Name == "list" {
			found = true
		}
	}
	if !found {
		t.Error("expected 'list' subcommand")
	}
}

func TestClapParserExtractsFlags(t *testing.T) {
	p := &ClapParser{}
	result, err := p.Parse(clapHelp)
	if err != nil {
		t.Fatal(err)
	}
	foundJSON := false
	for _, f := range result.Flags {
		if f.Long == "--json" {
			foundJSON = true
		}
	}
	if !foundJSON {
		t.Error("expected --json flag")
	}
}

func TestArgparseParserCanParse(t *testing.T) {
	p := &ArgparseParser{}
	if !p.CanParse(argparseHelp) {
		t.Error("ArgparseParser should recognize argparse-style help")
	}
	if p.CanParse(cobraHelp) {
		t.Error("ArgparseParser should not match cobra help")
	}
}

func TestArgparseParserExtractsSubcommands(t *testing.T) {
	p := &ArgparseParser{}
	result, err := p.Parse(argparseHelp)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Subcommands) < 2 {
		t.Errorf("expected >=2 subcommands, got %d", len(result.Subcommands))
	}
	found := false
	for _, s := range result.Subcommands {
		if s.Name == "list" {
			found = true
		}
	}
	if !found {
		t.Error("expected 'list' subcommand")
	}
}

func TestArgparseParserExtractsFlags(t *testing.T) {
	p := &ArgparseParser{}
	result, err := p.Parse(argparseHelp)
	if err != nil {
		t.Fatal(err)
	}
	foundJSON := false
	for _, f := range result.Flags {
		if f.Long == "--json" {
			foundJSON = true
		}
	}
	if !foundJSON {
		t.Error("expected --json flag")
	}
}

func TestGenericParserCanParse(t *testing.T) {
	p := &GenericParser{}
	// GenericParser should always return true as it's the fallback
	if !p.CanParse(genericHelp) {
		t.Error("GenericParser should always be able to parse")
	}
	if !p.CanParse("random text") {
		t.Error("GenericParser should always be able to parse")
	}
}

func TestGenericParserExtractsFlags(t *testing.T) {
	p := &GenericParser{}
	result, err := p.Parse(genericHelp)
	if err != nil {
		t.Fatal(err)
	}
	foundJSON := false
	for _, f := range result.Flags {
		if f.Long == "--json" {
			foundJSON = true
		}
	}
	if !foundJSON {
		t.Error("expected --json flag")
	}
}

func TestGenericParserExtractsSubcommands(t *testing.T) {
	p := &GenericParser{}
	result, err := p.Parse(genericHelp)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Subcommands) < 2 {
		t.Errorf("expected >=2 subcommands, got %d", len(result.Subcommands))
	}
	found := false
	for _, s := range result.Subcommands {
		if s.Name == "list" {
			found = true
		}
	}
	if !found {
		t.Error("expected 'list' subcommand")
	}
}

func TestParserCascade(t *testing.T) {
	tests := []struct {
		name     string
		helpText string
		wantName string
	}{
		{"cobra", cobraHelp, "cobra"},
		{"clap", clapHelp, "clap"},
		{"argparse", argparseHelp, "argparse"},
		{"generic", genericHelp, "generic"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsers := []HelpParser{
				&CobraParser{},
				&ClapParser{},
				&ArgparseParser{},
				&GenericParser{},
			}
			var matched HelpParser
			for _, p := range parsers {
				if p.CanParse(tt.helpText) {
					matched = p
					break
				}
			}
			if matched == nil {
				t.Fatal("no parser matched")
			}
			if matched.Name() != tt.wantName {
				t.Errorf("got parser %q, want %q", matched.Name(), tt.wantName)
			}
		})
	}
}
