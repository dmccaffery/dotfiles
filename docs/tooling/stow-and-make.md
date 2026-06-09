---
icon: lucide/link
---

# Stow & Make

Everything that gets symlinked into `$HOME` lives under a single
[`stow/`](https://github.com/dmccaffery/dotfiles/tree/main/stow) directory. Each tree inside it —
`stow/.config`, `stow/.claude`, `stow/.local`, … — is stowed into its own `$HOME` target by
[`setup/stow.sh`](../../setup/stow.sh). Anything that should _not_ land in `$HOME` (the docs site,
installer scripts, the Go sources, CI config) simply lives outside `stow/`, so `stow` never sees it.

## How `stow.sh` links each tree

`setup/stow.sh` does not run a single `stow .` over the repo root. It invokes `stow` once per tree,
each with an explicit `--dir` (the source under `stow/`) and `--target` (the matching directory in
`$HOME`):

```sh title="setup/stow.sh (excerpt)"
stow --dir="${INSTALL_DIR}/stow/.claude" --target="${HOME}/.claude" .
stow --dir="${INSTALL_DIR}/stow/.config" --target="${HOME}/.config" .
stow --dir="${INSTALL_DIR}/stow/.local" --target="${HOME}/.local" .
stow --dir="${INSTALL_DIR}/stow/.terminfo" --target="${HOME}/.terminfo" .
stow --dir="${INSTALL_DIR}/stow/.ssh" --target="${HOME}/.ssh" .
stow --dir="${INSTALL_DIR}/stow/Library" --target="${HOME}/Library" .
```

Two details make this safe:

- **Targets that hold runtime state are pre-created as real directories.** Before stowing,
  `stow.sh` `mkdir -p`s `~/.claude`, `~/.config/{codex,opencode,zsh}`, `~/.local/share/{scripts,wallpapers}`,
  and `~/.ssh`. Because the target already exists, `stow` links each child individually instead of
  folding the whole tree into one top-level symlink — keeping `~/.config`, `~/.claude`, `~/.ssh`
  as real directories so apps can keep writing into them (Claude Code's `~/.claude/projects/`, your
  own `~/.ssh/known_hosts`) without it landing in the repo. The trees with no runtime state
  (`~/.terminfo`, `~/Library`) are left to fold normally.
- **`.claude/plans` is not part of the stow source.** Claude Code's plan-mode runtime is written to
  the repo-root `.claude/plans/` (per the `plansDirectory` setting, which is relative to each repo),
  so it lives outside `stow/.claude/` entirely and is never symlinked into `$HOME`. It stays
  gitignored.

`~/.zshenv` is the one exception: it is a single file linked directly with `ln -Ffs`, not via `stow`
(Zsh must read it before `ZDOTDIR` is known).

## Makefile

`make` with no arguments prints the help (also `make help`). Each target carries a `##`
description that the help target parses with `awk`:

```makefile title="Makefile"
.DEFAULT_GOAL := help

.PHONY: help
help: ## Print this help
    @awk '...' $(MAKEFILE_LIST)   # see Makefile for the full awk script

.PHONY: requirements
requirements: ## Install Xcode CLI tools, Homebrew, and the base Brewfile
    ./install.sh requirements

# ... etc
```

| Target                | Maps to                                                                                                            |
| --------------------- | ------------------------------------------------------------------------------------------------------------------ |
| `make help` (default) | Self-documenting target list                                                                                       |
| `make packages`       | `./install.sh packages`                                                                                            |
| `make stow`           | `./install.sh stow`                                                                                                |
| `make build`          | `./install.sh build` — build the [`dot`](dot.md) CLI and link its applets                                          |
| `make fmt`            | `npm install` + `npx prettier --write 'docs/**/*.md'`                                                              |
| `make lint`           | `fmt` first, then `shellcheck --severity=warning` over every shell script + `markdownlint-cli2 'docs/**/*.md'`     |
| `make docs-serve`     | `uv sync` + `uv run zensical serve`                                                                                |
| `make docs-build`     | `lint` first, then `uv sync` + `uv run zensical build --clean` (output: `./site`)                                  |
| `make upgrade`        | Prompted bypass of dependabot cooldown — `npm update` + `uv sync --upgrade` (refreshes both lock files)            |
| `make requirements`   | `./install.sh requirements` — hidden from help (blank description); call directly if you need just the brew basics |

`install.sh` itself accepts any subset of stages — Make is just the curated front door.

!!! tip "Adding a new target"

    Add `## <one-line description>` after the target's colon-line. `make help` picks it up
    automatically — no changes to the help target needed. To hide a target from help while
    keeping it callable, leave the `## ` empty.

## Layout

```text
dotfiles/
├── stow/               # everything below is symlinked into $HOME
│   ├── .claude/        # → ~/.claude   (settings.json, themes, CLAUDE.md)
│   ├── .config/        # → ~/.config
│   ├── .local/         # → ~/.local
│   ├── .ssh/           # → ~/.ssh
│   ├── .terminfo/      # → ~/.terminfo
│   └── Library/        # → ~/Library
├── .zshenv             # → ~/.zshenv  (linked with ln, not stow)
├── setup/              # NOT stowed — installer scripts
├── docs/               # NOT stowed — docs site source
├── cmd/, internal/     # NOT stowed — the `dot` Go CLI
├── .github/            # NOT stowed — CI
└── …
```

Every subdirectory under `stow/.config/`, `stow/.claude/`, etc. is symlinked individually, not the
parent. That means partial installs work — you can stow a single tree, e.g.
`stow --dir=stow/.config --target=~/.config .` to link only your `~/.config`.

## Re-running stow

```sh
make stow
```

`make stow` re-runs every per-tree `stow` invocation shown above. Stow refuses to overwrite existing
files, so run [`./backup.sh`](../getting-started/backup.md) first if there are conflicts. To remove the
symlinks, [`./restore.sh`](../getting-started/backup.md#restoring-a-backup) runs the matching `stow -D`
over each tree — or do it by hand with `stow -D --dir=stow/<tree> --target=~/<target> .`.
