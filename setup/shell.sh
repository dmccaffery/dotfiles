#!/usr/bin/env bash

set -euo pipefail

SETUP_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
. "${SETUP_DIR}/printing.sh"

XDG_CONFIG_HOME=${XDG_CONFIG_HOME:-${HOME}/.config}

ZSH_BIN="$(command -v zsh)"

if [ "${ZSH_BIN:-}" == "${SHELL:-}" ]; then
	warning "zsh is already the default shell"
else
	info "setting default shell to zsh for current user"
	sudo sh -c "printf '%s\n' \"${ZSH_BIN}\" >> /etc/shells"
	chsh -s "${ZSH_BIN}"
fi

ZSHRC_CONFIG="${XDG_CONFIG_HOME}/zshrc/.zshrc"

if [ -f "${ZSHRC_CONFIG:-}" ]; then
	info "linking zshrc to config"
	ln -Ffs "${XDG_CONFIG_HOME}/zshrc/.zshrc" "${HOME}/.zshrc"
else
	warning "zshrc config is not present and cannot be symlinked; did you forget to
	stow?"
fi

# hush the last login prompt in tty
touch "${HOME}/.hushlogin"
