package sshaskpass

import (
	"testing"

	"github.com/dmccaffery/dotfiles/internal/cmd/cmdtest"
)

func TestBuildAssuanWithFingerprint(t *testing.T) {
	got := buildAssuan("Allow use of key SHA256:abc123:more")
	want := "SETDESC Allow use of key SHA256:abc123:more\n" +
		"OPTION allow-external-password-cache\n" +
		"SETKEYINFO s/abc123\n" +
		"GETPIN\n"
	if got != want {
		t.Fatalf("buildAssuan = %q, want %q", got, want)
	}
}

func TestBuildAssuanPlain(t *testing.T) {
	got := buildAssuan("Enter your PIN")
	if want := "SETDESC Enter your PIN\nGETPIN\n"; got != want {
		t.Fatalf("buildAssuan = %q, want %q", got, want)
	}
}

func TestExtractPin(t *testing.T) {
	for _, c := range []struct {
		name string
		in   string
		want string
	}{
		{"data line", "OK Pleased to meet you\nD 1234\nOK\n", "1234"},
		{"no data line", "OK\nERR 83886179 Operation cancelled\n", ""},
		{"empty", "", ""},
		{"bare marker", "D\n", ""},
	} {
		if got := extractPin(c.in); got != c.want {
			t.Errorf("%s: extractPin = %q, want %q", c.name, got, c.want)
		}
	}
}

func TestConfirmUserPresenceEchoesNewline(t *testing.T) {
	deps := cmdtest.NewDeps(t)
	out, _, err := cmdtest.Run(t, NewCmd(deps), "", "Confirm user presence for key ED25519-SK ...")
	if err != nil {
		t.Fatal(err)
	}
	if out != "\n" {
		t.Fatalf("stdout = %q, want a single newline", out)
	}
	// The presence path must never invoke pinentry.
	if len(cmdtest.Fake(deps).Calls) != 0 {
		t.Fatalf("presence path ran external commands: %v", cmdtest.Fake(deps).Lines())
	}
}

func TestPinPromptInvokesPinentry(t *testing.T) {
	deps := cmdtest.NewDeps(t)
	if _, _, err := cmdtest.Run(t, NewCmd(deps), "", "Enter your PIN"); err != nil {
		t.Fatal(err)
	}
	if !cmdtest.ContainsLine(cmdtest.Fake(deps).Lines(), pinentryPath) {
		t.Fatalf("expected pinentry-mac to be invoked, got %v", cmdtest.Fake(deps).Lines())
	}
}
