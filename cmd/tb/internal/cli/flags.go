package cli

import (
	"flag"
	"io"
)

// bindProfileFlag registers both spellings of the global --profile / -p
// flag (docs/design.md §5) onto the same variable. If both are given on
// one command line, the one parsed last wins — an unspecified edge case
// design.md leaves to the implementer.
func bindProfileFlag(fs *flag.FlagSet, dst *string) {
	fs.StringVar(dst, "profile", "", "toolbox profile name")
	fs.StringVar(dst, "p", "", "toolbox profile name (shorthand)")
}

// newSilentFlagSet returns a FlagSet that does not print its own usage
// or errors: every caller in this package prints its own usage message
// on a parse failure instead, matching the exact wording design.md gives
// for each command.
func newSilentFlagSet(name string) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	return fs
}

func parseInitFlags(args []string) (profile string, force bool, err error) {
	fs := newSilentFlagSet("init")
	bindProfileFlag(fs, &profile)
	fs.BoolVar(&force, "force", false, "overwrite AGENTS.md and toolbox.toml without prompting")
	err = fs.Parse(args)
	return profile, force, err
}

func parseDoctorFlags(args []string) (profile string, err error) {
	fs := newSilentFlagSet("doctor")
	bindProfileFlag(fs, &profile)
	err = fs.Parse(args)
	return profile, err
}

func parsePersonaShowFlags(args []string) (profile string, err error) {
	fs := newSilentFlagSet("persona show")
	bindProfileFlag(fs, &profile)
	err = fs.Parse(args)
	return profile, err
}
