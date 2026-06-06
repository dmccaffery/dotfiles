---
icon: lucide/download
---

# Install

The installer orchestrates a small set of staged shell scripts. Each stage is idempotent and may
be re-run independently.

## Fork first

!!! warning

    Create a fork before installing. This repository is tuned for one person — breaking changes
    land without a deprecation window. Tracking your own fork means you can pin to a known-good
    commit and pull updates on your schedule.

[Fork dmccaffery/dotfiles :octicons-repo-forked-24:](https://github.com/dmccaffery/dotfiles/fork){ .md-button }

## One-shot install

```sh
git clone git@github.com:<your-fork>/dotfiles.git ~/Repos/dotfiles
cd ~/Repos/dotfiles
./backup.sh    # move conflicting configs out of $HOME (optional but recommended)
./install.sh   # run all stages
```

The default invocation runs every stage in order:

```text
xdg → requirements → config → stow → build → packages → shell
```

To run a subset, pass stages as arguments:

```sh
./install.sh requirements stow   # only install brew basics, then symlink configs
```

## Stages

| Stage          | Script                         | Purpose                                                                                                                |
| -------------- | ------------------------------ | ---------------------------------------------------------------------------------------------------------------------- |
| `xdg`          | `setup/xdg.sh`                 | Create the XDG base directories under `$HOME`.                                                                         |
| `requirements` | `setup/darwin/requirements.sh` | Install Xcode CLI tools, Homebrew, and the core Brewfile.                                                              |
| `config`       | `setup/darwin/config.sh`       | Apply macOS system defaults (see [macOS](../macos/system-defaults.md)).                                                |
| `stow`         | `setup/stow.sh`                | Symlink configs from this repo into `$HOME` via GNU stow.                                                              |
| `build`        | `setup/build.sh`               | Build the [`dot`](../tooling/dot.md) Go CLI into `~/.local/bin` and link its applets into `~/.local/share/scripts`.    |
| `packages`     | `setup/darwin/packages.sh`     | Pick a Brewfile profile, merge required packages, confirm the required trust set, then `brew bundle install --global`. |
| `shell`        | `setup/darwin/shell.sh`        | Set Zsh from Homebrew as the default login shell.                                                                      |

The Makefile exposes a target for the full run and one for each stage:

```sh
make backup         # ./backup.sh
make install        # ./install.sh (runs every stage in order)
make xdg            # ./install.sh xdg
make requirements   # ./install.sh requirements
make config         # ./install.sh config
make stow           # ./install.sh stow
make build          # ./install.sh build
make packages       # ./install.sh packages
make shell          # ./install.sh shell
```

## What gets installed

The `requirements` stage installs Xcode CLI tools, Homebrew, and the locked-in baseline from
[`Brewfile.requirements`](../terminal/packages.md). The `packages` stage then
layers a **profile** on top — a Brewfile under `.config/homebrew/Brewfile.*`.

How profile selection works the first time `packages` runs:

1. If `$HOMEBREW_BUNDLE_FILE_GLOBAL` (`~/.local/share/homebrew/Brewfile`) already exists, the
   picker is skipped.
2. Otherwise `fzf` lists every `.config/homebrew/Brewfile.*` and lets you either pick one or
   type a new name. Typing a new name runs `brew bundle dump` of your current machine into
   `.config/homebrew/Brewfile.<name>`.
3. The chosen file is symlinked to `$HOMEBREW_BUNDLE_FILE_GLOBAL`.

After a profile is wired up, `packages.sh` rewrites the target file so it always starts with the
contents of `setup/darwin/Brewfile.requirements` (marked `# required packages -- do not edit`),
followed by everything else from the profile (marked `# profile packages`). Duplicates that are
already in requirements are stripped from the profile section. Then `brew bundle install --global`
syncs the merged file.

!!! warning "Review the required trust set before installing"

    Before the first `brew bundle install`, `packages.sh` prints the non-official taps, formulae, and casks
    trusted by the repo-shipped [`setup/darwin/homebrew/trust.json`](../../setup/darwin/homebrew/trust.json) — the
    set the requirements install loads, since the stage pins `XDG_CONFIG_HOME` to that directory — and waits for
    you to confirm. These run third-party code; if you forked this repo, review the entries (and the
    [trusted-taps list](../terminal/brew-bundle.md#trusted-taps)) and drop anything you don't want before
    answering `y`. Once the install finishes, those entries are merged into your primary
    `~/.config/homebrew/trust.json` so day-2 `brew bundle` runs trust them too.

The shipped profile is [`Brewfile.personal`](../../.config/homebrew/Brewfile.personal) — the
casks (`ghostty`, `hyperkey`, `yubico-authenticator`, the Nerd Fonts) and extras live there.

After everything finishes, `install.sh` runs [`fastfetch`](https://github.com/fastfetch-cli/fastfetch)
as a quick sanity check that the environment is ready.

## Re-running install

`install.sh` is safe to re-run. `brew bundle` is a no-op for already-installed formulae, the
`defaults write` calls are idempotent, and stow will refuse to overwrite existing files (so move
them out of the way with [backup](backup.md) first).
