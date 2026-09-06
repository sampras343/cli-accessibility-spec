// cmd/cli-acs/version.go
package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	suiteVersion = "1.0.0"
	specVersion  = "1.0.0"
)

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Long:  "Print the suite version and CLI-ACS specification version.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("cli-acs suite version: %s\n", suiteVersion)
			fmt.Printf("CLI-ACS spec version: %s\n", specVersion)
		},
	}
}
