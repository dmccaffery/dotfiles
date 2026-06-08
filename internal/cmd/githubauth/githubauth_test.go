package githubauth

import (
	"strings"
	"testing"

	"github.com/dmccaffery/dotfiles/internal/cmd/cmdtest"
	"github.com/dmccaffery/dotfiles/internal/execx"
	"github.com/dmccaffery/dotfiles/internal/ui"
)

const (
	loginsCmd = `gh auth status --json hosts --jq .hosts."github.com".[].login`
	activeCmd = `gh auth status --active --json hosts --jq .hosts["github.com"][0].login`
	scopesCmd = `gh auth status --active --json hosts --jq .hosts."github.com".[] | .scopes`
)

func TestFirstMissingScope(t *testing.T) {
	if s, missing := firstMissingScope(strings.Join(requiredScopes, ", "), requiredScopes); missing {
		t.Fatalf("none should be missing, got %q", s)
	}
	if s, missing := firstMissingScope("gist repo user", requiredScopes); !missing || s != "notifications" {
		t.Fatalf("got %q missing=%v, want notifications", s, missing)
	}
}

func allScopes() execx.Result { return execx.Result{Stdout: strings.Join(requiredScopes, ", ")} }

func TestGitHubAuthAlreadyGood(t *testing.T) {
	deps := cmdtest.NewDeps(t)
	fake := cmdtest.Fake(deps)
	fake.Responses[activeCmd] = execx.Result{Stdout: "alice\n"}
	fake.Responses[scopesCmd] = allScopes()

	if _, _, err := cmdtest.Run(t, NewCmd(deps), "", "alice"); err != nil {
		t.Fatal(err)
	}
	if cmdtest.ContainsLine(fake.Lines(), "auth login") || cmdtest.ContainsLine(fake.Lines(), "auth refresh") {
		t.Fatalf("should not log in or refresh when already authed with scopes: %v", fake.Lines())
	}
}

func TestGitHubAuthMissingScopeRefreshes(t *testing.T) {
	deps := cmdtest.NewDeps(t)
	fake := cmdtest.Fake(deps)
	fake.Responses[activeCmd] = execx.Result{Stdout: "alice\n"}
	fake.Responses[scopesCmd] = execx.Result{Stdout: "gist, repo, user"} // missing several

	if _, _, err := cmdtest.Run(t, NewCmd(deps), "", "alice"); err != nil {
		t.Fatal(err)
	}
	if !cmdtest.ContainsLine(fake.Lines(), "gh auth refresh --hostname github.com --scopes") {
		t.Fatalf("expected a scope refresh, got %v", fake.Lines())
	}
}

func TestGitHubAuthPicksExistingAccount(t *testing.T) {
	deps := cmdtest.NewDeps(t)
	deps.Prompt = &ui.Fake{Selections: []string{"bob"}}
	fake := cmdtest.Fake(deps)
	fake.Responses[loginsCmd] = execx.Result{Stdout: "alice\nbob\n"}
	fake.Responses[activeCmd] = execx.Result{Stdout: "bob\n"}
	fake.Responses[scopesCmd] = allScopes()

	if _, _, err := cmdtest.Run(t, NewCmd(deps), ""); err != nil {
		t.Fatal(err)
	}
	if !cmdtest.ContainsLine(fake.Lines(), "gh auth switch --user bob") {
		t.Fatalf("expected switch to bob, got %v", fake.Lines())
	}
}

func TestGitHubAuthNewAccountLogsIn(t *testing.T) {
	deps := cmdtest.NewDeps(t)
	deps.Prompt = &ui.Fake{Selections: []string{"new account"}}
	fake := cmdtest.Fake(deps)
	fake.Responses[loginsCmd] = execx.Result{Stdout: "alice\nbob\n"}

	if _, _, err := cmdtest.Run(t, NewCmd(deps), ""); err != nil {
		t.Fatal(err)
	}
	if !cmdtest.ContainsLine(fake.Lines(), "gh auth login") {
		t.Fatalf("expected gh auth login, got %v", fake.Lines())
	}
}

func TestGitHubAuthNoTTYErrors(t *testing.T) {
	deps := cmdtest.NewDeps(t)
	deps.Prompt = &ui.Fake{NoTTY: true}
	cmdtest.Fake(deps).Responses[loginsCmd] = execx.Result{Stdout: "alice\nbob\n"}

	if _, _, err := cmdtest.Run(t, NewCmd(deps), ""); err == nil {
		t.Fatal("expected an error with no account and no tty")
	}
}
