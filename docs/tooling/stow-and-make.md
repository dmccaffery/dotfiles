---
icon: lucide/link
---

# Stow & Make

The repository is laid out so that `stow` can symlink each top-level directory into `$HOME`
verbatim. Files that should _not_ be symlinked are listed in `.stowrc`.

## `.stowrc`

```text title=".stowrc"
--target=~/

--ignore=.DS_Store

--ignore=.github/
--ignore=setup/
--ignore=packages/

--ignore=.*.yaml
--ignore=.*.json
--ignore=.*ignore

--ignore=.editorconfig
--ignore=.gitattributes
--ignore=.gitignore
--ignore=.release-please-manifest.json
--ignore=.stowrc

--ignore=CHANGELOG.md
--ignore=README.md
--ignore=backup.sh
--ignore=install.sh
--ignore=release-please-config.json

--ignore=docs/
--ignore=site/
--ignore=zensical.toml
--ignore=pyproject.toml
--ignore=uv.lock
```

- **`--target=~/`** — `stow .` symlinks every non-ignored top-level entry into `$HOME`.
- **Ignored**: repo metadata, CI config, project-level tooling (markdownlint / prettier),
  installation scripts, and the docs site project.

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
| `make docs-serve`     | `uv sync` + `uv run zensical serve`                                                                                |
| `make docs-build`     | `uv sync` + `uv run zensical build --clean` (output: `./site`)                                                     |
| `make docs-upgrade`   | `uv sync --upgrade` (refresh `uv.lock` to latest matching versions and install)                                    |
| `make requirements`   | `./install.sh requirements` — hidden from help (blank description); call directly if you need just the brew basics |

`install.sh` itself accepts any subset of stages — Make is just the curated front door.

!!! tip "Adding a new target"
Add `## <one-line description>` after the target's colon-line. `make help` picks it up
automatically — no changes to the help target needed. To hide a target from help while
keeping it callable, leave the `## ` empty.

## Layout for stow

```text
dotfiles/
├── .claude/          # → ~/.claude
├── .config/          # → ~/.config
├── .local/           # → ~/.local
├── .ssh/             # → ~/.ssh
├── .terminfo/        # → ~/.terminfo
├── Library/          # → ~/Library
├── .zshrc, .zshenv?  # → ~/  (Zsh files live in .config/zsh/ in this repo,
│                            ZDOTDIR is set system-wide)
├── docs/             # NOT stowed — docs site source
├── setup/            # NOT stowed — installer scripts
├── .github/          # NOT stowed — CI
└── …
```

Every subdirectory under `.claude/`, `.config/`, `.local/`, etc. is symlinked individually,
not the parent. That means partial installs work — you can stow only `.config/nvim` if that's
all you want.

## Re-running stow

```sh
make stow
# or
cd ~/Repos/dotfiles && stow .
```

Stow refuses to overwrite existing files; use `./backup.sh` first if there are conflicts.
To remove symlinks: `stow -D .`.
