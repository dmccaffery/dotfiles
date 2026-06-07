package zs

import (
	"io"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/dmccaffery/dotfiles/internal/cmd/cmdtest"
	"github.com/dmccaffery/dotfiles/internal/cmd/cmdutil"
	"github.com/dmccaffery/dotfiles/internal/envx"
	"github.com/dmccaffery/dotfiles/internal/execx"
	"github.com/dmccaffery/dotfiles/internal/logx"
	"github.com/dmccaffery/dotfiles/internal/ui"
)

// swap sets *ptr to val and returns a restore func (for stubbing package vars).
func swap[T any](ptr *T, val T) func() {
	old := *ptr
	*ptr = val
	return func() { *ptr = old }
}

func zsDeps(t *testing.T, env map[string]string) (*cmdutil.Deps, *execx.Fake) {
	t.Helper()
	fake := execx.NewFake()
	return &cmdutil.Deps{
		Runner: fake,
		Env:    envx.New(t.TempDir(), env),
		Log:    logx.For(io.Discard, false),
		Prompt: &ui.Fake{},
	}, fake
}

func TestZsEnableLoadsStoppedDaemons(t *testing.T) {
	dir := t.TempDir()
	svc := filepath.Join(dir, "service.plist")
	tun := filepath.Join(dir, "tunnel.plist")
	for _, p := range []string{svc, tun} {
		if err := os.WriteFile(p, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	defer swap(&zscalerService, svc)()
	defer swap(&zscalerTunnel, tun)()
	defer swap(&tunnelWait, 0)()

	deps, fake := zsDeps(t, nil) // launchctl list is empty → not running → load both
	if _, _, err := cmdtest.Run(t, NewCmd(deps), "", "enable"); err != nil {
		t.Fatal(err)
	}
	if !slices.Contains(fake.Lines(), "sudo launchctl load "+svc) ||
		!slices.Contains(fake.Lines(), "sudo launchctl load "+tun) {
		t.Fatalf("expected both daemons loaded, got %v", fake.Lines())
	}
}

func TestZsCertsInjectsCertAndExecs(t *testing.T) {
	xdg := filepath.Join(t.TempDir(), "data")
	certDir := filepath.Join(xdg, "certificates")
	if err := os.MkdirAll(certDir, 0o755); err != nil {
		t.Fatal(err)
	}
	cert := filepath.Join(certDir, "zscaler.pem")
	if err := os.WriteFile(cert, []byte("CERT"), 0o644); err != nil {
		t.Fatal(err)
	}

	var gotArgv, gotEnv []string
	defer swap(&execReplace, func(argv, env []string) error {
		gotArgv, gotEnv = argv, env
		return nil
	})()

	deps, _ := zsDeps(t, map[string]string{"XDG_DATA_HOME": xdg})
	if _, _, err := cmdtest.Run(t, NewCmd(deps), "", "certs", "--", "echo", "hi"); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(gotArgv, []string{"echo", "hi"}) {
		t.Fatalf("exec argv = %v, want [echo hi]", gotArgv)
	}
	if !slices.Contains(gotEnv, "NODE_EXTRA_CA_CERTS="+cert) {
		t.Fatal("NODE_EXTRA_CA_CERTS was not injected into the child env")
	}
}

func TestZsCertsRequiresDashDash(t *testing.T) {
	deps, _ := zsDeps(t, nil)
	if _, _, err := cmdtest.Run(t, NewCmd(deps), "", "certs", "echo", "hi"); err == nil {
		t.Fatal("expected an error when '--' is missing")
	}
}
