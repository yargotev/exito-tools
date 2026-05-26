# Change: capability-execution-pipeline-foundation

## Why

The CLI can now emit standard metadata for discovery output, but real Capability execution still has no shared path. Before adding `exito run` or domain commands, Exito Tools needs a deep, surface-independent execution pipeline that can locate a registered Capability, create request metadata, call its Use Case handler with context, and return the standard JSON Envelope shape.

## What Changes

- Extend the capability contract with a minimal executable Capability entry and context-aware handler shape.
- Extend the registry with executable entries and lookup by stable Capability ID while preserving immutable inventory snapshots.
- Add a surface-independent execution Pipeline that wraps success, structured failure, unknown capability, request ID, correlation ID, profile, capability ID, and duration metadata.

## Out of Scope

- `exito run <capability-id>` CLI command.
- Explicit Orders or Geo domain commands.
- Input schema validation.
- Warning/pagination metadata.
- HTTP propagation and slog wiring.
