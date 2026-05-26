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
