package cli_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/dmccaffery/dotfiles/internal/cli"
	"github.com/dmccaffery/dotfiles/internal/envx"
	"github.com/dmccaffery/dotfiles/internal/logx"
)

// runRoot drives the dot root command with the given stdin and args, capturing
// stdout and stderr separately so the "path on stdout only" contract is testable.
// It is shared by the per-command test files in this package.
func runRoot(t *testing.T, deps *cli.Deps, stdin string, args ...string) (stdout, stderr string, err error) {
	t.Helper()
	if deps.Log == nil {
		deps.Log = logx.For(new(bytes.Buffer), false)
	}
	root := cli.NewRootCmd("test", deps)
	var out, errb bytes.Buffer
	root.SetArgs(args)
	root.SetIn(strings.NewReader(stdin))
	root.SetOut(&out)
	root.SetErr(&errb)
	err = root.Execute()
	return out.String(), errb.String(), err
}

func containsLine(lines []string, sub string) bool {
	for _, l := range lines {
		if strings.Contains(l, sub) {
			return true
		}
	}
	return false
}

func TestHasCommand(t *testing.T) {
	root := cli.NewRootCmd("test", &cli.Deps{Env: envx.New(t.TempDir(), nil)})
	for _, n := range []string{"worktree", "agent-tmux-status", "brewfile"} {
		if !cli.HasCommand(root, n) {
			t.Errorf("HasCommand(%q) = false, want true (argv[0] dispatch would miss it)", n)
		}
	}
	// "dot" and unknown names must NOT dispatch as applets, so an oddly-named
	// binary falls through to plain `dot` instead of erroring on an unknown command.
	for _, n := range []string{"dot", "bogus", "dotfinal", ""} {
		if cli.HasCommand(root, n) {
			t.Errorf("HasCommand(%q) = true, want false", n)
		}
	}
}

func TestAppletsListsSymlinkNames(t *testing.T) {
	out, _, err := runRoot(t, &cli.Deps{Env: envx.New(t.TempDir(), nil)}, "", "applets")
	if err != nil {
		t.Fatal(err)
	}
	got := strings.Fields(out)
	want := map[string]bool{"worktree": true, "agent-tmux-status": true, "brewfile": true}
	if len(got) != len(want) {
		t.Fatalf("applets = %v, want keys %v", got, want)
	}
	for _, n := range got {
		if !want[n] {
			t.Errorf("unexpected applet %q", n)
		}
	}
}
