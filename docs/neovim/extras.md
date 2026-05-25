---
icon: lucide/puzzle
---

# LazyVim extras

LazyVim extras are pre-packaged plugin bundles for specific languages or workflows. The
following are imported in `lua/config/lazy.lua`:

## Coding

| Extra                                         | Adds                                                         |
| --------------------------------------------- | ------------------------------------------------------------ |
| `lazyvim.plugins.extras.coding.blink`         | [blink.cmp](https://github.com/saghen/blink.cmp) completion. |
| `lazyvim.plugins.extras.coding.luasnip`       | LuaSnip snippet engine.                                      |
| `lazyvim.plugins.extras.coding.yanky`         | Enhanced yank/paste history.                                 |
| `lazyvim.plugins.extras.coding.mini-surround` | `mini.surround` for s/r/d-around motions.                    |

## Debugging

| Extra                             | Adds                                |
| --------------------------------- | ----------------------------------- |
| `lazyvim.plugins.extras.dap.core` | nvim-dap, dap-ui, dap-virtual-text. |
| `lazyvim.plugins.extras.dap.nlua` | Lua adapter (debug nvim itself).    |

## Formatting & linting

| Extra                                        | Adds                                       |
| -------------------------------------------- | ------------------------------------------ |
| `lazyvim.plugins.extras.formatting.prettier` | Prettier via conform.nvim.                 |
| `lazyvim.plugins.extras.linting.eslint`      | eslint_d via nvim-lint + auto-fix on save. |

## Testing

| Extra                              | Adds               |
| ---------------------------------- | ------------------ |
| `lazyvim.plugins.extras.test.core` | neotest framework. |

## Utilities

| Extra                                         | Adds                                                      |
| --------------------------------------------- | --------------------------------------------------------- |
| `lazyvim.plugins.extras.util.octo`            | Octo.nvim — GitHub PRs/issues inside nvim.                |
| `lazyvim.plugins.extras.util.mini-hipatterns` | `mini.hipatterns` — highlight TODO/FIXME and hex colours. |
| `lazyvim.plugins.extras.util.dot`             | Dotfile filetype support.                                 |

## Languages

| Extra                                    | Languages / tooling                                 |
| ---------------------------------------- | --------------------------------------------------- |
| `lazyvim.plugins.extras.lang.docker`     | Dockerfile, docker-compose.                         |
| `lazyvim.plugins.extras.lang.dotnet`     | C# / F# / .NET.                                     |
| `lazyvim.plugins.extras.lang.go`         | gopls, gotests, dap-go.                             |
| `lazyvim.plugins.extras.lang.helm`       | Helm chart templating.                              |
| `lazyvim.plugins.extras.lang.json`       | jsonls, schema completion.                          |
| `lazyvim.plugins.extras.lang.markdown`   | marksman, render-markdown.                          |
| `lazyvim.plugins.extras.lang.python`     | basedpyright, ruff, debugpy.                        |
| `lazyvim.plugins.extras.lang.sql`        | sqls + sqlfluff (see [SQL plugin](plugins.md#sql)). |
| `lazyvim.plugins.extras.lang.terraform`  | terraform-ls.                                       |
| `lazyvim.plugins.extras.lang.toml`       | taplo.                                              |
| `lazyvim.plugins.extras.lang.typescript` | vtsls + prettier + eslint.                          |
| `lazyvim.plugins.extras.lang.yaml`       | yamlls + schemastore.                               |

!!! tip "Pruning extras"
Each language extra installs an LSP, a formatter, and treesitter parsers. Removing extras
you don't use is the single biggest startup-time win — drop the import lines from
`lua/config/lazy.lua` and re-run `:Lazy sync`.
