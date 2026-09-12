---
name: rust-systems
description: Use when writing or reviewing Rust in a systems or low-level crate — ownership and lifetime choices, error handling in library code, concurrency, FFI, or any diff that touches `unsafe`; skip for non-Rust files and for pure `cargo`/CI plumbing.
---

# rust-systems

House style for the `rust-systems` profile. These are the four things
reviews here actually fail on.

## Ownership

Borrow before you clone. `&str` / `&[T]` in signatures, owned types only
when the callee genuinely stores the value. A `clone()` added to silence
the borrow checker is a design smell — say so and fix the shape, or leave
a comment explaining why the copy is deliberate. Lifetimes stay elided
unless an explicit one carries meaning.

## No `unwrap` in libraries

Library code returns `Result`; it does not `unwrap`, `expect`, `todo!`,
`unimplemented!`, or `panic!` on input it did not produce. Define an error
type (or use the crate's existing one) and propagate with `?`. `unwrap` is
acceptable only in tests, in `main`/binary top level, and on an invariant
the surrounding code just established — and there it gets
`expect("<why this cannot fail>")`, not a bare `unwrap`.

## Explicit concurrency

Say what shares what. Prefer message passing to shared mutable state; when
a lock is required, note in a comment what it protects and the ordering if
more than one is held. No `Arc<Mutex<...>>` reached for by default. Async
and blocking code do not mix silently — no blocking call inside an async
context without a `spawn_blocking`-equivalent.

## Document `unsafe`

Every `unsafe` block carries a `// SAFETY:` comment stating the invariant
that makes it sound and who upholds it. Every `unsafe fn` documents its
preconditions on the item. Any diff touching `unsafe` runs Miri — see the
`rust-verify` skill.

Changing any of the above at repo scope is a design fork: write an ADR
(`docs/adr/NNNN-title.md`) rather than deciding it inside a diff.
