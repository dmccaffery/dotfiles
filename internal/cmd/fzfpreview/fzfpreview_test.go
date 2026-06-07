package fzfpreview_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dmccaffery/dotfiles/internal/cmd/cmdtest"
	"github.com/dmccaffery/dotfiles/internal/cmd/fzfpreview"
	"github.com/dmccaffery/dotfiles/internal/envx"
	"github.com/dmccaffery/dotfiles/internal/execx"
)

func mustWrite(t *testing.T, name, content string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestFzfImagePreview(t *testing.T) {
	t.Run("directory lists", func(t *testing.T) {
		dir := t.TempDir()
		deps := cmdtest.NewDeps(t)
		if _, _, err := cmdtest.Run(t, fzfpreview.NewCmd(deps), "", dir); err != nil {
			t.Fatal(err)
		}
		if !cmdtest.ContainsLine(cmdtest.Fake(deps).Lines(), "ls -la --color "+dir) {
			t.Fatalf("expected directory listing, got %v", cmdtest.Fake(deps).Lines())
		}
	})

	t.Run("missing file errors", func(t *testing.T) {
		out, _, err := cmdtest.Run(t, fzfpreview.NewCmd(cmdtest.NewDeps(t)), "", "/no/such/file")
		if err == nil {
			t.Fatal("expected an error for a missing file")
		}
		if !strings.Contains(out, "file does not exist") {
			t.Fatalf("out = %q", out)
		}
	})

	t.Run("image renders with chafa", func(t *testing.T) {
		f := mustWrite(t, "x.png", "png")
		deps := cmdtest.NewDeps(t)
		deps.Env = envx.New(t.TempDir(), map[string]string{"FZF_PREVIEW_COLUMNS": "80"})
		cmdtest.Fake(deps).Responses["file --mime "+f] = execx.Result{Stdout: f + ": image/png; charset=binary"}
		if _, _, err := cmdtest.Run(t, fzfpreview.NewCmd(deps), "", f); err != nil {
			t.Fatal(err)
		}
		if !cmdtest.ContainsLine(cmdtest.Fake(deps).Lines(), "chafa --passthrough=auto --size=80 "+f) {
			t.Fatalf("expected chafa render, got %v", cmdtest.Fake(deps).Lines())
		}
	})

	t.Run("non-image binary reports its type", func(t *testing.T) {
		f := mustWrite(t, "a.out", "ELF")
		deps := cmdtest.NewDeps(t)
		cmdtest.Fake(deps).Responses["file --mime "+f] = execx.Result{Stdout: f + ": application/octet-stream; charset=binary"}
		out, _, err := cmdtest.Run(t, fzfpreview.NewCmd(deps), "", f)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(out, f+" is a binary file") {
			t.Fatalf("out = %q", out)
		}
	})

	t.Run("text renders with bat", func(t *testing.T) {
		f := mustWrite(t, "t.txt", "hello")
		deps := cmdtest.NewDeps(t)
		fake := cmdtest.Fake(deps)
		fake.Responses["file --mime "+f] = execx.Result{Stdout: f + ": text/plain; charset=us-ascii"}
		fake.Responses["bat --style=numbers --color=always "+f] = execx.Result{Stdout: "   1 hello\n"}
		out, _, err := cmdtest.Run(t, fzfpreview.NewCmd(deps), "", f)
		if err != nil {
			t.Fatal(err)
		}
		if !cmdtest.ContainsLine(fake.Lines(), "bat --style=numbers --color=always "+f) {
			t.Fatalf("expected bat render, got %v", fake.Lines())
		}
		if !strings.Contains(out, "1 hello") {
			t.Fatalf("out = %q", out)
		}
	})
}
