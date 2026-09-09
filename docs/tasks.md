# Toolbox — Tasks (v1)

Status: accepted
Date: 2026-09-09
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

- [ ] **2.1** `personas/default.md` exactly per design §9.
      Check: two `##` headings present, titles match design.md verbatim.
- [ ] **2.2** `profiles/rust-systems.toml` and `profiles/infra-go.toml` per
      design §3.3, with the exact skill lists from design §8 (five process +
      two or three domain skills, no extras).
      Check: manually diff each `skills` array against requirements §6.2/§6.3.
- [ ] **2.3** Five process skills, one `SKILL.md` each per design §10:
      `skills/spec/`, `skills/adr/`, `skills/onboard/`, `skills/scout/`,
      `skills/handoff/`. Bodies describe the workflow named in requirements
      §6.2 (one sentence per skill there is the minimum content bar).
      Check: each has `name`/`description` frontmatter matching its dir name.
- [ ] **2.4** Five domain skills: `skills/rust-verify/`,
      `skills/rust-systems/`, `skills/go-verify/`, `skills/infra-go/`,
      `skills/aws-guard/`. The two `*-verify` skills must state the exact
      commands from requirements §6.3 (`cargo test`; `cargo clippy -- -D
      warnings`; Miri when `unsafe` in the change / `go test ./...`; `go
      vet`; race when tests involve concurrency).
      Check: verify commands in the skill body match requirements §6.3
      word-for-word.
- [ ] **2.5** `templates/AGENTS.md` (the managed-region skeleton, design §4)
      and `templates/docs/{requirements,design,tasks,session}.md` as minimal
      starter stubs (headings + one-line "fill this in" placeholders) plus
      `templates/docs/adr/0000-template.md` (copy of this repo's own
      `docs/adr/0000-template.md`).
      Check: `templates/AGENTS.md` byte-matches design §4's template block.

## Group 3 — Binary core

- [ ] **3.1** `cmd/tb/go.mod` (`module toolbox/cmd/tb`, `go 1.27`, no
      `require` block per ADR 0001). `cmd/tb/main.go` with subcommand
      dispatch (`install`, `init`, `doctor`, `skill`, `persona`) and a
      top-level usage message on no/unknown args (exit 2).
      Check: `go build ./cmd/tb/...` succeeds with zero dependencies in
      `go.sum` (no `go.sum` file at all is fine).
- [ ] **3.2** `internal/paths`: `Resolve() (string, error)` per design §1.
      Check: table tests — env set to existing dir, env unset with
      `~/toolbox` present, env set to missing dir (error), env empty string
      (falls through to default).
- [ ] **3.3** `internal/config`: TOML-subset reader per design §3.1, plus
      typed loaders `LoadProfile(path) (Profile, error)` and
      `LoadPointer(path) (Pointer, error)` for the two schemas in §3.2/§3.3.
      Check: table tests covering every construct in the subset and at
      least three unsupported-syntax rejections with correct `file:line`.
- [ ] **3.4** `internal/agents`: `Detect() (claude, grok bool)` per ADR
      0003 (PATH lookup OR config-dir check, injectable for tests).
      Check: tests fake `PATH` and `$HOME`/`$USERPROFILE`-equivalent to
      cover all four presence combinations without touching the real
      environment.
- [ ] **3.5** `internal/render`: managed-region template execution and the
      three-state merge logic from design §4 (`Render(profile, persona)
      string`, `Merge(existing, region string) (result string, refused
      bool)`).
      Check: golden-file test for `Render`; `Merge` tests for "no markers",
      "markers present", confirming content outside markers survives
      byte-for-byte.

## Group 4 — Commands

- [ ] **4.1** `tb install` (design §5.1 + §6). Creates both skill dirs,
      links/relinks/skips per the symlink strategy, prints the summary,
      applies the exit-code policy from design §7.
      Check: integration test in `t.TempDir()` with a fake `$TOOLBOX_HOME`
      containing two skills and fake `$HOME`; asserts both dirs get correct
      symlinks and a pre-existing real file at one target is skipped with a
      warning, not fatal.
