---
icon: lucide/square-terminal
---

# Codex

[Codex](https://developers.openai.com/codex) is OpenAI's coding-agent CLI, configured as a stow-managed package
under `.config/codex/`. Like [OpenCode](../opencode/index.md), it deliberately shares this machine's guardrails with
[Claude Code](../claude/index.md): a sandbox that mirrors the Claude Code sandbox, an execpolicy + hook layer that
mirrors the Claude Code permissions, the cyberdream palette, and the tmux "needs you" status indicator. Switching
between the three agents keeps the same boundaries and the same look.

Codex defaults to `~/.codex`, but `CODEX_HOME` is set to `~/.config/codex` in
[`.zshenv`](../terminal/shell.md) so the whole package lives under XDG config like every other tool — only Claude
Code hardcodes its `~/.claude` home.

| Page                                            | Purpose                                                                          |
| ----------------------------------------------- | -------------------------------------------------------------------------------- |
| [Sandbox & permissions](sandbox-permissions.md) | `config.toml` + `rules/` + `hooks/` — the sandbox and permission mirror.         |
| [Hooks](hooks.md)                               | The lifecycle hooks that drive the tmux indicator and the credential-read guard. |
| [Theme](theme.md)                               | Why `theme = "ansi"` is the cyberdream mirror, and what Codex can't theme.       |
| [Memory](memory.md)                             | `AGENTS.md` — symlink to the shared `~/.claude/CLAUDE.md` agent rules.           |

## Files

| File                                      | Purpose                                                                            |
| ----------------------------------------- | ---------------------------------------------------------------------------------- |
| `.config/codex/config.toml`               | Sandbox, approvals, lifecycle hooks, and TUI theme (`$CODEX_HOME/config.toml`).    |
| `.config/codex/rules/default.rules`       | Execpolicy allow/forbidden rules — the command allowlist (Starlark).               |
| `.config/codex/hooks/pre-tool-use-policy` | PreToolUse deny-hook blocking credential-path reads + the macOS keychain.          |
| `.config/codex/AGENTS.md`                 | Symlink to `~/.claude/CLAUDE.md`; shares one set of global agent rules.            |
| `.config/codex/.gitignore`                | Fail-closed ignore — keeps Codex's runtime state (credentials, history) untracked. |

## Why a fail-closed `.gitignore`

`CODEX_HOME` points at this stow-managed directory, so Codex writes **all** of its runtime state here:
`auth.json` (an API token), session history, logs, caches, and lockfiles. The
[`.gitignore`](https://github.com/dmccaffery/dotfiles/blob/main/stow/.config/codex/.gitignore) ignores everything and
re-includes only the five hand-authored files, so a new state file Codex invents tomorrow is ignored by default,
not leaked into the repo.

## Reloading config

Codex loads its config at startup. After editing `config.toml`, the rules, or the hooks, quit and restart Codex for
the change to take effect. Validate any edit with `codex doctor` (it reports the parse result and the resolved
sandbox/approval summary).
