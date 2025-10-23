#! /usr/bin/env sh

# ensure that brew is configured
if command -v brew 1>/dev/null 2>&1; then
	eval "$(brew shellenv)"
elif [ -x '/opt/homebrew/bin/brew' ]; then
	eval "$(/opt/homebrew/bin/brew shellenv)"
fi

if [ -n "${HOMEBREW_PREFIX:-}" ]; then
	site_functions="${HOMEBREW_PREFIX}/share/zsh/site-functions"

	if [ "${FPATH#*"${site_functions}"}" = "${FPATH}" ]; then
		fpath=("${site_functions}" "${fpath[@]}")
	fi
fi

# load completion
autoload -U +X bashcompinit
bashcompinit

if ! (( $+functions[compdef] )); then
	autoload -U +X compinit
	compinit
fi

# set the zinit home
if [ -n "${HOMEBREW_PREFIX:-}" ]; then
	ZINIT_HOME="${HOMEBREW_PREFIX}/opt/zinit"
else
	ZINIT_HOME="${XDG_DATA_HOME:-${HOME}/.local/share}/zinit"
fi

# load zinit
source "${ZINIT_HOME}/zinit.zsh"

# add some zinit plugins
zinit light zsh-users/zsh-syntax-highlighting
zinit light zsh-users/zsh-autosuggestions
zinit light Aloxaf/fzf-tab

# add in some snippets
zinit snippet OMZP::sudo
zinit snippet OMZP::command-not-found

zinit cdreplay -q

# Keybindings
bindkey -e
bindkey '^p' history-search-backward
bindkey '^n' history-search-forward
bindkey '^[w' kill-region

# History
HISTSIZE=5000
HISTFILE=~/.zsh_history
SAVEHIST=$HISTSIZE
HISTDUP=erase
setopt appendhistory
setopt sharehistory
setopt hist_ignore_space
setopt hist_ignore_all_dups
setopt hist_save_no_dups
setopt hist_ignore_dups
setopt hist_find_no_dups

# Completion styling
zstyle ':completion:*' matcher-list 'm:{a-z}={A-Za-z}'
zstyle ':completion:*' list-colors "${(s.:.)LS_COLORS}"
zstyle ':completion:*' menu no
zstyle ':fzf-tab:complete:cd:*' fzf-preview 'ls --color $realpath'
zstyle ':fzf-tab:complete:__zoxide_z:*' fzf-preview 'ls --color $realpath'

# add aliases
alias ls='ls --color'
alias la='ls -la'
alias ll='ls -l'
alias sts='start-tmux-session'
alias vim='nvim'
alias vi='nvim'

export EDITOR=nvim
export POSH_THEME="${HOME}/.config/oh-my-posh/flags.toml"

# enable shell integrations
eval "$(fzf --zsh)"
eval "$(zoxide init --cmd cd zsh)"
eval "$(direnv hook zsh)"

if [ "$TERM_PROGRAM" != "Apple_Terminal" ]; then
	eval "$(oh-my-posh init zsh)"
fi

SCRIPTS_DIR="${XDG_CONFIG_HOME:-${HOME}/.config}/scripts"

if [ -d "${SCRIPTS_DIR:-}" ]; then
	export PATH="${PATH}:${SCRIPTS_DIR}"
fi
