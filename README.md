# Deavons' Dotfiles

> Configuration for macOS terminal environments

📖 **Full documentation:** <https://dmccaffery.github.io/dotfiles/>

## Demo

[![Watch the demo: tmux, NeoVim, lazygit, oh-my-posh, and a Claude Code agent in cyberdream][demo-poster]][demo-url]

A recorded walkthrough of the environment — tmux, NeoVim, lazygit, oh-my-posh, and a Claude Code agent, all in
cyberdream. The thumbnail links to the playable recording on the [docs site][demo-url].

## Features

- support for [ghostty][ghostty-url]
- custom prompt using [oh-my-posh][oh-my-posh-url]
- neovim configuration based on [LazyVim][lazyvim-url]
- tmux configuration
- [cyberdream][cyberdream-url] theme everywhere, including oh-my-posh, tmux, nvim, yazi, k9s, ghostty, and lazygit
- staged install script that sets up all dependencies
- uses [stow][stow-url] to link farm the configuration
- extensible git configuration

## Quick start

> [!IMPORTANT]
> [Create a fork][fork-url] before using this repo, even if you don't plan to customize it. This repo is primarily for
> my own use and may break without warning — a fork lets you manage breaking changes on your own schedule. See the
> [forking guide][forking-url] for the recommended workflow.

```sh
./backup.sh   # move existing dotfiles aside (stow refuses to overwrite)
./install.sh  # stage-by-stage, idempotent installer
```

See [Getting Started][getting-started-url] for the full walkthrough.

## Future Improvements

- add support for linux environments
- add support for WSL2 environments

[demo-poster]: docs/assets/images/demo-poster.png
[demo-url]: https://dmccaffery.github.io/dotfiles/#demo
[cyberdream-url]: https://github.com/scottmckendry/cyberdream.nvim
[ghostty-url]: https://ghostty.org/
[lazyvim-url]: https://www.lazyvim.org
[oh-my-posh-url]: https://ohmyposh.dev
[stow-url]: https://www.gnu.org/software/stow/
[fork-url]: https://github.com/dmccaffery/dotfiles/fork
[forking-url]: https://dmccaffery.github.io/dotfiles/getting-started/install/#fork-first
[getting-started-url]: https://dmccaffery.github.io/dotfiles/getting-started/
