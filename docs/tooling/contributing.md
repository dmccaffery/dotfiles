---
icon: lucide/git-pull-request
---

# Contributing

This is a personal dotfiles repository, but the docs site is structured so that anyone can
fork, edit, and serve locally.

## Dev loop

The docs site is a [Zensical](https://zensical.org) project at the repo root. The Python
runtime is managed by `uv`.

```sh
# install zensical + deps into .venv, then serve at http://localhost:8000
make docs-serve

# one-shot production build (run by CI), output: ./site
make docs-build

# upgrade zensical and friends to the latest matching versions
make docs-upgrade
```

`make docs-serve` wraps `uv sync` + `uv run zensical serve`; `make docs-build` wraps
`uv sync` + `uv run zensical build --clean`; `make docs-upgrade` runs `uv sync --upgrade`
to bump `uv.lock`. Drop down to raw `uv` commands when you need flags the targets don't
pass through.

The first `uv sync` pins `zensical>=0.0.43` (latest at time of writing — still alpha) and
writes `uv.lock`. CI uses `uv sync --frozen` so the lock file is authoritative.

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

## Deploy

Published GitHub Releases build via `.github/workflows/docs.yaml` and deploy to GitHub Pages.
Pull requests run the same build for validation but do not deploy. A manual
`workflow_dispatch` run is available as an escape hatch for off-cycle redeploys.

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
