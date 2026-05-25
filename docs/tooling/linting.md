---
icon: lucide/check-check
---

# Linting & formatting

Three tools enforce style across the repository.

## EditorConfig

`.editorconfig` ships baseline indentation, line-ending, and trailing-whitespace rules. All
modern editors (including LazyVim — `g.editorconfig = true`) pick it up automatically.

## Prettier

`.prettierrc.yaml` configures Prettier; `.prettierignore` lists what to skip.

```text title=".prettierignore"
CHANGELOG.md
docs/**/*.md
```

- **`CHANGELOG.md`** — owned by [release-please](release-please.md).
- **`docs/**/\*.md`\*\* — handled by markdownlint instead (prettier's wrap rules fight with
  markdownlint's 120-char line length).

## markdownlint-cli2

```yaml title=".markdownlint-cli2.yaml"
config:
    default: true
    line-length:
        line_length: 120
        heading_line_length: 120
        code_block_line_length: 120
    list-marker-space: false
ignores:
    - CHANGELOG.md
```

- **120-char wrap** across body, headings, and fenced code blocks.
- **`list-marker-space: false`** — turns off the rule requiring exactly one space after `-`
  list markers; lets you align bullets visually when useful.
- **`CHANGELOG.md` is ignored** — release-please writes it.

## Run them locally

```sh
# Prettier (skips CHANGELOG and docs)
npx prettier --check .

# markdownlint (covers docs and any other markdown)
npx markdownlint-cli2 "**/*.md"
```

Both are typically installed on demand via `npx`; they're not in the Brewfile so the per-repo
versions can drift independently.
