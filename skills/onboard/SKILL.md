---
name: onboard
description: Use on the first session in a brownfield repo the user has not worked in yet, when you need its layout, entry points, build and verify commands, and local conventions before touching anything; skip only when `docs/session.md` already contains an onboard map for this repo (shape, entry points, verify, conventions) and the area of work has not drifted. `AGENTS.md` existing is not a map.
---

# onboard

Map the repo before changing it. Cheap tools first (use the `scout` skill
for the searching itself). No code edits. Allowed writes: `docs/session.md`
and `docs/GLOSSARY.md` (see the `glossary` skill). `AGENTS.md` existing is
not a map; after `tb init` that file is only the managed region.

## Pass

1. **Shape** — top-level tree, build manifest (`Cargo.toml`, `go.mod`,
   workspace members), what is a binary vs a library, what is vendored or
   generated.
2. **Entry points** — `main`, exported API surface, the two or three files
   that everything else routes through.
3. **Verify recipe** — what CI actually runs, and what the profile's verify
   skill (`rust-verify` / `go-verify`) says. Note any gap between the two.
4. **Conventions in force** — error handling style, logging, test layout,
   whether `unsafe` appears, how config is loaded. Read the code for these;
   do not assume the profile's defaults are what this repo does.
5. **Existing harness state** — `toolbox.toml`, `AGENTS.md` managed region,
   `docs/` and `docs/adr/*`. Accepted ADRs bind you: read them before
   proposing anything they already settled.
6. **Landmines** — dead code, TODOs near the area of work, tests that are
   skipped or ignored.
7. **Domain language** — if a business domain is recognizable, use the
   `glossary` skill (business terms first: TTFT / NDVI, not "function").

## Output

Write the map into `docs/session.md` (the `handoff` skill's file) or, if
the work continues immediately, state it in one message before the first
edit. Name unknowns explicitly rather than guessing; an open question is a
better handoff than a confident wrong claim. Seeding or updating
`docs/GLOSSARY.md` is the other allowed write.

Do not refactor, reformat, or "fix" anything during an onboard pass.
