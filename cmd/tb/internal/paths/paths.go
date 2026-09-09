// Package paths resolves $TOOLBOX_HOME, the single clone that holds
// skills, profiles, personas, and templates (docs/design.md §1).
package paths

import (
	"fmt"
	"os"
	"path/filepath"
)

// Resolve returns the toolbox home directory: $TOOLBOX_HOME if the
// environment variable is set and non-empty, otherwise ~/toolbox. It
// fails if the resolved path does not exist or is not a directory — no
// subcommand runs without a valid toolbox home (docs/design.md §1).
//
// Resolve reads the environment exactly once; callers should call it a
// single time per invocation and pass the result down, per design.md §1.
func Resolve() (string, error) {
	home, err := resolveWith(os.Getenv("TOOLBOX_HOME"), os.UserHomeDir)
	if err != nil {
		return "", err
	}

	info, statErr := os.Stat(home)
	if statErr != nil || !info.IsDir() {
		return "", fmt.Errorf("tb: toolbox home %q not found (from $TOOLBOX_HOME or default ~/toolbox)", home)
	}
	return home, nil
}

// resolveWith computes the candidate toolbox home path without touching
// the filesystem, so the existence check in Resolve stays the only I/O.
// userHomeDir is injected for testability.
func resolveWith(envValue string, userHomeDir func() (string, error)) (string, error) {
	if envValue != "" {
		return filepath.Clean(envValue), nil
	}

	home, err := userHomeDir()
	if err != nil {
		return "", fmt.Errorf("tb: cannot determine user home directory: %w", err)
	}
	return filepath.Join(home, "toolbox"), nil
}
