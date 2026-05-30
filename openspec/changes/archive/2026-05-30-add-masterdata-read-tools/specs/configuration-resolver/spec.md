# Delta for Configuration Resolver

## ADDED Requirements

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
