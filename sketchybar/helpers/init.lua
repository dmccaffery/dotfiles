package.cpath = package.cpath
    .. ";"
    .. os.getenv("HOME")
    .. "/.local/share/sketchybar_lua/?.so"

os.execute("(cd helpers && make)")
