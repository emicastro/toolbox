# Glossary

Ubiquitous language for the teach-path delta (requirements §27, ADR 0012,
ADR 0013). Prefer these terms in chat. Harness words (profile, managed
region, verify recipe) stay defined in `docs/requirements.md` and
`docs/adr/`.

## Terms

**Course**:
One topic's directory, `<vault>/Learn/<slug>/`, holding the mission, the path, sources, and records.
_Avoid_: "workspace", "repo", "class"

**Vault**:
The Obsidian folder that contains `.obsidian`. Courses for this user live under its `Learn/` directory.
_Avoid_: "notes app", "second brain"

**Mission**:
The reason this course exists, written in `MISSION.md`. Lessons and the path are chosen against it.
_Avoid_: "goal statement", "learning objective"

**Path**:
`PATH.md` in the course. The node sections are the record; the Mermaid fence is the tree the learner looks at.
_Avoid_: "curriculum", "roadmap", "syllabus"

**Node**:
One concept, practice, or exam on the path, with a stable id (`n1`, `n2`, …).
_Avoid_: "lesson", "module", "card"

**Level**:
A band of nodes closed by that level's exam. The open level is fully written. Later levels are sketched as titles.
_Avoid_: "chapter", "unit", "tier"

**Frontier**:
The ready nodes in the open level. The session teaches one of them.
_Avoid_: "backlog", "queue", "next up"

**Quiz**:
The chat check that sets a concept or practice node to passed. It is not its own node.
_Avoid_: "test", "assessment", "HTML quiz"

**Exam**:
The node that closes a level, dated for a later session, covering that level's other nodes.
_Avoid_: "final", "boss quiz", "review session"

**Practice page**:
An HTML file linked from a node for a figure or a short in-browser task. It does not change node status.
_Avoid_: "lesson", "interactive lesson", "widget"

**Source**:
A high-trust entry in `RESOURCES.md` that a full node points at. Teaching a node comes from that source.
_Avoid_: "link", "reference material"

**Learning record**:
A short note of what the learner actually showed they know, or of a misconception that was corrected.
_Avoid_: "session log", "journal", "transcript"

**Reference**:
A compressed lookup sheet under `reference/`. This is the page to reopen later, not the session transcript.
_Avoid_: "cheatsheet dump", "notes"
