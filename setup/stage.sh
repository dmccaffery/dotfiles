#! /usr/bin/env sh

set -eu

SETUP_DIR=$(dirname "$(readlink -f -- "$0")")
. "${SETUP_DIR}/printing.sh"

PLATFORM=$(uname -s | tr '[:upper:]' '[:lower:]')

stage() {
	STAGE_NAME="${1:-}"

	STAGE_SCRIPT="${SETUP_DIR}/${STAGE_NAME}.sh"

	PLATFORM_DIR="${SETUP_DIR}/${PLATFORM}/${STAGE_NAME}"
	PLATFORM_BEFORE="${PLATFORM_DIR}/before.sh"
	PLATFORM_SCRIPT="${PLATFORM_DIR}.sh"
	PLATFORM_AFTER="${PLATFORM_DIR}/after.sh"

	info "running stage: ${STAGE_NAME}"

	if [ -x "${PLATFORM_BEFORE:-}" ]; then
		"${PLATFORM_BEFORE}"
	fi

	if [ -x "${STAGE_SCRIPT:-}" ]; then
		"${STAGE_SCRIPT}"
	fi

	if [ -x "${PLATFORM_SCRIPT:-}" ]; then
		"${PLATFORM_SCRIPT}"
	elif [ -x "${PLATFORM_DIR}/stage.sh" ]; then
		"${PLATFORM_DIR}/stage.sh"
	fi

	PLATFORM_AFTER="${PLATFORM_DIR}/after.sh"
	if [ -x "${PLATFORM_AFTER:-}" ]; then
		"${PLATFORM_AFTER}"
	fi
}

while [ -n "${1:-}" ]; do
	stage "${1}"
	shift
done
