# Repository instructions

> **Editing this file:** `CLAUDE.md` is a symlink to `AGENTS.md`. Tools that refuse to write
> through symlinks (Claude Code's `Edit`/`Write`, etc.) will error out on the `CLAUDE.md` path —
> always edit `AGENTS.md` directly. The symlink picks up the change automatically.

This repo is a stow-managed macOS dotfiles setup with a companion documentation site at
[`docs/`](docs/), published to <https://dmccaffery.github.io/dotfiles/> via the
`publish-docs` job in [`.github/workflows/release-main.yaml`](.github/workflows/release-main.yaml),
which runs when release-please cuts a release.
[`.github/workflows/pull-request.yaml`](.github/workflows/pull-request.yaml) is the PR
gate: it runs `shellcheck` over every shell script and a non-deploying docs build
(`prettier --check` + `markdownlint-cli2` + `zensical build`) so breakage surfaces before
merge.

## Keep the docs in sync with the configs

**When you change a configuration file, update the page that documents it in the same PR.**
The docs site is the user-facing surface for this repo; out-of-sync pages mislead readers
worse than missing pages do.

Use this map to find the right page (extend as new sections are added):

| If you change…                                           | Update…                                                               |
| -------------------------------------------------------- | --------------------------------------------------------------------- |
| `install.sh`, `backup.sh`, `setup/**`, `Makefile`        | `docs/getting-started/{install,backup,customize}.md`                  |
| `package.json`, `.prettierignore`, shellcheck rules      | `docs/tooling/linting.md`                                             |
| `.config/ghostty/**`                                     | `docs/terminal/ghostty.md`                                            |
| `.config/zsh/**`                                         | `docs/terminal/shell.md`                                              |
| `.config/oh-my-posh/**`                                  | `docs/terminal/oh-my-posh.md`                                         |
| `.config/tmux/**`                                        | `docs/terminal/tmux.md`                                               |
| `.config/opencode/{opencode.jsonc,tui.json,AGENTS.md}`   | `docs/terminal/opencode.md`                                           |
| `.config/opencode/themes/**`                             | `docs/theme/per-tool.md` + `docs/terminal/opencode.md`                |
| `setup/darwin/Brewfile.requirements`                     | `docs/terminal/packages.md`                                           |
| `brew bundle` lifecycle / `HOMEBREW_BUNDLE_*` env vars   | `docs/terminal/brew-bundle.md`                                        |
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
| `.claude/CLAUDE.md` (user-level agent memory)            | `docs/claude/memory.md`                                               |
| `.stowrc`, `Makefile`, release-please config, linting    | the matching page under `docs/tooling/`                               |
| `hack/**`, `docs/assets/demo.cast`                       | `docs/tooling/contributing.md`                                        |

When adding a brand-new top-level area, also wire it into the `nav = […]` block in
[`zensical.toml`](zensical.toml).

### Doc-sync check before every commit

Before creating any commit, inspect the staged changes and decide whether the docs need to
move with them. Don't rely on the user to remember — this is a hard prerequisite, not a
courtesy.

1. Run `git diff --staged --name-only` (and `git diff --staged` when the file list isn't
   enough to know what changed).
2. For each staged path, look it up in the map above. If it matches a row, the listed
   doc page is in scope and must be reviewed in the same commit.
3. Read the current doc page and confirm it still matches the new config. Update it when
   the change adds, removes, renames, or alters the behaviour of anything the page
   describes (env vars, keybinds, options, file paths, tables, code snippets).
4. If the staged set touches a path that isn't in the map but is user-facing (a new
   top-level area, a new script, a new tool), extend the map row in this file and add the
   page in the same commit.
5. Stage the doc updates and re-run `make lint` (or `make docs-build` for new pages)
   before finalising the commit.

Skip the doc step only when the diff is provably docs-irrelevant (e.g., a change to
`.claude/plans/`, a release-please-owned file, or a pure formatting pass already covered
by `make fmt`). When in doubt, update the page — out-of-sync docs are worse than missing
ones, per the rule above.

## Verify before committing

Never invoke `prettier` or `markdownlint-cli2` directly in this repo. Always use the
Makefile targets below so formatting and linting run through the same path as CI.

```sh
make fmt                                      # npm install + prettier --write
make lint                                     # fmt + shellcheck + markdownlint-cli2
make docs-serve                               # uv sync + zensical serve (live reload)
make docs-build                               # lint + uv sync + zensical build --clean
```

`make docs-build` is the single gate: it runs `lint` (which depends on `fmt`, so prettier
formats first, then shellcheck and markdownlint-cli2 both report 0 errors), then builds
with zensical (must finish with "No issues found").

To refresh dependency versions: `make upgrade` (runs `npm update` then `uv sync --upgrade`
after a typed confirmation). The target prints a warning that you're bypassing the 7-day
dependabot cooldown — prefer merging the dependabot PR when possible.

PRs run `.github/workflows/pull-request.yaml`, which shellchecks every shell script and
re-runs the docs build (`prettier --check` + `markdownlint-cli2` + `zensical build`) as a
smoke test (no deploy). The actual deploy runs only when release-please cuts a release on
`main`, via the `publish-docs` job in `.github/workflows/release-main.yaml`.

## Authoring conventions

- **Lead with a one-line summary** of what the page covers.
- **Wrap at 120 chars.** The `.markdownlint-cli2.yaml` config enforces it (and relaxes a few
  rules that conflict with the pymdownx extensions we use).
- **Tables for key/value mappings** (config options, aliases, scripts).
- **Reference the source-of-truth file path** instead of duplicating large config blocks —
  readers can click through.
- **Icons** must exist in Zensical's bundled set (`.venv/lib/python*/site-packages/zensical/templates/.icons/`).
  Stick to `lucide/*` unless you've verified a `simple/*` or `material/*` icon ships.

## Keep dependabot in sync

When a new package ecosystem is introduced to the repo (anything with a manifest like
`package.json`, `pyproject.toml`, `Cargo.toml`, `Gemfile`, `go.mod`, `composer.json`, a
new GitHub Actions workflow file, etc.), add a matching `package-ecosystem` block to
[`.github/dependabot.yaml`](.github/dependabot.yaml) **in the same PR**. Mirror the
existing entries: `directory: /`, `schedule.interval: daily`, `cooldown.default-days: 7`,
and an `all-minor-and-patch` group that batches minor + patch updates. Without this,
the new ecosystem gets no automated security or version updates.

## Don't touch

- `CHANGELOG.md` — owned by release-please.
- `.claude/plans/` — runtime artifacts from Claude Code's plan mode.
- `site/`, `.venv/` — build/runtime artifacts (gitignored).

## Useful entry points

- Repo overview & layout: [`docs/index.md`](docs/index.md)
- How the docs are built: [`docs/tooling/contributing.md`](docs/tooling/contributing.md)
- Cyberdream palette source-of-truth: [`docs/assets/extras.css`](docs/assets/extras.css)
