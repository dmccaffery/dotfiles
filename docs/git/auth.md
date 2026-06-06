---
icon: lucide/key
---

# Authentication

Git authentication is handled differently per host. The strategy:

- **GitHub** → OAuth via `gh` CLI (or `git-credential-manager` for HTTPS pushes).
- **Codeberg / Forgejo** → OAuth via `git-credential-manager` with a per-host client ID.
- **GitLab** → Browser OAuth via GCM.

## GitHub

GitHub credentials are served by `gh` itself rather than GCM. `.config/git/github.gitconfig`
(conditionally included for `github.com` remotes) configures this:

```ini title=".config/git/github.gitconfig"
[credential "https://github.com"]
    helper =
    helper = !gh auth git-credential

[credential "https://gist.github.com"]
    helper =
    helper = !gh auth git-credential
```

The empty `helper =` line clears any previously inherited credential helper before setting the
`gh`-backed one, so GCM never intercepts GitHub requests.

For programmatic CLI access, use [`git-github-auth`](../scripts/security-keys.md#git-github-auth) —
a helper script that validates `gh auth status` and ensures the required scopes are present:

> `gist`, `notifications`, `project`, `repo`, `user`, `workflow`, `read:org`,
> `read:public_key`, `read:ssh_signing_key`, `write:ssh_signing_key`

If any scope is missing the script re-runs `gh auth refresh` with the full set.

### Per-repo account switching

The [`gh-switch-user`](../scripts/security-keys.md#gh-switch-user) wrapper (aliased as `gh` in
`.zshrc`) reads `git config github.account` and switches the active `gh` session automatically
before every `gh` command. Set it once per clone:

```sh
git config github.account <login>
```

The same key also tells [`ssh-sk get --git`](signing-security-keys.md) which local YubiKey stub to
use for SSH commit signing.

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

Set the Forgejo signing account once per clone so `ssh-sk get --git` can choose the matching local
YubiKey stub without querying the Forgejo API:

```sh
git config forgejo.account <login>
```

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
git-github-auth                # interactive picker — choose an account or "new account"
git-github-auth <login>        # switch to a specific account
```

The picker lists every account already authenticated on this machine, plus a **`new account`**
entry. Select `new account` to run a fresh `gh auth login` for an account that has never been
authenticated here.

## See also

- [Signing & security keys](signing-security-keys.md) — separate from auth: SSH signing for
  commits, using a YubiKey resident key.
- [Private overlay](private-extension.md) — where to put your work email, work GitHub host,
  etc.
