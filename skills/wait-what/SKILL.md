---
name: wait-what
description: Use when the user says they didn't follow, asks to re-pitch, or the last explanation clearly didn't land; skip when they asked a specific question rather than a reset.
---

# wait-what

Stop. Re-pitch the current point: a few sentences of context, then the
claim. Short sentences, one idea each.

Use the terms already in this repo (`docs/design.md`, accepted ADRs,
`AGENTS.md`, `docs/GLOSSARY.md`). Do not invent a parallel vocabulary.
If the miss was a term, add it with the `glossary` skill, then re-pitch
with that word. Do not continue the previous plan until the user
confirms the re-pitch landed.
