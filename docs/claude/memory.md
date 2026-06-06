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

- **Temporary files** — always route `mktemp` through `$TMPDIR` (e.g. `mktemp -d "$TMPDIR/foo.XXXXXX"`); the bare
  default lands in `/var/folders/.../T`, which the sandbox blocks.
- **Commit messages** — [Conventional Commits](https://www.conventionalcommits.org/) (`type(scope): summary`),
  the format release-please and friends parse to derive versions and changelogs.
- **Creating commits** — the [`commit.sh` handoff](#the-commitsh-workflow) that works around the sandbox's inability
  to reach the SSH signing key. Detailed in its own section below.

See the source file for the authoritative wording.

## The `commit.sh` workflow

Commit signing on this machine runs through `ssh-agent`. Claude Code's sandbox refuses that connection — it gets
`EPERM` on `connect()` to the agent's `AF_UNIX` socket — so **anything `git commit`ed from inside the sandbox is
unsigned**. The workaround is to hand the real, signable commit back to you in a `commit.sh` script you run
**outside** the sandbox, where the signing key is reachable. Which form the script takes depends on whether Claude is
working in a git worktree:

=== "Outside a worktree"

    Claude does **not** run `git commit` at all. It writes the exact `git add` / `git commit` invocations it intended
    into `commit.sh` — one commit per `git commit` call, with real Conventional-Commit messages and any trailers —
    then tells you to run it. The signing happens when _you_ execute the script, key in hand.

=== "Inside a worktree"

    Claude commits normally at sensible stopping points; those land **unsigned** on the `agent/<name>` branch. The
    `commit.sh` it writes instead _re-signs_ the range it authored this session with
    [`git resign <base>`](../scripts/security-keys.md#git-resign). `<base>` is the parent of the first commit —
    `HEAD~3` for three commits, or `$(git merge-base HEAD <parent-branch>)` when the commit count is dynamic.

Both forms obey the same rules:

| Rule            | Detail                                                                                                                                                                                                  |
| --------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Location        | Written at the working-directory root (the repo root, or the worktree root).                                                                                                                            |
| Never committed | Kept out of version control by the global ignore at `~/.config/git/ignore` (wired up via `core.excludesFile` in [`.config/git/config`](../git/config.md)), so no per-repo `.gitignore` entry is needed. |
| Header          | Starts with `#!/usr/bin/env sh`, then `set -eu`.                                                                                                                                                        |
| Single batch    | Overwrites any prior `commit.sh` — the file is the _current_ batch, not history.                                                                                                                        |
| Executable      | `chmod +x`'d on write, so you can run it as `./commit.sh`.                                                                                                                                              |
| Self-deleting   | Ends with `rm -- "$0"`, so a successful run removes the script; under `set -eu` a failed commit aborts before the `rm`, leaving it in place to rerun.                                                   |

!!! note "Why the indirection"

    The sandbox is what keeps an autonomous agent from reaching your signing key (or anything else outside its
    allow-list) on its own. `commit.sh` doesn't poke a hole in that boundary — it moves the one privileged step,
    signing, back to a shell _you_ launch, so the agent never touches the key directly.

## Stowing it

Getting these files to `$HOME` correctly meant tightening [`.stowrc`](../tooling/stow-and-make.md). `stow`'s
`--ignore` patterns match the **full relative path, anchored at the end**, so two classes of pattern were too broad:

- A bare `--ignore=CLAUDE.md` matches _any_ path ending in `CLAUDE.md` — catching both the top-level `CLAUDE.md`
  symlink and `.claude/CLAUDE.md`. The top-level entry is pinned to the repo root with `^CLAUDE.md`; the transient
  `commit.sh` is kept out of `$HOME` the same way (`^commit.sh`).
- The broad `--ignore=.*.json` / `--ignore=.*.yaml` patterns were **removed**. By the same suffix rule they were
  silently skipping `.claude/settings.json` and `.claude/themes/cyberdream.json` whenever `stow` had to descend into
  an already-existing `~/.claude`. Root-level lockfiles and manifests are now ignored by explicit `^`-anchored
  entries instead.

See [Stow & Make](../tooling/stow-and-make.md#stowrc) for the full mechanism (and why the anchors use `^`, not
`\A`/`\z`).
