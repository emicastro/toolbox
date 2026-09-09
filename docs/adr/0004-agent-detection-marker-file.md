# 0004. Agent detection: PATH or a marker file inside the config dir

Status: accepted
Date: 2026-09-09

## Context

ADR 0003's config-dir signal ("`~/.claude` exists") is unsound: `tb
install` itself must create `~/.claude/skills` to hold the symlinks it
manages, and `os.MkdirAll` creates `~/.claude` as that path's parent in
the process. After the first `tb install` run on a machine, `~/.claude`
exists whether or not Claude Code is actually installed — the same for
`~/.grok`/Grok Build. From the second run onward, every `tb doctor` and
every `tb install` reports both agents present unconditionally, which
directly breaks requirements.md §13 item 8 ("doctor fails closed on zero
agents after install"). This was found by testing `tb install` and `tb
doctor` together during group 4 implementation (docs/tasks.md), not by
inspection — the bug does not show up in a single isolated `tb install`
run, only across two.

The fix needs a config-dir-adjacent signal that `tb` itself never
creates, so linking skills can never manufacture a false positive.

## Options

- **Keep checking the bare directory** — accept the permanent false
  positive after first install; document it as a known limitation.
- **Check a marker file inside the config dir that only the real client
  creates** — `~/.claude/settings.json` for Claude Code,
  `~/.grok/config.toml` for Grok Build. Confirmed present on this
  project's own development machine (this session's own `~/.claude` has
  `settings.json`; `~/.grok` has `config.toml`), and neither is a file
  `tb` writes or has any reason to write.
- **Drop the config-dir signal entirely, PATH only** — simplest, but
  reintroduces ADR 0003's original problem: an app installed but not on
  `PATH` reads as absent.

## Decision

Marker file inside the config dir, keeping the PATH-OR-config-dir shape
from ADR 0003 otherwise unchanged: either signal is still sufficient,
only the config-dir signal's target changes.

Claude Code: `exec.LookPath("claude")` OR `~/.claude/settings.json` exists.
Grok Build: `exec.LookPath("grok")` OR `~/.grok/config.toml` exists.

This keeps ADR 0003's actual tradeoff (bias toward not failing, PATH
gaps covered by a config-dir fallback) while removing the one signal `tb`
can accidentally manufacture itself. `tb install` never writes
`settings.json` or `config.toml`, so this signal cannot be self-inflicted
the way the bare directory was.

## Consequences

- `internal/agents`' per-agent signal changes from a directory path to a
  file path (`target.configDir` becomes a marker file path); the
  PATH-or-signal shape and the `Detect()`/`DetectIn()` API are unchanged,
  so this is a one-file diff in `internal/agents`, not a package redesign.
- A test now specifically covers the regression this ADR fixes: running
  `tb install`'s directory-creation step must not, by itself, flip
  detection to "present" on a subsequent call.
- Still not a full liveness check (a stale `settings.json` left behind by
  an uninstalled Claude Code still reads as present) — same accepted v1
  gap ADR 0003 already named, just no longer a gap `tb` triggers itself.
- If Claude Code or Grok Build ever stop writing these specific files, or
  a third agent is added later (§12) with no analogous marker file, this
  needs a new ADR, not a silent path change.
