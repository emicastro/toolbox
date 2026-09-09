package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"toolbox/cmd/tb/internal/config"
	"toolbox/cmd/tb/internal/render"
)

// PersonaShow implements `tb persona show` (docs/design.md §5.5).
// Read-only: it never writes anything.
func PersonaShow(toolboxHome string, args []string) int {
	profileName, parseErr := parsePersonaShowFlags(args)
	if parseErr != nil {
		fmt.Fprintln(os.Stderr, "usage: tb persona show [-p|--profile <name>]")
		return 2
	}

	if profileName == "" {
		pointer, err := config.LoadPointer("toolbox.toml")
		if err != nil {
			fmt.Fprintln(os.Stderr, "tb persona show: no -p/--profile given and no toolbox.toml in the current directory")
			return 2
		}
		profileName = pointer.Profile
	}

	profile, err := config.LoadProfileByName(toolboxHome, profileName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	personaPath := filepath.Join(toolboxHome, "personas", profile.Persona+".md")
	plan, implement, err := render.PersonaTitles(personaPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	fmt.Printf("persona file: %s\n", personaPath)
	fmt.Printf("Plan section: %s\n", plan)
	fmt.Printf("Implement section: %s\n", implement)
	return 0
}
