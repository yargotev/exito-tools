# Proposal: HTTP request metadata foundation

## Problem

Capability executions already generate request and correlation metadata for CLI envelopes, but outbound HTTP provider calls do not yet receive that metadata. ADR 0042 and ADR 0043 require request metadata to be propagated to provider requests when applicable.

## Scope

- Add a small shared HTTP infrastructure package for provider-agnostic request metadata.
- Attach execution request/correlation metadata to the context used by Capability handlers.
- Apply the shared metadata headers from Geo HTTP provider requests.

## Out of Scope

- Retries, backoff, provider-specific auth abstractions, or a full shared HTTP client factory.
- Changes to provider DTO mapping or public capability result shapes.
