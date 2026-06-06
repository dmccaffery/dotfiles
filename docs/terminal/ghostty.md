---
icon: simple/ghostty
---

# Ghostty

[Ghostty](https://ghostty.org) is the configured terminal emulator. Config lives in
`.config/ghostty/config` and themes in `.config/ghostty/themes/`.

## Config highlights

```ini title=".config/ghostty/config"
# theme — follows system appearance
theme = dark:cyberdream,light:cyberdream-light

# font
font-family = "Iosevka NF"
font-size = 15
font-thicken = true
font-thicken-strength = 1
adjust-cell-height = 1

# window
background-opacity = 0.9
background-blur = 16
window-padding-x = 5,5
window-padding-y = 5,5
window-colorspace = display-p3
window-decoration = true
window-save-state = always
macos-titlebar-style = hidden

# mouse / cursor
mouse-hide-while-typing = true
focus-follows-mouse = true
cursor-click-to-move = true

# splits
unfocused-split-opacity = 0.50

# clipboard
clipboard-read = allow
clipboard-write = allow
```

| Key                              | Effect                                                                                                      |
| -------------------------------- | ----------------------------------------------------------------------------------------------------------- |
| `theme = dark:…,light:…`         | Auto-switches between cyberdream variants with the system appearance.                                       |
| `font-thicken = true`            | Synthetic bolding fix for Iosevka on macOS — small, important.                                              |
| `background-opacity 0.9`         | 90% opacity with 16 px blur. Tweak to taste.                                                                |
| `window-colorspace = display-p3` | Wide-gamut color on supported displays.                                                                     |
| `macos-option-as-alt`            | Makes the Option key send Alt, which most TUIs and Vim expect.                                              |
| `clipboard-read/write = allow`   | Allow apps to read/write the system clipboard via OSC 52 — required for yank-over-SSH from tmux and Neovim. |

## Themes

`.config/ghostty/themes/` ships two custom themes:

<!-- markdownlint-disable MD046 -->

=== "cyberdream (dark)"

    ```ini
    background = #16181A
    foreground = #FFFFFF
    cursor-color = #FFFFFF
    selection-background = #3C4048

    palette = 4=#5EA1FF   # blue
    palette = 5=#BD5EFF   # purple
    palette = 6=#5EF1FF   # cyan
    # ...
    ```

=== "cyberdream-light"

    ```ini
    background = #FFFFFF
    foreground = #16181A
    cursor-color = #16181A
    selection-background = #ACACAC

    palette = 4=#0057D1
    palette = 5=#A018FF
    palette = 6=#008C99
    # ...
    ```

## Auto-updates

`auto-update = download` and `quit-after-last-window-closed-delay = 5s` make Ghostty hand-off
between launches feel like one continuous app.
