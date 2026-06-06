#! /usr/bin/env sh

set -eu

SETUP_DIR=$(dirname "$(readlink -f -- "$0")")
. "${SETUP_DIR}/printing.sh"

info "stowing config..."

INSTALL_DIR="$(realpath "${INSTALL_DIR:-${SETUP_DIR}/..}")"

if ! command -v stow 1> /dev/null 2>&1; then
	warn "stow could not be found; did you forget to install brews?"
else
	info "linking config symlinks to repository"
	mkdir -p "${HOME}/.claude"
	mkdir -p "${HOME}/.config/codex"
	mkdir -p "${HOME}/.config/opencode"
	mkdir -p "${HOME}/.config/zsh"
	mkdir -p "${HOME}/.local/share/scripts"
	mkdir -p "${HOME}/.local/share/wallpapers"
	mkdir -p "${HOME}/.ssh"

	stow --dir="${INSTALL_DIR}/.config" --target="${HOME}/.config" .
	stow --dir="${INSTALL_DIR}/.local" --target="${HOME}/.local" .
	stow --dir="${INSTALL_DIR}/.terminfo" --target="${HOME}/.terminfo" .
	stow --dir="${INSTALL_DIR}/.ssh" --target="${HOME}/.ssh" .
	stow --dir="${INSTALL_DIR}/Library" --target="${HOME}/Library" .

	ln -Ffs "${INSTALL_DIR}/.zshenv" "${HOME}/.zshenv"
fi
