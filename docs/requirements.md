# Toolbox — Requirements (v1.4)

Status: v1 accepted 2026-09-09; v1.1 delta accepted 2026-09-11; v1.2 delta accepted 2026-09-12; v1.3 delta accepted 2026-09-12; v1.4 delta accepted 2026-09-18  
Date: 2026-09-18  
Command: `tb`  
Source of truth for this product: this file. v1, v1.1, v1.2, and v1.3 text below is unchanged and remains accepted. v1.1 is the **Delta from v1** and **Acceptance (v1.1)** sections (ADR 0005). v1.2 is the **Delta from v1.1** and **Acceptance (v1.2)** sections (ADR 0007). v1.3 is the **Delta from v1.2** and **Acceptance (v1.3)** sections (ADR 0008). v1.4 is the **Delta from v1.3** and **Acceptance (v1.4)** sections at the end (ADR 0009). Design and tasks for each delta are produced in Plan Mode in the `toolbox` repo.

## 1. Problem

Coding agents (Claude Code, Grok Build, and others later) start empty every session and every repo. Skills, verify recipes, ADR conventions, and tone leak into ad-hoc `AGENTS.md` files and chats. That knowledge does not travel between an Arch Linux workstation and a MacBook Air M1, and it does not travel between greenfield and legacy repos.

Toolbox is a **developer harness**: a single git repo on the machine plus a small static Go binary. It is not an MCP server, not a daemon, not a marketplace, and not part of any product crate.

## 2. Goals

- One clone (`TOOLBOX_HOME`, default `~/toolbox`) shared by all product repos.
- One command to wire a product repo to a named profile (`tb init -p rust-systems`).
- Skills and persona live in toolbox; product repos keep a short pointer + generated `AGENTS.md` + their own `docs/` (including ADRs).
- Works on Arch Linux (x86_64) and macOS Apple Silicon (arm64).
- Agent-agnostic *content* (Markdown skills + AGENTS). v1 *install targets* are Claude Code and Grok Build only.
- Docs and code in English. Agent chat defaults to English; if the user writes Spanish, the agent may reply in Spanish.

## 3. Non-goals (v1)

- MCP server, tool schemas, or a local daemon.
- Engram / SQLite / automatic agent memory.
- Gentle-AI-style presets, TUI installer, RDD, GGA.
- Per-repo skill enable/disable CLI (`tb skill add` / `rm` on a product repo).
- Persona switch CLI.
- `back-go` profile (HTTP services, databases).
- Wrapping verify as `tb verify`.
- Sync besides `git pull` in `$TOOLBOX_HOME` + `tb install`.
- Installing or updating Claude Code / Grok Build.
- Other agents (Cursor, Codex, OpenCode, …).

## 4. Users and machines

Single user. Two machines: Arch Linux (custom), MacBook Air M1. SuperGrok (Grok Build) and Claude Pro (Claude Code) are available; the harness must not assume both binaries exist on every machine.

## 5. Layout

Clone path: environment variable `TOOLBOX_HOME` if set and non-empty, otherwise `~/toolbox`. The clone is a git repo. Not `~/.config/...` as the source of truth. `tb` must resolve this path once per invocation and use it for skills, profiles, templates, and `skill new`. Document the variable in `tb doctor` output.

Expected tree (names may be adjusted in design, not the roles):

```
$TOOLBOX_HOME/   # default: ~/toolbox
  AGENTS.md                 # global rules (short)
  toolbox.toml              # toolbox-side config if needed
  profiles/
    rust-systems.toml
    infra-go.toml
  personas/
    default.md              # two sections: Plan (Staff), Implement (Senior)
  skills/
    <name>/SKILL.md
  templates/                # files init may create if missing
    AGENTS.md
    docs/requirements.md
    docs/design.md
    docs/tasks.md
    docs/session.md
    docs/adr/0000-template.md
  bin/                      # optional user scripts; not MCP tools
  cmd/tb/                   # Go module for the tb binary
```

Product repo after init:

```
<product>/
  toolbox.toml              # pointer: profile = "rust-systems"
  AGENTS.md                 # generated, short
  docs/
    requirements.md         # created if missing
    design.md               # created if missing
    tasks.md                # created if missing
    session.md              # created if missing
    adr/                    # created if missing; never overwritten
```

