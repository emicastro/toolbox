# Toolbox — Tasks (v1.1)

Status: v1 groups 1–6 accepted 2026-09-09; v1.1 groups 7–10 accepted 2026-09-11; group 11 pending
Date: 2026-09-11
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

- [ ] **11.1** `tb init --force` in a rust-systems fixture and an infra-go
      fixture: `AGENTS.md` contains the five rules and `review` on
      `Skills:`; prose outside markers unchanged; `toolbox.toml` has
      `toolbox_version = "1.1.0"`. Check: `git diff` / file contents.
- [ ] **11.2** `tb install` links `review` into both agent skill dirs.
      Check: `readlink ~/.claude/skills/review` and `~/.grok/skills/review`.
- [ ] **11.3** `tb doctor` exit 0 with both agents present. Check: exit
      code and skill lines include `review` as `ok`.

## Explicitly out of scope for these tasks

Per requirements §3 and §15: no MCP server, no `tb verify` wrapper, no
third profile (`back-go`), no per-repo skill add/rm CLI, no persona switch
CLI, no other agents beyond Claude Code and Grok Build, no `staticcheck` /
`cargo deny` / `govulncheck` as required tools. Do not add tasks for any
of these without a new ADR recording why the non-goal changed.
