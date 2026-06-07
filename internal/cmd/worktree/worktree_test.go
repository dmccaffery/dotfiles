package worktree_test

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dmccaffery/dotfiles/internal/cmd/cmdtest"
	"github.com/dmccaffery/dotfiles/internal/cmd/worktree"
	"github.com/dmccaffery/dotfiles/internal/envx"
	"github.com/dmccaffery/dotfiles/internal/execx"
)

func TestWorktreeStartStdoutIsPathOnly(t *testing.T) {
	home := t.TempDir()
	deps := cmdtest.NewDeps(t)
	deps.Env = envx.New(home, nil)
	out, _, err := cmdtest.Run(t, worktree.NewCmd(deps), `{"name":"feature"}`, "start", "/repos/myrepo")
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(home, ".cache", "agent", "worktrees", "myrepo-feature")
	if got := strings.TrimSpace(out); got != want {
		t.Fatalf("stdout: got %q want %q", got, want)
	}
	if strings.Contains(strings.TrimSpace(out), "\n") {
		t.Fatalf("stdout must be exactly the path, got %q", out)
	}
	if !cmdtest.ContainsLine(cmdtest.Fake(deps).Lines(), "worktree add") {
		t.Fatalf("expected a git worktree add call, got %v", cmdtest.Fake(deps).Lines())
	}
}

func TestWorktreeStartCreatesBranchWhenMissing(t *testing.T) {
	home := t.TempDir()
	repo := "/repos/myrepo"
	branch := "agent/myrepo-feature"
	deps := cmdtest.NewDeps(t)
	deps.Env = envx.New(home, nil)
	cmdtest.Fake(deps).Errs["git -C "+repo+" show-ref --verify --quiet refs/heads/"+branch] = errors.New("absent")
	if _, _, err := cmdtest.Run(t, worktree.NewCmd(deps), "", "start", repo, "feature"); err != nil {
		t.Fatal(err)
	}
	wantPath := filepath.Join(home, ".cache", "agent", "worktrees", "myrepo-feature")
	if !cmdtest.ContainsLine(cmdtest.Fake(deps).Lines(), "worktree add -b "+branch+" "+wantPath) {
		t.Fatalf("expected create-branch call, got %v", cmdtest.Fake(deps).Lines())
	}
}

func TestWorktreeEndKeepsNonAgentBranch(t *testing.T) {
	wt := t.TempDir()
	deps := cmdtest.NewDeps(t)
	fake := cmdtest.Fake(deps)
	fake.Responses["git -C "+wt+" rev-parse --abbrev-ref HEAD"] = execx.Result{Stdout: "main\n"}
	fake.Responses["git -C "+wt+" rev-parse --git-common-dir"] = execx.Result{Stdout: "/repos/myrepo/.git\n"}
	if _, _, err := cmdtest.Run(t, worktree.NewCmd(deps), "", "end", wt); err != nil {
		t.Fatal(err)
	}
	if cmdtest.ContainsLine(fake.Lines(), "branch -D") {
		t.Fatalf("must not delete a non-agent branch, calls: %v", fake.Lines())
	}
	if !cmdtest.ContainsLine(fake.Lines(), "worktree remove --force "+wt) {
		t.Fatalf("expected worktree remove, got %v", fake.Lines())
	}
}

func TestWorktreeEndDeletesAgentBranch(t *testing.T) {
	wt := t.TempDir()
	deps := cmdtest.NewDeps(t)
	fake := cmdtest.Fake(deps)
	fake.Responses["git -C "+wt+" rev-parse --abbrev-ref HEAD"] = execx.Result{Stdout: "agent/foo\n"}
	fake.Responses["git -C "+wt+" rev-parse --git-common-dir"] = execx.Result{Stdout: "/repos/myrepo/.git\n"}
	if _, _, err := cmdtest.Run(t, worktree.NewCmd(deps), "", "end", wt); err != nil {
		t.Fatal(err)
	}
	if !cmdtest.ContainsLine(fake.Lines(), "git -C /repos/myrepo branch -D agent/foo") {
		t.Fatalf("expected agent branch delete, got %v", fake.Lines())
	}
}
