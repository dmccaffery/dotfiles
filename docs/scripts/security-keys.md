---
icon: lucide/shield-check
---

# Security key scripts

The `ssh-sk` command drives the YubiKey-resident SSH signing workflow end-to-end. A handful of
standalone helpers round it out: `gh-switch-user` and `git-github-auth` for GitHub accounts,
`git-resign` for re-signing committed history, and `ssh-askpass` for PIN entry.

## `ssh-sk` { #ssh-sk }

A single dispatcher with two verbs. `gen` creates a resident key; `get` loads resident keys and
updates the git allowed-signers file. The `get --github` / `get --forgejo` flags narrow `get` to
just resolving and printing one provider's signing key — the form git's `defaultKeyCommand`
consumes.

### `ssh-sk gen` { #ssh-sk-gen }

```sh
ssh-sk gen <user>
```

The `<user>` argument is **required** — it becomes the key comment (`-C`). The command exits with
a `usage: ssh-sk gen <user>` error if it's omitted.

Generates a resident `ed25519-sk` key on the YubiKey:

```sh
ssh-keygen \
    -t ed25519-sk \
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
Namespacing the application by user also gives `ssh-keygen -K` (run via [`ssh-sk get`](#ssh-sk-get))
distinct extraction filenames per user.

After the key is generated, the command calls [`ssh-sk get`](#ssh-sk-get) to load it.

### `ssh-sk get` { #ssh-sk-get }

```sh
ssh-sk get
```

1. If `ssh-add -L` shows no keys, runs `ssh-add -K` to extract resident keys from the
   YubiKey into the running agent.
2. Resolves the user's signing key by trying [`ssh-sk get --github`](#ssh-sk-get-github) first and
   falling back to [`ssh-sk get --forgejo`](#ssh-sk-get-forgejo).
3. Appends the public key to `~/.ssh/.git_allowed_signers` in the
   `<email> namespaces="git" <pubkey>` format git expects, skipping the write if an
   identical line is already present.

Run this after `ssh-sk gen`, after `ssh-agent` restarts, or after switching YubiKeys.

### `ssh-sk get --github` { #ssh-sk-get-github }

```sh
ssh-sk get --github        # prints `key::<pubkey>` on stdout
```

Used as `gpg.ssh.defaultKeyCommand` for github.com remotes. The flow:

1. `gh ssh-key list` → all SSH keys associated with the current GitHub user, filtered to
   `signing`-type keys.
2. `ssh-add -L` → all public keys currently in the local ssh-agent (loading resident keys with
   `ssh-add -K` first if the agent is empty).
3. Find the first agent key whose public key blob also appears in the GitHub list.
4. Emit it as `key::<line>` for git to consume.

Exits non-zero with a useful error if either set is empty or there's no match. Only the `key::`
line goes to stdout, so git reads it cleanly.

### `ssh-sk get --forgejo` { #ssh-sk-get-forgejo }

```sh
ssh-sk get --forgejo        # prints `key::<pubkey>` on stdout
```

The Forgejo counterpart to [`ssh-sk get --github`](#ssh-sk-get-github) — used as
`gpg.ssh.defaultKeyCommand` for Forgejo remotes. The flow:

1. `fj user key list --verbose` → all SSH keys associated with the current Forgejo user,
   filtered to entries whose key blob starts with `sk-ssh-` (i.e. security-key-backed keys).
2. `ssh-add -L` → all public keys currently in the local ssh-agent. If the agent is empty,
   tries `ssh-add -K` once to pull resident keys off the YubiKey.
3. Find the first agent key whose public key blob also appears in the Forgejo list.
4. Emit it as `key::<line>` for git to consume.

Exits non-zero with a useful error if either set is empty or there's no match.

## `gh-switch-user` { #gh-switch-user }

```sh
gh-switch-user <gh-args…>   # called automatically via alias gh='gh-switch-user'
```

A thin `gh` wrapper that reads `git config github.account` from the current repository and, if
the named account is not already the active `gh` session, calls `gh auth switch --user <account>`
before forwarding all arguments to `command gh`.

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
> `read:public_key`, `read:ssh_signing_key`

Behaviour:

- Not logged in → starts `gh auth login --web --git-protocol https --scopes <set>`.
- Logged in but wrong account → `gh auth switch --user <login>`.
- Missing scopes → `gh auth refresh --scopes <set>`.

When no `<login>` argument is provided, the script presents an fzf picker listing every account already authenticated
on this machine, plus a **`new account`** entry at the bottom. Selecting `new account` runs `gh auth login` so you can
authenticate a GitHub account that has never been used on this machine.

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
