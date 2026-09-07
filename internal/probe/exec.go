// internal/probe/exec.go
package probe

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

type ExecOpts struct {
	Args    []string
	Env     map[string]string
	Timeout time.Duration
	UsePTY  bool
}

type ExecResult struct {
	Stdout   []byte
	Stderr   []byte
	ExitCode int
	Duration time.Duration
	TimedOut bool
	Command  string
}

func Run(ctx context.Context, binary string, opts ExecOpts) (*ExecResult, error) {
	timeout := opts.Timeout
	if timeout == 0 {
		timeout = 10 * time.Second
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, binary, opts.Args...)

	cmd.Env = buildEnv(opts.Env)

	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	start := time.Now()
	err := cmd.Run()
	duration := time.Since(start)

	exitCode := 0
	timedOut := false

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			timedOut = true
			exitCode = -1
		} else if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			return nil, fmt.Errorf("exec %s: %w", binary, err)
		}
	}

	cmdStr := binary
	if len(opts.Args) > 0 {
		cmdStr += " " + strings.Join(opts.Args, " ")
	}

	return &ExecResult{
		Stdout:   stdout.Bytes(),
		Stderr:   stderr.Bytes(),
		ExitCode: exitCode,
		Duration: duration,
		TimedOut: timedOut,
		Command:  cmdStr,
	}, nil
}

func buildEnv(overrides map[string]string) []string {
	tmpDir, _ := os.MkdirTemp("", "cli-acs-*")
	homeDir := filepath.Join(tmpDir, "home")
	os.MkdirAll(homeDir, 0755)

	base := map[string]string{
		"PATH":   os.Getenv("PATH"),
		"HOME":   homeDir,
		"TMPDIR": tmpDir,
		"LANG":   "C.UTF-8",
		"TERM":   "xterm-256color",
	}

	for k, v := range overrides {
		base[k] = v
	}

	env := make([]string, 0, len(base))
	for k, v := range base {
		env = append(env, k+"="+v)
	}
	return env
}
