#! /usr/bin/env sh

set -eu

SCRIPT_DIR=$(dirname "$(readlink -f -- "$0")")
SETUP_DIR=$(realpath "${SCRIPT_DIR}/../")

# shellcheck source=./../printing.sh
. "${SETUP_DIR}/printing.sh"

info "setting up homebrew..."
brew_cmd=$(command -v brew 2>&1)

if [ -z "${brew_cmd:-}" ]; then
	brew_cmd=/opt/homebrew/bin/brew
fi

if [ ! -x "${brew_cmd:-}" ]; then
	warn "brew could not be found; please ensure homebrew is installed"
	return 0
fi

eval "$(${brew_cmd} shellenv)"

HOMEBREW_BUNDLE_FILE_GLOBAL="${XDG_DATA_HOME}/homebrew/Brewfile"
export HOMEBREW_BUNDLE_FILE_GLOBAL

HOMEBREW_BUNDLE_DIR="${XDG_CONFIG_HOME}/homebrew"

if [ ! -r "${HOMEBREW_BUNDLE_FILE_GLOBAL}" ]; then
	# allow user to pick an existing brewfile, or type a new name to dump current brews into
	selection=$(find "${HOMEBREW_BUNDLE_DIR}/" -mindepth 1 -maxdepth 1 -type f -name 'Brewfile.*' -print |
		awk -F/ '{ name=$NF; sub(/^Brewfile\./, "", name); print name "\t" $0 }' |
		fzf --header='Pick a brew bundle or type a new name' \
			--delimiter='\t' \
			--with-nth=1 \
			--preview='cat {2}' \
			--print-query || true)

	query=$(printf '%s\n' "${selection}" | sed -n '1p')
	chosen=$(printf '%s\n' "${selection}" | sed -n '2p' | cut -f2)

	if [ -n "${chosen:-}" ]; then
		brewfile="${chosen}"
	elif [ -n "${query:-}" ]; then
		brewfile="${HOMEBREW_BUNDLE_DIR}/Brewfile.${query}"
		brew bundle dump --brews --casks --mas --force --file="${brewfile}"
	else
		brewfile=""
	fi

	if [ -n "${brewfile:-}" ]; then
		mkdir -p "$(dirname "${HOMEBREW_BUNDLE_FILE_GLOBAL}")"
		ln -Ffs "${brewfile}" "${HOMEBREW_BUNDLE_FILE_GLOBAL}"
	fi
fi

requirements_file="${SCRIPT_DIR}/Brewfile.requirements"
if [ -r "${HOMEBREW_BUNDLE_FILE_GLOBAL}" ] && [ -r "${requirements_file}" ]; then
	info "merging required packages into brewfile..."
	target_brewfile=$(readlink -f -- "${HOMEBREW_BUNDLE_FILE_GLOBAL}")
	merged=$(mktemp)
	awk '
		function pkg_name(line) {
			if (match(line, /"[^"]*"/)) {
				return substr(line, RSTART + 1, RLENGTH - 2)
			}
			return ""
		}
		BEGIN {
			req_header = "# required packages -- do not edit"
			prof_header = "# profile packages"
			print req_header
			print ""
		}
		NR == FNR {
			name = pkg_name($0)
			if (name != "") {
				req[name] = 1
			}
			print
			next
		}
		FNR == 1 {
			print ""
			print prof_header
			print ""
			skip_next_blank = 1
		}
		{
			if ($0 == req_header || $0 == prof_header) {
				pending_blank = 0
				skip_next_blank = 1
				next
			}
			if ($0 == "") {
				if (skip_next_blank) {
					skip_next_blank = 0
					next
				}
				pending_blank++
				next
			}
			skip_next_blank = 0
			name = pkg_name($0)
			if (name != "" && name in req) {
				next
			}
			for (i = 0; i < pending_blank; i++) {
				print ""
			}
			pending_blank = 0
			print
		}
	' "${requirements_file}" "${target_brewfile}" > "${merged}"
	cat "${merged}" > "${target_brewfile}"
	rm -f "${merged}"
fi

# The requirements install below pins XDG_CONFIG_HOME to SCRIPT_DIR, so brew
# reads its trusted set from this repo-shipped file rather than the user's
# primary ~/.config/homebrew/trust.json. Review *that* file here so the prompt
# matches exactly what the install is about to load.
required_trust_file="${SCRIPT_DIR}/homebrew/trust.json"
if [ -r "${required_trust_file}" ]; then
	trusted_entries=$(grep -oE '"[^"]*/[^"]*"' "${required_trust_file}" | tr -d '"' || true)
	if [ -n "${trusted_entries:-}" ]; then
		warn "the required trust.json trusts these non-official taps/formulae/casks; their code loads during install:"
		printf '%s\n' "${trusted_entries}" | sed 's/^/      - /'

		if [ -r /dev/tty ]; then
			warn "review the entries above; proceed with install? [y/N]"
			read -r reply < /dev/tty || reply=""
			case "${reply}" in
			[Yy] | [Yy][Ee][Ss]) ;;
			*)
				error "aborted: required trust not confirmed"
				exit 1
				;;
			esac
		else
			warn "no tty available; skipping required-trust confirmation"
		fi
	fi
fi

info "installing / upgrading from brewfile..."
XDG_CONFIG_HOME="${SCRIPT_DIR}" brew bundle install --global
brew reinstall stow 1> /dev/null 2>&1

# Fold the required trust entries into the user's primary trust.json so later
# `brew bundle` / `brewfile` runs — which read ~/.config/homebrew/trust.json —
# already trust everything the requirements install pulled in. jq is on PATH now
# that the install (which provides it via Brewfile.requirements) has finished.
primary_trust_file="${HOMEBREW_BUNDLE_DIR}/trust.json"
if [ -r "${required_trust_file}" ] && command -v jq > /dev/null 2>&1; then
	info "merging required trust entries into ${primary_trust_file}..."
	mkdir -p "$(dirname "${primary_trust_file}")"
	[ -e "${primary_trust_file}" ] || printf '{}\n' > "${primary_trust_file}"

	merged_trust=$(mktemp)
	if jq -s '
		(.[0] // {}) as $primary
		| (.[1] // {}) as $required
		| reduce (($primary + $required) | keys_unsorted[]) as $key
			({}; .[$key] = ((($primary[$key] // []) + ($required[$key] // [])) | unique))
	' "${primary_trust_file}" "${required_trust_file}" > "${merged_trust}"; then
		cat "${merged_trust}" > "${primary_trust_file}"
		chmod 600 "${primary_trust_file}"
	else
		warn "could not merge required trust into ${primary_trust_file}; left unchanged"
	fi
	rm -f "${merged_trust}"
fi

info "cleaning up brew services"
brew services cleanup --quiet

info "cleaning up brew cache"
brew cleanup --quiet
