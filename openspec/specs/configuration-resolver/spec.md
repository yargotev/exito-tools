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

### Requirement: Credentials remain outside YAML

The system MUST model credential source precedence separately from YAML configuration so secret values can come from real process environment, profile-specific dotenv, then general dotenv.

#### Scenario: Dotenv source order is profile-aware

- GIVEN the Effective Profile is `staging`
- WHEN credential layers are described
- THEN real process environment has highest priority
- AND `.env.staging` is considered before `.env`

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
