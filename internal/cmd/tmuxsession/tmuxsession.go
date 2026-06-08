// Package tmuxsession implements the tmux-session command (start/attach/end).
package tmuxsession

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/dmccaffery/dotfiles/internal/cmd/cmdutil"
	worktreecmd "github.com/dmccaffery/dotfiles/internal/cmd/worktree"
	"github.com/dmccaffery/dotfiles/internal/execx"
	"github.com/dmccaffery/dotfiles/internal/ui"
	wt "github.com/dmccaffery/dotfiles/internal/worktree"
)

type tmuxCmd struct{ deps *cmdutil.Deps }

// NewCmd builds the tmux-session command.
func NewCmd(deps *cmdutil.Deps) *cobra.Command {
	t := &tmuxCmd{deps: deps}
	cmd := &cobra.Command{
		Use:   "tmux-session <command>",
		Short: "Tmux session manager for agent development",
		Long: "Picks a repo (or per-worktree checkout) and lays out a tmux session with nvim,\n" +
			"a shell, and a window per available coding agent. Subcommands: start, attach, end.",
	}
	cmd.AddCommand(
		&cobra.Command{
			Use:   "start [query] [worktree]",
			Short: "Pick a repo and create/attach a session (optionally per worktree)",
			Args:  cobra.MaximumNArgs(2),
			RunE:  t.start,
		},
		&cobra.Command{
			Use:   "attach [query]",
			Short: "Pick a running session and attach (or switch client inside tmux)",
			Args:  cobra.MaximumNArgs(1),
			RunE:  t.attach,
		},
	)
	end := &cobra.Command{
		Use:   "end [worktree...]",
		Short: "Pick agent worktrees and remove them (with dirty-state warnings)",
		RunE:  t.end,
	}
	end.Flags().BoolP("force", "f", false, "skip the confirmation when worktrees are dirty")
	cmd.AddCommand(end)
	return cmd
}

func (t *tmuxCmd) worktreesDir() string {
	return filepath.Join(t.deps.Env.Home(), ".cache", "agent", "worktrees")
}

func (t *tmuxCmd) start(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	deps := t.deps
	reposDir := deps.Env.GetOr("REPOS_DIR", filepath.Join(deps.Env.Home(), "Repos"))
	query := cmdutil.Arg(args, 0)
	wtName := cmdutil.Arg(args, 1)

	if !cmdutil.DirExists(reposDir) {
		reposDir = "."
	}
	if query == "." {
		reposDir = deps.Env.GetOr("PWD", cwd())
		query = ""
	}

	selected, err := cmdutil.PickOne(deps.Prompt, "select a repo", query, findGitRepos(reposDir, 4))
	if err != nil {
		return pickErr(deps, err)
	}
	if selected == "" {
		return nil
	}

	var name string
	if wtName != "" {
		path, err := worktreecmd.Start(ctx, deps, selected, wtName)
		if err != nil {
			return err
		}
		selected, name = path, filepath.Base(path)
	} else {
		name = wt.Sanitize(filepath.Base(selected))
	}

	if !t.hasSession(ctx, name) {
		if err := t.createSession(ctx, name, selected); err != nil {
			return err
		}
	}

	// Set the terminal title (OSC 0) so the tab shows the session, not the launcher.
	fmt.Fprintf(cmd.OutOrStdout(), "\033]0;%s\007", name)

	return deps.Runner.RunIO(ctx, cmdutil.Streams(cmd), "tmux", "-u", "attach-session", "-t", name, "-c", selected)
}

func (t *tmuxCmd) hasSession(ctx context.Context, name string) bool {
	_, err := t.deps.Runner.Run(ctx, "tmux", "has-session", "-t", name)
	return err == nil
}

// createSession lays out: nvim (editor pane + small shell pane), a window per
// available agent (opencode, codex, claude), then a plain zsh window.
func (t *tmuxCmd) createSession(ctx context.Context, name, dir string) error {
	r := t.deps.Runner
	editor := t.deps.Env.GetOr("EDITOR", "nvim")

	res, err := r.Run(ctx, "tmux", "-u", "new-session", "-d", "-P", "-F", "#{pane_id}",
		"-s", name, "-n", "  nvim", "-c", dir, "-x", "-", "-y", "-", editor, ".")
	if err != nil {
		t.deps.Log.Error("failed to create tmux session: " + err.Error())
		return cmdutil.ErrSilent
	}
	editorPane := strings.TrimSpace(res.Stdout)
	win, _ := r.Run(ctx, "tmux", "display-message", "-p", "-t", editorPane, "#{window_id}")
	editorWindow := strings.TrimSpace(win.Stdout)
	_, _ = r.Run(ctx, "tmux", "split-window", "-t", editorPane, "-v", "-l", "10%", "-c", dir)
	_, _ = r.Run(ctx, "tmux", "select-pane", "-t", editorPane)

	for _, agent := range []struct{ bin, window string }{
		{"opencode", "󰚩  opencode"},
		{"codex", "󱙺  codex"},
		{"claude", "󰯉  claude"},
	} {
		if path, err := r.Look(agent.bin); err == nil {
			_, _ = r.Run(ctx, "tmux", "new-window", "-a", "-d", "-t", editorWindow, "-c", dir, "-n", agent.window, path)
		}
	}
	_, _ = r.Run(ctx, "tmux", "new-window", "-a", "-d", "-t", editorWindow, "-n", "  zsh", "-c", dir)
	return nil
}

