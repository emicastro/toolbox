# 0005. This increment is v1.1, not v2

Status: accepted
Date: 2026-09-11

## Context

v1 of toolbox shipped: one clone, `tb`, two profiles, two agents, five
process skills, and the v1 non-goals (no MCP, no `tb verify`, no extra
profiles or agents). Acceptance group 6 is done.

A review of `AGENTS.md` and the skills found that the v1 quality bar is
optional in practice: the product-repo managed region is titles-only, skill
`skip` clauses are easy to take, and verify recipes are weaker than CI.
Closing those loopholes is new work. It needs a version label so product
repos' `toolbox_version` and this repo's docs do not pretend it is still
the v1 snapshot, and so v2 stays reserved for a different product surface.

`tb` already writes `toolbox_version` into each product `toolbox.toml`
(informational, docs/design.md §3.2). Today that string is `0.1.0`.

## Options

- **v1.1 (minor)** — same product; tighten the quality contract of content
  v1 already ships (skills, persona, managed-region payload, verify
  recipes). `toolbox_version` becomes `1.1.0`. `schema_version` stays `1`.
- **v2 (major)** — treat this as a new product: new commands, profiles,
  agents, or dropping a v1 non-goal.
- **Silent patch of v1** — edit skills and templates in place, leave
  version strings and the requirements header at v1.

## Decision

v1.1. The work does not change what toolbox *is*; it makes the v1
practices binding. v2 is reserved for non-goals becoming goals or a new
command/profile/agent surface.

## Consequences

- `docs/requirements.md` keeps all v1 text (do not rewrite history) and
  adds a **Delta from v1** section. Header becomes v1.1.
- Requirements §6.2 is amended in that delta: six process skills, adding
  `review` as the last gate (after `handoff`). This is a process skill, not
  a domain skill, and not a new `tb` subcommand.
- `const tbVersion` in `cmd/tb/internal/cli/init.go` becomes `"1.1.0"`.
  `tb init --force` records it. This repo's `toolbox.toml` `schema_version`
  stays `"1"`.
- v1 non-goals stay non-goals. Bringing any of them in is a new ADR, not a
  quiet edit of the delta.
- Existing product repos pick up v1.1 content on `git pull` in
  `$TOOLBOX_HOME` plus `tb install` and, for the managed region,
  `tb init --force`.