Skills are **not** copied into the product repo. `tb install` symlinks:

- `~/.claude/skills/<name>` → `$TOOLBOX_HOME/skills/<name>`
- `~/.grok/skills/<name>` → `$TOOLBOX_HOME/skills/<name>`

## 6. Profiles

A profile is a named pack: persona file, process skills, domain skills, init templates, verify *recipe documented in a skill* (not a `tb` subcommand).

### 6.1 v1 profiles

| Profile | Use |
|---|---|
| `rust-systems` | Systems / low-level Rust, including greenfield and legacy crates |
| `infra-go` | Go CLIs, scripts, AWS/infra glue, one-shot jobs — not long-running HTTP backends |

`back-go` is deferred.

### 6.2 Process skills (both profiles)

Must be listed in both profiles:

- `spec` — from idea/requirements to design + tasks (Plan Mode)
- `adr` — mandatory Architecture Decision Record before a design fork is treated as settled
- `onboard` — brownfield: map a repo the user has not touched
- `scout` — cheap exploration; do not spend a frontier model on `rg`
- `handoff` — write `docs/session.md` and stop (compact / next session)

### 6.3 Domain skills (starting list; Plan Mode may refine names, not drop process skills)

`rust-systems` also includes:

- `rust-verify` — `cargo test`; `cargo clippy -- -D warnings`; Miri when `unsafe` is in the change
- `rust-systems` — ownership, no unwrap in libraries, explicit concurrency, document unsafe

`infra-go` also includes:

- `go-verify` — `go test ./...`; `go vet`; race when tests involve concurrency
- `infra-go` — exit codes, stderr logs, `context` cancellation, explicit AWS profile, no prod by default
- `aws-guard` — refuse implicit prod, no secrets in repo, profile/region must be named

Plan Mode may add at most a small number of extra domain skills if they are load-bearing. It must not import a marketplace pack.

## 7. Persona

Single file: `personas/default.md`.

- **Plan section — Staff Engineer:** challenges weak specs, demands ADRs for forks, does not implement past the approved task list, prefers evidence (`cargo`/`go` output) over narration.
- **Implement section — Senior Engineer:** follows accepted ADRs and `docs/tasks.md`, does not reopen design without a new ADR, keeps diffs small, runs the profile verify recipe before claiming done.

The model tier (Opus, Grok 4.6, Sonnet) does **not** change the role. A cheaper model on implement still follows the Implement section.

Language: default English for chat and all artifacts. If the user writes in Spanish, replies may be Spanish; code, identifiers, comments, commit messages, ADRs, and docs stay English.

## 8. CLI

Static Go binary named `tb`, built for linux/amd64 and darwin/arm64.

Global flag: `--profile` with alias `-p`.

### 8.1 `tb install`

Machine-level. Run after clone and after `git pull` on toolbox.

- Create `~/.claude/skills` and `~/.grok/skills` if needed.
- Symlink each `$TOOLBOX_HOME/skills/<name>` into both skill dirs (idempotent).
- Do not modify product repos.
- If exactly one of Claude Code / Grok Build is missing: print a warning, exit 0 if install otherwise succeeded.
- If both are missing: fail (non-zero).
- Detection is “binary on PATH or known config dir present” — exact heuristic is a design choice.

### 8.2 `tb init [--profile|-p <name>]`

Run from a product repo root (greenfield or legacy). Must work after `cd` into any git repo.

- Require a profile name (`rust-systems` or `infra-go` in v1).
- Write `toolbox.toml` pointer (`profile = "..."`).
- Generate short `AGENTS.md` from the profile (English).
- Create missing `docs/` templates listed in §5. Never overwrite existing `docs/adr/*`.
- Existing toolbox-managed files (`AGENTS.md`, `toolbox.toml`):
  - interactive: offer **merge** (keep user prose, refresh toolbox-owned stubs) or **refuse**
  - `--force`: overwrite only `AGENTS.md` and `toolbox.toml`. Never delete or rewrite existing ADR bodies.
- Do not copy skills into the repo.

### 8.3 `tb doctor`

