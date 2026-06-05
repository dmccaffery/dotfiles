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

### Clipboard over SSH (OSC 52)

When `$SSH_TTY`/`$SSH_CONNECTION` is set, `options.lua` pins Neovim's clipboard provider to its bundled
[OSC 52](https://github.com/neovim/neovim/blob/master/runtime/lua/vim/ui/clipboard/osc52.lua) implementation so a
yank reaches the **local** terminal's clipboard instead of the remote host's:

```lua
if vim.env.SSH_TTY or vim.env.SSH_CONNECTION then
  local osc52 = require("vim.ui.clipboard.osc52")
  local paste = function()
    return { vim.split(vim.fn.getreg('"'), "\n"), vim.fn.getregtype('"') }
  end
  o.clipboard = "unnamedplus"                 -- override LazyVim's SSH blanking so plain `y` copies
  g.clipboard = {
    name = "OSC 52",
    copy = { ["+"] = osc52.copy("+"), ["*"] = osc52.copy("*") },
    paste = { ["+"] = paste, ["*"] = paste }, -- read from the register, never the terminal
  }
end
```

This is needed because LazyVim's default (`opt.clipboard = vim.env.SSH_CONNECTION and "" or "unnamedplus"`) relies on
Neovim auto-enabling OSC 52, and that branch in `provider/clipboard.vim` is only reached when no other provider matches.
On a **macOS** remote `has('mac')` selects `pbcopy` first, so the auto path never fires — the explicit `g.clipboard`
opt-in bypasses it. `paste` reads the unnamed register rather than issuing an OSC 52 read, which otherwise blocks for up
to 10s waiting on a terminal response over SSH + tmux. Use your terminal's own paste (++cmd+v++) to pull text copied
**outside** Neovim into a buffer.

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
