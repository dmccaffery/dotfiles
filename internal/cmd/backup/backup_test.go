package backup_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dmccaffery/dotfiles/internal/cmd/backup"
	"github.com/dmccaffery/dotfiles/internal/cmd/cmdtest"
	"github.com/dmccaffery/dotfiles/internal/cmd/cmdutil"
	"github.com/dmccaffery/dotfiles/internal/ui"
)

// env bundles the three directories a backup run spans: the stow source, the
// $HOME it mirrors against, and where backups land.
type env struct {
	deps   *cmdutil.Deps
	stow   string
	home   string
	backup string
}

func newEnv(t *testing.T) env {
	t.Helper()
	deps := cmdtest.NewDeps(t)
	deps.Prompt = &ui.Fake{} // not consulted with --yes
	return env{
		deps:   deps,
		stow:   filepath.Join(t.TempDir(), "stow"),
		home:   deps.Env.Home(),
		backup: filepath.Join(t.TempDir(), "backups", "snap"),
	}
}

func (e env) run(t *testing.T, extra ...string) error {
	t.Helper()
	args := append([]string{"--stow-dir", e.stow, "--backup-dir", e.backup}, extra...)
	_, _, err := cmdtest.Run(t, backup.NewCmd(e.deps), "", args...)
	return err
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func symlink(t *testing.T, target, link string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(link), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}
}

func exists(path string) bool {
	_, err := os.Lstat(path)
	return err == nil
}

func read(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestBackup(t *testing.T) {
	t.Run("regular file is moved into the backup", func(t *testing.T) {
		e := newEnv(t)
		writeFile(t, filepath.Join(e.stow, ".config/app/conf"), "src")
		writeFile(t, filepath.Join(e.home, ".config/app/conf"), "user-data")

		if err := e.run(t, "--yes"); err != nil {
			t.Fatal(err)
		}
		if exists(filepath.Join(e.home, ".config/app/conf")) {
			t.Fatal("target file should have been moved out of $HOME")
		}
		if got := read(t, filepath.Join(e.backup, ".config/app/conf")); got != "user-data" {
			t.Fatalf("backup content = %q, want user-data", got)
		}
	})

	t.Run("broken symlink is removed, not backed up", func(t *testing.T) {
		e := newEnv(t)
		writeFile(t, filepath.Join(e.stow, ".cache/foo"), "src")
		symlink(t, filepath.Join(e.home, "does-not-exist"), filepath.Join(e.home, ".cache/foo"))

		if err := e.run(t, "--yes"); err != nil {
			t.Fatal(err)
		}
		if exists(filepath.Join(e.home, ".cache/foo")) {
			t.Fatal("broken symlink should have been removed")
		}
		if exists(filepath.Join(e.backup, ".cache/foo")) {
			t.Fatal("broken symlink must not be backed up")
		}
	})

	t.Run("valid symlink has its contents backed up, then the link is removed", func(t *testing.T) {
		e := newEnv(t)
		external := filepath.Join(t.TempDir(), "real.txt")
		writeFile(t, external, "pointed-to")
		writeFile(t, filepath.Join(e.stow, ".ssh/rc"), "src")
		symlink(t, external, filepath.Join(e.home, ".ssh/rc"))

		if err := e.run(t, "--yes"); err != nil {
			t.Fatal(err)
		}
		if exists(filepath.Join(e.home, ".ssh/rc")) {
			t.Fatal("symlink should have been removed")
		}
		if !exists(external) {
			t.Fatal("a valid symlink's target must not be deleted")
		}
		if got := read(t, filepath.Join(e.backup, ".ssh/rc")); got != "pointed-to" {
			t.Fatalf("backup content = %q, want pointed-to", got)
		}
	})

	t.Run("a folded parent symlink is handled at the link, not descended", func(t *testing.T) {
		e := newEnv(t)
		// Two stow files share a parent; $HOME/.config/ghostty is one symlink to
		// an external dir (a "folded" stow link).
		writeFile(t, filepath.Join(e.stow, ".config/ghostty/config"), "src")
		writeFile(t, filepath.Join(e.stow, ".config/ghostty/theme"), "src")
		realDir := filepath.Join(t.TempDir(), "ghostty")
		writeFile(t, filepath.Join(realDir, "config"), "cfg")
		writeFile(t, filepath.Join(realDir, "theme"), "thm")
		// $HOME/.config must be a real dir so the walk descends into it.
		if err := os.MkdirAll(filepath.Join(e.home, ".config"), 0o755); err != nil {
			t.Fatal(err)
		}
		symlink(t, realDir, filepath.Join(e.home, ".config/ghostty"))

		if err := e.run(t, "--yes"); err != nil {
			t.Fatal(err)
		}
		if exists(filepath.Join(e.home, ".config/ghostty")) {
			t.Fatal("folded parent symlink should have been removed")
		}
		if !exists(realDir) {
			t.Fatal("the symlink's target directory must not be deleted")
		}
		if got := read(t, filepath.Join(e.backup, ".config/ghostty/config")); got != "cfg" {
			t.Fatalf("backup config = %q, want cfg", got)
		}
		if got := read(t, filepath.Join(e.backup, ".config/ghostty/theme")); got != "thm" {
			t.Fatalf("backup theme = %q, want thm", got)
		}
	})

	t.Run("a symlink in the stow source is not skipped; its target is cleared", func(t *testing.T) {
		e := newEnv(t)
		// stow/.config/codex/AGENTS.md is itself a symlink (-> the canonical
		// CLAUDE.md), exactly like the real repo. stow links it, so its $HOME
		// target must still be backed up.
		writeFile(t, filepath.Join(e.stow, ".claude/CLAUDE.md"), "canon")
		symlink(t, "../../.claude/CLAUDE.md", filepath.Join(e.stow, ".config/codex/AGENTS.md"))
		writeFile(t, filepath.Join(e.home, ".config/codex/AGENTS.md"), "user-agents")

		if err := e.run(t, "--yes"); err != nil {
			t.Fatal(err)
		}
		if exists(filepath.Join(e.home, ".config/codex/AGENTS.md")) {
			t.Fatal("target of a symlinked stow entry must not be skipped")
		}
		if got := read(t, filepath.Join(e.backup, ".config/codex/AGENTS.md")); got != "user-agents" {
			t.Fatalf("backup content = %q, want user-agents", got)
		}
	})

	t.Run("a missing target is left alone", func(t *testing.T) {
		e := newEnv(t)
		writeFile(t, filepath.Join(e.stow, ".config/missing/x"), "src")

		if err := e.run(t, "--yes"); err != nil {
			t.Fatal(err)
		}
		if exists(e.backup) {
			t.Fatal("nothing existed at the target; no backup dir should be created")
		}
	})

	t.Run("declining the confirmation leaves $HOME untouched", func(t *testing.T) {
		e := newEnv(t)
		fake := &ui.Fake{Replies: []bool{false}}
		e.deps.Prompt = fake
		writeFile(t, filepath.Join(e.stow, ".config/app/conf"), "src")
		writeFile(t, filepath.Join(e.home, ".config/app/conf"), "user-data")

		if err := e.run(t); err != nil { // no --yes -> prompts
			t.Fatal(err)
		}
		if !exists(filepath.Join(e.home, ".config/app/conf")) {
			t.Fatal("declining must not move anything")
		}
		if exists(e.backup) {
			t.Fatal("declining must not create a backup")
		}
		if len(fake.Asked) == 0 {
			t.Fatal("expected a confirmation prompt")
		}
	})
}
