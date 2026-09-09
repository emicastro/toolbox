# 0003. Agent detection: PATH or known config dir, either is presence

Status: accepted
Date: 2026-09-09

## Context

`tb install` and `tb doctor` must detect whether Claude Code and Grok Build
are present on the current machine (§8.1, §8.3). Requirements state the
exact heuristic is a design choice, but fix the failure policy: exactly one
missing is a warning with exit 0, both missing is a hard failure. The two
machines in scope (Arch custom install, MacBook Air M1) may have an agent
installed without its binary on `PATH` (e.g. a GUI-launched app), or a
binary on `PATH` without having been run yet (no config dir created).

## Options

- **PATH only** — `exec.LookPath("claude")` / `exec.LookPath("grok")`.
- **Config dir only** — check for `~/.claude` / `~/.grok`.
- **PATH or config dir (OR)** — either signal is sufficient to count the
  agent as present.

## Decision

PATH or config dir, OR'd together. PATH-only misses an installed app not on
`PATH`; config-dir-only misses a freshly installed binary that hasn't been
run yet. Since the cost of a false positive here is low (`tb install` links
a skill dir that an agent will simply not read if it isn't actually
installed) and the cost of a false negative is a wrongly hard-failed
`install`, the OR captures the union and biases toward not failing.

Claude Code: `exec.LookPath("claude")` OR `~/.claude` exists.
Grok Build: `exec.LookPath("grok")` OR `~/.grok` exists.

## Consequences

- `tb doctor` reports each agent as present/absent using the same check
  `install` used, so the two commands never disagree.
- A stale `~/.claude` or `~/.grok` directory left behind by an uninstalled
  agent reads as "present." This is accepted for v1; a real health check
  (e.g. running `claude --version`) is future work, not required by §8.1/§8.3.
- The exact binary names and config dirs are a small lookup table in
  `internal/agents`, so adding a third agent later (§12) is a one-line
  addition, not a design change.
