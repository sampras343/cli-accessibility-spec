// internal/engine/types_test.go
package engine

import "testing"

func TestLevelString(t *testing.T) {
	tests := []struct {
		level Level
		want  string
	}{
		{LevelA, "A"},
		{LevelAA, "AA"},
		{LevelAAA, "AAA"},
	}
	for _, tt := range tests {
		if got := tt.level.String(); got != tt.want {
			t.Errorf("Level(%d).String() = %q, want %q", tt.level, got, tt.want)
		}
	}
}

func TestOutcomeString(t *testing.T) {
	tests := []struct {
		outcome Outcome
		want    string
	}{
		{Supports, "Supports"},
		{PartiallySupports, "Partially Supports"},
		{DoesNotSupport, "Does Not Support"},
		{NotApplicable, "Not Applicable"},
		{NotEvaluated, "Not Evaluated"},
		{OutcomeError, "Error"},
		{OutcomeTimeout, "Timeout"},
	}
	for _, tt := range tests {
		if got := tt.outcome.String(); got != tt.want {
			t.Errorf("Outcome(%d).String() = %q, want %q", tt.outcome, got, tt.want)
		}
	}
}

func TestTestabilityString(t *testing.T) {
	tests := []struct {
		t    Testability
		want string
	}{
		{Auto, "AUTO"},
		{Semi, "SEMI"},
		{Manual, "MANUAL"},
	}
	for _, tt := range tests {
		if got := tt.t.String(); got != tt.want {
			t.Errorf("Testability(%d).String() = %q, want %q", tt.t, got, tt.want)
		}
	}
}

func TestResultOutcomeIsFailure(t *testing.T) {
	passing := []Outcome{Supports, NotApplicable, NotEvaluated}
	for _, o := range passing {
		r := Result{Outcome: o}
		if r.IsFailure() {
			t.Errorf("Outcome %v should not be a failure", o)
		}
	}
	failing := []Outcome{DoesNotSupport, PartiallySupports, OutcomeError, OutcomeTimeout}
	for _, o := range failing {
		r := Result{Outcome: o}
		if !r.IsFailure() {
			t.Errorf("Outcome %v should be a failure", o)
		}
	}
}
