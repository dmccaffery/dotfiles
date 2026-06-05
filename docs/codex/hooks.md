---
icon: lucide/webhook
---

# Hooks

Codex's lifecycle hooks (`[hooks]` in `config.toml`) do two jobs here, both mirrored from Claude Code: they drive the
tmux "needs you" indicator, and they enforce the [credential-read guard](sandbox-permissions.md#the-credential-read-guard-deny-hook).
Each hook receives the event as JSON on stdin and runs a shell command; a `PreToolUse` hook can return a decision
that blocks the tool.

## The tmux indicator

Codex reuses the **same** tool-agnostic [`agent-tmux-status`](../scripts/tmux.md#agent-tmux-status) script that
Claude Code's hooks and OpenCode's plugin drive — it sets a per-window `@agent_status` token that
[`theme.conf`](../terminal/tmux.md#agent-status) renders as a colour + glyph. The hook mapping mirrors the
[Claude Code hooks](../claude/hooks-skills.md#claude-is-waiting-indicator):

| Codex hook          | Runs                          | Claude Code equivalent | Effect                               |
| ------------------- | ----------------------------- | ---------------------- | ------------------------------------ |
| `Stop`              | `agent-tmux-status waiting`   | `Stop`                 | Turn finished — calm **peach** `●`.  |
| `PermissionRequest` | `agent-tmux-status attention` | `Notification`         | Approval needed — bold **red** `󰂚`.  |
| `PostToolUse`       | `agent-tmux-status clear`     | `PostToolUse`          | A tool ran after approval — cleared. |
| `UserPromptSubmit`  | `agent-tmux-status clear`     | `UserPromptSubmit`     | You replied — cleared.               |

Codex has no `SessionEnd` hook (Claude Code clears on that too); the four above cover the indicator lifecycle. Like
Claude Code's, the script is no-op-safe — every `tmux`/`printf` call is guarded, so a hook can never fail a turn.

## The credential-read guard

```toml title=".config/codex/config.toml"
[[hooks.PreToolUse]]
matcher = "^Bash$"
[[hooks.PreToolUse.hooks]]
type = "command"
command = "~/.config/codex/hooks/pre-tool-use-policy"
```

Scoped to the Bash tool, the hook reads the command from `.tool_input.command` on stdin and emits a
`permissionDecision = "deny"` when the command references a credential path or the macOS keychain — recovering the
protection Codex's Seatbelt sandbox can't provide. The full rationale and the blocked set are in
[Sandbox & permissions](sandbox-permissions.md#the-credential-read-guard-deny-hook).

## How Codex hooks run

`hooks` is a stable, on-by-default Codex feature (`codex features list`). Each handler is `type = "command"`; Codex
pipes the event JSON to the command's stdin. The status hooks ignore stdin and take their state from the argument
(`agent-tmux-status waiting`); the deny-hook parses stdin and returns its decision as JSON on stdout. Validate that
`config.toml` parses and the hooks are accepted with `codex doctor`.
