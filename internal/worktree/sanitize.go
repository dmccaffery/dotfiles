// Package worktree holds the pure logic of the agent-worktree lifecycle —
// name sanitization, identifier derivation and hook-JSON parsing — with no I/O,
// so it is exhaustively unit-testable. The git side effects live in the cli
// command that wraps this package.
package worktree

import "strings"

// Sanitize encodes a string into a tmux-safe session/worktree name, matching
// the original shell exactly:
//
//	sed -e 's/^\./dot-/' -e 's/\.$/-dot/' -e 's/\./-dot-/g' | tr -c 'A-Za-z0-9_-' '-'
//
// tmux >=3.5 rejects '.' in session names, so dots are encoded (leading -> "dot-",
// trailing -> "-dot", interior -> "-dot-") before every other disallowed
// character collapses to '-'.
func Sanitize(s string) string {
	if strings.HasPrefix(s, ".") {
		s = "dot-" + s[1:]
	}
	if strings.HasSuffix(s, ".") {
		s = s[:len(s)-1] + "-dot"
	}
	s = strings.ReplaceAll(s, ".", "-dot-")

	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if isAllowed(r) {
			b.WriteRune(r)
		} else {
			b.WriteByte('-')
		}
	}
	return b.String()
}

func isAllowed(r rune) bool {
	switch {
	case r >= 'a' && r <= 'z',
		r >= 'A' && r <= 'Z',
		r >= '0' && r <= '9',
		r == '_', r == '-':
		return true
	default:
		return false
	}
}
