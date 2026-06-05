---
icon: simple/gnubash
---

# Shell

The shell is Zsh. Configuration lives under `.config/zsh/` (XDG-compliant) and is loaded via
the standard `ZDOTDIR` mechanism. There are two files: `.zshenv` (always sourced) and `.zshrc`
(interactive shells only).

## `.zshenv`

`.zshenv` runs for every Zsh invocation — login, interactive, and scripts. Keep it minimal:

```sh title=".config/zsh/.zshenv"
export XDG_CONFIG_HOME="${HOME}/.config"
export XDG_CACHE_HOME="${HOME}/.cache"
export XDG_DATA_HOME="${HOME}/.local/share"
export XDG_STATE_HOME="${HOME}/.local/state"
export XDG_RUNTIME_DIR="${HOME}/.local/runtime"

export REPOS_DIR="${HOME}/Repos"

export HOMEBREW_BUNDLE_FILE_GLOBAL="${XDG_DATA_HOME}/homebrew/Brewfile"
export HOMEBREW_BUNDLE_FORCE_INSTALL_CLEANUP=1
export HOMEBREW_REQUIRE_TAP_TRUST=1

export GOPATH="${XDG_DATA_HOME}/go"
export GOCACHE="${XDG_CACHE_HOME}/go/build"
export GOMODCACHE="${XDG_CACHE_HOME}/go/mod"
export GOENV="${XDG_CACHE_HOME}/go/env"

export POSH_THEME="${XDG_CONFIG_HOME}/oh-my-posh/prompt.yaml"
export VIVID_THEME="${XDG_CONFIG_HOME}/vivid/themes/cyberdream.yaml"

export CODEX_HOME="${XDG_CONFIG_HOME}/codex"

export EDITOR=nvim
if [ "${TERM_PROGRAM}" = "vscode" ]; then
    export EDITOR='code --wait'
fi

SCRIPTS_DIR="${XDG_DATA_HOME}/scripts"
if [ -d "${SCRIPTS_DIR:-}" ]; then
    export PATH="${PATH}:${SCRIPTS_DIR}"
fi
```

Notable behaviors:

- **All XDG dirs are explicit.** Tools that respect XDG (most modern ones) end up here, keeping
  `$HOME` clean.
- **`SCRIPTS_DIR` is appended to `PATH`.** Anything dropped into
  [`.local/share/scripts/`](../scripts/index.md) is automatically callable.
- **`POSH_THEME` and `VIVID_THEME`** point oh-my-posh and vivid at cyberdream variants.
- **`CODEX_HOME`** relocates [Codex](../codex/index.md)'s home to `~/.config/codex` so its config lives under XDG
  config like every other tool (Codex defaults to `~/.codex`).
- **`HOMEBREW_BUNDLE_FILE_GLOBAL`** points `brew bundle --global` at the merged Brewfile under
  `$XDG_DATA_HOME/homebrew/Brewfile`, which [`setup/darwin/packages.sh`](../getting-started/install.md)
  populates by symlinking the chosen profile from `.config/homebrew/Brewfile.*`.
