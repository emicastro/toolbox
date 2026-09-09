package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadPointer(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "toolbox.toml")
	content := "profile = \"rust-systems\"\n" +
		"toolbox_version = \"0.1.0\"\n" +
		"generated_by = \"tb init\"\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := LoadPointer(path)
	if err != nil {
		t.Fatalf("LoadPointer() error = %v, want nil", err)
	}
	want := Pointer{Profile: "rust-systems", ToolboxVersion: "0.1.0", GeneratedBy: "tb init"}
	if got != want {
		t.Errorf("LoadPointer() = %#v, want %#v", got, want)
	}
}

func TestLoadPointerMissingProfile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "toolbox.toml")
	if err := os.WriteFile(path, []byte(`toolbox_version = "0.1.0"`), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := LoadPointer(path); err == nil {
		t.Fatal("LoadPointer() error = nil, want an error for a missing profile field")
	}
}

func TestLoadProfile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "rust-systems.toml")
	content := `name = "rust-systems"
description = "Systems / low-level Rust, greenfield and legacy"
persona = "default"
skills = [
  "spec", "adr", "onboard", "scout", "handoff",
  "rust-verify", "rust-systems",
]
templates = ["AGENTS.md", "docs/requirements.md"]

[verify]
summary = "cargo test; cargo clippy -- -D warnings; miri when unsafe changed"
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := LoadProfile(path)
	if err != nil {
		t.Fatalf("LoadProfile() error = %v, want nil", err)
	}

	want := Profile{
		Name:        "rust-systems",
		Description: "Systems / low-level Rust, greenfield and legacy",
		Persona:     "default",
		Skills: []string{
			"spec", "adr", "onboard", "scout", "handoff",
			"rust-verify", "rust-systems",
		},
		Templates:     []string{"AGENTS.md", "docs/requirements.md"},
		VerifySummary: "cargo test; cargo clippy -- -D warnings; miri when unsafe changed",
	}
	if got.Name != want.Name || got.Description != want.Description ||
		got.Persona != want.Persona || got.VerifySummary != want.VerifySummary {
		t.Errorf("LoadProfile() = %#v, want %#v", got, want)
	}
	if len(got.Skills) != len(want.Skills) {
		t.Fatalf("LoadProfile() Skills = %#v, want %#v", got.Skills, want.Skills)
	}
	for i := range want.Skills {
		if got.Skills[i] != want.Skills[i] {
			t.Errorf("LoadProfile() Skills[%d] = %q, want %q", i, got.Skills[i], want.Skills[i])
		}
	}
}

func TestLoadProfileMissingRequiredField(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{"missing name", `persona = "default"` + "\n" + `skills = ["spec"]`},
		{"missing persona", `name = "x"` + "\n" + `skills = ["spec"]`},
		{"missing skills", `name = "x"` + "\n" + `persona = "default"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "profile.toml")
			if err := os.WriteFile(path, []byte(tt.content), 0o644); err != nil {
				t.Fatal(err)
			}
			if _, err := LoadProfile(path); err == nil {
				t.Fatalf("LoadProfile() error = nil, want an error (%s)", tt.name)
			}
		})
	}
}

func TestLoadProfileByName(t *testing.T) {
	toolboxHome := t.TempDir()
	profilesDir := filepath.Join(toolboxHome, "profiles")
	if err := os.MkdirAll(profilesDir, 0o755); err != nil {
		t.Fatal(err)
	}
	content := `name = "infra-go"
persona = "default"
skills = ["spec", "adr"]
`
	if err := os.WriteFile(filepath.Join(profilesDir, "infra-go.toml"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := LoadProfileByName(toolboxHome, "infra-go")
	if err != nil {
		t.Fatalf("LoadProfileByName() error = %v, want nil", err)
	}
	if got.Name != "infra-go" {
		t.Errorf("LoadProfileByName() Name = %q, want %q", got.Name, "infra-go")
	}

	_, err = LoadProfileByName(toolboxHome, "does-not-exist")
	if err == nil {
		t.Fatal("LoadProfileByName() error = nil, want an error for an unknown profile")
	}
	wantSubstr := `profile "does-not-exist" not found in`
	if got, want := err.Error(), wantSubstr; !strings.Contains(got, want) {
		t.Errorf("LoadProfileByName() error = %q, want it to contain %q", got, want)
	}
}

// TestRealProfilesParse is a regression check tying this parser to the
// actual profiles this repo ships (profiles/rust-systems.toml,
// profiles/infra-go.toml) and its own toolbox.toml, so a future edit to
// either the parser or those files that breaks compatibility fails here
// rather than only at `tb init` time.
func TestRealProfilesParse(t *testing.T) {
	repoRoot := "../../../.."

	for _, name := range []string{"rust-systems", "infra-go"} {
		path := filepath.Join(repoRoot, "profiles", name+".toml")
		if _, err := os.Stat(path); err != nil {
			t.Skipf("real profile %s not found at %s (repo layout changed?): %v", name, path, err)
		}
		got, err := LoadProfile(path)
		if err != nil {
			t.Errorf("LoadProfile(%s) error = %v, want nil", path, err)
			continue
		}
		if got.Name != name {
			t.Errorf("LoadProfile(%s).Name = %q, want %q", path, got.Name, name)
		}
		if len(got.Skills) == 0 {
			t.Errorf("LoadProfile(%s).Skills is empty, want the process + domain skill list", path)
		}
		if got.VerifySummary == "" {
			t.Errorf("LoadProfile(%s).VerifySummary is empty, want [verify].summary", path)
		}
	}

	toolboxTOML := filepath.Join(repoRoot, "toolbox.toml")
	if _, err := os.Stat(toolboxTOML); err == nil {
		if _, err := parseFile(toolboxTOML); err != nil {
			t.Errorf("parsing this repo's own toolbox.toml failed: %v", err)
		}
	}
}
