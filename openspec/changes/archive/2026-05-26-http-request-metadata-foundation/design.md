# Design: HTTP request metadata foundation

## Approach

Introduce `internal/platform/httpclient` as the first narrow shared HTTP infrastructure slice. The package owns stable outbound header names and context helpers for provider-agnostic request metadata. The execution Pipeline stores the generated `requestId` and optional `correlationId` in the context before invoking a Capability handler. Domain HTTP clients can then call `httpclient.ApplyRequestMetadata` after constructing provider requests.

## Header Contract

- `X-Request-Id`: set when the execution request ID is available.
- `X-Correlation-Id`: set only when the caller supplied a correlation ID.

## Dependency Direction

`internal/execution` and domain HTTP clients may import `internal/platform/httpclient`. The platform package imports only standard library packages and does not know about domains, Cobra, presenters, or registry internals.
