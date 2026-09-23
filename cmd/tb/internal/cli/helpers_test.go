package cli

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// newFixtureToolboxHome builds a minimal but complete $TOOLBOX_HOME under
// t.TempDir(): one skill, one profile, the default persona, and the
// templates a `tb init` needs — everything the cli package's commands
// read.
func newFixtureToolboxHome(t *testing.T) string {
	t.Helper()
	home := t.TempDir()

	writeFile(t, filepath.Join(home, "skills", "spec", "SKILL.md"), "---\nname: spec\ndescription: test fixture\n---\n\n# spec\n")
	writeFile(t, filepath.Join(home, "skills", "adr", "SKILL.md"), "---\nname: adr\ndescription: test fixture\n---\n\n# adr\n")

	writeFile(t, filepath.Join(home, "personas", "default.md"), ""+
		"# Default persona\n\n"+
		"## Plan — Staff Engineer\n\nChallenges weak specs.\n\n"+
		"## Implement — Senior Engineer\n\nFollows accepted ADRs.\n")

	writeFile(t, filepath.Join(home, "profiles", "rust-systems.toml"), `name = "rust-systems"
description = "fixture"
persona = "default"
skills = ["spec", "adr"]
templates = ["AGENTS.md", "docs/requirements.md", "docs/adr/0000-template.md"]

[verify]
summary = "cargo test"
`)
	writeFile(t, filepath.Join(home, "profiles", "infra-go.toml"), `name = "infra-go"
description = "fixture"
persona = "default"
skills = ["spec", "adr"]
templates = ["AGENTS.md", "docs/requirements.md", "docs/adr/0000-template.md"]

[verify]
summary = "go test ./..."
`)
	writeFile(t, filepath.Join(home, "profiles", "back-go.toml"), `name = "back-go"
description = "fixture"
persona = "default"
skills = ["spec", "adr"]
templates = ["AGENTS.md", "docs/requirements.md", "docs/adr/0000-template.md"]

[verify]
summary = "go test ./..."
`)
	writeFile(t, filepath.Join(home, "profiles", "game-bevy.toml"), `name = "game-bevy"
description = "fixture"
persona = "default"
skills = ["spec", "adr"]
templates = ["AGENTS.md", "docs/requirements.md", "docs/adr/0000-template.md"]

[verify]
summary = "cargo test"
`)
	writeFile(t, filepath.Join(home, "profiles", "cpp-systems.toml"), `name = "cpp-systems"
description = "fixture"
persona = "default"
skills = ["spec", "adr"]
templates = ["AGENTS.md", "docs/requirements.md", "docs/adr/0000-template.md"]

[verify]
summary = "ctest --test-dir build"
`)
	writeFile(t, filepath.Join(home, "profiles", "c-cli.toml"), `name = "c-cli"
description = "fixture"
persona = "default"
skills = ["spec", "adr"]
templates = ["AGENTS.md", "docs/requirements.md", "docs/adr/0000-template.md"]

[verify]
summary = "make test"
`)

	writeFile(t, filepath.Join(home, "templates", "AGENTS.md"), ""+
		beginMarkerForTest+"\n"+
		"Profile: {{.Profile}}\n"+
		"Persona: {{.PersonaPlanTitle}} / {{.PersonaImplementTitle}}\n"+
		"Skills: {{.SkillList}}\n"+
		"Verify: {{.VerifySummary}}\n"+
		endMarkerForTest+"\n")
	writeFile(t, filepath.Join(home, "templates", "docs", "requirements.md"), "# Requirements\n\nFill this in.\n")
	writeFile(t, filepath.Join(home, "templates", "docs", "adr", "0000-template.md"), "# NNNN. Title\n")

	return home
}

// These mirror internal/render's unexported marker constants; kept local
// so this test fixture doesn't need to reach into another package.
const (
	beginMarkerForTest = "<!-- toolbox:begin -->"
	endMarkerForTest   = "<!-- toolbox:end -->"
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

// newGitRepo creates an empty git repository under t.TempDir() and
// returns its path. Tests that call Init or PersonaShow must t.Chdir
// into it first, since those commands discover the repo root / read
// toolbox.toml relative to the process's current directory.
func newGitRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	runGit(t, dir, "init", "-q")
	runGit(t, dir, "config", "user.email", "test@example.com")
	runGit(t, dir, "config", "user.name", "Test")
	return dir
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

// withFakeAgentsOnPath points PATH at a directory containing fake
// executables named after the given binaries, so agents.Detect() sees
// exactly those as present via PATH regardless of what is actually
// installed on the machine running the tests.
func withFakeAgentsOnPath(t *testing.T, binaries ...string) {
	t.Helper()
	dir := t.TempDir()
	for _, name := range binaries {
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte("#!/bin/sh\n"), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	t.Setenv("PATH", dir)
}

// captureOutput redirects os.Stdout and os.Stderr for the duration of fn
// and returns what was written to each, plus fn's return value. The cli
// package's commands write straight to os.Stdout/os.Stderr rather than
// an injected writer, so tests that need to assert on printed text go
// through this.
func captureOutput(t *testing.T, fn func() int) (stdout, stderr string, code int) {
	t.Helper()

	outR, outW, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	errR, errW, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	oldOut, oldErr := os.Stdout, os.Stderr
	os.Stdout, os.Stderr = outW, errW
	t.Cleanup(func() { os.Stdout, os.Stderr = oldOut, oldErr })

	var outBuf, errBuf bytes.Buffer
	outDone := make(chan struct{})
	errDone := make(chan struct{})
	go func() { io.Copy(&outBuf, outR); close(outDone) }()
	go func() { io.Copy(&errBuf, errR); close(errDone) }()

	code = fn()

	outW.Close()
	errW.Close()
	<-outDone
	<-errDone
	os.Stdout, os.Stderr = oldOut, oldErr

	return outBuf.String(), errBuf.String(), code
}
