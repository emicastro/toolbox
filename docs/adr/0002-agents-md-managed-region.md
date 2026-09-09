# 0002. Generated `AGENTS.md` uses a marker-delimited managed region

Status: accepted
Date: 2026-09-09

## Context

`tb init` writes a short `AGENTS.md` into the product repo (§8.2, §5). Users
will want to add their own prose to that file over time, and re-running
`tb init` (after a profile change, or a toolbox update) must refresh the
toolbox-owned content without destroying that prose. Requirements §8.2 also
require an interactive merge-or-refuse choice on an existing file, and a
`--force` mode that touches only `AGENTS.md` and `toolbox.toml`. Separately,
§14 asks whether the global `$TOOLBOX_HOME/AGENTS.md` is copied or only
referenced into the product repo.

## Options

- **Marker-delimited managed region** — `tb` owns everything between
  `<!-- toolbox:begin -->` and `<!-- toolbox:end -->`; everything outside
  those markers is the user's and is never touched. The region is rendered
  inline (profile, persona titles, skill list, verify recipe, language
  rule) rather than pointing at the global file.
- **Whole-file ownership** — `tb` owns the entire `AGENTS.md`; user notes go
  in a separate file `docs/session.md` or similar.
- **Reference-only stub** — the generated file is a few lines pointing at
  `$TOOLBOX_HOME/AGENTS.md` by absolute path; the agent reads both files.

## Decision

Marker-delimited managed region, rendered inline. Whole-file ownership
conflicts with the requirement that users can keep their own prose in
`AGENTS.md`. A reference-only stub makes the product repo's harness context
depend on a path outside the repo, which breaks if the repo is cloned
elsewhere or `TOOLBOX_HOME` differs across machines, and forces every agent
session to read a second file before it has any project context. Rendering
the content inline keeps the repo self-contained.

## Decision mechanics

- **No existing `AGENTS.md`**: write the full file, managed region plus a
  short human-editable footer, no prompt.
- **Existing file, no markers found**: treat as pre-toolbox content;
  interactive session offers *merge* (append a new managed region above the
  existing prose) or *refuse* (exit 1, print what would have changed);
  non-interactive session refuses.
- **Existing file, markers found**: *merge* replaces only the text between
  the markers, byte-for-byte, leaving everything outside untouched.
- **`--force`**: unconditionally rewrites the managed region (creating
  markers if absent) without prompting; never touches `docs/adr/*`, and
  never touches `AGENTS.md` content outside the markers.

## Consequences

- The renderer must produce byte-identical output for identical inputs
  (profile + persona + skills), so `merge` and `--force` are idempotent and
  diff-friendly in git.
- If a user hand-edits inside the markers, the next `tb init` silently
  overwrites that edit — this is documented in `tb init`'s help text and in
  design.md, not silently assumed.
- A future third install target (§12) only changes the rendered *content*
  of the skill list, not this mechanism.
