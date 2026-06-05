---
icon: lucide/palette
---

# Theme

Codex has **no custom hex palette** like Claude Code's [`cyberdream.json`](../claude/theme.md) or OpenCode's
[`cyberdream.json`](../opencode/theme.md) — so there is no Codex theme _file_ in this repo. Cyberdream still applies,
through two channels that need no fork:

1. **The TUI chrome** (borders, labels, spinners, diff colours) renders with the terminal's 16 ANSI colours, which
   [Ghostty already paints cyberdream](../theme/per-tool.md). Codex inherits the palette for free.
2. **Code-block syntax highlighting** is the one thing `[tui].theme` controls — a [syntect](https://github.com/trishume/syntect)
   theme name (the `bat` set: `catppuccin-mocha`, `dracula`, `gruvbox-dark`, …). Setting it to **`ansi`** ties
   highlighting to those same terminal ANSI colours, so the whole TUI stays on one palette instead of pulling in a
   fixed external one:

```toml title=".config/codex/config.toml"
[tui]
theme = "ansi"
```

That makes `theme = "ansi"` the closest mirror of Claude Code's `"theme": "custom:cyberdream"`: both keep the agent
on the cyberdream palette, Claude Code via an explicit hex theme and Codex via the terminal's ANSI colours. Update
Ghostty's [`cyberdream`](../theme/per-tool.md) palette and Codex follows automatically.

## What Codex can't theme

The trade-off is fidelity: an `ansi` highlight only has 16 colours to work with, so syntax highlighting is coarser
than Claude Code's or OpenCode's full hex themes. If you prefer richer code blocks over strict palette unity, swap
`theme` for the nearest fixed dark theme (`catppuccin-mocha` is the closest in spirit) — at the cost of code blocks
no longer tracking the terminal palette.
