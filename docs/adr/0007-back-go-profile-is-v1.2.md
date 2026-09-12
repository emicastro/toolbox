# 0007. `back-go` is a third profile; this increment is v1.2

Status: accepted
Date: 2026-09-12

## Context

v1 named `back-go` (HTTP services, databases) and deferred it: a non-goal
in requirements §3, "deferred" in §6.1, and `infra-go` explicitly skips
the listen loop. ADR 0005 shipped v1.1 as a quality-contract tightening
and reserved v2 for "non-goals becoming goals or a new command/profile/
agent surface". Bringing a non-goal in is a new ADR, not a quiet edit of
the v1.1 delta (0005 consequence).

A product repo that is a Go HTTP backend now needs `tb init` to wire a
house style. `LoadProfileByName` already loads any `profiles/<name>.toml`;
the coupling is the hardcoded init usage string, doctor's toolchain
switch, and the two shipped profile files.

`go-verify` already documents the Go recipe. The missing piece is a
profile pack plus a house-style skill. Router, ORM, and migrator stay
product ADRs — not a toolbox lock.

## Options

- **Skill only, stay on v1.1** — add `skills/back-go/SKILL.md`; `tb
  install` links it. The product still inits as `infra-go`, whose house
  style says skip for HTTP. The managed region's `Skills:` line does not
  list `back-go`.
- **Third profile, version v1.2** — add `profiles/back-go.toml` and the
  house-style skill. `tbVersion` / `toolbox_version` become `1.2.0`.
  `schema_version` stays `"1"`. Same CLI, same agents, same parser.
- **Third profile, version v2** — the same content pack, labelled a new
  product surface because 0005 reserved v2 for dropping a v1 non-goal.

## Decision

Third profile, version v1.2. This is not a new product: no subcommand, no
agent, no schema change. `back-go` was named and deferred in v1, not
rejected. 0005 is not rewritten; v2 stays reserved for MCP, `tb verify`,
other agents, or a CLI-surface change. A content pack the existing init
path already loads is a minor increment.

## Consequences

- `docs/requirements.md` keeps all v1 and v1.1 text and adds a **Delta
  from v1.1** section. Header becomes v1.2. The §3 non-goal listing
  `back-go` stays as history; the delta promotes it.
- A third profile file `profiles/back-go.toml`: process skills as of
  v1.1 (`spec`, `adr`, `onboard`, `scout`, `handoff`, `review`), domain
  skills `go-verify`, `back-go`, `aws-guard`. Verify summary is the same
  Go recipe as `infra-go`.
- House style lives in `skills/back-go/SKILL.md`: five rules (request
  context, graceful shutdown, errors at the HTTP edge, structured logs
  with no secrets, parameterized SQL and explicit migrations). It does
  not name a router, ORM, or migrator.
- `const tbVersion` becomes `"1.2.0"`. `tb init --force` records it.
  Init usage lists three profile names. Doctor treats `back-go` like
  `infra-go` for the `go` PATH warning.
- `infra-go` remains the CLI/jobs profile. Mutual skip in skill
  descriptions. `go-verify` is shared; it does not grow a second recipe.
- Remaining v1 non-goals stay non-goals (MCP, `tb verify`, per-repo skill
  add/rm, persona switch, other agents, wrapping extra linters).
- Existing product repos pick up v1.2 on `git pull` in `$TOOLBOX_HOME`
  plus `tb install` and, for the version stamp and region, `tb init
  --force`. A new HTTP backend inits with `tb init -p back-go`.
