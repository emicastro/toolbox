# Toolbox — Design (v1.1)

Status: v1 accepted 2026-09-09; v1.1 delta accepted 2026-09-11
Date: 2026-09-11
Source of truth for requirements: `docs/requirements.md`. This file answers
the "how" for every §14 open point and every command in §8, plus the v1.1
delta in requirements §15. It does not restate rationale already recorded
in `docs/adr/000{1,2,3,4,5,6}-*.md` — those are cited, not re-argued.
v1 sections below stay as accepted; v1.1 is §14 onward.

## 1. Path resolution

`tb` resolves `TOOLBOX_HOME` exactly once per invocation, in `main`, before
dispatching to any subcommand:

1. If the environment variable `TOOLBOX_HOME` is set and non-empty, use it
   verbatim (after `filepath.Clean`).
2. Otherwise use `filepath.Join(os.UserHomeDir(), "toolbox")`.
3. `os.Stat` the resolved path. If it does not exist or is not a directory,
   fail immediately: exit 1, message
   `tb: toolbox home "<path>" not found (from $TOOLBOX_HOME or default ~/toolbox)`.
   No subcommand runs without a valid toolbox home — even `tb doctor`, which
   instead reports this as its first failing check (doctor catches the
   error itself and renders it as a report line; it does not hard-exit
   before printing the rest of the one-screen report — see §5.3).

The resolved value is carried as a single `paths.Toolbox` string, passed
down to every package that needs it (`config`, `scaffold`, `render`,
`agents`). No package re-reads the environment.

## 2. Repo tree

```
$TOOLBOX_HOME/                     # this repo
  AGENTS.md                        # toolbox's own short rules (hand-written)
  toolbox.toml                     # toolbox-side config (schema: §3.2)
  profiles/
    rust-systems.toml
    infra-go.toml
  personas/
    default.md
  skills/
    spec/SKILL.md
    adr/SKILL.md
    onboard/SKILL.md
    scout/SKILL.md
    handoff/SKILL.md
    rust-verify/SKILL.md
    rust-systems/SKILL.md
    go-verify/SKILL.md
    infra-go/SKILL.md
    aws-guard/SKILL.md
  templates/
    AGENTS.md                      # managed-region skeleton, §4
    docs/requirements.md
    docs/design.md
    docs/tasks.md
    docs/session.md
    docs/adr/0000-template.md
  bin/                              # optional user scripts, not installed by tb
  cmd/tb/                           # Go module
    go.mod
    main.go
    internal/
      paths/                        # §1
      config/                       # §3: toolbox.toml + profiles/*.toml reader
      agents/                       # §7: detection
      render/                       # §4: AGENTS.md managed-region rendering
      scaffold/                     # §8.2: docs/ template materialization
      cli/                          # subcommand implementations, §5
```

Every directory's role matches §5 of requirements; nothing renamed. The Go
module lives under `cmd/tb/` specifically so `go build ./...` and `gofmt`
never walk `skills/`, `profiles/`, or `personas/` — those are markdown/TOML
content, not Go source.

## 3. Config schema

### 3.1 Parser subset (shared by both file kinds)

`internal/config` implements one small reader for this closed grammar:

- Line comments starting with `#`, blank lines ignored.
- Bare (unquoted) keys, `=`, then one of:
  - a double-quoted string: `profile = "rust-systems"`
  - a flat array of double-quoted strings, no nesting: `skills = ["spec",
    "adr", "onboard"]`. Like real TOML, newlines inside the brackets are
    insignificant, so the array may span multiple lines up to its closing
    `]` — `profiles/*.toml` (§3.3) relies on this to keep long skill lists
    readable.
- One `[section]` header form, used only by `profiles/*.toml` (`[verify]`
  holds free-form string keys for the verify recipe summary — see 3.3).
- No inline tables, no multi-line strings, no numbers, no dates, no nested
  arrays.

Anything else is a parse error: `<file>:<line>: unsupported syntax: <text>`.
This is intentionally strict — ADR 0001 accepts that `tb` never needs to
read arbitrary TOML, only the subset it itself writes and the profile files
this design defines below.

### 3.2 Product-repo `toolbox.toml`

```toml
# Written by `tb init`. Safe to hand-edit the `profile` line; the other
# two lines are refreshed by --force.
profile = "rust-systems"
toolbox_version = "0.1.0"
generated_by = "tb init"
```

`profile` is the only field `tb` reads back (on every subsequent command
run from that repo, e.g. `tb doctor`). The other two are informational.

### 3.3 `profiles/<name>.toml`

