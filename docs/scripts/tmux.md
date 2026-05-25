---
icon: simple/tmux
---

# Tmux scripts

Two scripts handle tmux session lifecycle: one for creating a fresh session per repo, one for
fuzzy-picking an existing one.

## `start-tmux-session` { #start-tmux-session }

```sh
sts                  # alias for start-tmux-session
sts <query>          # pre-fill the fzf query
sts .                # operate on $PWD instead of $REPOS_DIR
```

What it does:

1. Walks `$REPOS_DIR` (default `$HOME/Repos`) up to 4 levels deep looking for directories
   that contain a `.git/` entry.
2. Pipes the list into `fzf` for selection (`--select-1` auto-picks if there's only one).
3. Derives a session name from the repo's basename (with dots replaced by underscores).
4. If a session of that name doesn't exist, creates one with:
    - **Window 1 (` nvim`)** — nvim in the top pane (90%), shell in a small pane below (10%).
    - **Window 2 (`󰯉 claude (agent)`)** — runs `claude` (Claude Code) in the repo root.
    - **Window 3** — bare shell window.
5. Attaches to the session.

```sh title=".local/share/scripts/start-tmux-session (core)"
tmux -u new-session -d -s "${name}" -n ' nvim' -c "${selected}" -x - -y - "${EDITOR}" . \; \
    split-window -v -l '10%' -c "${selected}" \; \
    select-pane -t 1 \; \
    new-window -a -d -c "${selected}" -n '󰯉 claude (agent)' claude \; \
    new-window -a -d -c "${selected}"
```

## `attach-tmux-session` { #attach-tmux-session }

```sh
ats                  # alias for attach-tmux-session
ats <query>          # pre-fill the fzf query
```

Simpler: lists `tmux list-session -F '#S'`, fzf-picks, and either attaches (if running
outside tmux) or switches client (if inside).

The Snacks picker in nvim (++leader++ ++f++ ++s++) does the same thing without leaving the
editor — see [neovim/plugins](../neovim/plugins.md#snackslua).
