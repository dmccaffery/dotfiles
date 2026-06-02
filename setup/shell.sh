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

if command -v flux 1> /dev/null 2>&1; then
	flux completion zsh > "${COMPLETION_DIR}/_flux"
fi
