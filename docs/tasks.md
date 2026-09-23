# Toolbox — Tasks (v1.5)

Status: v1 groups 1–6 accepted 2026-09-09; v1.1 groups 7–11 accepted 2026-09-11; v1.2 groups 12–16 accepted 2026-09-12; v1.3 groups 17–21 accepted 2026-09-12; v1.4 groups 22–26 accepted 2026-09-18; v1.5 groups 27–31 accepted 2026-09-22
Date: 2026-09-22
Each task is scoped for one Implement session and ends in a runnable check.
Follow `docs/design.md` for shape; do not reopen a decision recorded there
or in `docs/adr/*` without a new ADR (persona rule, §7 of requirements).

## Group 1 — Repo skeleton

- [x] **1.1** `git init` this repo; add `.gitignore` (`dist/`, `*.local.toml`
      if that pattern is ever used, standard Go entries: none needed yet
      since `cmd/tb` has no build artifacts checked in).
      Check: `git status` shows a clean initial tree after first commit.
- [x] **1.2** Write toolbox's own `AGENTS.md` (hand-written, not
      `tb`-generated — this repo has no profile of its own). Short: point at
      `docs/requirements.md`, `docs/design.md`, `docs/tasks.md` as source of
      truth; state English-only content rule.
      Check: file exists, under ~30 lines.
