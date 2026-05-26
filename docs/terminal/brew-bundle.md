---
icon: simple/homebrew
---

# Brew bundle

Day-2 management of the global Brewfile: adding, removing, and keeping the installed state in sync after the
initial `./install.sh` run.

## How the global Brewfile is wired

The `packages` stage of [`install.sh`](../getting-started/install.md) symlinks the chosen profile from
`.config/homebrew/Brewfile.<profile>` to `$HOMEBREW_BUNDLE_FILE_GLOBAL`
(`~/.local/share/homebrew/Brewfile`) and rewrites the target so it always starts with
[`setup/darwin/Brewfile.requirements`](../../setup/darwin/Brewfile.requirements) followed by the profile's
extras. See [`setup/darwin/packages.sh`](../../setup/darwin/packages.sh) for the full merge.

After install, that symlinked file is the single source of truth — `brew bundle --global` reads and writes it
directly, so edits land in the underlying `Brewfile.<profile>` and get committed alongside any other dotfiles
change.

## Adding a package

```sh
brew bundle add <package> --global
brew bundle install --global
```

`brew bundle add` appends a line to `$HOMEBREW_BUNDLE_FILE_GLOBAL`; the follow-up `install` installs anything
not already present. Pass `--cask`, `--tap`, or `--mas` to add the package as that type instead of a formula:

```sh
brew bundle add --cask ghostty --global
brew bundle add --tap homebrew/cask-fonts --global
```

## Removing a package

```sh
brew bundle remove <package> --global
brew bundle install --global
```

`remove` strips the line from the Brewfile. Because `HOMEBREW_BUNDLE_INSTALL_CLEANUP=1` is exported (see
below), the follow-up `install` actually uninstalls the formula or cask — not just the Brewfile entry.

## The `brewfile` wrapper

[`brewfile`](../scripts/misc.md#brewfile) collapses the two-step flow into one command:

```sh
brewfile add jq
brewfile add --cask ghostty
brewfile remove jq
```

It runs `brew bundle <add|remove> --global` and then `brew bundle install --global`, so the Brewfile and
installed state stay in lockstep.

## Environment variables

The repo exports two `HOMEBREW_BUNDLE_*` variables that shape every `brew bundle` invocation:

| Variable                          | Set in                                             | Effect                                                                                               |
| --------------------------------- | -------------------------------------------------- | ---------------------------------------------------------------------------------------------------- |
| `HOMEBREW_BUNDLE_FILE_GLOBAL`     | [`.config/zsh/.zshenv`](../../.config/zsh/.zshenv) | Resolves `brew bundle … --global` to `~/.local/share/homebrew/Brewfile` (the merged symlink target). |
| `HOMEBREW_BUNDLE_INSTALL_CLEANUP` | [`.config/zsh/.zshenv`](../../.config/zsh/.zshenv) | `brew bundle install` auto-uninstalls anything not listed in the Brewfile (no separate `cleanup`).   |

A third variable is set only during install, scoped to the setup script:

| Variable              | Set in                                                       | Effect                                                                              |
| --------------------- | ------------------------------------------------------------ | ----------------------------------------------------------------------------------- |
| `HOMEBREW_BUNDLE_DIR` | [`setup/darwin/packages.sh`](../../setup/darwin/packages.sh) | Directory the install-time `fzf` picker searches for `Brewfile.<profile>` variants. |

Other `brew bundle` knobs you might reach for ad-hoc (none are exported by default):

- `HOMEBREW_BUNDLE_NO_LOCK=1` — skip writing `Brewfile.lock.json`.
- `HOMEBREW_BUNDLE_DUMP_NO_VSCODE=1` — keep `brew bundle dump` from emitting `vscode "<ext>"` lines.
- `HOMEBREW_BUNDLE_{BREW,CASK,TAP,MAS}_SKIP="<names>"` — exclude packages by name from `install`, `check`, or
  `cleanup`.

Run `brew bundle --help` for the full list.

## Related commands

| Command                             | Use when…                                                                 |
| ----------------------------------- | ------------------------------------------------------------------------- |
| `brew bundle check --global`        | You want a no-op verify that everything in the Brewfile is installed.     |
| `brew bundle cleanup --global`      | You want to preview what would be removed without `install`-ing.          |
| `brew bundle dump --global --force` | You installed something with plain `brew install` and want to capture it. |

## See also

- [Install](../getting-started/install.md) — how the profile + requirements Brewfiles are merged at install time.
- [CLI tools](cli-tools.md) — reference of what the shipped Brewfile installs.
- [Shell](shell.md) — the `.zshenv` block that exports the variables above.
- [`brewfile` script](../scripts/misc.md#brewfile) — wrapper that chains `add`/`remove` + `install`.
