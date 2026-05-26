# Proposal: Pagination metadata foundation

## Summary

Add the narrow shared contract foundation for cursor-based pagination metadata in JSON Envelope metadata.

## Motivation

ADR 0045 and the PRD require list-style Capabilities to expose explicit cursor-based pagination without hidden auto-pagination. The current Envelope metadata can carry request, profile, capability, duration, and warnings, but it cannot represent `nextCursor` and `hasMore` for future list results.

## Scope

- Add a shared Pagination Metadata contract with optional cursor and explicit `hasMore` state.
- Add optional pagination metadata to standard Envelope metadata.
- Let Capability handlers return pagination metadata alongside successful data through `ExecutionResult`.
- Preserve current failure behavior: failed envelopes do not expose handler result data, warnings, or pagination metadata.
- Add focused tests for serialization and Pipeline propagation.

## Non-goals

- Add a new list Capability.
- Add pagination flags to existing Orders or Geo commands.
- Implement auto-pagination.
- Interpret or transform provider cursors.
