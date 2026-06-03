---
icon: simple/tmux
---

# Tmux scripts

Three scripts handle tmux session lifecycle: one for creating a fresh session per repo (or per
worktree), one for fuzzy-picking an existing session, and one for tearing down agent worktree
sessions and the worktrees behind them. The worktree create/remove work itself is delegated
to [`start-worktree`](index.md) / [`end-worktree`](index.md), which are also wired in as
[Claude Code's `WorktreeCreate`/`WorktreeRemove` hooks](../claude/hooks-skills.md#hooks) so
the naming convention stays consistent regardless of who created the worktree.

## `start-tmux-session` { #start-tmux-session }

```sh
sts                          # alias for start-tmux-session
sts <query>                  # pre-fill the fzf query
sts .                        # operate on $PWD instead of $REPOS_DIR
sts <query> <worktree-name>  # create/attach a session inside a per-worktree checkout
```

What it does:

1. Walks `$REPOS_DIR` (default `$HOME/Repos`) up to 4 levels deep looking for directories
   that contain a `.git/` entry.
2. Pipes the list into `fzf` for selection (`--select-1` auto-picks if there's only one).
3. Sanitises the repo's basename to derive the bare-repo session name (see the
   [shared sanitizer](#sanitizer) below).
4. **If a second `<worktree-name>` argument is supplied**, hands off to
   [`start-worktree`](index.md), which creates `~/.cache/agent/worktrees/<repo>-<worktree>`
   on branch `agent/<repo>-<worktree>` via `git worktree add` (reusing the path or branch if
   either already exists), then prints the worktree path back. The tmux session is named
   after the worktree directory's basename (`<repo>-<worktree>`) so the
   [Snacks sessions picker](../neovim/plugins.md#snackslua) can nest it under the bare-repo
   parent by name prefix.
5. If a session of that name doesn't exist, creates one. The `nvim` window is created first;
   every other window is created detached and inserted with `-a` (immediately **after** the
   `nvim` window), so the windows end up in reverse creation order. The creation order is
   `opencode` → `claude` → `zsh`, which yields this layout:
    - **Window 1 (` nvim`)** — nvim in the top pane (90%), shell in a small pane below (10%).
    - **Window 2 (` zsh`)** — a plain login shell in the repo root.
    - **Window 3 (`󰯉  claude`)** — runs `claude` (Claude Code) in the repo root when
      it is available on `PATH`.
    - **Window 4 (`󰚩  opencode`)** — runs `opencode` in the repo root when it is
      available on `PATH`.

    The two agent windows are skipped entirely when their CLI isn't on `PATH`, so the trailing
    indices shift down accordingly — but `zsh` always lands at window 2, immediately after the
    editor.

6. Sets the terminal window/tab title to the session name via an `OSC 0` escape
   (`printf '\033]0;%s\007'`), so the tab reads e.g. `dotfiles` instead of the launching
   command `sts dotfiles`. tmux leaves this alone because `set-titles` is off.
7. Attaches to the session.

```sh title=".local/share/scripts/start-tmux-session (core)"
claude_bin=$(command -v claude 2> /dev/null || true)
opencode_bin=$(command -v opencode 2> /dev/null || true)

editor_pane=$(tmux -u new-session -d -P -F '#{pane_id}' -s "${name}" -n ' nvim' -c "${selected}" \
    -x - -y - "${EDITOR}" .)
editor_window=$(tmux display-message -p -t "${editor_pane}" '#{window_id}')
tmux split-window -t "${editor_pane}" -v -l '10%' -c "${selected}"
tmux select-pane -t "${editor_pane}"

[ -n "${opencode_bin}" ] && tmux new-window -a -d -t "${editor_window}" -c "${selected}" \
    -n '󰚩  opencode' "${opencode_bin}"
[ -n "${claude_bin}" ] && tmux new-window -a -d -t "${editor_window}" -c "${selected}" \
    -n '󰯉  claude' "${claude_bin}"

tmux new-window -a -d -t "${editor_window}" -n ' zsh' -c "${selected}"
```

### The shared sanitizer { #sanitizer }

Both [`start-tmux-session`](#start-tmux-session) and [`start-worktree`](index.md) run names
through the same `sanitize` helper, so `fix/stow symlinks` becomes `fix-stow-symlinks`. It
collapses any character outside `A-Za-z0-9_-` to `-`, with special handling for `.`: tmux
3.5+ rejects `.` in session names (it's the session/window/pane separator), so dots are
encoded rather than dropped — a leading `.` becomes `dot-`, a trailing `.` becomes `-dot`,
and an interior `.` becomes `-dot-`. So `next.js` becomes `next-dot-js` and `.config`
becomes `dot-config`, keeping each name unique and tmux-safe.

## `attach-tmux-session` { #attach-tmux-session }

```sh
ats                  # alias for attach-tmux-session
ats <query>          # pre-fill the fzf query
```

Simpler: lists `tmux list-session -F '#S'`, fzf-picks, and either attaches (if running
outside tmux) or switches client (if inside).

The Snacks picker in nvim (++leader++ ++f++ ++s++) does the same thing without leaving the
editor — see [neovim/plugins](../neovim/plugins.md#snackslua).

## `end-tmux-session` { #end-tmux-session }

```sh
ets                          # alias for end-tmux-session — fzf multi-select over agent worktrees
ets <worktree-name>...       # remove specific worktrees by name (or absolute path)
ets -f <worktree-name>...    # skip the confirmation prompt when worktrees are dirty
```

What it does:

1. Builds a selection list from positional args, or interactively via `fzf -m` over
   `~/.cache/agent/worktrees/*` (tab to mark, enter to confirm).
2. **Inspect pass** — for each selection, prints status and flags concerns:
    - `uncommitted` — count of working-tree changes (`git status --porcelain`).
    - `unpushed` — count of commits reachable from `HEAD` but absent from any remote ref
      (`git rev-list --count HEAD --not --remotes`). This catches both "no upstream set" and
      "upstream set but ahead".
3. If any selection had warnings and `--force` wasn't passed, prompts once before continuing.
4. **Destroy pass** — hands each selected path to [`end-worktree`](index.md), which:
    - Resolves the parent repo via `git rev-parse --git-common-dir`.
    - Kills the matching tmux session (`<repo>-<worktree>` — the worktree dir basename) if
      present.
    - `git worktree remove --force <path>` from the parent repo.
    - `git branch -D <branch>` if the branch is in the `agent/*` namespace (matches the
      convention used by [`start-worktree`](../claude/hooks-skills.md#worktreecreate)).

Remote branches are never touched — push before removing if you want to keep the work. The
matching PR (if any) keeps working off the remote branch even after the local one is gone.

## `agent-tmux-status` { #agent-tmux-status }

```sh
agent-tmux-status waiting     # turn finished — your turn (calm)
agent-tmux-status attention   # the agent needs you now — permission / notification (urgent)
agent-tmux-status clear       # lower the indicator (also the default with no/unknown arg)
```

A no-op-safe leaf script shared by two coding agents so the indicator tracks whether either is
waiting on you:

- **Claude Code** drives it through five
  [hooks](../claude/hooks-skills.md#claude-is-waiting-indicator) — `Stop` calls it with
  `waiting`, `Notification` with `attention`, and `PostToolUse`/`UserPromptSubmit`/`SessionEnd`
  with `clear` (`PostToolUse` clears the red after you approve a permission and the tool runs).
- **opencode** drives it through the
  [`agent-tmux-status` plugin](../opencode/plugins.md#status-indicator) — `session.idle` →
  `waiting`, `permission.updated` → `attention`, a new user message → `clear`.

It branches on `$TMUX`:

- **Inside tmux** — stores the **state token** (`waiting` or `attention`) in a per-window user
  option on the pane the agent runs in (`tmux set-window-option -t "$TMUX_PANE" @agent_status …`,
  unset on clear). tmux swallows `OSC` title sequences from inside a session (`set-titles` is
  off), so the option is the reliable channel. [`theme.conf`](../terminal/tmux.md#agent-status)
  maps the token to a colour and glyph in the `window-status` format — calm **peach `●`** for
  `waiting`, bold **red `󰂚`** for `attention` — and only that window changes, since user options
  resolve per-window.
- **Outside tmux** — falls back to an `OSC 0` terminal title written to `/dev/tty` (the same
  trick [`start-tmux-session`](#start-tmux-session) uses, since the agent may capture the
  caller's stdout), prefixing the cwd basename with `●` (waiting) or `󰂚` (attention) and
  dropping it on clear.

Every `tmux`/`printf` call is guarded with `|| true` and the script never exits non-zero, so a
missing tmux server or detached tty can't fail an agent turn. The fallback glyphs live in the
`waiting_glyph` / `attention_glyph` lines at the top of the script; the in-tmux colours and
glyphs live in [`theme.conf`](../terminal/tmux.md#agent-status).
