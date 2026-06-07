package printcolors_test

import (
	"strings"
	"testing"

	"github.com/dmccaffery/dotfiles/internal/cmd/cmdtest"
	"github.com/dmccaffery/dotfiles/internal/cmd/printcolors"
)

func TestPrintColors(t *testing.T) {
	out, _, err := cmdtest.Run(t, printcolors.NewCmd(cmdtest.NewDeps(t)), "")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(out, "\x1b[48;2;255;0;0m") {
		t.Fatalf("expected a truecolor gradient, got %q", out)
	}
	if !strings.HasSuffix(out, "\x1b[0m\n") {
		t.Fatal("expected a trailing reset + newline")
	}
}
