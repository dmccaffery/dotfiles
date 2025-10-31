return {
  {
    "mason-org/mason.nvim",
    opts = { ensure_installed = { "tofu-ls", "trivy" } },
  },
  {
    "neovim/nvim-lspconfig",
    opts = {
      servers = {
        tofu_ls = {},
        tflint = {},
      },
    },
  },
  {
    "mfussenegger/nvim-lint",
    optional = true,
    opts = {
      linters_by_ft = {
        tofu = { "tofu", "trivy" },
        opentofu = { "tofu", "trivy" },
        ["opentofu-vars"] = { "tofu", "trivy" },
      },
    },
  },
  {
    "stevearc/conform.nvim",
    optional = true,
    opts = {
      formatters_by_ft = {
        tofu = { "tofu_fmt" },
        opentofu = { "tofu_fmt" },
        ["opentofu-vars"] = { "tofu_fmt" },
      },
    },
  },
}
