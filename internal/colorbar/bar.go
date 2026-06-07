// Package colorbar renders the 24-bit truecolor gradient bar that the
// print-colors command prints — a pure port of the original awk one-liner, so
// it is unit-testable.
package colorbar

import (
	"fmt"
	"strings"
)

// Gradient returns the rainbow gradient bar (a single line ending in a newline).
// Each of the 77 cells sets a 24-bit background, an inverted foreground, and an
// alternating "/"/"\" glyph; integer truncation matches awk's %d.
func Gradient() string {
	var b strings.Builder
	for col := range 77 {
		c := float64(col)
		r := 255 - c*255/76
		g := c * 510 / 76
		bl := c * 255 / 76
		if g > 255 {
			g = 510 - g
		}
		glyph := byte('/')
		if col%2 == 1 {
			glyph = '\\'
		}
		fmt.Fprintf(&b, "\x1b[48;2;%d;%d;%dm\x1b[38;2;%d;%d;%dm%c\x1b[0m",
			int(r), int(g), int(bl), int(255-r), int(255-g), int(255-bl), glyph)
	}
	b.WriteByte('\n')
	return b.String()
}
