package worktree

import (
	"path/filepath"
	"strings"
)

// Names are the identifiers derived for a worktree from a repo path and suffix.
type Names struct {
	Name   string // sanitized "<repo-basename>-<suffix>"
	Branch string // "agent/<name>"
	Path   string // "<worktreesRoot>/<name>"
}

// Derive computes the worktree name, branch and path. worktreesRoot is the
// directory worktrees live under (the shell uses "$HOME/.cache/agent/worktrees").
func Derive(repoPath, suffix, worktreesRoot string) Names {
	name := Sanitize(filepath.Base(repoPath) + "-" + suffix)
	return Names{
		Name:   name,
		Branch: "agent/" + name,
		Path:   filepath.Join(worktreesRoot, name),
	}
}

// IsAgentBranch reports whether branch is an agent/* branch — the guard that
// stops `worktree end` from deleting a non-agent branch.
func IsAgentBranch(branch string) bool {
	return strings.HasPrefix(branch, "agent/")
}
