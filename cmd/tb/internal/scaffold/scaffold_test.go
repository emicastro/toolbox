package scaffold

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestMaterializeTemplatesCreatesMissingFiles(t *testing.T) {
	toolboxHome := t.TempDir()
	repoRoot := t.TempDir()

	writeFile(t, filepath.Join(toolboxHome, "templates", "docs", "requirements.md"), "# Requirements\n")
	writeFile(t, filepath.Join(toolboxHome, "templates", "docs", "adr", "0000-template.md"), "# NNNN. Title\n")

	created, skipped, err := MaterializeTemplates(toolboxHome, repoRoot,
		[]string{"docs/requirements.md", "docs/adr/0000-template.md"})
	if err != nil {
		t.Fatalf("MaterializeTemplates() error = %v, want nil", err)
	}
	if len(skipped) != 0 {
		t.Errorf("skipped = %v, want empty", skipped)
	}
	if len(created) != 2 {
		t.Fatalf("created = %v, want 2 entries", created)
	}

	got, err := os.ReadFile(filepath.Join(repoRoot, "docs", "requirements.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "# Requirements\n" {
		t.Errorf("docs/requirements.md content = %q, want %q", got, "# Requirements\n")
	}
}

func TestMaterializeTemplatesNeverOverwritesExisting(t *testing.T) {
	toolboxHome := t.TempDir()
	repoRoot := t.TempDir()

	writeFile(t, filepath.Join(toolboxHome, "templates", "docs", "adr", "0000-template.md"), "# NNNN. Title\n")
	existingADR := filepath.Join(repoRoot, "docs", "adr", "0001-something.md")
	writeFile(t, existingADR, "# 0001. Something already decided\n\nDo not touch me.\n")

	created, skipped, err := MaterializeTemplates(toolboxHome, repoRoot, []string{"docs/adr/0000-template.md"})
	if err != nil {
		t.Fatalf("MaterializeTemplates() error = %v, want nil", err)
	}
	if len(created) != 1 || created[0] != "docs/adr/0000-template.md" {
		t.Errorf("created = %v, want [docs/adr/0000-template.md]", created)
	}
	if len(skipped) != 0 {
		t.Errorf("skipped = %v, want empty (0000-template.md did not exist yet)", skipped)
	}

	got, err := os.ReadFile(existingADR)
	if err != nil {
		t.Fatal(err)
	}
	want := "# 0001. Something already decided\n\nDo not touch me.\n"
	if string(got) != want {
		t.Errorf("existing ADR was modified: got %q, want %q", got, want)
	}
}

func TestMaterializeTemplatesSkipsAgentsMD(t *testing.T) {
	toolboxHome := t.TempDir()
	repoRoot := t.TempDir()
	writeFile(t, filepath.Join(toolboxHome, "templates", "AGENTS.md"), "should never be copied verbatim\n")

	created, skipped, err := MaterializeTemplates(toolboxHome, repoRoot, []string{"AGENTS.md"})
	if err != nil {
		t.Fatalf("MaterializeTemplates() error = %v, want nil", err)
	}
	if len(created) != 0 || len(skipped) != 0 {
		t.Errorf("created = %v, skipped = %v, want both empty (AGENTS.md is always skipped)", created, skipped)
	}
	if _, err := os.Stat(filepath.Join(repoRoot, "AGENTS.md")); err == nil {
		t.Error("AGENTS.md was created by MaterializeTemplates, want it left to the caller")
	}
}

func TestMaterializeTemplatesLeavesExistingFileAlone(t *testing.T) {
	toolboxHome := t.TempDir()
	repoRoot := t.TempDir()
	writeFile(t, filepath.Join(toolboxHome, "templates", "docs", "session.md"), "template content\n")
	writeFile(t, filepath.Join(repoRoot, "docs", "session.md"), "user's own session notes\n")

	created, skipped, err := MaterializeTemplates(toolboxHome, repoRoot, []string{"docs/session.md"})
	if err != nil {
		t.Fatalf("MaterializeTemplates() error = %v, want nil", err)
	}
	if len(created) != 0 {
		t.Errorf("created = %v, want empty", created)
	}
	if len(skipped) != 1 || skipped[0] != "docs/session.md" {
		t.Errorf("skipped = %v, want [docs/session.md]", skipped)
	}

	got, err := os.ReadFile(filepath.Join(repoRoot, "docs", "session.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "user's own session notes\n" {
		t.Errorf("existing docs/session.md was modified: got %q", got)
	}
}
