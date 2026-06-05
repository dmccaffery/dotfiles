---
icon: lucide/shield-check
---

# Signing & security keys

Commits are SSH-signed using YubiKey-resident Ed25519 security keys, with local stubs saved by
YubiKey serial.

## How signing works

```ini title=".config/git/config"
[gpg]
    format = ssh

[gpg "ssh"]
    allowedSignersFile = ~/.config/private/git/allowed_signers
```

```ini title=".config/git/github.gitconfig and .config/git/forgejo.gitconfig"
[gpg "ssh"]
    defaultKeyCommand = "ssh-sk get --git"
```

When git needs a signing key, it runs `ssh-sk get --git` from the provider-specific include. That
command reads `github.account` or `forgejo.account` from git config, loads the saved stub for that
account from the currently inserted YubiKey serials, and returns it in the `key::<...>` format git
expects.

If neither account is configured, or both are configured to different usernames, `ssh-sk get --git`
refuses to sign. It never queries GitHub or Forgejo, so commit signing works offline and does not
depend on the active provider CLI session.

## One-time setup

### 1. Generate a resident key on the YubiKey

```sh
ssh-sk gen [user]
```

The optional `[user]` becomes the key comment (`-C`) and the local stub suffix. If it is omitted,
`ssh-sk gen` opens an fzf picker containing the unique usernames already authenticated with GitHub
or Forgejo. Exactly one YubiKey must be inserted so the script can save the stub under that key's
serial:

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

| Option                      | Effect                                                                                                 |
| --------------------------- | ------------------------------------------------------------------------------------------------------ |
| `-f …/<serial>/…`           | Saves the local Ed25519-SK stub under the inserted YubiKey's serial.                                   |
| `-O resident`               | Key lives on the YubiKey; can be re-extracted with `ssh-keygen -K`.                                    |
| `-O verify-required`        | PIN required to use the key.                                                                           |
| `-O no-touch-required`      | No physical touch is needed when OpenSSH loads the saved stub.                                         |
| `-O application=ssh:<user>` | Namespaces the resident credential by user so multiple keys coexist and extract to distinct filenames. |
| `-O user=<user>`            | Sets the FIDO2 user handle — the documented way to hold multiple resident keys for one application.    |

Repeat `ssh-sk gen [user]` once per YubiKey that should be able to sign. Each key gets its own
`~/.config/private/ssh/<serial>/` directory.

!!! warning "Keep the generated stub"

    OpenSSH's resident-key reload path (`ssh-add -K` / `ssh-keygen -K`) synthesizes resident keys
    as touch-required and does not preserve `-O no-touch-required`. The saved stub contains the key
    handle and OpenSSH flags needed to keep no-touch signing working.

### 2. Load saved stubs into the agent

`ssh-sk gen` finishes by calling `ssh-sk get` with the generated username, which uses
`ykman list --serials` to find the inserted YubiKey and loads matching stubs from
`~/.config/private/ssh/<serial>/` into the running ssh-agent.

`ssh-sk get` also appends the public key to `~/.config/private/git/allowed_signers`, so
`git log --show-signature` can verify your own commits. The append is skipped if an identical
signer line is already present, so re-running the script is safe.

### 3. Publish the key to authenticated providers

If `[user]` matches any authenticated account in `gh auth status` or `fj auth list`, `ssh-sk gen`
publishes the generated public key to every matching account. GitHub receives a signing key via
`gh ssh-key add --type signing`; Forgejo receives a regular SSH key via
`fj -H <host> user key upload`.

If no authenticated provider account matches `[user]`, `ssh-sk gen` warns that the key will not be
published and asks whether to continue. Without a tty for that confirmation, it refuses to generate
the unpublished key.

### 4. Verify

```sh
git-github-auth        # ensure gh can upload SSH signing keys
git commit -S          # the -S is implicit when gpg.format is ssh + a key is set
```

## Re-signing a chain of commits

If commits were authored without signing (e.g., from an older machine):

```sh
git-resign <target-ref>
```

Rebases from `<target-ref>..HEAD`, amending each commit with `--no-edit -n -S` to apply the
current signing key.

## ssh-agent + ssh-askpass

The Homebrew-built `ssh-agent` runs as a launch agent (see
[macOS → LaunchAgents](../macos/launchagents.md)) with `SSH_ASKPASS` pointing at the included
[`ssh-askpass`](../scripts/security-keys.md#ssh-askpass) script. The script delegates to
`pinentry-mac` for the PIN prompt, with SHA256 fingerprint detection so the prompt knows
_which_ key is being unlocked.

PIN caching and physical touch are separate. `pinentry-mac` can cache the PIN / user-verification
step, but the no-touch behavior depends on OpenSSH loading the saved stub that still carries the
`no-touch-required` flag.

## See also

- [scripts/security-keys](../scripts/security-keys.md) — full reference for each script.
- [macOS → LaunchAgents](../macos/launchagents.md) — how the ssh-agent runs at login.
