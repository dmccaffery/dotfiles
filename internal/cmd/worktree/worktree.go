// Package worktree implements the worktree command (the agent worktree lifecycle).
// Start and End are exported so tmux-session can drive the lifecycle in-process.
package worktree

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/dmccaffery/dotfiles/internal/cmd/cmdutil"
	"github.com/dmccaffery/dotfiles/internal/execx"
	"github.com/dmccaffery/dotfiles/internal/logx"
	wt "github.com/dmccaffery/dotfiles/internal/worktree"
)

// NewCmd builds the worktree command.
func NewCmd(deps *cmdutil.Deps) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "worktree",
		Short: "Manage agent git worktrees",
		Long: "Agent worktree lifecycle. `start` creates/reuses a worktree on an\n" +
			"agent/<name> branch and prints its path; `end` removes the worktree, its\n" +
			"tmux session and the agent/* branch. Both also accept hook JSON on stdin.",
	}
	cmd.AddCommand(
		&cobra.Command{
			Use:   "start [repo] [suffix]",
			Short: "Create or reuse an agent/* worktree and print its path",
			Args:  cobra.MaximumNArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				ctx := cmd.Context()
				repo := cmdutil.Arg(args, 0)
				if repo == "" {
					repo = deps.Env.Get("CLAUDE_PROJECT_DIR")
				}
				if repo == "" {
					deps.Log.Warn("no repository specified, attempting to fall back to repository root...")
					if res, err := deps.Runner.Run(ctx, "git", "rev-parse", "--show-toplevel"); err == nil {
						repo = strings.TrimSpace(res.Stdout)
					}
				}
				if repo == "" {
					deps.Log.Error("no repository found; did you forget to specify an argument?")
					return cmdutil.ErrSilent
				}

				suffix := cmdutil.Arg(args, 1)
				if suffix == "" {
					if data, ok := readIfPiped(cmd.InOrStdin()); ok {
						suffix = wt.ParseStartName(data)
					}
				}
				if suffix == "" {
					deps.Log.Warn("no suffix provided; using current timestamp")
					suffix = time.Now().UTC().Format("20060102-150405")
				}

				path, err := Start(ctx, deps, repo, suffix)
				if err != nil {
					return err
				}
				// The path is the command's result — stdout only, so hooks/callers can read it.
				fmt.Fprintln(cmd.OutOrStdout(), path)
				return nil
			},
		},
		&cobra.Command{
			Use:   "end [path]",
			Short: "Remove a worktree, its agent/* branch and tmux session",
			Args:  cobra.MaximumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				path := cmdutil.Arg(args, 0)
				if path == "" {
					if data, ok := readIfPiped(cmd.InOrStdin()); ok {
						path = wt.ParseEndPath(data)
					}
				}
				if path == "" {
					deps.Log.Error(`worktree end: no worktree path; pass as an argument or pipe JSON {"worktree_path":"..."}`)
					return cmdutil.ErrSilent
				}
				return End(cmd.Context(), deps, path)
			},
		},
	)
	return cmd
}

// worktreesRoot is where worktrees live: "$HOME/.cache/agent/worktrees", as in
// the original shell (a literal ~/.cache, not $XDG_CACHE_HOME).
func worktreesRoot(deps *cmdutil.Deps) string {
	return filepath.Join(deps.Env.Home(), ".cache", "agent", "worktrees")
}

