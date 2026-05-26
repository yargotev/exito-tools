# Design: Shared HTTP client foundation

## Approach

Evolve `internal/platform/httpclient` from metadata-only helpers into the shared low-level HTTP infrastructure package described by ADR 0037. The package remains domain-agnostic and imports only standard library packages.

The shared `Client` owns:

- trimmed base URL and bearer token configuration;
- default timeout creation when no `*http.Client` is injected;
- provider-relative endpoint joining;
- JSON request construction with `Accept`, `Content-Type`, `Authorization`, and execution metadata headers;
- 2xx success detection and bounded JSON response decoding.

Domain packages continue to own provider DTOs and error translation. The Geo HTTP geocoder now delegates shared mechanics to `httpclient.Client` and still maps provider failures to Geo structured error codes.

## Dependency Direction

`internal/platform/httpclient` has no dependency on domain, capability, execution, surface, or config packages. Domain HTTP clients may import it for cross-cutting HTTP mechanics.
