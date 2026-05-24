#! /usr/bin/env sh

set -eu

SCRIPT_DIR=$(dirname "$(readlink -f -- "$0")")
SETUP_DIR=$(realpath "${SCRIPT_DIR}/../")

# shellcheck source=./../printing.sh
. "${SETUP_DIR}/printing.sh"

info "installing requirements..."

if xcode-select -p > /dev/null; then
	warn "cli tools are already installed"
else
	xcode-select --install
	sudo xcodebuild -license accept
fi
if hash brew 1> /dev/null 2>&1; then
	warn "homebrew already installed"
else
	info "installing homebrew"
	sudo --validate
	NONINTERACTIVE=1 /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install.sh)"
fi

info "installing required packages..."
brew bundle check --file="${SCRIPT_DIR}/Brewfile.requirements"
