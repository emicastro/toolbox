// Package render renders the managed region of a product repo's
// AGENTS.md and merges it into an existing file (docs/design.md §4,
// docs/adr/0002-agents-md-managed-region.md).
package render

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	beginMarker = "<!-- toolbox:begin -->"
	endMarker   = "<!-- toolbox:end -->"
)

// TemplateData is the substitution set for $TOOLBOX_HOME/templates/AGENTS.md.
type TemplateData struct {
	Profile               string
	PersonaPlanTitle      string
	PersonaImplementTitle string
	SkillList             string
	VerifySummary         string
}

// Render executes toolboxHome's templates/AGENTS.md against data and
// returns the rendered managed region, markers included. The template
// file is the literal source of truth for the region's text — this
// function does not duplicate it in Go source, so templates/AGENTS.md
// and what tb writes can never drift apart. Rendering is deterministic:
// the same toolboxHome content and the same data always produce the same
// bytes, so re-running `tb init` with nothing changed makes no diff.
func Render(toolboxHome string, data TemplateData) (string, error) {
	tmplPath := filepath.Join(toolboxHome, "templates", "AGENTS.md")
	raw, err := os.ReadFile(tmplPath)
	if err != nil {
		return "", fmt.Errorf("tb: %w", err)
	}

	tmpl, err := template.New("AGENTS.md").Parse(string(raw))
	if err != nil {
		return "", fmt.Errorf("tb: %s: %w", tmplPath, err)
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("tb: %s: %w", tmplPath, err)
	}
	return buf.String(), nil
}

// SkillList joins a profile's skill names the way TemplateData.SkillList
// expects (docs/design.md §4: "the flattened skill list ... joined with
// ', '").
func SkillList(skills []string) string {
	return strings.Join(skills, ", ")
}

// Merge combines a freshly rendered region (markers included) into an
// existing AGENTS.md's content, implementing the structural half of the
// three-state table in docs/design.md §4:
//
//   - existing == "": region becomes the whole file.
//   - existing has both markers: only the text between them (inclusive)
//     is replaced; everything else survives byte-for-byte, and repeated
//     merges with the same region are idempotent (no growing blank lines).
//   - existing has no markers: region is prepended above the existing
//     content, separated by one blank line.
//
// refused is always false. Merge only computes a structural result from
// content; the policy decision in §4's table — prompt, refuse
// non-interactively, or force — depends on terminal interactivity and the
// --force flag, neither of which this package knows about. That decision
// belongs to the cli package (docs/tasks.md group 4): it decides whether
// to call Merge at all, and never calls it when the outcome is refusal.
// The bool return is kept so a future in-package refusal condition (none
// exists today) doesn't require a signature change.
func Merge(existing, region string) (result string, refused bool) {
	if existing == "" {
		return region, false
	}

	beginIdx := strings.Index(existing, beginMarker)
	endIdx := strings.Index(existing, endMarker)

	if beginIdx == -1 || endIdx == -1 || endIdx < beginIdx {
		return region + "\n" + existing, false
	}

	tailStart := endIdx + len(endMarker)
	if tailStart < len(existing) && existing[tailStart] == '\n' {
		tailStart++ // swallow the marker's own line ending so repeated merges don't grow a blank line
	}
	return existing[:beginIdx] + region + existing[tailStart:], false
}

// HasManagedRegion reports whether content contains a complete
// toolbox-managed region (used by `tb doctor`, design.md §5.3 step 5).
func HasManagedRegion(content string) bool {
	beginIdx := strings.Index(content, beginMarker)
	endIdx := strings.Index(content, endMarker)
	return beginIdx != -1 && endIdx != -1 && beginIdx < endIdx
}

// PersonaTitles reads a persona markdown file (docs/design.md §9) and
// returns its first two "## " headings as the Plan and Implement section
// titles. Parsing is generic — it does not require the heading text to
// match personas/default.md verbatim — so a future second persona file
// works without a code change.
func PersonaTitles(path string) (plan, implement string, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", "", fmt.Errorf("tb: %w", err)
	}

	var headings []string
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "## ") {
			headings = append(headings, strings.TrimSpace(strings.TrimPrefix(line, "## ")))
			if len(headings) == 2 {
				break
			}
		}
	}
	if len(headings) < 2 {
		return "", "", fmt.Errorf("tb: %s: expected two \"## \" headings, found %d", path, len(headings))
	}
	return headings[0], headings[1], nil
}
