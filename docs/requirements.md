# Toolbox — Requirements (v1)

Status: accepted  
Date: 2026-09-09  
Command: `tb`  
Source of truth for this product: this file. Design and tasks are produced later in Plan Mode, in the `toolbox` repo.

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
