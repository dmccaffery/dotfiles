---
icon: lucide/keyboard
---

# Auto-commands, keymaps, options

The three files under `lua/config/` shape day-to-day editing.

## `options.lua`

Loads before plugins (`event: BeforeLazy`). Highlights:

```lua
vim.loader.enable()                       -- byte-code cache
g.loaded_python3_provider = 0             -- no python3 provider
g.loaded_perl_provider   = 0
g.loaded_ruby_provider   = 0
g.loaded_node_provider   = 0

g.lazyvim_cmp = "auto"                    -- LazyVim default cmp (blink)
g.lazyvim_picker = "snacks"               -- Snacks picker over Telescope
g.autoformat = true                       -- format on save
g.editorconfig = true                     -- respect .editorconfig

g.gui_font_face = "Iosevka NF"
o.guifont = string.format("%s:h%s", g.gui_font_face, g.gui_font_size)

o.guicursor    = "n-v-c-sm:block,i-ci-ve:ver25,r-cr-o:hor20"
o.splitright   = true
o.splitbelow   = true
o.colorcolumn  = "80,120"                 -- rulers
o.list         = true                     -- show whitespace
o.spell        = true                     -- spell check on
o.showtabline  = 0                        -- no native bufferline
o.laststatus   = 3                        -- global statusline
```

The whitespace listchars use middle-dot for spaces and `>—` for tabs:

```lua
o.listchars:append({
  tab = ">—", multispace = "·",
  extends = ">", precedes = "<", lead = "·", trail = "·",
})
```

LSP logging is set to `off` because the log grows unboundedly when on.

## `keymaps.lua`

Loaded on `VeryLazy`. Notable additions on top of LazyVim defaults:

| Mode | Mapping                    | Action                                            |
| ---- | -------------------------- | ------------------------------------------------- |
| n    | ++tab++                    | Next buffer (uses bufferline if loaded).          |
| n    | ++shift+tab++              | Previous buffer.                                  |
| n    | ++less++                   | Deindent (`<<`).                                  |
| n    | ++greater++                | Indent (`>>`).                                    |
| n    | ++alt+s++                  | Save _without_ formatting (`:noautocmd w`).       |
| n    | ++plus++                   | Increment number (`<C-a>`).                       |
| n    | ++minus++                  | Decrement number (`<C-x>`).                       |
| n    | ++u++ (capital U)          | Redo (`<C-r>`).                                   |
| n    | ++ctrl+c++                 | Yank entire buffer.                               |
| n    | ++ctrl+e++                 | Select all text.                                  |
| n    | ++leader++ ++exclamation++ | `zg` — add word under cursor to spell dictionary. |
| n    | ++leader++ ++at++          | `zug` — remove word from spell dictionary.        |

Mappings are wrapped in a small `map()` helper that skips registration if a lazy.nvim handler
already owns the key, avoiding accidental overrides.

## `autocmds.lua`

Loaded on `VeryLazy`. Adds:

### Auto-close ephemeral UI

When buffers of these filetypes lose focus, their windows close:

> `lazy`, `mason`, `lspinfo`, `toggleterm`, `notify`

### Leader disabled in those filetypes

In the same set plus `floaterm`, both `<leader>` and `<localleader>` are bound to `<nop>` so
filetype-specific keybindings don't conflict with global leader mappings.

### Quickfix not buflisted

The quickfix window doesn't show in `:bnext` / bufferline rotation.

### Terminal sanity

In `:terminal` buffers: `listchars` cleared, no line numbers, no spell-check.

### Trailing whitespace stripped on save

```lua
autocmd("BufWritePre", {
  pattern = { "*" },
  command = [[%s/\s\+$//e]],
})
```

### Don't continue comments

`formatoptions-=cro` is applied on every `BufEnter` (both global and buffer-local) so newlines
after a comment line never auto-insert a comment leader.

### Relative number toggle

`relativenumber` is on in normal mode and off in insert mode, focus-lost, and command-line —
implemented as a pair of autocmds in the `numbertoggle` group.
