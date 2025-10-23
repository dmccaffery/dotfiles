-- event: VeryLazy
-- defaults: https://github.com/LazyVim/LazyVim/blob/main/lua/lazyvim/config/autocmds.lua

local autocmd = vim.api.nvim_create_autocmd

autocmd("FileType", {
  pattern = "qf",
  callback = function()
    vim.opt_local.buflisted = false
  end,
})

-- remove trailing whitespace
autocmd("BufWritePre", {
  pattern = { "*" },
  command = [[%s/\s\+$//e]],
})

autocmd("VimEnter", {
  callback = function()
    require("persistence").load()
  end,
  nested = true,
})