- [ ] **4.2** `tb init` (design §5.2). Profile load, git-repo check,
      `toolbox.toml` write/merge, `AGENTS.md` render/merge via
      `internal/render`, template materialization (create-if-missing,
      `docs/adr/*` never touched beyond the create-if-missing template
      copy).
      Check: integration tests for: empty repo (all files created); repo
      with existing `docs/adr/0001-something.md` (untouched, byte-identical
      before/after); re-run without `--force` on existing `AGENTS.md`
      (non-interactive refuses, exit 1, file unchanged); `--force` (only
      `AGENTS.md` + `toolbox.toml` change).
- [ ] **4.3** `tb doctor` (design §5.3). Read-only, one-screen output, exit
      code only on zero-agents.
      Check: tests for "clean" (both agents faked present), "one missing"
      (exit 0), "zero agents" (exit 1); a test with no `toolbox.toml` in
      cwd skips steps 5–6 without erroring.
- [ ] **4.4** `tb skill new <name>` (design §5.4).
      Check: test for valid name (scaffold created, next-steps text
      printed), invalid name (exit 2), already-exists (exit 1).
- [ ] **4.5** `tb persona show` (design §5.5).
      Check: test reading a fixture profile + `personas/default.md`,
      asserting both section titles print exactly.

## Group 5 — Build and CI

- [ ] **5.1** Build script (`Makefile` or `cmd/tb/build.sh` — pick one, keep
      it under ~20 lines) producing `dist/tb-linux-amd64` and
      `dist/tb-darwin-arm64` via the flags in design §11.
      Check: both binaries build locally on Arch (darwin/arm64 as a
      cross-compile check, `file dist/tb-darwin-arm64` reports the right
      target).
- [ ] **5.2** CI workflow: `gofmt -l`, `go vet ./...`, `go test ./...` on
      push, scoped to `cmd/tb/**` paths only.
      Check: workflow passes on this repo's initial commit containing all
      of the above.

## Group 6 — Acceptance (requirements §13, run manually)

- [ ] **6.1** On Arch: clone to `TOOLBOX_HOME` (or default), build/install
      `tb`, run `tb install` — confirm skills link, `tb doctor` is clean
      aside from any genuinely-missing agent.
- [ ] **6.2** `tb init -p rust-systems` in a fresh empty cargo repo and in a
      non-empty rust repo that already has `docs/adr/0001-*.md` — confirm
      pointer + `AGENTS.md` + `docs/` created, existing ADR untouched.
- [ ] **6.3** `tb init -p infra-go` in a Go module — same checks.
- [ ] **6.4** Re-run `tb init` without `--force` on each repo from 6.2/6.3 —
      confirm merge-or-refuse behavior (interactive prompt in a real TTY;
      refuse in a script).
- [ ] **6.5** Re-run with `--force` — confirm only `AGENTS.md` and
      `toolbox.toml` changed (`git diff --stat`).
- [ ] **6.6** `tb skill new demo` — confirm scaffold under
      `$TOOLBOX_HOME/skills/demo`.
- [ ] **6.7** Open Claude Code and/or Grok Build (whichever is installed on
      that machine) and confirm the symlinked skills are visible.
- [ ] **6.8** Repeat 6.1–6.7 on the MacBook Air M1.
- [ ] **6.9** Confirm `tb doctor` fails closed (exit 1) if both agent config
      dirs are temporarily renamed away, and warns (exit 0) with one
      renamed back.

## Explicitly out of scope for these tasks

Per requirements §3 and §14: no MCP server, no `tb verify` wrapper, no
third profile (`back-go`), no per-repo skill add/rm CLI, no persona switch
CLI, no other agents beyond Claude Code and Grok Build. Do not add tasks
for any of these without a new ADR recording why the non-goal changed.
