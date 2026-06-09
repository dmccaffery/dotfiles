#!/usr/bin/env sh
set -eu

# Strip Codex-generated runtime state from a tracked config.toml.
#
# Codex appends non-configuration state to $CODEX_HOME/config.toml as it runs —
# per-project trust ([projects], [projects."/abs/path"]), hook bookkeeping
# ([hooks.state]), model-availability NUX counters ([tui.model_availability_nux]),
# and whatever future versions decide to persist. None of that belongs in a
# public dotfiles repo.
#
# Rather than chase each new state section with a denylist (which silently misses
# anything new), this works the other way round: it keeps an ALLOWLIST of the
# real configuration tables and drops every other table outright. Any state form
# Codex invents later lands in a table that isn't on the list, so it is scrubbed
# without a script change. The flip side: when you add a genuinely new config
# table to config.toml, add its dotted header here too or scrub will delete it.
#
# Everything before the first table header (the `#:schema` line, the comment
# banner, and the top-level keys approval_policy / sandbox_mode) is always kept.
# Comments and blank lines are attributed to the table that follows them, so the
# banner above a dropped state table is dropped with it.

repo_root=$(CDPATH='' cd -- "$(dirname -- "$0")/.." && pwd)
config_file=${1:-"$repo_root/stow/.config/codex/config.toml"}
tmp_file=$(mktemp "$TMPDIR/codex-config.XXXXXX")
trap 'rm -f "$tmp_file"' EXIT HUP INT TERM

# Allowed configuration tables, one dotted header per line (no brackets). Both
# [table] and [[array.of.tables]] headers are matched by their dotted name.
allowed='
sandbox_workspace_write
tui
hooks.Stop
hooks.Stop.hooks
hooks.PermissionRequest
hooks.PermissionRequest.hooks
hooks.PostToolUse
hooks.PostToolUse.hooks
hooks.UserPromptSubmit
hooks.UserPromptSubmit.hooks
hooks.PreToolUse
hooks.PreToolUse.hooks
'
# Collapse to one whitespace-separated line: BSD awk rejects a multi-line -v value.
allowed=$(printf '%s' "$allowed" | tr '\n' ' ')

awk -v allowed="$allowed" '
	BEGIN {
		n = split(allowed, rows, " ")
		for (i = 1; i <= n; i++)
			if (rows[i] != "")
				allow[rows[i]] = 1
		keep = 1   # preamble (before the first table) is configuration
		buf_n = 0  # pending comment/blank lines, flushed once their table is known
	}

	function flush(   i) { for (i = 1; i <= buf_n; i++) print buf[i]; buf_n = 0 }
	function drop() { buf_n = 0 }

	# Table header: [name] or [[name]], optionally trailing whitespace/comment.
	/^[[:space:]]*\[\[?[^]]+\]\]?[[:space:]]*(#.*)?$/ {
		name = $0
		sub(/^[[:space:]]*/, "", name)
		sub(/[[:space:]]*(#.*)?$/, "", name)
		gsub(/^\[+|\]+$/, "", name)
		if (name in allow) { flush(); keep = 1; print } else { drop(); keep = 0 }
		next
	}

	# Comment or blank line: buffer it so it travels with the table it precedes.
	/^[[:space:]]*$/ || /^[[:space:]]*#/ {
		if (keep) buf[++buf_n] = $0
		next
	}

	# Any other line (a key = value, or a multi-line value continuation).
	{ if (keep) { flush(); print } }

	# Anything still buffered at EOF is trailing blank/comment lines with no
	# config line after them — drop them so the file ends on its last real line
	# (a single trailing newline), no matter how many blanks Codex left before
	# the state it appended.
	END { }
' "$config_file" >"$tmp_file"

mv "$tmp_file" "$config_file"