// Start creates or reuses an agent/<name> worktree for repo+suffix and returns its
// path. repo and suffix are already resolved by the caller.
func Start(ctx context.Context, deps *cmdutil.Deps, repo, suffix string) (string, error) {
	log, r := deps.Log, deps.Runner
	names := wt.Derive(repo, suffix, worktreesRoot(deps))

	switch {
	case cmdutil.DirExists(names.Path):
		log.Warn(fmt.Sprintf("worktree already exists at %s; reusing", names.Path))
	case branchExists(ctx, r, repo, names.Branch):
		log.Info(fmt.Sprintf("branch %s exists; checking it out at %s", names.Branch, names.Path))
		if _, err := r.Run(ctx, "git", "-C", repo, "worktree", "add", names.Path, names.Branch); err != nil {
			log.Error(fmt.Sprintf("failed to check out worktree: %v", err))
			return "", cmdutil.ErrSilent
		}
	default:
		log.Info(fmt.Sprintf("creating worktree %s at %s", names.Name, names.Path))
		if _, err := r.Run(ctx, "git", "-C", repo, "worktree", "add", "-b", names.Branch, names.Path); err != nil {
			log.Error(fmt.Sprintf("failed to create worktree: %v", err))
			return "", cmdutil.ErrSilent
		}
	}

	log.Info("worktree "+names.Name+" is ready",
		logx.WorktreeAttr(logx.Worktree{Name: names.Name, Path: names.Path, Branch: names.Branch, Repo: repo}))
	return names.Path, nil
}

// End removes a worktree, its agent/* branch and matching tmux session.
func End(ctx context.Context, deps *cmdutil.Deps, path string) error {
	log, r := deps.Log, deps.Runner

	if !cmdutil.DirExists(path) {
		log.Warn(fmt.Sprintf("worktree at %s no longer exists; exiting gracefully", path))
		return nil
	}

	branch := gitOutput(ctx, r, "git", "-C", path, "rev-parse", "--abbrev-ref", "HEAD")
	main := mainRepo(ctx, r, path)

	uncommitted := cmdutil.CountNonEmptyLines(gitOutput(ctx, r, "git", "-C", path, "status", "--porcelain"))
	unpushed := cmdutil.Atoi(gitOutput(ctx, r, "git", "-C", path, "rev-list", "--count", "HEAD", "--not", "--remotes"))
	if uncommitted != 0 || unpushed != 0 {
		shown := branch
		if shown == "" {
			shown = "?"
		}
		log.Warn(fmt.Sprintf("%s [%s] — uncommitted: %d, unpushed: %d", path, shown, uncommitted, unpushed))
	}

	session := filepath.Base(path)
	if _, err := r.Look("tmux"); err == nil {
		_, _ = r.Run(ctx, "tmux", "kill-session", "-t", session)
	}

	if main != "" {
		_, _ = r.Run(ctx, "git", "-C", main, "worktree", "remove", "--force", path)
	} else {
		_, _ = r.Run(ctx, "git", "worktree", "remove", "--force", path)
	}

	if wt.IsAgentBranch(branch) {
		if main != "" {
			_, _ = r.Run(ctx, "git", "-C", main, "branch", "-D", branch)
		} else {
			_, _ = r.Run(ctx, "git", "branch", "-D", branch)
		}
	}

	log.Info("removed " + session)
	return nil
}

func branchExists(ctx context.Context, r execx.Runner, repo, branch string) bool {
	_, err := r.Run(ctx, "git", "-C", repo, "show-ref", "--verify", "--quiet", "refs/heads/"+branch)
	return err == nil
}

// mainRepo resolves the main repository from a linked worktree (dirname of its
// --git-common-dir).
func mainRepo(ctx context.Context, r execx.Runner, path string) string {
	common := gitOutput(ctx, r, "git", "-C", path, "rev-parse", "--git-common-dir")
	if common == "" {
		return ""
	}
	if !filepath.IsAbs(common) {
		common = filepath.Join(path, common)
	}
	return filepath.Dir(common)
}

func gitOutput(ctx context.Context, r execx.Runner, name string, args ...string) string {
	res, err := r.Run(ctx, name, args...)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(res.Stdout)
}

// readIfPiped returns stdin's contents when it is not a terminal (a hook pipe),
// matching the shell's `[ ! -t 0 ]` guard. For a tty it returns false.
func readIfPiped(in io.Reader) ([]byte, bool) {
	if f, ok := in.(*os.File); ok && logx.IsTerminal(f) {
		return nil, false
	}
	data, err := io.ReadAll(in)
	if err != nil || len(data) == 0 {
		return nil, false
	}
	return data, true
}
