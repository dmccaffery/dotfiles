package ghswitch_test

import (
	"testing"

	"github.com/dmccaffery/dotfiles/internal/cmd/cmdtest"
	"github.com/dmccaffery/dotfiles/internal/cmd/ghswitch"
	"github.com/dmccaffery/dotfiles/internal/execx"
)

const ghStatusCmd = `gh auth status --active --json hosts --jq .hosts["github.com"][0].login`

func TestGhSwitchUser(t *testing.T) {
	t.Run("switches when the configured account is not active, then forwards", func(t *testing.T) {
		deps := cmdtest.NewDeps(t)
		fake := cmdtest.Fake(deps)
		fake.Responses["git config github.account"] = execx.Result{Stdout: "alice\n"}
		fake.Responses[ghStatusCmd] = execx.Result{Stdout: "bob\n"}
		if _, _, err := cmdtest.Run(t, ghswitch.NewCmd(deps), "", "pr", "list"); err != nil {
			t.Fatal(err)
		}
		if !cmdtest.ContainsLine(fake.Lines(), "gh auth switch --user alice") {
			t.Fatalf("expected account switch, got %v", fake.Lines())
		}
		if !cmdtest.ContainsLine(fake.Lines(), "gh pr list") {
			t.Fatalf("expected forward to gh, got %v", fake.Lines())
		}
	})

	t.Run("does not switch when already active", func(t *testing.T) {
		deps := cmdtest.NewDeps(t)
		fake := cmdtest.Fake(deps)
		fake.Responses["git config github.account"] = execx.Result{Stdout: "alice\n"}
		fake.Responses[ghStatusCmd] = execx.Result{Stdout: "alice\n"}
		if _, _, err := cmdtest.Run(t, ghswitch.NewCmd(deps), "", "repo", "view"); err != nil {
			t.Fatal(err)
		}
		if cmdtest.ContainsLine(fake.Lines(), "auth switch") {
			t.Fatalf("must not switch when already active: %v", fake.Lines())
		}
		if !cmdtest.ContainsLine(fake.Lines(), "gh repo view") {
			t.Fatalf("expected forward, got %v", fake.Lines())
		}
	})

	t.Run("no configured account just forwards", func(t *testing.T) {
		deps := cmdtest.NewDeps(t)
		fake := cmdtest.Fake(deps)
		if _, _, err := cmdtest.Run(t, ghswitch.NewCmd(deps), "", "auth", "status"); err != nil {
			t.Fatal(err)
		}
		if cmdtest.ContainsLine(fake.Lines(), "auth switch") {
			t.Fatalf("no switch expected: %v", fake.Lines())
		}
		if !cmdtest.ContainsLine(fake.Lines(), "gh auth status") {
			t.Fatalf("expected forward, got %v", fake.Lines())
		}
	})
}
