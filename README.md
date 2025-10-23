# Deavons' Dotfiles

> Configuration for macOS terminal environments

## Features

- support for [wezterm][wezterm-url](Default), [iTerm2][iterm2-url], and [alacritty][alacritty-url]
- custom prompt using [oh-my-posh][oh-my-posh-url].
- neovim configuration based on [LazyVim][lazyvim-url]
- tmux configuration
- [catppuccin][catppuccin-url] theme (mocha) everywhere, including terminal prompt, tmux, and neovim
- simple install script that sets up all dependencies
- uses [stow][stow-url] to link farm the configuration

## Setup

If you have any of the existing in your $HOME directory, you will need to move them out of the way first. The install
script does not currently perform a backup. It will refuse to overwrite anything that exists, so no worries there, but
installation will fail as a result.

- $HOME/.config
- $HOME/.zshrc

A simple way to back these up (assuming they are not symlinked / link farmed already):

```sh
mkdir -p "${HOME}/backup"
mv "${HOME}/.config" "${HOME}/backup"
mv "${HOME}/.zshrc" "${HOME}/backup"
```

Once this is done, simply run ./install:

```sh
./install
```

## Future Improvements

- add support for linux environments
- add support for WSL2 environments
- add support for creating backups
- add instructions for creating a fork to make your own dotfiles

[alacritty-url]: https://alacritty.org
[wezterm-url]: https://wezfurlong.org/wezterm/
[iterm2-url]: https://iterm2.com
[lazyvim-url]: https://www.lazyvim.org
[oh-my-posh-url]: https://ohmyposh.dev
[catppuccin-url]: https://catppuccin.com
[stow-url]: https://www.gnu.org/software/stow/
