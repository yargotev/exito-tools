# Design: capability-execution-pipeline-foundation

## Approach

Introduce the narrowest executable Capability contract and a reusable execution Pipeline below all Interaction Surfaces. Application wiring still registers capabilities explicitly through the Registry Builder. The Pipeline accepts a Capability ID, neutral input object, profile, and optional correlation ID, then looks up the immutable registry entry and invokes the handler with `context.Context`.

## Boundaries

- `internal/capability` owns neutral execution contracts and structured error shapes.
- `internal/registry` owns boot-time registration, duplicate ID protection, immutable inventory, and ID lookup.
- `internal/execution` owns pipeline orchestration and metadata wrapping, independent of Cobra or presenters.
- `internal/surface/cli` remains unchanged in this slice; it will call the Pipeline in a later `exito run` slice.

## Decisions

- Handlers return `capability.ExecutionResult` so domains return structured data, not JSON bytes or presentation models.
- Handlers may return `capability.StructuredError` for stable domain/technical errors; unknown Go errors are translated to `CAPABILITY_EXECUTION_FAILED` by the Pipeline.
- Unknown Capability IDs return an `ok:false` envelope with `CAPABILITY_NOT_FOUND` and still include standard metadata.
- Registry duplicate IDs are rejected during boot to protect stable Capability ID contracts.
