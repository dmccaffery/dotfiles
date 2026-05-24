#! /usr/bin/env sh

set -eu

SCRIPT_DIR=$(dirname "$(readlink -f -- "$0")")
SETUP_DIR=$(realpath "${SCRIPT_DIR}/../")

# shellcheck source=./../printing.sh
. "${SETUP_DIR}/printing.sh"

info "setting up homebrew..."
brew_cmd=$(command -v brew 2>&1)

if [ -z "${brew_cmd:-}" ]; then
	brew_cmd=/opt/homebrew/bin/brew
fi

if [ ! -x "${brew_cmd:-}" ]; then
	warn "brew could not be found; please ensure homebrew is installed"
	return 0
fi

eval "$(${brew_cmd} shellenv)"

HOMEBREW_BUNDLE_FILE="${XDG_DATA_HOME}/homebrew/Brewfile"
export HOMEBREW_BUNDLE_FILE

HOMEBREW_BUNDLE_DIR="${XDG_CONFIG_HOME}/homebrew"

if [ ! -r "${HOMEBREW_BUNDLE_FILE}" ]; then
	# allow user to pick an existing brewfile, or type a new name to dump current brews into
	selection=$(find "${HOMEBREW_BUNDLE_DIR}/" -mindepth 1 -maxdepth 1 -type f -name 'Brewfile.*' -print |
		awk -F/ '{ name=$NF; sub(/^Brewfile\./, "", name); print name "\t" $0 }' |
		fzf --header='Pick a brew bundle or type a new name' \
			--delimiter='\t' \
			--with-nth=1 \
			--preview='cat {2}' \
			--print-query || true)

	query=$(printf '%s\n' "${selection}" | sed -n '1p')
	chosen=$(printf '%s\n' "${selection}" | sed -n '2p' | cut -f2)

	if [ -n "${chosen:-}" ]; then
		brewfile="${chosen}"
	elif [ -n "${query:-}" ]; then
		brewfile="${HOMEBREW_BUNDLE_DIR}/Brewfile.${query}"
		touch "${brewfile}"
	else
		brewfile=""
	fi

	if [ -n "${brewfile:-}" ]; then
		mkdir -p "$(dirname "${HOMEBREW_BUNDLE_FILE}")"
		ln -Ffs "${brewfile}" "${HOMEBREW_BUNDLE_FILE}"
	fi
fi

info "installing / upgrading from brewfile..."
brew bundle dump --brews --casks --mas --force
brew bundle install --upgrade --zap --cleanup --force
brew reinstall stow 1> /dev/null 2>&1

# setup buildx
if hash docker-buildx 1> /dev/null 2>&1; then
	mkdir -p "${HOME}/.docker/cli-plugins"
	ln -fns "$(command -v docker-buildx 2> /dev/null)" "${HOME}/.docker/cli-plugins"
fi

info "cleaning up brew services"
brew services cleanup --quiet

info "cleaning up brew cache"
brew cleanup --quiet
