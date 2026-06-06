#! /usr/bin/env sh

set -eu

SETUP_DIR=$(dirname "$(readlink -f -- "$0")")
. "${SETUP_DIR}/printing.sh"

info "building dot..."

INSTALL_DIR="$(realpath "${INSTALL_DIR:-${SETUP_DIR}/..}")"

if ! command -v go 1> /dev/null 2>&1; then
	warn "go could not be found; did you forget to install brews?"
	exit 0
fi

bin_dir="${HOME}/.local/bin"
scripts_dir="${HOME}/.local/share/scripts"

mkdir -p "${bin_dir}"

# The applet symlinks are written as siblings of the stowed scripts. If the
# scripts dir is itself a (folded) stow symlink, writing into it would land
# inside the repo — refuse rather than pollute the working tree.
if [ -L "${scripts_dir}" ]; then
	error "${scripts_dir} is a symlink (folded stow); refusing to write applet links into the repo"
	exit 1
fi
mkdir -p "${scripts_dir}"

version=$(git -C "${INSTALL_DIR}" describe --tags --always --dirty 2> /dev/null || echo dev)

info "go build -> ${bin_dir}/dot (${version})"
(
	cd "${INSTALL_DIR}" || exit 1
	go build -trimpath -ldflags "-X main.version=${version}" -o "${bin_dir}/dot" ./cmd/dot
)

# Self-check: a wrong-arch or corrupt binary must fail loudly here, not later
# when it is gating a Claude hook (or, once ported, ssh signing).
if ! "${bin_dir}/dot" --version 1> /dev/null 2>&1; then
	error "built dot binary failed to execute; aborting"
	exit 1
fi

# One symlink per applet (the Go `dot applets` list is the source of truth), plus
# a `dot` self-link so the bare command is on PATH via the already-on-PATH
# scripts dir. New applets are linked automatically on the next build.
ln -sfn "${bin_dir}/dot" "${scripts_dir}/dot"
"${bin_dir}/dot" applets | while IFS= read -r applet; do
	if [ -n "${applet}" ]; then
		ln -sfn "${bin_dir}/dot" "${scripts_dir}/${applet}"
	fi
done

success "dot ${version} installed"
