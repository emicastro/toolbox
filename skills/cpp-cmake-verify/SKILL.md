---
name: cpp-cmake-verify
description: Use before claiming any task done in a plain CMake C++ project with no ggml/backend matrix, and after every change to C/C++ sources, CMake, or CI — this is the verify recipe (cmake -B build; cmake --build build -j; ctest --test-dir build -L main --output-on-failure; clang-format on added lines; a sanitizer build for memory-touching diffs); skip for ggml-family repos (use cpp-verify instead) and for non-C++ files.
---

# cpp-cmake-verify

The verify recipe for a plain CMake C++ project you own — not a
ggml-family repo (llama.cpp, whisper.cpp, ggml itself; use `cpp-verify`
there instead, since its recipe and CI expectations are different). Run it
in your own shell; success is the command output, not an assertion. There
is no `tb verify` wrapper and never will be.

## Commands

```sh
cmake -B build
cmake --build build -j
ctest --test-dir build -L main --output-on-failure
```

A configure is not a build and a build is not a test. Run all three.

## Formatting

Format **the lines you added**, never the whole file and never a
neighbouring function — a reformatted region buries your change in a diff a
reviewer cannot read. Follow the repo's own `.clang-format` if it has one.

## Memory-touching diffs

Anything touching allocation, buffers, raw pointers, or process/thread
lifetime gets a sanitizer build before you call it done. Use the repo's own
CMake sanitizer option — read its `CMakeLists.txt` and confirm the actual
option name and value rather than assuming one (a plain CMake project does
not use `-DLLAMA_SANITIZE_ADDRESS`, that is `cpp-verify`'s target repo, not
this one). A common shape:

```sh
cmake -B build-san -DCMAKE_BUILD_TYPE=Debug <the repo's own sanitizer option(s)>
cmake --build build-san -j
ctest --test-dir build-san -L main --output-on-failure
```

## Artifacts

After the recipe, `git status --short` shows no build output. A tracked
binary or object file means `.gitignore` is incomplete, which is a finding.

## Rules

- Run the recipe before ticking a checkbox in `docs/tasks.md` or reporting
  a task complete. A green build is not verification.
- Paste or summarize the actual output. Evidence over narration.
- A test that was already failing before your change is still a finding —
  report it, do not fold a fix into your diff without a task for it.
- No `ci/run.sh`, no `test-backend-ops`, no backend-parity checks — those
  belong to `cpp-verify` and a ggml-family repo's own backend matrix, which
  this profile does not have.
