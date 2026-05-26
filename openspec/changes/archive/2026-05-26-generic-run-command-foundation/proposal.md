# Proposal: Generic Run Command Foundation

## Summary

Add the first `exito run <capability-id>` CLI path so agents can execute any registered Capability through its stable neutral ID.

## Scope

- Add a `run <capability-id>` Cobra subcommand.
- Reuse the existing surface-independent execution Pipeline.
- Accept complete JSON input objects from `--input-json`, `--input-file`, or stdin.
- Emit the standard JSON Envelope with execution metadata.
- Keep explicit domain commands and real domain capabilities deferred.

## Non-goals

- Add schema validation beyond requiring a JSON object.
- Add domain-specific commands or real Orders/Geo handlers.
- Finalize exit-code mapping for failed envelopes.
