# Global agent instructions

Cross-repository conventions for this machine. These apply in every repo; a project's own
`CLAUDE.md` / `AGENTS.md` layers repo-specific rules on top.

## Temporary files

Always create temp files and dirs under `$TMPDIR`, e.g. `mktemp -d "$TMPDIR/foo.XXXXXX"`.
A bare `mktemp` / `mktemp -d` defaults to macOS's per-user `/var/folders/.../T`, which the
sandbox blocks (`Operation not permitted`); even bare `/tmp` resolves to `/private/tmp` outside
the allowed `…/claude` subdir and fails the same way. The Bash tool sets `$TMPDIR` to a
sandbox-writable directory, so routing every temp path through it is what actually works.

## Commit messages

Use [Conventional Commits](https://www.conventionalcommits.org/): `type(scope): summary`
(`feat:`, `fix:`, `docs:`, `chore:`, `refactor:`, `test:`, `ci:`, …). Tools like
release-please derive the next version and changelog from these prefixes, so the format is
load-bearing rather than cosmetic. Signal breaking changes with a `!` after the type
(`feat!:`) or a `BREAKING CHANGE:` footer.

## Creating commits

Commit signing on this machine runs through `ssh-agent`, which Claude Code's sandbox refuses
(it returns `EPERM` on `connect()` to the agent's `AF_UNIX` socket). Anything `git commit`ed
from inside the sandbox is therefore unsigned. Hand the real commit off to the user with a
`.commit.sh` script they run **outside** the sandbox, where the signing key is reachable:

- **Outside a worktree** — don't run `git commit` at all. Write `.commit.sh` at the repo root
  containing the exact `git add` / `git commit` invocations you intended (one commit per
  `git commit` call, real Conventional-Commit messages, trailers, etc.), then tell the user to
  run it.
- **Inside a worktree** — commit normally at sensible stopping points; those land unsigned on
  the `agent/<name>` branch. Still write `.commit.sh` at the worktree root, but its job is to
  re-sign: invoke `git resign <base>` over the range you authored this session. Pick `<base>`
  as the parent of your first commit (e.g. `HEAD~3` for three commits, or
  `$(git merge-base HEAD <parent-branch>)` when the count is dynamic).

Both variants:

- Live at the working-directory root and stay out of version control — add `.commit.sh` to
  the repo's `.gitignore` or `.git/info/exclude` if a `.*` rule doesn't already cover it.
- Start with `#!/usr/bin/env sh` + `set -eu` and overwrite any prior `.commit.sh` — the file
  is the _current_ batch, not history.
- Are `chmod +x`'d when written so the user can run them as `./.commit.sh`.
