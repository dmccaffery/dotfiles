local function setDisplays(PaperWM)
  local allScreens = hs.screen.allScreens()

  if #allScreens == 1 then
    PaperWM.window_filter:setDefaultFilter()
    return
  end

  local screens = {}

  for _, screen in ipairs(allScreens) do
    local name = screen:name()
    local builtIn = name:find("^Built%-in") ~= nil

    if not builtIn then
      table.insert(screens, screen:id())
      print("using screen: " .. name)
    else
      print("skipping: " .. name)
    end
  end

  PaperWM.window_filter:setScreens(screens)
end

local function watchDisplays()
  local PaperWM = hs.loadSpoon("PaperWM")
  setDisplays(PaperWM)
  PaperWM:start()
end

local screenWatcher = hs.screen.watcher.new(watchDisplays)
screenWatcher:start()

spoon.SpoonInstall.repos.PaperWM = {
  url = "https://github.com/mogenson/PaperWM.spoon",
  desc = "PaperWM.spoon",
  branch = "release",
}

local hyperKey = { "ctrl", "alt", "cmd", "shift" }

spoon.SpoonInstall:andUse("PaperWM", {
  repo = "PaperWM",
  config = { screen_margin = 16, window_gap = 2 },
  start = true,
  fn = setDisplays,
  hotkeys = {
    -- switch to a new focused window in tiled grid
    focus_left = { hyperKey, "left" },
    focus_right = { hyperKey, "right" },
    focus_up = { hyperKey, "up" },
    focus_down = { hyperKey, "down" },

    -- switch windows by cycling forward/backward
    -- (forward = down or right, backward = up or left)
    focus_prev = { hyperKey, "k" },
    focus_next = { hyperKey, "j" },

    -- move windows around in tiled grid
    swap_left = { { "alt", "cmd", "shift" }, "left" },
    swap_right = { { "alt", "cmd", "shift" }, "right" },
    swap_up = { { "alt", "cmd", "shift" }, "up" },
    swap_down = { { "alt", "cmd", "shift" }, "down" },

    -- alternative: swap entire columns, rather than
    -- individual windows (to be used instead of
    -- swap_left / swap_right bindings)
    -- swap_column_left = {{"alt", "cmd", "shift"}, "left"},
    -- swap_column_right = {{"alt", "cmd", "shift"}, "right"},

    -- position and resize focused window
    center_window = { hyperKey, "c" },
    full_width = { hyperKey, "f" },
    cycle_width = { { "alt", "cmd" }, "r" },
    reverse_cycle_width = { { "ctrl", "alt", "cmd" }, "r" },
    cycle_height = { { "alt", "cmd", "shift" }, "r" },
    reverse_cycle_height = { { "ctrl", "alt", "cmd", "shift" }, "r" },

    -- increase/decrease width
    increase_width = { hyperKey, "l" },
    decrease_width = { hyperKey, "h" },

    -- move focused window into / out of a column
    slurp_in = { hyperKey, "i" },
    barf_out = { hyperKey, "o" },

    -- move the focused window into / out of the tiling layer
    toggle_floating = { hyperKey, "`" },

    -- focus the first / second / etc window in the current space
    focus_window_1 = { { "cmd", "shift" }, "1" },
    focus_window_2 = { { "cmd", "shift" }, "2" },
    focus_window_3 = { { "cmd", "shift" }, "3" },
    focus_window_4 = { { "cmd", "shift" }, "4" },
    focus_window_5 = { { "cmd", "shift" }, "5" },
    focus_window_6 = { { "cmd", "shift" }, "6" },
    focus_window_7 = { { "cmd", "shift" }, "7" },
    focus_window_8 = { { "cmd", "shift" }, "8" },
    focus_window_9 = { { "cmd", "shift" }, "9" },

    -- switch to a new Mission Control space
    switch_space_l = { { "cmd" }, "," },
    switch_space_r = { { "cmd" }, "." },
    switch_space_1 = { { "cmd" }, "1" },
    switch_space_2 = { { "cmd" }, "2" },
    switch_space_3 = { { "cmd" }, "3" },
    switch_space_4 = { { "cmd" }, "4" },
    switch_space_5 = { { "cmd" }, "5" },
    switch_space_6 = { { "cmd" }, "6" },
    switch_space_7 = { { "cmd" }, "7" },
    switch_space_8 = { { "cmd" }, "8" },
    switch_space_9 = { { "cmd" }, "9" },

    -- move focused window to a new space and tile
    move_window_1 = { hyperKey, "1" },
    move_window_2 = { hyperKey, "2" },
    move_window_3 = { hyperKey, "3" },
    move_window_4 = { hyperKey, "4" },
    move_window_5 = { hyperKey, "5" },
    move_window_6 = { hyperKey, "6" },
    move_window_7 = { hyperKey, "7" },
    move_window_8 = { hyperKey, "8" },
    move_window_9 = { hyperKey, "9" },
  },
})

WarpMouse = hs.loadSpoon("WarpMouse")
WarpMouse:start()
