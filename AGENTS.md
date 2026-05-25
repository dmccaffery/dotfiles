# Repository instructions

This repo is a stow-managed macOS dotfiles setup with a companion documentation site at
[`docs/`](docs/), published to <https://dmccaffery.github.io/dotfiles/> via
[`.github/workflows/docs.yaml`](.github/workflows/docs.yaml).

## Keep the docs in sync with the configs

**When you change a configuration file, update the page that documents it in the same PR.**
The docs site is the user-facing surface for this repo; out-of-sync pages mislead readers
worse than missing pages do.

Use this map to find the right page (extend as new sections are added):

| If you change…                                           | Update…                                                               |
| -------------------------------------------------------- | --------------------------------------------------------------------- |
| `install.sh`, `backup.sh`, `setup/**`, `Makefile`        | `docs/getting-started/{install,backup,customize}.md`                  |
| `.config/ghostty/**`                                     | `docs/terminal/ghostty.md`                                            |
| `.config/zsh/**`                                         | `docs/terminal/shell.md`                                              |
| `.config/oh-my-posh/**`                                  | `docs/terminal/oh-my-posh.md`                                         |
| `.config/tmux/**`                                        | `docs/terminal/tmux.md`                                               |
| Brewfiles, new CLI tools                                 | `docs/terminal/cli-tools.md`                                          |
| `.config/nvim/lua/config/lazy.lua`                       | `docs/neovim/lazyvim.md` + `docs/neovim/extras.md`                    |
| `.config/nvim/lua/plugins/**`                            | `docs/neovim/plugins.md`                                              |
| `.config/nvim/lua/config/{autocmds,keymaps,options}.lua` | `docs/neovim/autocmds-keymaps.md`                                     |
| Any new tool's cyberdream theme                          | `docs/theme/per-tool.md`                                              |
| `.config/git/config` and friends                         | `docs/git/config.md` (+ `git-town.md` / `auth.md`)                    |
| Security-key flow scripts                                | `docs/git/signing-security-keys.md` + `docs/scripts/security-keys.md` |
| `.local/share/scripts/**`                                | the relevant `docs/scripts/*.md` page (+ `index.md` table)            |
| `setup/darwin/config.sh`                                 | `docs/macos/system-defaults.md`                                       |
| `Library/LaunchAgents/**`                                | `docs/macos/launchagents.md`                                          |
| `.claude/settings.json`, `.claude/themes/**`             | `docs/claude/{settings,theme}.md`                                     |
| `.stowrc`, `Makefile`, release-please config, linting    | the matching page under `docs/tooling/`                               |

When adding a brand-new top-level area, also wire it into the `nav = […]` block in
[`zensical.toml`](zensical.toml).

## Verify before committing

```sh
make docs-serve                               # uv sync + zensical serve (live reload)
make docs-build                               # uv sync + prettier --write + markdownlint-cli2 + zensical build --clean
```

`make docs-build` is the single gate: it formats `docs/**/*.md` with prettier, lints with
markdownlint-cli2 (must report 0 errors), then builds with zensical (must finish with "No
issues found"). No need to run the linter separately.

To refresh dependency versions: `make docs-upgrade` (wraps `uv sync --upgrade`).

The build runs on every push and PR via GitHub Actions; failures there block deploy.

## Authoring conventions

- **Lead with a one-line summary** of what the page covers.
- **Wrap at 120 chars.** `docs/.markdownlint-cli2.yaml` enforces it (and relaxes a few rules
  that conflict with the pymdownx extensions we use).
- **Tables for key/value mappings** (config options, aliases, scripts).
- **Reference the source-of-truth file path** instead of duplicating large config blocks —
  readers can click through.
- **Icons** must exist in Zensical's bundled set (`.venv/lib/python*/site-packages/zensical/templates/.icons/`).
  Stick to `lucide/*` unless you've verified a `simple/*` or `material/*` icon ships.

## Don't touch

- `CHANGELOG.md` — owned by release-please.
- `.claude/plans/` — runtime artifacts from Claude Code's plan mode.
- `site/`, `.venv/` — build/runtime artifacts (gitignored).

## Useful entry points

- Repo overview & layout: [`docs/index.md`](docs/index.md)
- How the docs are built: [`docs/tooling/contributing.md`](docs/tooling/contributing.md)
- Cyberdream palette source-of-truth: [`docs/assets/extras.css`](docs/assets/extras.css)