```toml
name = "rust-systems"
description = "Systems / low-level Rust, greenfield and legacy"
persona = "default"                       # personas/<persona>.md
skills = [
  "spec", "adr", "onboard", "scout", "handoff",
  "rust-verify", "rust-systems",
]
templates = ["AGENTS.md", "docs/requirements.md", "docs/design.md",
             "docs/tasks.md", "docs/session.md", "docs/adr/0000-template.md"]

[verify]
summary = "cargo test; cargo clippy -- -D warnings; miri when unsafe changed"
```

`templates` lists paths relative to `templates/`, materialized relative to
the product repo root (§6). `[verify].summary` is a one-line string
rendered into the AGENTS.md managed region (§4) — it is documentation only;
`tb` never executes it (§11, non-goal "wrapping verify as `tb verify`").

### 3.4 Resolution order failure

If `toolbox.toml` in a product repo names a `profile` with no matching
`profiles/<name>.toml` in `$TOOLBOX_HOME`, every command that needs the
profile (`init`, `doctor`) fails with
`tb: profile "<name>" not found in <TOOLBOX_HOME>/profiles`.

## 4. AGENTS.md rendering and merge

Mechanism fixed by ADR 0002. The managed region template
(`templates/AGENTS.md`) is:

```
<!-- toolbox:begin -->
<!-- Generated by tb init. Do not hand-edit between these markers;
     tb init --force replaces this region on every run. -->
Profile: {{.Profile}}
Persona: {{.PersonaPlanTitle}} / {{.PersonaImplementTitle}}
Skills: {{.SkillList}}
Verify: {{.VerifySummary}}
Language: English for code, identifiers, comments, commits, ADRs, and docs.
Chat may be Spanish if the user writes Spanish.
<!-- toolbox:end -->
```

Rendering is a plain `text/template` execution over a struct built from the
loaded profile + persona section titles (parsed from `personas/default.md`
by reading its two `##` headings) + the flattened skill list, joined with
`, `. Output must be deterministic: same profile + same toolbox content ⇒
byte-identical region, so re-running `init` twice with no changes produces
no diff.

Full decision table (mechanics already fixed by ADR 0002, restated here as
the implementation spec `internal/render` follows):

| State on disk | Interactive | Non-interactive | `--force` |
|---|---|---|---|
| No `AGENTS.md` | write full file, no prompt | write full file, no prompt | write full file |
| Exists, no markers | prompt: merge (prepend region) / refuse | refuse, exit 1 | insert region at top, keep rest |
| Exists, markers found | prompt: merge (replace region only) / refuse | refuse, exit 1 | replace region only |

"Refuse" prints the region that *would* have been written to stdout (so the
user can copy it by hand) and exits 1 without touching the file.
Interactivity is detected with stdlib only (ADR 0001 rules out
`golang.org/x/term`, which the original draft of this section
mistakenly named): `os.Stdin.Stat()` and a check of the
`os.ModeCharDevice` bit on the result, with one correction found while
implementing group 4 — `/dev/null` is itself a character device, so a
bare `ModeCharDevice` test misreads `tb init </dev/null` (a common way
scripts and CI mark a command non-interactive) as interactive, which
would print a prompt and then block reading stdin instead of refusing.
`os.SameFile` against `os.DevNull` rules that specific case out. This
is still not a full TTY check — a redirect from some other character
device would still read as interactive — which is an accepted v1 gap.

## 5. Command specs

Global flag: `--profile` / `-p`, parsed by each subcommand's own
`flag.FlagSet` that needs it (`init`; `doctor` accepts it optionally to
override profile detection when `toolbox.toml` is absent, for a
pre-`init` sanity check). Exit codes throughout: `0` success (incl. "one
agent missing" warning), `1` operational failure, `2` usage error (bad
flags, missing required flag).

### 5.1 `tb install`

1. Resolve toolbox home (§1).
2. Detect agents (§7) — **before** step 3 creates anything. ADR 0003's
   config-dir signal is "`$HOME/.claude` exists"; step 3's
   `os.MkdirAll(~/.claude/skills, ...)` creates that parent directory as
   a side effect, so detecting after step 3 would make every install
   self-fulfillingly report both agents present. This ordering bug was
   found while implementing and testing group 4 (docs/tasks.md) and is
   why detection is listed before directory creation here, unlike an
   earlier draft of this section.
3. For each of `~/.claude/skills`, `~/.grok/skills`: `os.MkdirAll` if
   missing.
