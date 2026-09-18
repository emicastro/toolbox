# 0009. `cpp-systems` is a fifth profile; this increment is v1.4

Status: accepted
Date: 2026-09-18

## Context

The four v1.3 profiles cover Rust and Go. A contribution to
`ggml-org/llama.cpp` (C/C++, CMake, ctest, a multi-backend matrix) has no
profile: `rust-verify` / `go-verify` name the wrong commands, and no
house-style skill names the failure modes of that codebase.

Two things make this different from every earlier profile, and both are
forks in their own right.

**The target repo is upstream, not the user's.** Every prior profile
assumed a repo where `tb init` may write `toolbox.toml`, `AGENTS.md`, and
`docs/` templates. `llama.cpp` already ships a tracked `AGENTS.md` and a
`CLAUDE.md` that points at it. `render.Merge` appends a managed region to
an existing file, so `tb init` there would modify a tracked upstream file
and inject toolbox binding rules on top of the project's own agent rules.
`scaffold.MaterializeTemplates` would additionally create
`docs/requirements.md`, `docs/design.md`, `docs/tasks.md`, and
`docs/session.md` inside the project's own `docs/`. `.git/info/exclude`
does not help for a file that is already tracked.

**The repo's agent rules are stricter than the persona and are
non-negotiable.** `llama.cpp/AGENTS.md` forbids an agent from writing PR
descriptions, commit messages, review comments, or replies; forbids
`git push`, `gh pr create`, `gh pr comment`; requires `Assisted-by:`
rather than `Co-authored-by:` when the user does ask for a commit;
requires AI-usage disclosure in the PR template; and tells fully
autonomous agents not to contribute at all. `CONTRIBUTING.md` adds that a
bug-fix PR needs a reproducible issue and a regression test that fails
before and passes after, that a new contributor keeps one open PR, and
that a first PR for a feature is CPU-only.

A generic C/C++ house-style pack (RAII, headers, ownership) would not
name any of that, nor the ggml specifics a review there actually fails
on: row-major tensor conventions, the transposed `ggml_mul_mat`,
`test-backend-ops` parity across two backends, the cost of a new
`ggml_type`. ADR 0008 already decided that case for Bevy: lock the
family when the generic pack cannot name the real failures.

ADR 0005 reserved v2 for a new profile surface; ADR 0007 narrowed v2 to
MCP, `tb verify`, other agents, or a CLI-surface change, and treated a
content pack `LoadProfileByName` already loads as a minor increment. A
fifth profile is the same shape as 0007 and 0008: content, no new
subcommand, no schema change.

## Options

- **No profile; use `rust-systems` or nothing** — the verify recipe is
  wrong (`cargo` against a CMake tree) and the `Skills:` line advertises
  Rust house style in a C++ repo. Rejected for the reason 0007 rejected
  a skill-only add.
- **Generic `cpp-systems` profile with a generic C/C++ house style** —
  portable across C++ repos, but cannot name tensor layout, backend
  parity, the AI policy, or the upstream-repo constraint, which is most
  of what this work needs.
- **`cpp-systems` profile whose house-style skill is ggml-locked
  (`cpp-ggml`), version v1.4** — the profile names the language and build
  surface (C/C++ with CMake, so it is reusable in the user's own C++
  repos); the house-style skill locks the family (ggml / llama.cpp /
  whisper.cpp) the way `game-bevy` locks the engine. `cpp-verify` is the
  CMake recipe and is shareable with any later C/C++ house-style skill.
- **`cpp-ggml` as the profile name** — honest about the lock, but burns
  the C/C++ profile slot on one upstream family and leaves nowhere for
  the user's own C++ repos to point.

## Decision

`cpp-systems` profile with skills `cpp-verify` and `cpp-ggml`, version
v1.4. The profile is the C/C++ + CMake surface; the lock lives in the
house-style skill, so a later generic `cpp-core` house style can join the
same profile without a sixth profile or a rename.

**`tb init` is not run in `llama.cpp`.** The profile is wired only in a
repo the user owns. For upstream work the two skills carry the whole
value: `tb install` symlinks them into `~/.claude/skills`, the agent
loads them by description, and the upstream tree stays byte-clean. This
is recorded as a rule in `cpp-ggml`, not enforced by `tb` — v1.4 adds no
guard, no `--no-write` flag, and no upstream detection.

The repo's own `AGENTS.md` and `CONTRIBUTING.md` outrank both the persona
and the `cpp-ggml` skill wherever they disagree. `cpp-ggml` says so in
its first rule rather than restating their content as toolbox policy.

## Consequences

- `docs/requirements.md` keeps all earlier text and adds a **Delta from
  v1.3** section; header becomes v1.4. `docs/design.md` gains §19.
- `profiles/cpp-systems.toml`: process skills (`spec`, `adr`, `onboard`,
  `scout`, `handoff`, `review`), domain skills `cpp-verify`, `cpp-ggml`.
  No `aws-guard`, no `rust-*`, no `*-go`. Templates are the same five as
  the other profiles.
- `skills/cpp-verify/SKILL.md` is the CMake recipe: configure with
  `-DLLAMA_FATAL_WARNINGS=ON`, build, `ctest -L main`, `clang-format` on
  added lines only, `test-backend-ops` when `ggml/` changed, sanitizer
  build for memory-touching diffs, `ci/run.sh` before publishing. Unlike
  `cargo`, there is no universal C++ recipe: the skill states that the
  repo's CI workflows are the source of truth and the commands are
  confirmed against them, not recalled.
- `skills/cpp-ggml/SKILL.md` is the house style: the repo's rules win;
  never speak for the contributor; understanding is the deliverable;
  blend in; comments last and short; ggml tensor facts.
- `const tbVersion` becomes `"1.4.0"`. Init usage lists five profile
  names. `reportToolchain` gains a `cpp-systems` branch warning on a
  missing `cmake` rather than `cargo` or `go`. `TestRealProfilesParse`
  gains `cpp-systems`; `helpers_test.go` gains the fixture; an init test
  and a doctor test mirror the `game-bevy` pair.
- `README.md` lists five profiles and the CMake verify block.
- `review` skill's profile list gains `cpp-systems`.
- Not in v1.4: a `tb` guard against initialising an upstream repo, a
  generic C/C++ house style, a CUDA/Metal/Vulkan sub-profile, a Python
  (`gguf-py`, `convert_*.py`) recipe, and any wrapper around `ci/run.sh`.
