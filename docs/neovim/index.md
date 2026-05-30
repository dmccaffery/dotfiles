---
icon: simple/neovim
---

# NeoVim

NeoVim configuration is built on LazyVim, with selected extras enabled, custom plugin overlays
in `lua/plugins/`, and the `options.lua` / `keymaps.lua` / `autocmds.lua` triplet under
`lua/config/` shaping day-to-day editing.

![NeoVim with LazyVim and the cyberdream colorscheme editing a buffer](../assets/images/neovim.png)

| Page                                                   | Purpose                                                          |
| ------------------------------------------------------ | ---------------------------------------------------------------- |
| [LazyVim](lazyvim.md)                                  | Layout of `.config/nvim/`, the base spec, and how LazyVim loads. |
| [Extras](extras.md)                                    | Enabled LazyVim extras (language servers, formatters, tooling).  |
| [Custom plugins](plugins.md)                           | Project-specific plugin specs and overrides.                     |
| [Auto-commands, keymaps, options](autocmds-keymaps.md) | The `lua/config/` triplet shaping day-to-day editing.            |
