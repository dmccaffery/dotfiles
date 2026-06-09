// Package backup implements the backup command: it clears the entries in $HOME
// that would otherwise make `stow` refuse to link. For every regular file under
// the stow tree it inspects the matching $HOME path and moves real files into a
// timestamped backup, copies the contents of valid symlinks in before removing
// them, and deletes broken symlinks outright. It replaces the old backup.sh.
package backup

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/dmccaffery/dotfiles/internal/cmd/cmdutil"
	"github.com/dmccaffery/dotfiles/internal/logx"
	"github.com/dmccaffery/dotfiles/internal/ui"
)

type backupCmd struct {
	deps      *cmdutil.Deps
	stowDir   string
	backupDir string
	yes       bool
}

// NewCmd builds the backup command. It is registered as a plain `dot backup`
// subcommand (not an applet), so no standalone `backup` symlink is created.
func NewCmd(deps *cmdutil.Deps) *cobra.Command {
	bc := &backupCmd{deps: deps}
	cmd := &cobra.Command{
		Use:   "backup [name]",
		Short: "Move conflicting $HOME entries into a backup so stow can link cleanly",
		Long: `Clear the entries in $HOME that would make stow refuse to overwrite, so the next
` + "`make stow`" + ` links cleanly. For every regular file under the stow tree (default
` + "`./stow`" + `, override with --stow-dir) the matching path under $HOME is handled:

  - a regular file        -> moved into the backup
  - a valid symlink       -> its resolved contents are copied into the backup, then
                             the symlink is removed
  - a broken symlink      -> removed (nothing to back up)
  - a symlinked directory -> handled at the link itself, never descended through

Real files are moved (not copied), so the originals leave $HOME the moment this runs.
The backup is written to <repo>/backups/<name> (name defaults to a timestamp), the same
layout restore.sh reads, so a backup round-trips.`,
		Args: cobra.MaximumNArgs(1),
		RunE: bc.run,
	}
	cmd.Flags().StringVar(&bc.stowDir, "stow-dir", "stow", "directory of stow trees to mirror against $HOME")
	cmd.Flags().StringVar(&bc.backupDir, "backup-dir", "", "destination directory (default <repo>/backups/<name>)")
	cmd.Flags().BoolVarP(&bc.yes, "yes", "y", false, "skip the confirmation prompt")
	return cmd
}

func (b *backupCmd) run(cmd *cobra.Command, args []string) error {
	log := b.deps.Log

	stowDir, err := filepath.Abs(b.stowDir)
	if err != nil {
		log.Error("resolving stow dir: " + err.Error())
		return cmdutil.ErrSilent
	}
	if fi, err := os.Stat(stowDir); err != nil || !fi.IsDir() {
		log.Error("stow dir not found: " + stowDir)
		return cmdutil.ErrSilent
	}

	name := cmdutil.Arg(args, 0)
	if name == "" {
		name = time.Now().Format("2006-01-02-150405")
	}
	backupDir := b.backupDir
	if backupDir == "" {
		backupDir = filepath.Join(filepath.Dir(stowDir), "backups", name)
	}
	if backupDir, err = filepath.Abs(backupDir); err != nil {
		log.Error("resolving backup dir: " + err.Error())
		return cmdutil.ErrSilent
	}

	home := b.deps.Env.Home()
	if home == "" {
		log.Error("could not determine home directory")
		return cmdutil.ErrSilent
	}

	if !b.yes {
		prompt := fmt.Sprintf("Back up conflicting entries from %s into %s? Real files are moved; symlinks are removed.",
			home, backupDir)
		ok, err := b.deps.Prompt.Confirm(prompt, false)
		switch {
		case errors.Is(err, ui.ErrNoTTY):
			log.Warn("no terminal to confirm on; re-run with --yes to proceed")
			return nil
		case err != nil:
			log.Error(err.Error())
			return cmdutil.ErrSilent
		case !ok:
			log.Warn("aborted")
			return nil
		}
	}

	res, err := backupTargets(home, stowDir, backupDir, log)
	if err != nil {
		log.Error(err.Error())
		return cmdutil.ErrSilent
	}

	if res.total() == 0 {
		log.Info("nothing to back up; no conflicting entries in " + home)
		return nil
	}
	log.Info(fmt.Sprintf("backed up %d file(s) and removed %d symlink(s) (%d broken) -> %s",
		res.filesMoved, res.symlinksRemoved, res.brokenRemoved, backupDir))
	return nil
}

// result tallies what backupTargets did, for the summary line.
type result struct {
	filesMoved      int // regular files moved into the backup
	symlinksRemoved int // symlinks removed (valid contents-backed-up + broken)
	brokenRemoved   int // subset of symlinksRemoved that pointed nowhere
}

func (r result) total() int { return r.filesMoved + r.symlinksRemoved }

