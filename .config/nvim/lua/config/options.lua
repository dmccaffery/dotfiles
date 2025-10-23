-- event: BeforeLazy
-- defaults: https://github.com/LazyVim/LazyVim/blob/main/lua/lazyvim/config/options.lua

local opt = vim.opt

-- tabs / indentation
opt.autoindent = true

-- line wrapping
opt.whichwrap:append("<>[]hl")

-- appearance
opt.background = "dark"

-- backspace
opt.backspace = "indent,eol,start"

-- split windows
opt.splitright = true
opt.splitbelow = true

-- special keys
opt.iskeyword:append("-")

-- add a ruler
opt.colorcolumn = "80,120"

-- disable providers
for _, provider in ipairs({ "node", "python3" }) do
  vim.g["loaded_" .. provider .. "_provider"] = 0
end

-- disable lsp logs -- this will grow infinitely so only enable it if you need it
vim.lsp.log_level = vim.lsp.log_levels.OFF
