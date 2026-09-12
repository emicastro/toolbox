# Toolbox — agent rules

This repo is `$TOOLBOX_HOME` itself, not a product repo, so it has no
generated managed region and no profile.

Source of truth, in order:

- `docs/requirements.md` — what v1.1 must do and must not do.
- `docs/design.md` — how, including every open point requirements left for
  Plan Mode.
- `docs/tasks.md` — the ordered work list; each task ends in a check.
- `docs/adr/*.md` — recorded decisions on load-bearing forks. A new
  architectural fork needs a new ADR before it is treated as settled; do not
  reopen a decision already recorded here without one.

Binding rules (ADR 0006, same wording as the product `AGENTS.md` managed
region):

- Do not implement past the accepted task list in `docs/tasks.md`.
- A design fork (two viable options, a dependency, a schema/protocol shape, or reversing an ADR) needs a new ADR before it is treated as settled.
- Before ticking a task or claiming done: run the profile verify recipe (success is the command output), then the `review` skill.
- Keep diffs small. No drive-by refactors or unrelated formatting.
- Do not assume profile defaults in a brownfield repo; map it first (`onboard`). `AGENTS.md` existing is not a map.

This repo additionally: one task group per commit at minimum. Docs and
code in English. Chat may reply in Spanish if the user writes Spanish,
but code, identifiers, comments, commit messages, ADRs, and docs stay
English.
