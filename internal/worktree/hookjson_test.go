package worktree

import "testing"

func TestParseStartName(t *testing.T) {
	cases := map[string]string{
		`{"name":"feature"}`:       "feature",
		`{"name":"feature","x":1}`: "feature",
		`{}`:                       "",
		`{"name":""}`:              "",
		``:                         "",
		`not json`:                 "",
		`{"worktree_path":"/a"}`:   "",
	}
	for in, want := range cases {
		if got := ParseStartName([]byte(in)); got != want {
			t.Errorf("ParseStartName(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestParseEndPath(t *testing.T) {
	cases := map[string]string{
		`{"worktree_path":"/a/b"}`: "/a/b",
		`{}`:                       "",
		``:                         "",
		`garbage`:                  "",
		`{"name":"x"}`:             "",
	}
	for in, want := range cases {
		if got := ParseEndPath([]byte(in)); got != want {
			t.Errorf("ParseEndPath(%q) = %q, want %q", in, got, want)
		}
	}
}
