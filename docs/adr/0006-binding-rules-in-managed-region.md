# 0006. Managed region includes a short binding-rules block

Status: accepted
Date: 2026-09-11

## Context

ADR 0002 decided *how* `AGENTS.md` is owned (marker-delimited region,
inline render, merge/`--force` table). It did not decide *what* the region
contains beyond profile, persona titles, skill list, verify summary, and
the language rule.

After `tb init`, a product repo's `AGENTS.md` therefore has job titles and
a skill index, not the rules. The persona body lives only in
`$TOOLBOX_HOME/personas/default.md`. Skills are global and each frontmatter
says when *not* to use them. An agent that does not open those files has
no in-repo contract. That is the main quality loophole of v1.

The region must stay short (v1 wanted a short `AGENTS.md`). The marker
mechanism is not up for debate (ADR 0002).

## Options

- **Titles only** — keep the v1 payload. Skills and persona remain the
  only place the rules live. No `tb init --force` churn.
- **Short rules block inside the markers** — five English lines, no new
  template variables, copied also into `personas/default.md` and toolbox's
  own hand-written `AGENTS.md`.
- **Full persona body in every product `AGENTS.md`** — the whole
  Plan/Implement text inline. Strongest contract; the file is no longer
  short, and `--force` rewrites a large region every time.

## Decision

Short rules block inside the markers. Titles-only leaves the loophole
open. The full persona body fights v1's "short `AGENTS.md`" and duplicates
a file that `tb persona show` already prints.

## Decision mechanics

The block is static English in `templates/AGENTS.md`, after the existing
`Verify:` / `Language:` lines, still between `<!-- toolbox:begin -->` and
`<!-- toolbox:end -->`. Wording to lock:

1. Do not implement past the accepted task list in `docs/tasks.md`.
2. A design fork (two viable options, a dependency, a schema/protocol
   shape, or reversing an ADR) needs a new ADR before it is treated as
   settled.
3. Before ticking a task or claiming done: run the profile verify recipe
   (success is the command output), then the `review` skill.
4. Keep diffs small. No drive-by refactors or unrelated formatting.
5. Do not assume profile defaults in a brownfield repo; map it first
   (`onboard`). `AGENTS.md` existing is not a map.

`internal/render` gains no new template fields. Merge and `--force`
behaviour is unchanged (ADR 0002).

## Consequences

- `tb init --force` on an existing product repo refreshes the region;
  prose outside the markers survives byte-for-byte.
- Hand-edits *inside* the markers still get overwritten (already an ADR
  0002 consequence).
- Persona and toolbox `AGENTS.md` must carry the same five lines, or the
  three sources drift.
- Tests that render the real template must assert the rules block is
  present, not only that markers exist.
