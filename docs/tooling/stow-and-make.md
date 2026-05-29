---
icon: lucide/link
---

# Stow & Make

The repository is laid out so that `stow` can symlink each top-level directory into `$HOME`
verbatim. Files that should _not_ be symlinked are listed in `.stowrc`.

## `.stowrc`

```text title=".stowrc"
--target=~/

--ignore=DS_Store

--ignore=.github
--ignore=docs
--ignore=packages
--ignore=setup
--ignore=site
--ignore=backups
--ignore=node_modules

--ignore=.*ignore

--ignore=__pycache__

--ignore=.editorconfig
--ignore=.gitattributes
--ignore=.gitignore
--ignore=.release-please-manifest.json
--ignore=.stowrc
--ignore=.venv

--ignore=^.claude/plans
--ignore=^CHANGELOG.md
--ignore=^CLAUDE.md
--ignore=^AGENTS.md
--ignore=^README.md

--ignore=^.commit.sh
--ignore=^backup.sh
--ignore=^install.sh
--ignore=^Makefile
--ignore=^package.json
--ignore=^package-lock.json
--ignore=^release-please-config.json
--ignore=^restore.sh

--ignore=^pyproject.toml
--ignore=^uv.lock
--ignore=^zensical.toml
```

- **`--target=~/`** — `stow .` symlinks every non-ignored top-level entry into `$HOME`.
- **Ignored**: repo metadata, CI config (`.github`), project-level tooling (`package.json`,
  `package-lock.json`, `node_modules`, the `.*ignore` files), installation scripts, the docs
  site project (`docs`, `site`, `.venv`, `zensical.toml`, `pyproject.toml`), and Claude's plan
  artifacts (`^.claude/plans`). Everything else — including `.claude/CLAUDE.md` and
  `.claude/settings.json` — is stowed into `$HOME`.

!!! warning "How `--ignore` matches, and why patterns are anchored"
Each `--ignore=X` compiles to `(?^:(X)\z)` and is tested against the path **relative to the package, anchored only
at the end**. A bare name therefore matches any path that _ends_ with it — a suffix match, not a top-level-only one.
Two consequences shaped the list above:

- A bare `--ignore=CLAUDE.md` catches both the top-level `CLAUDE.md` symlink _and_ `.claude/CLAUDE.md`. The latter is
  meant to stow to `~/.claude/CLAUDE.md` (see [Claude → Memory](../claude/memory.md)), so the top-level entry is
  pinned to the repo root with `^CLAUDE.md`. The transient
  [`.commit.sh`](../claude/memory.md#what-it-covers) is anchored the same way (`^.commit.sh`).
- The old broad `--ignore=.*.json` / `--ignore=.*.yaml` matched _every_ JSON/YAML by that suffix rule. Nested
  configs under `.config/**` survived only because `stow` **folds** a whole directory into a single symlink when the
  target dir doesn't exist yet — the per-file ignore never runs. But `~/.claude` already exists, so `stow`
  **descends** into it and applies the ignore file-by-file, which silently skipped `.claude/settings.json` and
  `.claude/themes/cyberdream.json`. Both patterns were removed; the root lockfiles and manifests they used to cover
  are now ignored by explicit `^`-anchored entries (`^package.json`, `^package-lock.json`, … ).

`.stowrc` strips backslashes when it parses each line, so the anchors use `^` rather than `\A`/`\z`/`\.` (the
surviving `.*ignore` pattern likewise relies on `.`-as-any-char).

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
