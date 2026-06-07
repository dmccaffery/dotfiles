---
icon: lucide/shield-check
---

# Security key scripts

The `ssh-sk` command drives the YubiKey-resident SSH signing workflow end-to-end. A handful of
standalone helpers round it out: `gh-switch-user` and `git-github-auth` for GitHub accounts,
`git-resign` for re-signing committed history, and `ssh-askpass` for PIN entry.

## `ssh-sk` { #ssh-sk }

A single dispatcher with two verbs. `gen` creates a resident Ed25519 security key and stores its
local stub by YubiKey serial; `get` loads saved stubs and updates the git allowed-signers file.
The `get --git` flag narrows `get` to just resolving and printing the configured signing key —
the form git's `defaultKeyCommand` consumes.

### `ssh-sk gen` { #ssh-sk-gen }

```sh
ssh-sk gen [user]
```

The optional `[user]` becomes the key comment (`-C`) and the local stub suffix. If it is omitted,
the command builds an fzf picker from the unique usernames already authenticated with GitHub
(`gh auth status`) and Forgejo (`fj auth list`). Exactly one YubiKey must be inserted so the script
can save the generated stub under that key's serial.

Generates a resident `ed25519-sk` key on the YubiKey:

```sh
ssh-keygen \
    -t ed25519-sk \
    -f ~/.config/private/ssh/<serial>/id_ed25519_sk_<user> \
    -O resident \
    -O verify-required \
    -O no-touch-required \
    -O application=ssh:<user> \
    -O user=<user> \
    -C "<user>"
```

The `<user>` is woven into `application`, `user`, and the comment so that **multiple** resident
keys can coexist on one YubiKey. FIDO2 resident credentials are keyed by `(application,
user-handle)`, so a fixed application string would make each new key overwrite the previous one.
Namespacing the application by user also gives each saved stub a distinct filename per user.

Stubs are saved under `~/.config/private/ssh/<serial>/`, where `<serial>` comes from
`ykman list --serials`. Only Ed25519 security-key stubs named `id_ed25519_sk_*` are generated,
loaded, and considered during git-account matching.

If the selected user matches any authenticated GitHub or Forgejo account, the generated public key
is published to every match. GitHub receives a signing key via `gh ssh-key add --type signing`;
Forgejo receives a regular SSH key via `fj -H <host> user key upload`. If the user does not match
any authenticated account, `ssh-sk gen` warns that the key will not be published and asks whether to
continue. Without a tty for that confirmation, it refuses to generate the unpublished key.

!!! warning "Do not replace stubs with resident reloads"

    OpenSSH's resident-key reload path (`ssh-add -K` / `ssh-keygen -K`) synthesizes loaded keys as
    touch-required and does not preserve `-O no-touch-required`. Keep the generated stub if you want
    PIN / user-verification caching without a physical touch for every Git signature.

