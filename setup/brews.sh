#!/usr/bin/env bash

set -euo pipefail

SETUP_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
. "${SETUP_DIR}/printing.sh"

info "setting up homebrew..."

INSTALL_DIR="${INSTALL_DIR:-${SETUP_DIR}/..}"

brew_cmd=$(command -v brew 2>&1)

if [ -z "${brew_cmd:-}" ]; then
	brew_cmd=/opt/homebrew/bin/brew
fi

if [ ! -x "${brew_cmd:-}" ]; then
	warn "brew could not be found; please ensure homebrew is installed"
	return 0
fi

eval "$(${brew_cmd} shellenv)"

info "installing common brews..."
brew bundle install --force --cleanup --file="${INSTALL_DIR}/Brewfile"

if uname -n | grep -F -q "bitwisemedia.uk" 1> /dev/null 2>&1; then
	info "install bitwise brews..."
	brew bundle install --force --file="${INSTALL_DIR}/Brewfile.bitwise"
fi

# always reinstall stow because it dynamically links to the version of perl
# included in the os, which can change
brew reinstall stow 1> /dev/null 2>&1

# setup buildx
if command -v docker-buildx 1> /dev/null 2>&1; then
	mkdir -p "${HOME}/.docker/cli-plugins"
	ln -fns $(command -v docker-buildx 2> /dev/null) "${HOME}/.docker/cli-plugins"
fi

# cleanup services
info "cleaning up brew services"
brew services cleanup
