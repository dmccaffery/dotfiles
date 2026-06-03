#!/usr/bin/env sh
# Re-record docs/assets/demo.cast at a compact, roughly-square Ghostty geometry.
#
# The published cast plays on the docs homepage via asciinema-player. This helper
# captures a fresh recording at a grid that renders ~800x800 px in Ghostty with
# Iosevka NF 15 (cell ~= 7.5 x 19.5 pt, window-padding 5) -> 105 cols x 40 rows.
# asciinema's --window-size forces that grid regardless of the physical window, and
# bakes it into the cast header, so the embed keeps its square shape.
#
# The recording runs in a *new* Ghostty window (opened with `open -na Ghostty.app`,
# the only supported way to launch it from the CLI on macOS). This window — the one you
# ran the script from — becomes a teleprompter: it walks the demo beats one at a time,
# waiting for Enter between each. Command steps are copied to the clipboard so you can
# paste them into the recording window with Cmd+V; quoted steps are directions to perform
# there and are not copied. When you exit the recorded shell, control returns here to
# publish the cast.
#
# When you publish a fresh take, the script also regenerates the README/docs poster
# (docs/assets/images/demo-poster.png) from the new cast with agg, so the static thumbnail
# stays in sync with the recording. The poster frame is the instant you stepped from step 8
# to step 9 in the teleprompter (the post-commit view): the script converts that wall-clock
# moment to the player's idle-collapsed npt timeline, renders it with agg, and rewrites
# poster: "npt:…" in asciinema-player-init.js to match. agg uses the same Iosevka Nerd Font as
# the terminal (see https://docs.asciinema.org/manual/agg/usage/#nerd-fonts) so the powerline /
# oh-my-posh glyphs render instead of tofu.
#
# Override the grid with: WINDOW_SIZE=120x45 ./hack/record-demo.sh
set -eu

WINDOW_SIZE="${WINDOW_SIZE:-105x40}"
# Pin the recording window's theme. Ghostty's config resolves the theme from macOS appearance
# (theme = dark:cyberdream,light:cyberdream-light), but a window spawned via `open -na … -e sh`
# queries the effective appearance before AppKit finishes initialising it and falls back to the
# light branch — asciinema then bakes cyberdream-light (white bg) into the cast header even in
# dark mode. Forcing --theme below removes that dependency so recordings are always cyberdream.
GHOSTTY_THEME="${GHOSTTY_THEME:-cyberdream}"

repo_root=$(cd "$(dirname "$0")/.." && pwd)
dest="$repo_root/docs/assets/demo.cast"
poster="$repo_root/docs/assets/images/demo-poster.png"
tmp=$(mktemp)

# Cyberdream palette for agg, as bg,fg,<16 ANSI colors> (comma-separated, no '#').
# Keep in sync with .config/ghostty/themes/cyberdream and docs/assets/extras.css.
AGG_THEME="16181a,ffffff,16181a,ff6e5e,5eff6c,f1ff5e,5ea1ff,bd5eff,5ef1ff,ffffff,3c4048,ff6e5e,5eff6c,f1ff5e,5ea1ff,bd5eff,5ef1ff,ffffff"
# Terminal font: "Iosevka NF" (matches font-family in .config/ghostty/config), rendered by agg
# from the installed Nerd Font in FONT_DIR so the poster matches the player.
FONT_DIR="${FONT_DIR:-$HOME/Library/Fonts}"
FONT_FAMILY="${FONT_FAMILY:-Iosevka NF}"
# Poster frame timestamp in seconds, on the player's idle-collapsed npt timeline. Normally
# derived from the step 8 -> 9 transition during recording (see update_poster_timestamp); this
# is the fallback used when that can't be computed.
POSTER_NPT="${POSTER_NPT:-120}"
player_init="$repo_root/docs/assets/asciinema-player-init.js"

if ! command -v asciinema > /dev/null 2>&1; then
	echo "error: asciinema not found (brew install asciinema)" >&2
	exit 1
fi

# The new-window recorder and the clipboard teleprompter are macOS-only.
for tool in open pbcopy; do
	if ! command -v "$tool" > /dev/null 2>&1; then
		echo "error: $tool not found — this script targets macOS" >&2
		exit 1
	fi
done

# Regenerate the static poster from a published cast. Renders the whole cast at the
# player's idle cap so the GIF timeline matches npt, grabs the frame shown at the poster
# timestamp, and flattens it onto the cyberdream background. Skips (with a warning) if the
# rendering tools are missing rather than failing the recording.
regenerate_poster() {
	cast="$1"
	for tool in agg ffmpeg magick; do
		if ! command -v "$tool" > /dev/null 2>&1; then
			printf 'warning: %s not found; skipping poster regen — update %s by hand.\n' \
				"$tool" "$poster" >&2
			return 0
		fi
	done
	workdir=$(mktemp -d)
	printf 'Rendering poster from %s with %s (idle cap 2s)…\n' "$cast" "$FONT_FAMILY"
	agg --idle-time-limit 2 --theme "$AGG_THEME" \
		--font-dir "$FONT_DIR" --font-family "$FONT_FAMILY" \
		"$cast" "$workdir/demo.gif"
	ffmpeg -v error -ss "$POSTER_NPT" -i "$workdir/demo.gif" -frames:v 1 -y "$workdir/frame.png"
	magick "$workdir/frame.png" -background "#16181a" -flatten -define png:color-type=2 "$poster"
	rm -rf "$workdir"
	printf 'Updated %s\n' "$poster"
}

