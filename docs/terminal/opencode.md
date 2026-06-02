# OpenCode

[OpenCode](https://opencode.ai) is configured as a stow-managed coding-agent TUI under
`.config/opencode/`.

## Files

| File                                      | Purpose                                                                 |
| ----------------------------------------- | ----------------------------------------------------------------------- |
| `.config/opencode/opencode.jsonc`         | Global runtime config and permission rules (JSONC, comments allowed).   |
| `.config/opencode/tui.json`               | TUI-only settings; selects the `cyberdream` theme.                      |
| `.config/opencode/AGENTS.md`              | Symlink to `~/.claude/CLAUDE.md`; shares one set of global agent rules. |
| `.config/opencode/themes/cyberdream.json` | Cyberdream palette for the OpenCode TUI.                                |
| `.config/opencode/.gitignore`             | Keeps local OpenCode package/plugin install artifacts out of the repo.  |

`AGENTS.md` is a symlink to the Claude Code user memory at `~/.claude/CLAUDE.md`, so OpenCode and
Claude Code read the same global agent instructions (temp-file, commit-message, and signing
conventions) from a single source — edit `.claude/CLAUDE.md` and both pick the change up.

## Permission Sandbox

OpenCode does **not** currently accept Claude Code's native `sandbox` block in
`opencode.jsonc`; unknown top-level keys make OpenCode fail config validation. The global
`.config/opencode/opencode.jsonc` therefore mirrors the Claude Code sandbox as closely as
OpenCode's supported `permission` schema allows.

The mapping follows the sandbox in `.claude/settings.json`:

| Claude Code setting             | OpenCode mapping                                                             |
| ------------------------------- | ---------------------------------------------------------------------------- |
| `filesystem.allowRead`          | `permission.external_directory` allow rules for `~/Repos`, `~/.config`,      |
|                                 | `~/.cache`, `~/.local/runtime`, `~/.local/share`, `~/.npm`, `/opt/homebrew`, |
|                                 | and `/tmp`.                                                                  |
| `filesystem.allowWrite`         | `permission.edit` prompts for repo/tool cache writes and allows scratch or   |
|                                 | agent-worktree writes.                                                       |
| `filesystem.denyRead`           | `permission.read`, `permission.list`, `permission.glob`, `permission.edit`,  |
|                                 | and `permission.external_directory` deny `~/.aws`, `~/.config/gcloud`,       |
|                                 | `~/.ssh`, `~/.gnupg`, and dotenv files.                                      |
| Claude `permissions.allow` Bash | `permission.bash` pre-approves the same inspection, Homebrew, Git, and       |
| allowlist                       | commit-script chmod patterns, then adds final deny patterns for credential   |
|                                 | and dotenv paths.                                                            |
| Claude `WebSearch` allow        | `permission.websearch = "allow"`; `webfetch` still prompts.                  |

OpenCode permission objects use last-match-wins ordering, so the config keeps broad rules
first and narrower allow/deny rules later. The catch-all `"*": "ask"` also changes
OpenCode's permissive default so any unmapped tool or command requires approval.

`grep` is allowed for routine searches even though OpenCode matches `grep` permissions against
the searched regex, not the target path. The path boundaries stay on `read`, `list`, `glob`,
`edit`, and `external_directory`.

## Limits

This is a permission-layer sandbox, not the same kernel-enforced boundary Claude Code runs.
Claude's `network.allowedDomains`, `allowMachLookup`, `allowUnixSockets`, and
`allowUnsandboxedCommands` fields have no OpenCode config equivalent today. Network-capable
Bash commands are therefore controlled by the Bash allowlist and approval prompts rather than
by domain-level egress rules.

After changing `opencode.jsonc`, quit and restart OpenCode; config is loaded at startup.

## TUI Theme

`.config/opencode/tui.json` selects the bundled `cyberdream` theme:

```json title=".config/opencode/tui.json"
{
    "$schema": "https://opencode.ai/tui.json",
    "theme": "cyberdream"
}
```

The theme keeps `defs.bg` set to `"none"` so OpenCode uses the terminal background instead of
painting its own main panel color.
