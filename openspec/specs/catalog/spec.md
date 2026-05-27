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

### Requirement: VTEX segment preparation capability

Exito Tools MUST expose a confirmation-required `catalog.create-vtex-segment` capability for explicitly preparing a VTEX segment from a caller-provided region ID and sales channel.

#### Scenario: Create segment from region and sales channel

- **Given** VTEX public sessions is configured for the requested brand
- **When** a caller executes `catalog.create-vtex-segment` with `regionId` set to `REGION_ID` and `salesChannel` set to `1`
- **Then** the capability MUST call `POST /io/api/sessions`
- **And** the request body MUST include `public.regionId.value` equal to `REGION_ID`
- **And** the request body MUST include `public.sc.value` equal to `1`
- **And** the result MUST include safe token metadata that indicates whether a segment token was returned.

#### Scenario: Segment creation requires confirmation

- **Given** `catalog.create-vtex-segment` mutates VTEX session state by creating a segment token
- **When** a caller executes it without explicit confirmation
- **Then** the shared Pipeline MUST return a structured `CONFIRMATION_REQUIRED` failure
- **And** the provider MUST NOT be called.

#### Scenario: Token diagnostics are redacted

- **Given** VTEX Sessions returns a segment token
- **When** `catalog.create-vtex-segment` succeeds or fails after receiving provider data
- **Then** diagnostics MUST NOT include the raw token value
- **And** provider payload token fields MUST be redacted.

#### Scenario: Optional cookie output is explicit

- **Given** VTEX Sessions returns a segment token
- **When** the caller sets `includeCookie` to `true`
- **Then** the result MAY include a `cookie` field formatted as `vtex_segment=<token>`
- **And** the capability MUST otherwise omit unredacted cookie output by default.

#### Scenario: CLI command emits JSON envelope

- **Given** the CLI is available
- **When** a caller runs `exito catalog create-vtex-segment --brand exito --region-id REGION_ID --sales-channel 1 --confirm`
- **Then** the command MUST execute `catalog.create-vtex-segment`
- **And** the command MUST emit the standard JSON envelope on stdout.
