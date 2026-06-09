---
icon: lucide/shield
---

# Permissions

`.config/opencode/opencode.jsonc` re-creates the Claude Code sandbox as closely as OpenCode's permission schema
allows. It is the OpenCode counterpart to the [Claude Code sandbox](../claude/settings.md#sandbox) — same intent
(let the agent roam the dev tree, keep it out of credentials), expressed in a different config language.

## Why a permission mirror, not a sandbox

OpenCode does **not** currently accept Claude Code's native `sandbox` block in `opencode.jsonc`; unknown top-level
keys make OpenCode fail config validation. So instead of a kernel-enforced boundary, the config leans on OpenCode's
supported `permission` schema to approximate the same allow/deny shape.

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

## Ordering and the catch-all

OpenCode permission objects use **last-match-wins** ordering, so the config keeps broad rules first and narrower
allow/deny rules later — the credential and dotenv denies sit at the end of each block so nothing earlier can
re-open them. The top-level `"*": "ask"` flips OpenCode's permissive default: any tool or command that isn't
explicitly mapped requires approval.

`"security *": "deny"` blocks the macOS keychain-dump CLI outright, mirroring Claude Code's `Bash(security *)`
deny and Codex's `forbidden` rule for `security` — the credential-exfiltration path is closed in all three agents.

`grep` is allowed for routine searches even though OpenCode matches `grep` permissions against the searched regex,
not the target path. The path boundaries therefore stay on `read`, `list`, `glob`, `edit`, and `external_directory`,
which do match against paths.

## Limits

This is a permission-layer approximation, not the kernel-enforced boundary Claude Code runs. Claude's
`network.allowedDomains`, `allowMachLookup`, `allowUnixSockets`, and `allowUnsandboxedCommands` fields have no
OpenCode config equivalent today. Network-capable Bash commands are therefore controlled by the `permission.bash`
allowlist and approval prompts rather than by domain-level egress rules.

See [`opencode.jsonc`](https://github.com/dmccaffery/dotfiles/blob/main/stow/.config/opencode/opencode.jsonc) for the
full rule set. After changing it, quit and restart OpenCode — config is loaded at startup.
