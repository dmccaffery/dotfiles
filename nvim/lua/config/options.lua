-- event: BeforeLazy
-- defaults: https://github.com/LazyVim/LazyVim/blob/main/lua/lazyvim/config/options.lua

local g = vim.g
local o = vim.opt
local space = "·"

-- optimise startup
vim.loader.enable()

-- defaults
g.lualine_info_extras = false
g.lazyvim_cmp = "blink"
g.lazyvim_picker = "snacks"

-- formatting
g.autoformat = true

-- editor config
g.editorconfig = true

-- ui
g.gui_font_default_size = 10
g.gui_font_fize = g.gui_font_default_size
g.gui_font_face = "FiraCode Nerd Font"

-- cursor
o.guicursor = "n-v-c-sm:block,i-ci-ve:ver25,r-cr-o:hor20"

-- line wrapping
o.whichwrap:append("<>[]hl")
o.autoindent = true

-- backspace
o.backspace = { "start", "eol", "indent" }
o.breakindent = true

-- split windows
o.splitright = true
o.splitbelow = true

-- add a ruler
o.colorcolumn = "80,120"

-- show whitespace characters
o.list = true
o.listchars:append({
  tab = ">—",
  multispace = space,
  extends = ">",
  precedes = "<",
  lead = space,
  trail = space,
})

-- disable native bufferline
o.showtabline = 0

-- command options
o.laststatus = 3

-- spell checking
o.spell = true

-- backspace
o.backspace = { "start", "eol", "indent" }
o.breakindent = true

-- smooth scrolling
o.smoothscroll = true

-- disable lsp logs -- this will grow infinitely so only enable it if you need it
vim.lsp.log.set_level("off")

-- enable tofu ls
vim.lsp.enable("tofu_ls")
