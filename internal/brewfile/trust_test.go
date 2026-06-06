package brewfile

import "testing"

func TestParseTrust(t *testing.T) {
	tr, err := ParseTrust([]byte(`{"trustedtaps":["u/t"],"trustedcasks":["ghostty"],"trustedformulae":["jq"]}`))
	if err != nil {
		t.Fatal(err)
	}
	if len(tr.TrustedTaps) != 1 || tr.TrustedTaps[0] != "u/t" {
		t.Errorf("taps = %v", tr.TrustedTaps)
	}
	if !tr.IsTrusted(KindCask, "ghostty") || !tr.IsTrusted(KindTap, "u/t") || !tr.IsTrusted(KindFormula, "jq") {
		t.Error("expected trusted entries to report trusted")
	}
	if tr.IsTrusted(KindCask, "jq") {
		t.Error("jq is a formula, not a trusted cask")
	}
}

func TestParseTrustEmpty(t *testing.T) {
	for _, in := range []string{"", "   ", "\n\t"} {
		if tr, err := ParseTrust([]byte(in)); err != nil || len(tr.TrustedTaps) != 0 {
			t.Errorf("ParseTrust(%q) = %+v, %v", in, tr, err)
		}
	}
}

func TestParseTrustInvalid(t *testing.T) {
	if _, err := ParseTrust([]byte(`{bad`)); err == nil {
		t.Error("expected error for invalid json")
	}
}

func TestKindFromFlags(t *testing.T) {
	cases := []struct {
		args []string
		want Kind
	}{
		{[]string{"jq"}, KindFormula},
		{[]string{"--cask", "ghostty"}, KindCask},
		{[]string{"--casks", "ghostty"}, KindCask},
		{[]string{"--tap", "u/t"}, KindTap},
		{[]string{"--go", "x"}, KindNone},
		{[]string{"--vscode", "x"}, KindNone},
		{[]string{"--cask", "--formula"}, KindFormula}, // last wins
	}
	for _, c := range cases {
		if got := KindFromFlags(c.args); got != c.want {
			t.Errorf("KindFromFlags(%v) = %v, want %v", c.args, got, c.want)
		}
	}
}

func TestKindString(t *testing.T) {
	cases := map[Kind]string{KindFormula: "formula", KindCask: "cask", KindTap: "tap", KindNone: "none"}
	for k, want := range cases {
		if got := k.String(); got != want {
			t.Errorf("Kind(%d).String() = %q, want %q", k, got, want)
		}
	}
}

func TestIsTapReference(t *testing.T) {
	cases := map[string]bool{
		"user/tap/foo": true,
		"u/t":          true,
		"jq":           false,
		"--cask":       false,
		"-v":           false,
	}
	for in, want := range cases {
		if got := IsTapReference(in); got != want {
			t.Errorf("IsTapReference(%q) = %v, want %v", in, got, want)
		}
	}
}