Read-only.

Reports: toolbox clone present; `tb` on PATH; skill symlinks healthy; Claude and/or Grok detected; if cwd has `toolbox.toml`, profile exists and `AGENTS.md` is present; if cwd looks like the profile’s language, warn when `cargo` or `go` is missing. One-screen output. No `--fix` in v1.

### 8.4 `tb skill new <name>`

Scaffolds `$TOOLBOX_HOME/skills/<name>/SKILL.md` (English template). Does not attach the skill to a profile. User edits the skill, adds the name to the chosen profile file, runs `tb install`, commits in `$TOOLBOX_HOME`.

### 8.5 `tb persona show`

Prints which persona file the current profile uses and the two section titles. No mutation.

### 8.6 Out of v1

`tb verify`, `tb skill add|rm` on a product repo, `tb persona set`, attach of extra agents, uninstall beyond “remove symlinks by hand / documented later”.

## 9. ADRs

Mandatory in v1 for both profiles.

- `tb init` ensures `docs/adr/` and a template exist.
- New architectural forks require a new `docs/adr/NNNN-title.md` (MADR-short: context, options, decision, consequences).
- Superseding a decision is a new ADR; old ADR marked superseded. Do not rewrite history.
- `--force` on init must not touch existing ADR files.
- Making ADRs optional is a later flag, not v1.

## 10. Memory rules (harness policy, not a database)

- Durable project knowledge: `docs/` and `docs/adr/` in the product repo.
- Durable *user* conventions: markdown in `$TOOLBOX_HOME`.
- Session bridge: `docs/session.md` via `handoff` (overwrite OK).
- Agents must not treat uncommitted chat summaries as source of truth.
- No SQLite memory in v1.

## 11. Verify

Documented only in profile domain skills (`rust-verify`, `go-verify`). The agent runs those commands in its own shell. Success is the command output, not a `tb` wrapper.

## 12. Compatibility and install targets

v1 writes skills for:

- Claude Code: `~/.claude/skills`
- Grok Build: `~/.grok/skills`

Content format: standard `SKILL.md` so a future agent can be added by teaching `tb install` a third symlink path.

## 13. Acceptance (v1 done when)

On both Arch and M1, after clone to `TOOLBOX_HOME` (or default `~/toolbox`) + build/install of `tb`:

1. `tb install` links skills; doctor is clean aside from optional missing agent.
2. `tb init -p rust-systems` in an empty cargo repo and in a non-empty rust repo produces pointer + `AGENTS.md` + `docs/` without destroying existing ADRs.
3. `tb init -p infra-go` same for a Go module.
4. Re-running init without `--force` prompts merge or refuse.
5. `--force` refreshes `AGENTS.md` + `toolbox.toml` only.
6. `tb skill new demo` creates a scaffold under `$TOOLBOX_HOME/skills/demo`.
7. Claude Code and Grok Build, when present, see the symlinked skills.
8. `tb doctor` fails closed on zero agents after install; warns on one missing agent.

## 14. Open points for Plan Mode (not requirements)

- Exact `toolbox.toml` schema and merge algorithm for `AGENTS.md`.
- Go module layout for `tb` (cobra vs std `flag`).
- How agent detection is implemented.
- Precise SKILL.md frontmatter required by Claude Code / Grok Build in 2026.
- Whether global `$TOOLBOX_HOME/AGENTS.md` is also copied or only referenced.
- Resolution order if `TOOLBOX_HOME` points at a missing directory (fail with a clear error).
- CI for `tb` itself (`go test`).

Plan Mode may refine skill *bodies* and add a few domain skills. It must not add MCP, extra profiles, or per-repo skill overrides.

## 15. Delta from v1 (v1.1)

Additive. Does not reopen §3 non-goals. Versioning: ADR 0005. Managed-region payload: ADR 0006.

### 15.1 Binding rules in the product repo

After `tb init`, the managed region of `AGENTS.md` must state five binding rules in English (wording locked in ADR 0006), not only profile, persona titles, skill list, verify summary, and language. The same five lines appear in `personas/default.md` and in toolbox's own `AGENTS.md`.

### 15.2 Process skills

