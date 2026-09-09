package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// isInteractive reports whether stdin looks like a real terminal, using
// only the standard library — docs/design.md §4's original text named
// golang.org/x/term for this, which docs/adr/0001-tb-dependency-policy.md
// rules out; design.md now reflects the stdlib-only check implemented
// here.
//
// The check is os.ModeCharDevice on the result of Stat, with one
// correction: /dev/null is also a character device, so a plain
// ModeCharDevice test misreads `tb init </dev/null` — a common way for
// scripts and CI to mark a command non-interactive — as interactive.
// Left uncorrected, that would print an interactive prompt and then
// block reading stdin instead of refusing immediately. os.SameFile
// against os.DevNull rules that specific, common case out; this is not a
// full TTY check (a redirect from some other character device would
// still read as interactive), which is an accepted v1 gap, not a claim
// of full terminal detection.
func isInteractive() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	if info.Mode()&os.ModeCharDevice == 0 {
		return false
	}
	if nullInfo, err := os.Stat(os.DevNull); err == nil && os.SameFile(info, nullInfo) {
		return false
	}
	return true
}

// confirmMerge implements the interactive/non-interactive/--force
// decision from docs/design.md §4's table: --force always proceeds
// without asking; a non-interactive session always refuses without
// asking; an interactive session is asked once, merge or refuse.
func confirmMerge(force bool, prompt string) bool {
	if force {
		return true
	}
	if !isInteractive() {
		return false
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprintf(os.Stderr, "%s [m]erge / [r]efuse: ", prompt)
		line, err := reader.ReadString('\n')
		if err != nil {
			return false
		}
		switch strings.ToLower(strings.TrimSpace(line)) {
		case "m", "merge":
			return true
		case "r", "refuse":
			return false
		}
	}
}
