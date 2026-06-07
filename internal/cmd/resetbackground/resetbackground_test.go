package resetbackground_test

import (
	"testing"

	"github.com/dmccaffery/dotfiles/internal/cmd/cmdtest"
	"github.com/dmccaffery/dotfiles/internal/cmd/resetbackground"
	"github.com/dmccaffery/dotfiles/internal/ui"
)

func TestResetBackgroundItems(t *testing.T) {
	t.Run("declining the first confirm does nothing", func(t *testing.T) {
		deps := cmdtest.NewDeps(t)
		deps.Prompt = &ui.Fake{Replies: []bool{false}}
		if _, _, err := cmdtest.Run(t, resetbackground.NewCmd(deps), ""); err != nil {
			t.Fatal(err)
		}
		if cmdtest.ContainsLine(cmdtest.Fake(deps).Lines(), "sfltool") {
			t.Fatalf("must not reset when declined: %v", cmdtest.Fake(deps).Lines())
		}
	})

	t.Run("confirm resets; declining the reboot skips it", func(t *testing.T) {
		deps := cmdtest.NewDeps(t)
		deps.Prompt = &ui.Fake{Replies: []bool{true, false}}
		if _, _, err := cmdtest.Run(t, resetbackground.NewCmd(deps), ""); err != nil {
			t.Fatal(err)
		}
		fake := cmdtest.Fake(deps)
		if !cmdtest.ContainsLine(fake.Lines(), "sudo sfltool resetbtm") {
			t.Fatalf("expected reset, got %v", fake.Lines())
		}
		if cmdtest.ContainsLine(fake.Lines(), "reboot") {
			t.Fatalf("must not reboot when declined: %v", fake.Lines())
		}
	})

	t.Run("confirming both reboots", func(t *testing.T) {
		deps := cmdtest.NewDeps(t)
		deps.Prompt = &ui.Fake{Replies: []bool{true, true}}
		if _, _, err := cmdtest.Run(t, resetbackground.NewCmd(deps), ""); err != nil {
			t.Fatal(err)
		}
		if !cmdtest.ContainsLine(cmdtest.Fake(deps).Lines(), "sudo reboot") {
			t.Fatalf("expected reboot, got %v", cmdtest.Fake(deps).Lines())
		}
	})
}
