---
icon: lucide/brain
---

# Global agent memory

`.claude/CLAUDE.md` is the **user-level** memory for [Claude Code](https://claude.com/claude-code) —
stowed to `~/.claude/CLAUDE.md`, where it is loaded in _every_ repository on the machine. It holds
cross-repository conventions; a project's own `CLAUDE.md` / `AGENTS.md` layers repo-specific rules on top.

`.claude/CLAUDE.md` and [`.claude/settings.json`](settings.md) are both stowed user-level files — they apply to
Claude Code in _every_ repo. The repo's own `CLAUDE.md` (→ `AGENTS.md`) is the project-scoped layer that sits on
top of them and is deliberately **not** stowed.

| File                               | Scope            | Stowed?                                  |
| ---------------------------------- | ---------------- | ---------------------------------------- |
| `.claude/CLAUDE.md`                | All repos (user) | Yes → `~/.claude/CLAUDE.md`              |
| `.claude/settings.json`            | All repos (user) | Yes → `~/.claude/settings.json`          |
| `<repo>/CLAUDE.md` (→ `AGENTS.md`) | This repo only   | No (`.stowrc` ignores top-level entries) |

## What it covers

The shipped file documents conventions that hold regardless of which repo Claude is working in:

- **Commit messages** — [Conventional Commits](https://www.conventionalcommits.org/) (`type(scope): summary`),
  the format release-please and friends parse to derive versions and changelogs.
- **Creating commits** — the `.commit.sh` handoff. SSH commit signing runs through `ssh-agent`, which the Claude
  Code sandbox refuses (`EPERM` on the agent socket), so commits made inside the sandbox are unsigned. Claude writes
  a `.commit.sh` at the working-directory root and the user runs it **outside** the sandbox, where the signing key is
  reachable. Inside a worktree, it instead re-signs the session's commits with
  [`git resign`](../scripts/security-keys.md#git-resign).

See the source file for the authoritative wording.

## Stowing it

Getting these files to `$HOME` correctly meant tightening [`.stowrc`](../tooling/stow-and-make.md). `stow`'s
`--ignore` patterns match the **full relative path, anchored at the end**, so two classes of pattern were too broad:

- A bare `--ignore=CLAUDE.md` matches _any_ path ending in `CLAUDE.md` — catching both the top-level `CLAUDE.md`
  symlink and `.claude/CLAUDE.md`. The top-level entry is pinned to the repo root with `^CLAUDE.md`; the transient
  `.commit.sh` is kept out of `$HOME` the same way (`^.commit.sh`).
- The broad `--ignore=.*.json` / `--ignore=.*.yaml` patterns were **removed**. By the same suffix rule they were
  silently skipping `.claude/settings.json` and `.claude/themes/cyberdream.json` whenever `stow` had to descend into
  an already-existing `~/.claude`. Root-level lockfiles and manifests are now ignored by explicit `^`-anchored
  entries instead.

See [Stow & Make](../tooling/stow-and-make.md#stowrc) for the full mechanism (and why the anchors use `^`, not
`\A`/`\z`).
