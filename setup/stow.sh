#!/usr/bin/env bash

set -euo pipefail

SETUP_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
. "${SETUP_DIR}/printing.sh"

info "stowing config..."

INSTALL_DIR="${INSTALL_DIR:-${SETUP_DIR}/..}"

if ! command -v stow 1> /dev/null 2>&1; then
	warn "stow could not be found; did you forget to install brews?"
else
	info "linking config symlinks to repository"
	mkdir -p "${HOME}/.config"
	mkdir -p "${HOME}/.ssh"
	mkdir -p "${HOME}/.gemini"

	stow --dir="${INSTALL_DIR}" .
fi
