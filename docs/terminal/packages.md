---
icon: lucide/package
---

# Packages

[`setup/darwin/Brewfile.requirements`](https://github.com/dmccaffery/dotfiles/blob/main/setup/darwin/Brewfile.requirements)
is the locked-in baseline of formulae, taps, and casks that every profile inherits.
The `packages` stage of [`install.sh`](../getting-started/install.md) merges it into the
chosen `.config/homebrew/Brewfile.<profile>` before each `brew bundle install --global`, so
removing one of these entries from the active Brewfile won't actually uninstall it on the
next sync — [`packages.sh`](https://github.com/dmccaffery/dotfiles/blob/main/setup/darwin/packages.sh)
re-inserts every entry below.

Treat this list as the floor: anything every machine needs, regardless of profile.
Profile-specific picks (extra casks, app-store apps, language SDKs) live in
[`.config/homebrew/Brewfile.personal`](https://github.com/dmccaffery/dotfiles/blob/main/.config/homebrew/Brewfile.personal)
or your own fork's profile — those aren't documented here because they're personal by design.

Each row links the package name to its Homebrew formula/cask page, plus its upstream source
so credit lands where it belongs.

## Formulae

### Core CLI

| Package                                               | Upstream                                                    | Purpose                                                       |
| ----------------------------------------------------- | ----------------------------------------------------------- | ------------------------------------------------------------- |
| [`bat`](https://formulae.brew.sh/formula/bat)         | [sharkdp/bat](https://github.com/sharkdp/bat)               | Syntax-highlighted `cat` replacement.                         |
| [`fd`](https://formulae.brew.sh/formula/fd)           | [sharkdp/fd](https://github.com/sharkdp/fd)                 | Fast user-friendly `find` replacement.                        |
| [`ripgrep`](https://formulae.brew.sh/formula/ripgrep) | [BurntSushi/ripgrep](https://github.com/BurntSushi/ripgrep) | Fast recursive grep.                                          |
| [`lsd`](https://formulae.brew.sh/formula/lsd)         | [lsd-rs/lsd](https://github.com/lsd-rs/lsd)                 | `ls` replacement with icons; aliased as `ls` in `.zshrc`.     |
| [`zoxide`](https://formulae.brew.sh/formula/zoxide)   | [ajeetdsouza/zoxide](https://github.com/ajeetdsouza/zoxide) | Smarter `cd` that learns jump targets.                        |
| [`fzf`](https://formulae.brew.sh/formula/fzf)         | [junegunn/fzf](https://github.com/junegunn/fzf)             | Interactive fuzzy finder; powers `Ctrl-T`, `Ctrl-R`, fzf-tab. |
| [`jq`](https://formulae.brew.sh/formula/jq)           | [jqlang/jq](https://github.com/jqlang/jq)                   | JSON query/transform.                                         |
| [`yq`](https://formulae.brew.sh/formula/yq)           | [mikefarah/yq](https://github.com/mikefarah/yq)             | YAML/XML/TOML query/transform.                                |
| [`curl`](https://formulae.brew.sh/formula/curl)       | [curl/curl](https://github.com/curl/curl)                   | HTTP client (Homebrew's, kept ahead of Apple's bundled copy). |
| [`wget`](https://formulae.brew.sh/formula/wget)       | [GNU wget](https://savannah.gnu.org/projects/wget/)         | HTTP/FTP downloader.                                          |

### Editor & shell

| Package                                             | Upstream                                                              | Purpose                                                                      |
| --------------------------------------------------- | --------------------------------------------------------------------- | ---------------------------------------------------------------------------- |
| [`neovim`](https://formulae.brew.sh/formula/neovim) | [neovim/neovim](https://github.com/neovim/neovim)                     | The editor. Configured via LazyVim — see [NeoVim](../neovim/index.md).       |
| [`zsh`](https://formulae.brew.sh/formula/zsh)       | [zsh-users/zsh](https://github.com/zsh-users/zsh)                     | Default shell; Homebrew build set as login shell by `setup/darwin/shell.sh`. |
| [`zinit`](https://formulae.brew.sh/formula/zinit)   | [zdharma-continuum/zinit](https://github.com/zdharma-continuum/zinit) | Zsh plugin manager — see [Shell](shell.md).                                  |
| [`tmux`](https://formulae.brew.sh/formula/tmux)     | [tmux/tmux](https://github.com/tmux/tmux)                             | Terminal multiplexer — see [tmux](tmux.md).                                  |
| [`lua`](https://formulae.brew.sh/formula/lua)       | [lua.org](https://www.lua.org/)                                       | Lua runtime required by NeoVim plugins and oh-my-posh.                       |

### Git ecosystem

| Package                                                               | Upstream                                                            | Purpose                                                                                                                                                                                    |
| --------------------------------------------------------------------- | ------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [`git`](https://formulae.brew.sh/formula/git)                         | [git/git](https://github.com/git/git)                               | Git CLI.                                                                                                                                                                                   |
| [`gh`](https://formulae.brew.sh/formula/gh)                           | [cli/cli](https://github.com/cli/cli)                               | GitHub CLI used by [`git-github-auth`](../scripts/security-keys.md#git-github-auth) and [`git-github-sk`](../scripts/security-keys.md#git-github-sk).                                      |
| [`git-delta`](https://formulae.brew.sh/formula/git-delta)             | [dandavison/delta](https://github.com/dandavison/delta)             | Side-by-side diff renderer (configured as the git pager).                                                                                                                                  |
| [`git-filter-repo`](https://formulae.brew.sh/formula/git-filter-repo) | [newren/git-filter-repo](https://github.com/newren/git-filter-repo) | Fast history rewriter.                                                                                                                                                                     |
| [`git-lfs`](https://formulae.brew.sh/formula/git-lfs)                 | [git-lfs/git-lfs](https://github.com/git-lfs/git-lfs)               | Git Large File Storage (used for the wallpapers under `.local/share/wallpapers/`).                                                                                                         |
| [`git-town`](https://formulae.brew.sh/formula/git-town)               | [git-town/git-town](https://github.com/git-town/git-town)           | Branch-chain workflow — see [git-town](../git/git-town.md).                                                                                                                                |
| [`lazygit`](https://formulae.brew.sh/formula/lazygit)                 | [jesseduffield/lazygit](https://github.com/jesseduffield/lazygit)   | Git TUI; configured with `delta` as pager and `nvim --remote` as editor (see [`.config/lazygit/config.yml`](https://github.com/dmccaffery/dotfiles/blob/main/.config/lazygit/config.yml)). |

### SSH & security

| Package                                                         | Upstream                                                                | Purpose                                                                                                              |
| --------------------------------------------------------------- | ----------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| [`openssh`](https://formulae.brew.sh/formula/openssh)           | [openssh/openssh-portable](https://github.com/openssh/openssh-portable) | Homebrew OpenSSH — paired with [`org.homebrew.ssh-agent`](../macos/launchagents.md) so `SSH_ASKPASS` actually works. |
| [`openssl@3`](https://formulae.brew.sh/formula/openssl@3)       | [openssl/openssl](https://github.com/openssl/openssl)                   | OpenSSL 3.x; required by the Homebrew OpenSSH build.                                                                 |
| [`pinentry-mac`](https://formulae.brew.sh/formula/pinentry-mac) | [GPGTools/pinentry-mac](https://github.com/GPGTools/pinentry-mac)       | PIN-entry GUI used by the [`ssh-askpass`](../scripts/security-keys.md#ssh-askpass) wrapper.                          |
| [`ykman`](https://formulae.brew.sh/formula/ykman)               | [Yubico/yubikey-manager](https://github.com/Yubico/yubikey-manager)     | YubiKey Manager CLI.                                                                                                 |

### Media & image

| Package                                                       | Upstream                                                                  | Purpose                                                                                      |
| ------------------------------------------------------------- | ------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- |
| [`chafa`](https://formulae.brew.sh/formula/chafa)             | [hpjansson/chafa](https://github.com/hpjansson/chafa)                     | Terminal image renderer used by [`fzf-image-preview`](../scripts/misc.md#fzf-image-preview). |
| [`ffmpeg`](https://formulae.brew.sh/formula/ffmpeg)           | [FFmpeg/FFmpeg](https://github.com/FFmpeg/FFmpeg)                         | Video/audio toolkit; required by yazi for media previews.                                    |
| [`imagemagick`](https://formulae.brew.sh/formula/imagemagick) | [ImageMagick/ImageMagick](https://github.com/ImageMagick/ImageMagick)     | Image manipulation; required by yazi for image previews.                                     |
| [`poppler`](https://formulae.brew.sh/formula/poppler)         | [freedesktop.org/poppler](https://gitlab.freedesktop.org/poppler/poppler) | PDF rendering library; required by yazi for PDF previews.                                    |
| [`resvg`](https://formulae.brew.sh/formula/resvg)             | [linebender/resvg](https://github.com/linebender/resvg)                   | SVG rasterizer; required by yazi for SVG previews.                                           |

### TUIs & system info

| Package                                                   | Upstream                                                              | Purpose                                                                         |
| --------------------------------------------------------- | --------------------------------------------------------------------- | ------------------------------------------------------------------------------- |
| [`btop`](https://formulae.brew.sh/formula/btop)           | [aristocratos/btop](https://github.com/aristocratos/btop)             | System monitor.                                                                 |
| [`fastfetch`](https://formulae.brew.sh/formula/fastfetch) | [fastfetch-cli/fastfetch](https://github.com/fastfetch-cli/fastfetch) | System info banner; runs at the end of `install.sh`.                            |
| [`yazi`](https://formulae.brew.sh/formula/yazi)           | [sxyazi/yazi](https://github.com/sxyazi/yazi)                         | TUI file manager; wrap with `y` (see [Shell](shell.md)) to inherit cwd on exit. |
| [`vivid`](https://formulae.brew.sh/formula/vivid)         | [sharkdp/vivid](https://github.com/sharkdp/vivid)                     | Generates `LS_COLORS` from a cyberdream YAML theme.                             |

### Utilities

| Package                                         | Upstream                                       | Purpose                                                                |
| ----------------------------------------------- | ---------------------------------------------- | ---------------------------------------------------------------------- |
| [`stow`](https://formulae.brew.sh/formula/stow) | [GNU Stow](https://www.gnu.org/software/stow/) | Symlink-farm manager — see [Stow & Make](../tooling/stow-and-make.md). |

### Third-party formulae

Pulled from non-`homebrew/core` taps, so they're not indexed on `formulae.brew.sh`:

| Package                                | Upstream                                                                  | Purpose                                                                                                     |
| -------------------------------------- | ------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------- |
| `derailed/k9s/k9s`                     | [derailed/k9s](https://github.com/derailed/k9s)                           | Kubernetes TUI.                                                                                             |
| `jandedobbeleer/oh-my-posh/oh-my-posh` | [JanDeDobbeleer/oh-my-posh](https://github.com/JanDeDobbeleer/oh-my-posh) | Prompt engine — see [oh-my-posh](oh-my-posh.md). Pinned with `args: ["formula"]` to avoid the cask variant. |

## Casks

### Fonts

| Cask                                                                                       | Upstream                                                        | Purpose                                                |
| ------------------------------------------------------------------------------------------ | --------------------------------------------------------------- | ------------------------------------------------------ |
| [`font-fira-code-nerd-font`](https://formulae.brew.sh/cask/font-fira-code-nerd-font)       | [ryanoasis/nerd-fonts](https://github.com/ryanoasis/nerd-fonts) | FiraCode Nerd Font (alternative coding font).          |
| [`font-iosevka-nerd-font`](https://formulae.brew.sh/cask/font-iosevka-nerd-font)           | [ryanoasis/nerd-fonts](https://github.com/ryanoasis/nerd-fonts) | Iosevka Nerd Font — the configured Ghostty font.       |
| [`font-symbols-only-nerd-font`](https://formulae.brew.sh/cask/font-symbols-only-nerd-font) | [ryanoasis/nerd-fonts](https://github.com/ryanoasis/nerd-fonts) | Symbols-only Nerd Font (icon glyphs for non-NF fonts). |

### Apps

| Cask                                                                             | Upstream                                                                                        | Purpose                                                          |
| -------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------- | ---------------------------------------------------------------- |
| [`ghostty`](https://formulae.brew.sh/cask/ghostty)                               | [ghostty-org/ghostty](https://github.com/ghostty-org/ghostty)                                   | Terminal emulator — see [Ghostty](ghostty.md).                   |
| [`git-credential-manager`](https://formulae.brew.sh/cask/git-credential-manager) | [git-ecosystem/git-credential-manager](https://github.com/git-ecosystem/git-credential-manager) | Credential helper for HTTPS git remotes (Codeberg/GitLab OAuth). |
| [`hyperkey`](https://formulae.brew.sh/cask/hyperkey)                             | [hyperkey.app](https://hyperkey.app/)                                                           | Caps Lock → hyper modifier.                                      |
| [`yubico-authenticator`](https://formulae.brew.sh/cask/yubico-authenticator)     | [Yubico/yubioath-flutter](https://github.com/Yubico/yubioath-flutter)                           | OATH/TOTP GUI for the YubiKey.                                   |

## See also

- [Install](../getting-started/install.md) — how `requirements` and `packages` stages run.
- [Brew bundle](brew-bundle.md) — day-2 add/remove flow once everything is installed.
- [Per-tool theming](../theme/per-tool.md) — where each tool's cyberdream config lives.
