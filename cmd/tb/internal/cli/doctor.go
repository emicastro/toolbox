package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"toolbox/cmd/tb/internal/agents"
	"toolbox/cmd/tb/internal/config"
	"toolbox/cmd/tb/internal/render"
)

// Doctor implements `tb doctor` (docs/design.md §5.3). It is read-only
// and, unlike every other command, must still print a full one-screen
// report when the toolbox home failed to resolve (homeErr != nil,
// toolboxHome == "") — that failure becomes this report's first line
// instead of a hard exit (docs/design.md §1).
func Doctor(toolboxHome string, homeErr error, args []string) int {
	profileOverride, parseErr := parseDoctorFlags(args)
	if parseErr != nil {
		fmt.Fprintln(os.Stderr, "usage: tb doctor [-p|--profile <name>]")
		return 2
	}

	fmt.Println("tb doctor")
	fmt.Println()

	// 1. toolbox home
	if homeErr != nil {
		fmt.Printf("toolbox home: FAIL - %v\n", homeErr)
	} else {
		fmt.Printf("toolbox home: ok (%s)\n", toolboxHome)
	}

	// 2. tb on PATH (informational only, never fails doctor)
	if _, err := exec.LookPath("tb"); err == nil {
		fmt.Println("tb on PATH: ok")
	} else {
		fmt.Println("tb on PATH: not found (informational only)")
	}

	// 3. skill symlink health
	if toolboxHome == "" {
		fmt.Println("skill symlinks: skipped (toolbox home missing)")
	} else {
		reportSkillLinks(toolboxHome)
	}

	// 4. agent detection - doctor's only source of a non-zero exit code
	claude, grok := agents.Detect()
	fmt.Printf("Claude Code: %s\n", presence(claude))
	fmt.Printf("Grok Build: %s\n", presence(grok))

	// 5. cwd toolbox.toml / profile
	profileName := reportCwdProfile(toolboxHome)

	// 6. language toolchain warning
	if profileOverride != "" {
		profileName = profileOverride
	}
	reportToolchain(profileName)

	if !claude && !grok {
		return 1
	}
	return 0
}

func reportSkillLinks(toolboxHome string) {
	names, err := listSkills(toolboxHome)
	if err != nil {
		fmt.Printf("skill symlinks: FAIL - %v\n", err)
		return
	}
	if len(names) == 0 {
		fmt.Println("skill symlinks: no skills found under $TOOLBOX_HOME/skills")
		return
	}

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("skill symlinks: FAIL - %v\n", err)
		return
	}

	for _, tgt := range installTargets(home) {
		for _, name := range names {
			src := filepath.Join(toolboxHome, "skills", name)
			dst := filepath.Join(tgt.dir, name)
			fmt.Printf("skill %s @ %s: %s\n", name, tgt.label, linkStatus(src, dst))
		}
	}
}

// linkStatus reports one of the four states docs/design.md §5.3 step 3
// names: ok, missing, broken, or "not managed". A symlink that resolves
// but points somewhere other than src is also reported broken — it is
// not correctly managed by tb even though it isn't dangling, and
// `tb install` would relink it on the next run.
func linkStatus(src, dst string) string {
	info, err := os.Lstat(dst)
	if err != nil {
		return "missing"
	}
	if info.Mode()&os.ModeSymlink == 0 {
		return "not managed"
	}
	target, err := os.Readlink(dst)
	if err != nil || target != src {
		return "broken"
	}
	if _, err := os.Stat(dst); err != nil {
		return "broken"
	}
	return "ok"
}

// reportCwdProfile implements docs/design.md §5.3 step 5 and returns the
// profile name found (if any), for step 6 to use.
func reportCwdProfile(toolboxHome string) string {
	pointer, err := config.LoadPointer("toolbox.toml")
	if err != nil {
		fmt.Println("cwd profile: no toolbox.toml in the current directory")
		return ""
	}

	fmt.Printf("cwd profile: %s\n", pointer.Profile)
	if toolboxHome == "" {
		fmt.Println("profile resolves: skipped (toolbox home missing)")
	} else if _, err := config.LoadProfileByName(toolboxHome, pointer.Profile); err != nil {
		fmt.Printf("profile resolves: FAIL - %v\n", err)
	} else {
		fmt.Println("profile resolves: ok")
	}

	data, err := os.ReadFile("AGENTS.md")
	switch {
	case err != nil:
		fmt.Println("AGENTS.md: missing")
	case render.HasManagedRegion(string(data)):
		fmt.Println("AGENTS.md managed region: ok")
	default:
		fmt.Println("AGENTS.md managed region: missing")
	}

	return pointer.Profile
}

// reportToolchain implements docs/design.md §5.3 step 6.
func reportToolchain(profileName string) {
	switch profileName {
	case "rust-systems":
		if _, err := exec.LookPath("cargo"); err != nil {
			fmt.Println("warn: cargo not found on PATH")
		}
	case "infra-go", "back-go":
		if _, err := exec.LookPath("go"); err != nil {
			fmt.Println("warn: go not found on PATH")
		}
	}
}
