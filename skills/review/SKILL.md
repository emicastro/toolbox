---
name: review
description: >
  Use after a diff is ready to land, after the profile verify recipe, and
  before ticking a task or claiming done — re-read the change against
  tasks, ADRs, and the profile house style; also when the user asks for a
  review. Skip only when this session produced no diff.
---

# review

Last gate. Not a substitute for verify; run the profile verify recipe
first. Do not tick a task while findings remain unless the user defers a
finding to a new task.

## Pass

1. **Scope** — every file in the diff is required by the current task.
   Drive-by refactors and unrelated formatting are findings.
2. **ADRs** — no fork treated as settled without an ADR (see the `adr`
   skill).
3. **House style** — load the profile domain skill (`rust-systems`,
   `infra-go`, `back-go`, `game-bevy`, `cpp-ggml` for `cpp-systems`, or
   `c-cli`) and `aws-guard` when the diff touches AWS or secrets. Do not restate
   those skills here. In a repo the user does not own, that repo's own
   `AGENTS.md` / `CONTRIBUTING.md` outrank the domain skill and this one.
4. **Verify** — the recipe was actually run this session; output is in
   the conversation. `Not run` is a finding. `#[allow]`, `t.Skip`, ignored
   tests, and empty `if err != nil {}` in the diff are findings.
5. **Check** — new behavior has a `Check:` that can fail (see the `spec`
   skill).

## Output

Findings with `file:line`, or `no findings` naming the files read. Do not
refactor during review; fix or open a task.

"Small change" and "tests passed" are not skips.
