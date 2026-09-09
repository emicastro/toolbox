package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSkillNew_CreatesScaffold(t *testing.T) {
	toolboxHome := newFixtureToolboxHome(t)

	stdout, _, code := captureOutput(t, func() int { return SkillNew(toolboxHome, []string{"demo"}) })
	if code != 0 {
		t.Fatalf("SkillNew() code = %d, want 0", code)
	}
	if !strings.Contains(stdout, "next steps:") {
		t.Errorf("stdout = %q, want the next-steps text", stdout)
	}

	skillFile := filepath.Join(toolboxHome, "skills", "demo", "SKILL.md")
	data, err := os.ReadFile(skillFile)
	if err != nil {
		t.Fatalf("expected %s to exist: %v", skillFile, err)
	}
	if !strings.Contains(string(data), "name: demo") {
		t.Errorf("SKILL.md = %q, want it to contain name: demo", data)
	}
}

func TestSkillNew_InvalidNameIsUsageError(t *testing.T) {
	toolboxHome := newFixtureToolboxHome(t)

	_, _, code := captureOutput(t, func() int { return SkillNew(toolboxHome, []string{"BadName"}) })
	if code != 2 {
		t.Errorf("SkillNew() code = %d, want 2 for an invalid name", code)
	}
}

func TestSkillNew_AlreadyExistsFails(t *testing.T) {
	toolboxHome := newFixtureToolboxHome(t)

	if _, _, code := captureOutput(t, func() int { return SkillNew(toolboxHome, []string{"demo"}) }); code != 0 {
		t.Fatalf("first SkillNew() code = %d, want 0", code)
	}

	_, stderr, code := captureOutput(t, func() int { return SkillNew(toolboxHome, []string{"demo"}) })
	if code != 1 {
		t.Errorf("second SkillNew() code = %d, want 1", code)
	}
	if !strings.Contains(stderr, "already exists") {
		t.Errorf("stderr = %q, want it to say already exists", stderr)
	}
}
