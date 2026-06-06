package logx

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
)

func TestTextOutputIsLeveledAndPlain(t *testing.T) {
	var buf bytes.Buffer
	lg := For(&buf, true)
	lg.Info("hello")
	lg.Warn("careful")
	out := buf.String()
	if !strings.Contains(out, "hello") || !strings.Contains(out, "careful") {
		t.Fatalf("messages missing from output: %q", out)
	}
	if strings.Contains(out, "==>") {
		t.Fatalf("the legacy ==> prefix should be gone: %q", out)
	}
}

func TestJSONOutputWhenPiped(t *testing.T) {
	var buf bytes.Buffer
	For(&buf, false).Info("ready", slog.String("path", "/x"))
	var m map[string]any
	if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
		t.Fatalf("not json: %v (%q)", err, buf.String())
	}
	if m["msg"] != "ready" || m["level"] != "info" || m["path"] != "/x" {
		t.Fatalf("unexpected json fields: %v", m)
	}
}

func TestWorktreeAttrGroupsInJSON(t *testing.T) {
	var buf bytes.Buffer
	For(&buf, false).Info("wt", WorktreeAttr(Worktree{Name: "n", Path: "/p", Branch: "agent/n", Repo: "/r"}))
	want := `"worktree":{"name":"n","path":"/p","branch":"agent/n","repo":"/r"}`
	if !strings.Contains(buf.String(), want) {
		t.Fatalf("grouped worktree attr missing\n got: %s\nwant substring: %s", buf.String(), want)
	}
}
