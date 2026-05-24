#! /usr/bin/env sh

# shellcheck disable=SC3028
id=${UID:-$(id -u)}

# disable agents
launchctl disable gui/"${id}"/com.openssh.ssh-agent 2> /dev/null || true
launchctl bootout gui/"${id}"/org.homebrew.ssh-agent 2> /dev/null || true

# bootstrap homebrew ssh agent
launchctl bootstrap gui/"${id}" ~/Library/LaunchAgents/org.homebrew.ssh-agent.plist
