---
icon: lucide/palette
---

# Cyberdream

[cyberdream](https://github.com/scottmckendry/cyberdream.nvim) is a vivid, neon-inflected
colour scheme by Scott McKendry. This repo uses it everywhere a tool exposes a theme.

## Palette reference

The canonical hex values used across this repo. Light variant values come from Ghostty's
`cyberdream-light` theme.

### Dark

| Role       | Hex       | Sample                                                                                                                         |
| ---------- | --------- | ------------------------------------------------------------------------------------------------------------------------------ |
| Background | `#16181A` | <span style="display:inline-block;width:1em;height:1em;vertical-align:middle;background:#16181A;border:1px solid #555"></span> |
| Foreground | `#FFFFFF` | <span style="display:inline-block;width:1em;height:1em;vertical-align:middle;background:#FFFFFF;border:1px solid #555"></span> |
| Subtle     | `#3C4048` | <span style="display:inline-block;width:1em;height:1em;vertical-align:middle;background:#3C4048;border:1px solid #555"></span> |
| Red        | `#FF6E5E` | <span style="display:inline-block;width:1em;height:1em;vertical-align:middle;background:#FF6E5E;border:1px solid #555"></span> |
| Green      | `#5EFF6C` | <span style="display:inline-block;width:1em;height:1em;vertical-align:middle;background:#5EFF6C;border:1px solid #555"></span> |
| Yellow     | `#F1FF5E` | <span style="display:inline-block;width:1em;height:1em;vertical-align:middle;background:#F1FF5E;border:1px solid #555"></span> |
| Blue       | `#5EA1FF` | <span style="display:inline-block;width:1em;height:1em;vertical-align:middle;background:#5EA1FF;border:1px solid #555"></span> |
| Purple     | `#BD5EFF` | <span style="display:inline-block;width:1em;height:1em;vertical-align:middle;background:#BD5EFF;border:1px solid #555"></span> |
| Cyan       | `#5EF1FF` | <span style="display:inline-block;width:1em;height:1em;vertical-align:middle;background:#5EF1FF;border:1px solid #555"></span> |
| Pink       | `#FF5EA0` | <span style="display:inline-block;width:1em;height:1em;vertical-align:middle;background:#FF5EA0;border:1px solid #555"></span> |
| Orange     | `#FFBD5E` | <span style="display:inline-block;width:1em;height:1em;vertical-align:middle;background:#FFBD5E;border:1px solid #555"></span> |

### Light

| Role       | Hex       | Sample                                                                                                                         |
| ---------- | --------- | ------------------------------------------------------------------------------------------------------------------------------ |
| Background | `#FFFFFF` | <span style="display:inline-block;width:1em;height:1em;vertical-align:middle;background:#FFFFFF;border:1px solid #555"></span> |
| Foreground | `#16181A` | <span style="display:inline-block;width:1em;height:1em;vertical-align:middle;background:#16181A;border:1px solid #555"></span> |
| Red        | `#D11500` | <span style="display:inline-block;width:1em;height:1em;vertical-align:middle;background:#D11500;border:1px solid #555"></span> |
| Green      | `#008B0C` | <span style="display:inline-block;width:1em;height:1em;vertical-align:middle;background:#008B0C;border:1px solid #555"></span> |
| Blue       | `#0057D1` | <span style="display:inline-block;width:1em;height:1em;vertical-align:middle;background:#0057D1;border:1px solid #555"></span> |
| Purple     | `#A018FF` | <span style="display:inline-block;width:1em;height:1em;vertical-align:middle;background:#A018FF;border:1px solid #555"></span> |
| Cyan       | `#008C99` | <span style="display:inline-block;width:1em;height:1em;vertical-align:middle;background:#008C99;border:1px solid #555"></span> |

## Source

Upstream: [scottmckendry/cyberdream.nvim](https://github.com/scottmckendry/cyberdream.nvim) —
includes ports for many tools in `extras/`.

## Adding a new tool

When integrating cyberdream into a new tool:

1. Check upstream `extras/` for an existing port.
2. If absent, derive from the hex table above. Map your tool's "info / link / accent /
   warning / error" roles to blue / purple / pink / yellow / red respectively.
3. Cross-reference [per-tool](per-tool.md) for examples already in this repo.
