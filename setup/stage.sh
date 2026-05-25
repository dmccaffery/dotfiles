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

	FOUND=0

	if [ -x "${PLATFORM_BEFORE:-}" ]; then
		"${PLATFORM_BEFORE}"
		FOUND=1
	fi

	if [ -x "${PLATFORM_SCRIPT:-}" ]; then
		"${PLATFORM_SCRIPT}"
		FOUND=1
	elif [ -x "${PLATFORM_DIR}/stage.sh" ]; then
		"${PLATFORM_DIR}/stage.sh"
		FOUND=1
	fi

	if [ -x "${STAGE_SCRIPT:-}" ]; then
		"${STAGE_SCRIPT}"
		FOUND=1
	fi

	if [ -x "${PLATFORM_AFTER:-}" ]; then
		"${PLATFORM_AFTER}"
		FOUND=1
	fi

	if [ "${FOUND}" -eq 0 ]; then
		error "no scripts found for stage: ${STAGE_NAME}"
		return 1
	fi
}

while [ -n "${1:-}" ]; do
	stage "${1}"
	shift
done
