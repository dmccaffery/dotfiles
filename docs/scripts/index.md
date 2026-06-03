---
icon: lucide/code
---

# Custom scripts

All scripts live in `.local/share/scripts/`. The directory is added to `PATH` by `.zshenv`
(only if it exists), so anything dropped in is immediately callable by name.

## Inventory

| Script                                                        | Purpose                                                                                                                                                                                                                    |
| ------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [`tmux-session start`](tmux.md#tmux-session-start)            | fzf-pick a repo and create/attach a tmux session — optionally per worktree.                                                                                                                                                |
| [`tmux-session attach`](tmux.md#tmux-session-attach)          | fzf-pick an existing tmux session and attach.                                                                                                                                                                              |
| [`tmux-session end`](tmux.md#tmux-session-end)                | fzf multi-select agent worktrees and remove them (with unpushed-commit warnings).                                                                                                                                          |
| [`worktree start`](../claude/hooks-skills.md#worktreecreate)  | Create a worktree at `~/.cache/agent/worktrees/<repo>-<suffix>`. Used by `tmux-session start` and Claude Code's `WorktreeCreate` hook.                                                                                     |
| [`worktree end`](../claude/hooks-skills.md#worktreeremove)    | Remove a worktree, its `agent/*` branch, and matching tmux session. Used by `tmux-session end` and Claude Code's `WorktreeRemove` hook.                                                                                    |
| [`agent-tmux-status`](tmux.md#agent-tmux-status)              | Flag the tmux window red (or set the terminal title) when a coding agent is waiting for input. Driven by Claude Code's `Stop`/`Notification`/`UserPromptSubmit`/`SessionEnd` hooks and opencode's status-indicator plugin. |
| [`ssh-sk gen`](security-keys.md#ssh-sk-gen)                   | Generate a resident ed25519-sk key on a YubiKey (requires a `<user>` argument).                                                                                                                                            |
| [`ssh-sk get`](security-keys.md#ssh-sk-get)                   | Load resident keys into ssh-agent and update allowed_signers.                                                                                                                                                              |
| [`ssh-sk get --github`](security-keys.md#ssh-sk-get-github)   | Resolve and print the matching GitHub signing key (git `defaultKeyCommand`).                                                                                                                                               |
| [`ssh-sk get --forgejo`](security-keys.md#ssh-sk-get-forgejo) | Resolve and print the matching Forgejo signing key (git `defaultKeyCommand`).                                                                                                                                              |
| [`gh-switch-user`](security-keys.md#gh-switch-user)           | `gh` wrapper that auto-switches to the account named in `git config github.account` before passing arguments through.                                                                                                      |
| [`git-github-auth`](security-keys.md#git-github-auth)         | Ensure `gh` is authenticated with the required scopes.                                                                                                                                                                     |
| [`git-resign`](security-keys.md#git-resign)                   | Re-sign every commit in a range with the current signing key.                                                                                                                                                              |
| [`ssh-askpass`](security-keys.md#ssh-askpass)                 | pinentry-mac wrapper used by the launch-managed ssh-agent.                                                                                                                                                                 |
| [`zs enable`](zscaler.md#zs-enable)                           | Load the Zscaler service + tunnel launch daemons.                                                                                                                                                                          |
| [`zs disable`](zscaler.md#zs-disable)                         | Unload the Zscaler service + tunnel.                                                                                                                                                                                       |
| [`zs certs`](zscaler.md#zs-certs)                             | Run a command with the Zscaler root CA injected as an extra trust anchor.                                                                                                                                                  |
| [`profile-shell`](misc.md#profile-shell)                      | Time Zsh startup with `zprof` enabled.                                                                                                                                                                                     |
| [`print-colors`](misc.md#print-colors)                        | Print a 24-bit truecolor gradient bar.                                                                                                                                                                                     |
| [`fzf-image-preview`](misc.md#fzf-image-preview)              | Preview handler for fzf — chafa for images, bat for text.                                                                                                                                                                  |
| [`reset-background-items`](misc.md#reset-background-items)    | `sfltool resetbtm` then reboot.                                                                                                                                                                                            |
| [`brewfile`](misc.md#brewfile)                                | `brew bundle <add\|remove> --global` then `brew bundle install --global`.                                                                                                                                                  |

## How scripts use color & logging

Every script uses the same minimal pattern based on `tput`:

```sh
default_color=$(tput sgr 0)
blue="$(tput setaf 4)"
info()    { printf "%s==> %s%s\n" "$blue"   "$1" "$default_color"; }
green="$(tput setaf 2)"
success() { printf "%s==> %s%s\n" "$green"  "$1" "$default_color"; }
yellow="$(tput setaf 3)"
warn()    { printf "%s==> %s%s\n" "$yellow" "$1" "$default_color"; }
red="$(tput setaf 1)"
error()   { printf "%s==> %s%s\n" "$red"    "$1" "$default_color"; }
```

This makes script output trivially scannable when chained together by the installer.
