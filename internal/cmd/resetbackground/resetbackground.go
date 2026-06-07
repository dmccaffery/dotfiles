// Package resetbackground implements the reset-background-items command.
package resetbackground

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/dmccaffery/dotfiles/internal/cmd/cmdutil"
	"github.com/dmccaffery/dotfiles/internal/ui"
)

// NewCmd builds the reset-background-items command.
func NewCmd(deps *cmdutil.Deps) *cobra.Command {
	return &cobra.Command{
		Use:   "reset-background-items",
		Short: "Reset macOS background-item state (sfltool resetbtm), then reboot",
		Long: "Runs `sudo sfltool resetbtm` to reset all macOS background-task-management state\n" +
			"(useful when login items get stuck), then offers to reboot — required for it to take effect.",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			log := deps.Log

			ok, err := deps.Prompt.Confirm("This resets the state of ALL background items. Continue?", false)
			switch {
			case errors.Is(err, ui.ErrNoTTY), !ok:
				log.Warn("aborted")
				return nil
			case err != nil:
				log.Error(err.Error())
				return cmdutil.ErrSilent
			}

			if err := deps.Runner.RunIO(ctx, cmdutil.Streams(cmd), "sudo", "sfltool", "resetbtm"); err != nil {
				return cmdutil.ErrSilent
			}

			reboot, err := deps.Prompt.Confirm("A reboot is required to complete this. Reboot now?", false)
			if err != nil || !reboot {
				log.Warn("skipping reboot; reboot manually to finish")
				return nil
			}
			if err := deps.Runner.RunIO(ctx, cmdutil.Streams(cmd), "sudo", "reboot"); err != nil {
				return cmdutil.ErrSilent
			}
			return nil
		},
	}
}
