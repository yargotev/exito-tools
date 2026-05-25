# Change: execution-metadata-foundation

## Why

The CLI now emits a JSON Envelope for `exito capabilities`, but its metadata is still incomplete. Agents need per-command request IDs, optional correlation IDs, and duration metadata so command output can be correlated with logs and future provider calls before real domain capabilities are added.

## What Changes

- Extend shared JSON Envelope metadata with `requestId`, optional `correlationId`, and `durationMs`.
- Add a small execution metadata helper for generated request IDs and duration measurement.
- Add a root `--correlation-id` flag and include it in `exito capabilities` output.
- Keep real capability execution, `exito run`, structured error translation, logging, and domain commands deferred.

## Out of Scope

- Generic `exito run` execution.
- Orders or Geo capabilities.
- slog wiring and outbound HTTP propagation.
- Full success/failure execution pipeline and error mapping.
