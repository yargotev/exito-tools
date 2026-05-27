# Catalog Specification

## Purpose

Define Catalog Domain capabilities and VTEX storefront product search behavior.

## Requirements

### Requirement: VTEX catalog product search capability

Exito Tools MUST expose a read-only `catalog.search-products` capability for VTEX public catalog product search.

#### Scenario: Search by a simple identifier

- **Given** VTEX Catalog is configured for the requested brand
- **When** a caller executes `catalog.search-products` with `by` set to `sku-id` and `value` set to a SKU ID
- **Then** the capability MUST query VTEX Search API using `fq=skuId:{value}`
- **And** the result MUST include matching products with stable product and SKU summary fields
- **And** each product result MUST preserve the provider payload under diagnostic details

#### Scenario: Search with advanced filters

- **Given** VTEX Catalog is configured for the requested brand
- **When** a caller executes `catalog.search-products` with one or more `fq` values
- **Then** the capability MUST send each filter as a repeated VTEX `fq` query parameter
- **And** optional `ft`, `order`, `from`, and `to` inputs MUST be forwarded when present

#### Scenario: Search unavailable provider

- **Given** VTEX Catalog is not configured for the requested brand
- **When** a caller executes `catalog.search-products`
- **Then** the capability MUST fail with a stable structured error code `CATALOG_NOT_CONFIGURED`

#### Scenario: CLI command exposes simple and advanced modes

- **Given** the CLI is available
- **When** a caller runs `exito catalog search-products --by sku-id --value 912350 --brand exito`
- **Then** the command MUST execute `catalog.search-products`
- **And** the command MUST emit the standard JSON envelope on stdout

#### Scenario: Pagination metadata is preserved

- **Given** VTEX Search API returns a `resources` header
- **When** `catalog.search-products` succeeds
- **Then** the result MUST expose the returned range and total count when parsable
