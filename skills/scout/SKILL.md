---
name: scout
description: Use when the next step is locating things — which file defines X, where Y is called, does this repo already have Z — anything a handful of `rg`/`find`/read calls answers; skip when the answer needs judgment over the results rather than the results themselves.
---

# scout

Cheap exploration. The rule is blunt: **do not spend a frontier model on
work equivalent to `rg`.** Search first, read second, reason last.

## How

- Grep before reading: `rg -n 'symbol'`, `rg -l`, `rg -t rust`, `rg -t go`.
  Narrow with a path prefix instead of widening the pattern.
- `find` / `ls` for shape questions ("is there a `docs/adr/`?"), not a full
  tree walk.
- Read the specific range you need (`Read` with an offset), not whole files,
  once search has told you where to look.
- Prefer three targeted searches over one broad one whose output you then
  skim; skimming a large dump is where the budget actually goes.
- Delegate a noisy fan-out sweep to a cheaper model or a subagent and keep
  only the findings, when the harness offers that.

## Stop conditions

- You have the file:line answers the task needs → stop searching, start
  working.
- Three searches found nothing → the name is wrong, not the repo. Ask, or
  widen along a different axis (a string in an error message, a config key,
  a test fixture) rather than repeating the same query.
- The question turned out to be "how does this system work", not "where is
  it" → switch to the `onboard` skill.

Report locations as absolute paths with line numbers. Never claim a symbol
does not exist on the strength of one pattern.