- **`HOMEBREW_BUNDLE_FORCE_INSTALL_CLEANUP`** and **`HOMEBREW_REQUIRE_TAP_TRUST`** shape every
  `brew bundle` run — auto-cleanup on install and the non-official-tap trust gate. See
  [Brew bundle](brew-bundle.md#environment-variables) and its [trusted-taps](brew-bundle.md#trusted-taps) section.
- **The `GO*` overrides** relocate Go's scattered, non-XDG default locations onto XDG paths:
  `GOPATH` to `$XDG_DATA_HOME/go` (workspace + `go install` binaries) and the build, module, and
  `env` caches under `$XDG_CACHE_HOME/go`. They live here rather than in Claude Code's
  [`settings.json` `env` block](../claude/settings.md#environment) because that block performs no
  variable expansion — a `~/...` or `$HOME/...` value reaches `go` verbatim and fails as a
  relative path. Here `${XDG_DATA_HOME}` / `${XDG_CACHE_HOME}` expand to absolute paths.
- **`EDITOR`** defaults to `nvim` but flips to `code --wait` inside VS Code's integrated terminal
  so `git commit` and friends pop a buffer in the host editor.

## `.zshrc`

`.zshrc` runs for interactive shells only. Key sections:

### Plugin manager (zinit)

```sh
ZINIT_HOME="${HOMEBREW_PREFIX}/opt/zinit"
source "${ZINIT_HOME}/zinit.zsh"

zinit light zsh-users/zsh-syntax-highlighting
zinit light zsh-users/zsh-autosuggestions
zinit light Aloxaf/fzf-tab

zinit snippet OMZP::sudo
zinit snippet OMZP::command-not-found
```

[zinit](https://github.com/zdharma-continuum/zinit) handles plugin lazy-loading. The snippets
pull in Oh My Zsh plugins without the full OMZ overhead.

### History and bindings

| Setting                       | Value                                               |
| ----------------------------- | --------------------------------------------------- |
| `HISTSIZE` / `SAVEHIST`       | 5000 / 5000                                         |
| `setopt sharehistory`         | History is shared across sessions in real time.     |
| `setopt hist_ignore_all_dups` | Duplicate commands are removed from history.        |
| `bindkey '^p'` / `'^n'`       | Prefix-search backwards / forwards through history. |
| `bindkey '^[w'`               | `kill-region` on `Alt-w`.                           |

### Aliases

```sh
alias ls='lsd'
alias la='ls -a'
alias lla='ls -la'
alias ll='ls -l'
alias lt='ls --tree'

alias sts='tmux-session start'
alias ats='tmux-session attach'
alias ets='tmux-session end'
alias lts='tmux list-session'
alias kts='tmux kill-server'

alias vim='nvim'
alias vi='nvim'
alias gh='gh-switch-user'
```

### Shell integrations

```sh
eval "$(oh-my-posh init zsh --config $POSH_THEME)"   # prompt — initialize first

eval "$(fzf --zsh)"                       # fzf keybinds and completion
eval "$(zoxide init --cmd cd zsh)"        # smart `cd` replacement
eval "$(direnv hook zsh)"                 # per-directory env vars
eval "$(fnm env --use-on-cd)"             # node version manager
```

`zoxide` is initialized with `--cmd cd`, so plain `cd` becomes the jump-aware version.
`oh-my-posh` is initialized **before** `zoxide` so that `zoxide`'s `chpwd` hook lands last in the
hook chain — this is what `zoxide doctor` expects and it silences the noisy startup warning that
appeared when the order was reversed.

### Yazi cwd helper

```sh
function y() {
    local tmp="$(mktemp -t "yazi-cwd.XXXXXX")" cwd
    command yazi "$@" --cwd-file="$tmp"
    IFS= read -r -d '' cwd < "$tmp"
    [ "$cwd" != "$PWD" ] && [ -d "$cwd" ] && builtin cd -- "$cwd"
    rm -f -- "$tmp"
}
```

Run `y` instead of `yazi` to have your shell `cd` into the directory you ended in.

### Profiling

```sh
if [ -n "${ZSHPROFILE:-}" ]; then
    zmodload zsh/zprof
fi
# ...
if [ -n "${ZSHPROFILE:-}" ]; then
    zprof
fi
```

`ZSHPROFILE=1 zsh -i -c exit` dumps a startup timing table. Wrapped by the
[`profile-shell`](../scripts/misc.md) helper.

## Where this is loaded from

The Zsh files live at `.config/zsh/` but Zsh looks for `$ZDOTDIR/.zshrc`. `ZDOTDIR` is set by
the shell-init stage of the installer (or your existing environment). On a fresh macOS install,
the simplest route is to add `export ZDOTDIR="$HOME/.config/zsh"` to `/etc/zshenv`.