§6.2's list of five is amended: both profiles also include `review`, after `handoff`. `review` is the last gate before ticking a task or claiming done. It is not a `tb` subcommand and not a domain skill.

### 15.3 Skill skip clauses

- `onboard` must not skip merely because `AGENTS.md` exists. Skip only when `docs/session.md` already contains an onboard map for this repo and the area of work has not drifted.
- `adr` options must be viable, not strawmen. Classifying a fork as "implementation detail" to avoid an ADR is itself a fork.
- `spec`: a task `Check:` must be a command whose non-zero exit fails the task, or a byte-level before/after assertion. "File exists" is valid only when the whole task is creating that file.
- `*-verify`: touching CI, scripts, IaC, or SQL is not a skip.
- `handoff`: `Verify status: Not run` is not allowed if the session edited code.
- `review`: skip only when this session produced no diff.

### 15.4 Verify recipes

Documented only in `rust-verify` / `go-verify` (still no `tb verify`).

Go: `gofmt -l .` (must print nothing); `go vet ./...`; `go test ./...`; `go test -race ./...` always. Run from the module root.

Rust: `cargo fmt --check`; `cargo test`; `cargo clippy --all-targets -- -D warnings`. Run `cargo miri test` when a **changed crate contains any `unsafe`**, even if this diff only edits safe code in that crate. Skip Miri only when the changed package(s) have no `unsafe` in `.rs` sources. If Miri is required and not installed, the recipe fails.

`#[allow(...)]` requires the lint name and a reason that is not "to make clippy pass".

Profile `[verify].summary` must match these recipes so the managed region's `Verify:` line is accurate.

### 15.5 House style

`rust-systems`: library code does not `todo!`, `unimplemented!`, or `panic!` on caller input (same bucket as `unwrap`).

`infra-go`: do not ignore `error`; wrap with `%w` when adding context.

`aws-guard`: the "no secrets in the repo" rule triggers on `.env`, keys, tokens, or connection strings even when the diff does not mention AWS.

### 15.6 Version stamp

`tb init` / `tb init --force` writes `toolbox_version = "1.1.0"`. This repo's `schema_version` stays `"1"`.

## 16. Acceptance (v1.1 done when)

On Arch, after `git pull` in `$TOOLBOX_HOME` and a rebuild of `tb`:

1. `tb init --force` in a rust-systems fixture and an infra-go fixture refreshes `AGENTS.md` with the five binding rules and `review` on the `Skills:` line; prose outside the markers is unchanged; `toolbox.toml` has `toolbox_version = "1.1.0"`.
2. `tb install` links `review` into `~/.claude/skills` and `~/.grok/skills`.
3. `tb doctor` is still clean aside from a genuinely missing agent (exit 0 with both present).
4. `skills/onboard/SKILL.md` does not treat a mere `AGENTS.md` as a map.
5. `go-verify` and `rust-verify` command blocks match §15.4, including the crate-level Miri trigger.
6. `skills/review/SKILL.md` exists with `name`/`description` frontmatter; both profiles list `review` after `handoff`.

M1 is a follow-up, not a v1.1 gate (v1 already proved both machines).

## 17. Delta from v1.1 (v1.2)

Additive. Does not rewrite §3; the `back-go` bullet there is history. Versioning and the third profile: ADR 0007. Remaining v1 non-goals stay non-goals.

### 17.1 Third profile

§6.1 is amended: three profiles.

| Profile | Use |
|---|---|
| `rust-systems` | Systems / low-level Rust, including greenfield and legacy crates |
| `infra-go` | Go CLIs, scripts, AWS/infra glue, one-shot jobs — not long-running HTTP backends |
| `back-go` | Go HTTP services and APIs, greenfield and legacy — not CLIs or one-shot jobs |

`back-go` is no longer deferred. `infra-go` does not absorb HTTP.

### 17.2 Domain skills (`back-go`)

Process skills are unchanged from §15.2 (six, including `review` after `handoff`).

`back-go` also includes:

