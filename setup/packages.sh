#! /usr/bin/env sh

set -eu

# setup buildx
if command -v docker-buildx 1> /dev/null 2>&1; then
	mkdir -p "${HOME}/.docker/cli-plugins"
	ln -fns "$(command -v docker-buildx 2> /dev/null)" "${HOME}/.docker/cli-plugins"
fi
