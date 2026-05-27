---
icon: lucide/component
---

# Hooks & skills

This repository ships a Claude Code [settings.json](settings.md), a [theme](theme.md), and a
pair of worktree-lifecycle hooks (`.claude/hooks/worktree-create.sh`,
`.claude/hooks/worktree-remove.sh`) wired up via `settings.json`. No user-level `skills/`
directory is checked in. The [`.claude/plans/`](https://claude.com/claude-code) directory is
present but is a runtime artifact for plan mode — not configuration.

## Hooks

Claude Code looks for hooks in `~/.claude/hooks/` (global) and `<repo>/.claude/hooks/`
(project-scoped). Both are picked up automatically once a matching event is declared in
`settings.json`. Upstream docs:
[docs.claude.com/en/docs/claude-code/hooks](https://docs.claude.com/en/docs/claude-code/hooks).

### WorktreeCreate

Registered in [`settings.json`](settings.md#hooks):

```json
"hooks": {
    "WorktreeCreate": [
        { "hooks": [{ "type": "command", "command": "~/.claude/hooks/worktree-create.sh" }] }
    ]
}
```

`.claude/hooks/worktree-create.sh` runs whenever Claude Code creates a worktree. It:

1. Picks a base name. Inside `tmux`, the active session name is used; otherwise the repo's
   directory name (resolved from `$CLAUDE_PROJECT_DIR`, falling back to
   `git rev-parse --show-toplevel`).
2. Picks a suffix. If Claude passes a `name` on stdin (JSON via the hook protocol), that
   wins; otherwise the script falls back to a UTC `YYYYMMDD-HHMMSS` timestamp.
3. Sanitises the combined `<name>-<suffix>` (anything outside `[A-Za-z0-9._-]` becomes `-`)
   and creates the worktree at `~/.cache/agent/worktrees/<name>` on a fresh branch
   `agent/<name>`.
4. Prints the worktree path to stdout (Claude Code consumes it) and logs progress to the
   controlling tty using `tput` colours.

The worktree root (`~/.cache/agent/worktrees/`) is XDG cache by design — worktrees are
throwaway work areas, not configuration. The path is hard-coded in the script because it
has to match the literal path in the [`includeIf "gitdir:…"` block](../git/config.md#includes)
that loads [`agent.gitconfig`](../git/config.md#includes), and git's `includeIf` can't expand
environment variables. Bonus: `~/.cache` is already in the sandbox
[`filesystem.allowRead`/`allowWrite`](settings.md#sandbox) lists, so no separate sandbox
extension is needed.

The tmux-session lookup runs `tmux display-message` against the active tmux socket. The
sandbox refuses `connect()` on any AF_UNIX path (see
[settings → sandbox](settings.md#sandbox)), so when the hook is invoked from inside Claude's
sandbox the lookup fails silently (the `2> /dev/null || true` swallow) and step 1 falls back
to the repo basename. The `agent/` branch prefix is deliberate: it's the signal
[`worktree-remove.sh`](#worktreeremove) uses to decide a branch is safe to delete.

### WorktreeRemove

Registered in [`settings.json`](settings.md#hooks):

```json
"hooks": {
    "WorktreeRemove": [
        { "hooks": [{ "type": "command", "command": "~/.claude/hooks/worktree-remove.sh" }] }
    ]
}
```

`.claude/hooks/worktree-remove.sh` runs when Claude Code finishes with a worktree. It:

1. Reads the hook payload from stdin and extracts `.worktree_path` with `jq`.
2. Exits cleanly if the path is already gone (`! -d`), so re-runs and manual cleanups don't
   error.
3. Captures the current branch with `git rev-parse --abbrev-ref HEAD` inside the worktree.
4. Warns to the TTY if the worktree has uncommitted changes (`git status --porcelain`) or
   unpushed commits (`git rev-list --count HEAD --not --remotes` — counts commits reachable
   from HEAD but absent from every remote ref). The warning doesn't block removal — Claude
   initiated this and the hook can't prompt — but the loud notice gives you a chance to
   recover lost commits from the reflog within ~30 days.
5. Removes the worktree (`git worktree remove --force`), swallowing failures so a half-broken
   worktree doesn't block the hook.
6. Deletes the branch (`git branch -D`) only if it starts with `agent/` — the prefix
   [`worktree-create.sh`](#worktreecreate) uses. Branches with any other prefix are left in
   place.

The `agent/`-prefix guard is intentional: branch deletion is destructive, so the hook only
touches branches it can identify as throwaway agent state.

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

## Adding more

Additional hooks or skills belong on this page. Add a fresh subsection alongside
[WorktreeCreate](#worktreecreate) and link to the script from
[`settings.json`](settings.md#hooks).
