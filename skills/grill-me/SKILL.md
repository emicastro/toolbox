---
name: grill-me
description: Use when a plan, decision, or idea is still too vague to write requirements, or the user asks to be grilled / stress-tested; skip once docs/requirements.md and docs/tasks.md already match the ask (the spec skill owns that), and skip when the answers are not the user's (use to-questionnaire).
---

# grill-me

Interview the user until you reach a shared understanding. Do not write
`docs/requirements.md`, `docs/design.md`, or ADRs here — that is the
`spec` / `adr` sequence after the frontier is empty.

Map the ask as a **design tree**: every decision branches into the
decisions that hang off it. Work the tree in **rounds**. The **frontier**
is every decision whose prerequisites are already settled: the questions
you can ask *now* without guessing at answers you have not heard yet.
Ask the whole frontier in one round; number each question and give a
recommended answer; then wait.

```
❓ **Q1** - **<title>**: <body, including choices if useful>

➡️ <recommended answer>
```

A question whose answer depends on another question still open in this
round belongs to a later round. Settled decisions push the frontier
outward; recompute and ask again.

Finding *facts* is your job, never the user's. Use the `scout` skill (or
a cheaper subagent) for anything in the environment; don't block the rest
of the frontier on it. The *decisions* are the user's: put each to them
and wait.

The session is done when the frontier is empty and the user confirms
shared understanding. Then offer the `spec` skill (requirements → design
→ tasks). Forks with two live options go to the `adr` skill, not into
chat.
