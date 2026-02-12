return {
  "folke/snacks.nvim",
  opts = {
    image = {},
    picker = {
      sources = {
        explorer = {
          hidden = true,
          ignored = false,
        },
        sessions = {
          supports_live = false,
          format = "text",
          layout = {
            layout = {
              box = "horizontal",
              width = 0.8,
              min_width = 120,
              height = 0.8,
              {
                box = "vertical",
                border = true,
                title = " Sessions ",
                { win = "input", height = 1, border = "bottom" },
                { win = "list", border = "none" },
              },
              { win = "preview", title = " Directory ", border = true, width = 0.5 },
            },
          },
          finder = function(_, ctx)
            local active = vim.fn.systemlist("tmux display-message -p '#S' 2>/dev/null")[1] or ""
            local raw = vim.fn.systemlist("tmux list-sessions -F '#{session_name} #{session_path}' 2>/dev/null")
            local sessions = {}

            for _, line in ipairs(raw) do
              local name, path = line:match("([^ ]+) (.+)")
              if name and path and name ~= active then
                table.insert(sessions, {
                  name = name,
                  data = { path = path },
                  text = string.format("%s (%s)", name, path),
                })
              end
            end

            local align_1 = 0
            for _, session in pairs(sessions) do
              align_1 = math.max(align_1, #session.text)
            end
            ctx.picker.align_1 = align_1

            return sessions
          end,
          confirm = function(picker, item)
            picker:close()
            local session = item.text
            if session then
              vim.fn.system(string.format("tmux switch-client -t '%s'", session.name))
            end
          end,
          preview = function(ctx)
            ctx.preview:reset()
            if not ctx.item then
              ctx.preview:set_title("No selection")
              return
            end

            local command = "lsd --tree {path}"
            command = command:gsub("{path}", ctx.item.data.path)

            local content = vim.fn.systemlist(command)
            ctx.preview:set_lines(content)
          end,
        },
        snippets = {
          supports_live = false,
          preview = "preview",
          format = function(item, picker)
            local name = Snacks.picker.util.align(item.name, picker.align_1 + 5)
            return {
              { name, item.ft == "" and "Conceal" or "DiagnosticWarn" },
              { item.description },
            }
          end,
          finder = function(_, ctx)
            local snippets = {}
            for _, snip in ipairs(require("luasnip").get_snippets().all) do
              snip.ft = ""
              table.insert(snippets, snip)
            end
            for _, snip in ipairs(require("luasnip").get_snippets(vim.bo.ft)) do
              snip.ft = vim.bo.ft
              table.insert(snippets, snip)
            end
            local align_1 = 0
            for _, snip in pairs(snippets) do
              align_1 = math.max(align_1, #snip.name)
            end
            ctx.picker.align_1 = align_1
            local items = {}
            for _, snip in pairs(snippets) do
              local docstring = snip:get_docstring()
              if type(docstring) == "table" then
                docstring = table.concat(docstring)
              end
              local name = snip.name
              local description = table.concat(snip.description)
              description = name == description and "" or description
              table.insert(items, {
                text = name .. " " .. description, -- search string
                name = name,
                description = description,
                trigger = snip.trigger,
                ft = snip.ft,
                preview = {
                  ft = snip.ft,
                  text = docstring,
                },
              })
            end
            return items
          end,
          confirm = function(picker, item)
            picker:close()
            --
            local expand = {}
            require("luasnip").available(function(snippet)
              if snippet.trigger == item.trigger then
                table.insert(expand, snippet)
              end
              return snippet
            end)
            if #expand > 0 then
              vim.cmd(":startinsert!")
              vim.defer_fn(function()
                require("luasnip").snip_expand(expand[1])
              end, 50)
            else
              Snacks.notify.warn("No snippet to expand")
            end
          end,
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
      "<leader>fs",
      function()
        Snacks.picker.sessions()
      end,
      desc = "Find Sessions (tmux)",
    },
    {
      "<leader>fx",
      function()
        Snacks.picker.snippets()
      end,
      desc = "Find Snippets",
    },
  },
}
