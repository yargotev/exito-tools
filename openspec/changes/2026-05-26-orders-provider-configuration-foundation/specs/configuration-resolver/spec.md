# Configuration Resolver Delta Spec

## Added Requirements

### Requirement: Orders provider configuration resolves from credential layers

The system MUST resolve Orders provider base URL and token values from credential layers using process environment first, profile-specific dotenv second, and general dotenv third.

#### Scenario: Environment values configure Orders provider

- **Given** `EXITO_ORDERS_BASE_URL` is set to `https://orders.example.test`
- **And** `EXITO_ORDERS_TOKEN` is set to `orders-token`
- **When** configuration is resolved
- **Then** `OrdersProvider.Configured` is true
- **And** `OrdersProvider.BaseURL` is `https://orders.example.test`
- **And** the Orders token is marked as set

#### Scenario: Environment token wins over dotenv token

- **Given** `.env.staging` contains `EXITO_ORDERS_BASE_URL` and `EXITO_ORDERS_TOKEN`
- **And** the process environment contains `EXITO_ORDERS_TOKEN`
- **When** configuration is resolved for profile `staging`
- **Then** the Orders base URL comes from `.env.staging`
- **And** the Orders token comes from the process environment

#### Scenario: Missing Orders token keeps provider unconfigured

- **Given** only `EXITO_ORDERS_BASE_URL` is present
- **When** configuration is resolved
- **Then** `OrdersProvider.Configured` is false
- **And** the Orders token is not marked as set

### Requirement: Orders provider token is not serialized

The system MUST omit the resolved Orders token value from JSON serialization of effective configuration.

#### Scenario: Effective configuration JSON does not expose Orders token

- **Given** `EXITO_ORDERS_TOKEN` is set to `super-secret-orders-token`
- **When** effective configuration is marshaled to JSON
- **Then** the JSON output does not contain `super-secret-orders-token`
- **And** the JSON output still includes non-sensitive Orders provider metadata
