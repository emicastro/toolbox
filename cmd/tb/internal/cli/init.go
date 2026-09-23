package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"toolbox/cmd/tb/internal/config"
	"toolbox/cmd/tb/internal/render"
	"toolbox/cmd/tb/internal/scaffold"
)

// tbVersion is written into a product repo's toolbox.toml
// (toolbox_version, docs/design.md §3.2) — informational only.
const tbVersion = "1.5.0"

// Init implements `tb init` (docs/design.md §5.2).
func Init(toolboxHome string, args []string) int {
	profileName, force, parseErr := parseInitFlags(args)
	if parseErr != nil || profileName == "" {
		fmt.Fprintln(os.Stderr, "usage: tb init -p <rust-systems|infra-go|back-go|game-bevy|cpp-systems|c-cli> [--force]")
		return 2
	}

	profile, err := config.LoadProfileByName(toolboxHome, profileName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	repoRoot, err := gitRepoRoot()
	if err != nil {
		fmt.Fprintln(os.Stderr, "tb init must run inside a git repository")
		return 1
	}

	personaPath := filepath.Join(toolboxHome, "personas", profile.Persona+".md")
	planTitle, implementTitle, err := render.PersonaTitles(personaPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	region, err := render.Render(toolboxHome, render.TemplateData{
		Profile:               profile.Name,
		PersonaPlanTitle:      planTitle,
		PersonaImplementTitle: implementTitle,
		SkillList:             render.SkillList(profile.Skills),
		VerifySummary:         profile.VerifySummary,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	pointerPath := filepath.Join(repoRoot, "toolbox.toml")
	newPointer := renderPointer(profile.Name)
	existingPointer, pointerExisted := readIfExists(pointerPath)
	pointerChanges := pointerExisted && existingPointer != newPointer

	agentsPath := filepath.Join(repoRoot, "AGENTS.md")
	existingAgents, agentsExisted := readIfExists(agentsPath)
	newAgents, _ := render.Merge(existingAgents, region)
	agentsChanges := agentsExisted && newAgents != existingAgents

	if (pointerChanges || agentsChanges) && !force {
		prompt := fmt.Sprintf("tb init: %s already exist and would change.", existingFilesLabel(pointerChanges, agentsChanges, pointerPath, agentsPath))
		if !confirmMerge(force, prompt) {
			printRefused(pointerPath, newPointer, agentsPath, region)
			return 1
		}
	}

	if err := os.WriteFile(pointerPath, []byte(newPointer), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, "tb:", err)
		return 1
	}
	if err := os.WriteFile(agentsPath, []byte(newAgents), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, "tb:", err)
		return 1
	}

	created, skipped, err := scaffold.MaterializeTemplates(toolboxHome, repoRoot, profile.Templates)
	if err != nil {
		fmt.Fprintln(os.Stderr, "tb:", err)
		return 1
	}

	printInitSummary(
		profile.Name,
		pointerPath, pointerExisted && existingPointer == newPointer,
		agentsPath, agentsExisted && existingAgents == newAgents,
		created, skipped,
	)
	return 0
}

func gitRepoRoot() (string, error) {
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func renderPointer(profileName string) string {
	return fmt.Sprintf(
		"# Written by `tb init`. Safe to hand-edit the `profile` line; the\n"+
			"# other two lines are refreshed by --force.\n"+
			"profile = %q\n"+
			"toolbox_version = %q\n"+
			"generated_by = \"tb init\"\n",
		profileName, tbVersion,
	)
}

func readIfExists(path string) (content string, existed bool) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", false
	}
	return string(data), true
}

func existingFilesLabel(pointerChanges, agentsChanges bool, pointerPath, agentsPath string) string {
	switch {
	case pointerChanges && agentsChanges:
		return pointerPath + " and " + agentsPath
	case pointerChanges:
		return pointerPath
	default:
		return agentsPath
	}
}

func printRefused(pointerPath, newPointer, agentsPath, region string) {
	fmt.Fprintln(os.Stderr, "refused: no files were changed. tb init would have written:")
	fmt.Fprintf(os.Stderr, "--- %s ---\n", pointerPath)
	fmt.Fprint(os.Stderr, newPointer)
	fmt.Fprintf(os.Stderr, "--- %s managed region ---\n", agentsPath)
	fmt.Fprint(os.Stderr, region)
}

func printInitSummary(profileName, pointerPath string, pointerUnchanged bool, agentsPath string, agentsUnchanged bool, created, skipped []string) {
	fmt.Printf("profile: %s\n", profileName)
	printWroteOrUnchanged(pointerPath, pointerUnchanged)
	printWroteOrUnchanged(agentsPath, agentsUnchanged)
	for _, c := range created {
		fmt.Printf("created: %s\n", c)
	}
	for _, s := range skipped {
		fmt.Printf("left alone (already exists): %s\n", s)
	}
}

func printWroteOrUnchanged(path string, unchanged bool) {
	if unchanged {
		fmt.Printf("unchanged: %s\n", path)
		return
	}
	fmt.Printf("wrote: %s\n", path)
}
