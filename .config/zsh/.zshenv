export XDG_CONFIG_HOME="${HOME}/.config"
export XDG_CACHE_HOME="${HOME}/.cache"
export XDG_DATA_HOME="${HOME}/.local/share"
export XDG_STATE_HOME="${HOME}/.local/state"
export XDG_RUNTIME_DIR="${HOME}/.local/runtime"

export REPOS_DIR="${HOME}/Repos"

export TMUX_SOCK="${TMUX%%,*}"

export HOMEBREW_BUNDLE_FILE_GLOBAL="${XDG_DATA_HOME}/homebrew/Brewfile"
export HOMEBREW_BUNDLE_INSTALL_CLEANUP=1

export POSH_THEME="${XDG_CONFIG_HOME}/oh-my-posh/prompt.yaml"
export VIVID_THEME="${XDG_CONFIG_HOME}/vivid/themes/cyberdream.yaml"

export EDITOR=nvim
if [ "${TERM_PROGRAM}" = "vscode" ]; then
	export EDITOR='code --wait'
fi

SCRIPTS_DIR="${XDG_DATA_HOME}/scripts"
if [ -d "${SCRIPTS_DIR:-}" ]; then
	export PATH="${PATH}:${SCRIPTS_DIR}"
	export SSH_ASKPASS="${SCRIPTS_DIR}/ssh-askpass"
	export SSH_ASKPASS_REQUIRE=force
fi
