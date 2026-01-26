return {
  {
    "neovim/nvim-lspconfig",
    opts = {
      servers = {
        gopls = {
          settings = {
            gopls = {
              directoryFilters = { "-.git", "-.vscode", "-.idea", "-.vscode-test", "-node_modules", "-**/.terraform" },
            },
          },
        },
      },
    },
  },
  {
    "stevearc/conform.nvim",
    optional = true,
    opts = {
      formatters_by_ft = {
        go = { "gofumpt" },
      },
    },
  },
}
