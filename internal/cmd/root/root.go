// Package root assembles the dot root command from the per-command packages and
// owns the applet registry that the build stage reads.
package root

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/dmccaffery/dotfiles/internal/cmd/agenttmux"
	"github.com/dmccaffery/dotfiles/internal/cmd/brewfile"
	"github.com/dmccaffery/dotfiles/internal/cmd/cmdutil"
	"github.com/dmccaffery/dotfiles/internal/cmd/fzfpreview"
	"github.com/dmccaffery/dotfiles/internal/cmd/ghswitch"
	"github.com/dmccaffery/dotfiles/internal/cmd/githubauth"
	"github.com/dmccaffery/dotfiles/internal/cmd/gitresign"
	"github.com/dmccaffery/dotfiles/internal/cmd/printcolors"
	"github.com/dmccaffery/dotfiles/internal/cmd/profileshell"
	"github.com/dmccaffery/dotfiles/internal/cmd/resetbackground"
	"github.com/dmccaffery/dotfiles/internal/cmd/tmuxsession"
	"github.com/dmccaffery/dotfiles/internal/cmd/worktree"
	"github.com/dmccaffery/dotfiles/internal/cmd/zs"
	"github.com/dmccaffery/dotfiles/internal/execx"
	"github.com/dmccaffery/dotfiles/internal/logx"
	"github.com/dmccaffery/dotfiles/internal/ui"
)

// appletAnnotation marks a command whose name the build stage exposes as a
// standalone symlink (enumerated by the hidden `dot applets` command).
const appletAnnotation = "dot.applet"

func applet(c *cobra.Command) *cobra.Command {
	if c.Annotations == nil {
		c.Annotations = map[string]string{}
	}
	c.Annotations[appletAnnotation] = "true"
	return c
}

// NewRootCmd builds the dot root command with every subcommand registered.
func NewRootCmd(version string, deps *cmdutil.Deps) *cobra.Command {
	root := &cobra.Command{
		Use:   "dot",
		Short: "dotfiles helper CLI",
		Long: "dot is a multi-call binary backing the dotfiles helper commands.\n" +
			"Run it as `dot <command>` or via a symlink named after the command\n" +
			"(e.g. `worktree`); it dispatches on argv[0].",
		Version:           version,
		SilenceUsage:      true,
		SilenceErrors:     true,
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			if deps.Runner == nil {
				deps.Runner = execx.Real{}
			}
			if deps.Log == nil {
				w := cmd.ErrOrStderr()
				deps.Log = logx.For(w, logx.IsTerminal(w))
			}
			if deps.Prompt == nil {
				deps.Prompt = ui.New()
			}
			return nil
		},
	}

	root.AddCommand(
		applet(worktree.NewCmd(deps)),
		applet(agenttmux.NewCmd(deps)),
		applet(brewfile.NewCmd(deps)),
		applet(printcolors.NewCmd(deps)),
		applet(profileshell.NewCmd(deps)),
		applet(gitresign.NewCmd(deps)),
		applet(ghswitch.NewCmd(deps)),
		applet(fzfpreview.NewCmd(deps)),
		applet(resetbackground.NewCmd(deps)),
		applet(zs.NewCmd(deps)),
		applet(githubauth.NewCmd(deps)),
		applet(tmuxsession.NewCmd(deps)),
	)
	root.AddCommand(newAppletsCmd(root))
	return root
}

// HasCommand reports whether name matches a registered subcommand. The argv[0]
// dispatch uses it so a binary whose name is neither "dot" nor a command still
// behaves as plain `dot` instead of failing on an unknown applet.
func HasCommand(root *cobra.Command, name string) bool {
	for _, c := range root.Commands() {
		if c.Name() == name {
			return true
		}
	}
	return false
}

// newAppletsCmd prints the names of commands marked as applets — the single
// source of truth the build stage reads to create the per-command symlinks.
func newAppletsCmd(root *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:    "applets",
		Short:  "List command names that should be symlinked to the dot binary",
		Hidden: true,
		Args:   cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			for _, c := range root.Commands() {
				if c.Annotations[appletAnnotation] == "true" {
					fmt.Fprintln(cmd.OutOrStdout(), c.Name())
				}
			}
			return nil
		},
	}
}
