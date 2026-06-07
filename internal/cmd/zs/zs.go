// Package zs implements the zs command (Zscaler tunnel control + cert injection).
package zs

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/dmccaffery/dotfiles/internal/cmd/cmdutil"
)

// Daemon plist paths — vars (not consts) so tests can point them at temp files.
var (
	zscalerService = "/Library/LaunchDaemons/com.zscaler.service.plist"
	zscalerTunnel  = "/Library/LaunchDaemons/com.zscaler.tunnel.plist"
)

// tunnelWait is how long enable pauses for the tunnel to come up (a var so tests
// can zero it out).
var tunnelWait = 10 * time.Second

// execReplace replaces the current process with argv (the shell `exec "$@"`).
// It is a var so tests can observe the call instead of replacing themselves.
var execReplace = func(argv, env []string) error {
	path, err := exec.LookPath(argv[0])
	if err != nil {
		return err
	}
	return syscall.Exec(path, argv, env)
}

type zsCmd struct{ deps *cmdutil.Deps }

// NewCmd builds the zs command.
func NewCmd(deps *cmdutil.Deps) *cobra.Command {
	z := &zsCmd{deps: deps}
	zs := &cobra.Command{
		Use:   "zs <command>",
		Short: "Control the Zscaler tunnel and inject its root CA",
	}
	zs.AddCommand(
		&cobra.Command{
			Use:   "enable",
			Short: "Load the zscaler service + tunnel launch daemons",
			Args:  cobra.NoArgs,
			RunE:  z.enable,
		},
		&cobra.Command{
			Use:   "disable",
			Short: "Unload the zscaler tunnel + service",
			Args:  cobra.NoArgs,
			RunE:  z.disable,
		},
		&cobra.Command{
			Use:                "certs -- <cmd> [args...]",
			Short:              "Run <cmd> with the zscaler root CA injected as an extra trust anchor",
			DisableFlagParsing: true, // pass the wrapped command (and its flags) through verbatim
			RunE:               z.certs,
		},
	)
	return zs
}

// running reports whether launchctl lists the given daemon label.
func (z *zsCmd) running(ctx context.Context, label string) bool {
	res, err := z.deps.Runner.Run(ctx, "sudo", "launchctl", "list")
	if err != nil {
		return false
	}
	return strings.Contains(res.Stdout, label)
}

func (z *zsCmd) enable(cmd *cobra.Command, _ []string) error {
	ctx, log, r := cmd.Context(), z.deps.Log, z.deps.Runner
	started := false

	for _, d := range []struct{ plist, label, what string }{
		{zscalerService, "com.zscaler.service", "service"},
		{zscalerTunnel, "com.zscaler.tunnel", "tunnel"},
	} {
		if cmdutil.FileExists(d.plist) && !z.running(ctx, d.label) {
			log.Warn("enabling zscaler " + d.what + "...")
			if err := r.RunIO(ctx, cmdutil.Streams(cmd), "sudo", "launchctl", "load", d.plist); err != nil {
				return cmdutil.ErrSilent
			}
			started = true
		}
	}
	if started {
		log.Warn("waiting for tunnel...")
		time.Sleep(tunnelWait)
	}
	return nil
}

func (z *zsCmd) disable(cmd *cobra.Command, _ []string) error {
	ctx, log, r := cmd.Context(), z.deps.Log, z.deps.Runner

	// Tunnel first, then service.
	for _, d := range []struct{ plist, label, what string }{
		{zscalerTunnel, "com.zscaler.tunnel", "tunnel"},
		{zscalerService, "com.zscaler.service", "service"},
	} {
		if cmdutil.FileExists(d.plist) && z.running(ctx, d.label) {
			log.Warn("disabling zscaler " + d.what + "...")
			if err := r.RunIO(ctx, cmdutil.Streams(cmd), "sudo", "launchctl", "unload", d.plist); err != nil {
				return cmdutil.ErrSilent
			}
		}
	}
	return nil
}

func (z *zsCmd) certs(cmd *cobra.Command, args []string) error {
	ctx, log := cmd.Context(), z.deps.Log

	if len(args) == 0 || args[0] != "--" {
		log.Error("expected '--' before the command to run")
		return cmdutil.ErrSilent
	}
	rest := args[1:]
	if len(rest) == 0 {
		log.Error("no command given after '--'")
		return cmdutil.ErrSilent
	}

	cert := filepath.Join(z.deps.Env.XDGDataHome(), "certificates", "zscaler.pem")
	if cmdutil.FileExists(cert) {
		log.Info("found existing certificate at: " + cert)
	} else {
		if err := os.MkdirAll(filepath.Dir(cert), 0o755); err != nil {
			log.Error(err.Error())
			return cmdutil.ErrSilent
		}
		res, err := z.deps.Runner.Run(ctx, "security", "find-certificate",
			"-c", "Zscaler Root CA", "-p", "/Library/Keychains/System.keychain")
		if err != nil || strings.TrimSpace(res.Stdout) == "" {
			log.Warn("no certificate could be found in keychain; running command without certificates...")
			return execReplace(rest, os.Environ())
		}
		if err := os.WriteFile(cert, []byte(res.Stdout), 0o644); err != nil {
			log.Error(err.Error())
			return cmdutil.ErrSilent
		}
	}

	env := append(os.Environ(), "ZSCALER_CERTIFICATE="+cert, "NODE_EXTRA_CA_CERTS="+cert)
	return execReplace(rest, env)
}
