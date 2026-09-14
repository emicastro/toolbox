---
name: glossary
description: Use when the user runs /glossary, when onboard/spec/grill-me/wait-what recognizes a business domain or a term the user does not master, or when chat would otherwise re-explain a project word; skip for general programming vocabulary and in a teaching workspace (the teach skill owns GLOSSARY.md there).
argument-hint: "[term to add]"
---

# glossary

`docs/GLOSSARY.md` is the ubiquitous language for this product repo.
**Business-domain terms first** (the words of the field, not of the
stack): TTFT, quantization, p99 in inference engineering; FMIS, NDVI,
VRA in AgTech. Code identifiers only when they *are* domain names in
this project. Prefer these terms in chat; do not re-explain them.

This is not the `teach` glossary (that one records terms after the user
already understands them). This file exists so the user can keep it
open beside the editor while still learning the field.

## /glossary

- **File exists, no term named:** open `docs/GLOSSARY.md` and stop. Do
  not rewrite it.
- **File exists, user named a term:** add or update that term, then open.
- **File missing:** create it (seed from the recognized domain and this
  repo), then open.

Open means: Read the file so the terms are in context, then make it
visible (`xdg-open`, `open`, or `$EDITOR docs/GLOSSARY.md`). Report the
path.

## When onboard / spec / grill-me / wait-what call this

If a business domain is recognizable from the repo, the ask, or the
conversation, add the load-bearing terms a newcomer in *that field*
would need beside this codebase. Cap a first seed at about 15. Skip
"function", "HTTP", "struct", and anything already in an ADR (cite the
ADR instead).

`wait-what`: if the miss was a term, that term is the entry; add it,
then re-pitch with that word.

Create `docs/` if needed. Update in place; do not leave stale entries.

## Format

```md
# Glossary

Ubiquitous language for this repo. Prefer these terms in chat.

## Terms

**TTFT (time to first token)**:
Latency from request arrival to the first generated token. This project's
SLO is on TTFT, not full-response latency.
_Avoid_: "startup time", "time to first byte"

**NDVI**:
Normalized difference vegetation index; a satellite/drone proxy for
canopy vigor used in this project's field maps.
_Avoid_: "greenness", "crop health score"
```

One or two sentences: what the term *is in this project*. `_Avoid_` lists
aliases the agent must not use. Group under extra `##` headings only when
clusters appear. A fork is still an ADR, not a glossary paragraph.
