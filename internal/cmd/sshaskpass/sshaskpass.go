// Package sshaskpass implements the ssh-askpass command: the SSH_ASKPASS bridge
// that hands ssh-agent's PIN prompts to pinentry-mac.
package sshaskpass

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"

	"github.com/dmccaffery/dotfiles/internal/cmd/cmdutil"
	"github.com/dmccaffery/dotfiles/internal/execx"
)

// pinentryPath is the pinentry-mac binary; a var so tests can point elsewhere.
// SSH_ASKPASS runs in a minimal environment, so the path is absolute rather than
// resolved on PATH.
var pinentryPath = "/opt/homebrew/bin/pinentry-mac"

type askCmd struct{ deps *cmdutil.Deps }

// NewCmd builds the ssh-askpass command.
func NewCmd(deps *cmdutil.Deps) *cobra.Command {
	a := &askCmd{deps: deps}
	return &cobra.Command{
		Use:   "ssh-askpass [prompt]",
		Short: "Bridge ssh-agent SSH_ASKPASS prompts to pinentry-mac",
		Long: "Invoked by the launch-managed ssh-agent as $SSH_ASKPASS. It acknowledges\n" +
			"FIDO2 user-presence prompts and forwards PIN prompts to pinentry-mac.",
		// The agent passes a single quoted prompt that may begin with '-'; take the
		// argument verbatim rather than parsing it as flags.
		DisableFlagParsing: true,
		RunE:               a.run,
	}
}

func (a *askCmd) run(cmd *cobra.Command, args []string) error {
	out := cmd.OutOrStdout()
	desc := cmdutil.Arg(args, 0)

	// FIDO2 user-presence: the agent only needs a non-error reply, not a PIN.
	if strings.HasPrefix(desc, "Confirm user presence") {
		fmt.Fprintln(out)
		return nil
	}

	var buf bytes.Buffer
	// Ignore pinentry's exit status (the shell original has no `set -e` here): a
	// cancel just yields no data line and an empty PIN, which the agent treats as
	// a failed unlock. ssh-askpass itself always exits 0.
	_ = a.deps.Runner.RunIO(cmd.Context(),
		execx.Streams{In: strings.NewReader(buildAssuan(desc)), Out: &buf, Err: io.Discard},
		pinentryPath)
	fmt.Fprintln(out, extractPin(buf.String()))
	return nil
}

// buildAssuan renders the pinentry Assuan command block for a PIN prompt. When
// the prompt carries a SHA256 fingerprint it is passed through as SETKEYINFO so
// pinentry can show — and cache the PIN for — the specific key being unlocked.
func buildAssuan(desc string) string {
	hashType, rest, _ := strings.Cut(desc, ":")
	if strings.HasSuffix(hashType, "SHA256") {
		sha, _, _ := strings.Cut(rest, ":")
		return "SETDESC " + desc + "\n" +
			"OPTION allow-external-password-cache\n" +
			"SETKEYINFO s/" + sha + "\n" +
			"GETPIN\n"
	}
	return "SETDESC " + desc + "\nGETPIN\n"
}

// extractPin pulls the PIN out of pinentry's Assuan response: the data lines
// (those containing "D"), joined and stripped of the leading "D " marker.
func extractPin(out string) string {
	var b strings.Builder
	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, "D") {
			b.WriteString(line)
		}
	}
	joined := b.String()
	if len(joined) < 2 {
		return ""
	}
	return joined[2:]
}
