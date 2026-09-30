# 0013. The terminal session is the lesson

Status: accepted
Date: 2026-09-29

## Context

`skills/teach` makes a self-contained HTML file the primary lesson, one per
thing taught. The learner reviews in Obsidian and is taught in a Claude Code
or Grok Build session. Writing that session out again as HTML spends the
turn on a page they may not reopen. Writing it out again as a full vault
note duplicates the same text. The vault already renders the path, the
glossary, the sources, and the reference sheets.

## Options

- **HTML lesson per node** — the current skill. The chat is a draft of that file.
- **Full Markdown note per node** — the vault holds the explanation, and HTML is extra.
- **The session is the lesson** — the chat teaches one node, a learning record stores what stuck, and HTML is written only as a practice page linked from the node.

## Decision

The session is the lesson. A practice page (HTML under the course's
`lessons/`, built from `assets/`) is written only when manipulating a figure
teaches the node, or when the node is practice that needs a small in-browser
task. The page does not write status. Learning records, reference sheets,
the glossary, sources, and `PATH.md` are what the learner rereads.

## Consequences

The verbatim session is not in the vault. A short learning record cannot
reconstruct the dialogue. Review of a node goes through its source, its
record, and its practice page when one exists. Replacing this with a full
note per node, or restoring HTML as the primary lesson, takes a new ADR.
Mission, sources, glossary, reference sheets, and learning records stay as
they are.
