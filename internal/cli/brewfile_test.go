package cli_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dmccaffery/dotfiles/internal/cli"
	"github.com/dmccaffery/dotfiles/internal/envx"
	"github.com/dmccaffery/dotfiles/internal/execx"
	"github.com/dmccaffery/dotfiles/internal/ui"
)

func TestBrewfileHelpRendersCobra(t *testing.T) {
	// DisableFlagParsing turns off Cobra's automatic --help, so run() routes it to
	// cmd.Help(); assert it renders the Use/Example metadata rather than a bespoke string.
	deps := &cli.Deps{Runner: execx.NewFake(), Env: envx.New(t.TempDir(), nil)}
	out, _, err := runRoot(t, deps, "", "brewfile", "--help")
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
	// add of an untrusted, non-official tap reference triggers a Confirm; the
	// answer decides whether `brew trust` runs. The injected Fake prompter makes
	// the whole flow testable without a tty.
	add := []string{"brewfile", "add", "--cask", "user/tap/foo"}

	t.Run("yes runs brew trust", func(t *testing.T) {
		fake := execx.NewFake()
		prompt := &ui.Fake{Replies: []bool{true}}
		deps := &cli.Deps{Runner: fake, Env: envx.New(t.TempDir(), nil), Prompt: prompt}
		if _, _, err := runRoot(t, deps, "", add...); err != nil {
			t.Fatal(err)
		}
		if len(prompt.Asked) != 1 {
			t.Fatalf("expected exactly one confirm, got %v", prompt.Asked)
		}
		if !containsLine(fake.Lines(), "brew trust --cask user/tap/foo") {
			t.Fatalf("expected brew trust call, got %v", fake.Lines())
		}
	})

	t.Run("no skips brew trust", func(t *testing.T) {
		fake := execx.NewFake()
		deps := &cli.Deps{Runner: fake, Env: envx.New(t.TempDir(), nil), Prompt: &ui.Fake{Replies: []bool{false}}}
		if _, _, err := runRoot(t, deps, "", add...); err != nil {
			t.Fatal(err)
		}
		if containsLine(fake.Lines(), "brew trust") {
			t.Fatalf("must not trust when declined: %v", fake.Lines())
		}
	})

	t.Run("no tty skips with warning", func(t *testing.T) {
		fake := execx.NewFake()
		deps := &cli.Deps{Runner: fake, Env: envx.New(t.TempDir(), nil), Prompt: &ui.Fake{NoTTY: true}}
		if _, _, err := runRoot(t, deps, "", add...); err != nil {
			t.Fatal(err)
		}
		if containsLine(fake.Lines(), "brew trust") {
			t.Fatalf("must not trust without a tty: %v", fake.Lines())
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
		prompt := &ui.Fake{}
		env := envx.New(home, map[string]string{"XDG_CONFIG_HOME": filepath.Join(home, ".config")})
		deps := &cli.Deps{Runner: execx.NewFake(), Env: env, Prompt: prompt}
		if _, _, err := runRoot(t, deps, "", add...); err != nil {
			t.Fatal(err)
		}
		if len(prompt.Asked) != 0 {
			t.Fatalf("a trusted entry must not prompt, got %v", prompt.Asked)
		}
	})
}

func TestBrewfileAddRunsBundleThenInstall(t *testing.T) {
	fake := execx.NewFake()
	// no XDG_CONFIG_HOME → trust file is ~/.homebrew/trust.json (absent); jq is a
	// bare name so no trust prompt is reached.
	deps := &cli.Deps{Runner: fake, Env: envx.New(t.TempDir(), nil)}
	if _, _, err := runRoot(t, deps, "", "brewfile", "add", "jq"); err != nil {
		t.Fatal(err)
	}
	lines := fake.Lines()
	if !containsLine(lines, "brew bundle add jq --global") {
		t.Fatalf("missing bundle add, got %v", lines)
	}
	if !containsLine(lines, "brew bundle install --global --zap --upgrade") {
		t.Fatalf("missing bundle install, got %v", lines)
	}
}
