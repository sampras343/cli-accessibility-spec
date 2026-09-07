// internal/engine/config.go
package engine

import "time"

// Config holds all configuration for the conformance runner
type Config struct {
	// DomainFilter restricts testing to specific domains (empty = all domains)
	DomainFilter []string

	// LevelFilter restricts testing to specific levels (empty = all levels)
	LevelFilter []Level

	// AutoOnly restricts testing to AUTO testability criteria only
	AutoOnly bool

	// Threshold is the minimum level that must pass for the test suite to succeed
	Threshold Level

	// SkipCriteria is a list of criterion IDs to skip during execution
	SkipCriteria []string

	// Timeout is the maximum duration for each criterion execution
	Timeout time.Duration

	// MaxSubcommands limits the number of subcommands probed
	MaxSubcommands int

	// CriteriaDir is the directory to load external criteria from (if applicable)
	CriteriaDir string

	// Format specifies the output format (json, yaml, etc.)
	Format string

	// OutputPath is the file path to write the report to
	OutputPath string

	// Quiet suppresses progress output
	Quiet bool

	// Plain disables color output
	Plain bool

	// Crosswalk enables crosswalk mapping output
	Crosswalk bool

	// NoColor disables color output (alias for Plain)
	NoColor bool
}
