#! /usr/bin/env sh

set -eu

tty=/dev/tty

default_color=$(tput sgr 0)

blue="$(tput setaf 4)"
info() {
	printf "%s==> %s%s\n" "$blue" "$1" "$default_color" > "${tty}" 2> /dev/null || true
}

green="$(tput setaf 2)"
success() {
	printf "%s==> %s%s\n" "$green" "$1" "$default_color" > "${tty}" 2> /dev/null || true
}

yellow="$(tput setaf 3)"
warn() {
	printf "%s==> %s%s\n" "$yellow" "$1" "$default_color" > "${tty}" 2> /dev/null || true
}

input=$(cat)
worktree_path=$(echo "${input:-}" | jq -r '.worktree_path')

if [ ! -d "$worktree_path" ]; then
	warn "worktree at: ${worktree_path} no longer exists; exiting gracefully..."
	exit 0
fi

branch=$(git -C "$worktree_path" rev-parse --abbrev-ref HEAD 2> /dev/null || true)

uncommitted=$(git -C "$worktree_path" status --porcelain 2> /dev/null | wc -l | tr -d ' ')
unpushed=$(git -C "$worktree_path" rev-list --count HEAD --not --remotes 2> /dev/null || echo 0)
if [ "${uncommitted}" != "0" ] || [ "${unpushed}" != "0" ]; then
	warn "${worktree_path} [${branch:-?}] — uncommitted: ${uncommitted}, unpushed: ${unpushed}"
fi

git worktree remove "$worktree_path" --force 2> /dev/null || true
if [ -n "$branch" ] && [ "${branch}" != "${branch##agent/*}" ]; then
	git branch -D "$branch" 2> /dev/null || true
fi