- `go-verify` — the same recipe as §15.4 (`gofmt -l .`; `go vet ./...`; `go test ./...`; `go test -race ./...`). One recipe, two Go profiles. No second verify skill.
- `back-go` — request context is the deadline; graceful shutdown; errors map at the HTTP edge; structured logs with no secrets; parameterized SQL and explicit migrations. Does not name a router, ORM, or migrator (product ADR).
- `aws-guard` — unchanged; backends carry connection strings and secrets even when the diff does not mention AWS.

Mutual skip: `infra-go` skips the listen loop (use `back-go`); `back-go` skips short-lived CLIs, jobs, and AWS/infra glue (use `infra-go`).

### 17.3 `tb init` profile names

§8.2 is amended: require a profile name `rust-systems`, `infra-go`, or `back-go`. Discovery of profile files on disk is not in this increment.

### 17.4 House style (`back-go`)

Five rules, same bar as `infra-go` / `rust-systems`. Changing any at repo scope is an ADR.

1. Request context is the deadline — `r.Context()`; every blocking call honours it; do not store a context in a struct; server and client timeouts are set. `ListenAndServe` defaults are not acceptable.
2. Graceful shutdown — SIGINT/SIGTERM → `http.Server.Shutdown` with a bounded context. No `os.Exit` from a handler. `http.ErrServerClosed` is success.
3. Errors map at the HTTP edge — internal layers return `error`; one mapper at the boundary; never panic on request input; wrap with `%w`; do not ignore `error`; do not leak internal strings to clients.
4. Structured logs, no secrets — method, path, status, duration, request id. Never Authorization, cookies, passwords, tokens, or connection strings. In a long-running service stdout and stderr are both diagnostics (the CLI "stdout is the result" rule does not apply). Secrets in the repo remain `aws-guard`.
5. Parameterized SQL, explicit migrations — bound parameters only; schema changes are migration files applied by a named command, not auto-migrate on boot in prod; a transaction has one owner. Default listen/config/DSN is local/dev; prod is named.

Auth, OpenAPI, health/ready, pagination, and middleware catalogs are not house style. Product ADRs if needed.

### 17.5 Version stamp

`tb init` / `tb init --force` writes `toolbox_version = "1.2.0"`. This repo's `schema_version` stays `"1"`.

## 18. Acceptance (v1.2 done when)

On Arch, after `git pull` in `$TOOLBOX_HOME` and a rebuild of `tb`:

1. `tb init -p back-go` in a fresh Go module writes `toolbox.toml` with `profile = "back-go"` and `toolbox_version = "1.2.0"`; `AGENTS.md` has `Profile: back-go`, `review` and `back-go` on `Skills:`, and `Verify:` matching the `go-verify` recipe.
2. `tb init --force` in an existing rust-systems fixture and an infra-go fixture writes `toolbox_version = "1.2.0"`; prose outside the managed-region markers is unchanged; the profile line does not change unless `-p` names a different profile.
3. `tb install` links `back-go` into `~/.claude/skills` and `~/.grok/skills`.
4. `tb doctor` is still clean aside from a genuinely missing agent (exit 0 with both present); with cwd profile `back-go` it warns on missing `go`, not on missing `cargo`.
5. `skills/back-go/SKILL.md` exists with `name`/`description` frontmatter and exactly the five house-style headings in §17.4.
6. `go-verify` command block is unchanged from §15.4; its description names both Go profiles. `infra-go` and `review` skill bodies mention `back-go`. `profiles/back-go.toml` lists `review` after `handoff` and domain skills `go-verify`, `back-go`, `aws-guard`.

M1 is a follow-up, not a v1.2 gate.

## 19. Delta from v1.2 (v1.3)

Additive. Does not rewrite §3, §15, or §17. Versioning and the fourth profile: ADR 0008. Remaining v1 non-goals stay non-goals.

### 19.1 Fourth profile

§6.1 / §17.1 are amended: four profiles.

| Profile | Use |
|---|---|
| `rust-systems` | Systems / low-level Rust, including greenfield and legacy crates — not Bevy games |
| `infra-go` | Go CLIs, scripts, AWS/infra glue, one-shot jobs — not long-running HTTP backends |
| `back-go` | Go HTTP services and APIs, greenfield and legacy — not CLIs or one-shot jobs |
| `game-bevy` | Bevy games in Rust, greenfield and legacy — not systems crates or non-Bevy engines |

