package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"toolbox/cmd/tb/internal/agents"
)

// installTarget is one of the two v1 install targets' skill directories
// (docs/design.md §12: Claude Code and Grok Build only in v1).
type installTarget struct {
	label string
	dir   string
}

func installTargets(home string) []installTarget {
	return []installTarget{
		{"Claude Code", filepath.Join(home, ".claude", "skills")},
		{"Grok Build", filepath.Join(home, ".grok", "skills")},
	}
}

// Install implements `tb install` (docs/design.md §5.1).
func Install(toolboxHome string, args []string) int {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "tb:", err)
		return 1
	}

	// Detect agent presence before creating anything: ADR 0003's
	// config-dir signal is "$HOME/.claude exists", and os.MkdirAll of
	// .claude/skills below would create that parent directory as a side
	// effect — detecting after that would make every install
	// self-fulfillingly report both agents present, defeating step 4's
	// whole point.
	claude, grok := agents.Detect()

	skillNames, err := listSkills(toolboxHome)
	if err != nil {
		fmt.Fprintln(os.Stderr, "tb:", err)
		return 1
	}

	targets := installTargets(home)
	linked := make(map[string]int, len(targets))
	for _, tgt := range targets {
		if err := os.MkdirAll(tgt.dir, 0o755); err != nil {
			fmt.Fprintln(os.Stderr, "tb:", err)
			return 1
		}
		for _, name := range skillNames {
			src := filepath.Join(toolboxHome, "skills", name)
			dst := filepath.Join(tgt.dir, name)
			ok, warn := linkSkill(src, dst)
			if warn != "" {
				fmt.Fprintln(os.Stderr, "warn:", warn)
			}
			if ok {
				linked[tgt.label]++
			}
		}
	}

	exitCode := reportAgentPresence(claude, grok)

	fmt.Printf("toolbox home: %s\n", toolboxHome)
	for _, tgt := range targets {
		fmt.Printf("%s: %d/%d skills linked at %s\n", tgt.label, linked[tgt.label], len(skillNames), tgt.dir)
	}
	fmt.Printf("Claude Code detected: %s\n", presence(claude))
	fmt.Printf("Grok Build detected: %s\n", presence(grok))

	return exitCode
}

// listSkills returns the names of directories under
// $TOOLBOX_HOME/skills that contain a SKILL.md, sorted for deterministic
// output.
func listSkills(toolboxHome string) ([]string, error) {
	dir := filepath.Join(toolboxHome, "skills")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("tb: %w", err)
	}

	var names []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if _, err := os.Stat(filepath.Join(dir, e.Name(), "SKILL.md")); err == nil {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	return names, nil
}

// linkSkill implements the symlink strategy from docs/design.md §5.1
// step 3 / §6: idempotent, and never destructive to a file tb doesn't
// manage. ok reports whether dst ends this call correctly linked to src;
// warn is non-empty only when a real file or directory blocked linking
// (that skill is skipped, the rest still proceed).
func linkSkill(src, dst string) (ok bool, warn string) {
	info, err := os.Lstat(dst)
	if err != nil {
		if symErr := os.Symlink(src, dst); symErr != nil {
			return false, fmt.Sprintf("%s: %v", dst, symErr)
		}
		return true, ""
	}

	if info.Mode()&os.ModeSymlink == 0 {
		return false, fmt.Sprintf("%s exists and is not managed by tb, skipping", dst)
	}

	if current, readErr := os.Readlink(dst); readErr == nil && current == src {
		return true, "" // already correct: idempotent no-op
	}

	if rmErr := os.Remove(dst); rmErr != nil {
		return false, fmt.Sprintf("%s: %v", dst, rmErr)
	}
	if symErr := os.Symlink(src, dst); symErr != nil {
		return false, fmt.Sprintf("%s: %v", dst, symErr)
	}
	return true, ""
}

// reportAgentPresence prints the install-time agent warning/error
// (docs/design.md §5.1 step 4, §7) and returns the exit code.
func reportAgentPresence(claude, grok bool) int {
	switch {
	case claude && grok:
		return 0
	case claude || grok:
		missing := "Grok Build"
		if grok {
			missing = "Claude Code"
		}
		fmt.Fprintf(os.Stderr, "warn: %s not detected on this machine\n", missing)
		return 0
	default:
		fmt.Fprintln(os.Stderr, "error: neither Claude Code nor Grok Build detected")
		return 1
	}
}

func presence(ok bool) string {
	if ok {
		return "present"
	}
	return "absent"
}
