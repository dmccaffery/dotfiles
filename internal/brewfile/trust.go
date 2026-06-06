// Package brewfile holds the pure logic behind the `brewfile` command: parsing
// Homebrew's trust.json and deciding which command-line arguments are
// non-official tap references that need a trust check before `brew bundle
// install` runs. No I/O — the cli command does the prompting and brew calls.
package brewfile

import (
	"bytes"
	"encoding/json"
	"slices"
	"strings"
)

// Trust mirrors the fields of Homebrew's trust.json that gate third-party taps.
type Trust struct {
	TrustedTaps     []string `json:"trustedtaps"`
	TrustedCasks    []string `json:"trustedcasks"`
	TrustedFormulae []string `json:"trustedformulae"`
}

// ParseTrust decodes trust.json. Empty input is a valid empty trust set.
func ParseTrust(data []byte) (Trust, error) {
	var t Trust
	if len(bytes.TrimSpace(data)) == 0 {
		return t, nil
	}
	err := json.Unmarshal(data, &t)
	return t, err
}

// Kind is what a `brew bundle add` selects, chosen by its flags.
type Kind int

const (
	KindFormula Kind = iota // default
	KindCask
	KindTap
	KindNone // --vscode/--go/--cargo/… : never a tap, no trust check
)

// String returns the trust noun used in messages ("formula"/"cask"/"tap").
func (k Kind) String() string {
	switch k {
	case KindCask:
		return "cask"
	case KindTap:
		return "tap"
	case KindNone:
		return "none"
	default:
		return "formula"
	}
}

// KindFromFlags determines the kind from the add arguments. As in the shell, the
// last type-selecting flag wins.
func KindFromFlags(args []string) Kind {
	k := KindFormula
	for _, a := range args {
		switch a {
		case "--cask", "--casks":
			k = KindCask
		case "--tap", "--taps":
			k = KindTap
		case "--formula", "--formulae", "--brews":
			k = KindFormula
		case "--vscode", "--go", "--cargo", "--uv", "--flatpak", "--krew", "--npm":
			k = KindNone
		}
	}
	return k
}

// IsTrusted reports whether name is already trusted for kind.
func (t Trust) IsTrusted(k Kind, name string) bool {
	var list []string
	switch k {
	case KindCask:
		list = t.TrustedCasks
	case KindTap:
		list = t.TrustedTaps
	default:
		list = t.TrustedFormulae
	}
	return slices.Contains(list, name)
}

// IsTapReference reports whether arg is a non-official tap reference (a
// "user/tap/…" path) that may need trusting — i.e. not a flag, and slashed.
// Bare names resolve to official taps and never need trust.
func IsTapReference(arg string) bool {
	if strings.HasPrefix(arg, "-") {
		return false
	}
	return strings.Contains(arg, "/")
}
