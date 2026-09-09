---
name: adr
description: Use before treating any design fork as settled — two or more viable options, a dependency or library choice, a file/schema/protocol shape, or a reversal of an earlier decision; skip for implementation details already covered by an accepted ADR.
---

# adr

An Architecture Decision Record is **mandatory** in both profiles before a
fork is treated as decided. If you catch yourself writing "we'll use X"
in `docs/design.md` and X had a real alternative, stop and write the ADR
first.

## Shape

New file `docs/adr/NNNN-title.md`, `NNNN` the next zero-padded number,
`title` a lowercase hyphenated slug (`0004-agents-md-merge.md`). Copy
`docs/adr/0000-template.md`. MADR-short, four sections plus a header:

- Header: `Status: proposed | accepted | superseded by 000X` and `Date:`.
- **Context** — the forces, and the problem that forced the fork. Only the
  facts that make the options legible.
- **Options** — at least two, one line each, named. An ADR with one option
  is a note, not a decision.
- **Decision** — which option, and the one or two sentences that tip it.
- **Consequences** — what gets easier, what gets harder, what this commits
  the project to or forecloses.

A decision with non-obvious mechanics may add a short "Decision mechanics"
section between Decision and Consequences (see `0002` in this repo).

## Rules

- Superseding is a **new** ADR; the old one gets `Status: superseded by
  000X` and keeps its body. Never rewrite an ADR's history.
- `tb init` creates `docs/adr/` and the template, and never overwrites an
  existing ADR — not even under `--force`. Treat the directory as append-only.
- `docs/design.md` cites ADRs; it does not duplicate their rationale.
- Write ADRs in English regardless of chat language.
