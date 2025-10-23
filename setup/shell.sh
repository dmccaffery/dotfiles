#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
. "${SCRIPT_DIR}/printing.sh"

ZSH_BIN="$(command -v zsh)"

if [ "${ZSH_BIN:-}" == "${SHELL:-}" ]; then
	warning "zsh is already the default shell"
else
	info "setting default shell to zsh for current user"
	sudo sh -c "printf '%s\n' \"${ZSH_BIN}\" >> /etc/shells"
	chsh -s "${ZSH_BIN}"
fi

ZSHRC_CONFIG="${HOME}/.config/zshrc/.zshrc"

if [ -f "${ZSHRC_CONFIG:-}" ]; then
	info "linking zshrc to config"
	ln -Ffs "${HOME}/.config/zshrc/.zshrc" "${HOME}/.zshrc"
else
	warning "zshrc config is not present and cannot be symlinked; did you forget to
	stow?"
fi

# hush the last login prompt in ttty
touch "${HOME}/.hushlogin"
