package worktree

import "testing"

func TestDerive(t *testing.T) {
	got := Derive("/home/user/myrepo", "feature", "/wt")
	if got.Name != "myrepo-feature" {
		t.Errorf("Name = %q", got.Name)
	}
	if got.Branch != "agent/myrepo-feature" {
		t.Errorf("Branch = %q", got.Branch)
	}
	if got.Path != "/wt/myrepo-feature" {
		t.Errorf("Path = %q", got.Path)
	}
}

func TestDeriveSanitizesDots(t *testing.T) {
	got := Derive("/x/my.repo", "feat", "/wt")
	if got.Name != "my-dot-repo-feat" {
		t.Errorf("Name = %q, want my-dot-repo-feat", got.Name)
	}
	if got.Branch != "agent/my-dot-repo-feat" {
		t.Errorf("Branch = %q", got.Branch)
	}
}

func TestIsAgentBranch(t *testing.T) {
	cases := map[string]bool{
		"agent/foo":     true,
		"agent/a/b":     true,
		"main":          false,
		"agentfoo":      false,
		"":              false,
		"feature/agent": false,
	}
	for in, want := range cases {
		if got := IsAgentBranch(in); got != want {
			t.Errorf("IsAgentBranch(%q) = %v, want %v", in, got, want)
		}
	}
}
