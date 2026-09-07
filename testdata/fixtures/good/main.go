// Package main implements a test fixture that passes all CLI-ACS conformance criteria.
package main

import (
	"fmt"
	"os"
)

func main() {
	noColor := os.Getenv("NO_COLOR") != ""
	termDumb := os.Getenv("TERM") == "dumb"

	// Parse flags
	for _, arg := range os.Args[1:] {
		switch arg {
		case "--help", "-h":
			printHelp(noColor || termDumb)
			os.Exit(0)
		case "--version":
			fmt.Println("fixture-good v1.0.0")
			os.Exit(0)
		case "--no-color":
			noColor = true
		case "--json":
			fmt.Println(`{"status":"ok"}`)
			os.Exit(0)
		case "--quiet", "-q":
			// Quiet mode: suppress output
			os.Exit(0)
		}
	}

	// If no valid command, show error
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Error: missing command. Run with --help for usage.")
		os.Exit(1)
	}

	// Unknown command
	fmt.Fprintln(os.Stderr, "Error: unknown command '"+os.Args[1]+"'. Run with --help for usage.")
	os.Exit(1)
}

func printHelp(plain bool) {
	if plain {
		// Plain text help (no ANSI color codes)
		fmt.Println("fixture-good - a test fixture for CLI-ACS")
		fmt.Println()
		fmt.Println("USAGE")
		fmt.Println("  fixture-good [flags]")
	} else {
		// Colored help with ANSI escape codes (using actual color codes)
		fmt.Println("\x1b[36mfixture-good\x1b[0m - a test fixture for CLI-ACS")
		fmt.Println()
		fmt.Println("\x1b[33mUSAGE\x1b[0m")
		fmt.Println("  fixture-good [flags]")
	}
	fmt.Println()
	fmt.Println("FLAGS")
	fmt.Println("  -h, --help       Show help")
	fmt.Println("      --version    Show version")
	fmt.Println("      --json       Output as JSON")
	fmt.Println("  -q, --quiet      Suppress output")
	fmt.Println("      --no-color   Disable color")
}
