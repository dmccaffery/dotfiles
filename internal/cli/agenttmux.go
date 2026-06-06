package cli

import (
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// Glyphs for the OSC (non-tmux) terminal-title fallback. Inside tmux the glyph
// and colour come from theme.conf via the @agent_status state token.
const (
	waitingGlyph   = "● "
	attentionGlyph = "󰂚 "
)

func newAgentTmuxCmd(deps *Deps) *cobra.Command {
	return &cobra.Command{
		Use:   "agent-tmux-status [waiting|attention|clear]",
		Short: "Set the coding-agent status indicator (tmux window option or terminal title)",
		Long: "Shared status indicator for coding agents. Inside tmux it sets a per-window\n" +
			"@agent_status token (waiting|attention) that theme.conf maps to a colour and\n" +
			"glyph; outside tmux it sets the terminal title via OSC 0. Always exits 0.",
		Args:          cobra.MaximumNArgs(1),
		SilenceErrors: true,
		// No-op-safe: this runs from agent lifecycle hooks and must never fail.
		RunE: func(cmd *cobra.Command, args []string) error {
			state := "clear"
			if len(args) > 0 {
				state = args[0]
			}
			agentStatus(cmd, deps, state)
			return nil
		},
	}
}

func agentStatus(cmd *cobra.Command, deps *Deps, state string) {
	if pane := deps.Env.Get("TMUX"); pane != "" {
		target := deps.Env.Get("TMUX_PANE")
		switch state {
		case "waiting", "attention":
			_, _ = deps.Runner.Run(cmd.Context(), "tmux", "set-window-option", "-t", target, "@agent_status", state)
		default:
			_, _ = deps.Runner.Run(cmd.Context(), "tmux", "set-window-option", "-t", target, "-u", "@agent_status")
		}
		return
	}

	title := filepath.Base(deps.Env.GetOr("PWD", cwd()))
	switch state {
	case "waiting":
		writeTTY("\033]0;" + waitingGlyph + title + "\007")
	case "attention":
		writeTTY("\033]0;" + attentionGlyph + title + "\007")
	default:
		writeTTY("\033]0;" + title + "\007")
	}
}

func cwd() string {
	d, err := os.Getwd()
	if err != nil {
		return ""
	}
	return d
}

// writeTTY writes s to the controlling terminal, swallowing every error (the
// agent may have captured the caller's stdout, so the title must go to /dev/tty).
func writeTTY(s string) {
	f, err := os.OpenFile("/dev/tty", os.O_WRONLY, 0)
	if err != nil {
		return
	}
	defer f.Close()
	_, _ = io.WriteString(f, s)
}
