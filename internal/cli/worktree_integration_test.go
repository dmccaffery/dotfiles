package cli_test

// These tests run the worktree command against REAL git in a temp repo, proving
// that worktrees and agent/* branches are actually created and removed (not just
// that the right commands were issued). They skip when git is unavailable; CI
// always has it. GIT_CONFIG_GLOBAL/SYSTEM are pinned to /dev/null so the run is
// hermetic (no signing config, no user settings leaking in).

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dmccaffery/dotfiles/internal/cli"
	"github.com/dmccaffery/dotfiles/internal/envx"
	"github.com/dmccaffery/dotfiles/internal/execx"
)

func gitRepo(t *testing.T) string {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not on PATH")
	}
	t.Setenv("GIT_CONFIG_GLOBAL", os.DevNull)
	t.Setenv("GIT_CONFIG_SYSTEM", os.DevNull)
	repo := t.TempDir()
	mustGit(t, repo, "init", "-q")
	mustGit(t, repo, "-c", "user.email=t@example.com", "-c", "user.name=t", "commit", "-q", "--allow-empty", "-m", "init")
	return repo
}

func mustGit(t *testing.T, repo string, args ...string) string {
	t.Helper()
	out, err := exec.Command("git", append([]string{"-C", repo}, args...)...).CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

func realDeps(home string) *cli.Deps {
	return &cli.Deps{Runner: execx.Real{}, Env: envx.New(home, nil)}
}

func TestWorktreeIntegrationLifecycle(t *testing.T) {
	repo := gitRepo(t)
	home := t.TempDir()

	out, _, err := runRoot(t, realDeps(home), "", "worktree", "start", repo, "feature")
	if err != nil {
		t.Fatal(err)
	}
	path := strings.TrimSpace(out)
	wantPath := filepath.Join(home, ".cache", "agent", "worktrees", filepath.Base(repo)+"-feature")
	if path != wantPath {
		t.Fatalf("stdout path = %q, want %q", path, wantPath)
	}
	if fi, err := os.Stat(path); err != nil || !fi.IsDir() {
		t.Fatalf("worktree dir was not created at %s: %v", path, err)
	}
	branch := "agent/" + filepath.Base(path)
	if mustGit(t, repo, "branch", "--list", branch) == "" {
		t.Fatalf("branch %s was not created", branch)
	}

	// Reuse: a second start of the same name returns the same path.
	out2, _, err := runRoot(t, realDeps(home), "", "worktree", "start", repo, "feature")
	if err != nil || strings.TrimSpace(out2) != path {
		t.Fatalf("reuse returned %q (err %v), want %q", strings.TrimSpace(out2), err, path)
	}

	// End: worktree dir and agent/* branch are both removed.
	if _, _, err := runRoot(t, realDeps(home), "", "worktree", "end", path); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("worktree dir still present after end: %v", err)
	}
	if got := mustGit(t, repo, "branch", "--list", branch); got != "" {
		t.Fatalf("agent branch %s not deleted (got %q)", branch, got)
	}
}

func TestWorktreeIntegrationStartFromStdinJSON(t *testing.T) {
	repo := gitRepo(t)
	home := t.TempDir()
	out, _, err := runRoot(t, realDeps(home), `{"name":"hookwt"}`, "worktree", "start", repo)
	if err != nil {
		t.Fatal(err)
	}
	path := strings.TrimSpace(out)
	if !strings.HasSuffix(path, "-hookwt") {
		t.Fatalf("expected suffix from stdin JSON .name, got %q", path)
	}
	if fi, err := os.Stat(path); err != nil || !fi.IsDir() {
		t.Fatalf("worktree from hook JSON not created: %v", err)
	}
}

func TestWorktreeIntegrationKeepsNonAgentBranch(t *testing.T) {
	repo := gitRepo(t)
	home := t.TempDir()
	wt := filepath.Join(t.TempDir(), "plain")
	mustGit(t, repo, "worktree", "add", "-q", "-b", "feature/plain", wt)

	if _, _, err := runRoot(t, realDeps(home), "", "worktree", "end", wt); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(wt); !os.IsNotExist(err) {
		t.Fatalf("worktree should be removed: %v", err)
	}
	if mustGit(t, repo, "branch", "--list", "feature/plain") == "" {
		t.Fatalf("non-agent branch feature/plain must be preserved")
	}
}
