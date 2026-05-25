---
icon: simple/neovim
---

# LazyVim

NeoVim configuration is built on [LazyVim](https://www.lazyvim.org/). LazyVim provides the
base plugin spec, sensible defaults, and an "extras" system for opt-in language and tooling
support.

## Layout

```text
.config/nvim/
├── init.lua                # entry point — just `require("config.lazy")`
├── lazy-lock.json          # pinned plugin versions
├── lazyvim.json            # LazyVim metadata
├── lua/
│   └── config/
│       ├── lazy.lua        # plugin spec, extras, performance
│       ├── options.lua     # vim options + UI
│       ├── keymaps.lua     # custom keymaps
│       └── autocmds.lua    # custom auto-commands
│   └── plugins/            # custom plugins (cyberdream, snacks, lualine, …)
│       └── lang/           # language-specific plugins
└── spell/
    └── en.utf-8.add        # personal spelling dictionary
```

## Boot path

```lua title="init.lua"
require("config.lazy")
```

`config.lazy` bootstraps `lazy.nvim` (clones it if missing), prepends it to `runtimepath`, and
then declares the full plugin spec:

```lua title="lua/config/lazy.lua (excerpt)"
require("lazy").setup({
  spec = {
    { "LazyVim/LazyVim", import = "lazyvim.plugins" },
    -- coding
    { import = "lazyvim.plugins.extras.coding.blink" },
    { import = "lazyvim.plugins.extras.coding.luasnip" },
    { import = "lazyvim.plugins.extras.coding.yanky" },
    { import = "lazyvim.plugins.extras.coding.mini-surround" },
    -- … (full list in extras.md)
    -- custom
    { import = "plugins" },
    { import = "plugins.lang" },
  },
  defaults = { lazy = true, version = false },
  checker = { enabled = false, notify = false },
  performance = { ... },
})
```

`defaults.lazy = true` makes every plugin lazy-loaded by default; opt-in eagerness is per-spec
via `lazy = false`.

## Performance tuning

`config.lazy.performance.rtp.disabled_plugins` strips a generous set of vendored Vim runtime
plugins that are essentially dead weight in a modern config:

> `gzip`, `netrw*`, `tar*`, `zip*`, `2html_plugin`, `tohtml`, `matchit`, `vimball*`,
> `spellfile_plugin`, `tutor`, `rplugin`, `syntax`, `synmenu`, `optwin`, `compiler`,
> `bugreport`, `ftplugin`, …

The combined effect is meaningful — startup drops by 30–50 ms on a warm cache.

## Plugin updates

`lazy-lock.json` is committed. Update with `:Lazy update` inside nvim and then commit the new
lock file. `checker.enabled = false` means LazyVim won't nag about updates in the background.
