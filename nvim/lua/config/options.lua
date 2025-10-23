-- event: BeforeLazy
-- defaults: https://github.com/LazyVim/LazyVim/blob/main/lua/lazyvim/config/options.lua

local opt = vim.opt
local space = "·"

-- cursor
opt.guicursor = "n-v-c-sm:block,i-ci-ve:ver25,r-cr-o:hor20"

-- tabs / indentation
opt.autoindent = true

-- line wrapping
opt.whichwrap:append("<>[]hl")

-- backspace
opt.backspace = "indent,eol,start"

-- split windows
opt.splitright = true
opt.splitbelow = true

-- add a ruler
opt.colorcolumn = "80,120"

-- disable lsp logs -- this will grow infinitely so only enable it if you need it
vim.lsp.log_level = vim.lsp.log_levels.OFF

-- show whitespace characters
opt.list = true
opt.listchars:append({
  tab = ">—",
  multispace = space,
  extends = ">",
  precedes = "<",
  lead = space,
  trail = space,
})

-- smooth scrolling
opt.smoothscroll = true
