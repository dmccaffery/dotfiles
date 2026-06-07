// Command dot is a multi-call binary backing the dotfiles helper scripts. It
// dispatches on argv[0], so a symlink named `worktree` runs the worktree
// command, while `dot worktree …` works too. Symlinks are created by the
// install build stage from the hidden `dot applets` list.
package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/dmccaffery/dotfiles/internal/cmd/cmdutil"
	"github.com/dmccaffery/dotfiles/internal/cmd/root"
	"github.com/dmccaffery/dotfiles/internal/envx"
)

// version is overridden via -ldflags "-X main.version=…" at build time.
var version = "dev"

func main() {
	deps := &cmdutil.Deps{Env: envx.System()}
	rootCmd := root.NewRootCmd(version, deps)

	// Multi-call dispatch: when invoked via a symlink whose name matches a
	// command (argv[0] is e.g. `worktree`), prepend it so Cobra routes there —
	// and `--help` still works at every level. A binary named anything else
	// (`dot`, a dev build, …) just runs as plain `dot`.
	args := os.Args[1:]
	if applet := filepath.Base(os.Args[0]); applet != "dot" && root.HasCommand(rootCmd, applet) {
		args = append([]string{applet}, args...)
	}
	rootCmd.SetArgs(args)

	if err := rootCmd.Execute(); err != nil {
		var exitErr *exec.ExitError
		switch {
		case errors.As(err, &exitErr):
			// A forwarded child (gh, git, …) failed; its output already went to the
			// terminal, so just mirror its exit code.
			os.Exit(exitErr.ExitCode())
		case errors.Is(err, cmdutil.ErrSilent):
			// The command logged its own diagnostics.
			os.Exit(1)
		default:
			// A cobra-level error (unknown command/flag) — print it.
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