`rust-systems` does not absorb Bevy gameplay.

### 19.2 Domain skills (`game-bevy`)

Process skills are unchanged from §15.2 (six, including `review` after `handoff`).

`game-bevy` also includes:

- `rust-verify` — the same recipe as §15.4 (`cargo fmt --check`; `cargo test`; `cargo clippy --all-targets -- -D warnings`; Miri when a **changed crate** contains any `unsafe` in its own `.rs` sources). One recipe, two Rust profiles. No second verify skill.
- `game-bevy` — ECS is the architecture; panic on broken world invariants, `Result` for I/O; the schedule is the concurrency model; handles and components are values; do not block the frame. Does not name a Bevy version, physics crate, net crate, or UI crate (product ADR).

No `aws-guard`. No `rust-systems` on this profile's skill list.

Mutual skip: `rust-systems` skips Bevy `App` / gameplay (use `game-bevy`); `game-bevy` skips systems/low-level crates, libraries, and non-Bevy engines (use `rust-systems`).

### 19.3 `tb init` profile names

§8.2 / §17.3 are amended: require a profile name `rust-systems`, `infra-go`, `back-go`, or `game-bevy`. Discovery of profile files on disk is not in this increment.

### 19.4 House style (`game-bevy`)

Five rules, same bar as `rust-systems` / `infra-go` / `back-go`. Changing any at repo scope is an ADR.

1. ECS is the architecture — Components are data, systems are behavior, plugins own a domain. `main` composes plugins. Do not model entities as objects with methods that reach into other entities. Automated tests do not take a window or a GPU: add the plugin under test to a headless `App`.
2. Panic on broken world invariants; `Result` for I/O — Missing required component, unique entity that is not unique, schedule that should have made a state impossible: panic with a message (`expect("…")`). File, asset, parse, and network failures return `Result` and are not `unwrap`'d. Do not thread `Result` through every system to avoid a panic that means the world is corrupt.
3. The schedule is the concurrency model — Systems run in parallel unless ordered. Do not share game state with `Arc<Mutex<_>>` or OS threads. Use `Commands`, events, and `Resources`. If two systems conflict, fix the schedule (`before` / `after` / `SystemSet`), do not add a lock. Exclusive `&mut World` is the rare case that needs the whole world.
4. Handles and components are values — Cloning `Handle<T>`, copying `Entity`, and owning components is normal. Borrow-across-systems is the smell, not clone. Do not clone a `Resource`'s inner collection to dodge the borrow checker — split the resource or the system.
5. Do not block the frame — No blocking I/O, asset decode, or network in `Update` / `FixedUpdate`. Load through Bevy assets (`Handle` + load/asset events). Motion uses `Time` deltas, not an assumed frame rate. Work that must block goes on Bevy's task pools and comes back as an event.

Bevy version, physics/net/UI crates, wasm, screenshot tests, `cargo run` as verify, and `[profile.dev]` opt-level are not house style. Product ADRs or product convention if needed.

### 19.5 Version stamp

`tb init` / `tb init --force` writes `toolbox_version = "1.3.0"`. This repo's `schema_version` stays `"1"`.

## 20. Acceptance (v1.3 done when)

On Arch, after `git pull` in `$TOOLBOX_HOME` and a rebuild of `tb`:

1. `tb init -p game-bevy` in a fresh cargo crate writes `toolbox.toml` with `profile = "game-bevy"` and `toolbox_version = "1.3.0"`; `AGENTS.md` has `Profile: game-bevy`, `review` and `game-bevy` on `Skills:`, no `rust-systems` on `Skills:`, and `Verify:` matching the `rust-verify` recipe.
2. `tb init --force` in an existing rust-systems fixture, an infra-go fixture, and a back-go fixture writes `toolbox_version = "1.3.0"`; prose outside the managed-region markers is unchanged; the profile line does not change unless `-p` names a different profile.
3. `tb install` links `game-bevy` into `~/.claude/skills` and `~/.grok/skills`.
4. `tb doctor` is still clean aside from a genuinely missing agent (exit 0 with both present); with cwd profile `game-bevy` it warns on missing `cargo`, not on missing `go`.
5. `skills/game-bevy/SKILL.md` exists with `name`/`description` frontmatter and exactly the five house-style headings in §19.4.
6. `rust-verify` command block is unchanged from §15.4; its description names both Rust profiles. `rust-systems` and `review` skill bodies mention `game-bevy`. `profiles/game-bevy.toml` lists `review` after `handoff` and domain skills `rust-verify`, `game-bevy`.