4. List `$TOOLBOX_HOME/skills/*` (directories containing `SKILL.md`). For
   each skill name, for each of the two target dirs:
   - target does not exist → `os.Symlink(source, target)`.
   - target is a symlink already pointing at `source` → leave it
     (idempotent, no-op).
   - target is a symlink pointing elsewhere → remove and relink.
   - target exists and is a real file/dir (not a symlink) → print
     `warn: <target> exists and is not managed by tb, skipping`, continue
     with the remaining skills (does not abort the run).
5. Using the presence from step 2: if both absent, print
   `error: neither Claude Code nor Grok Build detected`, exit 1 (skills
   were still linked — a later install with an agent present just
   succeeds). If exactly one absent: print
   `warn: <agent> not detected on this machine`, exit 0. If both present:
   proceed silently to the summary.
6. Print a summary: toolbox home, N skills linked per target, agent
   presence for both.

### 5.2 `tb init [-p|--profile <name>] [--force]`

1. Resolve toolbox home. Require `--profile`/`-p`; if absent, exit 2
   (`usage: tb init -p <rust-systems|infra-go>`).
2. Load `profiles/<name>.toml` (§3.4 for the not-found case).
3. Confirm cwd is inside a git repo (`git rev-parse --show-toplevel`); if
   not, exit 1: `tb init must run inside a git repository`.
4. Write/merge `toolbox.toml` at repo root per the same three-state table
   as AGENTS.md (no markers needed — it's a `tb`-owned file with one
   user-editable line, `profile`; interactive/refuse/`--force` behave the
   same as §4 but `--force` rewrites the whole file since only `profile` is
   ever hand-edited and it's preserved by reading the existing value first
   unless `-p` explicitly requests a different profile, in which case the
   new profile always wins).
5. Render/merge `AGENTS.md` per §4.
6. For each path in the profile's `templates` list: if the file does not
   exist at the corresponding repo-relative path, create parent dirs and
   copy the template verbatim. If it exists, leave it — untouched, no
   diffing, no prompt. `docs/adr/` gets its `0000-template.md` copied the
   same way (create-if-missing); no other file under `docs/adr/*` is ever
   written by `init`, satisfying "never overwrite existing ADRs" absolutely
   (not just under `--force`).
7. Print a summary: profile, files created, files left alone, files
   merged/refused.

### 5.3 `tb doctor` (read-only)

One-screen, fixed order:

1. Toolbox home: resolved path, exists Y/N (from §1's own error, caught
   and rendered as a line here rather than propagated).
