// Package logx is the dot CLI's diagnostics layer. It adapts charmbracelet/log
// to slog: styled, leveled output on a terminal and structured JSON when piped,
// so call sites use fully-typed slog attributes. Program results (e.g. a
// worktree path) are written to stdout directly and never logged.
package logx

import (
	"io"
	"log/slog"
	"os"

	"github.com/charmbracelet/log"
)

// newHandler builds a charmbracelet/log handler — styled text on a TTY, JSON
// otherwise — usable as an slog.Handler.
func newHandler(w io.Writer, tty bool) slog.Handler {
	l := log.NewWithOptions(w, log.Options{
		Level:           log.InfoLevel,
		ReportTimestamp: !tty, // timestamps are noise on a TTY, useful in JSON logs
	})
	if !tty {
		l.SetFormatter(log.JSONFormatter)
	}
	return l
}

// IsTerminal reports whether w is a character device (a TTY).
func IsTerminal(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}
