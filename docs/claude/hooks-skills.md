---
icon: lucide/component
---

# Hooks & skills

This repository ships a Claude Code [settings.json](settings.md), a [theme](theme.md), and
two worktree-lifecycle hooks. Both hooks point at the same scripts the manual
[`start-tmux-session`](../scripts/tmux.md#start-tmux-session) /
[`end-tmux-session`](../scripts/tmux.md#end-tmux-session) wrappers call, so the naming
convention stays in lockstep regardless of who triggered the worktree. No user-level
`skills/` directory is checked in. The [`.claude/plans/`](https://claude.com/claude-code)
directory is present but is a runtime artifact for plan mode — not configuration.

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
        { "hooks": [{ "type": "command", "command": "~/.local/share/scripts/start-worktree" }] }
    ]
}
```

[`start-worktree`](../scripts/index.md) runs whenever Claude Code creates a worktree, and is
also invoked directly by [`start-tmux-session`](../scripts/tmux.md#start-tmux-session) when a
worktree-name argument is passed. It:

1. Picks the repo path from `$1`, falling back to `$CLAUDE_PROJECT_DIR`, then
   `git rev-parse --show-toplevel`. Its basename becomes the worktree's base name.
2. Picks a suffix from `$2`, falling back to a `name` field in JSON on stdin (the hook
   protocol), then a UTC `YYYYMMDD-HHMMSS` timestamp.
3. Sanitises both the base name and the suffix (anything outside `[A-Za-z0-9._-]` becomes
   `-`), joins them as `<base>-<suffix>`, and creates the worktree at
   `~/.cache/agent/worktrees/<name>` on a branch `agent/<name>`. If the branch already
   exists, it's reused instead of recreated.
4. Prints the worktree path to stdout (Claude Code consumes it) and logs progress to the
   controlling tty using `tput` colours.

The worktree root (`~/.cache/agent/worktrees/`) is XDG cache by design — worktrees are
throwaway work areas, not configuration. The path is hard-coded in the script because it
has to match the literal path in the [`includeIf "gitdir:…"` block](../git/config.md#includes)
that loads [`agent.gitconfig`](../git/config.md#includes), and git's `includeIf` can't expand
environment variables. Bonus: `~/.cache` is already in the sandbox
[`filesystem.allowRead`/`allowWrite`](settings.md#sandbox) lists, so no separate sandbox
extension is needed.

The `agent/` branch prefix is deliberate: it's the signal
[`end-worktree`](#worktreeremove) uses to decide a branch is safe to delete.

### WorktreeRemove

Registered in [`settings.json`](settings.md#hooks):

```json
"hooks": {
    "WorktreeRemove": [
        { "hooks": [{ "type": "command", "command": "~/.local/share/scripts/end-worktree" }] }
    ]
}
```

[`end-worktree`](../scripts/index.md) runs when Claude Code finishes with a worktree, and is
also invoked per selection by [`end-tmux-session`](../scripts/tmux.md#end-tmux-session). It:

1. Reads the worktree path from `$1`, falling back to a `worktree_path` field in JSON on
   stdin (the hook protocol).
2. Exits cleanly if the path is already gone (`! -d`), so re-runs and manual cleanups don't
   error.
3. Captures the current branch and the parent repo (via `git rev-parse --git-common-dir`)
   inside the worktree before anything is destroyed.
4. Warns to the TTY if the worktree has uncommitted changes (`git status --porcelain`) or
   unpushed commits (`git rev-list --count HEAD --not --remotes` — counts commits reachable
   from HEAD but absent from every remote ref). The warning doesn't block removal — the loud
   notice gives you a chance to recover lost commits from the reflog within ~30 days.
5. Kills the tmux session whose name matches the worktree directory's basename, if one
   exists. With the unified naming convention that's the session created for this worktree.
6. Removes the worktree (`git worktree remove --force`), swallowing failures so a
   half-broken worktree doesn't block the hook.
7. Deletes the branch (`git branch -D`) only if it starts with `agent/` — the prefix
   [`start-worktree`](#worktreecreate) uses. Branches with any other prefix are left in
   place.

The `agent/`-prefix guard is intentional: branch deletion is destructive, so the hook only
touches branches it can identify as throwaway agent state.

### Claude-is-waiting indicator { #claude-is-waiting-indicator }

Four hooks drive a "Claude is waiting for input" indicator, all pointing at the same leaf
script [`claude-tmux-status`](../scripts/tmux.md#claude-tmux-status) with a state argument:

```json
"hooks": {
    "Stop": [
        { "hooks": [{ "type": "command", "command": "~/.local/share/scripts/claude-tmux-status waiting" }] }
    ],
    "Notification": [
        { "hooks": [{ "type": "command", "command": "~/.local/share/scripts/claude-tmux-status waiting" }] }
    ],
    "UserPromptSubmit": [
        { "hooks": [{ "type": "command", "command": "~/.local/share/scripts/claude-tmux-status clear" }] }
    ],
    "SessionEnd": [
        { "hooks": [{ "type": "command", "command": "~/.local/share/scripts/claude-tmux-status clear" }] }
    ]
}
```

| Event              | Arg       | When it fires                                |
| ------------------ | --------- | -------------------------------------------- |
| `Stop`             | `waiting` | Claude finished a turn and is awaiting input |
| `Notification`     | `waiting` | Claude needs permission or attention         |
| `UserPromptSubmit` | `clear`   | You submitted a reply — Claude is busy again |
| `SessionEnd`       | `clear`   | Session ended — don't leave the flag stuck   |

The script branches on `$TMUX`: inside tmux it sets a per-window `@claude_status` option that
[`theme.conf`](../terminal/tmux.md#claude-status) renders as a red status-bar entry; outside
tmux it falls back to an `OSC 0` terminal title. See
[scripts/tmux → `claude-tmux-status`](../scripts/tmux.md#claude-tmux-status) for the details.

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
