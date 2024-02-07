return {
  {
    "akinsho/bufferline.nvim",
    opts = {
      options = {
        always_show_bufferline = true,
        enforce_regular_tabs = true,
        tab_size = 20,
        offsets = {
          {
            filetype = "neo-tree",
            text = "EXPLORER",
            highlight = "Directory",
            text_align = "left",
          },
        },
      },
    },
  },
  {
    "nvim-lualine/lualine.nvim",
    opts = {
      options = {
        theme = "dracula",
        section_separators = "",
      },
    },
  },
}
