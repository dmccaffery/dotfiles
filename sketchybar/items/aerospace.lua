local function parse_string_to_table(s)
    local result = {}
    for line in s:gmatch("([^\n]+)") do
        table.insert(result, line)
    end
    return result
end

function GetWorkspaces()
    local handle = io.popen("aerospace list-workspaces --all")

    if handle == nil then
        return {}
    end

    local result = handle:read("*a")
    handle:close()

    return parse_string_to_table(result)
end

function GetCurrentWorkspace()
    local handle = io.popen("aerospace list-workspaces --focused")

    if handle == nil then
        return {}
    end

    local result = handle:read("*a")
    handle:close()

    return parse_string_to_table(result)[1]
end

function GetMonitors()
    local handle = io.popen("aerospace list-monitors | awk '{print $1}'")

    if handle == nil then
        return {}
    end

    local result = handle:read("*a")
    handle:close()

    return parse_string_to_table(result)
end

function GetWorkspacesOnMonitor(monitor)
    local handle = io.popen("aerospace list-workspaces --monitor " .. monitor)

    if handle == nil then
        return {}
    end

    local result = handle:read("*a")
    handle:close()

    return parse_string_to_table(result)
end

function GetVisibleWorkspaceOnMonitor(monitor)
    local handle = io.popen(
        "aerospace list-workspaces --monitor " .. monitor .. " --visible"
    )

    if handle == nil then
        return {}
    end

    local result = handle:read("*a")
    handle:close()

    return parse_string_to_table(result)[1]
end

function IsWorkspaceSelected(workspace)
    local available_monitors = GetMonitors()
    for _, monitor in ipairs(available_monitors) do
        local visible_workspace = GetVisibleWorkspaceOnMonitor(monitor)

        if visible_workspace == workspace then
            return true
        end
    end

    return false
end
