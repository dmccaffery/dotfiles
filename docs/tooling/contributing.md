---
icon: lucide/git-pull-request
---

# Contributing

This is a personal dotfiles repository, but the docs site is structured so that anyone can
fork, edit, and serve locally.

## Dev loop

The docs site is a [Zensical](https://zensical.org) project at the repo root. The Python
runtime is managed by `uv`; prettier and markdownlint-cli2 are pinned in `package.json`
(via `package-lock.json`). Prettier's config lives there too; markdownlint-cli2's config is
in a dedicated `.markdownlint-cli2.yaml`.

```sh
# install zensical + deps into .venv, then serve at http://localhost:8000
make docs-serve

# auto-format docs with prettier (npm install on first run)
make fmt

# fmt + shellcheck + markdownlint
make lint

# lint + zensical build --clean, output: ./site
make docs-build

# upgrade npm + uv deps to latest matching versions (prompts before running)
make upgrade
```

`make docs-serve` wraps `uv sync` + `uv run zensical serve`. `make fmt` runs
`npm install --silent` (idempotent against `package-lock.json`) then
`npx prettier --write 'docs/**/*.md'`. `make lint` depends on `fmt`, then runs
`shellcheck --severity=warning --external-sources …` over every shell script
(`install.sh`, `restore.sh`, `backup.sh`, `hack/*.sh`, `setup/**/*.sh`,
`.local/share/scripts/*`, the git template hooks, and `.ssh/rc`) followed by
`npx markdownlint-cli2 'docs/**/*.md'`.
`make docs-build` depends on `lint`, then runs `uv sync` + `uv run zensical build --clean`.
`make upgrade` prompts for confirmation (because it bypasses the 7-day dependabot
cooldown), then runs `npm update` to refresh `package-lock.json` and `uv sync --upgrade` to
refresh `uv.lock`. Prefer merging the matching dependabot PR when one is already open.
Drop down to raw `uv` / `npx` / `shellcheck` commands when you need flags the targets don't
pass through.

The first `uv sync` pins `zensical>=0.0.43` (latest at time of writing — still alpha) and
writes `uv.lock`. CI uses `uv sync --frozen` and `npm ci`, so both lock files are
authoritative.

## Editing content

```text
docs/
├── index.md
├── assets/extras.css        # cyberdream palette overrides for Zensical
├── getting-started/
├── terminal/
├── neovim/
├── theme/
├── git/
├── scripts/
├── macos/
├── claude/
└── tooling/
```

- **Markdown source** lives under `docs/`. Add a new section by creating a directory + an
  index page; new pages need to be wired into `nav = [...]` in [`zensical.toml`](https://github.com/dmccaffery/dotfiles/blob/main/zensical.toml).
- **120-char wrap** — `.markdownlint-cli2.yaml` enforces this for headings, body, and code.
- **Admonitions** (`!!! note`, `!!! warning`) are configured.
- **Mermaid diagrams** are configured (use ` ```mermaid `).
- **Content tabs** (`=== "Tab title"`) and **collapsible details** (`??? info "Title"`) are
  available.

## Re-recording the homepage demo

The homepage embeds [`docs/assets/demo.cast`](https://github.com/dmccaffery/dotfiles/blob/main/docs/assets/demo.cast)
through asciinema-player (mounted by
[`docs/assets/asciinema-player-init.js`](https://github.com/dmccaffery/dotfiles/blob/main/docs/assets/asciinema-player-init.js)).
A cast is baked to a fixed grid and can't be reflowed, so changing its size means re-recording.
[`hack/record-demo.sh`](https://github.com/dmccaffery/dotfiles/blob/main/hack/record-demo.sh)
captures a fresh take at a compact, roughly-square geometry and offers to publish it over the
existing asset.

```sh
# records to $TMPDIR, then prompts before replacing docs/assets/demo.cast
./hack/record-demo.sh

# override the default grid (105x40 ~= 800x800 px in Ghostty / Iosevka NF 15)
WINDOW_SIZE=120x45 ./hack/record-demo.sh
```

The grid is forced with asciinema's `--window-size`, so the physical Ghostty window only needs to
be at least that large. Run it inside Ghostty with the cyberdream theme active so the captured
header palette matches. The script prints a replay checklist of the original demo's beats — the
Claude-agent output naturally differs, so reproduce the structure, not the literal text.

!!! warning "Keep the take past 2:00"

    `asciinema-player-init.js` sets `poster: "npt:2:00"` (the lazygit-in-tmux frame). A recording
    shorter than two minutes blanks the poster — keep it longer, or drop the poster timestamp in
    that file. The script bakes idle at 2s (`--idle-time-limit 2`), matching the player's
    `idleTimeLimit: 2`.

## Deploy

The deploy lives in the `publish-docs` job of `.github/workflows/release-main.yaml`: when
release-please cuts a release on `main`, that job builds the site and ships it to GitHub
Pages in a single step. Pull requests run `.github/workflows/pull-request.yaml`, which
shellchecks every shell script and re-runs the docs build (`prettier --check` +
`markdownlint-cli2` + `zensical build --clean`) as a smoke test — it verifies the build but
does not deploy. A manual `workflow_dispatch` run of `release-main` is the only escape
hatch, and it only redeploys when release-please opens or merges a release on that run.

!!! note "One-time Pages enablement"

    The first deploy needs **Settings → Pages → Source = GitHub Actions** flipped on in the
    repository. The workflow can't enable that for you; the manual flip is one click.

## Style notes

- Lead each page with a one-line summary of what the page covers.
- Use tables to summarize key/value mappings (config options, aliases, scripts).
- Link to the source-of-truth file in the repo when describing config — readers can click
  through to see the canonical form.
- Avoid duplicating large blobs of config into the docs verbatim; show the _interesting_
  bits and link to the file for the rest.

## What NOT to do

- Don't edit `CHANGELOG.md` — release-please owns it.
- Don't add `docs/**/*.md` to prettier; markdownlint handles them.
- Don't add a `_navigation.md` or similar — Zensical reads navigation from `zensical.toml`.
