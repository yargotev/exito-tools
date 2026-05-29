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

### Requirement: Checkout orderForm base implementation

Exito Tools MUST implement the first Checkout orderForm slice with only `checkout.get-order-form` and `checkout.create-order-form` before adding broader purchase assembly mutations.

#### Scenario: First slice registers only base orderForm capabilities

- **Given** the application boots with the Checkout Domain enabled
- **When** a caller inspects capability discovery
- **Then** `checkout.get-order-form` MUST be registered
- **And** `checkout.create-order-form` MUST be registered
- **And** Checkout item, profile, shipping, place-order, and payment capabilities MUST NOT be introduced by this first slice.

#### Scenario: Create orderForm is confirmation-gated

- **Given** `checkout.create-order-form` creates provider-side Checkout state
- **When** a caller executes the capability without explicit confirmation
- **Then** the shared Pipeline MUST return `CONFIRMATION_REQUIRED`
- **And** the VTEX Checkout provider MUST NOT be called.

#### Scenario: OrderForm output is redacted by default

- **Given** VTEX Checkout returns an orderForm with client profile or shipping data
- **When** Exito Tools maps the provider response to the capability result
- **Then** the result MUST expose safe presence flags for client profile and shipping data
- **And** the result MUST NOT expose raw customer PII, raw cookies, ownership cookie values, or provider credentials.

#### Scenario: CLI commands use standard JSON envelopes

- **Given** the CLI Surface is available
- **When** a caller runs `exito checkout get-order-form` or `exito checkout create-order-form`
- **Then** the command MUST execute the registered Checkout capability through the shared Pipeline
- **And** stdout MUST contain the standard JSON envelope only.

### Requirement: Checkout add-items implementation

Exito Tools MUST implement `checkout.add-items` as a confirmation-required safe-write capability for adding selected SKU items to an existing VTEX orderForm.

#### Scenario: CLI adds selected SKU items

- **Given** VTEX Checkout is configured for the requested profile and brand
- **And** a caller has an orderForm ID and one or more selected SKU IDs
- **When** the caller runs `exito checkout add-items --brand exito --order-form-id <id> --item sku=<sku>,quantity=<qty>[,seller=<seller>] --confirm`
- **Then** Exito Tools MUST execute `checkout.add-items` through the shared Pipeline
- **And** the VTEX Checkout client MUST call the orderForm items endpoint with the selected items
- **And** the result MUST include the updated redacted orderForm summary and item-level diagnostics.

#### Scenario: Add-items requires confirmation

- **Given** adding items mutates provider-side orderForm state
- **When** a caller executes `checkout.add-items` without explicit confirmation
- **Then** the shared Pipeline MUST return `CONFIRMATION_REQUIRED`
- **And** VTEX Checkout MUST NOT be called.

#### Scenario: Reject empty add-items input

- **Given** a caller executes `checkout.add-items` without any items
- **When** Exito Tools validates the input
- **Then** it MUST fail with stable structured error code `CHECKOUT_INVALID_INPUT`
- **And** VTEX Checkout MUST NOT be called.
