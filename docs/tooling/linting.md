---
icon: lucide/check-check
---

# Linting & formatting

Four tools enforce style across the repository — one editor baseline, one formatter, and two
linters. The npm-based tools share a single `package.json` so both config and pinned
versions live in one place.

## EditorConfig

`.editorconfig` ships baseline indentation, line-ending, and trailing-whitespace rules. All
modern editors (including LazyVim — `g.editorconfig = true`) pick it up automatically.

## Prettier

Prettier is pinned in [`package.json`](https://github.com/dmccaffery/dotfiles/blob/main/package.json)
under `devDependencies`; its config lives in the same file under the `prettier` key, and
`.prettierignore` lists what to skip.

```json title="package.json (excerpt)"
"prettier": {
  "printWidth": 120
}
```

```text title=".prettierignore"
CHANGELOG.md
site/
```

- **120-char wrap** for everything prettier touches (matches markdownlint's line length).
- **`CHANGELOG.md`** is owned by [release-please](release-please.md).
- **`site/`** is the zensical build output.

`make fmt` runs `npm install` (idempotent) then `npx prettier --write 'docs/**/*.md'`. CI
runs the same with `--check` so unformatted docs fail the PR rather than getting silently
rewritten.

## markdownlint-cli2

markdownlint-cli2 is also pinned in `package.json`, and its config sits next to prettier's
under the `markdownlint-cli2` key. There is no longer a separate
`.markdownlint-cli2.yaml`.

```json title="package.json (excerpt)"
"markdownlint-cli2": {
  "config": {
    "default": true,
    "line-length": {
      "line_length": 120,
      "heading_line_length": 120,
      "code_block_line_length": 120,
      "tables": false
    },
    "list-marker-space": false,
    "table-column-style": { "style": "aligned" },
    "code-block-style": { "style": "fenced" },
    "code-fence-style": { "style": "backtick" },
    "no-inline-html": {
      "allowed_elements": ["span", "div", "br", "p"]
    },
    "no-space-in-code": false,
    "link-fragments": false
  },
  "ignores": ["CHANGELOG.md"]
}
```

- **120-char wrap** across body, headings, and fenced code blocks; tables are exempt because
  prettier auto-pads them past 120.
- **`list-marker-space: false`** — turns off the rule requiring exactly one space after `-`
  list markers; lets you align bullets visually when useful.
- **`no-inline-html` allow list** — span/div/br/p are needed by zensical's templates.
- **`no-space-in-code: false`** — inline code with intentional whitespace (e.g. `` `D ` ``)
  is needed for tmux window names with nerd-font prefixes.
- **`link-fragments: false`** — pymdownx `attr_list` IDs (`## Heading { #anchor }`) render
  in zensical but markdownlint can't resolve them.
- **`CHANGELOG.md` is ignored** — release-please writes it.

## shellcheck

[shellcheck](https://www.shellcheck.net) lints every shell script in the repo —
`install.sh`, `restore.sh`, `backup.sh`, `setup/**/*.sh`, every executable under
`.local/share/scripts/`, the git template hooks, and `.ssh/rc`. CI and `make lint` both run
it with `--severity=warning --external-sources` so style/info findings are skipped (most
notably SC1091 for dynamic `${SETUP_DIR}/printing.sh` sources) while real warnings and
errors still fail.

## Run them locally

```sh
make fmt    # prettier --write
make lint   # depends on fmt → shellcheck + markdownlint-cli2
make docs-build   # depends on lint → uv sync + zensical build --clean
```

`make fmt` runs `npm install` first, so the first invocation populates `node_modules/` from
`package-lock.json`. Re-runs are fast because `npm install` is a no-op when the lock is in
sync.

Need to run a single tool by hand?

```sh
npx prettier --check .            # prettier without writing
npx markdownlint-cli2 '**/*.md'   # markdownlint over the whole tree
shellcheck install.sh             # one-off shellcheck
```