// backupTargets walks stowDir for every entry stow would link and clears each
// matching $HOME path. The walk over $HOME for a given entry stops at the first
// symlink it meets (so a folded parent like ~/.config/x -> repo is handled at
// the link, never descended into and never mistaken for the repo's own file).
func backupTargets(home, stowDir, backupDir string, log *logx.Logger) (result, error) {
	var res result
	err := filepath.WalkDir(stowDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Process every non-directory entry stow links: regular files and the
		// symlinks the source itself holds (e.g. .config/{codex,opencode}/AGENTS.md
		// -> ../../.claude/CLAUDE.md). Skipping symlinks would leave their $HOME
		// targets in place to collide with stow.
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(stowDir, path)
		if err != nil {
			return err
		}
		return clearTarget(home, rel, backupDir, log, &res)
	})
	return res, err
}

// clearTarget walks the components of rel beneath home and acts on the first
// node that exists: a symlink (handled by handleSymlink), the leaf regular file
// (moved), or a regular file that stands where a directory is expected (also
// moved, since it blocks the subtree). A missing component means nothing is at
// the target, so there is nothing to do.
func clearTarget(home, rel, backupDir string, log *logx.Logger, res *result) error {
	parts := strings.Split(rel, string(os.PathSeparator))
	cur := home
	for i, p := range parts {
		cur = filepath.Join(cur, p)
		relSoFar := filepath.Join(parts[:i+1]...)

		fi, err := os.Lstat(cur)
		if err != nil {
			return nil // nothing exists at this point in the target path
		}

		switch {
		case fi.Mode()&os.ModeSymlink != 0:
			return handleSymlink(cur, relSoFar, backupDir, log, res)
		case fi.Mode().IsRegular():
			// Leaf file, or a file blocking a directory the subtree needs.
			return moveFile(cur, filepath.Join(backupDir, relSoFar), relSoFar, log, res)
		case fi.IsDir() && i < len(parts)-1:
			continue // descend through a real directory
		default:
			// A directory where a file is expected, or a special file: leave it
			// for stow to report rather than guessing.
			log.Warn(fmt.Sprintf("skipping %s: unexpected %s at target", relSoFar, fi.Mode().Type()))
			return nil
		}
	}
	return nil
}

// handleSymlink backs up and removes a symlink found at link. A broken link is
// just removed; a valid one has its resolved contents copied into the backup
// before the link itself (not its target) is removed.
func handleSymlink(link, rel, backupDir string, log *logx.Logger, res *result) error {
	if _, err := os.Stat(link); err != nil {
		if err := os.Remove(link); err != nil {
			return fmt.Errorf("removing broken symlink %s: %w", rel, err)
		}
		log.Info("removed broken symlink: " + rel)
		res.symlinksRemoved++
		res.brokenRemoved++
		return nil
	}

	resolved, err := filepath.EvalSymlinks(link)
	if err != nil {
		return fmt.Errorf("resolving symlink %s: %w", rel, err)
	}
	if err := copyPath(resolved, filepath.Join(backupDir, rel)); err != nil {
		return fmt.Errorf("backing up %s: %w", rel, err)
	}
	if err := os.Remove(link); err != nil {
		return fmt.Errorf("removing symlink %s: %w", rel, err)
	}
	log.Info("backed up and removed symlink: " + rel)
	res.symlinksRemoved++
	return nil
}

// moveFile relocates src into the backup at dest, falling back to copy+remove
// across filesystem boundaries (os.Rename's EXDEV).
func moveFile(src, dest, rel string, log *logx.Logger, res *result) error {
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	if err := os.Rename(src, dest); err != nil {
		if err := copyPath(src, dest); err != nil {
			return fmt.Errorf("moving %s: %w", rel, err)
		}
		if err := os.RemoveAll(src); err != nil {
			return fmt.Errorf("removing %s after copy: %w", rel, err)
		}
	}
	log.Info("backed up file: " + rel)
	res.filesMoved++
	return nil
}

// copyPath copies a file or directory tree from src to dest, preserving file
// permissions. Symlinks within a copied tree are recreated as symlinks.
func copyPath(src, dest string) error {
	fi, err := os.Lstat(src)
	if err != nil {
		return err
	}
	switch {
	case fi.IsDir():
		return copyDir(src, dest)
	case fi.Mode()&os.ModeSymlink != 0:
		target, err := os.Readlink(src)
		if err != nil {
			return err
		}
		return os.Symlink(target, dest)
	default:
		return copyFile(src, dest, fi.Mode().Perm())
	}
}

func copyDir(src, dest string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dest, 0o755); err != nil {
		return err
	}
	for _, e := range entries {
		if err := copyPath(filepath.Join(src, e.Name()), filepath.Join(dest, e.Name())); err != nil {
			return err
		}
	}
	return nil
}

func copyFile(src, dest string, perm os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	out, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}
