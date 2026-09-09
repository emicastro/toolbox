// Package agents detects whether the v1 install targets, Claude Code and
// Grok Build, are present on this machine (docs/design.md §7,
// docs/adr/0004-agent-detection-marker-file.md).
package agents

import (
	"os"
	"os/exec"
	"path/filepath"
)

// target describes one v1 install target's detection signals: a binary
// name looked up on PATH, and a marker file checked relative to the
// user's home. Either signal is sufficient (ADR 0004).
//
// The marker must be a file only the real client creates — not the bare
// config directory (ADR 0003, superseded): `tb install` itself creates
// `~/.claude/skills`, and MkdirAll creates `~/.claude` as that path's
// parent, so checking the directory's mere existence made every agent
// permanently read "present" after the first `tb install` run.
type target struct {
	binary     string
	markerFile string // relative to $HOME
}

var (
	claudeTarget = target{binary: "claude", markerFile: filepath.Join(".claude", "settings.json")}
	grokTarget   = target{binary: "grok", markerFile: filepath.Join(".grok", "config.toml")}
)

// Detect reports whether Claude Code and Grok Build are present on this
// machine, using the real PATH and home directory. `tb install` and
// `tb doctor` both call Detect so they can never disagree.
func Detect() (claude, grok bool) {
	home, _ := os.UserHomeDir() // an empty home just fails the marker-file check below
	return DetectIn(realLookPath, home)
}

// DetectIn is Detect with its PATH lookup and home directory injected,
// so tests can cover all four presence combinations without touching the
// real environment (docs/design.md §12).
func DetectIn(lookPath func(binary string) bool, home string) (claude, grok bool) {
	return present(claudeTarget, lookPath, home), present(grokTarget, lookPath, home)
}

func present(a target, lookPath func(string) bool, home string) bool {
	if lookPath(a.binary) {
		return true
	}
	if home == "" {
		return false
	}
	_, err := os.Stat(filepath.Join(home, a.markerFile))
	return err == nil
}

func realLookPath(binary string) bool {
	_, err := exec.LookPath(binary)
	return err == nil
}