M1 is a follow-up, not a v1.3 gate.

## 21. Delta from v1.3 (v1.4)

Additive. Does not rewrite §3, §15, §17, or §19. Versioning and the fifth profile: ADR 0009. Remaining v1 non-goals stay non-goals.

### 21.1 Fifth profile

§6.1 / §17.1 / §19.1 are amended: five profiles.

| Profile | Use |
|---|---|
| `rust-systems` | Systems / low-level Rust, including greenfield and legacy crates — not Bevy games |
| `infra-go` | Go CLIs, scripts, AWS/infra glue, one-shot jobs — not long-running HTTP backends |
| `back-go` | Go HTTP services and APIs, greenfield and legacy — not CLIs or one-shot jobs |
| `game-bevy` | Bevy games in Rust, greenfield and legacy — not systems crates or non-Bevy engines |
| `cpp-systems` | C/C++ built with CMake — ggml-family repos (llama.cpp) and the user's own C++ projects |

The profile names the language and build surface; the house-style skill carries the family lock (ADR 0009).

### 21.2 Domain skills (`cpp-systems`)

Process skills are unchanged from §15.2 (six, including `review` after `handoff`).

`cpp-systems` also includes:

- `cpp-verify` — CMake configure with `-DLLAMA_FATAL_WARNINGS=ON`, build, `ctest -L main`, `clang-format` on added lines only, `test-backend-ops` across two backends when `ggml/` changed, a sanitizer build for memory-touching diffs, `ci/run.sh` before publishing. Unlike `cargo` and `go` there is no universal C++ recipe: the skill states that the target repo's CI workflows are the source of truth and that flags are confirmed there, not recalled.
- `cpp-ggml` — the repo's own rules win; never speak for the contributor; understanding is the deliverable; blend in; comments last and short; ggml facts that reviews fail on.

No `aws-guard`. No `rust-*` and no `*-go` skill on this profile's skill list.

Mutual skip: `rust-systems` is Rust only and is unaffected; `cpp-ggml` skips non-ggml C/C++ and the target repo's Python tooling (`gguf-py`, `convert_*.py`).

### 21.3 Upstream repos are not initialised

New constraint, no earlier profile needed it. `tb init` is run only in a repo the user owns. In an upstream repo (`ggml-org/llama.cpp` is the motivating case) it would append a toolbox managed region to a tracked `AGENTS.md` and scaffold `docs/requirements.md`, `docs/design.md`, `docs/tasks.md`, and `docs/session.md` into the project's own `docs/`.

For upstream work the two domain skills carry the value on their own: `tb install` links them machine-wide and the agent loads them by description, with no pointer file and no change to the upstream tree.

This is a rule in the `cpp-ggml` skill. `tb` gains no guard, no `--no-write` flag, and no upstream detection in this increment.

### 21.4 Target-repo rules outrank toolbox

Where the target repo's `AGENTS.md` or `CONTRIBUTING.md` disagrees with the persona or with the `cpp-ggml` skill, the target repo wins. `cpp-ggml` says so in its first rule and points at those files rather than restating them as toolbox policy, because they change upstream.

Two consequences are load-bearing enough to name here: an agent does not write PR descriptions, commit messages, issues, review comments, or replies in such a repo, and does not run `git push` / `gh pr create` / `gh pr comment`; and when the user explicitly asks for a commit there, the trailer is `Assisted-by: <assistant name>`, never `Co-authored-by:`.

### 21.5 `tb init` profile names

§8.2 / §17.3 / §19.3 are amended: require a profile name `rust-systems`, `infra-go`, `back-go`, `game-bevy`, or `cpp-systems`. Discovery of profile files on disk is still not in this increment.

### 21.6 House style (`cpp-systems`)