After the key is generated, the command calls [`ssh-sk get`](#ssh-sk-get) with the selected user to
load it and update `allowed_signers` even when the current repository has no provider account set.

### `ssh-sk get` { #ssh-sk-get }

```sh
ssh-sk get
```

1. Uses `ykman list --serials` to identify the currently inserted YubiKey serials.
2. Loads saved `id_ed25519_sk_*` stubs from `~/.config/private/ssh/<serial>/` into the running
   agent.
3. Resolves the user's signing key with [`ssh-sk get --git`](#ssh-sk-get-git).
4. Appends the public key to `~/.config/private/git/allowed_signers` in the
   `<email> namespaces="git" <pubkey>` format git expects, skipping the write if an identical line
   is already present.

Run this after `ssh-sk gen`, after `ssh-agent` restarts, or after switching YubiKeys.

### `ssh-sk get --git` { #ssh-sk-get-git }

```sh
ssh-sk get --git        # prints `key::<pubkey>` on stdout
```

Used as `gpg.ssh.defaultKeyCommand` for GitHub and Forgejo remotes. The flow is fully offline:

1. Reads `git config --get github.account` and `git config --get forgejo.account`.
2. Refuses to sign if neither account is set, or if both are set to different usernames.
3. `ykman list --serials` → currently inserted YubiKey serials.
4. `~/.config/private/ssh/<serial>/id_ed25519_sk_<account>.pub` → the saved stub for the configured
   account on the inserted YubiKey.
5. `ssh-add -L` → all public keys currently in the local ssh-agent after loading the matching
   saved stubs.
6. Find the first agent key whose blob matches the configured account's saved stub.
7. Emit it as `key::<line>` for git to consume.

Because this resolver does not query GitHub or Forgejo, commit signing works without network access
and does not depend on the active `gh` / `fj` session. If no inserted YubiKey has a saved signing key
for the configured account, it prompts once for another key and retries. Only the `key::` line goes
to stdout, so git reads it cleanly.

## `gh-switch-user` { #gh-switch-user }

```sh
gh-switch-user <gh-args…>   # called automatically via alias gh='gh-switch-user'
```

A thin `gh` wrapper (now a [`dot`](../tooling/dot.md) applet) that reads `git config github.account`
from the current repository and, if the named account is not already the active `gh` session, calls
`gh auth switch --user <account>` before forwarding all arguments to `gh`.

Configured as the shell-wide `gh` alias in `.zshrc` so that every `gh` invocation in a repo with
`github.account` set automatically operates under the right identity — no manual switching needed.

Set a per-repo account with:

```sh
git config github.account <login>
```

If `github.account` is not set (e.g., outside a repo or in a repo without the key), `gh` runs
unmodified against whichever account is currently active.

## `git-github-auth` { #git-github-auth }

```sh
git-github-auth              # interactive picker if multiple accounts
git-github-auth <login>      # target a specific account
```

Ensures `gh` is logged in with the scopes this dotfiles setup needs:

> `gist`, `notifications`, `project`, `repo`, `user`, `workflow`, `read:org`,
> `read:public_key`, `read:ssh_signing_key`, `write:ssh_signing_key`

Behaviour:

- Not logged in → starts `gh auth login --web --git-protocol https --scopes <set>`.
- Logged in but wrong account → `gh auth switch --user <login>`.
- Missing scopes → `gh auth refresh --scopes <set>`.

When no `<login>` argument is provided, the script presents an fzf picker listing every account
already authenticated on this machine, plus a **`new account`** entry at the bottom. Selecting
`new account` runs `gh auth login` so you can authenticate a GitHub account that has never been
used on this machine.

## `git-resign` { #git-resign }

```sh
git-resign <target-ref>
```

Re-signs every commit from `<target-ref>..HEAD` with the currently-configured signing key (now a
[`dot`](../tooling/dot.md) applet, invoked as `git resign` via git's `git-*` subcommand dispatch):

```sh
git rebase --exec 'git commit --amend --no-edit -n -S' -i "${target}"
```

Common use: you authored commits on a machine without the YubiKey loaded, and now you want
the chain signed before opening a PR.

!!! warning "History rewrite"

    This rewrites commit SHAs. Don't run it on a branch that's already published unless you
    intend to force-push and notify collaborators.

## `ssh-askpass` { #ssh-askpass }

A bash wrapper that bridges `SSH_ASKPASS` callbacks from `ssh-agent` to
[`pinentry-mac`](https://gpgtools.org). Wired up via the
[launch agent plist](../macos/launchagents.md).

Key behaviour:

- If the prompt starts with `"Confirm user presence"` (FIDO2 user-verification flow), echo a
  newline immediately — no PIN needed.
- Otherwise, parse the SHA256 fingerprint out of the prompt and pass it to pinentry as
  `SETKEYINFO`, so the pinentry GUI can show _which_ key is being unlocked.
- Strip pinentry's `D ` prefix from the returned PIN before echoing it back.

This wrapper can satisfy OpenSSH's askpass callbacks, but it cannot bypass authenticator-enforced
touch. No-touch signing depends on loading the saved `id_ed25519_sk_*` stub that preserves
OpenSSH's `no-touch-required` flag.
