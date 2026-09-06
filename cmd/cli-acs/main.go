// cmd/cli-acs/main.go
package main

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	// Global noColor flag
	noColor bool
)

func main() {
	// Self-dogfooding: check NO_COLOR env var and TERM=dumb at startup
	if os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb" {
		noColor = true
	}

	rootCmd := &cobra.Command{
		Use:   "cli-acs",
		Short: "CLI Accessibility Conformance Suite",
		Long: `cli-acs is a conformance suite for evaluating CLI tools against
the CLI Accessibility Specification (CLI-ACS) v1.0.

It performs automated, semi-automated, and manual checks across
multiple domains including output, color, help, errors, input,
environment handling, timing, internationalization, and lifecycle.`,
	}

	// Global flags
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", noColor, "Disable color output")

	// Add subcommands
	rootCmd.AddCommand(newCheckCommand())
	rootCmd.AddCommand(newCoverageCommand())
	rootCmd.AddCommand(newValidateCommand())
	rootCmd.AddCommand(newVersionCommand())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(4) // Suite internal error
	}
}
