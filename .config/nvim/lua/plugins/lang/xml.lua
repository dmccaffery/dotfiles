return {
  recommended = {
    ft = { "xml", "plist" },
  },
  {
    "nvim-treesitter/nvim-treesitter",
    opts = { ensure_installed = { "xml" } },
  },
  {
    "neovim/nvim-lspconfig",
    opts = {
      servers = {
        lemminx = {},
      },
    },
  },
}
