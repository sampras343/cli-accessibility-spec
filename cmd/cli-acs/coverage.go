// cmd/cli-acs/coverage.go
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/sampras343/cli-accessibility-spec/internal/domains"
	_ "github.com/sampras343/cli-accessibility-spec/internal/domains/color"
	"github.com/spf13/cobra"
)

func newCoverageCommand() *cobra.Command {
	var specFile string

	cmd := &cobra.Command{
		Use:   "coverage",
		Short: "Check coverage of spec criteria by tests and docs",
		Long: `Check coverage of spec criteria by tests and documentation.

This command:
- Parses the spec file for all criterion IDs (e.g., OS-1, CV-2, HD-3)
- Scans domain testcases directories for YAML tests and testcase.md files
- Checks registered Go-based criteria from each domain
- Reports:
  - Missing tests: criterion IDs in spec with no YAML or Go check
  - Missing docs: criterion IDs with no testcase.md
  - Orphan tests: YAML/Go checks for IDs not in spec
  - Orphan docs: testcase.md for IDs not in spec

Exit code 0 if all clean, exit code 1 if gaps or orphans found.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCoverage(specFile)
		},
	}

	cmd.Flags().StringVar(&specFile, "spec", "spec/CLI_ACS_v1.0.md", "Path to spec file")

	return cmd
}

type coverageReport struct {
	specIDs      map[string]bool
	testedIDs    map[string]string // ID -> source (YAML file or "Go")
	documentedIDs map[string]string // ID -> testcase.md file
}

func runCoverage(specFile string) error {
	report := &coverageReport{
		specIDs:       make(map[string]bool),
		testedIDs:     make(map[string]string),
		documentedIDs: make(map[string]string),
	}

	// Step 1: Parse spec file for criterion IDs
	if err := parseSpecFile(specFile, report); err != nil {
		return fmt.Errorf("parse spec file: %w", err)
	}

	// Step 2: Scan for YAML tests and testcase.md files
	if err := scanTestcases(report); err != nil {
		return fmt.Errorf("scan testcases: %w", err)
	}

	// Step 3: Check registered Go criteria
	if err := checkRegisteredCriteria(report); err != nil {
		return fmt.Errorf("check registered criteria: %w", err)
	}

	// Step 4: Generate report
	hasIssues := generateReport(report)

	if hasIssues {
		return fmt.Errorf("coverage gaps or orphans found")
	}

	fmt.Println("\nCoverage check passed: all criteria have tests and documentation.")
	return nil
}

func parseSpecFile(specFile string, report *coverageReport) error {
	data, err := os.ReadFile(specFile)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	// Regex to match criterion IDs like OS-1, CV-2, HD-13, etc.
	// Pattern: 2-3 uppercase letters, hyphen, 1+ digits
	re := regexp.MustCompile(`\b([A-Z]{2,3})-(\d+)\b`)
	matches := re.FindAllStringSubmatch(string(data), -1)

	for _, match := range matches {
		criterionID := match[0] // Full match: "OS-1"
		report.specIDs[criterionID] = true
	}

	fmt.Printf("Found %d criterion IDs in spec\n", len(report.specIDs))
	return nil
}

func scanTestcases(report *coverageReport) error {
	domainsDir := "internal/domains"

	// Find all domain subdirectories
	entries, err := os.ReadDir(domainsDir)
	if err != nil {
		return fmt.Errorf("read domains dir: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		domainName := entry.Name()
		testcasesDir := filepath.Join(domainsDir, domainName, "testcases")

		// Check if testcases directory exists
		if _, err := os.Stat(testcasesDir); os.IsNotExist(err) {
			continue
		}

		// Scan testcases directory
		if err := scanTestcasesDir(testcasesDir, report); err != nil {
			return fmt.Errorf("scan %s: %w", testcasesDir, err)
		}
	}

	fmt.Printf("Found %d tested IDs (YAML)\n", len(report.testedIDs))
	fmt.Printf("Found %d documented IDs (testcase.md)\n", len(report.documentedIDs))
	return nil
}

func scanTestcasesDir(dir string, report *coverageReport) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read dir: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()

		// Check for YAML files
		if strings.HasSuffix(filename, ".yaml") || strings.HasSuffix(filename, ".yml") {
			// Extract criterion ID from filename (e.g., "CV-2.yaml" -> "CV-2")
			criterionID := strings.TrimSuffix(filename, ".yaml")
			criterionID = strings.TrimSuffix(criterionID, ".yml")

			// Validate ID format
			if matched, _ := regexp.MatchString(`^[A-Z]{2,3}-\d+$`, criterionID); matched {
				yamlPath := filepath.Join(dir, filename)
				report.testedIDs[criterionID] = yamlPath
			}
		}

		// Check for testcase.md files
		if strings.HasSuffix(filename, ".testcase.md") {
			// Extract criterion ID from filename (e.g., "CV-2.testcase.md" -> "CV-2")
			criterionID := strings.TrimSuffix(filename, ".testcase.md")

			// Validate ID format
			if matched, _ := regexp.MatchString(`^[A-Z]{2,3}-\d+$`, criterionID); matched {
				mdPath := filepath.Join(dir, filename)
				report.documentedIDs[criterionID] = mdPath
			}
		}
	}

	return nil
}

func checkRegisteredCriteria(report *coverageReport) error {
	// Get all registered domains
	allDomains := domains.AllDomains()

	for _, domain := range allDomains {
		criteria := domain.Criteria()
		for _, criterion := range criteria {
			id := criterion.ID()
			// Only mark as tested if not already found in YAML
			if _, exists := report.testedIDs[id]; !exists {
				report.testedIDs[id] = "Go (domain:" + domain.Name() + ")"
			}
		}
	}

	return nil
}

func generateReport(report *coverageReport) bool {
	hasIssues := false

	// Find missing tests
	var missingTests []string
	for id := range report.specIDs {
		if _, tested := report.testedIDs[id]; !tested {
			missingTests = append(missingTests, id)
		}
	}

	// Find missing docs
	var missingDocs []string
	for id := range report.specIDs {
		if _, documented := report.documentedIDs[id]; !documented {
			missingDocs = append(missingDocs, id)
		}
	}

	// Find orphan tests
	var orphanTests []string
	for id := range report.testedIDs {
		if _, inSpec := report.specIDs[id]; !inSpec {
			orphanTests = append(orphanTests, id)
		}
	}

	// Find orphan docs
	var orphanDocs []string
	for id := range report.documentedIDs {
		if _, inSpec := report.specIDs[id]; !inSpec {
			orphanDocs = append(orphanDocs, id)
		}
	}

	// Sort for consistent output
	sort.Strings(missingTests)
	sort.Strings(missingDocs)
	sort.Strings(orphanTests)
	sort.Strings(orphanDocs)

	// Print report
	fmt.Println("\n=== Coverage Report ===")

	if len(missingTests) > 0 {
		hasIssues = true
		fmt.Printf("\nMissing Tests (%d):\n", len(missingTests))
		for _, id := range missingTests {
			fmt.Printf("  - %s (in spec, no YAML or Go test)\n", id)
		}
	}

	if len(missingDocs) > 0 {
		hasIssues = true
		fmt.Printf("\nMissing Documentation (%d):\n", len(missingDocs))
		for _, id := range missingDocs {
			fmt.Printf("  - %s (in spec, no testcase.md)\n", id)
		}
	}

	if len(orphanTests) > 0 {
		hasIssues = true
		fmt.Printf("\nOrphan Tests (%d):\n", len(orphanTests))
		for _, id := range orphanTests {
			source := report.testedIDs[id]
			fmt.Printf("  - %s (has test at %s, not in spec)\n", id, source)
		}
	}

	if len(orphanDocs) > 0 {
		hasIssues = true
		fmt.Printf("\nOrphan Documentation (%d):\n", len(orphanDocs))
		for _, id := range orphanDocs {
			source := report.documentedIDs[id]
			fmt.Printf("  - %s (has doc at %s, not in spec)\n", id, source)
		}
	}

	if !hasIssues {
		fmt.Println("\nAll checks passed:")
		fmt.Printf("  - %d criteria in spec\n", len(report.specIDs))
		fmt.Printf("  - %d criteria with tests\n", len(report.testedIDs))
		fmt.Printf("  - %d criteria with documentation\n", len(report.documentedIDs))
	}

	return hasIssues
}