Six rules, same bar as the other house-style skills, except that rule 1 subordinates the rest to the target repo. Changing any at repo scope is an ADR in that product.

1. The repo's own rules win — read the target repo's `AGENTS.md` and `CONTRIBUTING.md` every session, not from memory; where they disagree with this skill or the persona, they win. Carries §21.3: do not run `tb init` here, plan in a scratch directory outside the repo, rely on `tb install` having linked the skills machine-wide. Check the repo's own `skills/` directory for one covering the task.
2. Never speak for the contributor — no PR description, commit message, issue, review comment, or reply, not even as a draft; no `git push` / `gh pr create` / `gh pr comment` / `gh issue create`; reading commands (`gh search issues`, `gh search prs`, `grep`) are encouraged and duplicates must be searched for first; when the user explicitly asks for a commit the trailer is `Assisted-by: <assistant name>`, never `Co-authored-by:`; the AI-usage disclosure in the PR template is the user's to write.
3. Understanding is the deliverable — a merged line is an indefinite maintenance obligation, so guide before solving, verify comprehension before writing a change, prefer the simpler change that does 90%. Features start as an issue; a bug fix needs a reproducible issue and a regression test that fails before and passes after; one PR per concern; a first PR for a new model or feature is CPU-only.
4. Blend in — read the surrounding code and match it; stop and warn the user when the change introduces a new pattern or is large. No third-party dependencies, no new subsystem, no fancy STL or templates; `snake_case`, longest-common-prefix naming, `<class>_<method>` public API, prefixed upper-case enum values, sized integer types in the public API, lowercase-dash filenames, cross-platform always.
5. Comments last, and short — write the code, then comment only where needed; one or two lines; simple English; explain a non-obvious invariant, never what the code says; no hard-wrapping mid-sentence; no comment that addresses the current task; copied code keeps the comments it had and gains none.
6. ggml facts that reviews fail on — row-major tensors with dimension 0 as columns; `ggml_mul_mat` is transposed; backend ops must match the CPU reference via `test-backend-ops`; a new `ggml_type` carries its own evidence bar; public headers are an ABI surface and a change there is a design fork.

Backend choice (CUDA / Metal / Vulkan / SYCL), a compiler version, a C++ standard level, and the target repo's Python tooling are not house style.

### 21.7 Version stamp

`tb init` / `tb init --force` writes `toolbox_version = "1.4.0"`. This repo's `schema_version` stays `"1"`.

## 22. Acceptance (v1.4 done when)

On Arch, after `git pull` in `$TOOLBOX_HOME` and a rebuild of `tb`:

1. `tb init -p cpp-systems` in a fresh git repo writes `toolbox.toml` with `profile = "cpp-systems"` and `toolbox_version = "1.4.0"`; `AGENTS.md` has `Profile: cpp-systems`, `review`, `cpp-verify`, and `cpp-ggml` on `Skills:`, no `rust-*` and no `*-go` skill on `Skills:`, and `Verify:` matching the `cpp-verify` recipe.
2. `tb init --force` in a rust-systems fixture, an infra-go fixture, a back-go fixture, and a game-bevy fixture writes `toolbox_version = "1.4.0"`; prose outside the managed-region markers is unchanged; the profile line does not change unless `-p` names a different profile.
3. `tb install` links `cpp-verify` and `cpp-ggml` into `~/.claude/skills` and `~/.grok/skills`.
4. `tb doctor` is still clean aside from a genuinely missing agent (exit 0 with both present); with cwd profile `cpp-systems` it warns on missing `cmake`, not on missing `cargo` or `go`.
5. `skills/cpp-ggml/SKILL.md` exists with `name`/`description` frontmatter and exactly the six house-style headings in §21.6; the first is the rule that the target repo wins and that `tb init` is not run there.
6. `skills/cpp-verify/SKILL.md` names the repo's CI workflows as the source of truth for flags. `review` skill body mentions `cpp-systems`. `profiles/cpp-systems.toml` lists `review` after `handoff` and domain skills `cpp-verify`, `cpp-ggml`.
7. `rust-verify`, `go-verify`, and the four earlier profile files are byte-unchanged by this increment.

M1 is a follow-up, not a v1.4 gate.
