# 0011. `trainer-maker` and `cpp-cmake-verify` are new skills; this increment is v1.6

Status: accepted
Date: 2026-09-27

## Context

`cpp-trials` is a rustlings-style C++ exercise trainer built entirely under
the `cpp-systems` profile. None of the design decisions that actually make it
a *trainer* — a learner workspace with earned solutions, the `watch`-loop
shape, tiered hints, a manifest-driven exercise set, reproducible-command
discipline, an embedded/bundled distribution model, layered docs — came from
a toolbox skill. All of it was invented from scratch across 27 ADRs inside
`cpp-trials` itself. The next trainer (any language, any subject domain — a
Rust distributed-database trainer was the motivating example) would have to
re-derive the same ~16 decisions with no toolbox help.

Separately, a real defect surfaced while tracing this: `skills/cpp-verify/
SKILL.md`'s worked commands are llama.cpp-shaped
(`-DLLAMA_FATAL_WARNINGS=ON`, `./build/bin/test-backend-ops`, `ci/run.sh`,
`-DLLAMA_SANITIZE_ADDRESS=ON`). `cpp-trials`' `CMakeLists.txt` defines none of
these (confirmed by `rg 'SANITIZE|LLAMA_FATAL' CMakeLists.txt` — zero
matches), yet `cpp-trials` has been pointed at `cpp-verify` as its verify
recipe since `tb init -p cpp-systems`. `cpp-verify`'s own text says the
repo's CI is the source of truth and to confirm flags rather than recall
them, so this is a mismatch the skill is honest about, not silently wrong —
but a plain CMake C++ repo with no ggml/backend matrix has had no skill that
actually fits it.

## Options

**Trainer-pattern fork:**

- **Bake the pattern into `cpp-systems`/`cpp-ggml` directly** — rejected.
  `cpp-ggml` is locked to ggml-family upstream contribution (its own
  description already skips non-ggml C/C++), and doing this would re-tie a
  pattern that is explicitly language- and domain-agnostic to C++, which
  defeats the point of extracting it.
- **A new full profile (e.g. `cpp-trainer`)** — rejected as speculative
  infrastructure for a second, currently hypothetical adopter (a Rust
  trainer does not exist yet). It is also unnecessary: `tb install` symlinks
  every skill in `skills/` into `~/.claude/skills` regardless of profile
  membership, and skills trigger off their own `description`, independent of
  any profile's `skills` array. A profile is not required to make a skill
  usable.
- **A standalone, language-agnostic skill (`trainer-maker`)** — chosen.
  Triggered by its own description whenever the task is building or
  reviewing a rustlings-style exercise trainer, in any language or subject
  domain. Meant to be paired with whatever verify/house-style skill the
  repo's actual profile already provides for its toolchain.

**`cpp-verify` mismatch fork:**

- **Patch `cpp-verify` with an if-not-ggml conditional branch** — rejected.
  Same shape ADR 0008 rejected for `rust-systems`/Bevy and ADR 0010 rejected
  for `cpp-systems`/`c-cli`: a shared skill does not grow a per-caller
  branch.
- **Reuse `cpp-verify` unmodified** — rejected. The mismatch is real, not
  cosmetic; a repo following it verbatim would compile with a CMake option
  that does not exist and be told to run `ci/run.sh`, which is not present.
- **A new sibling skill, `cpp-cmake-verify`** — chosen. A generic recipe for
  a plain CMake C++ project with no ggml/backend matrix. `cpp-verify`'s
  description gains one mutual-skip clause pointing at it; its body and
  recipe are untouched.

## Decision

Add `skills/trainer-maker/SKILL.md` and `skills/cpp-cmake-verify/SKILL.md`.
No new profile. No change to `cmd/tb` — nothing in the binary renders a
skill differently based on which one it is; skills are pure content,
symlinked wholesale by `tb install`. This increment is v1.6 in
`docs/requirements.md` / `docs/design.md`'s spec history; `tbVersion` (the
constant `tb init` stamps into a product's `toolbox.toml`) is **not**
bumped, since no profile or `tb` behavior changed. This is the first
increment where the spec version and `tbVersion` diverge, because it is the
first increment that adds skills without adding or changing a profile.

## Consequences

- `docs/requirements.md`: **Delta from v1.5** (§25) and **Acceptance
  (v1.6)** (§26); header becomes v1.6. `docs/design.md` gains §24–§25.
- `skills/trainer-maker/SKILL.md`: eight house-style rules (manifest and
  anti-spoiler workspace; a passing check is not "move on"; `watch` as the
  default command with a fixed loop shape; a minimal runner-owned check
  protocol; reproducible commands and one color policy; small atomic state
  independent of any subprocess; a standard CLI surface; bundled content
  with maintainer tooling and layered docs kept off the learner path).
  Frontmatter description names no language, toolchain, or subject domain.
- `skills/cpp-cmake-verify/SKILL.md`: `cmake -B build`; `cmake --build build
  -j`; `ctest --test-dir build -L main --output-on-failure`; clang-format
  on added lines; a sanitizer build (confirmed against the repo's own CMake
  options, not assumed) for memory-touching diffs; `git status` clean of
  build output. No `ci/run.sh`, no `test-backend-ops`, no backend-parity
  language.
- `skills/cpp-verify/SKILL.md`: one clause added to its `description`
  pointing at `cpp-cmake-verify` for a plain CMake C++ repo with no
  ggml/backend matrix. No other change.
- `skills/review/SKILL.md`'s house-style pass item names `trainer-maker`
  (paired with the repo's own verify/house-style skill, when building a
  trainer) and `cpp-cmake-verify` (the non-ggml CMake alternative) alongside
  the existing per-profile domain skills.
- `README.md` gains a new subsection documenting skills that exist outside
  any profile's fixed `skills` array — the first time this has happened;
  every prior skill belonged to exactly one profile's list.
- No `profiles/*.toml` file changes. `cmd/tb/**` is untouched: no
  `tbVersion` bump, no `init` usage-string edit, no `doctor` toolchain
  branch, no new Go tests. `TestRealProfilesParse` and every existing
  profile fixture are unaffected.
- Not in v1.6: a seventh profile, a generic C/C++ house-style skill beyond
  what these two skills need, any change to `cpp-systems.toml` or
  `cpp-ggml`'s body, any `tb` code change. Retrofitting `cpp-trials` to use
  these skills is a separate follow-up in that repo, done after this ADR is
  accepted and `tb install` has re-linked skills machine-wide — not part of
  this increment's scope.
