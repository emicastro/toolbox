---
name: c-cli
description: Use when writing or reviewing C in a program the user owns and runs — terminal tools, TUIs, exercise runners — including Makefiles, headers, and any change to memory, buffers, or error handling; skip C++ and ggml-family work (use the `cpp-ggml` skill) and non-C files.
---

# c-cli

House style for the `c-cli` profile: C programs the user owns and runs.
Six rules. Do not pick a dependency (ncurses or any other library), a C
standard later than C11, a test framework, or a debugger here — that is
an ADR in the product repo. Changing any of the six at repo scope is also
an ADR.

## The Makefile is the build

Every repo has a Makefile with at least `all`, `test`, and `clean`.
`make` builds; `make test` runs the tests and exits non-zero on failure.
Compiler flags are set once, in `CFLAGS`, not repeated per target.

Build commands do not live in README prose or in C string literals — both
drift, and nothing tells you when they have. A repo with no Makefile
builds through the `c-verify` fallback, and adding the Makefile is its
first task.

## Warnings are errors

`-std=c11 -Wall -Wextra -Werror`, plus `-g`. A diff that introduces a
warning fails.

Fix the cause with the narrowest change: `(void)param;` for a parameter
that is intentionally unused, the correct format specifier (`%zu` for
`size_t`, `%d` for `int`) rather than a cast that hides the mismatch. No
`#pragma GCC diagnostic` to make a build pass.

## Own every allocation and every buffer

Prefer caller-owned structs and fixed capacities behind a `#define` to
the heap. A program with no `malloc` has no leaks, no double frees, and
no use-after-free — that style is the default, not an accident to be
"improved".

When you do allocate:

- Check every `malloc` / `calloc` / `realloc` return.
- Route `realloc` through a temporary, so a failure does not leak the
  original: `tmp = realloc(p, n); if (!tmp) { ... } p = tmp;`.
- Give every allocation one owner and one `free` on every path out.

Bound every buffer write: `snprintf(buf, sizeof buf, ...)`. Never
`strcpy`, `strcat`, or `sprintf` into a fixed array. A diff touching any
of this gets the sanitizer build — see `c-verify`.

## One error convention, and exit codes that mean something

One convention per repo: functions return a status (`bool`, `0`/`-1`, or
an enum) and the caller decides. Only `main`, or a clearly top-level
function, prints the error and exits. A module below that never calls
`exit()`.

Errors go to `stderr` (`perror` or `fprintf(stderr, ...)`), never
`stdout`. Exit codes: `0` success, `1` failure, `2` usage error — the same
contract `tb` itself uses.

An ncurses program calls `endwin()` before printing any error and before
exiting, so the user gets their terminal back. The same applies after
shelling out: restore whatever state the program changed (terminal mode,
working directory) before carrying on.

## Headers declare, sources define

`#ifndef NAME_H` / `#define NAME_H` / `#endif` guards; do not mix in
`#pragma once`. Headers hold declarations, types, and `#define`
capacities; function bodies live in `.c` files.

Every file-local function and variable is `static`. Each `.c` includes
its own header first, so a header that is not self-contained fails to
compile. `snake_case` for functions and variables.

## Build artifacts are not source

Nothing the compiler produces is tracked: binaries, `*.o`, `*.a`,
`*.so`, sanitizer builds, `*.dSYM/`, coverage files. `.gitignore` names
every binary the Makefile produces.

A tracked binary is a finding. It is noise in every diff, and it is also
wrong for the next machine — a binary built on one architecture does not
run on another. Untrack it with `git rm --cached <file>` and add it to
`.gitignore`; do not rewrite history to remove it.
