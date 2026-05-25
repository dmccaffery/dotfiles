---
icon: lucide/bot
---

# Claude Code settings

`.claude/settings.json` is the user-level config for [Claude Code](https://claude.com/claude-code).
The shipped file is small and opinionated:

```json title=".claude/settings.json"
{
    "theme": "custom:cyberdream",
    "autoMemoryEnabled": true,
    "cleanupPeriodDays": 7,
    "editorMode": "vim",
    "attribution": { "commit": "", "pr": "" },
    "autoUpdatesChannel": "stable",
    "includeGitInstructions": false,
    "plansDirectory": ".claude/plans",
    "respectGitignore": true,
    "feedbackSurveyRate": 0,
    "sandbox": {
        "enabled": true,
        "filesystem": {
            "allowRead": ["~/Repos"],
            "allowWrite": ["~/Repos"]
        },
        "autoAllowBashIfSandboxed": false,
        "allowUnsandboxedCommands": false,
        "enableWeakerNetworkIsolation": false,
        "enableWeakerNestedSandbox": false
    },
    "statusLine": {
        "type": "command",
        "command": "oh-my-posh claude --config ~/.config/oh-my-posh/claude.yaml"
    }
}
```

## What each block does

### Theme

```json
"theme": "custom:cyberdream"
```

Points at `.claude/themes/cyberdream.json` (relative to `~/.claude/themes/`). See [Theme](theme.md).

### Memory & cleanup

```json
"autoMemoryEnabled": true,
"cleanupPeriodDays": 7
```

The model maintains persistent memory across sessions. Anything not touched for 7 days is
garbage-collected.

### Editor mode

```json
"editorMode": "vim"
```

Vim-style modal editing in the message composer.

### Attribution

```json
"attribution": { "commit": "", "pr": "" }
```

Empty strings disable Claude Code's default attribution footers on commits and PRs. The git
config's `Signed-off-by` trailer (from the prepare-commit-msg hook) is the only attribution
that lands.

### Plans

```json
"plansDirectory": ".claude/plans"
```

When Claude Code is in plan mode, plan files write to `<repo>/.claude/plans/`. The repo's
`.gitignore` excludes `.claude/plans/` by default.

### Sandbox

```json
"sandbox": {
  "enabled": true,
  "filesystem": {
    "allowRead":  ["~/Repos"],
    "allowWrite": ["~/Repos"]
  }
}
```

Hard-locks Claude Code's filesystem access to `~/Repos`. Anything outside requires explicit
permission. The three explicit `false` flags below disable common escape hatches:

- `autoAllowBashIfSandboxed: false` — Bash commands still require approval.
- `allowUnsandboxedCommands: false` — no commands can run outside the sandbox.
- `enableWeakerNetworkIsolation` / `enableWeakerNestedSandbox` — both off.

### Status line

```json
"statusLine": {
  "type": "command",
  "command": "oh-my-posh claude --config ~/.config/oh-my-posh/claude.yaml"
}
```

Claude Code runs `oh-my-posh claude …` and renders the output as a status line. The
`oh-my-posh claude` subcommand consumes Claude Code's session JSON on stdin.

See [Terminal → oh-my-posh](../terminal/oh-my-posh.md#claude-code-status-line).
