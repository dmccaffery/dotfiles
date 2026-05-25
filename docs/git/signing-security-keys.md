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
    defaultKeyCommand = "git-github-sk"
```

When git needs a signing key, it runs `git-github-sk` (only for github.com remotes per
`includeIf`). That script asks the `gh` CLI which of your SSH keys are marked for signing,
then checks `ssh-add -L` to find one currently loaded in the agent that matches. The match
is returned in the `key::<...>` format git expects.

## One-time setup

### 1. Generate resident keys on the YubiKey

```sh
gen-sk-ssh "Optional Comment"
```

The script generates two keys — `ecdsa-sk` and `ed25519-sk` — both as **resident** keys
(stored on the YubiKey itself). Each key is:

```sh
ssh-keygen \
    -t ecdsa-sk \
    -O resident \
    -O verify-required \
    -O no-touch-required \
    -O application=ssh:key-touch-required \
    "$@"
```

| Option                                  | Effect                                                                             |
| --------------------------------------- | ---------------------------------------------------------------------------------- |
| `-O resident`                           | Key lives on the YubiKey; can be re-extracted with `ssh-keygen -K`.                |
| `-O verify-required`                    | PIN required to use the key.                                                       |
| `-O no-touch-required`                  | No fingerprint touch needed on the YubiKey itself.                                 |
| `-O application=ssh:key-touch-required` | But the _application_ tag includes "touch-required" — many tools still gate on it. |

### 2. Load resident keys into the agent

`gen-sk-ssh` finishes by calling `get-sk-ssh`, which runs `ssh-add -K` to extract the resident
keys from the YubiKey into the running ssh-agent.

`get-sk-ssh` also appends the public key to `~/.ssh/.git_allowed_signers` (the per-user
allowed-signers file referenced by `~/.config/git/allowed_signers`), so `git log
--show-signature` can verify your own commits. The append is skipped if an identical
signer line is already present, so re-running the script is safe.

### 3. Mark the key as a signing key on GitHub

```sh
gh ssh-key add ~/.ssh/<your-key>.pub --type signing --title "YubiKey signing"
```

`git-github-sk` queries `gh ssh-key list` and looks for keys with the `signing` type — that's
how it discriminates between auth-only and signing-capable keys.

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
