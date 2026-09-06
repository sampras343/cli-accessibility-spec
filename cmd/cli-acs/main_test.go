// cmd/cli-acs/main_test.go
package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestVersionCommand(t *testing.T) {
	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Override os.Args to simulate 'cli-acs version'
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"cli-acs", "version"}

	// Reset cobra command state
	rootCmd := buildRootCommand()

	// Execute command
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("version command failed: %v", err)
	}

	// Restore stdout and read captured output
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Verify output
	if !strings.Contains(output, "cli-acs suite version") {
		t.Errorf("Expected suite version in output, got: %s", output)
	}
	if !strings.Contains(output, "CLI-ACS spec version") {
		t.Errorf("Expected spec version in output, got: %s", output)
	}
}

func TestCheckCommandHelp(t *testing.T) {
	// Override os.Args to simulate 'cli-acs check --help'
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"cli-acs", "check", "--help"}

	// Reset cobra command state
	rootCmd := buildRootCommand()

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Execute command (help doesn't return error)
	err := rootCmd.Execute()

	// Restore stdout
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	if err != nil {
		t.Fatalf("check --help failed: %v", err)
	}

	// Verify key flags are present
	expectedFlags := []string{
		"--domain",
		"--level",
		"--auto-only",
		"--threshold",
		"--format",
		"--output",
		"--crosswalk",
		"--plain",
		"--quiet",
		"--skip",
		"--timeout",
		"--max-subcommands",
		"--criteria-dir",
	}

	for _, flag := range expectedFlags {
		if !strings.Contains(output, flag) {
			t.Errorf("Expected flag %s in help output", flag)
		}
	}
}

func TestValidateCommandStub(t *testing.T) {
	// Override os.Args to simulate 'cli-acs validate'
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"cli-acs", "validate"}

	// Reset cobra command state
	rootCmd := buildRootCommand()

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Execute command
	err := rootCmd.Execute()

	// Restore stdout
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	if err != nil {
		t.Fatalf("validate command failed: %v", err)
	}

	// Verify it's a stub
	if !strings.Contains(output, "not yet implemented") {
		t.Errorf("Expected stub message in output, got: %s", output)
	}
}

// buildRootCommand builds a fresh root command for testing
func buildRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "cli-acs",
		Short: "CLI Accessibility Conformance Suite",
	}
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable color output")
	rootCmd.AddCommand(newCheckCommand())
	rootCmd.AddCommand(newValidateCommand())
	rootCmd.AddCommand(newVersionCommand())
	return rootCmd
}
