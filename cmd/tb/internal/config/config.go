package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// Pointer is the product-repo toolbox.toml schema (docs/design.md §3.2).
// GeneratedBy and ToolboxVersion are informational only; Profile is the
// only field tb reads back.
type Pointer struct {
	Profile        string
	ToolboxVersion string
	GeneratedBy    string
}

// LoadPointer reads a product repo's toolbox.toml.
func LoadPointer(path string) (Pointer, error) {
	doc, err := parseFile(path)
	if err != nil {
		return Pointer{}, err
	}

	profile := doc.strings["profile"]
	if profile == "" {
		return Pointer{}, fmt.Errorf("tb: %s: missing required field %q", path, "profile")
	}

	return Pointer{
		Profile:        profile,
		ToolboxVersion: doc.strings["toolbox_version"],
		GeneratedBy:    doc.strings["generated_by"],
	}, nil
}

// Profile is the profiles/<name>.toml schema (docs/design.md §3.3).
// Templates lists paths relative to $TOOLBOX_HOME/templates, materialized
// relative to a product repo's root by `tb init`. VerifySummary is the
// [verify].summary field: documentation rendered into AGENTS.md, never
// executed by tb (requirements.md §11).
type Profile struct {
	Name          string
	Description   string
	Persona       string
	Skills        []string
	Templates     []string
	VerifySummary string
}

// LoadProfile reads one profiles/<name>.toml file.
func LoadProfile(path string) (Profile, error) {
	doc, err := parseFile(path)
	if err != nil {
		return Profile{}, err
	}

	name := doc.strings["name"]
	if name == "" {
		return Profile{}, fmt.Errorf("tb: %s: missing required field %q", path, "name")
	}
	persona := doc.strings["persona"]
	if persona == "" {
		return Profile{}, fmt.Errorf("tb: %s: missing required field %q", path, "persona")
	}
	skills := doc.arrays["skills"]
	if len(skills) == 0 {
		return Profile{}, fmt.Errorf("tb: %s: missing required field %q", path, "skills")
	}

	return Profile{
		Name:          name,
		Description:   doc.strings["description"],
		Persona:       persona,
		Skills:        skills,
		Templates:     doc.arrays["templates"],
		VerifySummary: doc.sections["verify"]["summary"],
	}, nil
}

// LoadProfileByName loads $TOOLBOX_HOME/profiles/<name>.toml, using the
// not-found error format fixed by docs/design.md §3.4.
func LoadProfileByName(toolboxHome, name string) (Profile, error) {
	dir := filepath.Join(toolboxHome, "profiles")
	path := filepath.Join(dir, name+".toml")
	if _, err := os.Stat(path); err != nil {
		return Profile{}, fmt.Errorf("tb: profile %q not found in %s", name, dir)
	}
	return LoadProfile(path)
}
