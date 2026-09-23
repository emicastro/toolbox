---
name: c-verify
description: Use before claiming any task done in a `c-cli` repo, and after every change to C sources, headers, the Makefile, or CI — this is the verify recipe (`make`; `make test`; fallback `gcc -std=c11 -Wall -Wextra -Werror -g` when there is no Makefile yet; a `-fsanitize=address,undefined` build for anything touching memory; `git status` clean of build output); skip only for edits that touch none of those.
---

# c-verify

The C verify recipe (`c-cli` profile). Run it in your own shell; success
is the command output, not an assertion. There is no `tb verify` wrapper
and never will be.

## Commands

```sh
make
make test
```

`make test` must exit non-zero when a test fails. A test binary that
prints "ok" and returns 0 no matter what is not a test — if you find one,
that is a finding.

A build is not a test. Run both.

## Fallback when there is no Makefile yet

Until the repo has a Makefile, compile directly with the profile flags:

```sh
gcc -std=c11 -Wall -Wextra -Werror -g -o <bin> <sources> $(pkg-config --cflags --libs <lib>)
```

List `<sources>` from the tree (`ls *.c`), not from README prose — prose
build lines drift, and a source file silently missing from the list is
exactly how they drift. Add the `pkg-config` part only for libraries the
program actually links (for example `ncursesw`).

Needing the fallback means the Makefile task is still open. Say so in
your report; do not treat the fallback as the finished state.

## Sanitizers

Any diff that touches allocation, a buffer, a string copy, or pointer
arithmetic also gets a sanitizer build and a test run against it:

```sh
gcc -std=c11 -Wall -Wextra -Werror -g -fsanitize=address,undefined \
    -o <bin>-asan <sources> $(pkg-config --cflags --libs <lib>)
./<bin>-asan            # or the test binary, built the same way
```

Use `make asan` instead when the Makefile provides it. A sanitizer report
is a failure even if the program's own output looks right.

## Artifacts

After the recipe, check the tree:

```sh
git status --short
```

It must show no build output — no binaries, no `*.o`, no `<bin>-asan`. If
one appears, `.gitignore` is incomplete: that is a finding, and the fix is
a `.gitignore` entry, not deleting the file before you report.

## Rules

- Run the recipe before ticking a checkbox in `docs/tasks.md` or
  reporting a task complete. A clean compile is not verification.
- `-Werror` means every warning is a failure. Fix the cause; do not add a
  `#pragma GCC diagnostic` or drop the flag to make the build pass.
- Paste or summarize the actual output. Persona rule: evidence over
  narration.
- A test that was already failing before your change is still a finding —
  report it, do not fold a fix into your diff without a task for it.
