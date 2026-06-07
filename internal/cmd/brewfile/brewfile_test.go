package brewfile_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dmccaffery/dotfiles/internal/cmd/brewfile"
	"github.com/dmccaffery/dotfiles/internal/cmd/cmdtest"
	"github.com/dmccaffery/dotfiles/internal/envx"
	"github.com/dmccaffery/dotfiles/internal/ui"
)

func TestBrewfileHelpRendersCobra(t *testing.T) {
	// DisableFlagParsing turns off Cobra's automatic --help, so run() routes it to
	// cmd.Help(); assert it renders the Use/Example metadata rather than a bespoke string.
	out, _, err := cmdtest.Run(t, brewfile.NewCmd(cmdtest.NewDeps(t)), "", "--help")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"Usage:", "Examples:", "brewfile add --cask ghostty"} {
		if !strings.Contains(out, want) {
			t.Errorf("brewfile --help missing %q in:\n%s", want, out)
		}
	}
}

func TestBrewfileTrustPrompt(t *testing.T) {
	add := []string{"add", "--cask", "user/tap/foo"}

	t.Run("yes runs brew trust", func(t *testing.T) {
		deps := cmdtest.NewDeps(t)
		prompt := &ui.Fake{Replies: []bool{true}}
		deps.Prompt = prompt
		if _, _, err := cmdtest.Run(t, brewfile.NewCmd(deps), "", add...); err != nil {
			t.Fatal(err)
		}
		if len(prompt.Asked) != 1 {
			t.Fatalf("expected exactly one confirm, got %v", prompt.Asked)
		}
		if !cmdtest.ContainsLine(cmdtest.Fake(deps).Lines(), "brew trust --cask user/tap/foo") {
			t.Fatalf("expected brew trust call, got %v", cmdtest.Fake(deps).Lines())
		}
	})

	t.Run("no skips brew trust", func(t *testing.T) {
		deps := cmdtest.NewDeps(t)
		deps.Prompt = &ui.Fake{Replies: []bool{false}}
		if _, _, err := cmdtest.Run(t, brewfile.NewCmd(deps), "", add...); err != nil {
			t.Fatal(err)
		}
		if cmdtest.ContainsLine(cmdtest.Fake(deps).Lines(), "brew trust") {
			t.Fatalf("must not trust when declined: %v", cmdtest.Fake(deps).Lines())
		}
	})

	t.Run("no tty skips with warning", func(t *testing.T) {
		deps := cmdtest.NewDeps(t)
		deps.Prompt = &ui.Fake{NoTTY: true}
		if _, _, err := cmdtest.Run(t, brewfile.NewCmd(deps), "", add...); err != nil {
			t.Fatal(err)
		}
		if cmdtest.ContainsLine(cmdtest.Fake(deps).Lines(), "brew trust") {
			t.Fatalf("must not trust without a tty: %v", cmdtest.Fake(deps).Lines())
		}
	})

	t.Run("already trusted does not prompt", func(t *testing.T) {
		home := t.TempDir()
		cfg := filepath.Join(home, ".config", "homebrew")
		if err := os.MkdirAll(cfg, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(cfg, "trust.json"), []byte(`{"trustedcasks":["user/tap/foo"]}`), 0o644); err != nil {
			t.Fatal(err)
		}
		deps := cmdtest.NewDeps(t)
		deps.Env = envx.New(home, map[string]string{"XDG_CONFIG_HOME": filepath.Join(home, ".config")})
		prompt := &ui.Fake{}
		deps.Prompt = prompt
		if _, _, err := cmdtest.Run(t, brewfile.NewCmd(deps), "", add...); err != nil {
			t.Fatal(err)
		}
		if len(prompt.Asked) != 0 {
			t.Fatalf("a trusted entry must not prompt, got %v", prompt.Asked)
		}
	})
}

func TestBrewfileAddRunsBundleThenInstall(t *testing.T) {
	deps := cmdtest.NewDeps(t)
	if _, _, err := cmdtest.Run(t, brewfile.NewCmd(deps), "", "add", "jq"); err != nil {
		t.Fatal(err)
	}
	lines := cmdtest.Fake(deps).Lines()
	if !cmdtest.ContainsLine(lines, "brew bundle add jq --global") {
		t.Fatalf("missing bundle add, got %v", lines)
	}
	if !cmdtest.ContainsLine(lines, "brew bundle install --global --zap --upgrade") {
		t.Fatalf("missing bundle install, got %v", lines)
	}
}
