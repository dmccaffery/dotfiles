// Package execx is a small, mockable seam over os/exec so commands that shell
// out to git, tmux, brew and friends can be unit-tested with a fake runner.
package execx

import (
	"bytes"
	"context"
	"io"
	"os/exec"
)

// Result is the captured outcome of a Run.
type Result struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

// Streams bundles the three standard streams for RunIO.
type Streams struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
}

// Runner runs external commands and resolves their paths. Real wraps os/exec;
// Fake (fake.go) records calls and returns scripted results.
type Runner interface {
	// Run executes name+args, capturing stdout and stderr. The error is non-nil
	// when the process fails to start or exits non-zero; Result.ExitCode carries
	// the code in the latter case.
	Run(ctx context.Context, name string, args ...string) (Result, error)
	// RunIO executes name+args with the given streams wired straight through to
	// the child — for interactive children (fzf, pinentry, an editor).
	RunIO(ctx context.Context, s Streams, name string, args ...string) error
	// Look resolves name on PATH (like `command -v`).
	Look(name string) (string, error)
}

// Real is a Runner backed by os/exec and the real PATH.
type Real struct{}

// Run implements Runner.
func (Real) Run(ctx context.Context, name string, args ...string) (Result, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	var out, errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	err := cmd.Run()
	res := Result{Stdout: out.String(), Stderr: errb.String()}
	if cmd.ProcessState != nil {
		res.ExitCode = cmd.ProcessState.ExitCode()
	}
	return res, err
}

// RunIO implements Runner.
func (Real) RunIO(ctx context.Context, s Streams, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdin = s.In
	cmd.Stdout = s.Out
	cmd.Stderr = s.Err
	return cmd.Run()
}

// Look implements Runner.
func (Real) Look(name string) (string, error) { return exec.LookPath(name) }
