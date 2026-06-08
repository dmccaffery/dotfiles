package tmuxsession

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/dmccaffery/dotfiles/internal/cmd/cmdtest"
	"github.com/dmccaffery/dotfiles/internal/cmd/cmdutil"
	"github.com/dmccaffery/dotfiles/internal/envx"
	"github.com/dmccaffery/dotfiles/internal/execx"
	"github.com/dmccaffery/dotfiles/internal/logx"
	"github.com/dmccaffery/dotfiles/internal/ui"
)

func TestFindGitRepos(t *testing.T) {
	root := t.TempDir()
	mk := func(parts ...string) {
		if err := os.MkdirAll(filepath.Join(append([]string{root}, parts...)...), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	mk("a", ".git")                        // repo at depth 1
	mk("b", "c", ".git")                   // repo at depth 2
	mk("b", "c", "d")                      // inside repo c — must be pruned
	mk("w", "x", "y", "z", ".git")         // repo at depth 4
	mk("deep", "1", "2", "3", "4", ".git") // depth 5 — beyond maxdepth
	mk("e")                                // not a repo

	got := findGitRepos(root, 4)
	want := map[string]bool{
		filepath.Join(root, "a"):                true,
		filepath.Join(root, "b", "c"):           true,
		filepath.Join(root, "w", "x", "y", "z"): true,
	}
	if len(got) != len(want) {
		t.Fatalf("got %v, want keys %v", got, want)
	}
	for _, r := range got {
		if !want[r] {
			t.Errorf("unexpected repo %q", r)
		}
	}
}

func endDeps(t *testing.T) (*cmdutil.Deps, *execx.Fake, string) {
	t.Helper()
	home := t.TempDir()
	wt1 := filepath.Join(home, ".cache", "agent", "worktrees", "wt1")
	if err := os.MkdirAll(wt1, 0o755); err != nil {
		t.Fatal(err)
	}
	fake := execx.NewFake()
	deps := &cmdutil.Deps{
		Runner: fake,
		Env:    envx.New(home, nil),
		Log:    logx.For(io.Discard, false),
		Prompt: &ui.Fake{},
	}
	return deps, fake, wt1
}

func TestTmuxEndRemovesCleanWorktree(t *testing.T) {
	deps, fake, wt1 := endDeps(t)
	fake.Responses["git -C "+wt1+" rev-parse --abbrev-ref HEAD"] = execx.Result{Stdout: "agent/wt1\n"}
	fake.Responses["git -C "+wt1+" rev-parse --git-common-dir"] = execx.Result{Stdout: "/repos/r/.git\n"}

	if _, _, err := cmdtest.Run(t, NewCmd(deps), "", "end", "wt1"); err != nil {
		t.Fatal(err)
	}
	if !cmdtest.ContainsLine(fake.Lines(), "worktree remove --force "+wt1) {
		t.Fatalf("expected worktree removal, got %v", fake.Lines())
	}
}

func TestTmuxEndAbortsWhenDirtyAndDeclined(t *testing.T) {
	deps, fake, wt1 := endDeps(t)
	deps.Prompt = &ui.Fake{Replies: []bool{false}} // decline the confirmation
	fake.Responses["git -C "+wt1+" status --porcelain"] = execx.Result{Stdout: " M file\n"}
	fake.Responses["git -C "+wt1+" rev-parse --abbrev-ref HEAD"] = execx.Result{Stdout: "agent/wt1\n"}

	if _, _, err := cmdtest.Run(t, NewCmd(deps), "", "end", "wt1"); err == nil {
		t.Fatal("expected an error after declining")
	}
	if cmdtest.ContainsLine(fake.Lines(), "worktree remove") {
		t.Fatalf("must not remove when declined, got %v", fake.Lines())
	}
}

func TestTmuxEndForceSkipsConfirm(t *testing.T) {
	deps, fake, wt1 := endDeps(t)
	prompt := &ui.Fake{}
	deps.Prompt = prompt
	fake.Responses["git -C "+wt1+" status --porcelain"] = execx.Result{Stdout: " M file\n"} // dirty
	fake.Responses["git -C "+wt1+" rev-parse --abbrev-ref HEAD"] = execx.Result{Stdout: "agent/wt1\n"}

	if _, _, err := cmdtest.Run(t, NewCmd(deps), "", "end", "-f", "wt1"); err != nil {
		t.Fatal(err)
	}
	if !cmdtest.ContainsLine(fake.Lines(), "worktree remove --force "+wt1) {
		t.Fatalf("expected forced removal, got %v", fake.Lines())
	}
	if len(prompt.Asked) != 0 {
		t.Fatalf("--force must not prompt, asked %v", prompt.Asked)
	}
}

func TestTmuxEndMultiSelect(t *testing.T) {
	deps, fake, wt1 := endDeps(t)
	deps.Prompt = &ui.Fake{MultiSelections: [][]string{{wt1}}}
	fake.Responses["git -C "+wt1+" rev-parse --abbrev-ref HEAD"] = execx.Result{Stdout: "agent/wt1\n"}

	if _, _, err := cmdtest.Run(t, NewCmd(deps), "", "end"); err != nil {
		t.Fatal(err)
	}
	if !cmdtest.ContainsLine(fake.Lines(), "worktree remove --force "+wt1) {
		t.Fatalf("expected removal of the multi-selected worktree, got %v", fake.Lines())
	}
}
