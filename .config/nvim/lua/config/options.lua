-- event: BeforeLazy
-- defaults: https://github.com/LazyVim/LazyVim/blob/main/lua/lazyvim/config/options.lua

local opt = vim.opt

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
