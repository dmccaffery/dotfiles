#! /usr/bin/env sh

set -eu

SETUP_DIR=$(dirname "$(readlink -f -- "$0")")
. "${SETUP_DIR}/printing.sh"

info "setting up default shell..."

XDG_CONFIG_HOME="${XDG_CONFIG_HOME:-${HOME}/.config}"

ZSH_BIN="$(command -v zsh)"

if [ "${ZSH_BIN:-}" = "${SHELL:-}" ]; then
	warn "zsh is already the default shell"
else
	info "setting default shell to zsh for current user"
	sudo sh -c "printf '%s\n' \"${ZSH_BIN}\" >> /etc/shells"
	chsh -s "${ZSH_BIN}"
fi

ZSH_CONFIG="${XDG_CONFIG_HOME}/zsh"

ZSHRC_CONFIG="${ZSH_CONFIG}/.zshrc"
if [ -f "${ZSHRC_CONFIG:-}" ]; then
	info "linking zshrc to config"
	ln -Ffs "${ZSHRC_CONFIG}" "${HOME}/.zshrc"
else
	warn "zshrc config is not present and cannot be symlinked; did you forget to stow?"
fi

ZSHENV_CONFIG="${ZSH_CONFIG}/.zshenv"
if [ -f "${ZSHENV_CONFIG:-}" ]; then
	info "linking zshenv to config"
	ln -Ffs "${ZSHENV_CONFIG}" "${HOME}/.zshenv"
else
	warn "zshenv config is not present and cannot be symlinked; did you forget to stow?"
fi

# hush the last login prompt in tty
touch "${HOME}/.hushlogin"

# setup ssh ask pass (for pinentry)
if [ ! -f /usr/local/bin/ssh-askpass ]; then
	warn "linking ssh-askpass to /usr/local/bin"
	sudo ln -fs "${HOME}/.local/share/scripts/ssh-askpass" /usr/local/bin
fi

# setup completions
# shellcheck disable=SC2016
COMPLETION_DIR=$(${ZSH_BIN} -lc 'echo "${fpath// /\n}" | grep -i completion')

if hash flux 1> /dev/null 2>&1; then
	flux completion zsh > "${COMPLETION_DIR}/_flux"
fi
