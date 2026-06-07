// Package printcolors implements the print-colors command.
package printcolors

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/dmccaffery/dotfiles/internal/cmd/cmdutil"
	"github.com/dmccaffery/dotfiles/internal/colorbar"
)

// NewCmd builds the print-colors command.
func NewCmd(_ *cmdutil.Deps) *cobra.Command {
	return &cobra.Command{
		Use:   "print-colors",
		Short: "Print a 24-bit truecolor gradient bar",
		Long:  "Prints a horizontal 24-bit RGB gradient — a smoke test that the terminal supports true color.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			fmt.Fprint(cmd.OutOrStdout(), colorbar.Gradient())
			return nil
		},
	}
}
