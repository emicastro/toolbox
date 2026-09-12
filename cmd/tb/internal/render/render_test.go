package render

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTemplate(t *testing.T, toolboxHome, content string) {
	t.Helper()
	dir := filepath.Join(toolboxHome, "templates")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "AGENTS.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestRenderGolden(t *testing.T) {
	toolboxHome := t.TempDir()
	writeTemplate(t, toolboxHome, ""+
		beginMarker+"\n"+
		"Profile: {{.Profile}}\n"+
		"Persona: {{.PersonaPlanTitle}} / {{.PersonaImplementTitle}}\n"+
		"Skills: {{.SkillList}}\n"+
		"Verify: {{.VerifySummary}}\n"+
		endMarker+"\n")

	data := TemplateData{
		Profile:               "rust-systems",
		PersonaPlanTitle:      "Plan — Staff Engineer",
		PersonaImplementTitle: "Implement — Senior Engineer",
		SkillList:             SkillList([]string{"spec", "adr", "rust-verify"}),
		VerifySummary:         "cargo test; cargo clippy -- -D warnings",
	}

	got, err := Render(toolboxHome, data)
	if err != nil {
		t.Fatalf("Render() error = %v, want nil", err)
	}

	want := beginMarker + "\n" +
		"Profile: rust-systems\n" +
		"Persona: Plan — Staff Engineer / Implement — Senior Engineer\n" +
		"Skills: spec, adr, rust-verify\n" +
		"Verify: cargo test; cargo clippy -- -D warnings\n" +
		endMarker + "\n"

	if got != want {
		t.Errorf("Render() = %q, want %q", got, want)
	}
}

func TestRenderIsDeterministic(t *testing.T) {
	toolboxHome := t.TempDir()
	writeTemplate(t, toolboxHome, beginMarker+"\nProfile: {{.Profile}}\n"+endMarker+"\n")
	data := TemplateData{Profile: "infra-go"}

	first, err := Render(toolboxHome, data)
	if err != nil {
		t.Fatal(err)
	}
	second, err := Render(toolboxHome, data)
	if err != nil {
		t.Fatal(err)
	}
	if first != second {
		t.Errorf("Render() is not deterministic: %q != %q", first, second)
	}
}

// TestRenderAgainstRealTemplate ties Render to the actual
// templates/AGENTS.md this repo ships, so an edit to either the template
// or the renderer that breaks compatibility fails here.
func TestRenderAgainstRealTemplate(t *testing.T) {
	repoRoot := "../../../.."
	tmplPath := filepath.Join(repoRoot, "templates", "AGENTS.md")
	if _, err := os.Stat(tmplPath); err != nil {
		t.Skipf("real templates/AGENTS.md not found at %s (repo layout changed?): %v", tmplPath, err)
	}

	got, err := Render(repoRoot, TemplateData{
		Profile:               "rust-systems",
		PersonaPlanTitle:      "Plan — Staff Engineer",
		PersonaImplementTitle: "Implement — Senior Engineer",
		SkillList:             "spec, adr, onboard, scout, handoff, rust-verify, rust-systems",
		VerifySummary:         "cargo test; cargo clippy -- -D warnings; miri when unsafe changed",
	})
	if err != nil {
		t.Fatalf("Render(real template) error = %v, want nil", err)
	}
	if !strings.HasPrefix(got, beginMarker) || !strings.Contains(got, endMarker) {
		t.Errorf("Render(real template) = %q, want it wrapped in %s / %s", got, beginMarker, endMarker)
	}
	if !HasManagedRegion(got) {
		t.Error("Render(real template) output does not satisfy HasManagedRegion")
	}

	// ADR 0006: the real template carries the five binding rules, not
	// titles-only. This fails if the rules block is removed from
	// templates/AGENTS.md.
	for _, rule := range []string{
		"Do not implement past the accepted task list in `docs/tasks.md`.",
		"A design fork (two viable options, a dependency, a schema/protocol shape, or reversing an ADR) needs a new ADR before it is treated as settled.",
		"Before ticking a task or claiming done: run the profile verify recipe (success is the command output), then the `review` skill.",
		"Keep diffs small. No drive-by refactors or unrelated formatting.",
		"Do not assume profile defaults in a brownfield repo; map it first (`onboard`). `AGENTS.md` existing is not a map.",
	} {
		if !strings.Contains(got, rule) {
			t.Errorf("Render(real template) missing binding rule %q\noutput:\n%s", rule, got)
		}
	}
}

