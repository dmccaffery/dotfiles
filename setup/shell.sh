#! /usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SETUP_DIR="${SCRIPT_DIR}"
. "${SETUP_DIR}/printing.sh"

info "setting up default shell..."

XDG_CONFIG_HOME="${XDG_CONFIG_HOME:-${HOME}/.config}"

ZSH_BIN="$(command -v zsh)"

if [ "${ZSH_BIN:-}" == "${SHELL:-}" ]; then
	warn "zsh is already the default shell"
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
	warn "zshrc config is not present and cannot be symlinked; did you forget to
	stow?"
fi

# hush the last login prompt in tty
touch "${HOME}/.hushlogin"

# setup ssh ask pass (for pinentry)
if [ ! -f /usr/local/bin/ssh-askpass ]; then
	warn "linking ssh-askpass to /usr/local/bin"
	sudo ln -fs "${HOME}/.local/share/scripts/ssh-askpass" /usr/local/bin
fi

# setup completions
COMPLETION_DIR=$(${ZSH_BIN} -lc 'echo "${fpath// /\n}" | grep -i completion')

if command -v "flux" 1> /dev/null 2>&1; then
	flux completion zsh > "${COMPLETION_DIR}/_flux"
fi

PLATFORM_DIR=$(uname -s | tr '[:upper:]' '[:lower:]')
PLATFORM_SCRIPT="${SETUP_DIR}/${PLATFORM_DIR}/shell.sh"

if [ -x "${PLATFORM_SCRIPT}" ]; then
	"${PLATFORM_SCRIPT}"
fi
