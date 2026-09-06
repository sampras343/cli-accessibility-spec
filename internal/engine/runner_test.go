// internal/engine/runner_test.go
package engine

import "testing"

func TestComputeExitCode(t *testing.T) {
	tests := []struct {
		name      string
		results   []Result
		threshold Level
		want      int
	}{
		{
			name:      "all pass",
			results:   []Result{{Outcome: Supports, Level: LevelA}},
			threshold: LevelA,
			want:      0,
		},
		{
			name:      "level A failure",
			results:   []Result{{Outcome: DoesNotSupport, Level: LevelA}},
			threshold: LevelA,
			want:      1,
		},
		{
			name:      "level AA failure with AA threshold",
			results:   []Result{{Outcome: Supports, Level: LevelA}, {Outcome: DoesNotSupport, Level: LevelAA}},
			threshold: LevelAA,
			want:      2,
		},
		{
			name:      "level AA failure with A threshold passes",
			results:   []Result{{Outcome: Supports, Level: LevelA}, {Outcome: DoesNotSupport, Level: LevelAA}},
			threshold: LevelA,
			want:      0,
		},
		{
			name:      "level AAA failure with AAA threshold",
			results:   []Result{{Outcome: Supports, Level: LevelA}, {Outcome: Supports, Level: LevelAA}, {Outcome: DoesNotSupport, Level: LevelAAA}},
			threshold: LevelAAA,
			want:      3,
		},
		{
			name:      "partial support counts as failure",
			results:   []Result{{Outcome: PartiallySupports, Level: LevelA}},
			threshold: LevelA,
			want:      1,
		},
		{
			name:      "error counts as failure",
			results:   []Result{{Outcome: OutcomeError, Level: LevelA}},
			threshold: LevelA,
			want:      1,
		},
		{
			name:      "timeout counts as failure",
			results:   []Result{{Outcome: OutcomeTimeout, Level: LevelA}},
			threshold: LevelA,
			want:      1,
		},
		{
			name:      "not applicable does not count as failure",
			results:   []Result{{Outcome: NotApplicable, Level: LevelA}},
			threshold: LevelA,
			want:      0,
		},
		{
			name:      "not evaluated does not count as failure",
			results:   []Result{{Outcome: NotEvaluated, Level: LevelA}},
			threshold: LevelA,
			want:      0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report := &ConformanceReport{Results: tt.results, Threshold: tt.threshold.String()}
			got := ComputeExitCode(report, tt.threshold)
			if got != tt.want {
				t.Errorf("ComputeExitCode = %d, want %d", got, tt.want)
			}
		})
	}
}
