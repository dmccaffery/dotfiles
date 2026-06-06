package cli_test

import (
	"testing"

	"github.com/dmccaffery/dotfiles/internal/cli"
	"github.com/dmccaffery/dotfiles/internal/envx"
	"github.com/dmccaffery/dotfiles/internal/execx"
)

func TestAgentTmuxStatusInsideTmux(t *testing.T) {
	cases := []struct{ state, want string }{
		{"waiting", "tmux set-window-option -t %3 @agent_status waiting"},
		{"attention", "tmux set-window-option -t %3 @agent_status attention"},
		{"clear", "tmux set-window-option -t %3 -u @agent_status"},
		{"bogus", "tmux set-window-option -t %3 -u @agent_status"},
	}
	for _, c := range cases {
		fake := execx.NewFake()
		env := envx.New(t.TempDir(), map[string]string{"TMUX": "/tmp/tmux/default,1,0", "TMUX_PANE": "%3"})
		deps := &cli.Deps{Runner: fake, Env: env}
		if _, _, err := runRoot(t, deps, "", "agent-tmux-status", c.state); err != nil {
			t.Fatalf("state %s: %v", c.state, err)
		}
		if !containsLine(fake.Lines(), c.want) {
			t.Fatalf("state %s: want %q in %v", c.state, c.want, fake.Lines())
		}
	}
}
