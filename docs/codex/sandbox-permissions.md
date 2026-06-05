---
icon: lucide/shield
---

# Sandbox & permissions

`.config/codex/config.toml`, `.config/codex/rules/default.rules`, and `.config/codex/hooks/pre-tool-use-policy`
together re-create the [Claude Code sandbox](../claude/settings.md#sandbox) and its permission model as closely as
Codex's config surface allows — same intent (let the agent roam the dev tree, keep it out of credentials), expressed
in Codex's own three layers: an OS sandbox, an execution policy, and lifecycle hooks.

## What maps onto what

| Claude Code                                  | Codex equivalent                                                                      |
| -------------------------------------------- | ------------------------------------------------------------------------------------- |
| `sandbox.filesystem.allowWrite`              | `sandbox_mode = "workspace-write"` + `[sandbox_workspace_write].writable_roots`       |
| `sandbox.network.allowedDomains` (~20 hosts) | `[sandbox_workspace_write].network_access = true` (boolean — no per-domain allowlist) |
| `sandbox.filesystem.denyRead`                | `hooks/pre-tool-use-policy` (Seatbelt can't deny reads — blocked at the tool layer)   |
| `permissions.allow` Bash allowlist           | `rules/default.rules` — `prefix_rule(..., decision = "allow")`                        |
| `permissions.deny` (`Bash(security *)`)      | `rules/default.rules` `forbidden` + the deny-hook for credential paths                |
| default "ask"                                | `approval_policy = "untrusted"` — prompt unless a rule allows                         |

## Sandbox (`config.toml`)

`sandbox_mode = "workspace-write"` is the analog of Claude Code's filesystem sandbox: Codex reads broadly but may
only **write** to the working directory plus the listed roots; macOS enforces it with Seatbelt, the same primitive
Claude Code uses. `[sandbox_workspace_write].writable_roots` mirrors `allowWrite` — the agent worktrees and the
language/tool caches that live outside the current repo:

```toml title=".config/codex/config.toml"
sandbox_mode = "workspace-write"

[sandbox_workspace_write]
writable_roots = [
    "~/Repos", "~/.cache/agent/worktrees", "~/.cache/uv", "~/.cache/pip",
    "~/.cache/go", "~/.local/share/go", "~/.npm",
]
network_access = true
exclude_slash_tmp = false       # keep /tmp writable (Claude allowWrite /tmp)
exclude_tmpdir_env_var = false  # keep $TMPDIR writable (the global temp-file convention)
```

`cwd` is always writable, so working inside `~/Repos/<repo>` needs no entry — `writable_roots` only covers writes
that land outside the current repo. Verify the resolved policy with `codex doctor` (it prints
`restricted fs + enabled network · approval UnlessTrusted`).

## Permissions (`approval_policy` + execpolicy)

`approval_policy = "untrusted"` flips Codex to "prompt unless explicitly trusted" — the analog of Claude Code's
`autoAllowBashIfSandboxed = false`. The trust list lives in
[`rules/default.rules`](https://github.com/dmccaffery/dotfiles/blob/main/.config/codex/rules/default.rules), a
Starlark execpolicy that auto-loads from `$CODEX_HOME/rules/`. It mirrors the OpenCode/Claude Code Bash allowlist:

- **`decision = "allow"`** auto-runs read-only inspection (`ls`, `cat`, `grep`, `rg`, `find`, `jq`, …), read-only
  Homebrew and Git queries, the staging/branch Git verbs the agent routinely needs, and `chmod +x commit.sh`.
- **`decision = "forbidden"`** blocks `security` (the macOS keychain dump path — Claude Code's `Bash(security *)`).
- Anything unmatched (e.g. `rm -rf …`, `curl …`) falls through to an approval prompt.

`prefix_rule` matches a leading run of tokens, so `["git", "status"]` covers `git status --short`. Validate a rule
with `codex execpolicy check --rules ~/.config/codex/rules/default.rules -- <command>`.

## The credential-read guard (deny-hook)

Codex's Seatbelt sandbox only restricts **writes** in workspace-write mode — it grants broad filesystem reads, so it
has no equivalent of Claude Code's `denyRead`. The
[`pre-tool-use-policy`](https://github.com/dmccaffery/dotfiles/blob/main/.config/codex/hooks/pre-tool-use-policy)
hook recovers that protection at the tool layer: a `PreToolUse` hook scoped to the Bash tool inspects each command
and returns a `deny` decision when it touches `~/.ssh`, `~/.aws`, `~/.config/gcloud`, `~/.gnupg`, a `.env` file, or
the macOS `security` tool — the same set OpenCode denies in its `bash` rules. See [Hooks](hooks.md) for the wiring.

## Limits

This is the closest mirror Codex's config allows, not a byte-for-byte copy of the Claude Code sandbox:

- **Network is all-or-nothing.** Claude Code allow-lists ~20 hosts (npm, pypi, github, …); Codex's sandbox only
  toggles network as a whole. It's set to `true` to keep installs and fetches working — broader than Claude Code's
  egress allowlist. The deny-hook still blocks credential reads regardless of network state.
- **Read denial is tool-layer, not kernel-layer.** The deny-hook inspects Bash commands; it does not stop a non-Bash
  internal tool from reading a path the way Claude Code's kernel-enforced `denyRead` does. In practice Codex reads
  files through the shell, so the Bash-scoped hook covers the realistic paths.
