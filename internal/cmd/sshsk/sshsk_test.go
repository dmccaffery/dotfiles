package sshsk

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/dmccaffery/dotfiles/internal/cmd/cmdtest"
	"github.com/dmccaffery/dotfiles/internal/cmd/cmdutil"
	"github.com/dmccaffery/dotfiles/internal/execx"
	"github.com/dmccaffery/dotfiles/internal/ui"
)

const (
	ghAccountCmd = "git config --get github.account"
	fjAccountCmd = "git config --get forgejo.account"
	emailCmd     = "git config --get user.email"
	serialsCmd   = "ykman list --serials"
	agentKeysCmd = "ssh-add -L"
	ghStatusCmd  = `gh auth status --hostname github.com --json hosts --jq .hosts["github.com"] // [] | .[].login`

	agentLine = "sk-ssh-ed25519@openssh.com AAAAblob alice@host"
)

// writeStub creates a saved stub pair for user under serial, with the given
// public-key blob, rooted at the deps' XDG config home.
func writeStub(t *testing.T, deps *cmdutil.Deps, serial, user, blob string) {
	t.Helper()
	dir := filepath.Join(deps.Env.XDGConfigHome(), "private", "ssh", serial)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	base := filepath.Join(dir, "id_ed25519_sk_"+user)
	if err := os.WriteFile(base, []byte("PRIVATE"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(base+".pub", []byte("sk-ssh-ed25519@openssh.com "+blob+" "+user+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
}

// --- pure helpers ---

func TestSafeUser(t *testing.T) {
	for _, c := range []struct {
		in   string
		want bool
	}{
		{"alice", true},
		{"user.name@host+tag-1_2", true},
		{"", false},
		{".", false},
		{"..", false},
		{"a/b", false},
		{"a b", false},
		{"../etc", false},
	} {
		if got := safeUser(c.in); got != c.want {
			t.Errorf("safeUser(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestParseForgejoAccounts(t *testing.T) {
	got := parseForgejoAccounts("alice@git.example.com\nbob@code.test  extra\n\nnohost\n@nouser\nnotrail@\n")
	want := []account{
		{provider: "forgejo", host: "git.example.com", user: "alice"},
		{provider: "forgejo", host: "code.test", user: "bob"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseForgejoAccounts = %+v, want %+v", got, want)
	}
}

func TestMatchAgentKey(t *testing.T) {
	agent := "sk-ssh-ed25519@openssh.com AAA alice\nssh-ed25519 BBB bob\n"
	if got := matchAgentKey(agent, []string{"BBB"}); got != "ssh-ed25519 BBB bob" {
		t.Errorf("matchAgentKey hit = %q", got)
	}
	if got := matchAgentKey(agent, []string{"ZZZ"}); got != "" {
		t.Errorf("matchAgentKey miss = %q, want empty", got)
	}
}

func TestUniqueUsers(t *testing.T) {
	got := uniqueUsers([]account{{user: "a"}, {user: "b"}, {user: "a"}, {user: ""}})
	if want := []string{"a", "b"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("uniqueUsers = %v, want %v", got, want)
	}
}

func TestSignerLine(t *testing.T) {
	if got := signerLine("e@x", "KEY"); got != `e@x namespaces="git" KEY` {
		t.Fatalf("signerLine = %q", got)
	}
}

// --- the account picker (huh replaces fzf) ---

func TestSelectAuthenticatedUserAutoSelectsSole(t *testing.T) {
	deps := cmdtest.NewDeps(t)
	fake := cmdtest.Fake(deps)
	fake.Missing["fj"] = true
	fake.Responses[ghStatusCmd] = execx.Result{Stdout: "alice\n"}

	got, err := (&skCmd{deps: deps}).selectAuthenticatedUser(context.Background())
	if err != nil || got != "alice" {
		t.Fatalf("got %q, err %v", got, err)
	}
	if asked := deps.Prompt.(*ui.Fake).Asked; len(asked) != 0 {
		t.Fatalf("a sole account should not prompt, asked %v", asked)
	}
}

func TestSelectAuthenticatedUserPicksAmongMany(t *testing.T) {
	deps := cmdtest.NewDeps(t)
	deps.Prompt = &ui.Fake{Selections: []string{"bob"}}
	fake := cmdtest.Fake(deps)
	fake.Missing["fj"] = true
	fake.Responses[ghStatusCmd] = execx.Result{Stdout: "alice\nbob\n"}

	got, err := (&skCmd{deps: deps}).selectAuthenticatedUser(context.Background())
	if err != nil || got != "bob" {
		t.Fatalf("got %q, err %v", got, err)
	}
}

func TestSelectAuthenticatedUserNoneErrors(t *testing.T) {
	deps := cmdtest.NewDeps(t)
	fake := cmdtest.Fake(deps)
	fake.Missing["gh"], fake.Missing["fj"] = true, true

	if _, err := (&skCmd{deps: deps}).selectAuthenticatedUser(context.Background()); err == nil {
		t.Fatal("expected an error with no authenticated accounts")
	}
}

// --- get --git (the defaultKeyCommand resolver) ---

func TestGetGitResolvesConfiguredAccount(t *testing.T) {
	deps := cmdtest.NewDeps(t)
	fake := cmdtest.Fake(deps)
	fake.Responses[ghAccountCmd] = execx.Result{Stdout: "alice\n"}
	fake.Responses[serialsCmd] = execx.Result{Stdout: "12345\n"}
	fake.Responses[agentKeysCmd] = execx.Result{Stdout: agentLine + "\n"}
	writeStub(t, deps, "12345", "alice", "AAAAblob")

	out, _, err := cmdtest.Run(t, NewCmd(deps), "", "get", "--git")
	if err != nil {
		t.Fatal(err)
	}
	if want := "key::" + agentLine + "\n"; out != want {
		t.Fatalf("stdout = %q, want %q", out, want)
	}
}

func TestGetGitNoAccountErrors(t *testing.T) {
	deps := cmdtest.NewDeps(t)
	out, _, err := cmdtest.Run(t, NewCmd(deps), "", "get", "--git")
	if err == nil {
		t.Fatal("expected an error with no configured account")
	}
	if out != "" {
		t.Fatalf("stdout must stay empty for git, got %q", out)
	}
}

func TestGetGitConflictingAccountsError(t *testing.T) {
	deps := cmdtest.NewDeps(t)
	fake := cmdtest.Fake(deps)
	fake.Responses[ghAccountCmd] = execx.Result{Stdout: "alice\n"}
	fake.Responses[fjAccountCmd] = execx.Result{Stdout: "bob\n"}

	out, _, err := cmdtest.Run(t, NewCmd(deps), "", "get", "--git")
	if err == nil {
		t.Fatal("expected an error when the two accounts conflict")
	}
	if out != "" {
		t.Fatalf("stdout = %q, want empty", out)
	}
}

func TestGetGitNoKeyNoTTY(t *testing.T) {
	deps := cmdtest.NewDeps(t)
	deps.Prompt = &ui.Fake{NoTTY: true}
	cmdtest.Fake(deps).Responses[ghAccountCmd] = execx.Result{Stdout: "alice\n"}

	out, _, err := cmdtest.Run(t, NewCmd(deps), "", "get", "--git")
	if err == nil {
		t.Fatal("expected an error with no inserted key and no tty")
	}
	if out != "" {
		t.Fatalf("stdout = %q, want empty", out)
	}
}

// --- get (setup) ---

func TestSetupAppendsAllowedSigner(t *testing.T) {
	deps := cmdtest.NewDeps(t)
	fake := cmdtest.Fake(deps)
	fake.Responses[ghAccountCmd] = execx.Result{Stdout: "alice\n"}
	fake.Responses[emailCmd] = execx.Result{Stdout: "alice@example.com\n"}
	fake.Responses[serialsCmd] = execx.Result{Stdout: "12345\n"}
	fake.Responses[agentKeysCmd] = execx.Result{Stdout: agentLine + "\n"}
	writeStub(t, deps, "12345", "alice", "AAAAblob")

	if _, _, err := cmdtest.Run(t, NewCmd(deps), "", "get"); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(deps.Env.XDGConfigHome(), "private", "git", "allowed_signers"))
	if err != nil {
		t.Fatal(err)
	}
	want := `alice@example.com namespaces="git" ` + agentLine + "\n"
	if string(data) != want {
		t.Fatalf("allowed_signers = %q, want %q", data, want)
	}
}

func TestSetupSkipsDuplicateSigner(t *testing.T) {
	deps := cmdtest.NewDeps(t)
	fake := cmdtest.Fake(deps)
	fake.Responses[ghAccountCmd] = execx.Result{Stdout: "alice\n"}
	fake.Responses[emailCmd] = execx.Result{Stdout: "alice@example.com\n"}
	fake.Responses[serialsCmd] = execx.Result{Stdout: "12345\n"}
	fake.Responses[agentKeysCmd] = execx.Result{Stdout: agentLine + "\n"}
	writeStub(t, deps, "12345", "alice", "AAAAblob")

	signers := filepath.Join(deps.Env.XDGConfigHome(), "private", "git", "allowed_signers")
	if err := os.MkdirAll(filepath.Dir(signers), 0o700); err != nil {
		t.Fatal(err)
	}
	line := `alice@example.com namespaces="git" ` + agentLine + "\n"
	if err := os.WriteFile(signers, []byte(line), 0o600); err != nil {
		t.Fatal(err)
	}

	if _, _, err := cmdtest.Run(t, NewCmd(deps), "", "get"); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(signers)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != line {
		t.Fatalf("duplicate signer appended: %q", data)
	}
}

func TestSetupNoStubsErrors(t *testing.T) {
	deps := cmdtest.NewDeps(t)
	cmdtest.Fake(deps).Responses[serialsCmd] = execx.Result{Stdout: "12345\n"}
	// No stub files written → loadCurrentStubs finds nothing.
	if _, _, err := cmdtest.Run(t, NewCmd(deps), "", "get"); err == nil {
		t.Fatal("expected an error when no stubs are present")
	}
}

// --- gen ---

func TestGenRejectsUnsafeUser(t *testing.T) {
	deps := cmdtest.NewDeps(t)
	fake := cmdtest.Fake(deps)
	if _, _, err := cmdtest.Run(t, NewCmd(deps), "", "gen", ".."); err == nil {
		t.Fatal("expected an error for an unsafe user")
	}
	if cmdtest.ContainsLine(fake.Lines(), "ssh-keygen") {
		t.Fatal("must not keygen for an unsafe user")
	}
}

func TestGenRefusesExistingStub(t *testing.T) {
	deps := cmdtest.NewDeps(t)
	fake := cmdtest.Fake(deps)
	fake.Missing["fj"] = true
	fake.Responses[ghStatusCmd] = execx.Result{Stdout: "alice\n"} // matching account, no confirm
	fake.Responses[serialsCmd] = execx.Result{Stdout: "12345\n"}
	writeStub(t, deps, "12345", "alice", "AAAAblob")

	if _, _, err := cmdtest.Run(t, NewCmd(deps), "", "gen", "alice"); err == nil {
		t.Fatal("expected an error when a stub already exists")
	}
	if cmdtest.ContainsLine(fake.Lines(), "ssh-keygen") {
		t.Fatal("must not keygen over an existing stub")
	}
}

func TestGenDeclinedUnpublishedAborts(t *testing.T) {
	deps := cmdtest.NewDeps(t)
	deps.Prompt = &ui.Fake{Replies: []bool{false}} // decline the unpublished confirm
	fake := cmdtest.Fake(deps)
	fake.Missing["gh"], fake.Missing["fj"] = true, true // no accounts → unpublished

	if _, _, err := cmdtest.Run(t, NewCmd(deps), "", "gen", "alice"); err == nil {
		t.Fatal("expected an abort when the unpublished confirm is declined")
	}
	if cmdtest.ContainsLine(fake.Lines(), "ssh-keygen") {
		t.Fatal("a declined gen must not keygen")
	}
}
