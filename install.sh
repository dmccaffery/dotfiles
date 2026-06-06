#!/usr/bin/env sh

set -eu

INSTALL_DIR=$(dirname "$(readlink -f -- "$0")")
. "${INSTALL_DIR}/setup/printing.sh"

if [ -n "${1:-}" ]; then
	warn "custom stages defined: ${*}"
else
	set -- xdg requirements config stow build packages shell
fi

"${INSTALL_DIR}/setup/stage.sh" "$@"

fastfetch
printf '\n'
