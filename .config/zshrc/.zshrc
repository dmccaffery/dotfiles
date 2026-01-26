#! /usr/bin/env zsh

# profile zsh if enabled
if [ -n "${ZSHPROFILE:-}" ]; then
	zmodload zsh/zprof
fi

# set xdg config home (so lazygit picks it up)
export XDG_CONFIG_HOME="$HOME/.config"
export XDG_CACHE_HOME="$HOME/.cache"
export XDG_DATA_HOME="$HOME/.local/share"
export XDG_STATE_HOME="$HOME/.local/state"
export XDG_RUNTIME_DIR="$HOME/.local/runtime"

# set default editor
export EDITOR=nvim

# ensure that brew is configured
if command -v brew 1>/dev/null 2>&1; then
	eval "$(brew shellenv)"
elif [ -x '/opt/homebrew/bin/brew' ]; then
	eval "$(/opt/homebrew/bin/brew shellenv)"
fi

if [ -n "${HOMEBREW_PREFIX:-}" ]; then
	google_completion="${HOMEBREW_PREFIX:-}/share/zsh/site-functions/_google_cloud_sdk"

	if [ -f "${google_completion:-}" ]; then
		source "${google_completion:-}"
	fi
fi

if ! (( $+functions[compdef] )); then
	autoload -U +X compinit && compinit
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
# zinit light spaceship-prompt/spaceship-prompt
zinit light Aloxaf/fzf-tab

# add in some snippets
zinit snippet OMZP::sudo
zinit snippet OMZP::command-not-found

zinit cdreplay -q

# setup some key bindings
bindkey -e
bindkey '^p' history-search-backward
bindkey '^n' history-search-forward
bindkey '^[w' kill-region

# setup history
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

# completion styling
zstyle ':completion:*' matcher-list 'm:{a-z}={A-Za-z}'
zstyle ':completion:*' list-colors "${(s.:.)LS_COLORS}"
zstyle ':completion:*' menu no
zstyle ':fzf-tab:complete:cd:*' fzf-preview 'ls $realpath'
zstyle ':fzf-tab:complete:__zoxide_z:*' fzf-preview 'ls $realpath'

# add aliases
alias ls='lsd'
alias la='ls -a'
alias lla='ls -la'
alias ll='ls -l'
alias lt='ls --tree'

alias sts='start-tmux-session'
alias ats='attach-tmux-session'

alias vim='nvim'
alias vi='nvim'

# enable shell integrations
eval "$(fzf --zsh)"
eval "$(zoxide init --cmd cd zsh)"
eval "$(direnv hook zsh)"
eval "$(fnm env --use-on-cd)"

# setup the prompt
POSH_THEME="${HOME}/.config/oh-my-posh/flags.toml"
eval "$(oh-my-posh init zsh --config "${POSH_THEME}")"

if [ "${TERM_PROGRAM}" = "vscode" ]; then
	export EDITOR='code --wait'
fi

# add scripts to path
SCRIPTS_DIR="${XDG_DATA_HOME}/scripts"
if [ -d "${SCRIPTS_DIR:-}" ]; then
	export PATH="${PATH}:${SCRIPTS_DIR}"
fi

# start tmux if not already running
if [ -z "${TMUX:-}" ]; then
	tmux start-server 1>/dev/null 2>&1
fi

# detect gcloud path
if command -v 'gcloud' 1>/dev/null 2>&1; then
	GOOGLE_BIN_PATH=$(dirname $(readlink -f $(command -v 'gcloud')))
	export PATH="${GOOGLE_BIN_PATH}:${PATH}"
fi

# end zsh profiling if enabled
if [ -n "${ZSHPROFILE:-}" ]; then
	zprof
	printf '\n\nTIMINGS:\n\n'
fi
