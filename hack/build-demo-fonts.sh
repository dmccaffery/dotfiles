#!/usr/bin/env sh
# Rebuild the self-hosted Iosevka Nerd Font web fonts used by the asciinema-player demo.
#
# asciinema-player renders the homepage recording as HTML, so the browser needs a Nerd Font
# web font or the powerline / oh-my-posh glyphs fall back to tofu boxes. The full Iosevka Nerd
# Font is ~13 MB per face, so this script subsets each face down to the glyphs the demo renders
# (plus the common box-drawing, block, dingbat and powerline ranges) and writes ~50 KB woff2
# files to docs/assets/fonts/, which extras.css then @font-face's onto the player terminal.
#
# Re-run this after re-recording the demo (hack/record-demo.sh) if the new take introduces
# glyphs outside the ranges below. Requires uv (brew install uv) and python3.
#
# Override the source font directory with: FONT_DIR=/path/to/fonts ./hack/build-demo-fonts.sh
set -eu

repo_root=$(cd "$(dirname "$0")/.." && pwd)
cast="$repo_root/docs/assets/demo.cast"
dest="$repo_root/docs/assets/fonts"
FONT_DIR="${FONT_DIR:-$HOME/Library/Fonts}"

# Common terminal ranges kept so minor re-records don't immediately re-break: Latin-1,
# general punctuation, arrows, misc technical, box drawing, block elements, geometric shapes,
# misc symbols, dingbats, and Powerline. The --text-file below adds whatever else the cast uses.
RANGES="U+0000-00FF,U+2000-206F,U+2190-21FF,U+2300-23FF,U+2500-257F,U+2580-259F,U+25A0-25FF,U+2600-26FF,U+2700-27BF,U+E0A0-E0D7,U+E0B0-E0BF"

if ! command -v uv > /dev/null 2>&1; then
	echo "error: uv not found (brew install uv)" >&2
	exit 1
fi

mkdir -p "$dest"
text=$(mktemp)
trap 'rm -f "$text"' EXIT

# Pull every output chunk from the cast so the subset covers exactly the glyphs it renders.
python3 - "$cast" > "$text" << 'PY'
import json, sys

with open(sys.argv[1]) as f:
    next(f)  # skip the asciicast header
    for line in f:
        line = line.strip()
        if not line:
            continue
        ev = json.loads(line)
        if len(ev) >= 3 and ev[1] == "o":
            sys.stdout.write(ev[2])
PY

# "Iosevka NF" (the non-Mono Nerd Font, matching font-family in .config/ghostty/config).
for face in Regular Bold Italic BoldItalic; do
	src="$FONT_DIR/IosevkaNerdFont-$face.ttf"
	if [ ! -f "$src" ]; then
		echo "error: $src not found (install the Iosevka Nerd Font, or set FONT_DIR)" >&2
		exit 1
	fi
	# fonttools[woff] pulls in the brotli backend needed for the woff2 flavor.
	uvx --from "fonttools[woff]" pyftsubset "$src" \
		"--text-file=$text" "--unicodes=$RANGES" \
		--layout-features='' --no-hinting --desubroutinize \
		--flavor=woff2 \
		"--output-file=$dest/IosevkaNerdFont-$face.woff2"
	printf 'wrote %s\n' "$dest/IosevkaNerdFont-$face.woff2"
done
