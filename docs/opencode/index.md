---
icon: lucide/bot
---

# OpenCode

[OpenCode](https://opencode.ai) is a coding-agent TUI, configured as a stow-managed package under
`.config/opencode/`. It deliberately shares this machine's guardrails with [Claude Code](../claude/index.md):
the same global agent rules, a permission config that mirrors the Claude Code sandbox, the cyberdream palette, and
the tmux "needs you" status indicator. Switching between the two agents keeps the same boundaries and the same look.

| Page                          | Purpose                                                                   |
| ----------------------------- | ------------------------------------------------------------------------- |
| [Permissions](permissions.md) | `opencode.jsonc` — the permission-rule mirror of the Claude Code sandbox. |
| [Memory](memory.md)           | `AGENTS.md` — symlink to the shared `~/.claude/CLAUDE.md` agent rules.    |
| [Theme](theme.md)             | `cyberdream.json` + `tui.json` — the OpenCode TUI palette.                |
| [Plugins](plugins.md)         | `agent-tmux-status.js` — the tmux status indicator that flags your turn.  |

## Files

| File                                           | Purpose                                                                 |
| ---------------------------------------------- | ----------------------------------------------------------------------- |
| `.config/opencode/opencode.jsonc`              | Global runtime config and permission rules (JSONC, comments allowed).   |
| `.config/opencode/tui.json`                    | TUI-only settings; selects the `cyberdream` theme.                      |
| `.config/opencode/AGENTS.md`                   | Symlink to `~/.claude/CLAUDE.md`; shares one set of global agent rules. |
| `.config/opencode/themes/cyberdream.json`      | Cyberdream palette for the OpenCode TUI.                                |
| `.config/opencode/plugin/agent-tmux-status.js` | Status-indicator plugin; flags the tmux window when OpenCode needs you. |
| `.config/opencode/.gitignore`                  | Keeps local OpenCode package/plugin install artifacts out of the repo.  |

The `.gitignore` excludes the plugin's install artifacts — `node_modules`, `package.json`, `package-lock.json`,
and `bun.lock` — so only the hand-written config (`opencode.jsonc`, `tui.json`, the theme, and the plugin source)
is tracked. The [`@opencode-ai/plugin`](https://opencode.ai/docs/plugins/) dependency is installed locally by
OpenCode, not committed.

## Reloading config

OpenCode loads its config at startup. After editing `opencode.jsonc`, `tui.json`, or the theme, quit and restart
OpenCode for the change to take effect.
