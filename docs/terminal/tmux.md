---
icon: simple/tmux
---

# tmux

`.config/tmux/tmux.conf` is a thin loader; the real config is split across modular files under
`.config/tmux/conf/`.

```text title=".config/tmux/tmux.conf"
source ~/.config/tmux/conf/plugins.conf
source ~/.config/tmux/conf/keymaps.conf
source ~/.config/tmux/conf/options.conf

# install tpm if not already installed
if 'test ! -d #{TMUX_PLUGIN_MANAGER_PATH}/tpm' {
    run 'git clone https://github.com/tmux-plugins/tpm #{TMUX_PLUGIN_MANAGER_PATH}/tpm'
    run '#{TMUX_PLUGIN_MANAGER_PATH}/tpm/bin/install_plugins'
}

run -b '#{TMUX_PLUGIN_MANAGER_PATH}/tpm/tpm'
source ~/.config/tmux/conf/theme.conf
```

TPM (the Tmux Plugin Manager) is auto-installed on first launch. No manual bootstrap.

## Keymaps

```text title=".config/tmux/conf/keymaps.conf"
unbind C-a
set -g prefix C-a
bind C-a send-prefix
```

| Binding                                  | Action                                                                |
| ---------------------------------------- | --------------------------------------------------------------------- |
| ++ctrl+a++                               | Prefix (sent through with ++ctrl+a++ ++ctrl+a++).                     |
| ++ctrl+a++ ++bar++                       | Split window horizontally, preserve cwd.                              |
| ++ctrl+a++ ++minus++                     | Split window vertically, preserve cwd.                                |
| ++ctrl+a++ ++ctrl+p++                    | Previous window.                                                      |
| ++ctrl+a++ ++ctrl+n++                    | Next window.                                                          |
| ++ctrl+a++ ++shift+c++                   | New window running Claude Code in the current pane path if available. |
| ++ctrl+a++ ++shift+o++                   | New window running OpenCode in the current pane path if available.    |
| ++ctrl+a++ ++shift+x++                   | New window running Codex in the current pane path if available.       |
| ++ctrl+a++ ++h++ / ++j++ / ++k++ / ++l++ | Resize pane left / down / up / right (repeatable).                    |
| ++ctrl+a++ ++m++                         | Toggle pane zoom.                                                     |
| ++ctrl+a++ ++backspace++                 | Kill current session.                                                 |
| ++ctrl+a++ ++f++                         | Kill current session (alias).                                         |

Copy mode uses vi keys (`v` to start selection, `y` to copy). Mouse drag-select does _not_ exit
copy mode, so you can refine a selection after dragging.

The agent bindings show a tmux message instead of opening a broken window when their CLI is
missing from `PATH`.

## Options

```text title=".config/tmux/conf/options.conf"
set -g mouse on
set -g status-position top
set -g base-index 1
set -g renumber-windows 1
set -g default-shell $SHELL
set -g escape-time 0
set -g history-limit 50000
set -g extended-keys on
set -g extended-keys-format xterm
set -g set-clipboard on
```

Highlights:

- **Status bar at the top** to match the prompt direction.
- **`base-index 1`** + **`pane-base-index 1`** + **`renumber-windows`** — windows and panes
  number from 1 and stay contiguous.
- **`escape-time 0`** — kill the default 500ms delay that breaks Vim mode.
- **`extended-keys`** — full xterm key reporting so modified keys work.
- **`set-clipboard on`** — let programs set the system clipboard via OSC 52, so a yank inside
  tmux (including over SSH) reaches the **local** machine's clipboard through Ghostty.
- **Shift+Enter → ++ctrl+j++** — root-level binding lets Claude Code (which uses Shift+Enter
  for newline-without-submit) work inside tmux despite tmux not supporting the kitty keyboard
  protocol.

### Claude Code tmux environment

`options.conf` exports two `CLAUDE_CODE_*` variables via `set-environment -g` so they only
apply to processes spawned inside tmux, not to plain terminal shells:

| Variable                     | Value  | Effect                                                                    |
| ---------------------------- | ------ | ------------------------------------------------------------------------- |
| `CLAUDE_CODE_TMUX_TRUECOLOR` | `true` | Tells Claude Code to emit 24-bit colour even when `$TERM` advertises 256. |
| `CLAUDE_CODE_NO_FLICKER`     | `true` | Suppresses the redraw flicker Claude Code's TUI shows under tmux.         |

