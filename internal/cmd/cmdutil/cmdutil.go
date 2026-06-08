// Package cmdutil holds the shared command framework for the dot subcommands:
// the injectable Deps, the silent-exit sentinel, and small helpers used across
// commands. It is imported by every internal/cmd/<command> package and by
// internal/cmd/root.
package cmdutil

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/dmccaffery/dotfiles/internal/envx"
	"github.com/dmccaffery/dotfiles/internal/execx"
	"github.com/dmccaffery/dotfiles/internal/logx"
	"github.com/dmccaffery/dotfiles/internal/ui"
)

// Deps are the injectable dependencies shared by every command. main supplies
// the real ones; root defaults Runner/Log/Prompt when nil; tests supply fakes.
type Deps struct {
	Runner execx.Runner
	Env    envx.Env
	Log    *logx.Logger
	Prompt ui.Prompter
}

// ErrSilent signals a non-zero exit after the command has already logged its own
// diagnostics. main suppresses printing it (cobra's other errors still print).
var ErrSilent = errors.New("dot: command failed")

// Streams wires a command's stdio to the runner for interactive children
// (editors, fzf, an authenticating gh, sudo prompts).
func Streams(cmd *cobra.Command) execx.Streams {
	return execx.Streams{In: cmd.InOrStdin(), Out: cmd.OutOrStdout(), Err: cmd.ErrOrStderr()}
}

// Arg returns args[i], or "" when out of range.
func Arg(args []string, i int) string {
	if i < len(args) {
		return args[i]
	}
	return ""
}

// DirExists reports whether path is an existing directory.
func DirExists(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.IsDir()
}

// FileExists reports whether path is an existing regular file.
func FileExists(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && !fi.IsDir()
}

// Atoi parses a trimmed integer, returning 0 on failure.
func Atoi(s string) int {
	n, _ := strconv.Atoi(strings.TrimSpace(s))
	return n
}

// NonEmptyLines splits s into trimmed, non-empty lines.
func NonEmptyLines(s string) []string {
	var out []string
	for line := range strings.SplitSeq(s, "\n") {
		if line = strings.TrimSpace(line); line != "" {
			out = append(out, line)
		}
	}
	return out
}

// CountNonEmptyLines mirrors `... | wc -l` over non-empty output.
func CountNonEmptyLines(s string) int { return len(NonEmptyLines(s)) }

// PickOne resolves a selection like fzf's --query/--select-1/--exit-0: it filters
// options by a case-insensitive substring query, returns the sole match without
// prompting, returns "" when nothing matches, and otherwise prompts the user.
func PickOne(p ui.Prompter, title, query string, options []string) (string, error) {
	matches := options
	if query != "" {
		matches = nil
		q := strings.ToLower(query)
		for _, o := range options {
			if strings.Contains(strings.ToLower(o), q) {
				matches = append(matches, o)
			}
		}
	}
	switch len(matches) {
	case 0:
		return "", nil
	case 1:
		return matches[0], nil
	default:
		return p.Select(title, matches)
	}
}
