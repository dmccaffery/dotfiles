// Package sshsk implements the ssh-sk command: the YubiKey-resident SSH
// signing-key workflow (generate, load, publish, and resolve for git signing).
package sshsk

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/dmccaffery/dotfiles/internal/cmd/cmdutil"
	"github.com/dmccaffery/dotfiles/internal/ui"
)

// Internal sentinels — all surface to the user as cmdutil.ErrSilent after the
// relevant diagnostic is logged; callers only test them for nil.
var (
	errNoStubs      = errors.New("no saved security-key stubs")
	errNoSigningKey = errors.New("no signing key resolved")
)

type skCmd struct{ deps *cmdutil.Deps }

// account is one authenticated provider identity (GitHub or Forgejo).
type account struct {
	provider string // "github" | "forgejo"
	host     string
	user     string
}

// NewCmd builds the ssh-sk command.
func NewCmd(deps *cmdutil.Deps) *cobra.Command {
	s := &skCmd{deps: deps}
	cmd := &cobra.Command{
		Use:   "ssh-sk <command>",
		Short: "YubiKey-resident SSH signing-key workflow",
		Long: "Generates and loads resident ed25519-sk keys, publishes them to GitHub/Forgejo,\n" +
			"and resolves the configured git signing key for git's defaultKeyCommand.",
		// git reads `ssh-sk get --git` stdout for the signing key, so on error this
		// command must never spill usage/error text — keep stdout to the key line.
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	gen := &cobra.Command{
		Use:     "gen [user]",
		Short:   "Generate a resident ed25519-sk key and publish it when possible",
		Aliases: []string{"generate"},
		Args:    cobra.MaximumNArgs(1),
		RunE:    s.gen,
	}
	get := &cobra.Command{
		Use:   "get",
		Short: "Load saved stubs into the agent and update git allowed_signers",
		Args:  cobra.NoArgs,
		RunE:  s.get,
	}
	get.Flags().Bool("git", false, "print the configured git signing key as key::<blob> (defaultKeyCommand)")
	cmd.AddCommand(gen, get)
	return cmd
}

func (s *skCmd) keyRoot() string {
	return filepath.Join(s.deps.Env.XDGConfigHome(), "private", "ssh")
}

func (s *skCmd) allowedSigners() string {
	return filepath.Join(s.deps.Env.XDGConfigHome(), "private", "git", "allowed_signers")
}

// get implements `ssh-sk get [--git]`.
func (s *skCmd) get(cmd *cobra.Command, _ []string) error {
	if gitOnly, _ := cmd.Flags().GetBool("git"); gitOnly {
		line, err := s.resolveGit(cmd.Context(), "")
		if err != nil {
			return cmdutil.ErrSilent // git reads stdout; keep it empty on failure
		}
		fmt.Fprintf(cmd.OutOrStdout(), "key::%s\n", line)
		return nil
	}
	return s.setup(cmd, "")
}

// setup loads the inserted YubiKey's stubs into the agent and appends the
// resolved signing key to the git allowed_signers file.
func (s *skCmd) setup(cmd *cobra.Command, user string) error {
	ctx, log := cmd.Context(), s.deps.Log
	log.Info("setting up ssh agent...")

	serials := s.currentSerialsOrPrompt(ctx)
	if len(serials) == 0 {
		log.Error("no YubiKey is currently inserted")
		return cmdutil.ErrSilent
	}
	if err := s.loadCurrentStubs(ctx, serials); err != nil {
		log.Error("no saved security-key stubs found for the inserted YubiKey")
		return cmdutil.ErrSilent
	}
	log.Info("keys successfully added to agent!")

	log.Info("updating git allowed signers...")
	email := s.gitConfig(ctx, "user.email")
	if email == "" {
		log.Error("user email is not set")
		return cmdutil.ErrSilent
	}
	line, err := s.resolveGit(ctx, user)
	if err != nil {
		log.Error("user signing key is not set")
		return cmdutil.ErrSilent
	}
	return s.appendAllowedSigner(email, line)
}

// gen generates a resident ed25519-sk key namespaced by user, loads it, and
// publishes it to any matching authenticated account.
func (s *skCmd) gen(cmd *cobra.Command, args []string) error {
	ctx, log := cmd.Context(), s.deps.Log
	user := cmdutil.Arg(args, 0)

	if user == "" {
		picked, err := s.selectAuthenticatedUser(ctx)
		if err != nil {
			return err
		}
		user = picked
	}
	if !safeUser(user) {
		log.Error("user must be safe for a filename: " + user)
		return cmdutil.ErrSilent
	}

	matching := matchingAccounts(s.authenticatedAccounts(ctx), user)
	if len(matching) == 0 && !s.confirmUnpublishedGeneration(user) {
		return cmdutil.ErrSilent
	}

	serial, err := s.singleCurrentSerial(ctx)
	if err != nil {
		return err
	}
	keyDir := filepath.Join(s.keyRoot(), serial)
	keyFile := filepath.Join(keyDir, "id_ed25519_sk_"+user)

	if err := os.MkdirAll(keyDir, 0o700); err != nil {
		log.Error(err.Error())
		return cmdutil.ErrSilent
	}
	_ = os.Chmod(s.keyRoot(), 0o700)
	_ = os.Chmod(keyDir, 0o700)

	if pathExists(keyFile) || pathExists(keyFile+".pub") {
		log.Error("security-key stub already exists: " + keyFile)
		return cmdutil.ErrSilent
	}

	if err := s.deps.Runner.RunIO(ctx, cmdutil.Streams(cmd), "ssh-keygen",
		"-t", "ed25519-sk", "-f", keyFile,
		"-O", "resident", "-O", "verify-required", "-O", "no-touch-required",
		"-O", "application=ssh:"+user, "-O", "user="+user, "-C", user); err != nil {
		log.Error("failed to generate security key")
		return cmdutil.ErrSilent
	}
	if err := s.deps.Runner.RunIO(ctx, cmdutil.Streams(cmd), "ssh-add", keyFile); err != nil {
		log.Error("failed to add the generated key to the agent")
		return cmdutil.ErrSilent
	}

	publishFailed := false
	if len(matching) > 0 {
		publishFailed = !s.publishKeyToAccounts(cmd, keyFile+".pub", matching)
	}

	if err := s.setup(cmd, user); err != nil {
		return err
	}
	if publishFailed {
		return cmdutil.ErrSilent
	}
	return nil
}

// resolveGit returns the agent key line for user (or the configured account when
// user is ""), suitable for emitting as `key::<line>`.
func (s *skCmd) resolveGit(ctx context.Context, user string) (string, error) {
	if user == "" {
		acct, err := s.configuredGitAccount(ctx)
		if err != nil {
			return "", err
		}
		user = acct
	}
	if !safeUser(user) {
		s.deps.Log.Error("git signing account is not safe for a filename: " + user)
		return "", cmdutil.ErrSilent
	}
	return s.matchGitSigningKey(ctx, user)
}

// matchGitSigningKey loads the inserted YubiKey's stubs and returns the agent key
// line whose blob matches user's saved stub, prompting once to retry.
func (s *skCmd) matchGitSigningKey(ctx context.Context, user string) (string, error) {
	for attempt := 1; attempt <= 2; attempt++ {
		serials := s.currentSerialsOrPrompt(ctx)
		if len(serials) == 0 {
			return "", errNoSigningKey
		}
		_ = s.loadCurrentStubs(ctx, serials) // best-effort; matching is the real gate
		if blobs := s.currentUserKeyBlobs(user, serials); len(blobs) > 0 {
			if line := matchAgentKey(s.agentKeys(ctx), blobs); line != "" {
				return line, nil
			}
		}
		if attempt == 1 && s.promptForYubiKey("no inserted YubiKey has a saved signing key for "+user) {
			continue
		}
		return "", errNoSigningKey
	}
	return "", errNoSigningKey
}

// configuredGitAccount resolves the signing account from github.account /
// forgejo.account, erroring when they conflict or neither is set.
func (s *skCmd) configuredGitAccount(ctx context.Context) (string, error) {
	github := s.gitConfig(ctx, "github.account")
	forgejo := s.gitConfig(ctx, "forgejo.account")
	switch {
	case github != "" && forgejo != "":
		if github != forgejo {
			s.deps.Log.Error("github.account and forgejo.account are both set to different users")
			return "", cmdutil.ErrSilent
		}
		return github, nil
	case github != "":
		return github, nil
	case forgejo != "":
		return forgejo, nil
	default:
		s.deps.Log.Error("no git signing account configured; set github.account or forgejo.account")
		return "", cmdutil.ErrSilent
	}
}

// appendAllowedSigner adds `<email> namespaces="git" <publicKey>` to the
// allowed_signers file, skipping a write when the exact line already exists.
func (s *skCmd) appendAllowedSigner(email, publicKey string) error {
	path := s.allowedSigners()
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		s.deps.Log.Error(err.Error())
		return cmdutil.ErrSilent
	}
	line := signerLine(email, publicKey)

	existing, _ := os.ReadFile(path)
	for _, l := range strings.Split(string(existing), "\n") {
		if l == line {
			s.deps.Log.Info(email + ": " + publicKey + " already present in " + path)
			return nil
		}
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		s.deps.Log.Error(err.Error())
		return cmdutil.ErrSilent
	}
	defer f.Close()
	s.deps.Log.Info("adding " + email + ": " + publicKey + " to " + path)
	if _, err := fmt.Fprintln(f, line); err != nil {
		s.deps.Log.Error(err.Error())
		return cmdutil.ErrSilent
	}
	return nil
}

