// Package main implements a test fixture that fails CV-2 by ignoring NO_COLOR.
package main

import (
	"fmt"
	"os"
)

func main() {
	// Parse flags
	for _, arg := range os.Args[1:] {
		switch arg {
		case "--help", "-h":
			printHelp()
			os.Exit(0)
		case "--version":
			fmt.Println("fixture-no-nocolor v1.0.0")
			os.Exit(0)
		case "--json":
			fmt.Println(`{"status":"ok"}`)
			os.Exit(0)
		case "--quiet", "-q":
			os.Exit(0)
		}
	}

	// If no valid command, show error
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Error: missing command. Run with --help for usage.")
		os.Exit(1)
	}

	fmt.Fprintln(os.Stderr, "Error: unknown command '"+os.Args[1]+"'. Run with --help for usage.")
	os.Exit(1)
}

func printHelp() {
	// ALWAYS emit ANSI color codes, regardless of NO_COLOR or TERM
	// This fixture intentionally ignores NO_COLOR to fail CV-2
	// Using actual color codes (not just bold/reset)
	fmt.Println("\x1b[36mfixture-no-nocolor\x1b[0m - a test fixture that ignores NO_COLOR")
	fmt.Println()
	fmt.Println("\x1b[33mUSAGE\x1b[0m")
	fmt.Println("  fixture-no-nocolor [flags]")
	fmt.Println()
	fmt.Println("\x1b[33mFLAGS\x1b[0m")
	fmt.Println("  -h, --help       Show help")
	fmt.Println("      --version    Show version")
	fmt.Println("      --json       Output as JSON")
	fmt.Println("  -q, --quiet      Suppress output")
}
