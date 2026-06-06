package logx

import "log/slog"

// Worktree describes a git worktree for structured logging. It implements
// slog.LogValuer so it renders as a typed, grouped attribute set rather than a
// stringly-typed blob.
type Worktree struct {
	Name   string
	Path   string
	Branch string
	Repo   string
}

// LogValue implements slog.LogValuer.
func (w Worktree) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("name", w.Name),
		slog.String("path", w.Path),
		slog.String("branch", w.Branch),
		slog.String("repo", w.Repo),
	)
}

// WorktreeAttr returns a typed "worktree" attribute.
func WorktreeAttr(w Worktree) slog.Attr { return slog.Any("worktree", w) }

// Package describes a Homebrew package reference for structured logging.
type Package struct {
	Type string // formula | cask | tap
	Name string
}

// LogValue implements slog.LogValuer.
func (p Package) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("type", p.Type),
		slog.String("name", p.Name),
	)
}

// PackageAttr returns a typed "package" attribute.
func PackageAttr(p Package) slog.Attr { return slog.Any("package", p) }
