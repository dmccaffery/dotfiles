package cli

import (
	"errors"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/dmccaffery/dotfiles/internal/logx"
)

// ErrSilent signals a non-zero exit after the command has already logged its own
// diagnostics. main suppresses printing it (cobra's other errors still print).
var ErrSilent = errors.New("dot: command failed")

// readIfPiped returns stdin's contents when it is not a terminal (i.e. a hook
// pipe), matching the shell's `[ ! -t 0 ]` guard. For a tty it returns false.
func readIfPiped(in io.Reader) ([]byte, bool) {
	if f, ok := in.(*os.File); ok && logx.IsTerminal(f) {
		return nil, false
	}
	data, err := io.ReadAll(in)
	if err != nil || len(data) == 0 {
		return nil, false
	}
	return data, true
}

func dirExists(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.IsDir()
}

// countNonEmptyLines mirrors `... | wc -l` over git --porcelain output.
func countNonEmptyLines(s string) int {
	n := 0
	for line := range strings.SplitSeq(s, "\n") {
		if strings.TrimSpace(line) != "" {
			n++
		}
	}
	return n
}

func atoi(s string) int {
	n, _ := strconv.Atoi(strings.TrimSpace(s))
	return n
}
