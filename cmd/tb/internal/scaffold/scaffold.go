// Package scaffold materializes a profile's docs/ templates into a
// product repo (docs/design.md §5.2 step 6, §8.2).
package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
)

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
