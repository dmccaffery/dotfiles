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

| Binding                                  | Action                                             |
| ---------------------------------------- | -------------------------------------------------- |
| ++ctrl+a++                               | Prefix (sent through with ++ctrl+a++ ++ctrl+a++).  |
| ++ctrl+a++ ++bar++                       | Split window horizontally, preserve cwd.           |
| ++ctrl+a++ ++minus++                     | Split window vertically, preserve cwd.             |
| ++ctrl+a++ ++ctrl+p++                    | Previous window.                                   |
| ++ctrl+a++ ++ctrl+n++                    | Next window.                                       |
| ++ctrl+a++ ++h++ / ++j++ / ++k++ / ++l++ | Resize pane left / down / up / right (repeatable). |
| ++ctrl+a++ ++m++                         | Toggle pane zoom.                                  |
| ++ctrl+a++ ++backspace++                 | Kill current session.                              |
| ++ctrl+a++ ++f++                         | Kill current session (alias).                      |

Copy mode uses vi keys (`v` to start selection, `y` to copy). Mouse drag-select does _not_ exit
copy mode, so you can refine a selection after dragging.

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
```

Highlights:

- **Status bar at the top** to match the prompt direction.
- **`base-index 1`** + **`pane-base-index 1`** + **`renumber-windows`** — windows and panes
  number from 1 and stay contiguous.
- **`escape-time 0`** — kill the default 500ms delay that breaks Vim mode.
- **`extended-keys`** — full xterm key reporting so modified keys work.
- **Shift+Enter → ++ctrl+j++** — root-level binding lets Claude Code (which uses Shift+Enter
  for newline-without-submit) work inside tmux despite tmux not supporting the kitty keyboard
  protocol.

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

The status bar shows: current pane command + cwd (left) and session name (right, recoloured
red while the prefix is held).

## Companion scripts

See [scripts/tmux](../scripts/tmux.md) for `start-tmux-session` (fzf-driven repo browser that
creates a named session per repo) and `attach-tmux-session`.
