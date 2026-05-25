#! /usr/bin/env sh

set -eu

SCRIPT_DIR=$(dirname "$(readlink -f -- "$0")")
. "${SCRIPT_DIR}/setup/printing.sh"

BACKUPS_DIR="${SCRIPT_DIR}/backups"
query="${1:-}"

if ! hash fzf 1> /dev/null 2>&1; then
	error "fzf could not be found; install via 'brew install fzf' first"
	exit 1
fi

if [ ! -d "${BACKUPS_DIR}" ]; then
	error "no backups directory at ${BACKUPS_DIR}"
	exit 1
fi

info "select a backup to restore"
chosen=$(find "${BACKUPS_DIR}" -mindepth 1 -maxdepth 1 -type d -print |
	sort -r |
	awk -F/ '{ print $NF "\t" $0 }' |
	fzf --header='Pick a backup to restore (Ctrl-C to cancel)' \
		--delimiter='\t' \
		--with-nth=1 \
		--select-1 \
		--exit-0 \
		--query="${query}" \
		--preview='find {2} -mindepth 1 \( -type f -o -type l \) | sed "s|{2}/||" | sort' || true)

if [ -z "${chosen:-}" ]; then
	warn "no backup selected; aborting"
	exit 0
fi

backup_name=$(printf '%s\n' "${chosen}" | cut -f1)
backup_path=$(printf '%s\n' "${chosen}" | cut -f2)

info "selected backup: ${backup_name}"
warn "this will run 'stow -D .' and copy ${backup_name} back into \$HOME"
printf '==> continue? [y/N] '
read -r reply
case "${reply:-}" in
y | Y | yes | YES) ;;
*)
	warn "aborted"
	exit 0
	;;
esac

if hash stow 1> /dev/null 2>&1; then
	info "removing existing stow symlinks..."
	(cd "${SCRIPT_DIR}" && stow -D .)
else
	warn "stow could not be found; skipping unstow (collision check will catch leftovers)"
fi

file_list="${TMPDIR:-/tmp}/restore-files.$$"
collision_log="${TMPDIR:-/tmp}/restore-collisions.$$"
trap 'rm -f "${file_list}" "${collision_log}"' EXIT INT TERM

find "${backup_path}" -mindepth 1 \( -type f -o -type l \) -print > "${file_list}"

info "checking for collisions in \$HOME..."
: > "${collision_log}"
while IFS= read -r src; do
	rel="${src#"${backup_path}"/}"
	dst="${HOME}/${rel}"

	if [ -e "${dst}" ] || [ -L "${dst}" ]; then
		printf '%s\n' "${dst}" >> "${collision_log}"
	fi
done < "${file_list}"

if [ -s "${collision_log}" ]; then
	error "collisions detected in \$HOME; aborting restore:"
	while IFS= read -r line; do
		warn "  ${line}"
	done < "${collision_log}"
	exit 1
fi

info "copying ${backup_name} into \$HOME..."
while IFS= read -r src; do
	rel="${src#"${backup_path}"/}"
	dst="${HOME}/${rel}"
	dst_dir=$(dirname "${dst}")

	mkdir -p "${dst_dir}"
	cp -Pp "${src}" "${dst}"
done < "${file_list}"

success "restore complete"

printf '==> delete backup %s? [y/N] ' "${backup_name}"
read -r reply
case "${reply:-}" in
y | Y | yes | YES)
	rm -rf "${backup_path}"
	success "deleted ${backup_name}"
	;;
*)
	info "backup preserved at ${backup_path}"
	;;
esac
