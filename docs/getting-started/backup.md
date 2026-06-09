---
icon: lucide/archive
---

# Backup

`dot backup` clears the entries in `$HOME` that would make `stow` refuse to overwrite, so the next
`make stow` links cleanly. It is the Go replacement for the old `backup.sh` — a `dot` subcommand
(run as `dot backup`), deliberately **not** an applet, so there is no standalone `backup` symlink
on `PATH`.

## Usage

```sh
make backup                  # go run ./cmd/dot backup (prompts first)
dot backup                   # backups/YYYY-MM-DD-HHMMSS/    (once the binary is installed)
dot backup pre-nvim-bump     # backups/pre-nvim-bump/
dot backup --yes             # skip the confirmation
dot backup --stow-dir stow   # which tree decides the targets to clear (default: ./stow)
```

`dot backup` walks every file under the stow tree (`./stow` by default) and, for each, inspects
the matching path under `$HOME` (the target root is always `$HOME`). What it does depends on what
is there:

| At `$HOME/<path>`     | Action                                                                         |
| --------------------- | ------------------------------------------------------------------------------ |
| a regular file        | **moved** into the backup                                                      |
| a valid symlink       | its resolved contents are **copied** into the backup, then the link is removed |
| a broken symlink      | **removed** (nothing to back up)                                               |
| a symlinked directory | handled at the link itself — never descended into (see below)                  |
| nothing               | skipped                                                                        |

The backup is written to `<repo>/backups/<name>/<path>` — `<name>` defaults to a timestamp
(`YYYY-MM-DD-HHMMSS`) or the optional first argument — mirroring the `$HOME`-relative layout so
[`restore.sh`](#restoring-a-backup) can copy it straight back. The `backups/` directory is
gitignored, so nothing is committed.

Because broken symlinks are deleted, `dot backup` is also the cleanup step after the stow trees
change location: when the sources moved under `stow/`, every previously-stowed link in `$HOME` was
left pointing at the old path. `dot backup` removes those orphaned links (and backs up any real
files), clearing the way for `make stow`.

!!! note "Requires Go"

    Unlike the old shell script, `dot backup` is part of the `dot` Go CLI, so `make backup` runs
    it with `go run ./cmd/dot backup` and needs `go` on `PATH` (it ships in the base Brewfile /
    `make requirements`). Once `make build` has linked the binary, `dot backup` runs directly too.

## A folded parent is handled at the link

`stow` _folds_ a tree into a single symlink when the target directory doesn't exist yet — e.g.
`~/.config/ghostty` may be one symlink into the repo rather than a real directory of per-file
links. Walking down a target path, `dot backup` stops at the **first** symlink it meets, so it
backs up and removes that folded parent rather than descending through it into the repo's own
files. [`setup/stow.sh`](../../setup/stow.sh) pre-creates the trees that hold runtime state
(`~/.claude`, `~/.config/{codex,opencode,zsh}`, `~/.ssh`, `~/.local/share/wallpapers`) as real
directories so stow links their children individually; the rest are free to fold.

## Restoring a backup

```sh
./restore.sh                # interactive fzf picker (or: make restore)
./restore.sh pre-nvim-bump  # pre-fills the fzf query; auto-selects on a single match
./restore.sh 2026-05        # fuzzy-match against the timestamp layout
```

[`restore.sh`](../../restore.sh) inverts `dot backup`:

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
