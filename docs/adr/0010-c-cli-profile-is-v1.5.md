# 0010. `c-cli` is a sixth profile; this increment is v1.5

Status: proposed
Date: 2026-09-22

## Context

v1.4 added `cpp-systems` for C/C++ built with CMake. Its house-style
skill is `cpp-ggml`, written for contributing upstream to `llama.cpp`: do
not run `tb init` here, never write the PR description or the commit
message, the target repo's rules outrank the persona. Applied to a repo
the user owns, every one of those rules is either wrong or inert. Its
verify recipe is CMake plus `ctest -L main` with
`-DLLAMA_FATAL_WARNINGS=ON`, which no plain C program has.

The motivating repos are `~/src/algorithm-visualizer` (an ncurses
visualiser, nine `*_algorithms.c` modules plus a headless `algo_test.c`)
and `~/src/clings` (a rustlings-style exercise runner, single
`clings.c`). A survey of both establishes what this profile has to fix:

- **Neither repo has a build system.** algorithm-visualizer's only build
  instructions are `gcc` lines in `README.md`; clings' compile line
  exists only as `system("gcc -o main main.c")` inside `clings.c`.
- **The prose build lines have already drifted.** The README's
  application line compiles `utils.c`; its test line silently omits it.
- **Warning flags are weak or absent.** algorithm-visualizer uses
  `-std=c11 -Wall -Wextra`, with no `-Werror`, no `-g`, no sanitizers.
  clings passes no flags at all and contains a `%zu` format applied to an
  `int`, an unchecked `realloc`, and a `strcpy` into a `char[256]`.
- **Build artifacts are tracked.** `algo_test` is a 25 KB ELF in git;
  clings tracks a stale Mach-O arm64 binary — the wrong architecture for
  the machine it is checked out on — and has no `.gitignore`.
- **The house style is real and consistent.** No heap allocation anywhere
  in the visualiser, caller-owned structs, fixed-size arrays behind
  `#define` capacities, `#ifndef` header guards, file-locals `static`.

So the value of a C profile sits in the verify recipe and in artifact
hygiene, more than in prose about how to write C.

Two sub-forks are settled here rather than left to a diff. **The build:**
C has no `cargo` and no `go`, and these two repos have no build file at
all, so a recipe must either assume a build system or discover one.
**Doctor's toolchain probe:** `reportToolchain` does one `LookPath` per
profile, and a C profile could reasonably probe `gcc` or `make`.

ADR 0009's consequences list "no sixth profile" and "no generic C/C++
house style" as items that need a new ADR. This is that ADR. ADR 0005
reserved v2 for a new product surface; 0007 narrowed v2 to MCP, `tb
verify`, other agents, or a CLI-surface change. A content pack
`LoadProfileByName` already loads is a minor increment, as in 0007, 0008,
and 0009.

## Options

- **No profile; use `cpp-systems`** — one C/C++ profile. The `Skills:`
  line then advertises `cpp-ggml` in a repo the user owns, where its
  first two rules (the target repo's rules win; never speak for the
  contributor) are meaningless, and the verify recipe names CMake
  targets that do not exist. Rejected for the reason 0007 rejected a
  skill-only add.
- **Extend `cpp-systems` with a C section** — one profile, two recipes
  in `cpp-verify` and a C branch in the house style. Cheapest, but it
  makes the ggml lock conditional, which is exactly the shape 0008
  rejected when it declined to patch `rust-systems` with Bevy
  exceptions.
- **`c-systems` generic C profile, v1.5** — a sixth profile whose scope
  is "any C". A library and a terminal program fail review on different
  things, so the house style has to hedge on ownership, error
  conventions, and `main`, which is most of its content.
- **`c-cli` profile, v1.5** — a sixth profile scoped to C programs the
  user owns and runs. The axis is "runnable program", following
  `infra-go`. House style can be concrete: a Makefile, `-Werror`, owned
  allocations, exit codes, header conventions, no tracked artifacts.

## Decision

`c-cli` profile, version v1.5. Scope is C programs the user owns and
runs. ncurses is a dependency, not a profile axis, so the profile is not
named for it; a C library or embedded variant, if ever needed, is a
later profile and a later ADR.

**The build is a Makefile, with a direct-compiler fallback.** The
fallback is not a second-class path — it is the bootstrap, because a repo
with no build file must compile before it can acquire a Makefile, and
both motivating repos are in exactly that state. Writing the Makefile is
the first task the profile hands a repo, not a precondition for using
it. Discovery-only (the `cpp-verify` approach of deferring to the repo's
CI) is rejected here because there is nothing to discover: one repo keeps
its build line in prose that has already drifted, the other keeps it
inside a C string literal.

**Doctor probes `gcc`.** A missing compiler fails every path including
the fallback; a missing `make` only fails the preferred one. One probe
per profile stays the pattern.

## Consequences

- `docs/requirements.md` keeps all earlier text and adds a **Delta from
  v1.4** (§23) and **Acceptance (v1.5)** (§24); header becomes v1.5.
  `docs/design.md` gains §22–§23.
- `profiles/c-cli.toml`: process skills (`spec`, `adr`, `onboard`,
  `scout`, `handoff`, `review`), domain skills `c-verify`, `c-cli`. No
  `aws-guard`, no `cpp-*`, no `rust-*`, no `*-go`. Templates are the same
  five as every other profile.
- `skills/c-verify/SKILL.md`: `make` then `make test`, the fallback
  compiler line with `-std=c11 -Wall -Wextra -Werror -g`, a
  `-fsanitize=address,undefined` build required for any diff touching
  allocation, buffers, or pointer arithmetic, and `git status` clean of
  build output before a task is ticked.
- `skills/c-cli/SKILL.md`: six rules — the Makefile is the build;
  warnings are errors; own every allocation and every buffer; one error
  convention and exit codes that mean something; headers declare, sources
  define; build artifacts are not source.
- `const tbVersion` becomes `"1.5.0"`. Init usage lists six profile
  names. `reportToolchain` gains a `c-cli` branch on `gcc`.
  `TestRealProfilesParse` gains `c-cli`; `helpers_test.go` gains the
  fixture; an init test and a doctor test mirror the `cpp-systems` pair.
- `README.md` lists six profiles and the C verify block.
- `review` skill's profile list gains `c-cli`.
- `cpp-systems` is unchanged and remains the C++/CMake profile. Mutual
  skip in the skill descriptions: `c-cli` skips C++ and ggml-family work,
  `cpp-ggml` already skips non-ggml C/C++.
- Fixing the motivating repos is out of scope here. Adding their
  Makefiles and `.gitignore`s, untracking `algo_test` and the stale
  Mach-O binary, and the defects the survey found in `clings.c` are work
  in those repos, after `tb init -p c-cli`.
- Not in v1.5: a seventh profile, a C library or embedded variant, an
  ncurses or any other dependency pin, `valgrind` or `clang-tidy` as
  required tools, a `tb verify` wrapper, and any change to the five
  existing profiles' content.
