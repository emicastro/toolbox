---
name: handoff
description: Use when a session is ending, about to compact, or being paused mid-task and the next session needs the thread — write `docs/session.md` and stop; skip while work is in progress with context to spare, and never use it as a substitute for finishing a task.
---

# handoff

Session bridge. Write `docs/session.md` in the product repo, then stop.
Overwriting the previous contents is expected — this file is the *current*
bridge, not a log.

## What goes in `docs/session.md`

- **State** — which task from `docs/tasks.md` is in progress, by number, and
  which checkboxes were ticked this session.
- **Done** — what actually changed on disk, with paths. Committed or not.
- **Verify status** — the profile's verify recipe: run, passing, failing
  with which output. "Not run" is a valid and useful answer.
- **Next** — the single next action, concrete enough to start cold.
- **Open questions / blockers** — decisions the next session must not
  silently make. Anything that is a design fork goes to an ADR
  (`docs/adr/NNNN-title.md`), not into this file.

## Rules

- Durable knowledge belongs in `docs/` and `docs/adr/`, user-wide
  conventions in `$TOOLBOX_HOME`. `docs/session.md` is transient by design —
  do not park decisions here.
- An uncommitted chat summary is never the source of truth. If it matters,
  it is in a file.
- Leave the tree in a state the next session can read: no half-applied
  rename, no file written to a path the task list does not mention.
- Then actually stop. Handoff is a boundary, not a preamble to more work.
