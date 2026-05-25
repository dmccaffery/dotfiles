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

| Setting                     | Effect                                                                  |
| --------------------------- | ----------------------------------------------------------------------- |
| `core.autocrlf = input`     | Normalize line endings to LF in the repo; never auto-convert on out.    |
| `log.showSignature = true`  | `git log` displays signature verification status by default.            |
| `pull.ff = only`            | `git pull` only fast-forwards — refuses to create merge commits.        |
| `pull.rebase = true`        | When fast-forward isn't possible, rebase instead of merging.            |
| `rebase.autoStash`          | Auto-stash dirty working tree before rebasing.                          |
| `rebase.autoSquash`         | Honor `fixup!` / `squash!` commit prefixes when rebasing interactively. |
| `init.defaultBranch = main` | New repos use `main`.                                                   |
| `init.templatedir`          | New repos get the commit-msg hook (see [Template](#template)).          |

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
[get-sk-ssh](../scripts/security-keys.md#get-sk-ssh) script:

```ini
[gpg]
    format = ssh

[gpg "ssh"]
    allowedSignersFile = ~/.config/git/allowed_signers
```

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
```

Provider-specific overrides are scoped via `hasconfig:remote.*.url` matching — the right
`gitconfig` is loaded only when the repo's remote matches the host. See
[Private overlay](private-extension.md) for the unconditional include.

## Template { #template }

`.config/git/template/hooks/prepare-commit-msg` adds a `Signed-off-by:` trailer to commits
that don't already have one. New repositories pick it up automatically via
`init.templatedir`.

## Aliases

Most aliases delegate to [git-town](git-town.md). See that page for the full list.
