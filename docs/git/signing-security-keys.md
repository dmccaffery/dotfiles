---
icon: lucide/shield-check
---

# Signing & security keys

Commits are SSH-signed using a YubiKey resident key. The same key acts as the SSH auth key.

## How signing works

```ini title=".config/git/config"
[gpg]
    format = ssh

[gpg "ssh"]
    allowedSignersFile = ~/.config/git/allowed_signers
```

```ini title=".config/git/github.gitconfig"
[gpg "ssh"]
    defaultKeyCommand = "ssh-sk get --github"
```

When git needs a signing key, it runs `ssh-sk get --github` (only for github.com remotes per
`includeIf`). That command asks the `gh` CLI which of your SSH keys are marked for signing,
then checks `ssh-add -L` to find one currently loaded in the agent that matches. The match
is returned in the `key::<...>` format git expects.

## One-time setup

### 1. Generate a resident key on the YubiKey

```sh
ssh-sk gen <user>
```

The `<user>` argument is **required** and becomes the key comment (`-C`); the script exits with a
usage error if it's omitted. It generates a `ed25519-sk` **resident** key (stored on the YubiKey
itself):

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

| Option                      | Effect                                                                                                 |
| --------------------------- | ------------------------------------------------------------------------------------------------------ |
| `-O resident`               | Key lives on the YubiKey; can be re-extracted with `ssh-keygen -K`.                                    |
| `-O verify-required`        | PIN required to use the key.                                                                           |
| `-O no-touch-required`      | No fingerprint touch needed on the YubiKey itself.                                                     |
| `-O application=ssh:<user>` | Namespaces the resident credential by user so multiple keys coexist and extract to distinct filenames. |
| `-O user=<user>`            | Sets the FIDO2 user handle — the documented way to hold multiple resident keys for one application.    |

### 2. Load resident keys into the agent

`ssh-sk gen` finishes by calling `ssh-sk get`, which runs `ssh-add -K` to extract the resident
keys from the YubiKey into the running ssh-agent.

`ssh-sk get` also appends the public key to `~/.ssh/.git_allowed_signers` (the per-user
allowed-signers file referenced by `~/.config/git/allowed_signers`), so `git log
--show-signature` can verify your own commits. The append is skipped if an identical
signer line is already present, so re-running the script is safe.

### 3. Mark the key as a signing key on GitHub

```sh
gh ssh-key add ~/.ssh/<your-key>.pub --type signing --title "YubiKey signing"
```

`ssh-sk get --github` queries `gh ssh-key list` and looks for keys with the `signing` type —
that's how it discriminates between auth-only and signing-capable keys.

### 4. Verify

```sh
git-github-auth        # ensure gh has the right scopes
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

## See also

- [scripts/security-keys](../scripts/security-keys.md) — full reference for each script.
- [macOS → LaunchAgents](../macos/launchagents.md) — how the ssh-agent runs at login.
