package execx

import (
	"context"
	"fmt"
	"strings"
)

// Call records a single Run/RunIO invocation, in order.
type Call struct {
	Name string
	Args []string
}

// String renders the call as the command line that was run.
func (c Call) String() string {
	return strings.TrimSpace(c.Name + " " + strings.Join(c.Args, " "))
}

// Fake is a test Runner. Responses and Errs are keyed by the full command line
// ("name arg1 arg2 ..."); Missing names fail Look.
type Fake struct {
	Calls     []Call
	Responses map[string]Result
	Errs      map[string]error
	Missing   map[string]bool
}

// NewFake returns an initialised Fake.
func NewFake() *Fake {
	return &Fake{
		Responses: map[string]Result{},
		Errs:      map[string]error{},
		Missing:   map[string]bool{},
	}
}

func key(name string, args []string) string {
	return strings.TrimSpace(name + " " + strings.Join(args, " "))
}

// Run implements Runner.
func (f *Fake) Run(_ context.Context, name string, args ...string) (Result, error) {
	k := key(name, args)
	f.Calls = append(f.Calls, Call{Name: name, Args: args})
	return f.Responses[k], f.Errs[k]
}

// RunIO implements Runner.
func (f *Fake) RunIO(_ context.Context, _ Streams, name string, args ...string) error {
	k := key(name, args)
	f.Calls = append(f.Calls, Call{Name: name, Args: args})
	return f.Errs[k]
}

// Look implements Runner.
func (f *Fake) Look(name string) (string, error) {
	if f.Missing[name] {
		return "", fmt.Errorf("%s: not found", name)
	}
	return "/usr/bin/" + name, nil
}

// Lines returns each recorded call rendered as a command line.
func (f *Fake) Lines() []string {
	out := make([]string, len(f.Calls))
	for i, c := range f.Calls {
		out[i] = c.String()
	}
	return out
}
