---
name: onboard
description: Use on the first session in a brownfield repo the user has not worked in yet, when you need its layout, entry points, build and verify commands, and local conventions before touching anything; skip in a repo already mapped by an up-to-date `docs/session.md` or `AGENTS.md`.
---

# onboard

Map the repo before changing it. Read-only pass, cheap tools first (use the
`scout` skill for the searching itself), ending in a written map — not in a
diff.

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

## Output

Write the map into `docs/session.md` (the `handoff` skill's file) or, if
the work continues immediately, state it in one message before the first
edit. Name unknowns explicitly rather than guessing; an open question is a
better handoff than a confident wrong claim.

Do not refactor, reformat, or "fix" anything during an onboard pass.
