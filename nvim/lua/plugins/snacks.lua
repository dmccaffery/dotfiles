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
    {
      "<leader>@",
      desc = "Session Snacks (tmux)",
      function()
        local function get_sessions()
          local raw = vim.fn.system("tmux list-sessions -F '#{session_name}'")
          local sessions = {}

          for session in raw:gmatch("[^\r\n]+") do
            table.insert(sessions, { text = session })
          end

          return sessions
        end

        local available = get_sessions()

        Snacks.picker.pick({
          source = "tmux_sessions",
          items = available,
          format = "text",
          layout = {
            preset = "select",
            layout = {
              position = "bottom",
            },
          },
          confirm = function(picker, item)
            picker:close()
            local session = item.text
            if session then
              vim.fn.system(string.format("tmux switch -t '%s'", session))
            end
          end,
        })
      end,
    },
  },
}
