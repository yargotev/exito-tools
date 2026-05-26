# Design: CLI confirmation required foundation

## Approach

Extend `execution.ExecuteRequest` with a `Confirmed` boolean that surfaces can set after applying their own confirmation UX. The Pipeline checks the registered Capability definition after lookup and before input validation/handler execution. If `Definition.RequiresConfirmation` is true and `Confirmed` is false, it returns a failed envelope with stable code `CONFIRMATION_REQUIRED` and never calls the handler.

The generic CLI `run` command gains a local `--confirm` flag and passes it into the Pipeline. Read-only Capabilities and existing explicit domain commands continue passing the zero value (`false`) without behavior changes.

## Decisions

- Enforcement lives in the shared execution Pipeline so generic and explicit CLI paths share the same guard.
- The request carries confirmation intent, leaving future TUI confirmation prompts free to set the same field after a person confirms.
- Missing confirmation returns a normal failed JSON envelope on stdout and the existing generic failure exit category.
- The first flag is `--confirm`; stronger destructive target confirmation is deferred until a concrete destructive Capability exists.

## Risks

- Future destructive Capabilities may need stronger typed confirmation than a boolean. This slice keeps that out of scope and can be extended compatibly with additional request fields or metadata.
