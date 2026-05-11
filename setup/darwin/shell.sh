#! /usr/bin/env bash

# disable agents
launchctl disable gui/${UID}/com.openssh.ssh-agent 2> /dev/null || true
launchctl bootout gui/${UID}/org.homebrew.ssh-agent 2> /dev/null || true

# bootstrap homebrew ssh agent
launchctl bootstrap gui/${UID} ~/Library/LaunchAgents/org.homebrew.ssh-agent.plist
