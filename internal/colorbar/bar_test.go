package colorbar

import (
	"strings"
	"testing"
)

func TestGradient(t *testing.T) {
	g := Gradient()
	// 77 cells, each terminated by a reset; whole bar ends with reset + newline.
	if n := strings.Count(g, "\x1b[0m"); n != 77 {
		t.Fatalf("expected 77 reset codes, got %d", n)
	}
	if !strings.HasSuffix(g, "\x1b[0m\n") {
		t.Fatalf("must end with reset + newline")
	}
	// col 0: bg=(255,0,0), inverted fg=(0,255,255), glyph '/'.
	const wantFirst = "\x1b[48;2;255;0;0m\x1b[38;2;0;255;255m/\x1b[0m"
	if !strings.HasPrefix(g, wantFirst) {
		t.Fatalf("unexpected first cell:\n got %q\nwant prefix %q", g, wantFirst)
	}
	// glyphs alternate /\ ; cell 1 uses a backslash.
	if !strings.Contains(g, "m\\\x1b[0m") {
		t.Fatalf("expected an alternating backslash glyph")
	}
}
