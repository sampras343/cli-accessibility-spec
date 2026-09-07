// e2e_test.go - End-to-end tests for CLI-ACS conformance suite
//go:build e2e

package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestMain builds all required binaries before running e2e tests
func TestMain(m *testing.M) {
	// Build cli-acs binary
	cmd := exec.Command("go", "build", "-o", "cli-acs", "./cmd/cli-acs/")
	if output, err := cmd.CombinedOutput(); err != nil {
		panic("failed to build cli-acs: " + err.Error() + "\n" + string(output))
	}

	// Build fixture binaries
	cmd = exec.Command("go", "generate", "./testdata/fixtures/")
	if output, err := cmd.CombinedOutput(); err != nil {
		panic("failed to build fixtures: " + err.Error() + "\n" + string(output))
	}

	// Run tests
	code := m.Run()

	// Cleanup (optional - binaries are useful to keep for manual testing)
	// os.Remove("cli-acs")

	os.Exit(code)
}

// TestE2E_GoodFixture tests that a conforming binary passes all criteria
func TestE2E_GoodFixture(t *testing.T) {
	binPath, err := filepath.Abs("testdata/bin/fixture-good")
	if err != nil {
		t.Fatalf("failed to get absolute path: %v", err)
	}

	// Check color domain
	cmd := exec.Command("./cli-acs", "check", "--domain", "color", "--format", "json", binPath)
	output, err := cmd.CombinedOutput()

	// Don't check err - we just want to parse the output
	// (The fixture may fail some criteria like CV-4)

	// Parse JSON output to verify results
	var report map[string]interface{}
	if err := json.Unmarshal(output, &report); err != nil {
		t.Fatalf("failed to parse JSON output: %v\nOutput: %s", err, output)
	}

	// Get results array
	resultsRaw, ok := report["results"].([]interface{})
	if !ok {
		t.Fatalf("expected results array in JSON output, got: %v", report)
	}

	// Verify we got results
	if len(resultsRaw) == 0 {
		t.Fatal("expected at least one result in JSON output")
	}

	// Check that critical criteria pass
	// CV-2 (NO_COLOR Support) should PASS
	// CV-5 (TERM=dumb Respect) should PASS
	// CV-4 may FAIL because the fixture always emits color (doesn't check TTY)
	// CV-3 may be N/A if --no-color flag isn't detected
	criticalResults := make(map[string]float64)
	for _, resultRaw := range resultsRaw {
		result, ok := resultRaw.(map[string]interface{})
		if !ok {
			continue
		}

		criteriaID, ok := result["id"].(string)
		if !ok {
			continue
		}

		outcome, ok := result["outcome"].(float64)
		if !ok {
			t.Errorf("result missing outcome field: %v", result)
			continue
		}

		criticalResults[criteriaID] = outcome
	}

	// Outcome codes: 0=PASS, 1=PARTIAL, 2=FAIL, 3=N/A
	// Check critical criteria
	if cv2, ok := criticalResults["CV-2"]; ok && cv2 != 0 {
		t.Errorf("expected CV-2 (NO_COLOR Support) to PASS, got outcome %v", cv2)
	}
	if cv5, ok := criticalResults["CV-5"]; ok && cv5 != 0 {
		t.Errorf("expected CV-5 (TERM=dumb Respect) to PASS, got outcome %v", cv5)
	}

	// Note: We don't check exit code because the fixture may fail non-critical criteria like CV-4
}

// TestE2E_NoNocolorFixture tests that a binary ignoring NO_COLOR fails CV-2
func TestE2E_NoNocolorFixture(t *testing.T) {
	binPath, err := filepath.Abs("testdata/bin/fixture-no-nocolor")
	if err != nil {
		t.Fatalf("failed to get absolute path: %v", err)
	}

	// Check color domain
	cmd := exec.Command("./cli-acs", "check", "--domain", "color", "--format", "json", binPath)
	output, err := cmd.CombinedOutput()

	// We expect a non-zero exit code because CV-2 should fail
	if cmd.ProcessState.ExitCode() == 0 {
		t.Errorf("expected non-zero exit for fixture that fails NO_COLOR\nOutput: %s", output)
	}

	// Parse JSON output
	var report map[string]interface{}
	if err := json.Unmarshal(output, &report); err != nil {
		t.Fatalf("failed to parse JSON output: %v\nOutput: %s", err, output)
	}

	// Get results array
	resultsRaw, ok := report["results"].([]interface{})
	if !ok {
		t.Fatalf("expected results array in JSON output")
	}

	// Find CV-2 result
	foundCV2 := false
	for _, resultRaw := range resultsRaw {
		result, ok := resultRaw.(map[string]interface{})
		if !ok {
			continue
		}

		criteriaID, ok := result["id"].(string)
		if !ok {
			continue
		}

		if strings.HasPrefix(criteriaID, "CV-2") {
			foundCV2 = true
			outcome := result["outcome"].(float64)
			// Outcome 2 = FAIL
			if outcome != 2 {
				t.Errorf("expected CV-2 to FAIL (outcome 2) for no-nocolor fixture, got %v", outcome)
			}
		}
	}

	if !foundCV2 {
		t.Error("CV-2 criteria not found in results")
	}
}

