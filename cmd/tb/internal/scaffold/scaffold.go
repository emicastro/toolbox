// Package scaffold materializes a profile's docs/ templates into a
// product repo (docs/design.md §5.2 step 6, §8.2).
package scaffold

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const ignoreComment = "# Session handoff is local working state (tb init)"

// EnsureIgnored makes sure repoRoot/.gitignore has a line for entry
// (docs/design.md §28.1). A line equal to entry or "/"+entry, after
// trimming whitespace, counts as present. The file is only ever
// appended to or created, never rewritten.
func EnsureIgnored(repoRoot, entry string) (changed bool, err error) {
	path := filepath.Join(repoRoot, ".gitignore")
	data, readErr := os.ReadFile(path)
	if readErr != nil && !errors.Is(readErr, fs.ErrNotExist) {
		return false, fmt.Errorf("tb: %w", readErr)
	}

	for line := range strings.SplitSeq(string(data), "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == entry || trimmed == "/"+entry {
			return false, nil
		}
	}

	block := ignoreComment + "\n" + entry + "\n"
	if len(data) > 0 && !strings.HasSuffix(string(data), "\n") {
		block = "\n" + block
	}

	f, openErr := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if openErr != nil {
		return false, fmt.Errorf("tb: %w", openErr)
	}
	if _, writeErr := f.WriteString(block); writeErr != nil {
		f.Close()
		return false, fmt.Errorf("tb: %w", writeErr)
	}
	if closeErr := f.Close(); closeErr != nil {
		return false, fmt.Errorf("tb: %w", closeErr)
	}
	return true, nil
}

// MaterializeTemplates copies each of templatePaths (relative to
// $TOOLBOX_HOME/templates, per docs/design.md §3.3) into repoRoot at the
// same relative path, but only if that destination does not already
// exist. Files that already exist are left completely alone — no
// diffing, no prompt — which is what keeps existing docs/adr/*.md files
// (and any other pre-existing doc) safe on every rerun, --force
// included: `tb init` can only ever create these files, never overwrite
// them.
//
// "AGENTS.md" is skipped even if present in templatePaths: it is handled
// separately by the caller via internal/render, not copied verbatim.
func MaterializeTemplates(toolboxHome, repoRoot string, templatePaths []string) (created, skipped []string, err error) {
	for _, rel := range templatePaths {
		if rel == "AGENTS.md" {
			continue
		}

		dst := filepath.Join(repoRoot, rel)
		if _, statErr := os.Stat(dst); statErr == nil {
			skipped = append(skipped, rel)
			continue
		}

		src := filepath.Join(toolboxHome, "templates", rel)
		data, readErr := os.ReadFile(src)
		if readErr != nil {
			return created, skipped, fmt.Errorf("tb: %w", readErr)
		}
		if mkErr := os.MkdirAll(filepath.Dir(dst), 0o755); mkErr != nil {
			return created, skipped, fmt.Errorf("tb: %w", mkErr)
		}
		if writeErr := os.WriteFile(dst, data, 0o644); writeErr != nil {
			return created, skipped, fmt.Errorf("tb: %w", writeErr)
		}
		created = append(created, rel)
	}
	return created, skipped, nil
}
