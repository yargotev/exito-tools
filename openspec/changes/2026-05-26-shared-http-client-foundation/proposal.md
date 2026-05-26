# Proposal: Shared HTTP client foundation

## Problem

Domain-owned provider clients currently duplicate low-level HTTP concerns such as base URL joining, bearer authentication, default timeouts, JSON request creation, response success checks, and bounded JSON decoding. ADR 0037 calls for shared low-level HTTP infrastructure while keeping provider DTO mapping in the owning domain.

## Scope

- Extend `internal/platform/httpclient` with provider-agnostic client configuration and request helpers.
- Provide shared defaults for timeout and maximum response body decoding.
- Refactor the Geo HTTP geocoder to use the shared helpers without changing its domain result or error contract.

## Out of Scope

- Retries, backoff, circuit breakers, logging, or provider-specific auth schemes.
- Orders provider configuration or Orders HTTP getter implementation.
- Changes to public capability IDs, JSON envelope shapes, or provider DTO ownership.
