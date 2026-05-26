# Proposal: Structured warning metadata foundation

## Summary

Add the narrow shared contract foundation for non-fatal Structured Warnings in JSON Envelope metadata.

## Motivation

ADR 0046 and the PRD require warnings to be machine-readable metadata instead of text mixed into stdout. The current Envelope metadata has request, profile, capability, and duration fields but cannot carry partial-data, fallback, or deprecation warnings returned by use cases.

## Scope

- Add a shared `StructuredWarning` contract with stable code, message, and optional details.
- Add optional warnings to standard Envelope metadata.
- Let Capability handlers return warnings alongside successful data through `ExecutionResult`.
- Preserve current failure behavior: failed envelopes do not expose handler result data or warnings.
- Add focused tests for serialization and Pipeline propagation.

## Non-goals

- Add warnings to existing Orders or Geo provider paths.
- Add warning rendering to the TUI.
- Add pagination or deprecation policies beyond the metadata shape.
