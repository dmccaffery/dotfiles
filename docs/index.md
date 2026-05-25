---
icon: lucide/home
---

# Deavon's Dotfiles

> Configuration for macOS terminal environments

This repository is a curated, stow-managed dotfiles setup for macOS. It provisions a coherent terminal
environment around [Ghostty](terminal/ghostty.md), [Zsh](terminal/shell.md), [tmux](terminal/tmux.md),
[oh-my-posh](terminal/oh-my-posh.md), and [LazyVim](neovim/lazyvim.md), all skinned with the
[cyberdream](theme/cyberdream.md) palette, plus a security-key-driven
[Git workflow](git/signing-security-keys.md) and a small fleet of [custom scripts](scripts/index.md).

## What's inside

<div class="grid cards" markdown>

- :material-monitor-shimmer: **[Terminal](terminal/ghostty.md)**

    Ghostty + Iosevka, zinit-managed Zsh, oh-my-posh prompt, tmux with auto-installed plugins,
    and a curated set of CLI tools.

- :simple-neovim: **[NeoVim](neovim/lazyvim.md)**

    LazyVim with 12 language extras, custom plugins, auto-commands, and the cyberdream colourscheme.

- :material-palette: **[Theme](theme/cyberdream.md)**

    One palette, applied everywhere: nvim, oh-my-posh, ghostty, tmux, btop, lazygit, yazi, k9s,
    opencode, vivid, lsd, and Claude Code.

- :simple-git: **[Git](git/config.md)**

    Sensible defaults, [git-town](git/git-town.md) aliases, GitHub OAuth + Forgejo auth, YubiKey SSH
    signing, and an include hook for a private overlay.

- :material-script-text: **[Scripts](scripts/index.md)**

    Fifteen small scripts covering tmux session management, resident security keys, Zscaler
    toggling, shell profiling, and more.

- :simple-apple: **[macOS + Claude Code](macos/system-defaults.md)**

    System defaults, a homebrew-managed launch agent for ssh-agent, and a sandboxed, vim-mode
    Claude Code harness with a cyberdream theme.

</div>

## Getting started

1. [Fork this repository](getting-started/install.md#fork-first) — it's tuned for one person; expect
   breakage if you track upstream.
2. [Back up your existing dotfiles](getting-started/backup.md).
3. [Run `./install.sh`](getting-started/install.md) to install Homebrew, apply macOS defaults, and stow the configs.
4. [Customize](getting-started/customize.md) `REPOS_DIR` and add a private Git overlay if you need one.

!!! warning "Fork-first"
This repository is primarily for personal use; breaking changes land without ceremony. Forking
insulates you from those changes and lets you keep your own private overlay separately.

## Build the docs locally

```sh
uv sync
uv run zensical serve   # http://localhost:8000
```

See [Contributing](tooling/contributing.md) for the full dev loop.
