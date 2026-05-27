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
        <key>EnvironmentVariables</key>
        <dict>
            <key>SSH_ASKPASS</key>
            <string>/usr/local/bin/ssh-askpass</string>
            <key>SSH_ASKPASS_REQUIRE</key>
            <string>force</string>
            <key>DISPLAY</key>
            <string>:0</string>
            <key>SSH_AUTH_SOCK_LOCAL</key>
            <string>/tmp/ssh-agent.sock</string>
        </dict>
        <key>ProgramArguments</key>
        <array>
            <string>/bin/sh</string>
            <string>-c</string>
            <string>rm -f ${SSH_AUTH_SOCK}; killall ssh-agent; ln -fs ${SSH_AUTH_SOCK_LOCAL} ${SSH_AUTH_SOCK}; /opt/homebrew/bin/ssh-agent -D -a ${SSH_AUTH_SOCK_LOCAL};</string>
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

The `EnvironmentVariables` block exports values into the agent process and into every
shell descended from this launchd job:

| Key                   | Value                        | Purpose                                                                                                  |
| --------------------- | ---------------------------- | -------------------------------------------------------------------------------------------------------- |
| `SSH_ASKPASS`         | `/usr/local/bin/ssh-askpass` | The [ssh-askpass wrapper](../scripts/security-keys.md#ssh-askpass) that hands prompts to `pinentry-mac`. |
| `SSH_ASKPASS_REQUIRE` | `force`                      | OpenSSH 8.4+ flag that **always** uses `SSH_ASKPASS`, even when a TTY is attached.                       |
| `DISPLAY`             | `:0`                         | Set only to satisfy askpass programs that bail out without a display.                                    |
| `SSH_AUTH_SOCK_LOCAL` | `/tmp/ssh-agent.sock`        | The **stable** socket path that the agent listens on.                                                    |

The shell command:

1. Removes whatever socket launchd preallocated at `${SSH_AUTH_SOCK}` (the random
   per-boot `/var/run/com.apple.launchd.<random>/Listeners` path).
2. `killall ssh-agent` to clear any stale agent left behind by a previous boot or reload.
3. Symlinks the random launchd path → the stable path
   (`ln -fs ${SSH_AUTH_SOCK_LOCAL} ${SSH_AUTH_SOCK}`). Shells inherit the launchd
   `SSH_AUTH_SOCK`, follow the symlink, and end up talking to the agent on the stable
   socket. The stable path is also what the Claude Code sandbox allowlists — see
   [Claude → settings](../claude/settings.md#sandbox).
4. Runs Homebrew's `ssh-agent` in the foreground (`-D`) listening on the stable socket (`-a`).

> **Restart shells after reloading.** Already-running terminals keep their old
> `SSH_AUTH_SOCK`; only processes launched after the agent reloads pick up the new value.

## How it's installed

The `config` stage of `install.sh` (`setup/darwin/config.sh`) bootstraps it:

```sh title="setup/darwin/config.sh"
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
