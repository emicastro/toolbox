# 0008. `game-bevy` is a fourth profile; this increment is v1.3

Status: accepted
Date: 2026-09-12

## Context

v1 shipped `rust-systems` for libraries and low-level crates. Its house
style fails a Bevy game: clone of `Handle<T>` / `Entity` is treated as a
smell, `unwrap`/`expect` on a broken world invariant is treated as a
library bug, and `Arc<Mutex<_>>` / OS threads are reached for because
the skill does not name the Bevy schedule as the concurrency model.
Miri is not the mismatch — the crate-level trigger already skips a
gameplay crate with no `unsafe`.

A skill-only add while the repo still inits as `rust-systems` leaves
`rust-systems` on the managed region's `Skills:` line; `review` loads it.
ADR 0007 already rejected that shape for HTTP vs CLI Go.

ADR 0005 reserved v2 for a new profile surface; ADR 0007 narrowed v2 to
MCP, `tb verify`, other agents, or a CLI-surface change, and treated a
content pack `LoadProfileByName` already loads as a minor increment.
v1.2 groups 12–16 are accepted. A fourth profile is a new ADR, not a
quiet edit of 0007. Bevy version, physics, net, and UI crates are not
this fork — they stay product ADRs if a game needs to lock them.

## Options

- **Patch `rust-systems`** — add Bevy exceptions to the existing
  house-style skill. Systems crates then inherit game rules; Bevy repos
  still see "no unwrap in libraries" on the `Skills:` line.
- **Generic `game-rust` profile, version v1.3** — fourth profile; the
  engine is a product ADR. House style cannot name ECS, plugins, or the
  schedule, so it cannot correct the clone and concurrency failures.
- **`game-bevy` profile, version v1.3** — lock the engine (Bevy), not a
  Bevy version. House style is ECS / systems / plugins / schedule.
  `tbVersion` / `toolbox_version` become `1.3.0`. `schema_version` stays
  `"1"`. Same CLI, same agents, same parser. `rust-verify` is shared.
- **`game-bevy` profile, version v2** — the same content pack, labelled
  a new product surface because 0005 reserved v2 for a new profile.

## Decision

`game-bevy` profile, version v1.3. This is not a new product: no
subcommand, no agent, no schema change. The engine lock is the house
style; a generic `game-rust` pack cannot name the schedule. 0005 and
0007 are not rewritten; v2 stays reserved for MCP, `tb verify`, other
agents, or a CLI-surface change.

## Consequences

- `docs/requirements.md` keeps all v1, v1.1, and v1.2 text and adds a
  **Delta from v1.2** section. Header becomes v1.3.
- A fourth profile file `profiles/game-bevy.toml`: process skills as of
  v1.1 (`spec`, `adr`, `onboard`, `scout`, `handoff`, `review`), domain
  skills `rust-verify`, `game-bevy`. No `aws-guard`. No `rust-systems`
  on that array. Verify summary is the same Rust recipe as
  `rust-systems`.
- House style lives in `skills/game-bevy/SKILL.md`: five rules (ECS is
  the architecture; panic on broken world invariants, `Result` for I/O;
  the schedule is the concurrency model; handles and components are
  values; do not block the frame). It does not name a Bevy version,
  physics crate, net crate, or UI crate.
- `const tbVersion` becomes `"1.3.0"`. `tb init --force` records it.
  Init usage lists four profile names. Doctor treats `game-bevy` like
  `rust-systems` for the `cargo` PATH warning.
- `rust-systems` remains the systems/library profile. Mutual skip in
  skill descriptions. `rust-verify` is shared; it does not grow a
  second recipe. Miri stays crate-level: isolate renderer/FFI `unsafe`
  in its own crate.
- Remaining v1 non-goals stay non-goals (MCP, `tb verify`, per-repo
  skill add/rm, persona switch, other agents, wrapping extra linters).
  No fifth profile, no Bevy version pin, no `tb` Bevy template.
- Existing product repos pick up v1.3 on `git pull` in `$TOOLBOX_HOME`
  plus `tb install` and, for the version stamp and region, `tb init
  --force`. A new Bevy game inits with `tb init -p game-bevy`.
