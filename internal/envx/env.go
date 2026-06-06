// Package envx abstracts process-environment and home/XDG lookups so commands
// can be unit-tested without touching the real environment.
package envx

import (
	"os"
	"path/filepath"
)

// Env provides the subset of environment access the dot commands need. Build
// one with System for real use or New for tests.
type Env struct {
	lookup func(string) (string, bool)
	home   string
}

// System returns an Env backed by the real process environment.
func System() Env {
	home, _ := os.UserHomeDir()
	return Env{lookup: os.LookupEnv, home: home}
}

// New builds an Env from a fixed home and variable map (tests).
func New(home string, vars map[string]string) Env {
	return Env{
		home: home,
		lookup: func(k string) (string, bool) {
			v, ok := vars[k]
			return v, ok
		},
	}
}

// Get returns the value of key, or "" if unset.
func (e Env) Get(key string) string {
	if e.lookup == nil {
		return ""
	}
	v, _ := e.lookup(key)
	return v
}

// GetOr returns the value of key, or def if it is unset or empty.
func (e Env) GetOr(key, def string) string {
	if v := e.Get(key); v != "" {
		return v
	}
	return def
}

// Home returns the user's home directory.
func (e Env) Home() string { return e.home }

// XDGConfigHome returns $XDG_CONFIG_HOME, defaulting to $HOME/.config.
func (e Env) XDGConfigHome() string {
	return e.GetOr("XDG_CONFIG_HOME", filepath.Join(e.home, ".config"))
}

// XDGDataHome returns $XDG_DATA_HOME, defaulting to $HOME/.local/share.
func (e Env) XDGDataHome() string {
	return e.GetOr("XDG_DATA_HOME", filepath.Join(e.home, ".local", "share"))
}

// XDGCacheHome returns $XDG_CACHE_HOME, defaulting to $HOME/.cache.
func (e Env) XDGCacheHome() string {
	return e.GetOr("XDG_CACHE_HOME", filepath.Join(e.home, ".cache"))
}
