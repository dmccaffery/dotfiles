local colors = require("colors")
local icons = require("icons")
local settings = require("settings")
local app_icons = require("icon_map")

local spaces = {}

local workspaces = GetWorkspaces()
local current_workspace = GetCurrentWorkspace()

local function split(str, sep)
    local result = {}
    local regex = ("([^%s]+)"):format(sep)
    for each in str:gmatch(regex) do
        table.insert(result, each)
    end
    return result
end

for i, workspace in ipairs(workspaces) do
    local selected = workspace == current_workspace
    local space = SBar.add("item", "item." .. workspace, {
        icon = {
            font = {
                family = settings.font.numbers,
            },
            string = workspace,
            padding_left = settings.items.padding.left,
            padding_right = settings.items.padding.left / 2,
            color = settings.items.default_color(i),
            highlight_color = settings.items.highlight_color(i),
            highlight = selected,
        },
        label = {
            padding_right = settings.items.padding.right,
            color = settings.items.default_color(i),
            highlight_color = settings.items.highlight_color(i),
            font = settings.icons,
            y_offset = -1,
            highlight = selected,
        },
        padding_right = 1,
        padding_left = 1,
        background = {
            color = settings.items.colors.background,
            border_width = 1,
            height = settings.items.height,
            border_color = selected and settings.items.highlight_color(i)
                or settings.items.default_color(i),
        },
        popup = {
            background = {
                border_width = 5,
                border_color = colors.black,
            },
        },
    })

    spaces[i] = space

    -- Define the icons for open apps on each space initially
    SBar.exec(
        "aerospace list-windows --workspace "
            .. workspace
            .. " --format '%{app-name}' --json ",
        function(apps)
            local icon_line = ""
            local no_app = true
            for _, app in ipairs(apps) do
                no_app = false
                local app_name = app["app-name"]
                local lookup = app_icons[app_name]
                local icon = (
                    (lookup == nil) and app_icons["default"] or lookup
                )
                icon_line = icon_line .. " " .. icon
            end

            if no_app then
                icon_line = " —"
            end

            SBar.animate("tanh", 10, function()
                space:set({
                    label = icon_line,
                })
            end)
        end
    )

    -- Padding space between each item
    SBar.add("item", "item." .. workspace .. "padding", {
        script = "",
        width = settings.items.gap,
    })

    -- Item popup
    local space_popup = SBar.add("item", {
        position = "popup." .. space.name,
        padding_left = 5,
        padding_right = 0,
        background = {
            drawing = true,
            image = {
                corner_radius = 9,
                scale = 0.2,
            },
        },
    })

    space:subscribe("aerospace_workspace_change", function(env)
        local active = env.FOCUSED_WORKSPACE == workspace
        space:set({
            icon = {
                highlight = active,
            },
            label = {
                highlight = active,
            },
            background = {
                border_color = active and settings.items.highlight_color(i)
                    or settings.items.default_color(i),
            },
        })
    end)

    space:subscribe("mouse.clicked", function(env)
        local id = split(env.NAME, ".")[2]
        print("SID " .. env.NAME)
        if env.BUTTON == "other" then
            space_popup:set({
                background = {
                    image = "item." .. id,
                },
            })
            space:set({
                popup = {
                    drawing = "toggle",
                },
            })
        else
            SBar.exec("aerospace workspace " .. id)
        end
    end)

    space:subscribe("mouse.exited", function(_)
        space:set({
            popup = {
                drawing = false,
            },
        })
    end)
end

local space_window_observer = SBar.add("item", {
    drawing = false,
    updates = true,
})

-- Event handles
space_window_observer:subscribe("space_windows_change", function(_)
    for i, workspace in ipairs(workspaces) do
        SBar.exec(
            "aerospace list-windows --workspace "
                .. workspace
                .. " --format '%{app-name}' --json ",
            function(apps)
                local icon_line = ""
                local no_app = true
                for _, app in ipairs(apps) do
                    no_app = false
                    local app_name = app["app-name"]
                    local lookup = app_icons[app_name]
                    local icon = (
                        (lookup == nil) and app_icons["default"] or lookup
                    )
                    icon_line = icon_line .. " " .. icon
                end

                if no_app then
                    icon_line = " —"
                end

                SBar.animate("tanh", 10, function()
                    spaces[i]:set({
                        label = icon_line,
                    })
                end)
            end
        )
    end
end)

space_window_observer:subscribe("aerospace_focus_change", function(_)
    for i, workspace in ipairs(workspaces) do
        SBar.exec(
            "aerospace list-windows --workspace "
                .. workspace
                .. " --format '%{app-name}' --json ",
            function(apps)
                local icon_line = ""
                local no_app = true
                for _, app in ipairs(apps) do
                    no_app = false
                    local app_name = app["app-name"]
                    local lookup = app_icons[app_name]
                    local icon = (
                        (lookup == nil) and app_icons["default"] or lookup
                    )
                    icon_line = icon_line .. " " .. icon
                end

                if no_app then
                    icon_line = " —"
                end

                SBar.animate("tanh", 10, function()
                    spaces[i]:set({
                        label = icon_line,
                    })
                end)
            end
        )
    end
end)
