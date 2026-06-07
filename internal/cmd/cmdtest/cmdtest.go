// Package cmdtest provides helpers for driving dot commands in tests. It is
// imported only by _test files, so it never reaches the binary.
package cmdtest

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/dmccaffery/dotfiles/internal/cmd/cmdutil"
	"github.com/dmccaffery/dotfiles/internal/envx"
	"github.com/dmccaffery/dotfiles/internal/execx"
	"github.com/dmccaffery/dotfiles/internal/logx"
	"github.com/dmccaffery/dotfiles/internal/ui"
)

// NewDeps returns a Deps wired for tests: a fresh Fake runner, a temp-home Env, a
// discarding logger and a Fake prompter. Commands are driven directly (not through
// root), so Log/Prompt are populated here rather than by root's PersistentPreRunE.
func NewDeps(t *testing.T) *cmdutil.Deps {
	t.Helper()
	return &cmdutil.Deps{
		Runner: execx.NewFake(),
		Env:    envx.New(t.TempDir(), nil),
		Log:    logx.For(io.Discard, false),
		Prompt: &ui.Fake{},
	}
}

// Fake returns deps.Runner as *execx.Fake so tests can script Responses/Errs.
func Fake(deps *cmdutil.Deps) *execx.Fake {
	return deps.Runner.(*execx.Fake)
}

// Run drives cmd with the given stdin and args, capturing stdout and stderr.
func Run(t *testing.T, cmd *cobra.Command, stdin string, args ...string) (stdout, stderr string, err error) {
	t.Helper()
	var out, errb bytes.Buffer
	cmd.SetArgs(args)
	cmd.SetIn(strings.NewReader(stdin))
	cmd.SetOut(&out)
	cmd.SetErr(&errb)
	err = cmd.Execute()
	return out.String(), errb.String(), err
}

// ContainsLine reports whether any line contains sub.
func ContainsLine(lines []string, sub string) bool {
	for _, l := range lines {
		if strings.Contains(l, sub) {
			return true
		}
	}
	return false
}
