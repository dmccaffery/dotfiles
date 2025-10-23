#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
. "${SCRIPT_DIR}/printing.sh"

if ! command -v sketchybar 1>/dev/null 2>&1; then
	warn "sketchybar is not installed or not available; did you forget to
	install brews?"
else
	# install sketchybar lua
	info "installing sketchybar lua"
	(git clone https://github.com/FelixKratz/SbarLua.git /tmp/SbarLua && cd /tmp/SbarLua/ && make install && rm -rf /tmp/SbarLua/)

	# install sketchybar app font
	info "installing sketchybar app font"
	curl -Lfs https://github.com/kvndrsslr/sketchybar-app-font/releases/latest/download/sketchybar-app-font.ttf -o "${HOME}/Library/Fonts/sketchybar-app-font.ttf"
	curl -Lfs https://github.com/kvndrsslr/sketchybar-app-font/releases/latest/download/icon_map.lua -o "${HOME}/.local/share/sketchybar_lua/icon_map.lua"

	mkdir -p "${HOME}/.local/share/sketchybar"
fi
