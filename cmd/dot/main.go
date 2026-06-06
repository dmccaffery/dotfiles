// Command dot is a multi-call binary backing the dotfiles helper scripts. It
// dispatches on argv[0], so a symlink named `worktree` runs the worktree
// command, while `dot worktree …` works too. Symlinks are created by the
// install build stage from the hidden `dot applets` list.
package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dmccaffery/dotfiles/internal/cli"
	"github.com/dmccaffery/dotfiles/internal/envx"
)

// version is overridden via -ldflags "-X main.version=…" at build time.
var version = "dev"

func main() {
	deps := &cli.Deps{Env: envx.System()}
	root := cli.NewRootCmd(version, deps)

	// Multi-call dispatch: when invoked via a symlink whose name matches a
	// command (argv[0] is e.g. `worktree`), prepend it so Cobra routes there —
	// and `--help` still works at every level. A binary named anything else
	// (`dot`, a dev build, …) just runs as plain `dot`.
	args := os.Args[1:]
	if applet := filepath.Base(os.Args[0]); applet != "dot" && cli.HasCommand(root, applet) {
		args = append([]string{applet}, args...)
	}
	root.SetArgs(args)

	if err := root.Execute(); err != nil {
		// Commands log their own diagnostics and return ErrSilent; print only
		// cobra-level errors (unknown command/flag) so they aren't swallowed.
		if !errors.Is(err, cli.ErrSilent) {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}
}
