// Package agents detects whether the v1 install targets, Claude Code and
// Grok Build, are present on this machine (docs/design.md §7,
// docs/adr/0003-agent-detection.md).
package agents

import (
	"os"
	"os/exec"
	"path/filepath"
)

// target describes one v1 install target's detection signals: a binary
// name looked up on PATH, and a config directory checked relative to the
// user's home. Either signal is sufficient (ADR 0003).
type target struct {
	binary    string
	configDir string
}

var (
	claudeTarget = target{binary: "claude", configDir: ".claude"}
	grokTarget   = target{binary: "grok", configDir: ".grok"}
)

// Detect reports whether Claude Code and Grok Build are present on this
// machine, using the real PATH and home directory. `tb install` and
// `tb doctor` both call Detect so they can never disagree.
func Detect() (claude, grok bool) {
	home, _ := os.UserHomeDir() // an empty home just fails the config-dir check below
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
	info, err := os.Stat(filepath.Join(home, a.configDir))
	return err == nil && info.IsDir()
}

func realLookPath(binary string) bool {
	_, err := exec.LookPath(binary)
	return err == nil
}
