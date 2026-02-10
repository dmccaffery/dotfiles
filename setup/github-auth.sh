#!/usr/bin/env bash

set -euo pipefail

SETUP_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
. "${SETUP_DIR}/printing.sh"

info "checking github auth..."

scopes=$(gh auth status --json 'hosts' --jq '.hosts."github.com".[].scopes' 2> /dev/null)

if ! gh auth status 1> /dev/null 2>&1; then
	warning "not logged into github, will login now..."
	gh auth login \
		--git-protocol https \
		--hostname github.com \
		--scopes gist,workflow,repo,user,read:org,read:public_key,read:ssh_signing_key \
		--web \
		--clipboard

	return 0
fi

scopes=$(gh auth status --json 'hosts' --jq '.hosts."github.com".[].scopes' 2> /dev/null)

if echo "${scopes:-}" | grep -F -q "read:ssh_signing_key" 1> /dev/null 2>&1 &&
	echo "${scopes:-}" | grep -F -q "read:public_key" 1> /dev/null 2>&1; then
	info "already logged into github with the appropriate scopes"
	return 0
fi

warn "not logged into github with the read:public_key and read:ssh_signing_key scopes..."
gh auth refresh \
	--hostname github.com \
	--scopes repo,user,read:org,read:public_key,read:ssh_signing_key \
	--clipboard
