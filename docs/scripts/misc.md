---
icon: lucide/sparkles
---

# Miscellaneous scripts

Small utilities that don't fit a larger category.

## `profile-shell` { #profile-shell }

```sh
profile-shell
```

One-liner that times Zsh startup with `zprof` enabled:

```sh title=".local/share/scripts/profile-shell"
ZSHPROFILE=1 time -p zsh -i -c exit
```

`.zshrc` checks `$ZSHPROFILE` at the top and bottom — when set, it loads `zsh/zprof` early
and prints a timing table on exit. The wrapper also gives you wall/user/system numbers via
`time -p`.

## `print-colors` { #print-colors }

```sh
print-colors
```

Prints a horizontal 24-bit RGB gradient bar to test that your terminal supports true color.
The bar smoothly transitions through the full red/green/blue spectrum.

Useful as a smoke test after switching terminal, font, or `TERM` settings.

## `fzf-image-preview` { #fzf-image-preview }

```sh
fzf-image-preview <file>
```

Preview handler intended for use with fzf's `--preview`. Behaviour:

- **Directory** → `ls -la --color`.
- **Binary image** → renders with `chafa --passthrough=<auto|tmux> --size=<cols>`.
- **Other binary** → prints `<file> is a binary file`.
- **Text** → `bat --style=numbers --color=always | head -100`.

The passthrough mode is automatically `tmux` when running inside a tmux session, otherwise
`auto` (chafa picks kitty graphics, sixel, or character art based on terminal capability).

Use as:

```sh
fzf --preview 'fzf-image-preview {}'
```

## `reset-background-items` { #reset-background-items }

```sh
reset-background-items
```

Wraps `sudo sfltool resetbtm` (resets all macOS background-task management state) with two
"press any key" prompts: one before resetting, one before the required reboot.

Use this when macOS gets confused about which background items are authorized — usually
manifests as repeated "Background Items Added" notifications, or login items that can't be
disabled in System Settings.

!!! warning "Reboot required"
`sfltool resetbtm` only takes effect after a reboot. The script prompts and runs
`sudo reboot` for you.
