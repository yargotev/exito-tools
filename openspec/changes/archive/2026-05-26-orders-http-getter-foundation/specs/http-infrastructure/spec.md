# http-infrastructure Spec Delta

## ADDED Requirements

### Requirement: Provider domain clients use shared HTTP infrastructure

Domain-owned provider clients MUST use the shared HTTP infrastructure for base URL joining, bearer authentication, JSON requests, bounded JSON decoding, and execution metadata propagation.

#### Scenario: Orders HTTP getter sends authenticated metadata-bearing request

- **Given** an Orders HTTP getter has base URL `https://orders.example.test/api` and token `token-123`
- **And** request metadata contains request ID `req_orders` and correlation ID `corr-orders`
- **When** it gets order `A123`
- **Then** it sends a JSON `POST` request to `https://orders.example.test/api/orders/get`
- **And** the request body contains `{"id":"A123"}`
- **And** the request includes `Authorization: Bearer token-123`, `X-Request-Id: req_orders`, and `X-Correlation-Id: corr-orders`
