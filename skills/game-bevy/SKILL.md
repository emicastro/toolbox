---
name: game-bevy
description: Use when writing or reviewing Bevy games, ECS gameplay, or anything with an `App` schedule — plugins, systems, components, assets, or frame-loop work; skip for systems/low-level crates, libraries, and non-Bevy engines (use the `rust-systems` skill) and for non-Rust files.
---

# game-bevy

House style for the `game-bevy` profile: Bevy games. Five rules. Do not
pick a Bevy version, physics crate, net crate, or UI crate here — that is
an ADR in the product repo. Changing any of the five at repo scope is
also an ADR. `[profile.dev]` opt-level is product convention, not a
sixth rule.

## ECS is the architecture

Components are data, systems are behavior, plugins own a domain. `main`
composes plugins. Do not model entities as objects with methods that
reach into other entities. Automated tests do not take a window or a
GPU: add the plugin under test to a headless `App`.

## Panic on broken world invariants; `Result` for I/O

A missing required component, a unique entity that is not unique, or a
state the schedule should have made impossible is a broken world: panic
with a message (`expect("…")`). File, asset, parse, and network failures
return `Result` and are not `unwrap`'d. Do not thread `Result` through
every system to avoid a panic that means the world is corrupt.

## The schedule is the concurrency model

Systems run in parallel unless ordered. Do not share game state with
`Arc<Mutex<_>>` or OS threads. Use `Commands`, events, and `Resources`.
If two systems conflict, fix the schedule (`before` / `after` /
`SystemSet`), do not add a lock. Exclusive `&mut World` is the rare case
that needs the whole world.

## Handles and components are values

Cloning `Handle<T>`, copying `Entity`, and owning components is normal.
Borrow-across-systems is the smell, not clone. Do not clone a
`Resource`'s inner collection to dodge the borrow checker — split the
resource or the system.

## Do not block the frame

No blocking I/O, asset decode, or network in `Update` / `FixedUpdate`.
Load through Bevy assets (`Handle` + load/asset events). Motion uses
`Time` deltas, not an assumed frame rate. Work that must block goes on
Bevy's task pools and comes back as an event.
