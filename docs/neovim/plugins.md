---
icon: lucide/plug
---

# Custom plugins

Beyond LazyVim and its extras, `.config/nvim/lua/plugins/` holds project-specific plugin
specs and overrides.

## Theme

### `cyberdream.lua`

```lua
return {
  { "scottmckendry/cyberdream.nvim", lazy = false, priority = 1000 },
  { "nvim-lualine/lualine.nvim", opts = { theme = "auto" } },
}
```

Loads [cyberdream.nvim](https://github.com/scottmckendry/cyberdream.nvim) eagerly with the
highest priority so the colourscheme is in place before any other UI renders.

### `lua-line.lua`

Drops the `lualine_z` section so the status line ends at branch/diagnostics instead of the
default position/percent block.

### `header.lua`

Custom Snacks dashboard header — ASCII art that spells "Deavon's Terminal". See
[Customize → NeoVim dashboard header](../getting-started/customize.md#neovim-dashboard-header)
for how to regenerate the art with `figlet`.

## UX

### `snacks.lua`

A big Snacks configuration covering:

- **Explorer** — hidden files visible, gitignored files visible.
- **Sessions picker** — lists tmux sessions, previews their cwd with `lsd --tree`, and
  switches to the picked one with `tmux switch-client`. Bind: ++leader++ ++f++ ++s++.
- **Snippets picker** — fuzzy-find LuaSnip snippets with live preview and expansion. Bind:
  ++leader++ ++f++ ++x++.
- **Image previews** — enabled (`image = {}`), so `:Snacks.picker.files` previews images
  inline in the supported terminals.
- **Files default to cwd** — ++leader++ ++space++ uses `root = false` so it searches the
  current directory rather than the project root.

### `auto-save.lua`

```lua
return {
  "Pocco81/auto-save.nvim",
  opts = {
    trigger_events = { "InsertLeave" },
    debounce_delay = 500,
  },
  keys = { { "<leader>uv", "<cmd>ASToggle<CR>", desc = "Toggle Auto-save" } },
}
```

Saves on `InsertLeave` with a 500 ms debounce. Toggle with ++leader++ ++u++ ++v++.

### `blink.lua`

Configures `blink.cmp` so the completion menu _doesn't_ preselect the first item — pressing
Enter inserts a literal newline instead of accepting an unwanted completion.

### `grug-far.lua`

Search-and-replace TUI scoped to the current buffer. Bind: ++leader++ ++s++ ++f++.

### `vim-tmux-navigator.lua`

Maps ++ctrl+h++ / ++ctrl+j++ / ++ctrl+k++ / ++ctrl+l++ for seamless navigation between Vim
splits and tmux panes. Paired with the matching tmux plugin (see [tmux](../terminal/tmux.md)).

### `ghostty.lua`

Loads the [vimfiles bundled with Ghostty](https://ghostty.org/docs/help/editors) so the
`ghostty` config filetype gets syntax highlighting:

```lua
return { "ghostty", dir = "/Applications/Ghostty.app/Contents/Resources/vim/vimfiles/" }
```

### `sql.lua` { #sql }

Wires `sqlfluff` into the LazyVim formatting and linting machinery (`format -` and
`lint --format=json`).

## Language plugins

`lua/plugins/lang/` adds support for languages not covered by LazyVim extras:

| File           | Adds                                                        |
| -------------- | ----------------------------------------------------------- |
| `go.lua`       | Extra Go tooling / debug configuration.                     |
| `markdown.lua` | Markdown-specific plugin tweaks (treesitter, render, etc.). |
| `tofu.lua`     | OpenTofu (Terraform fork) support.                          |
| `xml.lua`      | XML LSP / formatting.                                       |
| `jinja.lua`    | Jinja2 templating filetype + LSP.                           |
| `protobuf.lua` | Protocol Buffers LSP / formatting.                          |
