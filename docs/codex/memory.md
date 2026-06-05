---
icon: lucide/brain
---

# Memory

`.config/codex/AGENTS.md` is a symlink to the Claude Code user memory at `~/.claude/CLAUDE.md`, so Codex, OpenCode,
and Claude Code all read the **same** global agent instructions from a single source of truth. Codex loads
`$CODEX_HOME/AGENTS.md` as its top-level agent guidance, and `CODEX_HOME` resolves to `~/.config/codex`.

```text
.config/codex/AGENTS.md -> ../../.claude/CLAUDE.md
```

That shared file carries the cross-repo conventions every agent must follow:

- **Temp-file policy** — always create temporary files and dirs under `$TMPDIR`.
- **Commit messages** — Conventional Commits (`type(scope): summary`).
- **Commit signing** — hand the real, signed commit off to the user via a `commit.sh` script run outside the
  sandbox.

Because it's one file behind three names, there's nothing to keep in sync: edit
[`.claude/CLAUDE.md`](../claude/memory.md) and Codex, OpenCode, and Claude Code all pick the change up. The companion
[`docs/claude/memory.md`](../claude/memory.md) documents the contents of that file in full.
