local wezterm = require("wezterm")
local config = wezterm.config_builder()

wezterm.on("gui-startup", function(cmd)
    local _, _, window = wezterm.mux.spawn_window(cmd or {})
    window:gui_window():toggle_fullscreen()
end)

config = {
    -- theme
    color_scheme = "Catppuccin Mocha",
    font = wezterm.font(
        "FiraCode Nerd Font",
        { weight = 450, stretch = "Normal", style = "Normal" }
    ),
    font_size = 15,

    -- cursor
    default_cursor_style = "SteadyBlock",
    cursor_blink_ease_out = "Constant",
    cursor_blink_ease_in = "Constant",
    cursor_blink_rate = 0,

    -- window
    enable_tab_bar = false,
    window_decorations = "RESIZE",
    window_background_opacity = 0.90,
    window_padding = {
        left = 0,
        right = 0,
        top = 15,
        bottom = 0,
    },
    macos_window_background_blur = 10,
    native_macos_fullscreen_mode = false,

    -- term
    term = "xterm-kitty",
    enable_kitty_graphics = true,

    -- performance
    scrollback_lines = 2500,
    webgpu_power_preference = "HighPerformance",
    max_fps = 120,

    -- keyboard
    keys = {
        {
            key = "d",
            mods = "CMD",
            action = wezterm.action.SplitHorizontal({
                domain = "CurrentPaneDomain",
            }),
        },
        {
            key = "d",
            mods = "CMD|SHIFT",
            action = wezterm.action.SplitVertical({
                domain = "CurrentPaneDomain",
            }),
        },
        -- Make Option-Left equivalent to Alt-b which many line editors interpret as backward-word
        {
            key = "LeftArrow",
            mods = "OPT",
            action = wezterm.action({ SendString = "\x1bb" }),
        },
        -- Make Option-Right equivalent to Alt-f; forward-word
        {
            key = "RightArrow",
            mods = "OPT",
            action = wezterm.action({ SendString = "\x1bf" }),
        },
        -- Select next tab with cmd-opt-left/right arrow
        {
            key = "LeftArrow",
            mods = "CMD|OPT",
            action = wezterm.action.ActivateTabRelative(-1),
        },
        {
            key = "RightArrow",
            mods = "CMD|OPT",
            action = wezterm.action.ActivateTabRelative(1),
        },
        -- Select next pane with cmd-left/right arrow
        {
            key = "LeftArrow",
            mods = "CMD",
            action = wezterm.action({ ActivatePaneDirection = "Prev" }),
        },
        {
            key = "RightArrow",
            mods = "CMD",
            action = wezterm.action({ ActivatePaneDirection = "Next" }),
        },
    },

    -- mouse
    mouse_bindings = {
        {
            event = { Down = { streak = 1, button = "Left" } },
            mods = "CMD|ALT",
            action = wezterm.action.SelectTextAtMouseCursor("Block"),
            alt_screen = "Any",
        },
        {
            event = { Down = { streak = 4, button = "Left" } },
            action = wezterm.action.SelectTextAtMouseCursor("SemanticZone"),
            mods = "NONE",
        },
    },
}

return config
