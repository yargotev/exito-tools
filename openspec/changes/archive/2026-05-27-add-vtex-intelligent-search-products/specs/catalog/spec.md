# Catalog Specification Delta

## ADDED Requirements

### Requirement: VTEX Intelligent Search product capability

Exito Tools MUST expose a read-only `catalog.intelligent-search-products` capability for VTEX Intelligent Search REST product listings.

#### Scenario: Search by text with required trade policy

- **Given** VTEX Intelligent Search is configured for the requested brand
- **When** a caller executes `catalog.intelligent-search-products` with `tradePolicy` set to `1` and `text` set to `leche`
- **Then** the capability MUST call `GET /api/io/_v/api/intelligent-search/product_search/trade-policy/1`
- **And** the request query MUST include `query=leche`
- **And** the result MUST include domain-owned product and SKU summary fields
- **And** provider payloads MUST be preserved only under diagnostic details.

#### Scenario: Search by typed SKU lookup

- **Given** VTEX Intelligent Search is configured for the requested brand
- **When** a caller executes `catalog.intelligent-search-products` with `tradePolicy` set to `1`, `by` set to `sku-id`, and `value` containing `912350`
- **Then** the request query MUST include `query=sku.id:912350`.

#### Scenario: Search by repeated same-type IDs

- **Given** VTEX Intelligent Search is configured for the requested brand
- **When** a caller executes `catalog.intelligent-search-products` with `by` set to `sku-id` and multiple `value` entries
- **Then** the request query MUST join values with semicolons using one typed expression such as `sku.id:123;456`.

#### Scenario: Search with additional path facets

- **Given** VTEX Intelligent Search is configured for the requested brand
- **When** a caller executes `catalog.intelligent-search-products` with repeated `facet` entries
- **Then** the capability MUST encode each facet as escaped path segments after `trade-policy/<id>`
- **And** repeated keys MUST be preserved as repeated path segments.

#### Scenario: Search unavailable provider

- **Given** VTEX Intelligent Search is not configured for the requested brand
- **When** a caller executes `catalog.intelligent-search-products`
- **Then** the capability MUST fail with stable structured error code `CATALOG_NOT_CONFIGURED`.

#### Scenario: Reject ambiguous query modes

- **Given** the capability input includes more than one of `text`, raw `query`, or `by` plus `value`
- **When** a caller executes `catalog.intelligent-search-products`
- **Then** the capability MUST fail with stable structured error code `CATALOG_INVALID_INPUT`.

#### Scenario: Cookie values are redacted

- **Given** a caller provides a VTEX cookie for segment/session diagnostics
- **When** `catalog.intelligent-search-products` succeeds or fails
- **Then** Exito Tools MUST NOT include the cookie value in stdout JSON, structured diagnostics, or logs.
