# capability-contract-foundation Spec Delta

## ADDED Requirements

### Requirement: orders.get domain execution

The system MUST execute `orders.get` through a domain-owned getter and return the stable `orders.GetResult` result shape.

#### Scenario: Configured Orders provider returns a mapped order

- **Given** the Orders provider is configured with a base URL and token
- **And** the provider returns `{"order":{"id":"A123","status":"created","createdAt":"2026-05-26T00:00:00Z"}}`
- **When** `orders.get` is executed with input `{"id":"A123"}`
- **Then** the capability succeeds
- **And** the data is an Orders domain result containing order ID `A123`, status `created`, and the provider `createdAt` value

#### Scenario: Orders provider returns not found

- **Given** the Orders provider is configured
- **And** the provider responds with HTTP 404
- **When** `orders.get` is executed for the missing ID
- **Then** the capability fails with structured error code `ORDER_NOT_FOUND`
