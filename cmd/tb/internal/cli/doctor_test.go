package cli

import (
	"os"
	"strings"
	"testing"
)

func TestDoctor_CleanWhenBothAgentsPresent(t *testing.T) {
	toolboxHome := newFixtureToolboxHome(t)
	t.Setenv("HOME", t.TempDir())
	withFakeAgentsOnPath(t, "claude", "grok")
	t.Chdir(t.TempDir()) // no toolbox.toml here

	stdout, _, code := captureOutput(t, func() int { return Doctor(toolboxHome, nil, nil) })
	if code != 0 {
		t.Errorf("Doctor() code = %d, want 0", code)
	}
	if !strings.Contains(stdout, "Claude Code: present") || !strings.Contains(stdout, "Grok Build: present") {
		t.Errorf("stdout = %q, want both agents reported present", stdout)
	}
}

func TestDoctor_WarnsButSucceedsWithOneMissingAgent(t *testing.T) {
	toolboxHome := newFixtureToolboxHome(t)
	t.Setenv("HOME", t.TempDir())
	withFakeAgentsOnPath(t, "claude")
	t.Chdir(t.TempDir())

	stdout, _, code := captureOutput(t, func() int { return Doctor(toolboxHome, nil, nil) })
	if code != 0 {
		t.Errorf("Doctor() code = %d, want 0 (one missing agent is not fatal)", code)
	}
	if !strings.Contains(stdout, "Grok Build: absent") {
		t.Errorf("stdout = %q, want Grok Build reported absent", stdout)
	}
}

func TestDoctor_FailsClosedWithZeroAgents(t *testing.T) {
	toolboxHome := newFixtureToolboxHome(t)
	t.Setenv("HOME", t.TempDir())
	withFakeAgentsOnPath(t /* none */)
	t.Chdir(t.TempDir())

	_, _, code := captureOutput(t, func() int { return Doctor(toolboxHome, nil, nil) })
	if code != 1 {
		t.Errorf("Doctor() code = %d, want 1 (zero agents is doctor's only hard failure)", code)
	}
}

func TestDoctor_ReportsMissingToolboxHomeWithoutHardExit(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	withFakeAgentsOnPath(t, "claude", "grok")
	t.Chdir(t.TempDir())

	homeErr := os.ErrNotExist
	stdout, _, code := captureOutput(t, func() int { return Doctor("", homeErr, nil) })
	if code != 0 {
		t.Errorf("Doctor() code = %d, want 0 (missing toolbox home alone is not the zero-agent failure)", code)
	}
	if !strings.Contains(stdout, "toolbox home: FAIL") {
		t.Errorf("stdout = %q, want a toolbox-home failure line", stdout)
	}
	if !strings.Contains(stdout, "skill symlinks: skipped") {
		t.Errorf("stdout = %q, want doctor to keep going and skip the skill-link check", stdout)
	}
	if !strings.Contains(stdout, "Claude Code: present") {
		t.Errorf("stdout = %q, want doctor to still report agent detection", stdout)
	}
}

func TestDoctor_BackGoWarnsOnMissingGoNotCargo(t *testing.T) {
	toolboxHome := newFixtureToolboxHome(t)
	t.Setenv("HOME", t.TempDir())
	// PATH has agents but neither go nor cargo, so the toolchain warning
	// is deterministic regardless of the host's real PATH.
	withFakeAgentsOnPath(t, "claude", "grok")
	t.Chdir(t.TempDir())

	stdout, _, code := captureOutput(t, func() int {
		return Doctor(toolboxHome, nil, []string{"-p", "back-go"})
	})
	if code != 0 {
		t.Errorf("Doctor() code = %d, want 0", code)
	}
	if !strings.Contains(stdout, "warn: go not found on PATH") {
		t.Errorf("stdout = %q, want a missing-go warning for back-go", stdout)
	}
	if strings.Contains(stdout, "cargo") {
		t.Errorf("stdout = %q, want no cargo warning for back-go", stdout)
	}
}

func TestDoctor_SkipsCwdProfileStepsWithoutToolboxTOML(t *testing.T) {
	toolboxHome := newFixtureToolboxHome(t)
	t.Setenv("HOME", t.TempDir())
	withFakeAgentsOnPath(t, "claude", "grok")
	t.Chdir(t.TempDir()) // no toolbox.toml

	stdout, stderr, code := captureOutput(t, func() int { return Doctor(toolboxHome, nil, nil) })
	if code != 0 {
		t.Fatalf("Doctor() code = %d, stderr = %q, want 0", code, stderr)
	}
	if !strings.Contains(stdout, "no toolbox.toml in the current directory") {
		t.Errorf("stdout = %q, want it to say no toolbox.toml found", stdout)
	}
}
