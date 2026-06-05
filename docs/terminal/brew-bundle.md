---
icon: simple/homebrew
---

# Brew bundle

Day-2 management of the global Brewfile: adding, removing, and keeping the installed state in sync after the
initial `./install.sh` run.

## How the global Brewfile is wired

The `packages` stage of [`install.sh`](../getting-started/install.md) symlinks the chosen profile from
`.config/homebrew/Brewfile.<profile>` to `$HOMEBREW_BUNDLE_FILE_GLOBAL`
(`~/.local/share/homebrew/Brewfile`) and rewrites the target so it always starts with the
[`Brewfile.requirements`](packages.md) baseline followed by the profile's extras. See
[`setup/darwin/packages.sh`](https://github.com/dmccaffery/dotfiles/blob/main/setup/darwin/packages.sh)
for the full merge.

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

`remove` strips the line from the Brewfile. Because `HOMEBREW_BUNDLE_FORCE_INSTALL_CLEANUP=1` is exported (see
below), the follow-up `install` actually uninstalls the formula or cask — not just the Brewfile entry — without
pausing to confirm the cleanup.

## The `brewfile` wrapper

[`brewfile`](../scripts/misc.md#brewfile) collapses the two-step flow into one command:

```sh
brewfile add jq
brewfile add --cask ghostty
brewfile remove jq
```

It runs `brew bundle <add|remove> --global` and then `brew bundle install --global --zap --upgrade`, so the
Brewfile and installed state stay in lockstep — `--zap` fully removes any casks dropped from the Brewfile (their
support files included) instead of leaving them behind, and `--upgrade` brings outdated formulae up to date in the
same pass.

On `add`, the wrapper also guards the [trust check](#trusted-taps): any referenced formula, cask, or tap that comes
from a non-official tap (a `user/tap/...` path) is looked up in `trust.json`, and if it is missing you are prompted
to `brew trust --<type> <name>` it before the install runs — so a non-official entry never blocks `brew bundle
install` mid-run. Bare names resolve to `homebrew/core`/`homebrew/cask` and are trusted by default, so they are
skipped without a prompt.

## Environment variables

The repo exports two `HOMEBREW_BUNDLE_*` variables that shape every `brew bundle` invocation:

| Variable                                | Set in                                             | Effect                                                                                                                                                                               |
| --------------------------------------- | -------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `HOMEBREW_BUNDLE_FILE_GLOBAL`           | [`.config/zsh/.zshenv`](../../.config/zsh/.zshenv) | Resolves `brew bundle … --global` to `~/.local/share/homebrew/Brewfile` (the merged symlink target).                                                                                 |
| `HOMEBREW_BUNDLE_FORCE_INSTALL_CLEANUP` | [`.config/zsh/.zshenv`](../../.config/zsh/.zshenv) | With `--global`, `brew bundle install` runs cleanup automatically — uninstalling anything not in the Brewfile — and skips the confirmation prompt (the `--force-cleanup` behaviour). |

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

## Trusted taps

Recent Homebrew gates formulae, casks, and commands from **non-official taps** behind a trust check
(`$HOMEBREW_REQUIRE_TAP_TRUST`, on its way to becoming the default). Until a tap is trusted, `brew bundle install`
refuses to load anything from it. The trusted set lives in
[`.config/homebrew/trust.json`](../../.config/homebrew/trust.json) — Homebrew reads it from
`$XDG_CONFIG_HOME/homebrew/trust.json` — and is managed with `brew trust --tap <name>` / `brew untrust --tap <name>`.
The [`brewfile`](#the-brewfile-wrapper) wrapper also offers to run `brew trust` for you when an `add` references a
non-official tap, formula, or cask that is not yet listed here.

The install flow keeps a second, install-time trust set in
[`setup/darwin/homebrew/trust.json`](../../setup/darwin/homebrew/trust.json) covering exactly what
[`Brewfile.requirements`](packages.md) pulls from non-official taps. The [`packages`](../getting-started/install.md)
stage pins `XDG_CONFIG_HOME` to that directory so the requirements install reads it instead of your primary file,
prints its entries for you to confirm, and — once the install finishes — merges them into
`~/.config/homebrew/trust.json` so day-2 runs trust them too.

This repo pre-trusts exactly the third-party taps its Brewfiles pull from, so a fresh `brew bundle install --global`
is never blocked mid-run:

| Tap                         | Provides        |
| --------------------------- | --------------- |
| `anomalyco/tap`             | `opencode`      |
| `bufbuild/buf`              | `buf`           |
| `controlplaneio-fluxcd/tap` | `flux-operator` |
| `derailed/k9s`              | `k9s`           |
| `fluxcd/tap`                | `flux`          |
| `hashicorp/tap`             | `terraform`     |
| `jandedobbeleer/oh-my-posh` | `oh-my-posh`    |
| `oven-sh/bun`               | `bun`           |

!!! warning "Review before you install"

    These taps run code from third parties. If you fork this repo, audit `trust.json` and drop any tap you don't
    actually want before running `./install.sh`. The [`packages`](../getting-started/install.md) stage separately
    prints the repo-shipped required trust set and asks you to confirm it before the first `brew bundle install`.

## Related commands

| Command                             | Use when…                                                                 |
| ----------------------------------- | ------------------------------------------------------------------------- |
| `brew bundle check --global`        | You want a no-op verify that everything in the Brewfile is installed.     |
| `brew bundle cleanup --global`      | You want to preview what would be removed without `install`-ing.          |
| `brew bundle dump --global --force` | You installed something with plain `brew install` and want to capture it. |

## See also

- [Install](../getting-started/install.md) — how the profile + requirements Brewfiles are merged at install time.
- [Packages](packages.md) — full reference of what `Brewfile.requirements` installs.
- [Shell](shell.md) — the `.zshenv` block that exports the variables above.
- [`brewfile` script](../scripts/misc.md#brewfile) — wrapper that chains `add`/`remove` + `install`.
- [Install](../getting-started/install.md) — the `packages` stage that confirms the [trusted taps](#trusted-taps)
  before the first `brew bundle install`.
