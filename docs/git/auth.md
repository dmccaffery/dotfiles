---
icon: lucide/key
---

# Authentication

Git authentication is handled differently per host. The strategy:

- **GitHub** → OAuth via `gh` CLI (or `git-credential-manager` for HTTPS pushes).
- **Codeberg / Forgejo** → OAuth via `git-credential-manager` with a per-host client ID.
- **GitLab** → Browser OAuth via GCM.

## GitHub

The repo's main config sets:

```ini title=".config/git/config"
[credential]
    gitHubAuthModes = oauth
    namespace = personal
    helper =
    helper = /usr/local/share/gcm-core/git-credential-manager

[credential "https://github.com"]
    provider = github
```

For programmatic CLI access, use [`git-github-auth`](../scripts/security-keys.md#git-github-auth) —
a helper script that validates `gh auth status` and ensures the required scopes are present:

> `gist`, `workflow`, `repo`, `user`, `read:org`, `read:public_key`, `read:ssh_signing_key`,
> `delete_repo`

If any scope is missing the script re-runs `gh auth refresh` with the full set.

## Codeberg / Forgejo

`.config/git/forgejo.gitconfig` is conditionally included when the remote URL points to
`codeberg.org` or `code.forgejo.org`:

```ini title=".config/git/forgejo.gitconfig"
[credential "https://codeberg.org"]
    oauthClientId = a4792ccc-144e-407e-86c9-5e7d8d9c3269
    oauthAuthURL = /login/oauth/authorize
    oauthTokenURL = /login/oauth/access_token
    provider = generic

[credential "https://code.forgejo.org"]
    oauthClientId = a4792ccc-144e-407e-86c9-5e7d8d9c3269
    oauthAuthURL = /login/oauth/authorize
    oauthTokenURL = /login/oauth/access_token
    provider = generic

[git-town]
    github-connector = gh
    forge-type = forgejo
```

The OAuth client ID is shared between the two hosts (the same registered application). GCM
treats them as `provider = generic` (no host-specific helper exists for Forgejo) but the OAuth
flow works because Forgejo implements the standard endpoints.

For `git-town`, `forge-type = forgejo` switches the workflow integrations (PRs become
"proposals" mapped to Forgejo PRs).

## GitLab

```ini
[credential]
    gitlabAuthModes = browser
```

GCM opens a browser tab on first auth, walks the OAuth flow, and stores the token in the
keychain.

## Switching accounts

`git-github-auth` can also switch between authenticated GitHub accounts:

```sh
git-github-auth                # interactive picker, switches via `gh auth switch`
git-github-auth <login>        # switch to a specific account
```

## See also

- [Signing & security keys](signing-security-keys.md) — separate from auth: SSH signing for
  commits, using a YubiKey resident key.
- [Private overlay](private-extension.md) — where to put your work email, work GitHub host,
  etc.
