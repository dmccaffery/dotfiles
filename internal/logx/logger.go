package logx

import (
	"io"
	"log/slog"
)

// Logger wraps *slog.Logger. Its methods (Info, Warn, Error, With, …) come from
// the embedded slog.Logger, so call sites pass typed slog.Attr values for
// structured logging.
type Logger struct {
	*slog.Logger
}

// NewLogger builds a Logger from a slog.Handler.
func NewLogger(h slog.Handler) *Logger {
	return &Logger{slog.New(h)}
}

// For builds a Logger that writes to w, choosing styled text (tty) or JSON.
func For(w io.Writer, tty bool) *Logger {
	return NewLogger(newHandler(w, tty))
}

// With returns a Logger with the given attributes bound to every record.
func (l *Logger) With(args ...any) *Logger {
	return &Logger{l.Logger.With(args...)}
}
