# application-bootstrap Spec Delta

## ADDED Requirements

### Requirement: Explicit domain dependency wiring

The Application bootstrap MUST explicitly wire domain capability dependencies from resolved configuration.

#### Scenario: Configured Orders provider is wired to orders.get

- **Given** `EXITO_ORDERS_BASE_URL` and `EXITO_ORDERS_TOKEN` are resolved
- **When** the Application boots
- **Then** `orders.get` uses an Orders HTTP getter
- **And** executing `orders.get` contacts the configured provider with bearer authentication

#### Scenario: Missing Orders provider config keeps not-configured behavior

- **Given** the Orders provider is not fully configured
- **When** the Application boots
- **Then** `orders.get` remains registered
- **And** execution fails with `ORDERS_NOT_CONFIGURED` instead of contacting a provider