# Convert a wall-clock instant (epoch seconds) captured during recording into the player's
# idle-collapsed npt timeline, set POSTER_NPT to it, and rewrite poster: "npt:…" in the player
# init so the live poster and the regenerated PNG agree. The cast header's `timestamp` is the
# recording start (npt 0), and each event's interval is real time since the previous one —
# capped at the idle limit, exactly as the player and agg collapse idle gaps. Leaves the
# fallback POSTER_NPT in place (with a warning) if python3 is missing or the value can't be
# derived.
update_poster_timestamp() {
	cast="$1"
	mark="$2"
	if ! command -v python3 > /dev/null 2>&1; then
		printf 'warning: python3 not found; keeping the default poster frame.\n' >&2
		return 0
	fi
	npt=$(python3 -c '
import json, sys

cast, mark, cap = sys.argv[1], float(sys.argv[2]), 2.0
with open(cast) as f:
    t0 = float(json.loads(f.readline()).get("timestamp") or 0)
    if t0 <= 0:
        print(-1)
        sys.exit(0)
    target = mark - t0
    real = npt = 0.0
    for line in f:
        line = line.strip()
        if not line:
            continue
        d = float(json.loads(line)[0])
        if real + d >= target:
            npt += min(target - real, cap)
            break
        real += d
        npt += min(d, cap)
print(max(0, int(round(npt))))
' "$cast" "$mark" 2> /dev/null || true)
	case "$npt" in
	'' | *[!0-9]*)
		printf 'warning: could not derive a poster timestamp from the recording; keeping the existing frame.\n' >&2
		return 0
		;;
	esac
	POSTER_NPT="$npt"
	label=$(printf '%d:%02d' "$((npt / 60))" "$((npt % 60))")
	if [ -f "$player_init" ]; then
		sed -i '' -E 's/(poster: "npt:)[0-9:.]+(")/\1'"$label"'\2/' "$player_init"
		printf 'Poster frame set to npt:%s (%ss); updated %s\n' "$label" "$npt" "$player_init"
	fi
}

# Teleprompter helpers. Each step waits for Enter; command steps are copied to the
# clipboard with pbcopy, quoted instruction steps are printed as directions only.
copy_step() {
	printf '\n  [%s] %s\n' "$1" "$2"
	printf '      ↳ Enter to copy to the clipboard… '
	read -r _
	printf '%s' "$2" | pbcopy
	printf '      ✓ copied — paste into the recording window with Cmd+V, then run it\n'
}

note_step() {
	printf '\n  [%s] (do this in the recording window)\n\n %s\n' "$1" "$2"
}

teleprompter() {
	cat << 'EOF'

Teleprompter — a new Ghostty window is now recording.

Press Enter here for each step. Command steps land on your clipboard; paste them
into the recording window with Cmd+V and run them. Quoted steps are directions to
perform in that window. The window is pinned to the cyberdream theme; the frame on
screen as you cross from step 8 to step 9 becomes the README/docs poster.
EOF
	copy_step 1 'mkdir demo && cd demo'
	copy_step 2 'git init'
	copy_step 3 'git remote add origin https://github.com/dmccaffery/demo'
	copy_step 4 'tmux-session start .'
	note_step 5 'Navigate to each of the windows within the tmux session; returning to either claude or opencode.'
	copy_step 6 'Create a golang cli tool that will tell jokes.'
	note_step 7 'Return to the neovim window and open lazygit to show main.go.'
	copy_step 8 './commit.sh'
	# The frame showing here — after commit.sh, before the exit — is the poster.
	poster_mark=$(date +%s)
	note_step 9 'Exit the session and the recording.'
}

# Launch the recorder in a new Ghostty window. `open -na Ghostty.app` is the only
# supported way to start Ghostty from the CLI on macOS. The window runs a tiny script
# that records to $tmp and then touches $flag, so this (teleprompter) window can tell
# when the recording has finished. Sizing the window to the grid keeps the visible
# window in step with the geometry asciinema bakes into the cast.
flag="$tmp.done"
rec_script="$tmp.rec.sh"
poster_mark=0 # set by the teleprompter at the step 8 -> 9 boundary
rm -f "$flag"
trap 'rm -f "$rec_script" "$flag"' EXIT

cat > "$rec_script" << EOF
#!/usr/bin/env sh
cd "$repo_root/.." || exit 1
asciinema rec --window-size "$WINDOW_SIZE" --idle-time-limit 2 --title "dotfiles demo" "$tmp"
: > "$flag"
EOF

cols=${WINDOW_SIZE%x*}
rows=${WINDOW_SIZE#*x}
printf 'Opening a new Ghostty window to record at %s -> %s\n' "$WINDOW_SIZE" "$tmp"
open -na Ghostty.app --args --theme="$GHOSTTY_THEME" --window-width="$cols" --window-height="$rows" -e sh "$rec_script"

teleprompter

printf '\nWaiting for the recording to finish (exit the recorded shell in the other window)…\n'
while [ ! -e "$flag" ]; do
	sleep 1
done

printf '\nRecorded %s\n' "$tmp"
printf 'Replace the published cast at %s? [y/N] ' "$dest"
read -r ans
case "$ans" in
y | Y)
	mv "$tmp" "$dest"
	printf 'Updated %s\n' "$dest"
	update_poster_timestamp "$dest" "$poster_mark"
	regenerate_poster "$dest"
	;;
*)
	printf 'Left the new recording at %s (not published).\n' "$tmp"
	;;
esac
