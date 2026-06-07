// Package brewfile implements the brewfile command.
package brewfile

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	bf "github.com/dmccaffery/dotfiles/internal/brewfile"
	"github.com/dmccaffery/dotfiles/internal/cmd/cmdutil"
	"github.com/dmccaffery/dotfiles/internal/ui"
)

type brewfileCmd struct{ deps *cmdutil.Deps }

// NewCmd builds the brewfile command.
func NewCmd(deps *cmdutil.Deps) *cobra.Command {
	bc := &brewfileCmd{deps: deps}
	return &cobra.Command{
		Use:   "brewfile <add|remove> <package>... [brew bundle flags...]",
		Short: "Sync a package into/out of the global Brewfile and installed state",
		Long: `Sync a package into or out of the global Brewfile and the installed state in one
command: "brew bundle <add|remove> --global", then "brew bundle install --global".

Flags after the action pass through to brew bundle, so "brewfile add --cask ghostty"
and "brewfile add --tap user/tap" work. On add, any non-official tap reference
("user/tap/...") is checked against trust.json and you are prompted to brew trust it
before the install runs.`,
		Example: `  brewfile add jq
  brewfile add --cask ghostty
  brewfile remove jq`,
		// brew bundle flags (--cask, --tap, …) must reach brew verbatim, so Cobra's
		// flag parsing — and with it the automatic --help — is disabled; run() routes
		// -h/--help to cmd.Help(), which renders Use/Long/Example above.
		DisableFlagParsing: true,
		RunE:               bc.run,
	}
}

func (b *brewfileCmd) run(cmd *cobra.Command, args []string) error {
	log := b.deps.Log

	switch {
	case len(args) == 0:
		_ = cmd.Usage()
		return cmdutil.ErrSilent
	case args[0] == "-h", args[0] == "--help", args[0] == "help":
		return cmd.Help()
	}

	action := args[0]
	rest := args[1:]
	if action != "add" && action != "remove" {
		log.Error("unknown action: " + action)
		_ = cmd.Usage()
		return cmdutil.ErrSilent
	}
	if len(rest) == 0 {
		log.Error("missing package argument")
		_ = cmd.Usage()
		return cmdutil.ErrSilent
	}

	if action == "add" {
		b.ensureTrusted(cmd, rest)
	}

	ctx := cmd.Context()
	log.Info("brew bundle " + action + " --global " + strings.Join(rest, " "))
	bundleArgs := append([]string{"bundle", action}, rest...)
	bundleArgs = append(bundleArgs, "--global")
	if err := b.brew(ctx, cmd, bundleArgs...); err != nil {
		log.Error("brew bundle " + action + " failed: " + err.Error())
		return cmdutil.ErrSilent
	}

	log.Info("brew bundle install --global --zap")
	if err := b.brew(ctx, cmd, "bundle", "install", "--global", "--zap", "--upgrade"); err != nil {
		log.Error("brew bundle install failed: " + err.Error())
		return cmdutil.ErrSilent
	}

	log.Info("Brewfile in sync")
	return nil
}

// ensureTrusted trust-checks every non-official tap reference named on an add,
// prompting to `brew trust` any that are not already trusted.
func (b *brewfileCmd) ensureTrusted(cmd *cobra.Command, addArgs []string) {
	kind := bf.KindFromFlags(addArgs)
	if kind == bf.KindNone {
		return
	}
	trust, _ := bf.ParseTrust(b.readTrust())
	for _, a := range addArgs {
		if !bf.IsTapReference(a) || trust.IsTrusted(kind, a) {
			continue
		}
		b.promptTrust(cmd, kind, a)
	}
}

// trustFile resolves Homebrew's trust.json: $XDG_CONFIG_HOME/homebrew/trust.json
// when set, else ~/.homebrew/trust.json.
func (b *brewfileCmd) trustFile() string {
	if x := b.deps.Env.Get("XDG_CONFIG_HOME"); x != "" {
		return filepath.Join(x, "homebrew", "trust.json")
	}
	return filepath.Join(b.deps.Env.Home(), ".homebrew", "trust.json")
}

func (b *brewfileCmd) readTrust() []byte {
	data, _ := os.ReadFile(b.trustFile())
	return data
}

func (b *brewfileCmd) promptTrust(cmd *cobra.Command, kind bf.Kind, name string) {
	log := b.deps.Log
	log.Warn(fmt.Sprintf("%s '%s' is not trusted; brew will run its third-party code during install", kind, name))

	ok, err := b.deps.Prompt.Confirm(fmt.Sprintf("Trust %s %q now?", kind, name), false)
	switch {
	case errors.Is(err, ui.ErrNoTTY):
		log.Warn(fmt.Sprintf("%s '%s' is untrusted and no tty is available; brew bundle install may refuse it", kind, name))
		return
	case err != nil:
		log.Error(fmt.Sprintf("trust prompt failed: %v", err))
		return
	case !ok:
		log.Warn(fmt.Sprintf("leaving %s '%s' untrusted", kind, name))
		return
	}

	log.Info(fmt.Sprintf("brew trust --%s %s", kind, name))
	if _, err := b.deps.Runner.Run(cmd.Context(), "brew", "trust", "--"+kind.String(), name); err != nil {
		log.Error(fmt.Sprintf("brew trust failed: %v", err))
	}
}

// brew runs `brew <args>` with stdio wired straight through (install is long
// and informational; brewfile's stdout is not a captured value).
func (b *brewfileCmd) brew(ctx context.Context, cmd *cobra.Command, args ...string) error {
	return b.deps.Runner.RunIO(ctx, cmdutil.Streams(cmd), "brew", args...)
}
