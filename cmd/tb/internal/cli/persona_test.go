package cli

import (
	"strings"
	"testing"
)

func TestPersonaShow_ReadsProfileFromCwdToolboxTOML(t *testing.T) {
	toolboxHome := newFixtureToolboxHome(t)
	repo := newGitRepo(t)
	t.Chdir(repo)
	writeFile(t, "toolbox.toml", `profile = "rust-systems"`+"\n")

	stdout, _, code := captureOutput(t, func() int { return PersonaShow(toolboxHome, nil) })
	if code != 0 {
		t.Fatalf("PersonaShow() code = %d, want 0", code)
	}
	if !strings.Contains(stdout, "Plan section: Plan — Staff Engineer") {
		t.Errorf("stdout = %q, want the exact Plan section title", stdout)
	}
	if !strings.Contains(stdout, "Implement section: Implement — Senior Engineer") {
		t.Errorf("stdout = %q, want the exact Implement section title", stdout)
	}
}

func TestPersonaShow_ProfileFlagOverridesCwd(t *testing.T) {
	toolboxHome := newFixtureToolboxHome(t)
	repo := newGitRepo(t)
	t.Chdir(repo) // no toolbox.toml at all

	stdout, _, code := captureOutput(t, func() int { return PersonaShow(toolboxHome, []string{"-p", "infra-go"}) })
	if code != 0 {
		t.Fatalf("PersonaShow() code = %d, want 0", code)
	}
	if !strings.Contains(stdout, "Plan section: Plan — Staff Engineer") {
		t.Errorf("stdout = %q, want the persona section title", stdout)
	}
}

func TestPersonaShow_NoProfileAndNoToolboxTOMLIsUsageError(t *testing.T) {
	toolboxHome := newFixtureToolboxHome(t)
	repo := newGitRepo(t)
	t.Chdir(repo)

	_, stderr, code := captureOutput(t, func() int { return PersonaShow(toolboxHome, nil) })
	if code != 2 {
		t.Errorf("PersonaShow() code = %d, want 2 (usage error)", code)
	}
	if !strings.Contains(stderr, "toolbox.toml") {
		t.Errorf("stderr = %q, want it to mention toolbox.toml", stderr)
	}
}
