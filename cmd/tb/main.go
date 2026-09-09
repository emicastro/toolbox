// Command tb is the toolbox CLI: a static Go binary that wires a product
// repo to a toolbox profile and keeps skills symlinked into the agents
// installed on this machine (docs/requirements.md, docs/design.md).
package main

import (
	"fmt"
	"os"

	"toolbox/cmd/tb/internal/cli"
	"toolbox/cmd/tb/internal/paths"
)

func main() {
	os.Exit(run(os.Args[1:]))
}

// run resolves $TOOLBOX_HOME exactly once (docs/design.md §1) and
// dispatches to the requested subcommand. Every subcommand except
// `doctor` treats a resolution failure as fatal before it does anything
// else; `doctor` is handed the error instead, since it must still print
// the rest of its one-screen report (docs/design.md §1, §5.3).
func run(args []string) int {
	if len(args) == 0 {
		usage()
		return 2
	}

	cmd, rest := args[0], args[1:]
	home, homeErr := paths.Resolve()

	switch cmd {
	case "install":
		if homeErr != nil {
			return fatal(homeErr)
		}
		return cli.Install(home, rest)

	case "init":
		if homeErr != nil {
			return fatal(homeErr)
		}
		return cli.Init(home, rest)

	case "doctor":
		return cli.Doctor(home, homeErr, rest)

	case "skill":
		if len(rest) == 0 || rest[0] != "new" {
			usage()
			return 2
		}
		if homeErr != nil {
			return fatal(homeErr)
		}
		return cli.SkillNew(home, rest[1:])

	case "persona":
		if len(rest) == 0 || rest[0] != "show" {
			usage()
			return 2
		}
		if homeErr != nil {
			return fatal(homeErr)
		}
		return cli.PersonaShow(home, rest[1:])

	default:
		usage()
		return 2
	}
}

func fatal(err error) int {
	fmt.Fprintln(os.Stderr, err)
	return 1
}

func usage() {
	fmt.Fprintln(os.Stderr, `usage: tb <command> [flags]

commands:
  install                 link skills into Claude Code and Grok Build
  init -p|--profile NAME   wire the current repo to a toolbox profile
  doctor                   report toolbox health (read-only)
  skill new NAME           scaffold a new skill under $TOOLBOX_HOME/skills
  persona show             print the active profile's persona sections

$TOOLBOX_HOME (default ~/toolbox) must exist and be a git clone of the
toolbox repo. See docs/requirements.md and docs/design.md in that repo.`)
}