// publishKeyToAccounts uploads pubFile to each account, returning false if any
// upload fails (or there were none to publish to).
func (s *skCmd) publishKeyToAccounts(cmd *cobra.Command, pubFile string, accounts []account) bool {
	ctx, log, r := cmd.Context(), s.deps.Log, s.deps.Runner
	published, ok := false, true
	for _, a := range accounts {
		published = true
		switch a.provider {
		case "github":
			log.Info("publishing signing key to GitHub as " + a.user + "...")
			_, switchErr := r.Run(ctx, "gh", "auth", "switch", "--hostname", a.host, "--user", a.user)
			if switchErr == nil && r.RunIO(ctx, cmdutil.Streams(cmd), "gh", "ssh-key", "add", pubFile,
				"--type", "signing", "--title", "YubiKey signing") == nil {
				log.Info("published signing key to GitHub as " + a.user)
			} else {
				log.Error("failed to publish signing key to GitHub as " + a.user)
				ok = false
			}
		case "forgejo":
			log.Info("publishing SSH key to Forgejo " + a.host + " as " + a.user + "...")
			if r.RunIO(ctx, cmdutil.Streams(cmd), "fj", "-H", a.host, "user", "key", "upload", pubFile,
				"--title", "YubiKey signing") == nil {
				log.Info("published SSH key to Forgejo " + a.host + " as " + a.user)
			} else {
				log.Error("failed to publish SSH key to Forgejo " + a.host + " as " + a.user)
				ok = false
			}
		}
	}
	return published && ok
}

