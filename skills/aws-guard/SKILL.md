---
name: aws-guard
description: Use before running, writing, or reviewing anything that talks to AWS or handles credentials — CLI invocations, SDK calls, IaC, CI steps, or a change that could put a secret in the repo; triggers on any command whose target account, profile, or region is not stated in the command itself.
---

# aws-guard

Three refusals. They apply to commands you run yourself and to code you
write for someone else to run.

## Refuse implicit prod

If a command could act on production and nothing in it says so, do not run
it. Ask which environment, or run the non-prod equivalent. This covers the
default credential chain, an ambient `AWS_PROFILE` you did not set, an
assumed role inherited from the shell, and any tool whose default target is
prod. Destructive verbs (`delete`, `terminate`, `put`, `apply`, `destroy`,
force-flags) additionally get a dry-run or a plan first, shown to the user
before the real call.

## No secrets in the repo

Access keys, session tokens, `.env` files with real values, private keys,
connection strings with passwords: never written to a tracked file, never
pasted into `docs/`, never embedded in a commit message. Use the named
profile, a secrets manager reference, or an environment variable read at
runtime. If you find a secret already committed, stop and report it — do
not "fix" it with a follow-up commit, which does not remove it from history.

## Profile and region must be named

Every AWS invocation states both, explicitly:

```sh
aws --profile <name> --region <region> <command>
```

In Go, construct the config with an explicit profile and region rather than
relying on defaults. In CI, the profile/role and region come from a named
configuration step, not from whatever the runner happened to have.

When a rule here blocks a task, that is the correct outcome: say what is
blocked and what you need (which account, which profile) rather than
choosing a default and proceeding.
