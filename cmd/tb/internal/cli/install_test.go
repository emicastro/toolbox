package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInstall_LinksSkillsIntoBothTargets(t *testing.T) {
	toolboxHome := newFixtureToolboxHome(t)
	fakeHome := t.TempDir()
	t.Setenv("HOME", fakeHome)
	withFakeAgentsOnPath(t, "claude", "grok")

	_, _, code := captureOutput(t, func() int { return Install(toolboxHome, nil) })
	if code != 0 {
		t.Fatalf("Install() code = %d, want 0", code)
	}

	for _, dir := range []string{".claude", ".grok"} {
		for _, skill := range []string{"spec", "adr"} {
			link := filepath.Join(fakeHome, dir, "skills", skill)
			target, err := os.Readlink(link)
			if err != nil {
				t.Errorf("Readlink(%s) error = %v, want a symlink", link, err)
				continue
			}
			want := filepath.Join(toolboxHome, "skills", skill)
			if target != want {
				t.Errorf("%s -> %q, want %q", link, target, want)
			}
		}
	}
}

func TestInstall_IsIdempotent(t *testing.T) {
	toolboxHome := newFixtureToolboxHome(t)
	fakeHome := t.TempDir()
	t.Setenv("HOME", fakeHome)
	withFakeAgentsOnPath(t, "claude", "grok")

	_, _, code := captureOutput(t, func() int { return Install(toolboxHome, nil) })
	if code != 0 {
		t.Fatalf("first Install() code = %d, want 0", code)
	}
	_, _, code = captureOutput(t, func() int { return Install(toolboxHome, nil) })
	if code != 0 {
		t.Fatalf("second Install() code = %d, want 0", code)
	}

	link := filepath.Join(fakeHome, ".claude", "skills", "spec")
	target, err := os.Readlink(link)
	if err != nil {
		t.Fatalf("Readlink(%s) error = %v", link, err)
	}
	if want := filepath.Join(toolboxHome, "skills", "spec"); target != want {
		t.Errorf("%s -> %q, want %q", link, target, want)
	}
}

func TestInstall_SkipsUnmanagedFileWithoutAborting(t *testing.T) {
	toolboxHome := newFixtureToolboxHome(t)
	fakeHome := t.TempDir()
	t.Setenv("HOME", fakeHome)
	withFakeAgentsOnPath(t, "claude", "grok")

	// A real file already sits where the "spec" skill would be linked
	// under Claude Code.
	blocked := filepath.Join(fakeHome, ".claude", "skills", "spec")
	writeFile(t, blocked, "not managed by tb\n")

	stdout, stderr, code := captureOutput(t, func() int { return Install(toolboxHome, nil) })
	if code != 0 {
		t.Fatalf("Install() code = %d, want 0 (a skipped skill should not abort the run)", code)
	}
	if !strings.Contains(stderr, "not managed by tb") {
		t.Errorf("stderr = %q, want a warning mentioning the unmanaged file", stderr)
	}

	// The blocked file must survive untouched.
	got, err := os.ReadFile(blocked)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "not managed by tb\n" {
		t.Errorf("blocked file was modified: got %q", got)
	}

	// The other skill, and this skill under the other target, must still
	// link correctly — one conflict does not block the rest.
	other, err := os.Readlink(filepath.Join(fakeHome, ".claude", "skills", "adr"))
	if err != nil {
		t.Errorf("Readlink(adr under Claude Code) error = %v, want it linked despite the spec conflict", err)
	}
	if want := filepath.Join(toolboxHome, "skills", "adr"); other != want {
		t.Errorf("adr link = %q, want %q", other, want)
	}
	if _, err := os.Readlink(filepath.Join(fakeHome, ".grok", "skills", "spec")); err != nil {
		t.Errorf("Readlink(spec under Grok Build) error = %v, want it linked (only Claude Code's target was blocked)", err)
	}

	_ = stdout
}

func TestInstall_WarnsOnOneMissingAgentButSucceeds(t *testing.T) {
	toolboxHome := newFixtureToolboxHome(t)
	t.Setenv("HOME", t.TempDir())
	withFakeAgentsOnPath(t, "claude") // grok absent

	_, stderr, code := captureOutput(t, func() int { return Install(toolboxHome, nil) })
	if code != 0 {
		t.Errorf("Install() code = %d, want 0 (one missing agent is a warning, not a failure)", code)
	}
	if !strings.Contains(stderr, "Grok Build") {
		t.Errorf("stderr = %q, want a warning naming Grok Build", stderr)
	}
}

func TestInstall_FailsClosedWhenNoAgentDetected(t *testing.T) {
	toolboxHome := newFixtureToolboxHome(t)
	t.Setenv("HOME", t.TempDir())
	withFakeAgentsOnPath(t /* no binaries */)

	_, stderr, code := captureOutput(t, func() int { return Install(toolboxHome, nil) })
	if code != 1 {
		t.Errorf("Install() code = %d, want 1 (neither agent detected)", code)
	}
	if !strings.Contains(stderr, "neither Claude Code nor Grok Build") {
		t.Errorf("stderr = %q, want the neither-detected error", stderr)
	}
}
