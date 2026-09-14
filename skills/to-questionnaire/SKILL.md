---
name: to-questionnaire
description: Use when a decision needs answers from someone who is not in this chat; skip when the user can answer (use grill-me) or when spec/adr already cover the fork.
---

# to-questionnaire

Turn something the user cannot answer alone into a **questionnaire**: a
Markdown document they hand to one person to fill in async, or fill out
together in a meeting. The recipient holds knowledge the user lacks.

**Grill the send, not the subject.** Interview the user only about the
*send* (they can always answer that): who it goes to, and what they need
back. The questions in the document then target the **gap** between what
the recipient knows and what the user needs.

1. **Who is it going to?** One exchange: role, expertise, relationship
   to the user. Done when you know who the recipient is and what they
   know that the user does not.
2. **What do you need back?** One exchange: the decisions or facts the
   user cannot resolve alone. Done when you have a concrete list of what
   the user must walk away able to do or decide.
3. **Write the questionnaire.** Follow the structure below. Write it to
   `docs/<slug>-questionnaire.md` (create `docs/` if needed; slug from
   the topic) and report the path. Done when the file exists and every
   item from step 2 is covered.

## Document structure

Discovery questionnaire: the user lacks context, the recipient holds it.
Order questions most-important-first (async may get one pass). Group
under `##` headings by theme once there are more than a handful.

```md
# <Questionnaire title>

**Purpose:** why this exists and the decision riding on it.

**From:** <the user>, **To:** <the recipient>, **How your answers will be used:** <where they go>

## Context

One paragraph orienting a recipient who was not in the user's head.
Enough to answer well, not a page.

## How to answer

Deadline and rough effort. Partial answers and "I don't know" are
useful: flag anything unsure rather than skipping it.

## <Theme heading>

### What load is the system expected to handle at launch?

_Why this matters: it decides whether we provision for burst traffic now or defer it._

**Answer:**

## Anything else?

Anything we did not ask that we should know?
```

Every question is one idea, never compound, with an answer stub directly
beneath. A one-line *why this matters* only where the question could be
misread or invite a throwaway answer.
