# Configuration Resolver Specification

## Purpose

Define shared Application Configuration precedence for Exito Tools.

## Requirements

### Requirement: Effective Profile precedence

The system MUST resolve the Effective Profile using explicit profile input, then `EXITO_PROFILE`, then saved Default Profile, then `staging` fallback.

#### Scenario: Explicit profile wins

- GIVEN explicit profile input and `EXITO_PROFILE` are both set
- WHEN configuration is resolved
- THEN the Effective Profile is the explicit profile

#### Scenario: Environment profile wins over saved default

- GIVEN no explicit profile is set
- AND `EXITO_PROFILE` is set
- AND a saved Default Profile exists
- WHEN configuration is resolved
- THEN the Effective Profile is `EXITO_PROFILE`

#### Scenario: Staging fallback is used

- GIVEN no explicit profile, no `EXITO_PROFILE`, and no saved Default Profile exist
- WHEN configuration is resolved
- THEN the Effective Profile is `staging`

### Requirement: Configuration File discovery precedence

The system MUST discover non-sensitive YAML configuration using explicit config path, then `EXITO_CONFIG`, then local `./exito.yaml`, then user `~/.config/exito-tools/config.yaml`, then internal defaults.

#### Scenario: Explicit config path wins without probing lower sources

- GIVEN an explicit config path is set
- WHEN configuration is resolved
- THEN that path is selected as the Configuration File source

#### Scenario: Local project config wins over user config

- GIVEN no explicit config path or `EXITO_CONFIG` is set
- AND local `./exito.yaml` exists
- AND user config exists
- WHEN configuration is resolved
- THEN local `./exito.yaml` is selected

#### Scenario: Defaults are selected when no config file exists

- GIVEN no explicit path, no `EXITO_CONFIG`, no local config, and no user config exist
- WHEN configuration is resolved
- THEN no Configuration File path is selected
- AND defaults remain the selected configuration source

### Requirement: YAML configuration can provide saved Default Profile

The Configuration Resolver MUST load a saved Default Profile from the selected non-sensitive YAML Configuration File when no saved Default Profile is supplied directly to the resolver.

#### Scenario: Local YAML default profile is used

- GIVEN no explicit profile and no `EXITO_PROFILE` are set
- AND local `./exito.yaml` contains `defaultProfile: prod`
- WHEN configuration is resolved
- THEN the Effective Profile is `prod`
- AND the Effective Profile source is saved default

#### Scenario: Explicit profile overrides YAML default profile

- GIVEN explicit profile input is `qa`
- AND local `./exito.yaml` contains `defaultProfile: prod`
- WHEN configuration is resolved
- THEN the Effective Profile is `qa`
- AND the Effective Profile source is explicit

#### Scenario: Environment profile overrides YAML default profile

- GIVEN no explicit profile is set
- AND `EXITO_PROFILE` is `dev`
- AND local `./exito.yaml` contains `defaultProfile: prod`
- WHEN configuration is resolved
- THEN the Effective Profile is `dev`
- AND the Effective Profile source is environment

### Requirement: Credentials remain outside YAML

The system MUST model credential source precedence separately from YAML configuration so secret values can come from real process environment, profile-specific dotenv, then general dotenv.

#### Scenario: Dotenv source order is profile-aware

- GIVEN the Effective Profile is `staging`
- WHEN credential layers are described
- THEN real process environment has highest priority
- AND `.env.staging` is considered before `.env`

### Requirement: Orders provider configuration

The Configuration Resolver MUST resolve Orders GEOMS configuration from non-sensitive YAML base URLs and sensitive environment/dotenv credentials, using process environment first, profile-specific dotenv second, and general dotenv third.

#### Scenario: Environment values configure Orders provider

- **Given** `EXITO_ORDERS_BASE_URL` is set to `https://orders.example.test`
- **And** `EXITO_ORDERS_TOKEN` is set to `orders-token`
- **When** configuration is resolved
- **Then** `OrdersProvider.Configured` is true
- **And** `OrdersProvider.BaseURL` is `https://orders.example.test`
- **And** the Orders token is marked as set

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

#### Scenario: Environment token wins over dotenv token

- **Given** `.env.staging` contains `EXITO_ORDERS_BASE_URL` and `EXITO_ORDERS_TOKEN`
- **And** the process environment contains `EXITO_ORDERS_TOKEN`
- **When** configuration is resolved for profile `staging`
- **Then** the Orders base URL comes from `.env.staging`
- **And** the Orders token comes from the process environment

#### Scenario: Missing Orders token and GEOMS credentials keeps provider unconfigured

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

### Requirement: Saved Default Profile can be persisted to YAML

