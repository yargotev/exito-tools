# Configuration Resolver Delta Specification

## ADDED Requirements

### Requirement: YAML profiles can provide provider base URLs

The Configuration Resolver MUST load non-sensitive provider base URLs from the selected YAML Configuration File for the Effective Profile when no higher-priority environment or dotenv value provides that provider base URL.

#### Scenario: YAML profile base URLs configure providers with environment tokens

- GIVEN local `./exito.yaml` contains provider base URLs under `profiles.staging.geo.baseUrl` and `profiles.staging.orders.baseUrl`
- AND `EXITO_GEO_TOKEN` and `EXITO_ORDERS_TOKEN` are set in the process environment
- WHEN configuration is resolved for profile `staging`
- THEN the Geo and Orders providers are configured
- AND their base URLs come from the Configuration File
- AND their token values come from the process environment

#### Scenario: Environment base URL overrides YAML profile base URL

- GIVEN local `./exito.yaml` contains `profiles.staging.geo.baseUrl`
- AND `EXITO_GEO_BASE_URL` is set in the process environment
- WHEN configuration is resolved for profile `staging`
- THEN the Geo provider base URL is `EXITO_GEO_BASE_URL`
- AND the Geo provider base URL source is environment

#### Scenario: Effective Profile selects matching YAML profile

- GIVEN local `./exito.yaml` contains base URLs for profiles `staging` and `prod`
- AND the Effective Profile is `prod`
- WHEN configuration is resolved
- THEN provider base URLs are loaded from the `prod` YAML profile

#### Scenario: YAML token-like keys are ignored

- GIVEN local `./exito.yaml` contains provider base URLs and token-like keys under a profile
- AND no provider token exists in environment or dotenv layers
- WHEN configuration is resolved
- THEN the providers are not configured
- AND token values are not read from YAML
