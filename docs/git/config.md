---
icon: simple/git
---

# Git config

`.config/git/config` is the canonical user-level git config (loaded via XDG —
`~/.config/git/config`).

## Defaults

```ini
[core]
    autocrlf = input
    excludesFile = ~/.config/git/ignore

[checkout]
    defaultRemote = origin

[log]
    showSignature = true

[pull]
    ff = only
    rebase = true

[rebase]
    autoStash = true
    autoSquash = true
    interactive = true

[init]
    defaultBranch = main
    templatedir = ~/.config/git/template
```

| Setting                     | Effect                                                                      |
| --------------------------- | --------------------------------------------------------------------------- |
| `core.autocrlf = input`     | Normalize line endings to LF in the repo; never auto-convert on out.        |
| `core.excludesFile`         | Global ignore at `~/.config/git/ignore` (e.g. agent-generated `commit.sh`). |
| `log.showSignature = true`  | `git log` displays signature verification status by default.                |
| `pull.ff = only`            | `git pull` only fast-forwards — refuses to create merge commits.            |
| `pull.rebase = true`        | When fast-forward isn't possible, rebase instead of merging.                |
| `rebase.autoStash`          | Auto-stash dirty working tree before rebasing.                              |
| `rebase.autoSquash`         | Honor `fixup!` / `squash!` commit prefixes when rebasing interactively.     |
| `init.defaultBranch = main` | New repos use `main`.                                                       |
| `init.templatedir`          | New repos get the commit-msg hook (see [Template](#template)).              |

## LFS

```ini
[lfs]
    locksverify = true

[filter "lfs"]
    clean = git-lfs clean -- %f
    smudge = git-lfs smudge -- %f
    process = git-lfs filter-process
    required = true
```

LFS is required (`required = true`) — clones will fail if the LFS pointer can't be resolved.

## Signing

SSH signing using a hardware security key, with the allowed-signers file managed by the
[`ssh-sk get`](../scripts/security-keys.md#ssh-sk-get) command:

```ini
[gpg]
    format = ssh

[gpg "ssh"]
    allowedSignersFile = ~/.config/private/git/allowed_signers
```

Provider-specific includes set the default signing-key resolver:

```ini
[gpg "ssh"]
    defaultKeyCommand = "ssh-sk get --git"
```

`ssh-sk get --git` reads `github.account` or `forgejo.account` from the current repository and
matches that account to the saved YubiKey stub, so signing does not require provider API access.

See [Git → Signing](signing-security-keys.md) for the full flow.

## Credentials

```ini
[credential]
    gitHubAuthModes = oauth
    gitlabAuthModes = browser
    gitHubAccountFiltering = false
    namespace = personal
    helper =
    helper = /usr/local/share/gcm-core/git-credential-manager

[credential "https://github.com"]
    provider = github
```

- **`helper = `** (empty) clears any inherited helpers, then sets exactly one:
  Git Credential Manager.
- **`gitHubAuthModes = oauth`** — for github.com, prefer OAuth flow over PAT.
- **`namespace = personal`** — keychain entries are namespaced so a work overlay can use a
  different namespace.

See [Git → Auth](auth.md) for the full GitHub / Codeberg / GitLab story.

## Includes

```ini
[includeIf "hasconfig:remote.*.url:https://github.com/**"]
    path = ~/.config/git/github.gitconfig

[includeIf "hasconfig:remote.*.url:https://codeberg.org/**"]
    path = ~/.config/git/forgejo.gitconfig

[includeIf "hasconfig:remote.*.url:https://code.forgejo.org/**"]
    path = ~/.config/git/forgejo.gitconfig

[include]
    path = ~/.config/private/git/config

[includeIf "gitdir:~/.cache/agent/worktrees/**"]
    path = ~/.config/git/agent.gitconfig
```

Provider-specific overrides are scoped via `hasconfig:remote.*.url` matching — the right
`gitconfig` is loaded only when the repo's remote matches the host. See
[Private overlay](private-extension.md) for the unconditional include.

The trailing `includeIf "gitdir:…"` block loads
[`agent.gitconfig`](https://github.com/dmccaffery/dotfiles/blob/main/stow/.config/git/agent.gitconfig)
when the repository's `.git` directory sits under `~/.cache/agent/worktrees/` — i.e. a
worktree created by Claude Code's [WorktreeCreate hook](../claude/hooks-skills.md#worktreecreate).
Worktrees live in XDG cache (not config) because they're throwaway work areas. The same
literal path is hard-coded in the hook script; both sides match by convention rather than via
an environment variable, because git's `includeIf` can't expand env vars. The agent overlay
sets `commit.gpgSign = false` and `tag.gpgSign = false` so a sandboxed agent (which can't
reach the SSH key on the security key) can commit without blocking on a signing prompt. It
is loaded last so it wins ties against the private overlay.

## Template { #template }

`.config/git/template/hooks/prepare-commit-msg` adds a `Signed-off-by:` trailer to commits
that don't already have one. New repositories pick it up automatically via
`init.templatedir`.

## Aliases

Most aliases delegate to [git-town](git-town.md). See that page for the full list.
