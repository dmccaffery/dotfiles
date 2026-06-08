package root_test

import (
	"strings"
	"testing"

	"github.com/dmccaffery/dotfiles/internal/cmd/cmdtest"
	"github.com/dmccaffery/dotfiles/internal/cmd/root"
)

func TestHasCommand(t *testing.T) {
	r := root.NewRootCmd("test", cmdtest.NewDeps(t))
	for _, n := range []string{"worktree", "agent-tmux-status", "brewfile", "zs"} {
		if !root.HasCommand(r, n) {
			t.Errorf("HasCommand(%q) = false, want true (argv[0] dispatch would miss it)", n)
		}
	}
	// "dot" and unknown names must NOT dispatch as applets, so an oddly-named
	// binary falls through to plain `dot` instead of erroring on an unknown command.
	for _, n := range []string{"dot", "bogus", "dotfinal", ""} {
		if root.HasCommand(r, n) {
			t.Errorf("HasCommand(%q) = true, want false", n)
		}
	}
}

func TestAppletsListsSymlinkNames(t *testing.T) {
	out, _, err := cmdtest.Run(t, root.NewRootCmd("test", cmdtest.NewDeps(t)), "", "applets")
	if err != nil {
		t.Fatal(err)
	}
	got := strings.Fields(out)
	want := map[string]bool{
		"worktree": true, "agent-tmux-status": true, "brewfile": true,
		"print-colors": true, "profile-shell": true, "git-resign": true,
		"gh-switch-user": true, "fzf-image-preview": true,
		"reset-background-items": true, "zs": true,
		"git-github-auth": true, "tmux-session": true,
	}
	if len(got) != len(want) {
		t.Fatalf("applets = %v, want keys %v", got, want)
	}
	for _, n := range got {
		if !want[n] {
			t.Errorf("unexpected applet %q", n)
		}
	}
}