- [x] **1.3** Write toolbox-side `toolbox.toml` per design §3.1 note ("if
      needed") — a single `schema_version = "1"` line is sufficient for v1;
      do not invent fields design.md doesn't define.
      Check: file parses under the task-group-3 TOML reader once it exists
      (revisit after 3.3 if needed).

## Group 2 — Content (skills, profiles, persona, templates)

- [x] **2.1** `personas/default.md` exactly per design §9.
      Check: two `##` headings present, titles match design.md verbatim.
- [x] **2.2** `profiles/rust-systems.toml` and `profiles/infra-go.toml` per
      design §3.3, with the exact skill lists from design §8 (five process +
      two or three domain skills, no extras).
      Check: manually diff each `skills` array against requirements §6.2/§6.3.
- [x] **2.3** Five process skills, one `SKILL.md` each per design §10:
      `skills/spec/`, `skills/adr/`, `skills/onboard/`, `skills/scout/`,
      `skills/handoff/`. Bodies describe the workflow named in requirements
      §6.2 (one sentence per skill there is the minimum content bar).
      Check: each has `name`/`description` frontmatter matching its dir name.
- [x] **2.4** Five domain skills: `skills/rust-verify/`,
      `skills/rust-systems/`, `skills/go-verify/`, `skills/infra-go/`,
      `skills/aws-guard/`. The two `*-verify` skills must state the exact
      commands from requirements §6.3 (`cargo test`; `cargo clippy -- -D
      warnings`; Miri when `unsafe` in the change / `go test ./...`; `go
      vet`; race when tests involve concurrency).
      Check: verify commands in the skill body match requirements §6.3
      word-for-word.
- [x] **2.5** `templates/AGENTS.md` (the managed-region skeleton, design §4)
      and `templates/docs/{requirements,design,tasks,session}.md` as minimal
      starter stubs (headings + one-line "fill this in" placeholders) plus
      `templates/docs/adr/0000-template.md` (copy of this repo's own
      `docs/adr/0000-template.md`).
      Check: `templates/AGENTS.md` byte-matches design §4's template block.

## Group 3 — Binary core

- [x] **3.1** `cmd/tb/go.mod` (`module toolbox/cmd/tb`, `go 1.27`, no
      `require` block per ADR 0001). `cmd/tb/main.go` with subcommand
      dispatch (`install`, `init`, `doctor`, `skill`, `persona`) and a
      top-level usage message on no/unknown args (exit 2).
      Check: `go build ./cmd/tb/...` succeeds with zero dependencies in
      `go.sum` (no `go.sum` file at all is fine).
- [x] **3.2** `internal/paths`: `Resolve() (string, error)` per design §1.
      Check: table tests — env set to existing dir, env unset with
      `~/toolbox` present, env set to missing dir (error), env empty string
      (falls through to default).
- [x] **3.3** `internal/config`: TOML-subset reader per design §3.1, plus
      typed loaders `LoadProfile(path) (Profile, error)` and
      `LoadPointer(path) (Pointer, error)` for the two schemas in §3.2/§3.3.
      Check: table tests covering every construct in the subset and at
      least three unsupported-syntax rejections with correct `file:line`.
- [x] **3.4** `internal/agents`: `Detect() (claude, grok bool)` per ADR
      0003 (PATH lookup OR config-dir check, injectable for tests).
      Check: tests fake `PATH` and `$HOME`/`$USERPROFILE`-equivalent to
      cover all four presence combinations without touching the real
      environment.
- [x] **3.5** `internal/render`: managed-region template execution and the
      three-state merge logic from design §4 (`Render(profile, persona)
      string`, `Merge(existing, region string) (result string, refused
      bool)`).
      Check: golden-file test for `Render`; `Merge` tests for "no markers",
      "markers present", confirming content outside markers survives
      byte-for-byte.

## Group 4 — Commands

- [x] **4.1** `tb install` (design §5.1 + §6). Creates both skill dirs,
      links/relinks/skips per the symlink strategy, prints the summary,
      applies the exit-code policy from design §7.
      Check: integration test in `t.TempDir()` with a fake `$TOOLBOX_HOME`
      containing two skills and fake `$HOME`; asserts both dirs get correct
      symlinks and a pre-existing real file at one target is skipped with a
      warning, not fatal.
- [x] **4.2** `tb init` (design §5.2). Profile load, git-repo check,
      `toolbox.toml` write/merge, `AGENTS.md` render/merge via
      `internal/render`, template materialization (create-if-missing,
      `docs/adr/*` never touched beyond the create-if-missing template
      copy).
      Check: integration tests for: empty repo (all files created); repo
      with existing `docs/adr/0001-something.md` (untouched, byte-identical
      before/after); re-run without `--force` on existing `AGENTS.md`
      (non-interactive refuses, exit 1, file unchanged); `--force` (only
      `AGENTS.md` + `toolbox.toml` change).
- [x] **4.3** `tb doctor` (design §5.3). Read-only, one-screen output, exit
      code only on zero-agents.
      Check: tests for "clean" (both agents faked present), "one missing"
      (exit 0), "zero agents" (exit 1); a test with no `toolbox.toml` in
      cwd skips steps 5–6 without erroring.
- [x] **4.4** `tb skill new <name>` (design §5.4).
      Check: test for valid name (scaffold created, next-steps text
      printed), invalid name (exit 2), already-exists (exit 1).
- [x] **4.5** `tb persona show` (design §5.5).
      Check: test reading a fixture profile + `personas/default.md`,
      asserting both section titles print exactly.

## Group 5 — Build and CI

- [x] **5.1** Build script (`Makefile` or `cmd/tb/build.sh` — pick one, keep
      it under ~20 lines) producing `dist/tb-linux-amd64` and
      `dist/tb-darwin-arm64` via the flags in design §11.
      Check: both binaries build locally on Arch (darwin/arm64 as a
      cross-compile check, `file dist/tb-darwin-arm64` reports the right
      target).
- [x] **5.2** CI workflow: `gofmt -l`, `go vet ./...`, `go test ./...` on
      push, scoped to `cmd/tb/**` paths only.
      Check: workflow passes on this repo's initial commit containing all
      of the above.

## Group 6 — Acceptance (requirements §13, run manually)

- [x] **6.1** On Arch: clone to `TOOLBOX_HOME` (or default), build/install
      `tb`, run `tb install` — confirm skills link, `tb doctor` is clean
      aside from any genuinely-missing agent.
- [x] **6.2** `tb init -p rust-systems` in a fresh empty cargo repo and in a
      non-empty rust repo that already has `docs/adr/0001-*.md` — confirm
      pointer + `AGENTS.md` + `docs/` created, existing ADR untouched.
- [x] **6.3** `tb init -p infra-go` in a Go module — same checks.
- [x] **6.4** Re-run `tb init` without `--force` on each repo from 6.2/6.3 —
      confirm merge-or-refuse behavior (interactive prompt in a real TTY;
      refuse in a script).
- [x] **6.5** Re-run with `--force` — confirm only `AGENTS.md` and
      `toolbox.toml` changed (`git diff --stat`).
- [x] **6.6** `tb skill new demo` — confirm scaffold under
      `$TOOLBOX_HOME/skills/demo`.
- [x] **6.7** Open Claude Code and/or Grok Build (whichever is installed on
      that machine) and confirm the symlinked skills are visible.
- [x] **6.8** Repeat 6.1–6.7 on the MacBook Air M1.
- [x] **6.9** Confirm `tb doctor` fails closed (exit 1) if both agent config
      dirs are temporarily renamed away, and warns (exit 0) with one
      renamed back.

## Group 7 — Spec (v1.1)

- [x] **7.1** ADR `docs/adr/0005-increment-is-v1.1.md` as `proposed` (v1.1
      vs v2 vs silent patch). Check: file exists, four MADR sections, two
      or more named options.
- [x] **7.2** ADR `docs/adr/0006-binding-rules-in-managed-region.md` as
      `proposed` (titles-only vs short rules block vs full persona).
      Check: wording of the five rules is in Decision mechanics.
- [x] **7.3** `docs/requirements.md` header v1.1; v1 body unchanged; §15
      delta + §16 acceptance. Check: every §16 item is a line you can fail.
- [x] **7.4** `docs/design.md` §14–§15 cite 0005/0006 and map §16 →
      mechanism. Check: traceability table has one row per §16 item.
- [x] **7.5** User marks 0005 and 0006 `accepted` (and this file's v1.1
      groups accepted). Check: both ADR headers say `Status: accepted`.
      Do not start Group 8 until this box is ticked.

## Group 8 — Content (v1.1)

- [x] **8.1** `templates/AGENTS.md` rules block per ADR 0006. Check: all
      five rules present between the markers; no new `{{.Field}}`.
- [x] **8.2** `personas/default.md`: same five rules; Implement names
      `review` as the last gate. Check: both `##` headings still parse as
      Plan / Implement titles.
- [x] **8.3** `skills/review/SKILL.md` per design §14.2. Check: frontmatter
      `name: review`; skip only when this session produced no diff.
- [x] **8.4** Both `profiles/*.toml`: `review` after `handoff`; `[verify].summary`
      matches design §14.3. Check: `rg 'handoff", "review' profiles/`.
- [x] **8.5** Process + verify + house-style + `aws-guard` skill bodies
      per requirements §15.3–§15.5. Check: `rg 'or \`AGENTS.md\`' skills/`
      prints nothing; `go-verify` lists `go vet ./...` and `go test -race ./...`;
      `rust-verify` lists `cargo fmt --check`, `clippy --all-targets`, and
      crate-level Miri.

## Group 9 — Binary stamp

- [x] **9.1** `const tbVersion = "1.1.0"` in `cmd/tb/internal/cli/init.go`.
      Check: `rg 'tbVersion' cmd/tb` shows `1.1.0`.
- [x] **9.2** `TestRenderAgainstRealTemplate` asserts the rules block is
      in the rendered region. Check: `cd cmd/tb && go test ./internal/render/`
      fails if the five rules are removed from the template.
- [x] **9.3** `gofmt -l .` silent; `go vet ./...`; `go test ./...` from
      `cmd/tb`. Check: all three pass.

## Group 10 — Toolbox-own + README

- [x] **10.1** This repo's `AGENTS.md` carries the five binding rules
      (no managed region here). Check: all five lines present.
- [x] **10.2** `README.md` verify commands, `review` skill, version 1.1.0.
      Check: README command lists match the `*-verify` skills.

## Group 11 — Acceptance (requirements §16, run manually on Arch)

- [x] **11.1** `tb init --force` in a rust-systems fixture and an infra-go
      fixture: `AGENTS.md` contains the five rules and `review` on
      `Skills:`; prose outside markers unchanged; `toolbox.toml` has
      `toolbox_version = "1.1.0"`. Check: `git diff` / file contents.
- [x] **11.2** `tb install` links `review` into both agent skill dirs.
      Check: `readlink ~/.claude/skills/review` and `~/.grok/skills/review`.
- [x] **11.3** `tb doctor` exit 0 with both agents present. Check: exit
      code and skill lines include `review` as `ok`.

## Explicitly out of scope for groups 1–11

Per requirements §3 and §15 as of v1.1: no MCP server, no `tb verify`
wrapper, no third profile (`back-go`), no per-repo skill add/rm CLI, no
persona switch CLI, no other agents beyond Claude Code and Grok Build, no
`staticcheck` / `cargo deny` / `govulncheck` as required tools.

## Group 12 — Spec (v1.2)

- [x] **12.1** ADR `docs/adr/0007-back-go-profile-is-v1.2.md` as `proposed`
      (skill-only v1.1 vs third profile v1.2 vs third profile v2).
      Check: file exists, four MADR sections, three named options,
      Decision is v1.2 + third profile.
- [x] **12.2** `docs/requirements.md` header v1.2; v1 and v1.1 body
      unchanged; §17 delta + §18 acceptance. Check: every §18 item is a
      line you can fail; `rg 'back-go' docs/requirements.md` shows the v1
      non-goal *and* the promotion in the delta.
- [x] **12.3** `docs/design.md` §16–§17 cite 0007 and map §18 → mechanism.
      Check: traceability table has one row per §18 item.
- [x] **12.4** User marks 0007 `accepted` (and this file's v1.2 groups
      accepted). Check: ADR header says `Status: accepted`. Do not start
      Group 13 until this box is ticked.

## Group 13 — Content (v1.2)

- [x] **13.1** `profiles/back-go.toml` per design §16.1. Check:
      `rg 'go-verify", "back-go", "aws-guard' profiles/back-go.toml`;
      `rg 'handoff", "review' profiles/back-go.toml`.
- [x] **13.2** `skills/back-go/SKILL.md` per design §16.3. Check:
      frontmatter `name: back-go`; `rg '^## ' skills/back-go/SKILL.md`
      prints exactly the five headings in requirements §17.4.
- [x] **13.3** Pointers in `go-verify`, `infra-go`, `review` per design
      §16.3. Check: `rg 'back-go' skills/go-verify/SKILL.md
      skills/infra-go/SKILL.md skills/review/SKILL.md` matches all three;
      `go-verify` command block still matches requirements §15.4.

## Group 14 — Binary stamp (v1.2)

- [x] **14.1** `const tbVersion = "1.2.0"` in
      `cmd/tb/internal/cli/init.go`; init usage lists `back-go`; doctor
      treats `back-go` like `infra-go` for `go` on PATH; tests per design
      §16.4. Check: `rg 'tbVersion' cmd/tb` shows `1.2.0`;
      `cd cmd/tb && go test ./internal/config/ ./internal/cli/` covers
      `TestRealProfilesParse` and init `-p back-go`.
- [x] **14.2** From `cmd/tb`: `gofmt -l .` silent; `go vet ./...`;
      `go test ./...`; `go test -race ./...`. Check: all four pass.

## Group 15 — README

- [x] **15.1** `README.md` lists three profiles, `tb init -p back-go`,
      version 1.2.0; drop “There is no `back-go` profile”. Check: the
      profiles table has exactly three rows; `back-go` verify block
      matches `infra-go`.

## Group 16 — Acceptance (requirements §18, run manually on Arch)

- [x] **16.1** `tb init -p back-go` in a fresh Go module: pointer
      `profile = "back-go"` and `toolbox_version = "1.2.0"`; `AGENTS.md`
      has `Profile: back-go`, `review` and `back-go` on `Skills:`,
      `Verify:` matching `go-verify`. Check: file contents.
- [x] **16.2** `tb init --force` in a rust-systems fixture and an
      infra-go fixture: version stamp 1.2.0; prose outside markers
      unchanged; profile unchanged unless `-p`. Check: `git diff` / file
      contents.
- [x] **16.3** `tb install` links `back-go` into both agent skill dirs.
      Check: `readlink ~/.claude/skills/back-go` and
      `~/.grok/skills/back-go`.
- [x] **16.4** `tb doctor` exit 0 with both agents; cwd profile `back-go`
      warns on missing `go`, not `cargo`. Check: exit code and output.

## Explicitly out of scope for groups 12–16

Per requirements §17 and ADR 0007: no MCP server, no `tb verify` wrapper,
no per-repo skill add/rm CLI, no persona switch CLI, no other agents, no
`staticcheck` / `govulncheck` as required tools, no lock of router/ORM/
migrator in the `back-go` skill, no fourth profile. Do not add tasks for
any of these without a new ADR.

## Group 17 — Spec (v1.3)

- [x] **17.1** ADR `docs/adr/0008-game-bevy-profile-is-v1.3.md` as
      `proposed` (patch rust-systems vs generic `game-rust` v1.3 vs
      `game-bevy` v1.3 vs `game-bevy` v2).
      Check: file exists, four MADR sections, four named options,
      Decision is v1.3 + `game-bevy`.
- [x] **17.2** `docs/requirements.md` header v1.3; v1, v1.1, and v1.2
      body unchanged; §19 delta + §20 acceptance. Check: every §20 item
      is a line you can fail; `rg 'game-bevy' docs/requirements.md`
      shows the new profile.
- [x] **17.3** `docs/design.md` §18–§19 cite 0008 and map §20 →
      mechanism. Check: traceability table has one row per §20 item.
- [x] **17.4** User marks 0008 `accepted` (and this file's v1.3 groups
      accepted). Check: ADR header says `Status: accepted`. Do not start
      Group 18 until this box is ticked.

## Group 18 — Content (v1.3)

- [x] **18.1** `profiles/game-bevy.toml` per design §18.1. Check:
      `rg 'rust-verify", "game-bevy' profiles/game-bevy.toml`;
      `rg 'handoff", "review' profiles/game-bevy.toml`;
      `rg 'aws-guard|rust-systems' profiles/game-bevy.toml` does not
      match the `skills` array.
- [x] **18.2** `skills/game-bevy/SKILL.md` per design §18.3. Check:
      frontmatter `name: game-bevy`; `rg '^## ' skills/game-bevy/SKILL.md`
      prints exactly the five headings in requirements §19.4.
- [x] **18.3** Pointers in `rust-verify`, `rust-systems`, `review` per
      design §18.3. Check: `rg 'game-bevy' skills/rust-verify/SKILL.md
      skills/rust-systems/SKILL.md skills/review/SKILL.md` matches all
      three; `rust-verify` command block still matches requirements
      §15.4.

## Group 19 — Binary stamp (v1.3)

- [x] **19.1** `const tbVersion = "1.3.0"` in
      `cmd/tb/internal/cli/init.go`; init usage lists `game-bevy`; doctor
      treats `game-bevy` like `rust-systems` for `cargo` on PATH; tests
      per design §18.4. Check: `rg 'tbVersion' cmd/tb` shows `1.3.0`;
      `cd cmd/tb && go test ./internal/config/ ./internal/cli/` covers
      `TestRealProfilesParse` and init `-p game-bevy`.
- [x] **19.2** From `cmd/tb`: `gofmt -l .` silent; `go vet ./...`;
      `go test ./...`; `go test -race ./...`. Check: all four pass.

## Group 20 — README

- [x] **20.1** `README.md` lists four profiles, `tb init -p game-bevy`,
      version 1.3.0. Check: the profiles table has exactly four rows;
      `game-bevy` verify block matches `rust-systems`.

## Group 21 — Acceptance (requirements §20, run manually on Arch)

- [x] **21.1** `tb init -p game-bevy` in a fresh cargo crate: pointer
      `profile = "game-bevy"` and `toolbox_version = "1.3.0"`; `AGENTS.md`
      has `Profile: game-bevy`, `review` and `game-bevy` on `Skills:`,
      no `rust-systems` on `Skills:`, `Verify:` matching `rust-verify`.
      Check: file contents.
- [x] **21.2** `tb init --force` in a rust-systems fixture, an infra-go
      fixture, and a back-go fixture: version stamp 1.3.0; prose outside
      markers unchanged; profile unchanged unless `-p`. Check:
      `git diff` / file contents.
- [x] **21.3** `tb install` links `game-bevy` into both agent skill dirs.
      Check: `readlink ~/.claude/skills/game-bevy` and
      `~/.grok/skills/game-bevy`.
- [x] **21.4** `tb doctor` exit 0 with both agents; cwd profile
      `game-bevy` warns on missing `cargo`, not `go`. Check: exit code
      and output.

## Explicitly out of scope for groups 17–21

Per requirements §19 and ADR 0008: no MCP server, no `tb verify` wrapper,
no per-repo skill add/rm CLI, no persona switch CLI, no other agents, no
fifth profile, no Bevy version pin, no physics/net/UI crate lock in the
`game-bevy` skill, no wasm as required verify, no screenshot tests, no
`tb` scaffolding of a Bevy template. Do not add tasks for any of these
without a new ADR.

## Group 22 — Spec (v1.4)

- [x] **22.1** ADR `docs/adr/0009-cpp-systems-profile-is-v1.4.md` as
      `proposed` (no profile vs generic `cpp-systems` house style vs
      `cpp-systems` with a ggml-locked house style vs `cpp-ggml` as the
      profile name).
      Check: file exists, four MADR sections, four named options,
      Decision is v1.4 + `cpp-systems` + `cpp-ggml` house style.
- [x] **22.2** `docs/requirements.md` header v1.4; v1 through v1.3 body
      unchanged; §21 delta + §22 acceptance. Check: every §22 item is a
      line you can fail; `rg 'cpp-systems' docs/requirements.md` shows
      the new profile.
- [x] **22.3** `docs/design.md` §20–§21 cite 0009 and map §22 →
      mechanism. Check: traceability table has one row per §22 item.
- [x] **22.4** User marks 0009 `accepted` (and this file's v1.4 groups
      accepted). Check: ADR header says `Status: accepted`. Do not start
      Group 23 until this box is ticked.

## Group 23 — Content (v1.4)

- [x] **23.1** `profiles/cpp-systems.toml` per design §20.1. Check:
      `rg 'cpp-verify", "cpp-ggml' profiles/cpp-systems.toml`;
      `rg 'handoff", "review' profiles/cpp-systems.toml`;
      `rg 'aws-guard|rust-|-go' profiles/cpp-systems.toml` does not match
      the `skills` array.
- [x] **23.2** `skills/cpp-verify/SKILL.md` per design §20.3. Check:
      frontmatter `name: cpp-verify`; the command block is the three
      CMake/ctest lines; the body tells the agent to confirm flags
      against the target repo's CI workflows.
- [x] **23.3** `skills/cpp-ggml/SKILL.md` per design §20.4. Check:
      frontmatter `name: cpp-ggml`; `rg '^## ' skills/cpp-ggml/SKILL.md`
      prints exactly the six headings in requirements §21.6; heading 1
      contains the `tb init` prohibition; heading 2 contains
      `Assisted-by:` and no `Co-authored-by:` recommendation.
- [x] **23.4** Pointer in `review` per design §20.4. Check:
      `rg 'cpp-systems|cpp-ggml' skills/review/SKILL.md` matches;
      `git diff --stat skills/rust-verify skills/go-verify` is empty.

## Group 24 — Binary stamp (v1.4)

- [x] **24.1** `const tbVersion = "1.4.0"` in
      `cmd/tb/internal/cli/init.go`; init usage lists `cpp-systems`;
      doctor warns on missing `cmake` for `cpp-systems`; tests per design
      §20.6. Check: `rg 'tbVersion' cmd/tb` shows `1.4.0`;
      `cd cmd/tb && go test ./internal/config/ ./internal/cli/` covers
      `TestRealProfilesParse` and init `-p cpp-systems`.
- [x] **24.2** From `cmd/tb`: `gofmt -l .` silent; `go vet ./...`;
      `go test ./...`; `go test -race ./...`. Check: all four pass.

## Group 25 — README

- [x] **25.1** `README.md` lists five profiles, `tb init -p cpp-systems`,
      version 1.4.0, and the CMake verify block. Check: the profiles
      table has exactly five rows; the layout block names
      `cpp-systems.toml`.

## Group 26 — Acceptance (requirements §22, run manually on Arch)

- [x] **26.1** `tb init -p cpp-systems` in a fresh git repo: pointer
      `profile = "cpp-systems"` and `toolbox_version = "1.4.0"`;
      `AGENTS.md` has `Profile: cpp-systems`, `review`, `cpp-verify`, and
      `cpp-ggml` on `Skills:`, no Rust or Go skill on `Skills:`,
      `Verify:` matching `cpp-verify`. Check: file contents.
- [x] **26.2** `tb init --force` in a rust-systems, infra-go, back-go,
      and game-bevy fixture: version stamp 1.4.0; prose outside markers
      unchanged; profile unchanged unless `-p`. Check: `git diff` / file
      contents.
- [x] **26.3** `tb install` links `cpp-verify` and `cpp-ggml` into both
      agent skill dirs. Check: `readlink ~/.claude/skills/cpp-verify` and
      `~/.grok/skills/cpp-ggml`.
- [x] **26.4** `tb doctor` exit 0 with both agents; cwd profile
      `cpp-systems` warns on missing `cmake`, not `cargo` or `go`.
      Check: exit code and output.

## Explicitly out of scope for groups 22–26

Per requirements §21 and ADR 0009: no MCP server, no `tb verify` wrapper,
no per-repo skill add/rm CLI, no persona switch CLI, no other agents, no
sixth profile, no `tb` guard against initialising an upstream repo, no
generic C/C++ house style, no CUDA/Metal/Vulkan sub-profile, no Python
recipe for `gguf-py` / `convert_*.py`, and no wrapper around `ci/run.sh`.
Do not add tasks for any of these without a new ADR.

## Group 27 — Spec (v1.5)

- [x] **27.1** ADR `docs/adr/0010-c-cli-profile-is-v1.5.md` as `proposed`
      (no profile / extend `cpp-systems` / `c-systems` generic / `c-cli`
      v1.5), recording the Makefile-with-fallback build and the `gcc`
      doctor probe. Check: file exists, four MADR sections, four named
      options, Decision is v1.5 + `c-cli`.
- [x] **27.2** `docs/requirements.md` header v1.5; v1 through v1.4 body
      unchanged; §23 delta + §24 acceptance. Check: every §24 item is a
      line you can fail; `rg 'c-cli' docs/requirements.md` shows the new
      profile.
- [x] **27.3** `docs/design.md` §22–§23 cite 0010 and map §24 →
      mechanism. Check: traceability table has one row per §24 item.
- [x] **27.4** User marks 0010 `accepted` (and this file's v1.5 groups
      accepted). Check: ADR header says `Status: accepted`. Do not start
      Group 28 until this box is ticked.

## Group 28 — Content (v1.5)

- [x] **28.1** `profiles/c-cli.toml` per design §22.1. Check:
      `rg 'c-verify", "c-cli' profiles/c-cli.toml`;
      `rg 'handoff", "review' profiles/c-cli.toml`;
      `rg 'aws-guard|cpp-|rust-|-go"' profiles/c-cli.toml` does not match
      the `skills` array.
- [x] **28.2** `skills/c-verify/SKILL.md` per design §22.3. Check:
      frontmatter `name: c-verify`; the body contains `make test`, a
      fallback line with `-Werror`, `-fsanitize=address,undefined`, and
      `git status`.
- [x] **28.3** `skills/c-cli/SKILL.md` per design §22.4. Check:
      frontmatter `name: c-cli`; `rg '^## ' skills/c-cli/SKILL.md` prints
      exactly the six headings in requirements §23.5.
- [x] **28.4** Pointer in `review` per design §22.4. Check:
      `rg 'c-cli' skills/review/SKILL.md` matches;
      `git diff --stat skills/rust-verify skills/go-verify skills/cpp-verify
      skills/cpp-ggml` is empty.

## Group 29 — Binary stamp (v1.5)

- [x] **29.1** `const tbVersion = "1.5.0"` in
      `cmd/tb/internal/cli/init.go`; init usage lists `c-cli`; doctor
      warns on missing `gcc` for `c-cli`; tests per design §22.5. Check:
      `rg 'tbVersion' cmd/tb` shows `1.5.0`;
      `cd cmd/tb && go test ./internal/config/ ./internal/cli/` covers
      `TestRealProfilesParse` and init `-p c-cli`.
- [x] **29.2** From `cmd/tb`: `gofmt -l .` silent; `go vet ./...`;
      `go test ./...`; `go test -race ./...`. Check: all four pass.

## Group 30 — README

- [x] **30.1** `README.md` lists six profiles, `tb init -p c-cli`,
      version 1.5.0, and the C verify block. Check: the profiles table has
      exactly six rows; the layout block names `c-cli.toml`.

## Group 31 — Acceptance (requirements §24, run manually on Arch)

- [x] **31.1** `tb init -p c-cli` in a fresh git repo: pointer
      `profile = "c-cli"` and `toolbox_version = "1.5.0"`; `AGENTS.md` has
      `Profile: c-cli`, `review`, `c-verify`, and `c-cli` on `Skills:`, no
      C++, Rust, or Go skill on `Skills:`, `Verify:` matching `c-verify`.
      Check: file contents.
- [x] **31.2** `tb init --force` in a rust-systems, infra-go, back-go,
      game-bevy, and cpp-systems fixture: version stamp 1.5.0; prose
      outside markers unchanged; profile unchanged unless `-p`. Check:
      `git diff` / file contents.
- [x] **31.3** `tb install` links `c-verify` and `c-cli` into both agent
      skill dirs. Check: `readlink ~/.claude/skills/c-verify` and
      `~/.grok/skills/c-cli`.
- [x] **31.4** `tb doctor` exit 0 with both agents; cwd profile `c-cli`
      warns on missing `gcc`, not `cargo`, `go`, or `cmake`. Check: exit
      code and output — on a host with `gcc` installed the manual run
      shows no warning, and the branch is proven by the unit test.

## Explicitly out of scope for groups 27–31

Per requirements §23 and ADR 0010: no MCP server, no `tb verify` wrapper,
no per-repo skill add/rm CLI, no persona switch CLI, no other agents, no
seventh profile, no C library or embedded variant, no ncurses or other
dependency pin, no `valgrind` or `clang-tidy` as required tools, and no
change to the five existing profiles' content. Adding Makefiles and
`.gitignore`s to `algorithm-visualizer` and `clings`, untracking their
committed binaries, and fixing the defects the survey found in
`clings.c` is work in those repos, not here. Do not add tasks for any of
these without a new ADR.
