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

	# The stowed trees live under stow/ in the repo; each is symlinked into its
	# own $HOME target so ~/.config, ~/.claude, etc. stay real directories that
	# stow descends into (rather than folding the whole tree into one symlink).
	# Claude Code's plan-mode runtime stays at the repo-root .claude/plans (per
	# the plansDirectory setting), so it is outside this stow source by design.
	stow --dir="${INSTALL_DIR}/stow/.claude" --target="${HOME}/.claude" .
	stow --dir="${INSTALL_DIR}/stow/.config" --target="${HOME}/.config" .
	stow --dir="${INSTALL_DIR}/stow/.local" --target="${HOME}/.local" .
	stow --dir="${INSTALL_DIR}/stow/.terminfo" --target="${HOME}/.terminfo" .
	stow --dir="${INSTALL_DIR}/stow/.ssh" --target="${HOME}/.ssh" .
	stow --dir="${INSTALL_DIR}/stow/Library" --target="${HOME}/Library" .

	ln -Ffs "${INSTALL_DIR}/.zshenv" "${HOME}/.zshenv"
fi
