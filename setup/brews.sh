#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
. "${SCRIPT_DIR}/printing.sh"

brew_cmd=$(command -v brew 2>&1)

if [ -z "${brew_cmd:-}" ]; then
	brew_cmd=/opt/homebrew/bin/brew
fi

if [ ! -x "${brew_cmd:-}" ]; then
	warn "brew could not be found; please ensure homebrew is installed"
else
	eval "$(${brew_cmd} shellenv)"

	info "installing brews"
	brew bundle install --force --cleanup --file="${INSTALL_DIR}/Brewfile"

	# always reinstall stow because it dynamically links to the version of perl
	# included in the os, which can change
	brew reinstall stow 1>/dev/null 2>&1
fi
