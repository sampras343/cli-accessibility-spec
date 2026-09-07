// Package main implements a test fixture that fails HD-1 by not providing --help.
package main

import (
	"fmt"
	"os"
)

func main() {
	// Parse flags - but --help and -h are NOT supported
	for _, arg := range os.Args[1:] {
		switch arg {
		case "--version":
			fmt.Println("fixture-no-help v1.0.0")
			os.Exit(0)
		case "--json":
			fmt.Println(`{"status":"ok"}`)
			os.Exit(0)
		case "--quiet", "-q":
			os.Exit(0)
		case "--help", "-h":
			// Intentionally fail on --help to test HD-1
			fmt.Fprintln(os.Stderr, "Error: unknown flag: "+arg)
			os.Exit(1)
		}
	}

	// If no valid command, show error
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Error: missing command")
		os.Exit(1)
	}

	fmt.Fprintln(os.Stderr, "Error: unknown command '"+os.Args[1]+"'")
	os.Exit(1)
}
