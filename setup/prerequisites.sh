#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
. "${SCRIPT_DIR}/printing.sh"

info "installing cli tools"
if xcode-select -p >/dev/null; then
	warning "cli tools are already installed"
else
	xcode-select --install
	sudo xcodebuild -license accept
fi

info "installing homebrew"
export HOMEBREW_CASK_OPTS="--appdir=/Applications"

if hash brew &>/dev/null; then
	warning "homebrew already installed"
else
	sudo --validate
	NONINTERACTIVE=1 /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install.sh)"
fi
