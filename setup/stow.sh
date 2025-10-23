#!/usr/bin/env bash

set -euo pipefail

SETUP_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
. "${SETUP_DIR}/printing.sh"

INSTALL_DIR="${INSTALL_DIR:-${SETUP_DIR}/..}"

if ! command -v stow 1>/dev/null 2>&1; then
	warn "stow could not be found; did you forget to install brews?"
else
	info "linking config symlinks to repository"
	mkdir -p "${HOME}/.config"
	stow --dir="${INSTALL_DIR}" .
fi
