---
name: go-verify
description: Use before claiming any task done in an `infra-go` or `back-go` repo, and after every change to Go code, CI, scripts, IaC, or SQL — this is the verify recipe (`gofmt -l .`; `go vet ./...`; `go test ./...`; `go test -race ./...`); skip only for edits that touch none of those.
---

# go-verify

The Go verify recipe (`infra-go` and `back-go` profiles). Run it in your
own shell; success is the command output, not an assertion. There is no
`tb verify` wrapper — `tb` does not shell out to `go`.

## Commands

```sh
gofmt -l .
go vet ./...
go test ./...
go test -race ./...
```

`gofmt -l .` must print nothing. `-race` is always required, not only when
the change "looks concurrent".

If the module lives under a subdirectory (this repo's `tb` lives in
`cmd/tb/`), run the commands from that module root so `./...` covers the
module and nothing else.

Touching CI, scripts, IaC, or SQL is not a skip: run this recipe, or the
repo's own tests/CI, or report that you could not.

## Rules

- Run the full recipe before ticking a checkbox in `docs/tasks.md` or
  reporting a task complete. `go build` succeeding is not verification.
- `go vet` findings are failures, not suggestions.
- Paste or summarize the actual output. Persona rule: evidence over
  narration.
- A pre-existing failure is a finding to report, not something to fold into
  your diff without a task for it.
