---
icon: lucide/archive
---

# Backup

`backup.sh` moves existing dotfile entries out of `$HOME` so that `stow` won't refuse to
symlink. It is intentionally minimal — no copying, no compression, no remote upload.

## Usage

```sh
./backup.sh                # backups/YYYY-MM-DD-HHMMSS/
./backup.sh pre-nvim-bump  # backups/pre-nvim-bump/
```

Creates `backups/<name>/` inside the repo and moves only the specific files and
subdirectories that this repo will symlink. Anything else under `$HOME` is left alone.

`<name>` defaults to a timestamp (`YYYY-MM-DD-HHMMSS`). Pass any name as the first
argument to label the backup — handy when you're snapshotting before a known-risky change.

Entries (skipped if absent or already a symlink):

- `~/.claude/settings.json`, `~/.claude/themes`
- every top-level entry under `~/.config/` that this repo also ships (iterated from the repo's
  own `stow/.config/` listing — i.e. `~/.config/ghostty`, `~/.config/nvim`, `~/.config/zsh`, …)
- `~/.local/share/scripts`, `~/.local/share/wallpapers`
- `~/.ssh/rc`
- `~/.terminfo/67/ghostty`, `~/.terminfo/78/xterm-ghostty`
- `~/Library/LaunchAgents/org.homebrew.ssh-agent.plist`
- `~/.zshrc`

This is deliberately narrower than a wholesale `~/.config` backup. [`setup/stow.sh`](../../setup/stow.sh)
pre-creates `~/.claude`, `~/.config`, and `~/.ssh` as real directories so that `stow` folds
children in rather than replacing the whole directory — keeping these parents intact means apps
that write into them (e.g. Claude Code's `.claude/projects/`, your own `.ssh/known_hosts`) don't
get displaced. Files (`settings.json`, `.ssh/rc`) get individual parent directories created in
the backup tree so paths round-trip if you ever want to restore them.

The `backups/` directory is already ignored by `.gitignore`, so nothing committed.

## What it does NOT do

- **Doesn't follow symlinks.** If an entry is already a symlink (e.g., from a prior stow run),
  it's skipped — stow can replace it cleanly.
- **Doesn't copy.** Files are _moved_ (`mv`), not copied. The originals are gone from `$HOME`
  the moment the script runs.

## Restoring a backup

```sh
./restore.sh                # interactive fzf picker (or: make restore)
./restore.sh pre-nvim-bump  # pre-fills the fzf query; auto-selects on a single match
./restore.sh 2026-05        # fuzzy-match against the timestamp layout
```

[`restore.sh`](../../restore.sh) inverts `backup.sh`:

1. Lists every directory under `./backups/` via `fzf` (newest first) so you can pick one.
   An optional first argument is passed to `fzf --query` for fuzzy filtering; combined
   with `--select-1`, a unique match skips the picker entirely.
2. Asks for confirmation, then runs `stow -D` over each stowed tree (`stow/.config`, `stow/.claude`,
   `stow/.local`, …) to remove the existing symlink layer.
3. Walks the chosen backup tree and **aborts** if any entry would collide with something
   still present in `$HOME` (after the unstow). Nothing has been written at this point —
   resolve the listed collisions manually and re-run.
4. Copies the backup tree back into `$HOME` (`cp -Pp`, preserving symlinks, mode, and
   timestamps). The backup folder is **preserved** so you can re-restore or compare.
5. Offers to delete the backup folder at the end.

Requires `fzf` on `PATH` (installed by the default Brewfile).
