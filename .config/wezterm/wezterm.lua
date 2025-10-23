local wezterm = require("wezterm")
local config = wezterm.config_builder()

config.color_scheme = "Catppuccin Mocha"
config.font = wezterm.font(
    "FiraCode Nerd Font",
    { weight = 450, stretch = "Normal", style = "Normal" }
)
config.font_size = 15

config.enable_tab_bar = false

config.window_decorations = "RESIZE"
config.window_background_opacity = 0.90
config.macos_window_background_blur = 10

config.native_macos_fullscreen_mode = false
config.scrollback_lines = 2500

wezterm.on("gui-startup", function()
    local tab, pane, window = wezterm.mux.spawn_window(cmd or {})
    window:gui_window():toggle_fullscreen()
end)

config.keys = {
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
        action = wezterm.action.SplitVertical({ domain = "CurrentPaneDomain" }),
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
}

config.mouse_bindings = {
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
}

return config
