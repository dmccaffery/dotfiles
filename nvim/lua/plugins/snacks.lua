return {
  "folke/snacks.nvim",
  opts = {
    dashboard = { enabled = true },
    picker = {
      sources = {
        explorer = {
          hidden = true,
          ignored = false,
        },
      },
    },
  },
  keys = {
    -- swap the defaults to use cwd instead of root
    { "<leader><space>", LazyVim.pick("files", { root = false }), desc = "Find Files (cwd)" },
    { "<leader>E", "<leader>fe", desc = "Explorer Snacks (cwd)", remap = true },
    { "<leader>e", "<leader>fE", desc = "Explorer Snacks (root dir)", remap = true },
  },
}
