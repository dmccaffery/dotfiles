---
icon: lucide/component
---

# Hooks & skills

This repository ships a Claude Code [settings.json](settings.md) and a [theme](theme.md), but
no `hooks/` or user-level `skills/` directories. The
[`.claude/plans/`](https://claude.com/claude-code) directory is present but is a runtime
artifact for plan mode — not configuration.

## Why no hooks?

Hooks are user-specific automation; what makes sense for one workflow doesn't generalize. The
relevant configuration is documented at
[docs.claude.com/en/docs/claude-code/hooks](https://docs.claude.com/en/docs/claude-code/hooks).

If you want to add hooks in your own fork, the conventional location is `~/.claude/hooks/`
(global) or `<repo>/.claude/hooks/` (project-scoped). Both are picked up by Claude Code
automatically.

## Skills

Same story for skills — none ship in this repo. The default Claude Code install provides
several built-in skills (`init`, `review`, `security-review`, etc.); they're discovered
automatically and don't need to be checked in.

Project-specific skills go in `<repo>/.claude/skills/`. See
[Claude Code skills docs](https://docs.claude.com/en/docs/claude-code/skills) for the
declaration format.

## Plans

`.claude/plans/` (excluded from stow via the conventions in [stow](../tooling/stow-and-make.md))
holds plan files written by Claude Code while in plan mode. The directory is intentionally
gitignored at the repo level via the matching `.gitignore` entry — plans are session
artifacts, not configuration.

## Future home

If you do add hooks or skills to your fork, they belong here in the docs alongside this
page. The intent of this page is to be the entry point for those additions when they exist.
