---
name: cpp-ggml
description: Use when reading, writing, or reviewing C/C++ in a ggml-family repo — llama.cpp, whisper.cpp, ggml itself — including tensor and backend code, and any work aimed at an upstream PR there; skip for non-ggml C/C++ and for the repo's Python tooling (`gguf-py`, `convert_*.py`).
---

# cpp-ggml

House style for the `cpp-systems` profile in ggml-family repos. Six
rules. The first one outranks the other five.

## The repo's own rules win

Read `AGENTS.md` and `CONTRIBUTING.md` in the target repo before writing
anything, every session — not from memory, and not from this skill. They
change, and they are policy, not advice. Where they disagree with this
skill, the persona, or the `AGENTS.md` of any other repo, **they win**.
`llama.cpp` also ships `skills/` (currently `add-new-model`,
`code-review`); check for one covering your task before starting.

**Do not run `tb init` here.** The repo has its own tracked `AGENTS.md`;
`tb init` would append a toolbox managed region to it and scaffold
`docs/requirements.md`, `docs/design.md`, `docs/tasks.md`, and
`docs/session.md` into the project's `docs/`. The upstream tree stays
byte-clean apart from your actual change. This skill and `cpp-verify` are
installed machine-wide by `tb install` and work without a pointer file.

Every toolbox write goes outside the clone, to `~/src/notes/<repo>/` —
`~/src/notes/llama.cpp/session.md` for `onboard` and `handoff`,
`~/src/notes/llama.cpp/GLOSSARY.md` for `glossary`. Create the directory
if it is missing. Those skills default to `docs/session.md` and
`docs/GLOSSARY.md`; here that is the project's own `docs/`, so redirect
them rather than writing there. Do not use `.git/info/exclude` to hide a
file inside the tree instead: it hides the file from your own
`git status` too, and `git clean -xdf` — the usual recovery from a bad
CMake build — deletes exactly what it covers.

## Never speak for the contributor

Non-overridable, and the stated penalty is a project ban:

- Do **not** write the PR description, the commit message, an issue, a
  review comment, or a reply to a reviewer. Not a draft, not "here is
  something to adapt". If asked, decline and say why.
- Do **not** run `git push`, `gh pr create`, `gh pr comment`, or
  `gh issue create`.
- If the user explicitly asks for a commit on their behalf: concise
  subject in `module : summary` form, and the trailer is
  `Assisted-by: <assistant name>`. **Never `Co-authored-by:`** — this
  overrides the default attribution trailer used in every other repo.
- Reading commands are encouraged: `gh search issues`, `gh search prs`,
  `grep`. Search before proposing anything; duplicates are closed.
- The PR template requires an AI-usage disclosure. The user writes it.
  Remind them it is required and that undisclosed AI use gets accounts
  banned.

A fully autonomous agent with no human in the loop is asked not to
contribute to this project at all.

## Understanding is the deliverable

A merged line is a maintenance obligation for a small team across a large
platform matrix, so the bar is not "the patch works" — it is that the
contributor understands it, can defend it in review without an assistant,
and will maintain it.

Guide before solving. If the user cannot yet explain the problem area,
point at the code and the docs and let them form the approach. Verify
comprehension before writing a change for them. A simpler change that
does 90% beats a complex one that does 100%.

Scope: features start as an issue, not a PR. A bug fix needs a
reproducible issue and a regression test that fails before the change and
passes after. New contributors keep one open PR and skip trivial fixes.
One PR per concern. First PR for a new model or feature is CPU-only;
CUDA and friends are follow-ups.

## Blend in

Read the surrounding code first and match it. If the change introduces a
new pattern or is large, **stop and tell the user** it will likely be
rejected without prior discussion.

- No third-party dependencies, no extra headers, no new subsystem. Reuse
  existing infrastructure.
- No fancy modern STL, no templates, no cleverness. Plain `for` loops.
- `snake_case` for functions, variables, types. Names optimize for the
  longest common prefix (`number_small`, not `small_number`).
- Public API pattern `<class>_<method>` (`llama_sampler_get_seed`);
  enum values upper case and prefixed with the enum name; sized integer
  types (`int32_t`) in the public API; `struct foo {}`, not a `typedef`;
  omit `struct` / `enum` keywords in C++ where optional.
- Filenames lowercase with dashes; `.h` / `.c` / `.cpp`.
- Cross-platform and cross-architecture, always.

## Comments last, and short

Write the code, then add a comment only where one is genuinely needed.
Writing comments first produces the redundant narration reviewers hate.

One or two lines. Simple English. Explain a non-obvious invariant, never
what the code already says. Do not hard-wrap to a fixed column and do not
split a sentence across lines. Never write a comment that addresses the
current task ("this fixes the phantom content you mentioned") — it is
meaningless to the next reader. Code copied from elsewhere keeps the
comments it had, and gains none.

## ggml facts that reviews fail on

- Tensors are row-major. Dimension 0 is columns, 1 is rows, 2 is
  matrices.
- `ggml_mul_mat` is unconventional: `C = ggml_mul_mat(ctx, A, B)` means
  `C^T = A B^T`, i.e. `C = B A^T`. Check this before trusting a shape.
- A backend op must match the CPU reference. Changing `ggml/` means
  `test-backend-ops` across two backends — see `cpp-verify`.
- Extending `ggml_type` with a new quantization type carries a
  disproportionate maintenance burden and its own evidence bar
  (converted model uploaded, perplexity and KL-divergence comparisons,
  CPU performance data). Do not start one casually.
- Public headers (`include/llama.h`, `ggml/include/*.h`) are an ABI
  surface. A change there is a design fork: raise it with the user and
  the issue tracker, not inside a diff.