The Configuration Resolver MUST provide a narrow persistence path for updating the saved Default Profile in the non-sensitive YAML Configuration File without writing credentials.

#### Scenario: Existing YAML default profile is updated

- GIVEN the selected Configuration File contains `defaultProfile: staging`
- WHEN the saved Default Profile is persisted as `prod`
- THEN the same Configuration File contains `defaultProfile: prod`
- AND no credential keys are written

#### Scenario: Missing YAML default profile is appended

- GIVEN the selected Configuration File exists without `defaultProfile`
- WHEN the saved Default Profile is persisted as `qa`
- THEN the file contains a top-level `defaultProfile: qa` entry

#### Scenario: No configuration file creates local project config

- GIVEN no explicit path, no environment path, no local config, and no user config exist
- WHEN the saved Default Profile is persisted as `dev`
- THEN local `./exito.yaml` is created with `defaultProfile: dev`

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


### Requirement: VTEX OMS provider configuration resolves by profile and brand

The Configuration Resolver MUST resolve non-sensitive VTEX OMS base URLs by Effective Profile and brand, while resolving VTEX OMS app key/token credentials only from process environment or non-committed dotenv files.

#### Scenario: Non-production Exito VTEX OMS uses QA credentials

- **Given** the Effective Profile is `staging`
- **And** Exito VTEX OMS base URL and QA app key/token values are available
- **When** configuration is resolved
- **Then** the Exito VTEX OMS provider is configured
- **And** it uses the QA credential variable names

#### Scenario: Production Carulla VTEX OMS uses production credentials

- **Given** the Effective Profile is `prod`
- **And** Carulla VTEX OMS base URL and production app key/token values are available
- **When** configuration is resolved
- **Then** the Carulla VTEX OMS provider is configured
- **And** it uses `CARULLA_APP_KEY` and `CARULLA_APP_TOKEN`

#### Scenario: VTEX OMS credentials are not serialized

- **Given** VTEX OMS app key and app token are configured
- **When** effective configuration is marshaled to JSON
- **Then** the JSON output does not contain the VTEX OMS app key or app token
- **And** the JSON output still includes credential presence/source metadata

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

### Requirement: VTEX Checkout provider configuration

The Configuration Resolver MUST resolve non-sensitive VTEX Checkout base URLs per profile and brand, while keeping cookies, orderForm ownership values, customer PII, and any credentials out of committed YAML and serialized effective configuration.

#### Scenario: Resolve VTEX Checkout base URL by brand

- **Given** the Effective Profile has `vtexCheckout.exito.baseUrl` configured
- **When** Checkout capabilities request the Exito Checkout provider
- **Then** the resolver MUST provide the configured base URL for Exito
- **And** the provider MUST be marked configured for public orderForm operations.

#### Scenario: Environment overrides YAML Checkout base URL

- **Given** YAML and environment/dotenv both provide a VTEX Checkout base URL for a brand
- **When** the Configuration Resolver builds effective configuration
- **Then** the environment/dotenv value MUST take precedence over YAML.

#### Scenario: Sensitive Checkout execution state is not serialized

- **Given** a Checkout execution uses cookies, orderForm ownership values, or customer profile data
- **When** effective configuration is serialized
- **Then** those sensitive execution values MUST NOT appear in the serialized configuration.

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

### Requirement: VTEX Master Data provider configuration resolves by profile and brand

The Configuration Resolver MUST resolve non-sensitive VTEX Master Data base URLs by Effective Profile and brand, while resolving VTEX app key/token credentials only from process environment or non-committed dotenv files.

#### Scenario: YAML configures Master Data brand providers

- GIVEN selected YAML contains `profiles.staging.vtexMasterData.exito.baseUrl`
- AND VTEX app key/token values are available for non-production Exito
- WHEN configuration is resolved for profile `staging`
- THEN the Exito Master Data provider is configured
- AND the base URL source is the Configuration File

#### Scenario: Environment overrides YAML Master Data endpoint

- GIVEN YAML provides a Master Data base URL for Carulla
- AND the matching `CARULLA_VTEX_MASTERDATA_BASE_URL_PROD` value is set
- WHEN configuration is resolved for profile `prod`
- THEN the environment value MUST take precedence over YAML

#### Scenario: Master Data credentials are not serialized

- GIVEN a Master Data app key and app token are configured
- WHEN effective configuration is marshaled to JSON
- THEN the JSON output MUST NOT contain the app key or app token
- AND it MAY expose only credential presence/source metadata

#### Scenario: Missing credentials leaves brand unconfigured

- GIVEN a Master Data base URL exists for a brand
- AND no VTEX app key/token values exist in environment or dotenv layers
- WHEN configuration is resolved
- THEN that brand provider MUST be marked unconfigured
