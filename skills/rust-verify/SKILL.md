---
name: rust-verify
description: Use before claiming any task done in a `rust-systems` or `game-bevy` repo, and after every change to Rust code, CI, scripts, IaC, or SQL — this is the verify recipe (`cargo fmt --check`; `cargo test`; `cargo clippy --all-targets -- -D warnings`; Miri when the changed crate has `unsafe`); skip only for edits that touch none of those.
---

# rust-verify

The Rust verify recipe (`rust-systems` and `game-bevy` profiles). Run it
in your own shell; success is the command output, not an assertion. There
is no `tb verify` wrapper and never will be — `tb` does not shell out to
`cargo`.

## Commands

```sh
cargo fmt --check
cargo test
cargo clippy --all-targets -- -D warnings
```

If `cargo fmt --check` fails, format, include the formatting in the diff,
and re-run `--check`.

Workspace repos: `cargo test --workspace` and
`cargo clippy --workspace --all-targets -- -D warnings` when the change
crosses crate boundaries.

Touching CI, scripts, IaC, or SQL is not a skip: run this recipe, or the
repo's own tests/CI, or report that you could not.

## Miri

Run `cargo miri test` when a **changed crate contains any `unsafe`** in
its `.rs` sources, even if this diff only edits safe code in that crate.
Skip Miri only when the changed package(s) have no `unsafe` token.

If Miri is required and not installed, the recipe **fails** — do not tick
the task. Print `rustup +nightly component add miri`. Workspace: run Miri
on each changed crate that has `unsafe` (or `--workspace` if that still
covers them).

## Rules

- Run the full recipe before ticking a checkbox in `docs/tasks.md` or
  reporting a task complete. A green build is not verification.
- `-D warnings` means clippy warnings are failures. Fix them. `#[allow(...)]`
  requires the lint name and a reason that is not "to make clippy pass".
- Paste or summarize the actual output. Persona rule: evidence over
  narration.
- A failing test that was already failing before your change is still a
  finding — report it, do not fold it into your diff without a task for it.
