package gitresign_test

import (
	"testing"

	"github.com/dmccaffery/dotfiles/internal/cmd/cmdtest"
	"github.com/dmccaffery/dotfiles/internal/cmd/gitresign"
)

func TestGitResign(t *testing.T) {
	t.Run("missing target errors", func(t *testing.T) {
		if _, _, err := cmdtest.Run(t, gitresign.NewCmd(cmdtest.NewDeps(t)), ""); err == nil {
			t.Fatal("expected an error for a missing target")
		}
	})

	t.Run("runs interactive rebase that re-signs each commit", func(t *testing.T) {
		deps := cmdtest.NewDeps(t)
		if _, _, err := cmdtest.Run(t, gitresign.NewCmd(deps), "", "HEAD~3"); err != nil {
			t.Fatal(err)
		}
		if !cmdtest.ContainsLine(cmdtest.Fake(deps).Lines(), "git rebase --exec git commit --amend --no-edit -n -S -i HEAD~3") {
			t.Fatalf("unexpected rebase invocation: %v", cmdtest.Fake(deps).Lines())
		}
	})
}
