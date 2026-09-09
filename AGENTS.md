# Toolbox — agent rules

This repo is `$TOOLBOX_HOME` itself, not a product repo, so it has no
generated managed region and no profile.

Source of truth, in order:

- `docs/requirements.md` — what v1 must do and must not do.
- `docs/design.md` — how, including every open point requirements left for
  Plan Mode.
- `docs/tasks.md` — the ordered work list; each task ends in a check.
- `docs/adr/*.md` — recorded decisions on load-bearing forks. A new
  architectural fork needs a new ADR before it's treated as settled; do not
  reopen a decision already recorded here without one.

Rules:

- Do not implement past the approved task list in `docs/tasks.md`.
- Keep diffs small; one task group per commit at minimum.
- Docs and code in English. Chat may reply in Spanish if the user writes
  Spanish, but code, identifiers, comments, commit messages, ADRs, and docs
  stay English.
