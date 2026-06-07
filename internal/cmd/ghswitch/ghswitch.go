// Package ghswitch implements the gh-switch-user command.
package ghswitch

import (
	"context"
	"strings"

	"github.com/spf13/cobra"

	"github.com/dmccaffery/dotfiles/internal/cmd/cmdutil"
	"github.com/dmccaffery/dotfiles/internal/execx"
)

// NewCmd builds the gh-switch-user command.
func NewCmd(deps *cmdutil.Deps) *cobra.Command {
	return &cobra.Command{
		Use:   "gh-switch-user [gh args...]",
		Short: "gh wrapper that switches to the account in `git config github.account` first",
		Long: "Reads `git config github.account`; if it is set and not already the active gh account,\n" +
			"runs `gh auth switch` to it, then forwards every argument to gh. Aliased as `gh`.",
		DisableFlagParsing: true, // forward all args/flags (incl. --help/--version) to gh untouched
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			r := deps.Runner

			if target := gitConfigValue(ctx, r, "github.account"); target != "" {
				if target != ghActiveLogin(ctx, r) {
					_, _ = r.Run(ctx, "gh", "auth", "switch", "--user", target)
				}
			}
			// Transparent forward: stdio and the exit code pass straight through.
			return r.RunIO(ctx, cmdutil.Streams(cmd), "gh", args...)
		},
	}
}

func gitConfigValue(ctx context.Context, r execx.Runner, key string) string {
	res, err := r.Run(ctx, "git", "config", key)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(res.Stdout)
}

func ghActiveLogin(ctx context.Context, r execx.Runner) string {
	res, err := r.Run(ctx, "gh", "auth", "status", "--active", "--json", "hosts", "--jq", `.hosts["github.com"][0].login`)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(res.Stdout)
}
