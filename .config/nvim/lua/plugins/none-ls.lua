return {
  {
    "nvimtools/none-ls.nvim",
    optional = true,
    opts = function(_, opts)
      local null_ls = require("null-ls")
      opts.sources = vim.list_extend(opts.sources or {}, {
        null_ls.builtins.diagnostics.trivy,
        null_ls.builtins.diagnostics.actionlint,
        null_ls.builtins.diagnostics.buf,
      })
    end,
  },
}
