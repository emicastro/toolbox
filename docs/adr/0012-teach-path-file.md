# 0012. A course path is one Markdown note with generated Mermaid

Status: accepted
Date: 2026-09-29

## Context

The `teach` skill stores a course as files in a directory. Claude Code and
Grok Build both read that directory during a subscription session. There is
no API app and no Pi extension between them. The learner wants the tree of
what they will study visible in their Obsidian vault, one level specified
and later levels titled only. Status, prerequisites, sources, and exam dates
have to survive a switch of agent. The view they chose is Mermaid inside a
vault note. A hand-edited diagram kept beside a status table will drift. A
second data file will drift from the note Obsidian has open.

## Options

- **Mermaid as the only record** — status, sources, and dates live in node labels.
- **`path.yaml` plus a generated note** — a structured file is the record; a note is rendered from it.
- **One `PATH.md`** — structured node sections are the record; the Mermaid fence in that same file is regenerated from those sections.

## Decision

One `PATH.md`. Node sections hold kind, status, requires, source, due, and
an optional page. The map fence is rewritten from those sections in the same
edit, and again at the start of every session before teaching. Obsidian
renders the fence. There is no second file to sync.

## Decision mechanics

Node ids are `n1`, `n2`, … assigned once, never reused or renumbered.

Kinds are `concept`, `practice`, and `exam`. A quiz is the chat check that
sets a concept or practice node to `passed`. It is not a node.

Statuses are `locked`, `ready`, `in-progress`, and `passed`. A node is
`ready` only when every id in `requires` is `passed`, or `requires` is
empty. Only one node is `in-progress`.

An open level contains `###` nodes. A sketched level contains a bullet list
of titles and no ids. The map shows a sketched level as one node, with an
edge from the previous level's exam once that exam exists. Expanding a level
replaces the bullet list with `###` nodes and replaces the collapsed map
node in the same edit. The exam's due date is the next calendar day after
the last non-exam node of that level passes.

## Consequences

Both agents have to follow the format file or the map lies. A sketched level
cannot store per-topic prerequisites until it is expanded. Quiz results show
up as node status, not as their own boxes. Splitting the record into two
files, or making a quiz its own node, takes a new ADR.
