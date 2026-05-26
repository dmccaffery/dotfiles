---
icon: lucide/git-merge
---

# Release Please

[Release Please](https://github.com/googleapis/release-please) automates versioning and the
`CHANGELOG.md`. It opens a "Release PR" against `main` that bumps the version and the
changelog; merging that PR publishes a release.

## Configuration

```json title="release-please-config.json"
{
    "$schema": "https://raw.githubusercontent.com/googleapis/release-please/main/schemas/config.json",
    "release-type": "simple",
    "include-component-in-tag": true,
    "packages": { ".": {} }
}
```

- **`release-type: simple`** — no language-specific version bumping (this is dotfiles, not a
  package). Versions come from conventional-commit messages.
- **`include-component-in-tag`** — tag format is `<component>-v<version>`.
- **`packages: { ".": {} }`** — single package at the repo root.

The companion `.release-please-manifest.json` tracks the current version.

## Workflow

`.github/workflows/release-main.yaml` runs on push to `main`:

```yaml
on:
    push:
        branches: [main]
    workflow_dispatch:

permissions:
    contents: write
    pull-requests: write

jobs:
    release:
        runs-on: ubuntu-latest
        steps:
            - uses: googleapis/release-please-action@…
              id: release
              with:
                  token: ${{ secrets.GITHUB_TOKEN }}

            - uses: actions/checkout@…
              if: ${{ steps.release.outputs.release_created }}

            - name: add-tags
              if: ${{ steps.release.outputs.release_created }}
              run: |
                  # delete & re-push v<major> and v<major>.<minor> tags so they always
                  # point at the latest matching release
```

When a release is created, the workflow also moves the floating `v<major>` and
`v<major>.<minor>` tags forward — useful for consumers that pin to a major or minor branch.

The same workflow then runs a `publish-docs` job (gated on
`needs.release.outputs.release_created`) that builds the site with `zensical build --clean`
and ships it to GitHub Pages in a single job. That is why `release-main.yaml` carries the
`pages: write` / `id-token: write` permissions and the `pages` concurrency group alongside
the release-please permissions. PR-time smoke-test builds (plus shellcheck and the
markdown linters) live in `.github/workflows/pull-request.yaml`.

## CHANGELOG ownership

`CHANGELOG.md` is owned by release-please. Don't hand-edit it; commit conventional commit
messages (`feat:`, `fix:`, `chore:`, etc.) and let the action regenerate the changelog when
it next runs.

`.markdownlint-cli2.yaml` and `.prettierignore` already exclude `CHANGELOG.md` from local
linting for this reason.
