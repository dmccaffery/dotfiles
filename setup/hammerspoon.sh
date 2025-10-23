#!/usr/bin/env bash

set -euo pipefail

SETUP_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
. "${SETUP_DIR}/printing.sh"

XDG_CONFIG_HOME=${XDG_CONFIG_HOME:-${HOME}/.config}

warning "shutting down hammerspoon"
killall Hammerspoon

info "setting hammerspoon configuration path to use xdg config"
defaults write org.hammerspoon.Hammerspoon MJConfigFile "${XDG_CONFIG_HOME}/hammerspoon/init.lua"

info "downloading spoon install"
temp=$(mktemp)
curl -fsSL \
	https://github.com/Hammerspoon/Spoons/raw/master/Spoons/SpoonInstall.spoon.zip \
	-o "${temp}"

info "unzipping spoon install"
spoons="${XDG_CONFIG_HOME}"/hammerspoon/spoons

mkdir -p "${spoons}"
rm -rf "${spoons}"/SpoonInstall.spoon
unzip "${temp}" -d "${XDG_CONFIG_HOME}"/hammerspoon/spoons
rm -f "${temp}"

info "launching hammerspoon"
open /Applications/Hammerspoon.app