2. `tb` on PATH: `exec.LookPath("tb")` — informational, never fails doctor.
3. Skill symlink health: for each skill in `$TOOLBOX_HOME/skills`, for each
   of the two target dirs, report `ok` / `missing` / `broken` (symlink
   exists but source doesn't) / `not managed` (real file present).
4. Agent detection (§7): Claude Code present Y/N, Grok Build present Y/N.
5. If cwd has `toolbox.toml`: profile name, whether `profiles/<name>.toml`
   resolves, whether `AGENTS.md` exists and has a managed region.
6. If the profile (from 5, or `-p` override) is `rust-systems`, warn if
   `cargo` is not on PATH; if `infra-go`, warn if `go` is not on PATH.

Exit code: `0` unless step 4 finds *zero* agents present, in which case
`1` (mirrors `install`'s fail-closed rule, §7). All other findings are
warnings only — doctor never fails on a missing skill link or a missing
`cargo`/`go`. No `--fix` (explicitly out of v1, §8.6).

### 5.4 `tb skill new <name>`

1. Resolve toolbox home. Validate `<name>` matches `^[a-z][a-z0-9-]*$`;
   else exit 2.
2. If `$TOOLBOX_HOME/skills/<name>/SKILL.md` already exists, exit 1
   (`already exists`).
3. `os.MkdirAll($TOOLBOX_HOME/skills/<name>)`, write `SKILL.md` from a
   built-in template (§8):
   ```
   ---
   name: <name>
   description: TODO — one line, when an agent should reach for this skill.
   ---

   # <name>

   TODO: body.
   ```
4. Print next steps verbatim (matches §8.4): edit the skill, add its name
   to the chosen profile's `skills` array, run `tb install`, commit in
   `$TOOLBOX_HOME`. `tb skill new` does not touch any profile file itself.

### 5.5 `tb persona show [-p|--profile <name>]`

1. Resolve toolbox home. If `-p` given, load that profile; else read
   `toolbox.toml` in cwd for `profile`; else exit 2 asking for one or the
   other.
2. Load the profile's `persona` field, read `personas/<persona>.md`, parse
   its two `##` section headings.
3. Print: persona file path, both section titles. Read-only, no mutation.

## 6. Symlink strategy for `install`

Covered inline in §5.1 step 3. The key property: **idempotent and
non-destructive to unmanaged files, best-effort across skills** — one
skill's conflict never blocks the others from linking.

## 7. Agent detection and failure policy

Detection heuristic and lookup table fixed by ADR 0004 (superseding ADR
0003): `exec.LookPath` OR a marker file inside the agent's config
directory, per agent — `~/.claude/settings.json` for Claude Code,
`~/.grok/config.toml` for Grok Build. ADR 0003's original "config dir
exists" signal is not used: `tb install` itself creates
`~/.claude/skills`, so the bare directory existing is not a reliable
presence signal (ADR 0004 has the story). Failure policy (fixed by
requirements §8.1/§8.3, restated for implementers):

- Both present → success, exit 0.
- Exactly one present → warning, exit 0.
- Neither present → error, exit 1 (`install`) / doctor's only hard failure
  (`doctor`).

`internal/agents` exposes one function,
`Detect() (claude, grok bool)`, used identically by `install` and `doctor`
so the two commands can never disagree (ADR 0004 consequence).

## 8. Profile content (v1, exact)

Both profiles carry the five mandatory process skills unchanged (§6.2):
`spec`, `adr`, `onboard`, `scout`, `handoff`.

`rust-systems` adds: `rust-verify`, `rust-systems`.
`infra-go` adds: `go-verify`, `infra-go`, `aws-guard`.

No profile gets an skill beyond this list — task group 2 writes exactly
these ten `SKILL.md` files, nothing more, matching the "Plan Mode may add a
small number of extra domain skills if load-bearing" ceiling by adding
zero, since all ten are already named in requirements §6.2/§6.3.

## 9. Persona

`personas/default.md`:

```markdown
# Default persona

## Plan — Staff Engineer

Challenges weak specs. Demands an ADR before treating a design fork as
settled. Does not implement past the approved task list. Prefers evidence
(`cargo`/`go` command output) over narration.

## Implement — Senior Engineer

Follows accepted ADRs and `docs/tasks.md`. Does not reopen design without a
new ADR. Keeps diffs small. Runs the profile's verify recipe before
claiming a task done.

---

Model tier (Opus, Grok 4.6, Sonnet, ...) does not change the role above. A
cheaper model on Implement still follows the Implement section in full.

Language: default English for chat and all artifacts. Replies may be
Spanish if the user writes Spanish; code, identifiers, comments, commit
messages, ADRs, and docs stay English regardless.
```

`internal/render` parses this file generically (first two `##` headings
become the section titles used in §4 and §5.5) so a future second persona
file works without code changes — though v1 ships only `default.md`.

## 10. Skill authoring contract

Every `SKILL.md` under `skills/<name>/`:

```markdown
---
name: <name>
description: One line — when an agent should reach for this skill.
---

# <name>

<body: what the skill does, and for the two `*-verify` skills, the exact
verify commands, e.g. `cargo test`, `cargo clippy -- -D warnings`,
`cargo miri test` when the diff touches `unsafe`>
```

Frontmatter is exactly `name` + `description` (ADR-level minimal common
denominator, §14 answer) — nothing Claude-Code-specific or
Grok-Build-specific, since `tb install` symlinks the same file into both
tools' skill directories (§12: "standard `SKILL.md` so a future agent can
be added by teaching `tb install` a third symlink path"). Verify recipes
live only inside `rust-verify` / `go-verify` bodies — `tb` itself never
shells out to `cargo` or `go` (§11 non-goal).

## 11. Cross-platform notes

- Build targets: `GOOS=linux GOARCH=amd64` and `GOOS=darwin GOARCH=arm64`,
  both `CGO_ENABLED=0 go build -trimpath -o dist/tb-<os>-<arch> ./cmd/tb`.
- `os.Symlink` / `os.Readlink` behave identically on both targets; no
  Windows-only or POSIX-only path assumed beyond what `path/filepath`
  already normalizes.
- `os.UserHomeDir()` resolves correctly on both (`$HOME` on Linux, `$HOME`
  on macOS).
- No platform-conditional code paths anticipated for v1; if one turns out
  to be needed (e.g. macOS Gatekeeper on the built binary), that's a task
  in group 6 (acceptance), not a design change.

## 12. Testing strategy

- `internal/config`: table tests over the TOML subset (valid inputs, and
  each unsupported-syntax case producing the expected `file:line` error).
- `internal/paths`: table tests over `TOOLBOX_HOME` set/unset/missing-dir.
- `internal/render`: golden-file tests — given a fixed profile + persona
  fixture, the rendered managed region matches a checked-in expected
  string; a second test asserts merge preserves content outside markers.
- `internal/agents`: detection logic tested by injecting a fake `PATH` and
  fake home dir (no real dependency on `claude`/`grok` being installed in
  CI).
- `internal/cli` (per command): each test runs inside `t.TempDir()` acting
  as both a fake `$TOOLBOX_HOME` and a fake product repo (`git init`'d in
  the temp dir); no test writes outside `t.TempDir()`, no test needs
  network access.
- Acceptance (task group 6) is run manually on Arch and M1 per §13 of
  requirements — not part of `go test ./...`.

## 13. Traceability: §13 acceptance → design mechanism

| §13 item | Mechanism |
|---|---|
| 1. install links skills; doctor clean aside from optional agent | §5.1, §5.3, §7 |
| 2/3. init produces pointer + AGENTS.md + docs/, doesn't destroy ADRs | §5.2 steps 4–6 |
| 4. re-init without `--force` prompts merge/refuse | §4 table, §5.2 step 4 |
| 5. `--force` refreshes only AGENTS.md + toolbox.toml | §4, §5.2 step 4 |
| 6. `skill new demo` scaffolds under `$TOOLBOX_HOME/skills/demo` | §5.4 |
| 7. both agents see symlinked skills | §5.1, §6 |
| 8. doctor fails closed on zero agents, warns on one | §5.3, §7 |

No item in §13 lacks a mechanism above.

## 14. v1.1 quality contract

Requirements §15. Versioning ADR 0005; region payload ADR 0006. No new
`tb` subcommands. `schema_version` stays `"1"`. `tbVersion` becomes
`"1.1.0"` (the informational `toolbox_version` in product `toolbox.toml`).

### 14.1 Managed region payload (ADR 0006)

§4's template stays the mechanism (markers, `text/template`, merge table).
The file `templates/AGENTS.md` grows by a **Rules** block after the
`Language:` line, still inside the markers, no new template variables.
Wording is locked in ADR 0006. `personas/default.md` and toolbox's own
`AGENTS.md` carry the same five lines.

`internal/render.TemplateData` is unchanged. Tests that render the real
template must assert the rules block is present.

### 14.2 Process skill `review`

Both profiles' `skills` arrays become:

```
spec, adr, onboard, scout, handoff, review,
<domain skills unchanged>
```

`review` sits after `handoff`. `tb install` already links every
`skills/*/SKILL.md`; the array is what the region's `Skills:` line shows.

Skill contract (body in `skills/review/SKILL.md`): last gate after verify,
before ticking a task. Skip only when the session produced no diff. Pass:
scope, ADRs, profile house style (load `rust-systems` / `infra-go` /
`aws-guard`, do not restate them), verify actually run, `Check:` can fail.
Output: findings with `file:line`, or `no findings` naming files read.
No refactor during review.

Persona Implement names this skill; it does not duplicate house style.

### 14.3 Verify recipes (requirements §15.4)

`profiles/rust-systems.toml` `[verify].summary`:

```
cargo fmt --check; cargo test; cargo clippy --all-targets -- -D warnings; miri when the changed crate has unsafe
```

`profiles/infra-go.toml` `[verify].summary`:

```
gofmt -l .; go vet ./...; go test ./...; go test -race ./...
```

Skill bodies list the same commands. Miri: run `cargo miri test` when a
changed crate contains any `unsafe` token in `.rs` sources, even if this
diff is only safe code. Skip Miri only when the changed package(s) have
no `unsafe`. Missing Miri while required fails the recipe.

§10's example body still says recipes live in the `*-verify` skills; the
v1 phrase "when the diff touches `unsafe`" is superseded by the crate-level
trigger above for v1.1.

### 14.4 Skill skip and house-style deltas

Documented in the skill files (group 8), matching requirements §15.3 and
§15.5. `onboard` in particular: `AGENTS.md` existing is not a map.

## 15. Traceability: v1.1 acceptance → design mechanism

| Requirements §16 item | Mechanism |
|---|---|
| 1. `--force` refreshes rules + `review` on Skills; version 1.1.0; outside markers intact | §14.1, §14.2, `tbVersion`, ADR 0002 merge |
| 2. `tb install` links `review` | §5.1, `skills/review/` |
| 3. `tb doctor` still exit 0 with both agents | §5.3, §7 unchanged |
| 4. onboard does not skip on mere `AGENTS.md` | `skills/onboard/SKILL.md` |
| 5. verify command blocks match §15.4 | §14.3, `*-verify` skills, profile summaries |
| 6. `review` skill + both profiles list it after `handoff` | §14.2 |

No item in requirements §16 lacks a mechanism above.
