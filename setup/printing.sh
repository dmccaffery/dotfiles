#! /usr/bin/env sh

set -eu

default_color=$(tput sgr 0)

blue="$(tput setaf 4)"
info() {
	printf "%s==> %s%s\n" "$blue" "$1" "$default_color"
}

yellow="$(tput setaf 3)"
warn() {
	printf "%s==> %s%s\n" "$yellow" "$1" "$default_color" >&2
}

green="$(tput setaf 2)"
success() {
	printf "%s==> %s%s\n" "$green" "$1" "$default_color"
}

red="$(tput setaf 1)"
error() {
	printf "%s==> %s%s\n" "$red" "$1" "$default_color" >&2
}
