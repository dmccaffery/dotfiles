package.path = package.path
    .. ";"
    .. os.getenv("HOME")
    .. "/.local/share/sketchybar_lua/?.lua"

-- Require the sketchybar module
SBar = require("sketchybar")

-- Set the bar name, if you are using another bar instance than sketchybar
-- sbar.set_bar_name("bottom_bar")

-- Bundle the entire initial configuration into a single message to sketchybar
SBar.begin_config()
require("bar")
require("default")
-- require("test")
require("items")
SBar.end_config()

-- Run the event loop of the sketchybar module (without this there will be no
-- callback functions executed in the lua module)
SBar.event_loop()
