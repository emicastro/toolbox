---
name: spec
description: Use in Plan Mode when work starts from an idea or a rough requirements note and `docs/design.md` + `docs/tasks.md` do not exist yet or no longer match the ask; skip once a task list is accepted and the job is executing it.
---

# spec

Turn an idea into three accepted documents before any implementation
starts. Run under the persona's **Plan — Staff Engineer** section: challenge
the spec, do not implement past the list you are producing.

## Sequence

1. **`docs/requirements.md`** — the source of truth for *what* and *why*.
   Problem, goals, explicit non-goals, users/machines, acceptance criteria
   the work is done against. No implementation choices. Push back on vague
   goals here; a requirement you cannot write an acceptance line for is not
   a requirement yet. If the ask is still a tree of unsettled decisions,
   switch to the `grill-me` skill, then resume this sequence from step 1.
2. **`docs/design.md`** — answers *how* for every open point requirements
   left, section by section. It cites decisions, it does not re-argue them:
   any fork with two live options goes to an ADR first (see the `adr`
   skill), and design.md references `docs/adr/NNNN-title.md` instead of
   repeating the rationale.
3. **`docs/tasks.md`** — numbered groups, each task scoped to one Implement
   session and ending in a runnable check. The `Check:` line must be a
   command whose non-zero exit fails the task, or a byte-level before/after
   assertion. "File exists" is valid only when the whole task is creating
   that file. Checkboxes `- [ ]`, ticked only by the session that finishes
   the task.

## Rules

- Stop at the task list. Producing code in the same pass defeats the point.
- Every acceptance item in requirements must map to a mechanism in design;
  add a traceability table if there is more than a handful.
- Non-goals are load-bearing — write them down, and treat later scope creep
  as needing a new ADR, not a quiet edit.
- End the session by getting the three files accepted (Status/Date header),
  then hand off with the `handoff` skill.

This repo's own `docs/requirements.md`, `docs/design.md`, `docs/tasks.md`
and `docs/adr/000{1,2,3}-*.md` are a worked example of the output shape.
