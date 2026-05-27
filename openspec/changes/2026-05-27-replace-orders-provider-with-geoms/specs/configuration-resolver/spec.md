## MODIFIED Requirements

### Requirement: Orders provider configuration

The Configuration Resolver MUST resolve Orders GEOMS configuration from non-sensitive YAML base URLs and sensitive environment/dotenv credentials.

#### Scenario: GEOMS credential bundle configures Orders

- **Given** `EXITO_ORDERS_BASE_URL` is set
- **And** `GEOMS_CREDENTIALS_QA` contains `client_id`, `client_secret`, and `scope`
- **When** configuration is resolved for the `staging` profile
- **Then** `OrdersProvider.Configured` is true
- **And** the client ID, secret, and scope are available to application wiring but omitted from JSON serialization

#### Scenario: Prod profile reads PDN GEOMS bundle

- **Given** `EXITO_ORDERS_BASE_URL` is set
- **And** `GEOMS_CREDENTIALS_PDN` contains `client_id`, `client_secret`, and `scope`
- **When** configuration is resolved for the `prod` profile
- **Then** `OrdersProvider.Configured` is true
