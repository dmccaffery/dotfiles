---
icon: lucide/zap
---

# Launch agents

One launch agent ships with this repo: `org.homebrew.ssh-agent`. It replaces Apple's
`com.openssh.ssh-agent` with the Homebrew build so that `SSH_ASKPASS` actually works.

## The plist

<!-- markdownlint-disable MD013 -->

```xml title="Library/LaunchAgents/org.homebrew.ssh-agent.plist"
<?xml version="1.0" encoding="UTF-8"?>
<plist version="1.0">
    <dict>
        <key>Label</key>
        <string>org.homebrew.ssh-agent</string>
        <key>ProgramArguments</key>
        <array>
            <string>/bin/sh</string>
            <string>-c</string>
            <string>rm -f $SSH_AUTH_SOCK; SSH_ASKPASS=/usr/local/bin/ssh-askpass DISPLAY=':0' /opt/homebrew/bin/ssh-agent -D -a $SSH_AUTH_SOCK</string>
        </array>
        <key>RunAtLoad</key>
        <true />
        <key>StandardOutPath</key>
        <string>/tmp/org.homebrew.ssh-agent.out.log</string>
        <key>StandardErrorPath</key>
        <string>/tmp/org.homebrew.ssh-agent.err.log</string>
    </dict>
</plist>
```

The single shell command:

1. Removes any stale `$SSH_AUTH_SOCK` file from a previous run.
2. Sets `SSH_ASKPASS=/usr/local/bin/ssh-askpass` (the [wrapper script](../scripts/security-keys.md#ssh-askpass)).
3. Sets `DISPLAY=:0` — required for ssh-agent to call `SSH_ASKPASS` on macOS even though
   there's no X server.
4. Runs Homebrew's `ssh-agent` in the foreground (`-D`) listening on the auth socket (`-a`).

## How it's installed

The `shell` stage of `install.sh` (`setup/darwin/shell.sh`) bootstraps it:

```sh title="setup/darwin/shell.sh"
id=${UID:-$(id -u)}

launchctl disable gui/"${id}"/com.openssh.ssh-agent 2> /dev/null || true
launchctl bootout  gui/"${id}"/org.homebrew.ssh-agent 2> /dev/null || true

launchctl bootstrap gui/"${id}" ~/Library/LaunchAgents/org.homebrew.ssh-agent.plist
```

- Apple's ssh-agent is _disabled_ (not just unloaded) so it doesn't race to claim the auth
  socket.
- The Homebrew ssh-agent is bootstrapped into the GUI session — that means it runs once per
  user-login, not per Terminal launch.

## Logs

Both stdout and stderr are tee'd into `/tmp/`:

- `/tmp/org.homebrew.ssh-agent.out.log`
- `/tmp/org.homebrew.ssh-agent.err.log`

If `ssh-add` fails or pinentry hangs, the error log is where to look.

## Reloading after a config change

```sh
launchctl bootout  gui/$(id -u)/org.homebrew.ssh-agent
launchctl bootstrap gui/$(id -u) ~/Library/LaunchAgents/org.homebrew.ssh-agent.plist
```

Note that you'll lose any keys currently loaded into the agent. Re-run [`get-sk-ssh`](../scripts/security-keys.md#get-sk-ssh)
to load them again.

## See also

- [Git → Signing & security keys](../git/signing-security-keys.md)
- [Scripts → ssh-askpass](../scripts/security-keys.md#ssh-askpass)