// selectAuthenticatedUser picks one unique authenticated username (huh replaces
// the fzf picker; a sole account is auto-selected).
func (s *skCmd) selectAuthenticatedUser(ctx context.Context) (string, error) {
	users := uniqueUsers(s.authenticatedAccounts(ctx))
	if len(users) == 0 {
		s.deps.Log.Error("no authenticated GitHub or Forgejo accounts found")
		return "", cmdutil.ErrSilent
	}
	selected, err := cmdutil.PickOne(s.deps.Prompt, "select git account", "", users)
	if err != nil {
		if errors.Is(err, ui.ErrNoTTY) {
			s.deps.Log.Error("no tty available to select an account")
		}
		return "", cmdutil.ErrSilent
	}
	if selected == "" {
		return "", cmdutil.ErrSilent // aborted
	}
	return selected, nil
}

// authenticatedAccounts lists every GitHub (gh) and Forgejo (fj) identity logged
// in on this machine.
func (s *skCmd) authenticatedAccounts(ctx context.Context) []account {
	var accounts []account
	r := s.deps.Runner

	if _, err := r.Look("gh"); err == nil {
		res, err := r.Run(ctx, "gh", "auth", "status",
			"--hostname", "github.com", "--json", "hosts",
			"--jq", `.hosts["github.com"] // [] | .[].login`)
		if err == nil {
			for _, login := range cmdutil.NonEmptyLines(res.Stdout) {
				accounts = append(accounts, account{provider: "github", host: "github.com", user: strings.Fields(login)[0]})
			}
		}
	}
	if _, err := r.Look("fj"); err == nil {
		res, err := r.Run(ctx, "fj", "--style", "minimal", "auth", "list")
		if err == nil {
			accounts = append(accounts, parseForgejoAccounts(res.Stdout)...)
		}
	}
	return accounts
}

