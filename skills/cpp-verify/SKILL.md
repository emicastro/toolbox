---
name: cpp-verify
description: Use before claiming any task done in a `cpp-systems` repo, and after every change to C/C++ sources, CMake, CI, or shaders — this is the verify recipe (`cmake -B build`; `cmake --build build -j`; `ctest --test-dir build -L main --output-on-failure`; `clang-format` on added lines; `test-backend-ops` when `ggml/` changed); skip only for edits that touch none of those.
---

# cpp-verify

The C/C++ verify recipe (`cpp-systems` profile). Run it in your own
shell; success is the command output, not an assertion. There is no
`tb verify` wrapper and never will be.

Unlike `cargo` or `go`, C++ has no universal recipe. The commands below
are the llama.cpp shape. **The repo's CI workflows are the source of
truth** — read `.github/workflows/build-cpu.yml` (or the workflow for the
backend you touched) and `docs/build.md` and confirm the flags before you
run them. Do not recall flags from memory; a flag that was renamed
upstream produces a build that proves nothing.

## Commands

```sh
cmake -B build -DCMAKE_BUILD_TYPE=RelWithDebInfo -DLLAMA_FATAL_WARNINGS=ON
cmake --build build -j $(nproc)
ctest --test-dir build -L main --output-on-failure --timeout 900
```

`-DLLAMA_FATAL_WARNINGS=ON` is this profile's `-D warnings`: compiler
warnings are failures. Fix them; do not silence them with a pragma.

A configure is not a build and a build is not a test. Run all three.

## Formatting

The repo has `.clang-format`, `.editorconfig`, and `.ecrc`. Format **the
lines you added**, never the file and never a neighbouring function — a
reformatted region buries your change in a diff a reviewer cannot read.
clang-format v15+.

Non-negotiable from `.editorconfig` and the project's agent rules: 4
spaces, no tabs, no trailing whitespace, brackets on the same line,
`void * ptr` and `int & a` spacing, ASCII only. No em dash, no `→`, no
`×`, no `…` anywhere in code, comments, or commit text — use `-`, `->`,
`x`, `...`.

## When you changed `ggml/`

`ctest` is not enough. Backend implementations of an operator must agree
with the CPU reference:

```sh
./build/bin/test-backend-ops -o <OP_NAME>   # e.g. MUL_MAT
./build/bin/test-backend-ops                # full sweep
```

This needs **at least two backends** built (CPU plus CUDA / Metal /
Vulkan / SYCL). If you only have CPU, say so in your report — you did not
verify the change, you verified half of it.

New or modified operator: add its case to `tests/test-backend-ops.cpp`.
That is an edit to an existing test file, which is expected. Adding a
**new file** under `tests/` needs maintainer approval first.

## Memory-touching diffs

Anything touching allocation, buffers, graph allocation, or raw pointers
gets a sanitizer build before you call it done:

```sh
cmake -B build-san -DCMAKE_BUILD_TYPE=Debug -DLLAMA_SANITIZE_ADDRESS=ON \
      -DLLAMA_SANITIZE_UNDEFINED=ON
cmake --build build-san -j $(nproc)
ctest --test-dir build-san -L main -E tokenizer --output-on-failure
```

(`LLAMA_SANITIZE_THREAD` is separate and does not combine with address.)

## Before publishing

`ci/README.md` documents the full local CI, which is what the project
asks contributors to run before opening a PR:

```sh
mkdir -p tmp && bash ./ci/run.sh ./tmp/results ./tmp/mnt
```

It downloads models and takes a long time. Run it, or report that you
could not and which parts of the matrix are therefore unproven.

For a change that could affect quality or speed, `llama-perplexity` and
`llama-bench` before/after are part of verification, not a nice-to-have.

## Rules

- Run the recipe before ticking a checkbox in `docs/tasks.md` or
  reporting a task complete. A green build is not verification.
- Paste or summarize the actual output. Persona rule: evidence over
  narration.
- A test that was already failing before your change is still a finding —
  report it, do not fold a fix into your diff without a task for it.
- Never report "tests pass" when you built one backend and the change
  spans several. Name the backend you built.
