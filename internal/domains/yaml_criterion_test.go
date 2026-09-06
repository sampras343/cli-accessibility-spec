package domains

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sampras343/cli-accessibility-spec/internal/engine"
)

func TestLoadYAMLCriteria(t *testing.T) {
	dir := filepath.Join("testdata")
	criteria, err := LoadYAMLCriteria(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(criteria) == 0 {
		t.Fatal("expected at least one criterion")
	}
	c := criteria[0]
	if c.ID() != "TEST-1" {
		t.Errorf("ID = %q, want %q", c.ID(), "TEST-1")
	}
	if c.Name() != "Sample Test Criterion" {
		t.Errorf("Name = %q, want %q", c.Name(), "Sample Test Criterion")
	}
	if c.Domain() != "test" {
		t.Errorf("Domain = %q, want %q", c.Domain(), "test")
	}
}

func TestYAMLCriterionRunPassesWithEcho(t *testing.T) {
	// Create a temp YAML that tests "echo" (which always exits 0)
	dir := t.TempDir()
	yaml := `id: "ECHO-1"
domain: "test"
level: "A"
testability: "AUTO"
spec_version: "1.0"
name: "Echo test"
steps:
  - name: "echo exits 0"
    exec:
      args: ["-c", "echo hello"]
      env: {}
    assert:
      exit_code: 0
      stdout_contains:
        - "hello"
result_on_all_pass: "Supports"
result_on_any_fail: "Does Not Support"
`
	os.WriteFile(filepath.Join(dir, "ECHO-1.yaml"), []byte(yaml), 0644)

	criteria, err := LoadYAMLCriteria(dir)
	if err != nil {
		t.Fatal(err)
	}

	probe := &ProbeResult{HasHelp: true}
	result := criteria[0].Run(context.Background(), "sh", probe)
	if result.Outcome != engine.Supports {
		t.Errorf("outcome = %v, want Supports. Remarks: %s", result.Outcome, result.Remarks)
	}
}

func TestYAMLCriterionAssertions(t *testing.T) {
	tests := []struct {
		name        string
		yaml        string
		expectPass  bool
		expectError string
	}{
		{
			name: "stdout_matches regex",
			yaml: `id: "TEST-REGEX"
domain: "test"
level: "A"
testability: "AUTO"
spec_version: "1.0"
name: "Test regex"
steps:
  - name: "test"
    exec:
      args: ["-c", "echo 'hello world'"]
    assert:
      stdout_matches:
        - "^hello.*world"
result_on_all_pass: "Supports"
result_on_any_fail: "Does Not Support"
`,
			expectPass: true,
		},
		{
			name: "stdout_not_contains",
			yaml: `id: "TEST-NOT-CONTAINS"
domain: "test"
level: "A"
testability: "AUTO"
spec_version: "1.0"
name: "Test not contains"
steps:
  - name: "test"
    exec:
      args: ["-c", "echo 'hello'"]
    assert:
      stdout_not_contains:
        - "goodbye"
result_on_all_pass: "Supports"
result_on_any_fail: "Does Not Support"
`,
			expectPass: true,
		},
		{
			name: "stderr_contains failure",
			yaml: `id: "TEST-STDERR"
domain: "test"
level: "A"
testability: "AUTO"
spec_version: "1.0"
name: "Test stderr"
steps:
  - name: "test"
    exec:
      args: ["-c", "echo hello"]
    assert:
      stderr_contains:
        - "error"
result_on_all_pass: "Supports"
result_on_any_fail: "Does Not Support"
`,
			expectPass:  false,
			expectError: "stderr_contains",
		},
		{
			name: "exit_code_nonzero",
			yaml: `id: "TEST-NONZERO"
domain: "test"
level: "A"
testability: "AUTO"
spec_version: "1.0"
name: "Test nonzero"
steps:
  - name: "test"
    exec:
      args: ["-c", "exit 1"]
    assert:
      exit_code_nonzero: true
result_on_all_pass: "Supports"
result_on_any_fail: "Does Not Support"
`,
			expectPass: true,
		},
		{
			name: "stdout_empty",
			yaml: `id: "TEST-EMPTY"
domain: "test"
level: "A"
testability: "AUTO"
spec_version: "1.0"
name: "Test empty"
steps:
  - name: "test"
    exec:
      args: ["-c", "true"]
    assert:
      stdout_empty: true
result_on_all_pass: "Supports"
result_on_any_fail: "Does Not Support"
`,
			expectPass: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			os.WriteFile(filepath.Join(dir, "test.yaml"), []byte(tt.yaml), 0644)

			criteria, err := LoadYAMLCriteria(dir)
			if err != nil {
				t.Fatal(err)
			}

			probe := &ProbeResult{HasHelp: true}
			result := criteria[0].Run(context.Background(), "sh", probe)

			if tt.expectPass {
				if result.Outcome != engine.Supports {
					t.Errorf("expected pass but got %v. Remarks: %s", result.Outcome, result.Remarks)
				}
			} else {
				if result.Outcome == engine.Supports {
					t.Errorf("expected failure but got Supports")
				}
				if tt.expectError != "" {
					found := strings.Contains(result.Remarks, tt.expectError)
					for _, ev := range result.Evidence {
						if strings.Contains(ev.Note, tt.expectError) {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("expected error containing %q, got Remarks: %s", tt.expectError, result.Remarks)
					}
				}
			}
		})
	}
}

func TestYAMLCriterionPreconditions(t *testing.T) {
	yaml := `id: "TEST-PRECOND"
domain: "test"
level: "A"
testability: "AUTO"
spec_version: "1.0"
name: "Test preconditions"
requires:
  - "has_help"
  - "has_version"
steps:
  - name: "test"
    exec:
      args: ["-c", "echo hello"]
    assert:
      exit_code: 0
result_on_all_pass: "Supports"
result_on_any_fail: "Does Not Support"
`
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "test.yaml"), []byte(yaml), 0644)

	criteria, err := LoadYAMLCriteria(dir)
	if err != nil {
		t.Fatal(err)
	}

	// Test with preconditions not met
	probe := &ProbeResult{HasHelp: false}
	result := criteria[0].Run(context.Background(), "sh", probe)
	if result.Outcome != engine.NotApplicable {
		t.Errorf("expected NotApplicable when preconditions not met, got %v", result.Outcome)
	}

	// Test with preconditions met
	probe = &ProbeResult{HasHelp: true, HasVersion: true}
	result = criteria[0].Run(context.Background(), "sh", probe)
	if result.Outcome != engine.Supports {
		t.Errorf("expected Supports when preconditions met, got %v. Remarks: %s", result.Outcome, result.Remarks)
	}
}

func TestYAMLCriterionLevelParsing(t *testing.T) {
	tests := []struct {
		level    string
		expected engine.Level
		wantErr  bool
	}{
		{"A", engine.LevelA, false},
		{"AA", engine.LevelAA, false},
		{"AAA", engine.LevelAAA, false},
		{"a", engine.LevelA, false},
		{"invalid", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.level, func(t *testing.T) {
			yaml := fmt.Sprintf(`id: "TEST"
domain: "test"
level: "%s"
testability: "AUTO"
spec_version: "1.0"
name: "Test"
steps:
  - name: "test"
    exec:
      args: ["-c", "true"]
    assert:
      exit_code: 0
result_on_all_pass: "Supports"
result_on_any_fail: "Does Not Support"
`, tt.level)
			dir := t.TempDir()
			os.WriteFile(filepath.Join(dir, "test.yaml"), []byte(yaml), 0644)

			criteria, err := LoadYAMLCriteria(dir)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if criteria[0].Level() != tt.expected {
				t.Errorf("level = %v, want %v", criteria[0].Level(), tt.expected)
			}
		})
	}
}
