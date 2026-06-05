---
icon: lucide/grid
---

# Per-tool theming

Cyberdream is applied in every tool that exposes a theme system. The configurations live next
to each tool's other config.

| Tool        | Config path                                            | Notes                                                                    |
| ----------- | ------------------------------------------------------ | ------------------------------------------------------------------------ |
| NeoVim      | `.config/nvim/lua/plugins/cyberdream.lua`              | `cyberdream.nvim`, loaded eagerly with priority 1000.                    |
| oh-my-posh  | `.config/oh-my-posh/prompt.yaml` (`palette:` block)    | Plus `claude.yaml` for the Claude Code status line.                      |
| Ghostty     | `.config/ghostty/themes/{cyberdream,cyberdream-light}` | Auto-switches with system appearance.                                    |
| tmux        | `.config/tmux/conf/theme.conf`                         | Catppuccin engine with cyberdream flavor (auto-fetched on first launch). |
| btop        | `.config/btop/themes/cyberdream.theme`                 | —                                                                        |
| lazygit     | `.config/lazygit/config.yml`                           | Colors + delta pager.                                                    |
| yazi        | `.config/yazi/theme.toml`                              | Includes bat syntax theme.                                               |
| k9s         | `.config/k9s/skins/cyberdream.yaml`                    | Selected as the active skin.                                             |
| opencode    | `.config/opencode/themes/cyberdream.json`              | Terminal-default background.                                             |
| vivid       | `.config/vivid/themes/cyberdream.yaml`                 | Generates `LS_COLORS`; `VIVID_THEME` env points here.                    |
| lsd         | `.config/lsd/config.yaml`                              | `color: { theme: custom }` driven by `LS_COLORS`.                        |
| Claude Code | `.claude/themes/cyberdream.json`                       | Selected via `"theme": "custom:cyberdream"`.                             |
| Codex       | `.config/codex/config.toml` (`[tui].theme`)            | No theme file; `theme = "ansi"` inherits cyberdream from the terminal.   |
| Docs site   | `docs/assets/extras.css`                               | CSS overrides on Zensical's `slate`/`default` schemes.                   |

## Pattern

The recurring approach is:

1. **Use the upstream's theme system, not a hex monkey-patch.** Most tools have a way to
   declare a theme by name or import a file — that's the path of least churn when the tool
   updates.
2. **Centralize palette in the canonical place.** When a theme system supports it (oh-my-posh,
   yazi, Claude Code), define the palette once in a `palette:` / `colors:` block and reference
   colours symbolically elsewhere.
3. **Match the system appearance when possible.** Ghostty does this natively
   (`dark:cyberdream,light:cyberdream-light`). For tools that don't, the dark variant wins.
