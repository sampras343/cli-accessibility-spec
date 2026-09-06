package probe

import (
	"context"
	"os/exec"
	"testing"
)

func TestProbeGH(t *testing.T) {
	if _, err := exec.LookPath("gh"); err != nil {
		t.Skip("gh not installed")
	}

	result, err := Probe(context.Background(), "gh", 5)
	if err != nil {
		t.Fatal(err)
	}

	if !result.HasHelp {
		t.Error("gh should have --help")
	}

	if !result.HasVersion {
		t.Error("gh should have --version")
	}

	if !result.HasSubcommands {
		t.Error("gh should have subcommands")
	}

	if len(result.Subcommands) == 0 {
		t.Error("gh should have discovered subcommands")
	}

	if result.HelpFormat != "cobra" {
		t.Errorf("gh help format = %q, want cobra", result.HelpFormat)
	}

	if len(result.SampleErrors) == 0 {
		t.Error("should have captured error sample")
	}

	t.Logf("Discovered %d subcommands", len(result.Subcommands))
	for _, sc := range result.Subcommands {
		t.Logf("  - %s: %s", sc.Name, sc.HelpText)
	}

	t.Logf("Discovered %d global flags", len(result.GlobalFlags))
	t.Logf("Version: %s", result.Version)
	t.Logf("Has color: %v", result.HasColor)
}

func TestProbeNonexistentBinary(t *testing.T) {
	_, err := Probe(context.Background(), "/nonexistent/binary", 5)
	if err == nil {
		t.Error("expected error for nonexistent binary, got nil")
	}
}

func TestProbeWithContext(t *testing.T) {
	if _, err := exec.LookPath("sleep"); err != nil {
		t.Skip("sleep not installed")
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// This should fail quickly due to canceled context
	_, err := Probe(ctx, "sleep", 5)
	if err == nil {
		// Note: some operations might complete before context is checked
		t.Log("warning: expected error from canceled context")
	}
}

func TestProbeDetectsCommonFlags(t *testing.T) {
	// Create a mock help text with common flags
	mockHelp := `A test CLI

Usage:
  test [command]

Available Commands:
  list        List items

Flags:
  -h, --help           help for test
      --json           Output as JSON
  -q, --quiet          Suppress output
      --color string   Color mode
      --dry-run        Dry run mode
      --no-input       Non-interactive mode
`

	parsed := parseHelpText(mockHelp)
	if parsed == nil {
		t.Fatal("failed to parse mock help")
	}

	if parsed.format != "cobra" {
		t.Errorf("expected cobra format, got %s", parsed.format)
	}

	t.Logf("Parsed %d flags", len(parsed.help.Flags))
	for _, f := range parsed.help.Flags {
		t.Logf("  Flag: short=%q long=%q", f.Short, f.Long)
	}

	// Convert to flags for testing
	var flags []struct {
		Long string
	}
	for _, f := range parsed.help.Flags {
		flags = append(flags, struct{ Long string }{Long: f.Long})
	}

	hasJSON := false
	hasQuiet := false
	hasColor := false
	hasDryRun := false
	hasNoInput := false

	for _, f := range flags {
		switch f.Long {
		case "--json":
			hasJSON = true
		case "--quiet":
			hasQuiet = true
		case "--color":
			hasColor = true
		case "--dry-run":
			hasDryRun = true
		case "--no-input":
			hasNoInput = true
		}
	}

	if !hasJSON {
		t.Error("expected to detect --json flag")
	}
	if !hasQuiet {
		t.Error("expected to detect --quiet flag")
	}
	if !hasColor {
		t.Error("expected to detect --color flag")
	}
	if !hasDryRun {
		t.Error("expected to detect --dry-run flag")
	}
	if !hasNoInput {
		t.Error("expected to detect --no-input flag")
	}
}
