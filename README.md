# Toolbox

A developer harness for coding agents: one git clone on the machine plus a
small static Go binary (`tb`). Skills, persona, and verify recipes live here;
each product repo keeps a short pointer, a generated `AGENTS.md`, and its
own `docs/` (including ADRs).

Toolbox is not an MCP server, not a daemon, not a marketplace, and not part
of any product crate.

Source of truth for v1.4: [`docs/requirements.md`](docs/requirements.md).
How it is built: [`docs/design.md`](docs/design.md). Recorded forks:
[`docs/adr/`](docs/adr/). `tb init` writes `toolbox_version = "1.4.0"`.

## What it does

- One clone (`$TOOLBOX_HOME`, default `~/toolbox`) shared by every product
  repo.
- `tb init -p <profile>` wires a git repo to a named profile.
- `tb install` symlinks skills into Claude Code and Grok Build (v1's only
  install targets). Skills are **not** copied into product repos.
- Works on Arch Linux (x86_64) and macOS Apple Silicon (arm64).

## Install

Clone this repo. If it is not at `~/toolbox`, set `TOOLBOX_HOME` to the
clone path (must exist and be a directory):

```sh
export TOOLBOX_HOME="$HOME/src/toolbox"   # example; skip if you cloned to ~/toolbox
```

Build the binary for this machine:

```sh
make linux    # dist/tb-linux-amd64
# make darwin # dist/tb-darwin-arm64  (cross-compile check, or native on M1)
```

Put `tb` on `PATH`. Example:

```sh
ln -sf "$TOOLBOX_HOME/dist/tb-linux-amd64" ~/.local/bin/tb
# on M1: ln -sf "$TOOLBOX_HOME/dist/tb-darwin-arm64" ~/.local/bin/tb
```

Link skills into the agents on this machine (run after clone and after
`git pull` in `$TOOLBOX_HOME`):

```sh
tb install
tb doctor
```

`tb install` creates `~/.claude/skills` and `~/.grok/skills` if needed and
symlinks each `$TOOLBOX_HOME/skills/<name>` into both. It is idempotent. A
real file already at a target is skipped with a warning, not overwritten.

Agent detection (ADR 0004): binary on `PATH`, **or** a marker file
(`~/.claude/settings.json` / `~/.grok/config.toml`). One missing agent is a
warning (exit 0). Neither present is an error (exit 1). Skills are still
linked in the error case.

`tb` does not install or update Claude Code or Grok Build.

## Wire a product repo

From inside a git repository:

```sh
tb init -p rust-systems   # or -p infra-go, -p back-go, -p game-bevy, -p cpp-systems
```

This writes `toolbox.toml` (the profile pointer), renders `AGENTS.md`, and
creates missing `docs/` templates listed by the profile. Existing
`docs/adr/*` files are never touched, including under `--force`. Skills stay
in `$TOOLBOX_HOME`.

Re-running `tb init` on a repo that already has `AGENTS.md` /
`toolbox.toml`:

| | Interactive | Script / non-TTY | `--force` |
|---|---|---|---|
| Files would change | prompt: merge or refuse | refuse, exit 1, files untouched | refresh `AGENTS.md` (managed region only) and `toolbox.toml` |
| Bytes already match | no prompt; prints `unchanged:` | same | same |

`--force` never rewrites ADR bodies or other `docs/` files that already
exist. It does refresh the managed region, including the binding rules
and the `Skills:` line (`review` is a process skill in every profile).

## Commands

| Command | Role |
|---|---|
| `tb install` | Machine-level: symlink skills into Claude Code and Grok Build |
| `tb init -p NAME [--force]` | Wire the current git repo to a profile |
| `tb doctor` | Read-only health report (one screen). No `--fix` |
| `tb skill new NAME` | Scaffold `$TOOLBOX_HOME/skills/<name>/SKILL.md`. Does not attach it to a profile |
| `tb persona show [-p NAME]` | Print the persona file and its two section titles |

Exit codes: `0` success (including “one agent missing”), `1` operational
failure, `2` usage error.

`tb doctor` only hard-fails when **zero** agents are detected. Missing skill
links and a missing `cargo`/`go`/`cmake` are warnings.

## Profiles (v1.4)

