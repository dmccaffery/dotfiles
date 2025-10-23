return {
  {
    "neovim/nvim-lspconfig",
    version = "*",
    opts = {
      codelens = {
        enabled = true,
      },

      -- servers = {
      --   terraformls = {
      --     cmd = {
      --       "terraform-ls",
      --       "serve",
      --       "-logfile",
      --       vim.fs.dirname(require("vim.lsp.log").get_filename()) .. "/terraform-ls.log",
      --     },
      --   },
      -- },
    },
  },
}
