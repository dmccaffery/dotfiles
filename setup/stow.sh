#! /usr/bin/env sh

set -eu

SETUP_DIR=$(dirname "$(readlink -f -- "$0")")
. "${SETUP_DIR}/printing.sh"

info "stowing config..."

INSTALL_DIR="${INSTALL_DIR:-${SETUP_DIR}/..}"

if ! hash stow 1> /dev/null 2>&1; then
	warn "stow could not be found; did you forget to install brews?"
else
	info "linking config symlinks to repository"
	mkdir -p "${HOME}/.claude"
	mkdir -p "${HOME}/.config"
	mkdir -p "${HOME}/.config/opencode"
	mkdir -p "${HOME}/.ssh"

	stow --dir="${INSTALL_DIR}" .
fi
