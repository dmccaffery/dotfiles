package worktree

import "testing"

func TestSanitize(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"dotfiles", "dotfiles"},
		{"my-repo_9", "my-repo_9"},
		{"UPPER", "UPPER"},
		{".config", "dot-config"},            // leading dot
		{"build.", "build-dot"},              // trailing dot
		{"a.b", "a-dot-b"},                   // interior dot
		{"a.b.c", "a-dot-b-dot-c"},           // multiple interior dots
		{".", "dot-"},                        // lone dot
		{"..", "dot--dot"},                   // two dots
		{".a.", "dot-a-dot"},                 // leading + trailing
		{"feature/branch", "feature-branch"}, // slash collapses
		{"a b", "a-b"},                       // space collapses
		{"café", "caf-"},                     // non-ASCII rune collapses
	}
	for _, c := range cases {
		if got := Sanitize(c.in); got != c.want {
			t.Errorf("Sanitize(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}
