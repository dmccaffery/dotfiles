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
# opens a new Ghostty window to record in, walks you through the beats, then
# prompts before replacing docs/assets/demo.cast
./hack/record-demo.sh

# override the default grid (105x40 ~= 800x800 px in Ghostty / Iosevka NF 15)
WINDOW_SIZE=120x45 ./hack/record-demo.sh
```

The recording runs in a **new Ghostty window** (opened with `open -na Ghostty.app`, the only
supported way to launch Ghostty from the CLI on macOS), sized to the recording grid and pinned to
the dark `cyberdream` theme with `--theme`. The pin matters: a window spawned via `open -na … -e sh`
resolves its appearance before AppKit finishes initialising and can fall back to the light branch of
`theme = dark:cyberdream,light:cyberdream-light`, which bakes a white background into the cast header
even when macOS is in dark mode. Override the pinned theme with `GHOSTTY_THEME=… ./hack/record-demo.sh`
if you ever need a different palette. The window you ran the script from becomes a **teleprompter**: it walks the
demo beats one at a time, waiting for Enter between each. Command steps are copied to the clipboard
so you can paste them into the recording window with ++cmd+v++; the quoted steps are directions to
perform there (navigate the tmux windows, drive the agent, open lazygit) and are not copied. Reproduce
the structure, not the literal text — the Claude-agent output naturally differs. When you exit the
recorded shell, control returns to the teleprompter window to publish the cast (`open` and `pbcopy`
make this step macOS-only). The grid is forced with asciinema's `--window-size`, so the geometry
holds regardless of how you resize the window.

When you publish a fresh take, the script also regenerates the README/docs poster
([`docs/assets/images/demo-poster.png`](https://github.com/dmccaffery/dotfiles/blob/main/docs/assets/images/demo-poster.png))
from the new cast with [`agg`](https://docs.asciinema.org/manual/agg/): it renders the whole cast
at the player's idle cap, grabs the poster frame, and flattens it onto the cyberdream background.
`agg` renders with the same Iosevka Nerd Font as the terminal (via `--font-dir` / `--font-family`,
per the [agg Nerd Fonts docs](https://docs.asciinema.org/manual/agg/usage/#nerd-fonts)) so the
powerline / oh-my-posh glyphs render instead of tofu. This step needs `agg`, `ffmpeg`, and `magick`
on `PATH`; if any is missing it warns and skips, leaving the poster for you to update by hand.

The poster frame is **whatever was on screen as you crossed from step 8 to step 9** in the
teleprompter (the post-commit view). The script records that wall-clock instant, converts it to the
player's idle-collapsed `npt` timeline using the cast's start `timestamp` and per-event intervals
(each capped at the 2s idle limit, exactly as `agg` and the player collapse idle gaps), then both
seeks `agg` to that frame and rewrites `poster: "npt:…"` in `asciinema-player-init.js` so the live
poster and the static PNG agree. Converting through the cast — rather than using raw elapsed
seconds — keeps the frame correct even when the take has long idle pauses (e.g. while the agent
works). Deriving the timestamp needs `python3`; without it the script keeps the existing
`POSTER_NPT` fallback.

### Nerd Font web fonts

asciinema-player renders the recording as HTML, so the browser needs a Nerd Font web font or the
powerline / oh-my-posh glyphs fall back to tofu boxes. The full Iosevka Nerd Font is ~13 MB per
face, so [`hack/build-demo-fonts.sh`](https://github.com/dmccaffery/dotfiles/blob/main/hack/build-demo-fonts.sh)
subsets each face to the glyphs the cast renders (plus common box-drawing, block, dingbat and
powerline ranges) and writes ~50 KB `woff2` files to `docs/assets/fonts/`. `extras.css` then
`@font-face`s them, and `asciinema-player-init.js` points the player at the family with its
`terminalFontFamily` option (the player measures its own glyph metrics, so a CSS `font-family`
override alone leaves it on the default font) and waits for the font to load before mounting so the
grid is measured against the right metrics.

```sh
# re-run after a re-record only if the new take uses glyphs outside the kept ranges
./hack/build-demo-fonts.sh
```

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
