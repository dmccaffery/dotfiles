#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
. "${SCRIPT_DIR}/printing.sh"

XDG_CONFIG_HOME="$HOME/.config"

if [ ! -d "${XDG_CONFIG_HOME}" ]; then
	info "creating config home at: ${XDG_CONFIG_HOME}"
	mkdir -p "${XDG_CONFIG_HOME}"
fi

XDG_CACHE_HOME="$HOME/.cache"

if [ ! -d "${XDG_CACHE_HOME}" ]; then
	info "creating cache home at: ${XDG_CACHE_HOME}"
	mkdir -p "${XDG_CACHE_HOME}"
fi

XDG_DATA_HOME="$HOME/.local/share"

if [ ! -d "${XDG_DATA_HOME}" ]; then
	info "creating data home at: ${XDG_DATA_HOME}"
	mkdir -p "${XDG_DATA_HOME}"
fi

XDG_STATE_HOME="$HOME/.local/state"

if [ ! -d "${XDG_STATE_HOME}" ]; then
	info "creating state home at: ${XDG_STATE_HOME}"
	mkdir -p "${XDG_STATE_HOME}"
fi

XDG_RUNTIME_DIR="$HOME/.local/runtime"

if [ ! -d "${XDG_RUNTIME_DIR}" ]; then
	info "creating runtime dir at: ${XDG_RUNTIME_DIR}"
	mkdir -p "${XDG_RUNTIME_DIR}"
	chmod u=rwx,go= "${XDG_RUNTIME_DIR}"
fi
