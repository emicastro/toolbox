---
name: rust-verify
description: Use before claiming any task done in a `rust-systems` repo, and after every change to Rust code — this is the verify recipe (`cargo test`; `cargo clippy -- -D warnings`; Miri when `unsafe` is in the change); skip only for edits that touch no Rust source.
---

# rust-verify

The `rust-systems` verify recipe. Run it in your own shell; success is the
command output, not an assertion. There is no `tb verify` wrapper and never
will be — `tb` does not shell out to `cargo`.

## Commands

```sh
cargo test
cargo clippy -- -D warnings
```

And, when `unsafe` is in the change:

```sh
cargo miri test
```

Miri is required whenever the diff adds, moves, or modifies an `unsafe`
block or an `unsafe fn` — not only when new `unsafe` is introduced. If Miri
is not installed on the machine, say so explicitly instead of skipping
silently: `rustup +nightly component add miri`.

## Rules

- Run the full recipe before ticking a checkbox in `docs/tasks.md` or
  reporting a task complete. A green build is not verification.
- `-D warnings` means clippy warnings are failures. Fix them; do not add
  `#[allow(...)]` without a one-line comment saying why.
- Paste or summarize the actual output. Persona rule: evidence over
  narration.
- A failing test that was already failing before your change is still a
  finding — report it, do not fold it into your diff without a task for it.
- Workspace repos: `cargo test --workspace` and
  `cargo clippy --workspace --all-targets -- -D warnings` when the change
  crosses crate boundaries.
