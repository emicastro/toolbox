package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

var skillNamePattern = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)

// SkillNew implements `tb skill new <name>` (docs/design.md §5.4). It
// does not attach the new skill to any profile — that stays a manual
// edit of profiles/<name>.toml followed by `tb install`, per
// requirements.md §8.4.
func SkillNew(toolboxHome string, args []string) int {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: tb skill new <name>")
		return 2
	}

	name := args[0]
	if !skillNamePattern.MatchString(name) {
		fmt.Fprintf(os.Stderr, "tb skill new: invalid name %q (must match ^[a-z][a-z0-9-]*$)\n", name)
		return 2
	}

	skillDir := filepath.Join(toolboxHome, "skills", name)
	skillFile := filepath.Join(skillDir, "SKILL.md")
	if _, err := os.Stat(skillFile); err == nil {
		fmt.Fprintf(os.Stderr, "tb skill new: %s already exists\n", skillFile)
		return 1
	}

	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		fmt.Fprintln(os.Stderr, "tb:", err)
		return 1
	}

	content := fmt.Sprintf(
		"---\nname: %s\ndescription: TODO — one line, when an agent should reach for this skill.\n---\n\n# %s\n\nTODO: body.\n",
		name, name,
	)
	if err := os.WriteFile(skillFile, []byte(content), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, "tb:", err)
		return 1
	}

	fmt.Printf("created %s\n\n", skillFile)
	fmt.Println("next steps:")
	fmt.Println("  1. edit the skill")
	fmt.Println("  2. add its name to the chosen profile's skills array")
	fmt.Println("  3. run tb install")
	fmt.Println("  4. commit in $TOOLBOX_HOME")

	return 0
}
