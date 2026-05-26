# Design: Structured warning metadata foundation

## Approach

Introduce `capability.StructuredWarning` next to `StructuredError` so warning metadata remains shared and surface-neutral. `capability.ExecutionResult` gains a `Warnings` field for handlers to return non-fatal issues with successful data. `capability.EnvelopeMeta` gains optional `warnings` so JSON command output can carry those warnings without changing the envelope top-level success shape.

The execution Pipeline copies handler warnings into success envelope metadata. Failure envelopes remain focused on `error` metadata and do not propagate handler result warnings, because handlers that return an error do not have a successful result contract.

## Decisions

- Warning details use `map[string]any` for optional machine-readable context without introducing provider-specific DTOs.
- `warnings` uses `omitempty` to preserve existing JSON output when no warnings are present.
- Pipeline propagation is one-way from handler result to success metadata; validation and not-found failures do not synthesize warnings.

## Risks

- Details values must remain JSON-serializable when commands write envelopes. This mirrors existing `data any` behavior and can be tightened later if needed.
