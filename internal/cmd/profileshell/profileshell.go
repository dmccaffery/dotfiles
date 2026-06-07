// Package profileshell implements the profile-shell command.
package profileshell

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/spf13/cobra"

	"github.com/dmccaffery/dotfiles/internal/cmd/cmdutil"
)

// NewCmd builds the profile-shell command.
func NewCmd(_ *cmdutil.Deps) *cobra.Command {
	return &cobra.Command{
		Use:   "profile-shell",
		Short: "Time Zsh startup with zprof enabled",
		Long: "Runs `zsh -i -c exit` with ZSHPROFILE=1 so .zshrc loads zprof and prints a timing\n" +
			"table, then reports real/user/sys time (like `time -p`).",
		Args: cobra.NoArgs,
		// Times the real zsh, so it runs zsh directly rather than through the
		// mockable Runner — there is nothing meaningful to fake here.
		RunE: func(cmd *cobra.Command, _ []string) error {
			c := exec.CommandContext(cmd.Context(), "zsh", "-i", "-c", "exit")
			c.Env = append(os.Environ(), "ZSHPROFILE=1")
			c.Stdin = cmd.InOrStdin()
			c.Stdout = cmd.OutOrStdout()
			c.Stderr = cmd.ErrOrStderr()

			start := time.Now()
			runErr := c.Run()
			real := time.Since(start)

			var user, sys time.Duration
			if c.ProcessState != nil {
				user = c.ProcessState.UserTime()
				sys = c.ProcessState.SystemTime()
			}
			fmt.Fprintf(cmd.ErrOrStderr(), "real %.2f\nuser %.2f\nsys %.2f\n",
				real.Seconds(), user.Seconds(), sys.Seconds())

			if runErr != nil {
				return cmdutil.ErrSilent
			}
			return nil
		},
	}
}
