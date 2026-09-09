// Package cli implements tb's subcommands (docs/design.md §5).
//
// The five functions below are placeholders: task group 3
// (docs/tasks.md) only wires up the packages each command will need
// (paths, config, agents, render) and the main.go dispatch that calls
// into this package. Task group 4 replaces each body with the real
// command logic described in docs/design.md §5.1-§5.5.
package cli

import (
	"fmt"
	"os"
)

// Install implements `tb install` (docs/design.md §5.1).
func Install(toolboxHome string, args []string) int {
	return notImplemented("install")
}

// Init implements `tb init` (docs/design.md §5.2).
func Init(toolboxHome string, args []string) int {
	return notImplemented("init")
}

// Doctor implements `tb doctor` (docs/design.md §5.3). homeErr is the
// error from resolving $TOOLBOX_HOME, if any: doctor is the one command
// that must still run and print a one-screen report when the toolbox
// home is missing, rendering that failure as its first report line
// instead of exiting before it prints anything (docs/design.md §1, §5.3).
func Doctor(toolboxHome string, homeErr error, args []string) int {
	return notImplemented("doctor")
}

// SkillNew implements `tb skill new <name>` (docs/design.md §5.4).
func SkillNew(toolboxHome string, args []string) int {
	return notImplemented("skill new")
}

// PersonaShow implements `tb persona show` (docs/design.md §5.5).
func PersonaShow(toolboxHome string, args []string) int {
	return notImplemented("persona show")
}

func notImplemented(cmd string) int {
	fmt.Fprintf(os.Stderr, "tb %s: not implemented yet (docs/tasks.md group 4)\n", cmd)
	return 1
}