// confirmUnpublishedGeneration warns and asks whether to generate a key that
// won't be published; false (incl. no tty) aborts.
func (s *skCmd) confirmUnpublishedGeneration(user string) bool {
	s.deps.Log.Warn("no authenticated GitHub or Forgejo account found for " + user + "; key will not be published")
	ok, err := s.deps.Prompt.Confirm("continue generating an unpublished key?", false)
	if err != nil {
		if errors.Is(err, ui.ErrNoTTY) {
			s.deps.Log.Error("cannot confirm unpublished key generation without a tty")
		}
		return false
	}
	return ok
}

// currentSerials lists the serials of the inserted YubiKeys (empty if none).
func (s *skCmd) currentSerials(ctx context.Context) []string {
	res, err := s.deps.Runner.Run(ctx, "ykman", "list", "--serials")
	if err != nil {
		return nil
	}
	var serials []string
	for _, line := range cmdutil.NonEmptyLines(res.Stdout) {
		serials = append(serials, strings.Fields(line)[0])
	}
	return serials
}

// currentSerialsOrPrompt returns inserted serials, prompting once to insert a
// YubiKey when none are present.
func (s *skCmd) currentSerialsOrPrompt(ctx context.Context) []string {
	if serials := s.currentSerials(ctx); len(serials) > 0 {
		return serials
	}
	if s.promptForYubiKey("no YubiKey is currently inserted") {
		if serials := s.currentSerials(ctx); len(serials) > 0 {
			return serials
		}
	}
	return nil
}

// singleCurrentSerial requires exactly one inserted YubiKey (for generation).
func (s *skCmd) singleCurrentSerial(ctx context.Context) (string, error) {
	serials := s.currentSerialsOrPrompt(ctx)
	switch len(serials) {
	case 1:
		return serials[0], nil
	case 0:
		s.deps.Log.Error("no YubiKey is currently inserted")
	default:
		s.deps.Log.Error("insert exactly one YubiKey when generating a key")
	}
	return "", cmdutil.ErrSilent
}

// currentUserKeyBlobs reads the public-key blobs of user's saved stubs across the
// given serials.
func (s *skCmd) currentUserKeyBlobs(user string, serials []string) []string {
	var blobs []string
	for _, serial := range serials {
		pub := filepath.Join(s.keyRoot(), serial, "id_ed25519_sk_"+user+".pub")
		data, err := os.ReadFile(pub)
		if err != nil {
			continue
		}
		for _, line := range cmdutil.NonEmptyLines(string(data)) {
			if fields := strings.Fields(line); len(fields) >= 2 {
				blobs = append(blobs, fields[1])
			}
		}
	}
	return blobs
}

