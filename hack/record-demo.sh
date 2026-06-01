#!/usr/bin/env sh
# Re-record docs/assets/demo.cast at a compact, roughly-square Ghostty geometry.
#
# The published cast plays on the docs homepage via asciinema-player. This helper
# captures a fresh recording at a grid that renders ~800x800 px in Ghostty with
# Iosevka NF 15 (cell ~= 7.5 x 19.5 pt, window-padding 5) -> 105 cols x 40 rows.
# asciinema's --window-size forces that grid regardless of the physical window, and
# bakes it into the cast header, so the embed keeps its square shape.
#
# Override the grid with: WINDOW_SIZE=120x45 ./hack/record-demo.sh
set -eu

WINDOW_SIZE="${WINDOW_SIZE:-105x40}"

repo_root=$(cd "$(dirname "$0")/.." && pwd)
dest="$repo_root/docs/assets/demo.cast"
tmp=$(mktemp)

if ! command -v asciinema > /dev/null 2>&1; then
	echo "error: asciinema not found (brew install asciinema)" >&2
	exit 1
fi

cat << 'EOF'
Replay checklist (reproduce the beats; the Claude-agent output will differ):

  1. Scaffold the repo at the oh-my-posh prompt:
       mkdir demo && cd demo
       git init
       git remote add origin https://github.com/dmccaffery/demo
  2. Enter your themed tmux session: 1: nvim (LazyVim), 2: zsh, 3: claude (agent).
  3. NeoVim window loads (LazyVim dashboard).
  4. Claude Code agent (Opus 4.8) builds a small Go "joke-telling CLI"
     (go.mod, main.go, math/rand/v2), then writes .commit.sh and excludes it
     via .git/info/exclude for the first commit.
  5. Show the LazyVim which-key +git menu (the ~2:00 poster frame).
  6. Open Lazygit; stage (space) and commit (c).
  7. Wrap up: git remote -v ; git log --oneline ; exit

Notes:
  - Record inside Ghostty with the cyberdream theme active so the header palette matches.
  - Keep the take >= 2:00 long; asciinema-player-init.js posters at npt:2:00.
  - Idle is baked at 2s (--idle-time-limit 2), matching the player's idleTimeLimit: 2.

EOF

printf 'Recording at %s -> %s\n' "$WINDOW_SIZE" "$tmp"
printf 'Exit the recorded shell to stop.\n\n'

(
	cd "${repo_root}/../"
	asciinema rec --window-size "$WINDOW_SIZE" --idle-time-limit 2 --title "dotfiles demo" "$tmp"
)

printf '\nRecorded %s\n' "$tmp"
printf 'Replace the published cast at %s? [y/N] ' "$dest"
read -r ans
case "$ans" in
y | Y)
	mv "$tmp" "$dest"
	printf 'Updated %s\n' "$dest"
	;;
*)
	printf 'Left the new recording at %s (not published).\n' "$tmp"
	;;
esac
