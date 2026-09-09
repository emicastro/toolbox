// Package cli implements tb's five subcommands (docs/design.md §5):
// install, init, doctor, skill new, and persona show. Each function
// takes the already-resolved $TOOLBOX_HOME (docs/design.md §1, resolved
// once in main) plus its own remaining argv, and returns the process
// exit code — 0 success, 1 operational failure, 2 usage error.
package cli
