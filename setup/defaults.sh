#! /usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SETUP_DIR="${SCRIPT_DIR}"

PLATFORM_DIR=$(uname -s | tr '[:upper:]' '[:lower:]')
PLATFORM_SCRIPT="${SETUP_DIR}/${PLATFORM_DIR}/defaults.sh"

if [ -x "${PLATFORM_SCRIPT}" ]; then
	"${PLATFORM_SCRIPT}"
fi
