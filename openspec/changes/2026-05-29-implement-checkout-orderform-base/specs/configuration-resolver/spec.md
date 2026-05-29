# Configuration Resolver Specification Delta

## ADDED Requirements

### Requirement: Implement VTEX Checkout profile and brand configuration

The Configuration Resolver MUST expose non-sensitive VTEX Checkout base URLs by Effective Profile and brand for the first Checkout orderForm slice.

#### Scenario: YAML configures Checkout brand providers

- **Given** the selected YAML configuration contains `profiles.staging.vtexCheckout.exito.baseUrl`
- **And** it contains `profiles.staging.vtexCheckout.carulla.baseUrl`
- **When** configuration is resolved for the `staging` profile
- **Then** the effective configuration MUST expose configured Checkout providers for Exito and Carulla.

#### Scenario: Environment overrides YAML Checkout endpoint

- **Given** YAML provides a Checkout base URL for a brand
- **And** the matching `*_VTEX_CHECKOUT_BASE_URL_*` environment or dotenv variable is set
- **When** configuration is resolved
- **Then** the environment or dotenv value MUST take precedence using standard credential layer order.

#### Scenario: Checkout sensitive execution state stays out of configuration

- **Given** Checkout operations may use cookies, ownership values, or customer fields in future execution inputs
- **When** effective configuration is serialized
- **Then** those sensitive values MUST NOT be read from committed YAML
- **And** those sensitive values MUST NOT appear in serialized effective configuration.
