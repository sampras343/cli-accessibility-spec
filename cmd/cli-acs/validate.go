// cmd/cli-acs/validate.go
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func newValidateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "validate [path]",
		Short: "Validate YAML criteria and testcase.md files",
		Long: `Validate YAML criteria and testcase.md files for correctness.

This command checks:
- YAML files parse correctly and have required fields
- Each YAML file has a corresponding testcase.md file
- testcase.md frontmatter has correct criterion_id matching filename
- testcase.md has required sections (What, How, Pass, Fail)

If no path is provided, validates all testcases in internal/domains/*/testcases/.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			paths := []string{"internal/domains"}
			if len(args) > 0 {
				paths = []string{args[0]}
			}
			return runValidate(paths)
		},
	}
}

type validationError struct {
	file    string
	message string
}

type validationReport struct {
	errors []validationError
}

func (r *validationReport) addError(file, message string) {
	r.errors = append(r.errors, validationError{file: file, message: message})
}

func (r *validationReport) hasErrors() bool {
	return len(r.errors) > 0
}

func (r *validationReport) print() {
	if !r.hasErrors() {
		fmt.Println("Validation passed: all YAML and testcase.md files are valid.")
		return
	}

	fmt.Printf("\nValidation failed with %d error(s):\n\n", len(r.errors))
	for _, err := range r.errors {
		fmt.Printf("  %s:\n    %s\n", err.file, err.message)
	}
}

func runValidate(paths []string) error {
	report := &validationReport{}

	for _, path := range paths {
		if err := validatePath(path, report); err != nil {
			return err
		}
	}

	report.print()

	if report.hasErrors() {
		return fmt.Errorf("validation errors found")
	}

	return nil
}

func validatePath(path string, report *validationReport) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat %s: %w", path, err)
	}

	if info.IsDir() {
		// If it's a domains directory, scan for testcases subdirectories
		if strings.Contains(path, "domains") {
			return validateDomainsDir(path, report)
		}
		// Otherwise, treat as a single testcases directory
		return validateTestcasesDir(path, report)
	}

	return fmt.Errorf("%s is not a directory", path)
}

func validateDomainsDir(domainsDir string, report *validationReport) error {
	entries, err := os.ReadDir(domainsDir)
	if err != nil {
		return fmt.Errorf("read dir %s: %w", domainsDir, err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		testcasesDir := filepath.Join(domainsDir, entry.Name(), "testcases")
		if _, err := os.Stat(testcasesDir); os.IsNotExist(err) {
			continue
		}

		if err := validateTestcasesDir(testcasesDir, report); err != nil {
			return err
		}
	}

	return nil
}

func validateTestcasesDir(dir string, report *validationReport) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read dir %s: %w", dir, err)
	}

	yamlFiles := make(map[string]string)     // criterionID -> filepath
	testcaseMDs := make(map[string]string)   // criterionID -> filepath

	// First pass: collect all files
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		filepath := filepath.Join(dir, filename)

		// YAML files
		if strings.HasSuffix(filename, ".yaml") || strings.HasSuffix(filename, ".yml") {
			criterionID := strings.TrimSuffix(filename, ".yaml")
			criterionID = strings.TrimSuffix(criterionID, ".yml")

			if matched, _ := regexp.MatchString(`^[A-Z]{2,3}-\d+$`, criterionID); matched {
				yamlFiles[criterionID] = filepath
			} else {
				report.addError(filepath, "filename does not match criterion ID pattern [A-Z]{2,3}-\\d+")
			}
		}

		// Testcase markdown files
		if strings.HasSuffix(filename, ".testcase.md") {
			criterionID := strings.TrimSuffix(filename, ".testcase.md")

			if matched, _ := regexp.MatchString(`^[A-Z]{2,3}-\d+$`, criterionID); matched {
				testcaseMDs[criterionID] = filepath
			} else {
				report.addError(filepath, "filename does not match criterion ID pattern [A-Z]{2,3}-\\d+")
			}
		}
	}

	// Second pass: validate each YAML file
	for criterionID, yamlPath := range yamlFiles {
		validateYAMLFile(yamlPath, criterionID, report)

		// Check for corresponding testcase.md
		if _, hasMD := testcaseMDs[criterionID]; !hasMD {
			report.addError(yamlPath, fmt.Sprintf("missing corresponding %s.testcase.md file", criterionID))
		}
	}

	// Third pass: validate each testcase.md file
	for criterionID, mdPath := range testcaseMDs {
		validateTestcaseMD(mdPath, criterionID, report)
	}

	return nil
}

type yamlCriterionFile struct {
	ID              string   `yaml:"id"`
	Domain          string   `yaml:"domain"`
	Level           string   `yaml:"level"`
	Testability     string   `yaml:"testability"`
	SpecVersion     string   `yaml:"spec_version"`
	Name            string   `yaml:"name"`
	ResultOnAllPass string   `yaml:"result_on_all_pass"`
	ResultOnAnyFail string   `yaml:"result_on_any_fail"`
	Steps           []interface{} `yaml:"steps"`
}

func validateYAMLFile(path, expectedID string, report *validationReport) {
	data, err := os.ReadFile(path)
	if err != nil {
		report.addError(path, fmt.Sprintf("failed to read file: %v", err))
		return
	}

	var yf yamlCriterionFile
	if err := yaml.Unmarshal(data, &yf); err != nil {
		report.addError(path, fmt.Sprintf("failed to parse YAML: %v", err))
		return
	}

	// Validate required fields
	if yf.ID == "" {
		report.addError(path, "missing required field: id")
	} else if yf.ID != expectedID {
		report.addError(path, fmt.Sprintf("id mismatch: filename says %q but YAML has %q", expectedID, yf.ID))
	}

	if yf.Domain == "" {
		report.addError(path, "missing required field: domain")
	}

	if yf.Level == "" {
		report.addError(path, "missing required field: level")
	} else {
		validLevels := map[string]bool{"A": true, "AA": true, "AAA": true}
		if !validLevels[strings.ToUpper(yf.Level)] {
			report.addError(path, fmt.Sprintf("invalid level %q (must be A, AA, or AAA)", yf.Level))
		}
	}

	if yf.Testability == "" {
		report.addError(path, "missing required field: testability")
	} else {
		validTestability := map[string]bool{"AUTO": true, "SEMI": true, "MANUAL": true}
		if !validTestability[strings.ToUpper(yf.Testability)] {
			report.addError(path, fmt.Sprintf("invalid testability %q (must be AUTO, SEMI, or MANUAL)", yf.Testability))
		}
	}

	if yf.SpecVersion == "" {
		report.addError(path, "missing required field: spec_version")
	}

	if yf.Name == "" {
		report.addError(path, "missing required field: name")
	}

	if yf.ResultOnAllPass == "" {
		report.addError(path, "missing required field: result_on_all_pass")
	} else {
		validOutcomes := map[string]bool{
			"Supports": true, "Partially Supports": true, "Does Not Support": true,
			"Not Applicable": true, "Not Evaluated": true, "Error": true, "Timeout": true,
		}
		if !validOutcomes[yf.ResultOnAllPass] {
			report.addError(path, fmt.Sprintf("invalid result_on_all_pass %q", yf.ResultOnAllPass))
		}
	}

	if yf.ResultOnAnyFail == "" {
		report.addError(path, "missing required field: result_on_any_fail")
	} else {
		validOutcomes := map[string]bool{
			"Supports": true, "Partially Supports": true, "Does Not Support": true,
			"Not Applicable": true, "Not Evaluated": true, "Error": true, "Timeout": true,
		}
		if !validOutcomes[yf.ResultOnAnyFail] {
			report.addError(path, fmt.Sprintf("invalid result_on_any_fail %q", yf.ResultOnAnyFail))
		}
	}

	if len(yf.Steps) == 0 {
		report.addError(path, "missing or empty field: steps")
	}
}

type testcaseMDFrontmatter struct {
	CriterionID string `yaml:"criterion_id"`
	Domain      string `yaml:"domain"`
	Level       string `yaml:"level"`
	Testability string `yaml:"testability"`
	SpecVersion string `yaml:"spec_version"`
	Type        string `yaml:"type"`
}

func validateTestcaseMD(path, expectedID string, report *validationReport) {
	data, err := os.ReadFile(path)
	if err != nil {
		report.addError(path, fmt.Sprintf("failed to read file: %v", err))
		return
	}

	content := string(data)

	// Extract frontmatter
	if !strings.HasPrefix(content, "---\n") {
		report.addError(path, "missing frontmatter (should start with ---)")
		return
	}

	// Find end of frontmatter
	endIndex := strings.Index(content[4:], "\n---\n")
	if endIndex == -1 {
		report.addError(path, "malformed frontmatter (missing closing ---)")
		return
	}

	frontmatterStr := content[4 : 4+endIndex]
	bodyContent := content[4+endIndex+5:]

	// Parse frontmatter
	var fm testcaseMDFrontmatter
	if err := yaml.Unmarshal([]byte(frontmatterStr), &fm); err != nil {
		report.addError(path, fmt.Sprintf("failed to parse frontmatter: %v", err))
		return
	}

	// Validate frontmatter fields
	if fm.CriterionID == "" {
		report.addError(path, "frontmatter missing required field: criterion_id")
	} else if fm.CriterionID != expectedID {
		report.addError(path, fmt.Sprintf("criterion_id mismatch: filename says %q but frontmatter has %q", expectedID, fm.CriterionID))
	}

	if fm.Domain == "" {
		report.addError(path, "frontmatter missing required field: domain")
	}

	if fm.Level == "" {
		report.addError(path, "frontmatter missing required field: level")
	}

	if fm.Testability == "" {
		report.addError(path, "frontmatter missing required field: testability")
	}

	if fm.SpecVersion == "" {
		report.addError(path, "frontmatter missing required field: spec_version")
	}

	// Check for required sections in body
	requiredSections := []string{
		"## What This Tests",
		"## How The Test Works",
		"## Pass Criteria",
		"## Fail Criteria",
	}

	for _, section := range requiredSections {
		if !strings.Contains(bodyContent, section) {
			report.addError(path, fmt.Sprintf("missing required section: %s", section))
		}
	}
}

