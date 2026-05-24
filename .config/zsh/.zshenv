export XDG_CONFIG_HOME="${HOME}/.config"
export XDG_CACHE_HOME="${HOME}/.cache"
export XDG_DATA_HOME="${HOME}/.local/share"
export XDG_STATE_HOME="${HOME}/.local/state"
export XDG_RUNTIME_DIR="${HOME}/.local/runtime"

export REPO_DIR="${HOME}/Repos"

export CLAUDE_CODE_NO_FLICKER=1

export HOMEBREW_BUNDLE_FILE="${XDG_DATA_HOME}/homebrew/Brewfile"
export HOMEBREW_BUNDLE_INSTALL_CLEANUP=1

export POSH_THEME="${XDG_CONFIG_HOME}/oh-my-posh/prompt.yaml"
export VIVID_THEME="${XDG_CONFIG_HOME}/vivid/themes/cyberdream.yaml"

SCRIPTS_DIR="${XDG_DATA_HOME}/scripts"
if [ -d "${SCRIPTS_DIR:-}" ]; then
	export PATH="${PATH}:${SCRIPTS_DIR}"
fi