// loadCurrentStubs adds every non-.pub id_ed25519_sk_* stub for the given serials
// into the agent (stdin is /dev/null so ssh-add relies on the askpass/pinentry).
func (s *skCmd) loadCurrentStubs(ctx context.Context, serials []string) error {
	found, loaded := false, false
	for _, serial := range serials {
		keyDir := filepath.Join(s.keyRoot(), serial)
		if !cmdutil.DirExists(keyDir) {
			continue
		}
		matches, _ := filepath.Glob(filepath.Join(keyDir, "id_ed25519_sk_*"))
		for _, keyFile := range matches {
			if strings.HasSuffix(keyFile, ".pub") || !cmdutil.FileExists(keyFile) {
				continue
			}
			found = true
			if _, err := s.deps.Runner.Run(ctx, "ssh-add", keyFile); err == nil {
				loaded = true
			}
		}
	}
	if !found || !loaded {
		return errNoStubs
	}
	return nil
}

// agentKeys returns the current `ssh-add -L` listing (empty on error).
func (s *skCmd) agentKeys(ctx context.Context) string {
	res, err := s.deps.Runner.Run(ctx, "ssh-add", "-L")
	if err != nil {
		return ""
	}
	return res.Stdout
}

// gitConfig returns `git config --get <key>` trimmed (empty when unset).
func (s *skCmd) gitConfig(ctx context.Context, key string) string {
	res, err := s.deps.Runner.Run(ctx, "git", "config", "--get", key)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(res.Stdout)
}

// promptForYubiKey asks whether to insert a YubiKey and retry; false (incl. no
// tty) stops the retry.
func (s *skCmd) promptForYubiKey(message string) bool {
	ok, err := s.deps.Prompt.Confirm(message+"; insert a YubiKey and retry?", true)
	if err != nil {
		if errors.Is(err, ui.ErrNoTTY) {
			s.deps.Log.Warn(message)
		}
		return false
	}
	return ok
}

// --- pure helpers (unit-tested directly) ---

// safeUser reports whether u is safe to embed in a stub filename.
func safeUser(u string) bool {
	if u == "" || u == "." || u == ".." {
		return false
	}
	for _, r := range u {
		switch {
		case r >= 'A' && r <= 'Z', r >= 'a' && r <= 'z', r >= '0' && r <= '9':
		case r == '.' || r == '_' || r == '@' || r == '+' || r == '-':
		default:
			return false
		}
	}
	return true
}

// signerLine renders an allowed_signers entry for the git namespace.
func signerLine(email, publicKey string) string {
	return email + ` namespaces="git" ` + publicKey
}

// parseForgejoAccounts maps `fj --style minimal auth list` lines (user@host …)
// to forgejo accounts, splitting on the first '@'.
func parseForgejoAccounts(out string) []account {
	var accounts []account
	for _, line := range cmdutil.NonEmptyLines(out) {
		first := strings.Fields(line)[0]
		at := strings.Index(first, "@")
		if at <= 0 || at == len(first)-1 {
			continue
		}
		accounts = append(accounts, account{provider: "forgejo", host: first[at+1:], user: first[:at]})
	}
	return accounts
}

// matchAgentKey returns the first `ssh-add -L` line whose key blob (field 2) is
// in blobs, or "" when none match.
func matchAgentKey(agentKeys string, blobs []string) string {
	set := make(map[string]bool, len(blobs))
	for _, b := range blobs {
		if b != "" {
			set[b] = true
		}
	}
	for _, line := range cmdutil.NonEmptyLines(agentKeys) {
		if fields := strings.Fields(line); len(fields) >= 2 && set[fields[1]] {
			return line
		}
	}
	return ""
}

// uniqueUsers returns the distinct, in-order usernames across accounts.
func uniqueUsers(accounts []account) []string {
	var users []string
	seen := map[string]bool{}
	for _, a := range accounts {
		if a.user != "" && !seen[a.user] {
			seen[a.user] = true
			users = append(users, a.user)
		}
	}
	return users
}

// matchingAccounts returns the accounts whose user equals user.
func matchingAccounts(accounts []account, user string) []account {
	var out []account
	for _, a := range accounts {
		if a.user == user {
			out = append(out, a)
		}
	}
	return out
}

// pathExists reports whether path exists (of any type), like the shell's `-e`.
func pathExists(p string) bool {
	_, err := os.Lstat(p)
	return err == nil
}
