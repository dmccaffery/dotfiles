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

## `brewfile` { #brewfile }

```sh
brewfile add <package> [brew bundle flags...]
brewfile remove <package>
```

A [`dot`](../tooling/dot.md) applet (same interface as the former shell script) that wraps the
two-step day-2 Brewfile flow into one command:

1. `brew bundle <add|remove> "$@" --global` — edit `$HOMEBREW_BUNDLE_FILE_GLOBAL`.
2. `brew bundle install --global --zap --upgrade` — install, upgrade outdated formulae, and (because
   `HOMEBREW_BUNDLE_FORCE_INSTALL_CLEANUP=1` is exported) uninstall/zap anything no longer in the Brewfile, so the
   Brewfile and the installed state stay in lockstep.

Flags after the action are passed through to `brew bundle`, so `brewfile add --cask ghostty` and
`brewfile add --tap user/tap` work. See [Brew bundle](../terminal/brew-bundle.md) for the underlying
commands and the environment variables that shape them.

On `add`, the wrapper trust-checks anything that comes from a non-official tap (a `user/tap/...` reference)
before the install runs. The `--cask`/`--tap`/`--formula` flag decides which `trust.json` list to look in
(`trustedcasks`, `trustedtaps`, `trustedformulae`); bare names resolve to `homebrew/core`/`homebrew/cask` and
are trusted by default, so they are skipped. If a referenced entry is missing from
[`trust.json`](../../.config/homebrew/trust.json), you are prompted to `brew trust --<type> <name>` it — keeping
`brew bundle install` from stalling on an untrusted tap mid-run. Decline to leave it untrusted; with no tty the
prompt is skipped with a warning.