| Profile | Use | Extra skills |
|---|---|---|
| `rust-systems` | Systems / low-level Rust, greenfield and legacy | `rust-verify`, `rust-systems` |
| `infra-go` | Go CLIs, scripts, AWS/infra glue — not long-running HTTP backends | `go-verify`, `infra-go`, `aws-guard` |
| `back-go` | Go HTTP services and APIs, greenfield and legacy — not CLIs or one-shot jobs | `go-verify`, `back-go`, `aws-guard` |
| `game-bevy` | Bevy games in Rust, greenfield and legacy — not systems crates or non-Bevy engines | `rust-verify`, `game-bevy` |
| `cpp-systems` | C/C++ built with CMake — ggml-family repos (llama.cpp) and your own C++ projects | `cpp-verify`, `cpp-ggml` |

All five profiles include the process skills `spec`, `adr`, `onboard`,
`scout`, `handoff`, `review`. `review` is the last gate after verify,
before ticking a task.

Verify recipes live in `rust-verify` / `go-verify` / `cpp-verify`. The
agent runs those commands in its own shell. There is no `tb verify`.

`rust-systems`:

```sh
cargo fmt --check
cargo test
cargo clippy --all-targets -- -D warnings
# cargo miri test  — when the changed crate contains any unsafe
```

`infra-go`:

```sh
gofmt -l .          # must print nothing
go vet ./...
go test ./...
go test -race ./...
```

`back-go`:

```sh
gofmt -l .          # must print nothing
go vet ./...
go test ./...
go test -race ./...
```

`game-bevy`:

```sh
cargo fmt --check
cargo test
cargo clippy --all-targets -- -D warnings
# cargo miri test  — when the changed crate contains any unsafe
```

`cpp-systems`:

```sh
cmake -B build -DCMAKE_BUILD_TYPE=RelWithDebInfo -DLLAMA_FATAL_WARNINGS=ON
cmake --build build -j $(nproc)
ctest --test-dir build -L main --output-on-failure --timeout 900
# clang-format the added lines only
# ./build/bin/test-backend-ops  — when the diff touches ggml/
```

C++ has no universal recipe: `cpp-verify` says to confirm the flags
against the target repo's CI workflows rather than recall them. In a repo
you do not own (llama.cpp is the motivating case) do **not** run
`tb init` — it would append a managed region to that project's tracked
`AGENTS.md` and scaffold toolbox docs into its `docs/`. `tb install` has
already linked the skills machine-wide; that is enough.

## Layout

```
$TOOLBOX_HOME/          # this repo (default ~/toolbox)
  profiles/             # rust-systems.toml, infra-go.toml, back-go.toml,
                        # game-bevy.toml, cpp-systems.toml
  personas/default.md
  skills/<name>/SKILL.md
  templates/            # files tb init copies if missing
  cmd/tb/               # Go module for the binary
```

After `tb init`, a product repo looks like:

```
<product>/
  toolbox.toml          # pointer: profile = "..."
  AGENTS.md             # generated managed region; user prose outside the markers is kept
  docs/                 # created if missing; docs/adr/* never overwritten
```

`AGENTS.md` is marker-delimited (`<!-- toolbox:begin -->` …
`<!-- toolbox:end -->`). Hand-edits inside the markers are replaced on the
next merge/`--force`. Content outside the markers is the repo's.

## Adding a skill

```sh
tb skill new my-skill
# edit $TOOLBOX_HOME/skills/my-skill/SKILL.md
# add "my-skill" to the chosen profiles/<name>.toml skills array
tb install
git -C "$TOOLBOX_HOME" add skills/my-skill profiles && git -C "$TOOLBOX_HOME" commit
```

`tb skill new` does not edit any profile. There is no per-repo `tb skill
add` / `rm` in v1.

## Building and testing `tb`

From this repo:

```sh
make all                          # linux/amd64 and darwin/arm64 into dist/
cd cmd/tb && gofmt -l . && go vet ./... && go test ./... && go test -race ./...
```

CI on `cmd/tb/**` runs `gofmt -l`, `go vet ./...`, and `go test ./...`. The
module is under `cmd/tb/` so those tools never walk `skills/` or
`profiles/`. `tb` has no third-party Go dependencies (ADR 0001).

## Not in v1.4

MCP, a local daemon, automatic agent memory, extra agents (Cursor, Codex,
…), a persona-switch CLI, per-repo skill enable/disable, wrapping verify as
`tb verify`, installing the coding agents themselves, a `tb` guard against
initialising an upstream repo. Sync is `git pull` in `$TOOLBOX_HOME` plus
`tb install`.
