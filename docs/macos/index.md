---
icon: simple/apple
---

# macOS

The macOS-specific bits: a curated set of `defaults write` calls applied by the installer, and
the launch agent that replaces Apple's `ssh-agent` with the Homebrew build.

| Page                                  | Purpose                                                             |
| ------------------------------------- | ------------------------------------------------------------------- |
| [System defaults](system-defaults.md) | The `defaults write` calls applied by `setup/darwin/config.sh`.     |
| [Launch agents](launchagents.md)      | `org.homebrew.ssh-agent` replacing Apple's `com.openssh.ssh-agent`. |
