# Deavons' Dotfiles

> Configuration for macOS terminal environments

## Features

-   support for [wezterm][wezterm-url]
-   custom prompt using [oh-my-posh][oh-my-posh-url].
-   neovim configuration based on [LazyVim][lazyvim-url]
-   tmux configuration
-   [catppuccin][catppuccin-url] theme (mocha) everywhere, including terminal prompt, tmux, and neovim
-   simple install script that sets up all dependencies
-   uses [stow][stow-url] to link farm the configuration
-   extensible git configuration

## Setup

> [!IMPORTANT]
> Please [Create a fork][fork-url] of this repository before using it, even if you do not intend to make any
> customizations for yourself. This repository is primarily meant for me to use for my own needs and I cannot guarantee
> that I will not break something at some point for you. Having a fork of your own ensures that one can manage any
> breaking changes as necessary.

If you have any of the existing in your $HOME directory, you will need to move them out of the way first. The install
script does not currently perform a backup. It will refuse to overwrite anything that exists, so no worries there, but
installation will fail as a result.

-   $HOME/.config/git
-   $HOME/.config/nvim
-   $HOME/.config/oh-my-posh
-   $HOME/.config/scripts
-   $HOME/.config/tmux
-   $HOME/.config/wezterm
-   $HOME/.zshrc

A backup script is provided that can perform this backup, which copies the current config / zshrc to a backup folder
in this repository (ignored by git).

```sh
./backup
```

Once this is done, simply run ./install:

```sh
./install
```

## Customization

The [`start-tmux-session`][sts-url] cmd assumes that you store git repositories in `${HOME}/Repos/` and uses this path
to search for git repositories in order to create named sessions. To set a different path, update the .zshrc to export
a variable called `REPOS_DIR` or modify the [`start-tmux-session`][sts-url] cmd itself to use a different default.

[Extending the git configuration][git-config-include] can be achieved by adding a `~/.config/private/git/config`. I
handle private configuration by using a separate (obviously private) repository that also uses `stow` to symlink the
`${HOME}/.config/private` directory.

> [!WARNING]
> DO NOT PUT CREDENTIALS IN ANY REPOSITORY, EVEN A PRIVATE ONE. The use of a private repository is to hide personal
> information, such as company git urls, email addresses. It is not intended to hide privileged information such as
> credentials!

## Future Improvements

-   add support for linux environments
-   add support for WSL2 environments
-   add support for creating backups
-   add instructions for creating a fork to make your own dotfiles

[wezterm-url]: https://wezfurlong.org/wezterm/
[lazyvim-url]: https://www.lazyvim.org
[oh-my-posh-url]: https://ohmyposh.dev
[catppuccin-url]: https://catppuccin.com
[stow-url]: https://www.gnu.org/software/stow/
[git-config-include]: https://github.com/dmccaffery/dotfiles/blob/main/.config/git/config#L73
[sts-url]: https://github.com/dmccaffery/dotfiles/blob/main/.config/scripts/start-tmux-session#L3
[fork-url]: https://github.com/dmccaffery/dotfiles/fork