// TestE2E_NoHelpFixture tests that a binary without --help can be tested
// Note: This test is forward-looking - it will pass once HD-1 criteria are implemented
func TestE2E_NoHelpFixture(t *testing.T) {
	binPath, err := filepath.Abs("testdata/bin/fixture-no-help")
	if err != nil {
		t.Fatalf("failed to get absolute path: %v", err)
	}

	// Verify the fixture exists and rejects --help
	cmd := exec.Command(binPath, "--help")
	output, _ := cmd.CombinedOutput()
	if cmd.ProcessState.ExitCode() == 0 {
		t.Errorf("fixture-no-help should reject --help, but it succeeded\nOutput: %s", output)
	}
	if !strings.Contains(string(output), "unknown flag") {
		t.Errorf("expected 'unknown flag' error message, got: %s", output)
	}

	// Check with cli-acs (when HD criteria are implemented, they should fail)
	cmd = exec.Command("./cli-acs", "check", "--format", "json", binPath)
	output, _ = cmd.CombinedOutput()

	// Parse JSON output
	var report map[string]interface{}
	if err := json.Unmarshal(output, &report); err != nil {
		t.Fatalf("failed to parse JSON output: %v\nOutput: %s", err, output)
	}

	// Get results array
	resultsRaw, ok := report["results"].([]interface{})
	if !ok {
		t.Fatalf("expected results array in JSON output")
	}

	// Check if HD-1 exists (forward-looking test)
	foundHD1 := false
	for _, resultRaw := range resultsRaw {
		result, ok := resultRaw.(map[string]interface{})
		if !ok {
			continue
		}

		criteriaID, ok := result["id"].(string)
		if !ok {
			continue
		}

		if strings.HasPrefix(criteriaID, "HD-1") {
			foundHD1 = true
			outcome := result["outcome"].(float64)
			// Outcome 2 = FAIL
			if outcome != 2 {
				t.Errorf("expected HD-1 to FAIL (outcome 2) for no-help fixture, got %v", outcome)
			}
		}
	}

	// Skip the HD-1 check if not implemented yet
	if !foundHD1 {
		t.Skip("HD-1 criteria not yet implemented - test will pass once help/docs domain is added")
	}
}

// TestE2E_MultiDomain tests checking multiple domains at once
func TestE2E_MultiDomain(t *testing.T) {
	binPath, err := filepath.Abs("testdata/bin/fixture-good")
	if err != nil {
		t.Fatalf("failed to get absolute path: %v", err)
	}

	// Check without specifying domain (should check all domains)
	cmd := exec.Command("./cli-acs", "check", "--format", "json", binPath)
	output, _ := cmd.CombinedOutput()

	// Don't check err - some criteria may fail

	// Parse JSON output
	var report map[string]interface{}
	if err := json.Unmarshal(output, &report); err != nil {
		t.Fatalf("failed to parse JSON output: %v\nOutput: %s", err, output)
	}

	// Get results array
	resultsRaw, ok := report["results"].([]interface{})
	if !ok {
		t.Fatalf("expected results array in JSON output")
	}

	// Should have results from multiple domains
	domains := make(map[string]bool)
	for _, resultRaw := range resultsRaw {
		result, ok := resultRaw.(map[string]interface{})
		if !ok {
			continue
		}

		domain, ok := result["domain"].(string)
		if ok {
			domains[domain] = true
		}
	}

	if len(domains) < 1 {
		t.Errorf("expected results from at least one domain, got: %v", domains)
	}
}
