// Package gitresign implements the git-resign command.
package gitresign

import (
	"github.com/spf13/cobra"

	"github.com/dmccaffery/dotfiles/internal/cmd/cmdutil"
)

// NewCmd builds the git-resign command.
func NewCmd(deps *cmdutil.Deps) *cobra.Command {
	return &cobra.Command{
		Use:   "git-resign <target>",
		Short: "Re-sign every commit in a range with the current signing key",
		Long: "Runs `git rebase -i <target>` with `--exec 'git commit --amend --no-edit -n -S'`,\n" +
			"re-signing each commit from <target> to HEAD. Also invoked as `git resign <target>`\n" +
			"via the git-resign symlink (git's git-* subcommand dispatch).",
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target := cmdutil.Arg(args, 0)
			if target == "" {
				deps.Log.Error("a target must be set; exiting...")
				return cmdutil.ErrSilent
			}
			// Interactive rebase: stdio passes through so the editor gets the tty,
			// and git's own exit code/output reaches the caller.
			return deps.Runner.RunIO(cmd.Context(), cmdutil.Streams(cmd), "git",
				"rebase", "--exec", "git commit --amend --no-edit -n -S", "-i", target)
		},
	}
}
