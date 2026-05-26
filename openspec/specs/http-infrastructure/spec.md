# Http Infrastructure Specification

## Requirements

### Requirement: Outbound HTTP requests propagate execution metadata

Provider HTTP requests created during a Capability execution MUST include the generated request ID and SHOULD include the caller-supplied correlation ID when present.

#### Scenario: Geo provider receives request and correlation IDs

- **Given** a `geo.geocode-address` execution with request ID `req_test` and correlation ID `corr-123`
- **When** the Geo HTTP geocoder calls the provider
- **Then** the outbound HTTP request includes `X-Request-Id: req_test`
- **And** the outbound HTTP request includes `X-Correlation-Id: corr-123`

#### Scenario: Correlation header is omitted when no correlation ID is supplied

- **Given** an execution with request ID `req_test` and no correlation ID
- **When** a shared HTTP request metadata helper applies metadata
- **Then** the outbound HTTP request includes `X-Request-Id: req_test`
- **And** the outbound HTTP request does not include `X-Correlation-Id`

### Requirement: Shared HTTP clients provide provider-agnostic request foundations

Domain-owned HTTP provider clients MUST be able to use shared infrastructure for base URL resolution, bearer authentication, JSON request headers, configured timeouts, and execution metadata propagation without centralizing provider DTOs or domain error translation.

#### Scenario: Shared client builds authenticated JSON provider request

- **Given** a shared HTTP client with base URL `https://provider.test/api/` and token `token-123`
- **And** request metadata with request ID `req_shared` and correlation ID `corr-shared`
- **When** it builds a JSON `POST` request for `/orders/get`
- **Then** the request URL is `https://provider.test/api/orders/get`
- **And** the request includes `Accept: application/json`
- **And** the request includes `Content-Type: application/json`
- **And** the request includes `Authorization: Bearer token-123`
- **And** the request includes `X-Request-Id: req_shared`
- **And** the request includes `X-Correlation-Id: corr-shared`

#### Scenario: Shared client rejects non-absolute base URLs

- **Given** a shared HTTP client with base URL `provider.test`
- **When** it resolves a provider endpoint
- **Then** endpoint resolution fails before sending an outbound request

### Requirement: Shared HTTP clients provide provider-agnostic response helpers

Domain-owned HTTP provider clients MUST be able to use shared infrastructure to identify successful HTTP responses and decode bounded JSON response bodies while preserving domain-owned error translation.

#### Scenario: Successful JSON response is decoded

- **Given** a provider response with a 2xx status and JSON body
- **When** shared response helpers evaluate and decode the response
- **Then** the response is treated as successful
- **And** its JSON body is decoded into the caller-owned DTO

#### Scenario: Non-2xx response is not successful

- **Given** a provider response with status `502`
- **When** shared response helpers evaluate the response
- **Then** the response is not treated as successful
