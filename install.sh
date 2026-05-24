#!/usr/bin/env sh

set -eu

INSTALL_DIR=$(dirname "$(readlink -f -- "$0")")
. "${INSTALL_DIR}/setup/printing.sh"

if [ -n "${1:-}" ]; then
	warn "custom stages defined: ${*}"
else
	set -- xdg requirements defaults stow packages shell
fi

"${INSTALL_DIR}/setup/stage.sh" "$@"

# SCRIPT_DIR="${INSTALL_DIR}/.local/share/scripts"
# "${SCRIPT_DIR}/git-github-auth"
# "${SCRIPT_DIR}/get-sk-ssh"

fastfetch
printf '\n'