func (t *tmuxCmd) attach(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	deps := t.deps
	query := cmdutil.Arg(args, 0)
	if query == "." {
		query = ""
	}

	var sessions []string
	if res, err := deps.Runner.Run(ctx, "tmux", "list-session", "-F", "#S"); err == nil {
		sessions = cmdutil.NonEmptyLines(res.Stdout)
	}
	selected, err := cmdutil.PickOne(deps.Prompt, "select a session", query, sessions)
	if err != nil {
		return pickErr(deps, err)
	}
	if selected == "" {
		return nil
	}

	name := strings.ReplaceAll(filepath.Base(selected), ".", "_")
	if deps.Env.Get("TMUX") != "" {
		return deps.Runner.RunIO(ctx, cmdutil.Streams(cmd), "tmux", "switch-client", "-t", name)
	}
	return deps.Runner.RunIO(ctx, cmdutil.Streams(cmd), "tmux", "-u", "attach-session", "-t", name)
}

func (t *tmuxCmd) end(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	deps := t.deps
	force, _ := cmd.Flags().GetBool("force")
	dir := t.worktreesDir()

	if !cmdutil.DirExists(dir) {
		deps.Log.Warn("no agent worktrees at " + dir)
		return nil
	}

	var selections []string
	if len(args) > 0 {
		for _, a := range args {
			if filepath.IsAbs(a) {
				selections = append(selections, a)
			} else {
				selections = append(selections, filepath.Join(dir, a))
			}
		}
	} else {
		picked, err := deps.Prompt.MultiSelect("worktrees to remove", childDirs(dir))
		if err != nil {
			return pickErr(deps, err)
		}
		selections = picked
	}
	if len(selections) == 0 {
		return nil
	}

	needsConfirm := false
	for _, w := range selections {
		if !cmdutil.DirExists(w) {
			deps.Log.Warn("missing: " + w)
			continue
		}
		branch := gitField(ctx, deps.Runner, w, "rev-parse", "--abbrev-ref", "HEAD")
		if branch == "" {
			branch = "?"
		}
		uncommitted := cmdutil.CountNonEmptyLines(gitField(ctx, deps.Runner, w, "status", "--porcelain"))
		unpushed := cmdutil.Atoi(gitField(ctx, deps.Runner, w, "rev-list", "--count", "HEAD", "--not", "--remotes"))
		if uncommitted != 0 || unpushed != 0 {
			deps.Log.Warn(fmt.Sprintf("%s [%s] — uncommitted: %d, unpushed: %d", filepath.Base(w), branch, uncommitted, unpushed))
			needsConfirm = true
		} else {
			deps.Log.Info(fmt.Sprintf("%s [%s] — clean", filepath.Base(w), branch))
		}
	}

	if needsConfirm && !force {
		ok, err := deps.Prompt.Confirm("proceed with removal?", false)
		if errors.Is(err, ui.ErrNoTTY) || (err == nil && !ok) {
			deps.Log.Error("aborted")
			return cmdutil.ErrSilent
		}
		if err != nil {
			deps.Log.Error(err.Error())
			return cmdutil.ErrSilent
		}
	}

	for _, w := range selections {
		if !cmdutil.DirExists(w) {
			continue
		}
		if err := worktreecmd.End(ctx, deps, w); err != nil {
			return err
		}
	}
	return nil
}

func gitField(ctx context.Context, r execx.Runner, dir string, args ...string) string {
	res, err := r.Run(ctx, "git", append([]string{"-C", dir}, args...)...)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(res.Stdout)
}

// findGitRepos walks root up to maxDepth and returns directories that contain a
// `.git` directory, pruning each match (it does not descend into a repo).
func findGitRepos(root string, maxDepth int) []string {
	var repos []string
	base := depthOf(root)
	_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || !d.IsDir() {
			return nil
		}
		if cmdutil.DirExists(filepath.Join(path, ".git")) {
			repos = append(repos, path)
			return fs.SkipDir
		}
		if depthOf(path)-base >= maxDepth {
			return fs.SkipDir
		}
		return nil
	})
	return repos
}

func depthOf(p string) int { return strings.Count(filepath.Clean(p), string(filepath.Separator)) }

func childDirs(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var dirs []string
	for _, e := range entries {
		if e.IsDir() {
			dirs = append(dirs, filepath.Join(dir, e.Name()))
		}
	}
	return dirs
}

func cwd() string {
	d, _ := os.Getwd()
	return d
}

func pickErr(deps *cmdutil.Deps, err error) error {
	if errors.Is(err, ui.ErrNoTTY) {
		deps.Log.Error("no tty available to make a selection")
	} else if err != nil {
		deps.Log.Error(err.Error())
	}
	return cmdutil.ErrSilent
}
