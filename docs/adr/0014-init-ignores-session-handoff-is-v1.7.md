# 0014. `tb init` ignores the session handoff; this increment is v1.7

Status: proposed
Date: 2026-10-15

## Context

Every profile's `templates` list includes `docs/session.md`, so `tb init`
scaffolds it into every product repo. The `handoff` and `onboard` skills
then write into it: the in-progress task, what changed, what is blocked,
the next step. The `handoff` skill already calls the file "transient by
design — do not park decisions here".

Nothing stops that file from being committed. In a public repo the
handoff becomes public: working notes, half-formed plans, and whatever a
session wrote about the user's environment. This happened in the toolbox
repo itself — `docs/session.md` was committed, and was caught only
because that commit had not yet been pushed. The fix there was manual (a
`.gitignore` entry plus dropping the unpushed commit). Product repos get
no such protection.

A `.gitignore` entry does nothing for a file that is already tracked, so
a repo that has already committed its handoff stays exposed even after
the entry is added. `tb` writes files; it does not touch the git index.

ADR 0011 stopped bumping `tbVersion` for content-only increments, so the
spec is at v1.6 while `tbVersion` is still `1.5.0`. This increment
changes `tb init`'s behavior, which is exactly when `tbVersion` moves.

## Options

- **Document it; no `tb` change** — tell the user in README and in the
  `handoff` skill to ignore `docs/session.md`. Relies on remembering, which
  is the failure that already happened once.
- **New profile field (`ignore = [...]`)** — each profile lists paths to
  ignore. Flexible, but a `profiles/*.toml` schema change that ADR 0005
  reserves for v2, for a list that has one entry and no profile that
  wants a different one.
- **`tb init` appends `docs/session.md` to the repo's `.gitignore`
  whenever the profile scaffolds `docs/session.md`** — no schema change.
  The rule follows from data the profile already has: if `tb init` is
  what puts the handoff file there, `tb init` is what keeps it out of
  git. Create-if-missing, append-if-absent, never rewrite existing lines.
- **Same, plus `git rm --cached` an already-tracked handoff** — closes the
  tracked-file gap automatically, but makes `tb init` modify the git
  index, which no `tb` command does today, and stages a deletion the user
  did not review.

## Decision

`tb init` appends `docs/session.md` to the product repo's root
`.gitignore` when the profile's `templates` list contains
`docs/session.md`. Every current profile does.

- **Idempotent by exact line.** If `.gitignore` already has a line equal
  to `docs/session.md` or `/docs/session.md` (surrounding whitespace
  ignored), nothing is written. `tb` does not evaluate glob patterns or
  consult `git check-ignore`: a match against global excludes on one
  machine would leave the file tracked on the next.
- **Append-only.** Missing `.gitignore` is created. An existing one gets a
  newline if it lacks a trailing one, then a one-line comment and the
  entry. Existing lines are never reordered or rewritten.
- **Same write gate as the rest of init.** When `tb init` refuses (non-TTY
  with changes, or the user declines), `.gitignore` is not touched either:
  "no files were changed" stays true. It needs no prompt of its own and
  no `--force`, because it can only add a line, like template creation.
- **Tracked handoff is warned, not fixed.** If `docs/session.md` is already
  tracked (`git ls-files --error-unmatch`), init prints a warning naming
  `git rm --cached docs/session.md`. `tb` still never touches the index.
- **`tbVersion` becomes `"1.7.0"`**, matching the spec version that
  introduced the behavior, rather than `1.6.0`, which would collide with
  spec v1.6 (ADR 0011), a content-only increment.

## Consequences

- `docs/requirements.md` gains **Delta from v1.6 + teach path (v1.7)** (§29)
  and **Acceptance (v1.7)** (§30); header becomes v1.7. `docs/design.md`
  gains §28–§29.
- `cmd/tb/internal/scaffold` gains the `.gitignore` append; `cli/init.go`
  calls it after `MaterializeTemplates` and reports `wrote:` /
  `unchanged:` for `.gitignore` and the tracked-file warning. Tests cover:
  missing file, file without trailing newline, entry already present (both
  spellings), profile without the template, refusal path untouched, and
  the tracked-file warning.
- Existing product repos pick this up with `tb init --force`.
- `README.md` documents the new write next to the existing init
  description.
- `docs/GLOSSARY.md` and every other `docs/` file stay tracked: they are
  durable project knowledge, not session state.
- No profile file, skill, or template changes. The `handoff` skill already
  describes the file as transient.
- Not in v1.7: a profile `ignore` field, `tb` modifying the git index,
  ignoring any other path, a doctor check for a tracked handoff, and
  Windows support.
