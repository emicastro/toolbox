---
name: trainer-maker
description: Use when writing or reviewing a rustlings-style exercise trainer — an ordered set of small broken programs a learner fixes and unlocks one at a time, checked on save, in any language or subject domain; skip for the exercises' own subject-matter content and for any single language's compiler/toolchain mechanics (pair this with your profile's own verify/house-style skill for that).
---

# trainer-maker

House style for building or reviewing a rustlings-style exercise trainer.
Eight rules. It is language- and domain-agnostic: it does not know C++ from
Rust, or inference engineering from distributed databases. Pair it with
whatever verify/house-style skill your actual profile already gives you for
the toolchain (`cpp-cmake-verify`, `rust-verify`, `go-verify`, ...) — that
skill owns the compiler/interpreter mechanics; this one owns the trainer
shape. Changing any of the eight at repo scope is an ADR in that product.

## Manifest and anti-spoiler workspace

One ordered, declarative manifest (chapters → exercises → metadata →
hints); hints never live inside learner-visible exercise source. A separate
learner workspace, populated by an `init`-style command, starts with zero
solutions; a solution is revealed into it only after that exercise's check
passes.

The manifest's file format (TOML, YAML, JSON, ...) is a repo-scope design
fork — write an ADR rather than picking one here.

## A passing check is not "move on"

A check flips a state bit; the learner explicitly dismisses the exercise
(for example, removing a marker string from the file) before the runner
advances to the next one. Do not conflate "the check passed" with "the
learner is done looking at this" — they are different signals with
different owners.

## `watch` is the default command, with a fixed loop shape

A no-args invocation finds the first unfinished exercise, runs its check
once immediately, prints a chapter banner only on a chapter transition,
then blocks for either a file-save event or a single-letter interactive
command (hint / list / re-run / quit), always reprinting the command legend
afterward so it is never lost off-screen.

Debounce file-watch events so a burst of editor writes collapses into one
recheck. The watch backend (inotify, kqueue, polling, ...) is a repo-scope
fork.

## The check protocol is minimal and runner-owned

The runner and the exercise agree on a trivial pass/fail contract embedded
in the exercise language itself. Do not pull in the host language's full
test framework for this; a heavyweight framework's own failure output
competes with the runner's for the learner's attention.

## Reproducible commands, one color policy

Every subprocess invocation made on the learner's behalf is recorded and
printed verbatim, so it can be rerun by hand outside the trainer. Raw tool
output (compiler errors, test failures) is shown through unmodified, never
rewritten.

One TTY-and-`NO_COLOR` policy governs both the trainer's own output and any
color flag it passes to an invoked tool — what is on screen must match what
was printed as the reproducible command.

## State is small, atomic, and independent of any subprocess it starts

Progress (which exercises are done, each exercise's hint tier capped at its
hint count) lives in one small file, written via temp-file-then-rename,
located by walking up from the current directory to a workspace marker.
Build/scratch output is a separate cache, never mixed into that state file.

On interrupt, the trainer's own exit must not leave a learner-owned
subprocess running past it. The signal/process mechanism used to guarantee
that is a repo-scope fork.

## Standard CLI surface

`init`, `watch` (the default, no-args command), `run <name>` (check once),
`hint <name>` (reveal the next tier), `list` (progress by chapter),
`verify` (check every exercise in manifest order, stop at the first
failure), `reset <name>` (restore the original file and clear that
exercise's state, warning about anything downstream that depends on it),
`--help`/`-h`, `--version`/`-V`.

Misuse is exit code 2; an exercise's check failing as expected is not a
crash.

## Ship bundled; keep maintainer tooling and docs separate from the learner path

Exercises and solutions are embedded or bundled into the distributed
artifact, so the installed tool and the learner's workspace do not depend
on a live source checkout.

A maintainer-only lint (manifest integrity) and grade (the solution passes,
the untouched exercise fails) pair stays off the learner's command surface
entirely. Course documentation is layered — a short goal comment on the
exercise itself, a per-chapter README with fixed sections, one
whole-course roadmap doc — and `watch` only ever auto-surfaces the
per-chapter layer, on a chapter transition.

Compiler/toolchain invocation mechanics, sanitizers, a "the toolchain must
refuse this" exercise mechanism, and the course content itself are not
house style here — they belong to the profile's own verify/house-style
skill, or to a product ADR.
