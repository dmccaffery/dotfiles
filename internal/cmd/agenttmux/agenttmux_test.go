package agenttmux_test

import (
	"testing"

	"github.com/dmccaffery/dotfiles/internal/cmd/agenttmux"
	"github.com/dmccaffery/dotfiles/internal/cmd/cmdtest"
	"github.com/dmccaffery/dotfiles/internal/envx"
)

func TestAgentTmuxStatusInsideTmux(t *testing.T) {
	cases := []struct{ state, want string }{
		{"waiting", "tmux set-window-option -t %3 @agent_status waiting"},
		{"attention", "tmux set-window-option -t %3 @agent_status attention"},
		{"clear", "tmux set-window-option -t %3 -u @agent_status"},
		{"bogus", "tmux set-window-option -t %3 -u @agent_status"},
	}
	for _, c := range cases {
		deps := cmdtest.NewDeps(t)
		deps.Env = envx.New(t.TempDir(), map[string]string{"TMUX": "/tmp/tmux/default,1,0", "TMUX_PANE": "%3"})
		if _, _, err := cmdtest.Run(t, agenttmux.NewCmd(deps), "", c.state); err != nil {
			t.Fatalf("state %s: %v", c.state, err)
		}
		if !cmdtest.ContainsLine(cmdtest.Fake(deps).Lines(), c.want) {
			t.Fatalf("state %s: want %q in %v", c.state, c.want, cmdtest.Fake(deps).Lines())
		}
	}
}
