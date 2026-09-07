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
	// Create a temp dir with a valid YAML criterion for testing
	tmpDir := t.TempDir()
	yamlContent := `id: "TX-1"
domain: "test"
level: "A"
testability: "AUTO"
spec_version: "1.0"
name: "Test Criterion"
steps:
  - name: "test step"
    exec:
      args: ["--help"]
    assert:
      exit_code: 0
result_on_all_pass: "Supports"
result_on_any_fail: "Does Not Support"
`
	os.WriteFile(tmpDir+"/TX-1.yaml", []byte(yamlContent), 0644)
	tcContent := "---\ncriterion_id: TX-1\ndomain: test\nlevel: A\ntestability: AUTO\nspec_version: \"1.0\"\n---\n# TX-1\n## What This Tests\ntest\n## How The Test Works\ntest\n## Pass Criteria\ntest\n## Fail Criteria\ntest\n"
	os.WriteFile(tmpDir+"/TX-1.testcase.md", []byte(tcContent), 0644)

	os.Args = []string{"cli-acs", "validate", tmpDir}

	rootCmd := buildRootCommand()
	err := rootCmd.Execute()

	if err != nil {
		t.Fatalf("validate command failed: %v", err)
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
