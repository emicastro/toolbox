# 0001. `tb` uses the Go standard library only

Status: accepted
Date: 2026-09-09

## Context

`tb` is a small, static binary that must build reproducibly on Arch Linux
(x86_64) and macOS Apple Silicon (arm64) with nothing more than a Go
toolchain. It has five subcommands and a config format it fully controls
(`toolbox.toml`, `profiles/*.toml`) — a small, closed TOML subset, not
arbitrary user TOML. There is no requirement to consume third-party TOML
files, and no plan to grow `tb` into a large CLI surface.

## Options

- **Stdlib only** — `flag.NewFlagSet` per subcommand; a hand-rolled reader
  for the TOML subset `tb` actually emits and reads.
- **Stdlib CLI + `BurntSushi/toml`** — keep `flag`, add a real TOML library
  for config parsing.
- **`cobra` + `BurntSushi/toml`** — full subcommand framework plus a TOML
  library.

## Decision

Stdlib only. `tb`'s own config format is a closed subset (comments, bare
keys, quoted strings, flat string arrays) that a ~150-line reader parses
completely; a general TOML library would parse far more than `tb` ever
writes, which is unneeded surface for a two-machine, single-user tool. Five
subcommands do not justify a subcommand framework — `flag.NewFlagSet` per
command plus a manual dispatch in `main` is enough, and it keeps `go.mod`
free of a `require` block.

## Consequences

- `go build` needs no module download; the binary is trivially reproducible
  and auditable on both machines.
- The TOML reader must reject anything outside its subset with a clear
  "file:line: unsupported syntax" error rather than silently misparsing —
  covered by table tests in task group 3.
- If a future profile or config need grows past the subset (nested tables,
  inline tables, multi-line strings), that is a new ADR, not a silent
  parser patch.
