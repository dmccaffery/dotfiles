---
icon: lucide/palette
---

# Theme

OpenCode's TUI wears the same cyberdream skin as the rest of the terminal stack, defined in
`.config/opencode/themes/cyberdream.json` and selected in `.config/opencode/tui.json`. See
[Per-tool theming](../theme/per-tool.md) for how it fits the wider palette.

## Activation

`.config/opencode/tui.json` selects the bundled `cyberdream` theme by name:

```json title=".config/opencode/tui.json"
{
    "$schema": "https://opencode.ai/tui.json",
    "theme": "cyberdream"
}
```

OpenCode resolves the name against `~/.config/opencode/themes/<name>.json`, so dropping `cyberdream.json` in that
directory is all the wiring the custom theme needs.

## Structure

The theme file splits into a `defs` block of named colors and a `theme` block that maps every visual role to one of
those names:

```json title=".config/opencode/themes/cyberdream.json"
{
    "$schema": "https://opencode.ai/theme.json",
    "defs": {
        "bg": "none",
        "fg": "#FFFFFF",
        "blue": "#5EA1FF",
        "green": "#5EFF6C",
        "cyan": "#5EF1FF",
        "red": "#FF6E5E",
        "yellow": "#F1FF5E",
        "magenta": "#FF5EF1",
        "pink": "#FF5EA0",
        "orange": "#FFBD5E",
        "purple": "#BD5EFF"
        // … plus bg_alt, bg_highlight, grey, and diff backgrounds
    },
    "theme": {
        "primary": "cyan",
        "secondary": "purple",
        "accent": "cyan"
        // … plus border, diff, markdown, and syntax roles
    }
}
```

## Transparent background

`defs.bg` is set to `"none"` so OpenCode uses the **terminal** background rather than painting its own main panel
color. That keeps the TUI consistent with the transparent/blurred Ghostty background instead of stamping a solid
fill over it.

## Notable mappings

| Role                              | Color                | Meaning                                        |
| --------------------------------- | -------------------- | ---------------------------------------------- |
| `primary` / `accent`              | `#5EF1FF` cyan       | Primary accents and highlights.                |
| `secondary` / `syntaxType`        | `#BD5EFF` purple     | Secondary accents, type tokens.                |
| `borderActive`                    | `#5EA1FF` blue       | Border of the focused panel.                   |
| `error` / `warning` / `success`   | red / orange / green | Standard status colors.                        |
| `diffAdded` / `diffRemoved`       | green / red          | Added and removed lines in diffs.              |
| `markdownEmph` / `markdownStrong` | magenta / orange     | Emphasis and strong text in rendered markdown. |
| `syntaxKeyword` / `syntaxString`  | orange / green       | Syntax highlighting in code blocks.            |

Config is loaded at startup — quit and restart OpenCode after editing the theme.
