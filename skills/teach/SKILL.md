---
name: teach
description: Use when the user wants to learn a topic over multiple sessions — a vault course with a path of levels, chat quizzes, and an exam per level; skip in a product repo (toolbox.toml, Cargo.toml, go.mod at cwd) unless they name a directory that is not the product root.
argument-hint: "What would you like to learn about?"
---

# teach

The user has asked you to teach them something. This is a stateful request — they intend to learn the topic over multiple sessions. The terminal session is the lesson. The course files are what the next session, in either agent, reads.

## Teaching Workspace

A product repo is never the course. If cwd has `toolbox.toml`, `Cargo.toml`, `go.mod`, or a generated `AGENTS.md`, do not write the course there.

Resolve the vault in this order:

1. Walk ancestors of cwd for `.obsidian`. The nearest one is the vault.
2. Else, if `~/Documents/Obsidian Vault/.obsidian` exists, that is the vault.
3. Else ask for the vault root and stop.

Do not scan `$HOME` for other vaults. The slug is the lowercase dash-case of the topic. The course is `<vault>/Learn/<slug>/`. Create that directory only after the mission interview has a concrete why.

The state of their learning is captured in this directory in several files:

- `MISSION.md`: The reason the user is learning this topic. Ground every session in it. Use the format in [MISSION-FORMAT.md](./MISSION-FORMAT.md).
- `GLOSSARY.md`: Canonical terms for this course. Use the format in [GLOSSARY-FORMAT.md](./GLOSSARY-FORMAT.md). Once a term is here, every session uses it.
- `PATH.md`: The path. Node sections are the record; the Mermaid map is rebuilt from them. Use the format in [PATH-FORMAT.md](./PATH-FORMAT.md). Do not repeat that schema here.
- `./reference/*.html`: Compressed reference (cheat sheets, algorithms, syntax). Beautiful documents which print well, designed for quick lookup.
- `RESOURCES.md`: Trusted sources. Use the format in [RESOURCES-FORMAT.md](./RESOURCES-FORMAT.md).
- `./learning-records/*.md`: What the user has actually shown they know. Use these, with the path, to judge the zone of proximal development. Titled `0001-<dash-case-name>.md`. Use the format in [LEARNING-RECORD-FORMAT.md](./LEARNING-RECORD-FORMAT.md).
- `./lessons/*.html`: Optional practice pages. Not one per node. See [Practice page](#practice-page).
- `./assets/*`: Reusable components shared across practice pages. See [Assets](#assets).
- `NOTES.md`: Preferences and working notes.

## Philosophy

To learn at a deep level, the user needs three things:

- **Knowledge**, captured from high-quality, high-trust resources
- **Skills**, acquired through practice on the path, based on that knowledge
- **Wisdom**, which comes from interacting with other learners and practitioners

Before `RESOURCES.md` is well-populated, find high-quality resources. Never trust parametric knowledge.

Some topics may require more skills than knowledge. Learning more about theoretical physics might be more knowledge-based. For physical work, more skills-based.

### Fluency vs Storage Strength

Split two kinds of learning:

- **Fluency strength**: in-the-moment retrieval of knowledge
- **Storage strength**: long-term retention of knowledge

Fluency can feel like mastery. Storage strength is the goal. Build it with desirable difficulty:

- Using retrieval practice (recall from memory)
- Spacing (the level exam is the spaced check)
- Interleaving (mixing related topics in practice — for skills practice only)

## Session

Read `MISSION.md`, `PATH.md`, and `./learning-records/` once the course exists. Rebuild the map from the node sections, per [PATH-FORMAT.md](./PATH-FORMAT.md), before any teaching.

1. Resolve the course as above. If `MISSION.md` is empty or missing, interview until the why is concrete, then create the directory and write it.
2. If `PATH.md` is missing: probe, write sources into `RESOURCES.md`, write Level 1 in full and later levels as sketched titles, show `PATH.md`, and wait. Do not teach the first node of a level until the user accepts the map.
3. If the open level's exam is `ready` and `due` is today or earlier, or `override` is `user`, give the exam before any new node. Pass and failure follow [PATH-FORMAT.md](./PATH-FORMAT.md).
4. Else if a node is `in-progress`, resume it.
5. Else teach one `ready` concept or practice. If several are ready, pick the one whose source serves the mission. If still tied, ask. Set it `in-progress` when teaching starts.
6. If every concept and practice in the open level is `passed` and the exam is not yet due, tell the learner the exam date and stop. Do not expand the next level.
7. After an exam passes, close the level. Expanding the next level is a new plan step — probe, sources, nodes, show `PATH.md`, and wait — not the same turn as the exam.

The chat quiz is the only writer of node status. A practice page does not set status.

Before creating `PATH.md`, and again before expanding a sketched level: for each strand that level will use, ask one question. If the answer is right, ask one harder question, then stop that strand. Do not probe a later level in that session unless the user asks.

Teach one node in four moves, then the quiz. Motivate why this node is next. Establish it from its source. Connect it to the nodes in `requires`. Quiz in chat. A miss leaves the node `in-progress`; repair it before building on it. A pass sets `passed` and writes a learning record only when [LEARNING-RECORD-FORMAT.md](./LEARNING-RECORD-FORMAT.md) says to.

Quiz options are bare claims of the same shape. Write the correct claim, then mutate it into each distractor. The explanation comes after the answer. Use the host's multiple-choice tool when it has one (`ask_user_question` on Grok Build). Otherwise ask in prose and wait. Do not mark the correct option.

## Practice page

Write a practice page only when a manipulable figure teaches the node, or the practice needs a short in-browser task. It is HTML under `lessons/`, built from `assets/`. Link it from the node's `page` field. The page does not write `PATH.md`.

A claim drawn as a picture is Mermaid in the map, or HTML/SVG whose source you can re-read. Do not send that claim through an image generator. When the host can open the page, look at it before the quiz. When it cannot, keep the page to a few elements.

Open the page for the user when you write one (`xdg-open`, `open`, or the editor). Remind them they can ask the agent about anything unclear.

## Assets

Practice pages are built from reusable **components**, stored in `./assets/`: stylesheets, simulators, diagram helpers, and anything else a second page could reuse.

Reuse is the default. Before authoring a practice page, read `./assets/` and build from the components already there. When a page needs something new and reusable, write it as a component in `./assets/` and link to it.

A shared stylesheet is the first component every course earns: every practice page links it.

## The Mission

Every node should tie into the mission — the reason the user is learning this topic.

If the user is unclear about the mission, or `MISSION.md` is empty, question them on why they want to learn this before writing the course.

A mission that is still vague will make the path abstract, and you will have no way to pick among ready nodes.

Missions change. Update `MISSION.md` and add a learning record when the reason changes. Confirm with the user before changing the mission.

## Zone Of Proximal Development

Each node should challenge the user just enough.

The bounded probe in [Session](#session) is how you locate that edge before writing or expanding a level. Between sessions, read `learning-records` and the path, and teach the ready node that best serves the mission.

## Knowledge

The knowledge in a node is only what that node needs. Gather it from the node's source in `RESOURCES.md` before teaching. Cite that source. For acquiring knowledge, difficulty eats the working memory you need for understanding.

## Skills

Skills are durability. The chat quiz and the level exam are retrieval. A practice page or a real-world sequence (for instance, assembling a drone) is how a practice node is exercised. A practice node may point at a trainer exercise; do not absorb the trainer into this skill.

Feedback is immediate: the quiz answer, or the result of the step they just did.

## Acquiring Wisdom

Wisdom comes from using the skill outside the course.

When a question needs wisdom, answer it, and also point at a community: a forum, a subreddit, a class, or a local group with a strong reputation. If the user does not want to join a community, record that in `RESOURCES.md` and stop proposing one.

## Reference Documents

While teaching nodes, also write reference documents. Sessions are rarely reread in full. Reference documents are. They are the compressed essence, for quick lookup.

Some topics lend themselves to reference:

- Syntax and code snippets for programming
- Algorithms and flowcharts for processes
- System design patterns and architecture
- Glossaries (`GLOSSARY.md`, not a file under `reference/`)

Once a glossary exists, every session uses it.

## `NOTES.md`

Record how the user wants to be taught, so the next session can follow it.
