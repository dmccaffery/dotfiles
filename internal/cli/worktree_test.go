package cli_test

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dmccaffery/dotfiles/internal/cli"
	"github.com/dmccaffery/dotfiles/internal/envx"
	"github.com/dmccaffery/dotfiles/internal/execx"
)

func TestWorktreeStartStdoutIsPathOnly(t *testing.T) {
	home := t.TempDir()
	fake := execx.NewFake()
	deps := &cli.Deps{Runner: fake, Env: envx.New(home, nil)}
	out, _, err := runRoot(t, deps, `{"name":"feature"}`, "worktree", "start", "/repos/myrepo")
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
	if !containsLine(fake.Lines(), "worktree add") {
		t.Fatalf("expected a git worktree add call, got %v", fake.Lines())
	}
}

func TestWorktreeStartCreatesBranchWhenMissing(t *testing.T) {
	home := t.TempDir()
	repo := "/repos/myrepo"
	branch := "agent/myrepo-feature"
	fake := execx.NewFake()
	fake.Errs["git -C "+repo+" show-ref --verify --quiet refs/heads/"+branch] = errors.New("absent")
	deps := &cli.Deps{Runner: fake, Env: envx.New(home, nil)}
	if _, _, err := runRoot(t, deps, "", "worktree", "start", repo, "feature"); err != nil {
		t.Fatal(err)
	}
	wantPath := filepath.Join(home, ".cache", "agent", "worktrees", "myrepo-feature")
	if !containsLine(fake.Lines(), "worktree add -b "+branch+" "+wantPath) {
		t.Fatalf("expected create-branch call, got %v", fake.Lines())
	}
}

func TestWorktreeEndKeepsNonAgentBranch(t *testing.T) {
	wt := t.TempDir()
	fake := execx.NewFake()
	fake.Responses["git -C "+wt+" rev-parse --abbrev-ref HEAD"] = execx.Result{Stdout: "main\n"}
	fake.Responses["git -C "+wt+" rev-parse --git-common-dir"] = execx.Result{Stdout: "/repos/myrepo/.git\n"}
	deps := &cli.Deps{Runner: fake, Env: envx.New(t.TempDir(), nil)}
	if _, _, err := runRoot(t, deps, "", "worktree", "end", wt); err != nil {
		t.Fatal(err)
	}
	if containsLine(fake.Lines(), "branch -D") {
		t.Fatalf("must not delete a non-agent branch, calls: %v", fake.Lines())
	}
	if !containsLine(fake.Lines(), "worktree remove --force "+wt) {
		t.Fatalf("expected worktree remove, got %v", fake.Lines())
	}
}

func TestWorktreeEndDeletesAgentBranch(t *testing.T) {
	wt := t.TempDir()
	fake := execx.NewFake()
	fake.Responses["git -C "+wt+" rev-parse --abbrev-ref HEAD"] = execx.Result{Stdout: "agent/foo\n"}
	fake.Responses["git -C "+wt+" rev-parse --git-common-dir"] = execx.Result{Stdout: "/repos/myrepo/.git\n"}
	deps := &cli.Deps{Runner: fake, Env: envx.New(t.TempDir(), nil)}
	if _, _, err := runRoot(t, deps, "", "worktree", "end", wt); err != nil {
		t.Fatal(err)
	}
	if !containsLine(fake.Lines(), "git -C /repos/myrepo branch -D agent/foo") {
		t.Fatalf("expected agent branch delete, got %v", fake.Lines())
	}
}
