package cmdutil_test

import (
	"testing"

	"github.com/dmccaffery/dotfiles/internal/cmd/cmdutil"
	"github.com/dmccaffery/dotfiles/internal/ui"
)

func TestPickOne(t *testing.T) {
	opts := []string{"dotfiles", "dotfiles-docs", "other"}

	t.Run("a query with a single match auto-selects without prompting", func(t *testing.T) {
		p := &ui.Fake{}
		got, err := cmdutil.PickOne(p, "repo", "other", opts)
		if err != nil || got != "other" {
			t.Fatalf("got %q err %v", got, err)
		}
		if len(p.Asked) != 0 {
			t.Fatal("a single match must not prompt")
		}
	})

	t.Run("a query with no match returns empty without prompting", func(t *testing.T) {
		p := &ui.Fake{}
		got, err := cmdutil.PickOne(p, "repo", "zzz", opts)
		if err != nil || got != "" {
			t.Fatalf("got %q err %v", got, err)
		}
		if len(p.Asked) != 0 {
			t.Fatal("no match must not prompt")
		}
	})

	t.Run("a query with multiple matches prompts over the matches", func(t *testing.T) {
		p := &ui.Fake{Selections: []string{"dotfiles-docs"}}
		got, err := cmdutil.PickOne(p, "repo", "dotfiles", opts)
		if err != nil || got != "dotfiles-docs" {
			t.Fatalf("got %q err %v", got, err)
		}
		if len(p.Asked) != 1 {
			t.Fatalf("multiple matches should prompt once, asked %v", p.Asked)
		}
	})

	t.Run("no query prompts over all options", func(t *testing.T) {
		p := &ui.Fake{Selections: []string{"dotfiles"}}
		got, err := cmdutil.PickOne(p, "repo", "", opts)
		if err != nil || got != "dotfiles" {
			t.Fatalf("got %q err %v", got, err)
		}
	})
}
