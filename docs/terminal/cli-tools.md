---
icon: lucide/wrench
---

# CLI tools

The Brewfile pulls in a curated set of modern terminal tools. Many ship with cyberdream
configuration in this repo.

## Replacements & power tools

| Tool        | Replaces            | Config                                                        |
| ----------- | ------------------- | ------------------------------------------------------------- |
| `lsd`       | `ls`                | `.config/lsd/config.yaml` — aliased as `ls` in `.zshrc`.      |
| `bat`       | `cat`               | (default config)                                              |
| `fd`        | `find`              | (default config)                                              |
| `ripgrep`   | `grep -r`           | (default config)                                              |
| `zoxide`    | `cd`                | Initialised with `--cmd cd` so `cd <fuzzy>` jumps.            |
| `fzf`       | (interactive fuzzy) | Shell integration via `eval "$(fzf --zsh)"`; used by fzf-tab. |
| `delta`     | `diff` / pager      | Used by lazygit and as the git pager.                         |
| `jq` / `yq` | (JSON / YAML)       | —                                                             |

## TUIs

| Tool        | Purpose            | Config                                                                                                                  |
| ----------- | ------------------ | ----------------------------------------------------------------------------------------------------------------------- |
| `lazygit`   | Git porcelain      | `.config/lazygit/config.yml` — cyberdream, delta pager, `nvim --remote` editor.                                         |
| `yazi`      | File manager       | `.config/yazi/theme.toml` (cyberdream), `.config/yazi/yazi.toml`. Wrap with `y` (see [Shell](shell.md)) to inherit cwd. |
| `btop`      | System monitor     | `.config/btop/themes/cyberdream.theme`.                                                                                 |
| `k9s`       | Kubernetes UI      | `.config/k9s/` with cyberdream skin and aliases.                                                                        |
| `fastfetch` | System info banner | Runs at the end of `install.sh`.                                                                                        |

## Dev runtime managers

| Tool           | Purpose                                                                                                   |
| -------------- | --------------------------------------------------------------------------------------------------------- |
| `fnm`          | Node version manager. `eval "$(fnm env --use-on-cd)"` auto-switches on `cd` into a project with `.nvmrc`. |
| `direnv`       | Per-directory env vars via `.envrc`.                                                                      |
| `uv`           | Python project manager (used by this docs site).                                                          |
| `pinentry-mac` | PIN entry GUI used by [ssh-askpass](../scripts/security-keys.md#ssh-askpass) wrapper.                     |

## Color & terminal

| Tool    | Purpose                                                                                                                                                                             |
| ------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `vivid` | Generates `LS_COLORS` from a YAML theme. The cyberdream theme is at `.config/vivid/themes/cyberdream.yaml`; `.zshenv` activates it via `vivid generate ${VIVID_THEME:-cyberdream}`. |
| `chafa` | Image renderer used by `fzf-image-preview` (kitty/sixel/auto passthrough).                                                                                                          |

## GitHub & signing

| Tool                     | Purpose                                                                                                                                                |
| ------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `gh`                     | GitHub CLI. Used by [`git-github-auth`](../scripts/security-keys.md#git-github-auth) and [`git-github-sk`](../scripts/security-keys.md#git-github-sk). |
| `ykman`                  | YubiKey Manager CLI.                                                                                                                                   |
| `yubico-authenticator`   | OATH/TOTP GUI for YubiKey.                                                                                                                             |
| `git-credential-manager` | Credential helper for HTTPS git remotes.                                                                                                               |

## Casks

- **`ghostty`** — terminal emulator (see [Ghostty](ghostty.md)).
- **`hyperkey`** — Caps Lock → hyper key, freeing keybindings for window managers and TUIs.
- **`yubico-authenticator`** — paired with the YubiKey.
