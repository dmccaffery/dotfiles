---
icon: simple/apple
---

# macOS system defaults

`setup/darwin/config.sh` (the `defaults` stage of `install.sh`) applies a curated set of
`defaults write` calls. Each is idempotent — re-running the stage doesn't churn state.

## Input & trackpad

```sh
defaults write -g ApplePressAndHoldEnabled -bool false
defaults write com.apple.AppleMultitouchTrackpad TrackpadThreeFingerDrag -bool true
```

- **Press-and-hold disabled** so key repeat works (essential for Vim).
- **Three-finger drag** for window dragging without click-and-hold.

## Menu bar

```sh
defaults write -g _HIHideMenuBar -bool false
defaults write -g AppleMenuBarVisibleInFullscreen -bool false
defaults write com.apple.controlcenter AutoHideMenuBarOption -int 2
killall SystemUIServer
```

Menu bar always visible _except_ in full-screen apps.

## Finder & desktop

```sh
defaults write com.apple.finder ShowExternalHardDrivesOnDesktop -bool false
defaults write com.apple.finder ShowHardDrivesOnDesktop -bool false
defaults write com.apple.finder ShowRemovableMediaOnDesktop -bool false
defaults write com.apple.finder ShowMountedServersOnDesktop -bool false
defaults write com.apple.finder CreateDesktop -bool false           # hide ALL desktop icons
defaults write com.apple.desktopservices DSDontWriteNetworkStores -bool true
defaults write com.apple.finder ShowPathbar -bool true
defaults write com.apple.finder ShowStatusBar -bool true
defaults write com.apple.finder FXEnableExtensionChangeWarning -bool false
```

- **Desktop icons hidden** entirely.
- **No `.DS_Store` on network volumes.**
- **Path bar + status bar** always shown.

## Screenshots

```sh
defaults write com.apple.screencapture type -string "png"
```

## Software update

```sh
defaults write com.apple.SoftwareUpdate ScheduleFrequency -int 1
```

Daily check.

## Spaces & Dock

```sh
defaults write com.apple.spaces "spans-displays" -bool false
defaults write com.apple.dock "mru-spaces" -bool false
defaults write com.apple.dock autohide -bool false
defaults write com.apple.dock largesize -float 96
defaults write com.apple.dock "minimize-to-application" -bool true
defaults write com.apple.dock tilesize -float 48
killall Dock
```

- Spaces are per-display, not spanning.
- Spaces don't reorder themselves.
- Dock visible, minimize-to-app, 48 px tiles → 96 px magnified.

## Window behaviour

```sh
defaults write -g NSWindowShouldDragOnGesture -bool true
defaults write -g NSAutomaticWindowAnimationsEnabled -bool false
```

- Drag windows by clicking anywhere with the appropriate gesture.
- No window animations — they slow everything down for no benefit.

## GPG & smartcard

```sh
defaults write org.gpgtools.common UseKeychain -bool yes
defaults write org.gpgtools.common DisableKeychain -bool no

sudo defaults write /Library/Preferences/com.apple.security.smartcard enforceSmartCard -bool true
sudo defaults write /Library/Preferences/com.apple.security.smartcard allowUnmappedUsers -int 1
```

The last two enforce smartcard authentication for the login user — paired with the YubiKey
setup. `allowUnmappedUsers = 1` keeps unmapped users from being locked out, important if you
test with multiple cards.
