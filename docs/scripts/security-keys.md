---
icon: lucide/shield-check
---

# Security key scripts

Six scripts cover the YubiKey-resident SSH workflow end-to-end, plus one helper for
re-signing committed history.

## `gen-sk-ssh` { #gen-sk-ssh }

```sh
gen-sk-ssh "optional comment"
```

Generates **two** resident keys on the YubiKey — `ecdsa-sk` and `ed25519-sk`. Each is created
with:

```sh
ssh-keygen \
    -t ecdsa-sk \
    -O resident \
    -O verify-required \
    -O no-touch-required \
    -O application=ssh:key-touch-required \
    -C "<comment>"
```

After both keys are generated, the script calls [`get-sk-ssh`](#get-sk-ssh) to load them.

!!! note "Why two keys?"
Different services accept different key types. Shipping both means the same YubiKey works
everywhere without re-running `gen-sk-ssh`.

## `get-sk-ssh` { #get-sk-ssh }

```sh
get-sk-ssh
```

1. If `ssh-add -L` shows no keys, runs `ssh-add -K` to extract resident keys from the
   YubiKey into the running agent.
2. Resolves the user's signing key by trying [`git-github-sk`](#git-github-sk) first and
   falling back to [`git-forgejo-sk`](#git-forgejo-sk).
3. Appends the public key to `~/.ssh/.git_allowed_signers` in the
   `<email> namespaces="git" <pubkey>` format git expects, skipping the write if an
   identical line is already present.

Run this after `gen-sk-ssh`, after `ssh-agent` restarts, or after switching YubiKeys.

## `git-github-sk` { #git-github-sk }

```sh
git-github-sk        # prints `key::<pubkey>` on stdout
```

Used as `git.gpg.ssh.defaultKeyCommand` for github.com remotes. The flow:

1. `gh ssh-key list` → all SSH keys associated with the current GitHub user, filtered to
   `signing`-type keys.
2. `ssh-add -L` → all public keys currently in the local ssh-agent.
3. Find the first agent key whose public key blob also appears in the GitHub list.
4. Emit it as `key::<line>` for git to consume.

Exits non-zero with a useful error if either set is empty or there's no match.

## `git-forgejo-sk` { #git-forgejo-sk }

```sh
git-forgejo-sk        # prints `key::<pubkey>` on stdout
```

The Forgejo counterpart to [`git-github-sk`](#git-github-sk) — used as
`git.gpg.ssh.defaultKeyCommand` for Forgejo remotes. The flow:

1. `fj user key list --verbose` → all SSH keys associated with the current Forgejo user,
   filtered to entries whose key blob starts with `sk-ssh-` (i.e. security-key-backed keys).
2. `ssh-add -L` → all public keys currently in the local ssh-agent. If the agent is empty,
   tries `ssh-add -K` once to pull resident keys off the YubiKey.
3. Find the first agent key whose public key blob also appears in the Forgejo list.
4. Emit it as `key::<line>` for git to consume.

Exits non-zero with a useful error if either set is empty or there's no match.

## `git-github-auth` { #git-github-auth }

```sh
git-github-auth              # interactive picker if multiple accounts
git-github-auth <login>      # target a specific account
```

Ensures `gh` is logged in with the scopes this dotfiles setup needs:

> `gist`, `workflow`, `repo`, `user`, `read:org`, `read:public_key`, `read:ssh_signing_key`,
> `delete_repo`

Behaviour:

- Not logged in → starts `gh auth login --web --git-protocol https --scopes <set>`.
- Logged in but wrong account → `gh auth switch --user <login>`.
- Missing scopes → `gh auth refresh --scopes <set>`.

## `git-resign` { #git-resign }

```sh
git-resign <target-ref>
```

Re-signs every commit from `<target-ref>..HEAD` with the currently-configured signing key:

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
