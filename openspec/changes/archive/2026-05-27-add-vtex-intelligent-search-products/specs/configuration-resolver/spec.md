# Configuration Resolver Specification Delta

## ADDED Requirements

### Requirement: VTEX Intelligent Search provider configuration

The Configuration Resolver MUST resolve non-sensitive VTEX Intelligent Search base URLs per profile and brand.

#### Scenario: YAML profile configures brand base URLs

- **Given** the selected profile contains `vtexIntelligentSearch.exito.baseUrl` and `vtexIntelligentSearch.carulla.baseUrl`
- **When** Application Configuration is resolved
- **Then** the effective configuration MUST expose configured Intelligent Search providers for Exito and Carulla.

#### Scenario: Environment overrides YAML base URL

- **Given** YAML and environment/dotenv both provide a VTEX Intelligent Search base URL for a brand
- **When** Application Configuration is resolved
- **Then** the environment/dotenv value MUST take precedence according to the standard credentials source order.

#### Scenario: Cookie secrets are not persisted in YAML

- **Given** Intelligent Search supports caller-provided VTEX cookies
- **When** Application Configuration is resolved or serialized
- **Then** cookie/session/segment values MUST NOT be read from committed YAML fields or exposed in JSON serialization.
