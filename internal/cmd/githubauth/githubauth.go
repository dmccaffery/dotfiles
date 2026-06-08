// Package githubauth implements the git-github-auth command.
package githubauth

import (
	"context"
	"errors"
	"strings"

	"github.com/spf13/cobra"

	"github.com/dmccaffery/dotfiles/internal/cmd/cmdutil"
	"github.com/dmccaffery/dotfiles/internal/execx"
	"github.com/dmccaffery/dotfiles/internal/ui"
)

var requiredScopes = []string{
	"gist", "notifications", "project", "repo", "user", "workflow",
	"read:org", "read:public_key", "read:ssh_signing_key", "write:ssh_signing_key",
}

// NewCmd builds the git-github-auth command.
func NewCmd(deps *cmdutil.Deps) *cobra.Command {
	return &cobra.Command{
		Use:   "git-github-auth [account]",
		Short: "Ensure gh is authenticated with the required scopes",
		Long: "Ensures the GitHub CLI is logged in (optionally for a chosen account) with the\n" +
			"scopes this setup needs — logging in, switching account, or refreshing scopes as required.",
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			log, r := deps.Log, deps.Runner
			scopeCSV := strings.Join(requiredScopes, ",")

			log.Info("checking github auth...")

			account := cmdutil.Arg(args, 0)
			if account == "" {
				picked, err := cmdutil.PickOne(deps.Prompt, "select github account", "",
					append(ghLogins(ctx, r), "new account"))
				if err != nil {
					if errors.Is(err, ui.ErrNoTTY) {
						log.Error("no account specified and no tty to pick one")
					}
					return cmdutil.ErrSilent
				}
				switch picked {
				case "":
					return nil // aborted
				case "new account":
					return ghLogin(ctx, cmd, r, scopeCSV)
				default:
					account = picked
					if err := r.RunIO(ctx, cmdutil.Streams(cmd), "gh", "auth", "switch", "--user", account); err != nil {
						return err
					}
				}
			}

			// gh auth status exits non-zero when not authenticated at all → log in.
			if _, err := r.Run(ctx, "gh", "auth", "status", "--json", "hosts",
				"--jq", `.hosts."github.com".[] | select(.login == "`+account+`") | .active`); err != nil {
				log.Warn("not logged into github, will login now...")
				return ghLogin(ctx, cmd, r, scopeCSV)
			}

			if active := ghActive(ctx, r); account != active {
				_, _ = r.Run(ctx, "gh", "auth", "switch", "--user", account)
			}

			if scope, missing := firstMissingScope(ghScopes(ctx, r), requiredScopes); missing {
				log.Warn("not logged into github with the required scope: " + scope)
				if err := r.RunIO(ctx, cmdutil.Streams(cmd), "gh", "auth", "refresh",
					"--hostname", "github.com", "--scopes", scopeCSV, "--clipboard"); err != nil {
					return err
				}
			}

			log.Info("logged into github with the required scopes")
			return nil
		},
	}
}

func ghLogin(ctx context.Context, cmd *cobra.Command, r execx.Runner, scopeCSV string) error {
	return r.RunIO(ctx, cmdutil.Streams(cmd), "gh", "auth", "login",
		"--git-protocol", "https", "--hostname", "github.com", "--scopes", scopeCSV, "--web", "--clipboard")
}

func ghLogins(ctx context.Context, r execx.Runner) []string {
	res, err := r.Run(ctx, "gh", "auth", "status", "--json", "hosts", "--jq", `.hosts."github.com".[].login`)
	if err != nil {
		return nil
	}
	return cmdutil.NonEmptyLines(res.Stdout)
}

func ghActive(ctx context.Context, r execx.Runner) string {
	res, err := r.Run(ctx, "gh", "auth", "status", "--active", "--json", "hosts", "--jq", `.hosts["github.com"][0].login`)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(res.Stdout)
}

func ghScopes(ctx context.Context, r execx.Runner) string {
	res, err := r.Run(ctx, "gh", "auth", "status", "--active", "--json", "hosts", "--jq", `.hosts."github.com".[] | .scopes`)
	if err != nil {
		return ""
	}
	return res.Stdout
}

// firstMissingScope reports the first required scope absent from the current
// scopes string (a substring match, like the shell's `grep -F`).
func firstMissingScope(current string, required []string) (string, bool) {
	for _, s := range required {
		if !strings.Contains(current, s) {
			return s, true
		}
	}
	return "", false
}
