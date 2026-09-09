---
name: infra-go
description: Use when writing or reviewing Go CLIs, one-shot jobs, scripts, or AWS/infra glue — anything touching exit codes, where logs go, cancellation, or which account a command acts on; skip for long-running HTTP service code, which this profile deliberately does not cover.
---

# infra-go

House style for the `infra-go` profile: short-lived tools that other things
script around. Five rules.

## Exit codes

The exit code is the API. `0` success, `1` operational failure, `2` usage
error (bad flags, missing required argument). Never exit `0` on a failed
operation, and never print an error and fall through. Document any
additional code in the command's help text.

## Logs to stderr, output to stdout

stdout carries the program's *result* — the thing a pipe consumes. Every
diagnostic, progress line, warning, and error goes to stderr. This is what
makes `tool | jq` work; violating it is a bug even when it looks fine
interactively.

## `context` cancellation

Anything that can block takes a `context.Context` as its first parameter
and honours it. Wire `signal.NotifyContext` for SIGINT/SIGTERM in `main` so
Ctrl-C actually stops in-flight work. Do not store a context in a struct.
Set a timeout on every network and AWS call.

## Explicit AWS profile

No implicit credential chain. The profile and region are named — by flag,
by an explicitly-read environment variable, or by config the tool prints
back — and the tool says which account it is about to act on before it
acts. See the `aws-guard` skill.

## No prod by default

The default target of any destructive command is the safe one: dry-run,
non-prod, or nothing at all. Prod requires an explicit flag or argument
naming it. A default that reaches prod is never acceptable, however
convenient.
