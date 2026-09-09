---
name: go-verify
description: Use before claiming any task done in an `infra-go` repo, and after every change to Go code — this is the verify recipe (`go test ./...`; `go vet`; race when tests involve concurrency); skip only for edits that touch no Go source.
---

# go-verify

The `infra-go` verify recipe. Run it in your own shell; success is the
command output, not an assertion. There is no `tb verify` wrapper — `tb`
does not shell out to `go`.

## Commands

```sh
go test ./...
go vet
```

And, when tests involve concurrency:

```sh
go test -race ./...
```

The race detector is required whenever the code under test spawns
goroutines, uses channels, or shares state across them — including tests
that only *exercise* concurrent code paths indirectly.

## Rules

- Run the full recipe before ticking a checkbox in `docs/tasks.md` or
  reporting a task complete. `go build` succeeding is not verification.
- `go vet` findings are failures, not suggestions.
- `gofmt -l .` should print nothing; unformatted code never lands.
- Paste or summarize the actual output. Persona rule: evidence over
  narration.
- A pre-existing failure is a finding to report, not something to fold into
  your diff without a task for it.
- If the module lives under a subdirectory (this repo's `tb` lives in
  `cmd/tb/`), run the commands from that module root so `./...` covers the
  module and nothing else.
