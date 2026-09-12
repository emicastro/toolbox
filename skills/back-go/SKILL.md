---
name: back-go
description: Use when writing or reviewing Go HTTP services, APIs, or anything with a listen loop and a database — timeouts, shutdown, handler error mapping, request logs, or SQL; skip for short-lived CLIs, one-shot jobs, and AWS/infra glue (use the `infra-go` skill) and for non-Go files.
---

# back-go

House style for the `back-go` profile: long-running HTTP services. Five
rules. Do not pick a router, ORM, or migrator here — that is an ADR in
the product repo. Changing any of the five at repo scope is also an ADR.

## Request context is the deadline

Handlers use `r.Context()`. Every call that can block (DB, HTTP client,
queue) takes that context as its first parameter and honours it. Do not
store a context in a struct. Set server timeouts (`ReadHeaderTimeout`,
`ReadTimeout`, `WriteTimeout`, `IdleTimeout`) and `http.Client` Timeout.
The defaults of `http.ListenAndServe` are not acceptable.

## Graceful shutdown

SIGINT/SIGTERM runs `http.Server.Shutdown` with a bounded context. Do
not `os.Exit` from a handler. `ListenAndServe` returning
`http.ErrServerClosed` is success. In-flight requests get the drain
window; new connections are refused.

## Errors map at the HTTP edge

Internal layers return `error`. One mapper at the handler or middleware
boundary turns them into a status and a stable client body. Never panic
on request input. Wrap with `%w` when adding context. Do not ignore
`error` (`_ =`, empty `if err != nil {}`). Do not leak internal error
strings to clients: log them; the client gets a public code or message.

## Structured logs, no secrets

Every request log: method, path, status, duration, request id. Never log
Authorization, cookies, passwords, tokens, or connection strings. In a
long-running service stdout and stderr are both diagnostics — the
`infra-go` rule that stdout is the pipe result does not apply here.
Secrets in the repo are the `aws-guard` skill; do not restate it.

## Parameterized SQL, explicit migrations

No query-string concatenation with request or user data; bound
parameters only. Schema changes are migration files applied by a named
command, not auto-migrate on process boot in production. A transaction
has one owner (who Begin/Commit/Rollback). Default listen address,
config, and DSN target local or dev; prod is named.
