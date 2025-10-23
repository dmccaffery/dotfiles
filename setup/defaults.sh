#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
. "${SCRIPT_DIR}/printing.sh"

info "applying macOS system defaults"

# enable key repeats
defaults write -g ApplePressAndHoldEnabled -bool false

# enable three finger drag
defaults write com.apple.AppleMultitouchTrackpad TrackpadThreeFingerDrag -bool true

# disable prompting to use new external drives for time machine
defaults write com.apple.TimeMachine DoNotOfferNewDisksForBackup -bool true

# hide external hard drives on desktop
defaults write com.apple.finder ShowExternalHardDrivesOnDesktop -bool false

# hide hard drives on desktop
defaults write com.apple.finder ShowHardDrivesOnDesktop -bool false

# hide removable media hard drives on desktop
defaults write com.apple.finder ShowRemovableMediaOnDesktop -bool false

# hide mounted servers on desktop
defaults write com.apple.finder ShowMountedServersOnDesktop -bool false

# hide icons on desktop
defaults write com.apple.finder CreateDesktop -bool false

# avoid creating .DS_Store files on network volumes
defaults write com.apple.desktopservices DSDontWriteNetworkStores -bool true

# show path bar
defaults write com.apple.finder ShowPathbar -bool true

# hide status bar
defaults write com.apple.finder "ShowStatusBar" -bool false

# do not show warning when changing the file extension
defaults write com.apple.finder FXEnableExtensionChangeWarning -bool false

# save screenshots in png format (other options: BMP, GIF, JPG, PDF, TIFF)
defaults write com.apple.screencapture type -string "png"

# set weekly software update checks
defaults write com.apple.SoftwareUpdate ScheduleFrequency -int 7

# spaces span all displays
defaults write com.apple.spaces "spans-displays" -bool false

# do not rearrange spaces automatically
defaults write com.apple.dock "mru-spaces" -bool false

# setup dock to autohide
defaults write com.apple.dock autohide -bool true
defaults write com.apple.dock largesize -float 128
defaults write com.apple.dock "minimize-to-application" -bool true
defaults write com.apple.dock tilesize -float 32

# drag windows on gesture
defaults write -g NSWindowShouldDragOnGesture -bool true

# disable window animations
defaults write -g NSAutomaticWindowAnimationsEnabled -bool false
