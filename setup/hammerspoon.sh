#!/usr/bin/env bash

set -euo pipefail

SETUP_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
. "${SETUP_DIR}/printing.sh"

XDG_CONFIG_HOME=${XDG_CONFIG_HOME:-${HOME}/.config}

warning "shutting down hammerspoon"
killall Hammerspoon 2>/dev/null || true

info "setting hammerspoon configuration path to use xdg config"
defaults write org.hammerspoon.Hammerspoon MJConfigFile "${XDG_CONFIG_HOME}/hammerspoon/init.lua"

spoons="${XDG_CONFIG_HOME}"/hammerspoon/spoons
mkdir -p "${spoons}"

info "downloading spoon install"
temp=$(mktemp)
curl -fsSL \
	https://github.com/Hammerspoon/Spoons/raw/master/Spoons/SpoonInstall.spoon.zip \
	-o "${temp}"

info "unzipping spoon install"
spoon="${spoons}/SpoonInstall.spoon"
rm -rf "${spoon}" 1>/dev/null 2>&1
mkdir -p "${spoon}"
tar -xvf "${temp}" --strip-components=1 -C "${spoon}"
rm -f "${temp}"

info "downloading warp mouse"
temp=$(mktemp)
curl -fsSL \
	https://github.com/mogenson/WarpMouse.spoon/archive/refs/heads/main.zip \
	-o "${temp}"

info "unzipping warp mouse"
spoon="${spoons}/WarpMouse.spoon"
rm -rf "${spoon}" 1>/dev/null 2>&1
mkdir -p "${spoon}"
tar -xvf "${temp}" --strip-components=1 -C "${spoon}"
rm -f "${temp}"

info "launching hammerspoon"
open /Applications/Hammerspoon.app
