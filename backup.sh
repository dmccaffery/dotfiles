#! /usr/bin/env sh
set -eu

BACKUP_DIR=$(dirname "$(readlink -f -- "$0")")
. "${BACKUP_DIR}/setup/printing.sh"

name="${1:-$(date +'%Y-%M-%d-%H%M%S')}"
backup_path="${BACKUP_DIR}/backups/${name}"

mkdir -p "${backup_path}" 1> /dev/null 2>&1 || true

__backup() {
	name="${1:-}"

	if [ -z "${name:-}" ]; then
		return 0
	fi

	path="${HOME}/${name}"

	if [ ! -e "${path:-}" ]; then
		info "${path} does not exist, skipping..."
		return 0
	fi

	# do not backup if its already a symlink
	if [ -h "${path:-}" ]; then
		info "${path} is already stowed, skipping..."
		return 0
	fi

	# create the directory in the backup and mv just the file
	if [ -f "${path:-}" ]; then
		dir="${backup_path}/$(dirname "${name}")"
		warn "backing up file ${name} to ${dir}..."

		mkdir -p "${dir}"
		mv "${path}" "${dir}"

		return 0
	fi

	if [ -d "${path:-}" ]; then
		warn "backing up directory ${name} to ${backup_path}..."
		mv "${path}" "${backup_path}"

		return 0
	fi

	error "unknown path ${path:-} or directory entry type"
	return 1
}

__backup .claude/settings.json
__backup .claude/themes

(
	# The stowed trees live under stow/ now; cd into it so each discovered entry
	# stays $HOME-relative (.config/<name>) for __backup to resolve under $HOME.
	cd "${BACKUP_DIR}/stow"
	find .config/ -mindepth 1 -maxdepth 1 |
		while read -r config; do
			if [ -z "${config:-}" ]; then
				break
			fi

			__backup "${config}"
		done
)

__backup .local/share/scripts
__backup .local/share/wallpapers

__backup .ssh/rc

__backup .terminfo/67/ghostty
__backup .terminfo/78/xterm-ghostty

__backup Library/LaunchAgents/org.homebrew.ssh-agent.plist

__backup .zshrc
