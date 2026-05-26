# Design: execution-metadata-foundation

## Approach

Add the smallest shared foundation needed for standard JSON Envelope metadata. The CLI surface remains responsible for command wiring; a new deep helper package owns request ID generation and duration-to-metadata conversion so future capability execution can reuse the same behavior.

## Boundaries

- `internal/capability` owns the Envelope metadata contract shape.
- `internal/execution` owns request metadata helpers that are independent from Cobra and presenters.
- `internal/surface/cli` owns parsing `--correlation-id` and applying metadata to implemented JSON commands.
- `internal/presenter` continues to only encode JSON.

## Decisions

- Request IDs use a generated `req_` prefix plus random hex to stay opaque and recognizable.
- `durationMs` is always present for JSON commands and may be `0` for very fast commands.
- `correlationId` is omitted when not supplied.
- `capabilityId` remains omitted for `exito capabilities` because it is a discovery command, not a registered neutral Capability Execution.
