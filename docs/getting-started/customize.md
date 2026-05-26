---
icon: lucide/settings
---

# Customize

A handful of environment variables and overlay paths cover the most common customizations
without forking changes upstream.

## `REPOS_DIR`

[`start-tmux-session`](../scripts/tmux.md) (and `attach-tmux-session`) discover Git repositories
by walking a directory tree. `.zshenv` exports `REPOS_DIR="${HOME}/Repos"` as the default;
override it by editing `.config/zsh/.zshenv` (or layering it in via your own private overlay)
to point elsewhere:

```sh
export REPOS_DIR="$HOME/code"
```

!!! tip "Prefer `.zshenv` for exported customizations"
Environment variables (anything you'd `export`) belong in `.zshenv` so they're visible to
every Zsh invocation — scripts, non-interactive shells, and editor subprocesses — not just
interactive sessions. Reserve `.zshrc` for things that only matter when a human is typing
(aliases, key bindings, prompt init, plugin loaders).

## Private Git overlay

The shipped `~/.config/git/config` ends with:

```ini
[include]
    path = ~/.config/private/git/config
```

Put your `user.name`, `user.email`, and any company-specific Git config in
`~/.config/private/git/config`. Common pattern: keep this file in a separate, private repository
that uses `stow` to symlink into `~/.config/private`.

!!! warning "No credentials in private repos"
Even private repositories are not a credential store. Use the system keychain
(`git-credential-manager`) or a hardware token. The private overlay is for _personal_
information (email, work URLs), not secrets.

## Theme variants

The cyberdream palette has both dark and light variants; the system theme drives Ghostty's
choice automatically:

```text title=".config/ghostty/config"
theme = dark:cyberdream,light:cyberdream-light
```

Other tools follow whichever variant they were configured for in this repo. See
[Theme → per-tool](../theme/per-tool.md) for the matrix.

## Disabling tools you don't use

The Brewfile and nvim plugins are a take-it-or-leave-it bundle, but the easy escape hatches are:

- **Brewfile**: edit your profile under `.config/homebrew/Brewfile.<profile>` (e.g.
  `Brewfile.personal`) and re-run `make packages`. The required baseline in
  `setup/darwin/Brewfile.requirements` is reapplied every run, so removing a required package
  from the profile won't drop it — edit the requirements file instead if you really mean to
  remove something everyone gets.
- **NeoVim extras**: edit `.config/nvim/lua/config/lazy.lua` and remove `extras.lang.*` imports
  for languages you don't touch — startup time drops noticeably.
- **Tmux plugins**: edit `.config/tmux/conf/plugins.conf`.

## NeoVim dashboard header

The Snacks dashboard header in [`.config/nvim/lua/plugins/header.lua`](../neovim/plugins.md#headerlua)
is a plain heredoc string — swap it for any ASCII art you like. The shipped art was generated
with [`figlet`](http://www.figlet.org/) (available via `brew install figlet`):

```sh
figlet -f graffiti -w 80 "Deavon's Terminal" | awk '{printf "%-80s\n", $0}'
```

`-w 80` sets the output width so the rendered art wraps to two lines instead of one very long
one. The `awk` pipe right-pads every line to exactly 80 characters — Snacks centres each line
independently, so unequal line lengths will visibly stagger the art. Pick a font with
`figlet -f <font> -w 80 "preview"` after browsing `ls $(brew --prefix)/share/figlet/fonts`,
or use `showfigfonts` to list every installed font with a sample. Paste the output between the
`[[` and `]]` delimiters in `header.lua` — keep the surrounding indentation so Lua's
long-string literal stays valid, and the dashboard picks it up on the next NeoVim launch.

## Profile shell startup

Suspect Zsh is slow? Use the included profiler:

```sh
profile-shell
```

This wraps `time -p zsh -i -c exit` with `ZSHPROFILE=1`, which causes `.zshrc` to load
`zsh/zprof` and dump a timing table on exit. See [scripts/misc](../scripts/misc.md).
