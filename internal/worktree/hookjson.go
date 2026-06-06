package worktree

import "encoding/json"

// ParseStartName extracts ".name" from WorktreeCreate hook JSON on stdin,
// replacing the shell's `jq -r '.name // empty'`. Invalid/empty input yields "".
func ParseStartName(data []byte) string {
	var h struct {
		Name string `json:"name"`
	}
	_ = json.Unmarshal(data, &h)
	return h.Name
}

// ParseEndPath extracts ".worktree_path" from WorktreeRemove hook JSON on stdin,
// replacing `jq -r '.worktree_path // empty'`. Invalid/empty input yields "".
func ParseEndPath(data []byte) string {
	var h struct {
		WorktreePath string `json:"worktree_path"`
	}
	_ = json.Unmarshal(data, &h)
	return h.WorktreePath
}
