# Checkout Specification

## Purpose

Define VTEX Checkout orderForm capabilities for guided purchase assembly before order placement.

## Requirements

### Requirement: Checkout domain owns VTEX orderForm assembly

Exito Tools MUST expose VTEX Checkout cart and orderForm operations through a dedicated Checkout Domain.

#### Scenario: Checkout is separate from Catalog and Orders

- **Given** product discovery is handled by Catalog capabilities
- **And** created order lookup is handled by Orders capabilities
- **When** a caller needs to create or mutate a VTEX orderForm
- **Then** Exito Tools MUST route that behavior through Checkout capabilities
- **And** Catalog search MUST NOT mutate cart state as a side effect
- **And** Orders lookup MUST NOT own pre-order cart state.

### Requirement: Get or create orderForm

Exito Tools MUST expose `checkout.get-order-form` as a read-only capability for retrieving a known VTEX orderForm and `checkout.create-order-form` as a confirmation-required safe-write capability for creating a new cart/orderForm.

#### Scenario: Create a fresh orderForm

- **Given** VTEX Checkout is configured for the requested profile and brand
- **When** a caller executes `checkout.create-order-form` with a sales channel and explicit confirmation
- **Then** Exito Tools MUST call VTEX Checkout current-cart creation with a new-cart request
- **And** the result MUST include the orderForm ID, sales channel, totals, items summary, client profile presence, shipping data presence, and safe diagnostics.

#### Scenario: Creation requires confirmation

- **Given** creating a VTEX orderForm mutates Checkout session state
- **When** a caller executes `checkout.create-order-form` without explicit confirmation
- **Then** the shared Pipeline MUST return `CONFIRMATION_REQUIRED`
- **And** VTEX Checkout MUST NOT be called.

#### Scenario: Retrieve an orderForm by ID

- **Given** VTEX Checkout is configured for the requested profile and brand
- **When** a caller executes `checkout.get-order-form` with an orderForm ID
- **Then** Exito Tools MUST retrieve that orderForm
- **And** the result MUST redact sensitive personal data unless the caller explicitly opts into approved PII output for a non-production diagnostic flow.

### Requirement: Add selected products to orderForm

Exito Tools MUST expose `checkout.add-items` as a confirmation-required safe-write capability for adding SKU items to an existing VTEX orderForm.

#### Scenario: Add products selected from Catalog search

- **Given** a caller has selected one or more SKU IDs from Catalog search results
- **And** a VTEX orderForm ID is available
- **When** the caller executes `checkout.add-items` with SKU IDs, quantities, sellers when needed, and explicit confirmation
- **Then** Exito Tools MUST call VTEX Checkout add cart items for the selected orderForm
- **And** the result MUST return the updated orderForm summary and item-level diagnostics.

#### Scenario: Reject empty item updates

- **Given** a caller executes `checkout.add-items` without any items
- **When** Exito Tools validates the input
- **Then** it MUST fail with stable structured error code `CHECKOUT_INVALID_INPUT`
- **And** VTEX Checkout MUST NOT be called.

### Requirement: Update client profile attachment

Exito Tools MUST expose `checkout.update-client-profile` as a confirmation-required safe-write capability for attaching customer profile data to a VTEX orderForm.

#### Scenario: Add customer profile data

- **Given** a VTEX orderForm ID is available
- **When** a caller executes `checkout.update-client-profile` with required customer fields and explicit confirmation
- **Then** Exito Tools MUST call the VTEX Checkout client profile attachment endpoint
- **And** the result MUST return the updated orderForm summary
- **And** stdout JSON, logs, and diagnostics MUST redact sensitive customer values by default.

### Requirement: Update shipping data and logistics selection

Exito Tools MUST expose `checkout.update-shipping-data` as a confirmation-required safe-write capability for attaching delivery address data and selected logistics options to a VTEX orderForm.

#### Scenario: Add address and selected delivery options

- **Given** a VTEX orderForm ID with items is available
- **When** a caller executes `checkout.update-shipping-data` with a delivery address, item logistics selections, and explicit confirmation
- **Then** Exito Tools MUST call the VTEX Checkout shipping data attachment endpoint
- **And** the request MUST preserve VTEX geoCoordinates order as longitude before latitude when coordinates are supplied
- **And** the result MUST return the updated orderForm summary, shipping total, and available/selected SLA diagnostics.

### Requirement: Checkout writes are sequential

Exito Tools MUST execute orderForm-mutating Checkout capabilities sequentially within a guided purchase assembly flow.

#### Scenario: No parallel orderForm mutation

- **Given** a Purchase Assembly Flow needs to create an orderForm, add items, update client profile data, and update shipping data
- **When** the flow executes
- **Then** each Checkout write MUST wait for the previous write result
- **And** the next write MUST use the latest returned orderForm state
- **And** Exito Tools MUST NOT issue parallel Checkout write requests for the same orderForm.

### Requirement: Final order placement is out of the first Checkout slice

Exito Tools MUST NOT place or pay a VTEX order as part of the first Checkout purchase assembly slice.

#### Scenario: Prepared orderForm stops before placement

- **Given** a Purchase Assembly Flow has prepared items, client profile data, shipping data, and logistics selections
- **When** the first Checkout slice completes
- **Then** the result MUST identify the orderForm as prepared for inspection
- **And** Exito Tools MUST NOT call place-order, process-order, or payment endpoints.
