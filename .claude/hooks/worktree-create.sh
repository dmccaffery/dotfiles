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

sanitize() {
	printf '%s' "$1" | tr -c 'A-Za-z0-9._-' '-'
}

if [ -n "${TMUX:-}" ] && hash tmux > /dev/null 2>&1; then
	info "using tmux session name as worktree name"
	name=$(tmux display-message -p '#S' 2> /dev/null || true)
fi

repo_path="${CLAUDE_PROJECT_DIR:-}"
if [ -z "${repo_path:-}" ]; then
	warn "missing CLAUDE_PROJECT_DIR env-var; falling back to repo top-level path..."
	repo_path=$(git rev-parse --show-toplevel)
fi

if [ -z "${name:-}" ]; then
	warn "not in tmux session; using project dir name as worktree name"
	name=$(basename -- "${repo_path}")
fi

input=$(cat)

if [ -n "${input:-}" ]; then
	info "suffix: using name provided by claude"
	suffix=$(echo "${input}" | jq -r '.name')
fi

if [ -z "${suffix:-}" ]; then
	warn "suffix: no name specified; falling back to current timestamp"
	suffix="$(date -u +%Y%m%d-%H%M%S)"
fi

name=$(sanitize "${name}-${suffix}")
worktree_path="${HOME}/.cache/agent/worktrees/${name}"
branch="agent/${name}"

info "creating worktree ${name} at ${worktree_path}"
mkdir -p "${worktree_path}" > /dev/null 2>&1
git worktree add -b "${branch}" "${worktree_path}" > /dev/null 2>&1

success "worktree ${name} is ready"

echo "${worktree_path}"