A third toggle, `CLAUDE_CODE_DISABLE_MOUSE=true`, is kept commented out — turning it on hands
mouse events to tmux instead of Claude Code, but the trade-off (no in-app scroll / no
selection inside the TUI) didn't carry its weight. Uncomment it in `options.conf` if you'd
rather let tmux's `mouse on` win drag-select.

## Plugins

```text title=".config/tmux/conf/plugins.conf"
set -g @plugin 'tmux-plugins/tmux-sensible'
set -g @plugin 'wfxr/tmux-fzf-url'
set -g @plugin 'christoomey/vim-tmux-navigator'
set -g @plugin 'joshmedeski/tmux-nerd-font-window-name'
set -g @plugin 'tmux-plugins/tpm'
```

| Plugin                       | What it does                                                                                         |
| ---------------------------- | ---------------------------------------------------------------------------------------------------- |
| `tmux-sensible`              | Sane defaults that don't conflict with custom options.                                               |
| `tmux-fzf-url`               | ++prefix++ ++u++ to pick a URL from the visible buffer with fzf.                                     |
| `vim-tmux-navigator`         | ++ctrl+h++ / ++ctrl+j++ / ++ctrl+k++ / ++ctrl+l++ move between Vim splits and tmux panes seamlessly. |
| `tmux-nerd-font-window-name` | Automatic window names use Nerd Font icons for known processes.                                      |
| `tpm`                        | The plugin manager itself.                                                                           |

## Theme

`conf/theme.conf` uses [Catppuccin for tmux](https://github.com/catppuccin/tmux) with a
cyberdream flavour. The flavour file isn't shipped by upstream Catppuccin, so the config
self-bootstraps by `curl`-ing it from the
[cyberdream.nvim extras directory](https://github.com/scottmckendry/cyberdream.nvim/tree/main/extras/tmux)
on first run.

The status bar shows: session name plus current pane command on the left, window list in the
centre, and the short hostname on the right. The session name doubles as the prefix-key
indicator — green while idle, red and bold while the ++ctrl+a++ prefix is held. When tmux has
`SSH_CONNECTION`, the hostname is followed by `| 󰣀` with the SSH glyph in red; local sessions
omit both the separator and glyph.

### Agent-is-waiting indicator { #agent-status }

`theme.conf` makes the `window-status` format react to a per-window `@agent_status` user
option that holds a **state token** so a window can flag what a coding agent (Claude Code,
opencode, or Codex) needs. Both `window-status-format` and `window-status-current-format` gain a
leading conditional segment that maps the token to a style:

```text title=".config/tmux/conf/theme.conf (style segment)"
#{?@agent_status,#{?#{==:#{@agent_status},attention},#[bg=#{@thm_red}#,fg=#{@thm_bg}#,bold],#[bg=#{@thm_peach}#,fg=#{@thm_bg}]},}
```

| Token       | Set when                          | Look                           |
| ----------- | --------------------------------- | ------------------------------ |
| _(empty)_   | cleared                           | base `window-status-style`     |
| `waiting`   | turn finished (your move)         | calm **peach** background, `●` |
| `attention` | blocked on permission / attention | bold **red** background, `󰂚`   |

A matching glyph segment appends `●` or `󰂚` after the window name via the same
`#{==:…,attention}` test. The literal commas inside `#[…]` are escaped as `#,` because `,`
separates the conditional's branches; the commas inside nested `#{…}` are protected by the
braces and left bare. The active and last-window styles use `@thm_blue` (not the indicator's
`@thm_peach`/`@thm_red`), so a `waiting` or `attention` window pops by colour even when it is the
focused one — no overlap between the "this is the current window" cue and the "this window needs
you" cue.

The option is toggled by the [`agent-tmux-status`](../scripts/tmux.md#agent-tmux-status)
script, wired into Claude Code's
[`Stop`/`Notification`/`UserPromptSubmit`/`SessionEnd` hooks](../claude/hooks-skills.md#claude-is-waiting-indicator),
opencode's [status-indicator plugin](../opencode/plugins.md#status-indicator), and Codex's
[`Stop`/`PermissionRequest`/`PostToolUse`/`UserPromptSubmit` hooks](../codex/hooks.md#the-tmux-indicator).

## Companion scripts

See [scripts/tmux](../scripts/tmux.md) for `tmux-session start` (fzf-driven repo browser that
creates a named session per repo) and `tmux-session attach`.
