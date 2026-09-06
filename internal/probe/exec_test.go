// internal/probe/exec_test.go
package probe

import (
	"context"
	"testing"
	"time"
)

func TestRunCapturesStdoutStderr(t *testing.T) {
	result, err := Run(context.Background(), "sh", ExecOpts{
		Args: []string{"-c", "echo hello; echo err >&2"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if string(result.Stdout) != "hello\n" {
		t.Errorf("stdout = %q, want %q", result.Stdout, "hello\n")
	}
	if string(result.Stderr) != "err\n" {
		t.Errorf("stderr = %q, want %q", result.Stderr, "err\n")
	}
	if result.ExitCode != 0 {
		t.Errorf("exit code = %d, want 0", result.ExitCode)
	}
}

func TestRunCapturesNonZeroExit(t *testing.T) {
	result, err := Run(context.Background(), "sh", ExecOpts{
		Args: []string{"-c", "exit 42"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.ExitCode != 42 {
		t.Errorf("exit code = %d, want 42", result.ExitCode)
	}
}

func TestRunTimesOut(t *testing.T) {
	result, err := Run(context.Background(), "sleep", ExecOpts{
		Args:    []string{"60"},
		Timeout: 500 * time.Millisecond,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.ExitCode == 0 {
		t.Error("expected non-zero exit on timeout")
	}
	if result.TimedOut != true {
		t.Error("expected TimedOut = true")
	}
}

func TestRunIsolatesEnv(t *testing.T) {
	result, err := Run(context.Background(), "sh", ExecOpts{
		Args: []string{"-c", "echo $HOME"},
		Env:  map[string]string{"HOME": "/tmp/fake"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if string(result.Stdout) != "/tmp/fake\n" {
		t.Errorf("HOME not isolated: got %q", result.Stdout)
	}
}