func TestMergeNoExistingFile(t *testing.T) {
	region := beginMarker + "\ncontent\n" + endMarker + "\n"
	got, refused := Merge("", region)
	if refused {
		t.Error("Merge() refused = true, want false")
	}
	if got != region {
		t.Errorf("Merge() = %q, want %q", got, region)
	}
}

func TestMergeNoMarkersPrepends(t *testing.T) {
	existing := "# My notes\n\nHand-written prose here.\n"
	region := beginMarker + "\nProfile: rust-systems\n" + endMarker + "\n"

	got, refused := Merge(existing, region)
	if refused {
		t.Error("Merge() refused = true, want false")
	}
	if !strings.HasPrefix(got, region) {
		t.Errorf("Merge() = %q, want it to start with the region", got)
	}
	if !strings.HasSuffix(got, existing) {
		t.Errorf("Merge() = %q, want it to end with the original content byte-for-byte", got)
	}
}

func TestMergeReplacesOnlyTheRegion(t *testing.T) {
	before := "# Header\n\n"
	oldRegion := beginMarker + "\nProfile: infra-go\n" + endMarker + "\n"
	after := "\nHand-written prose that must survive.\n"
	existing := before + oldRegion + after

	newRegion := beginMarker + "\nProfile: rust-systems\n" + endMarker + "\n"
	got, refused := Merge(existing, newRegion)
	if refused {
		t.Error("Merge() refused = true, want false")
	}

	want := before + newRegion + after
	if got != want {
		t.Errorf("Merge() = %q, want %q", got, want)
	}
	if !strings.HasPrefix(got, before) {
		t.Error("Merge() did not preserve content before the markers")
	}
	if !strings.HasSuffix(got, after) {
		t.Error("Merge() did not preserve content after the markers")
	}
}

func TestMergeIsIdempotent(t *testing.T) {
	region := beginMarker + "\nProfile: rust-systems\n" + endMarker + "\n"
	existing := "# Header\n\nprose\n"

	once, _ := Merge(existing, region)
	twice, _ := Merge(once, region)

	if once != twice {
		t.Errorf("Merge() is not idempotent:\n first = %q\nsecond = %q", once, twice)
	}
}

func TestHasManagedRegion(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    bool
	}{
		{"both markers present", beginMarker + "\nx\n" + endMarker, true},
		{"no markers", "just prose", false},
		{"only begin marker", beginMarker + "\nx\n", false},
		{"markers reversed", endMarker + "\n" + beginMarker, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasManagedRegion(tt.content); got != tt.want {
				t.Errorf("HasManagedRegion(%q) = %v, want %v", tt.content, got, tt.want)
			}
		})
	}
}

func TestPersonaTitles(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "default.md")
	content := "# Default persona\n\n" +
		"## Plan — Staff Engineer\n\nChallenges weak specs.\n\n" +
		"## Implement — Senior Engineer\n\nFollows accepted ADRs.\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	plan, implement, err := PersonaTitles(path)
	if err != nil {
		t.Fatalf("PersonaTitles() error = %v, want nil", err)
	}
	if plan != "Plan — Staff Engineer" {
		t.Errorf("PersonaTitles() plan = %q, want %q", plan, "Plan — Staff Engineer")
	}
	if implement != "Implement — Senior Engineer" {
		t.Errorf("PersonaTitles() implement = %q, want %q", implement, "Implement — Senior Engineer")
	}
}

func TestPersonaTitlesTooFewHeadings(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "default.md")
	if err := os.WriteFile(path, []byte("# Only a title\n\n## One heading\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, _, err := PersonaTitles(path); err == nil {
		t.Fatal("PersonaTitles() error = nil, want an error for a file with fewer than two headings")
	}
}

// TestPersonaTitlesAgainstRealFile ties PersonaTitles to the actual
// personas/default.md this repo ships.
func TestPersonaTitlesAgainstRealFile(t *testing.T) {
	path := filepath.Join("../../../..", "personas", "default.md")
	if _, err := os.Stat(path); err != nil {
		t.Skipf("real personas/default.md not found at %s (repo layout changed?): %v", path, err)
	}

	plan, implement, err := PersonaTitles(path)
	if err != nil {
		t.Fatalf("PersonaTitles(real file) error = %v, want nil", err)
	}
	if !strings.Contains(plan, "Staff Engineer") {
		t.Errorf("PersonaTitles(real file) plan = %q, want it to mention Staff Engineer", plan)
	}
	if !strings.Contains(implement, "Senior Engineer") {
		t.Errorf("PersonaTitles(real file) implement = %q, want it to mention Senior Engineer", implement)
	}
}
