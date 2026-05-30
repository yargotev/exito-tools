# Delta for Checkout

## ADDED Requirements

### Requirement: Checkout shipping data implementation

Exito Tools MUST implement `checkout.update-shipping-data` as a confirmation-required safe-write capability for attaching delivery address data and selected logistics options to an existing VTEX orderForm.

#### Scenario: CLI attaches shipping data and selected SLA

- **Given** VTEX Checkout is configured for the requested profile and brand
- **And** a caller has an orderForm ID, address data, and logistics selections
- **When** the caller runs `exito checkout update-shipping-data --brand exito --order-form-id <id> --input-json <shipping-json> --confirm`
- **Then** Exito Tools MUST execute `checkout.update-shipping-data` through the shared Pipeline
- **And** the VTEX Checkout client MUST call the orderForm shipping data attachment endpoint
- **And** the result MUST include an updated redacted orderForm summary with shipping total and selected SLA diagnostics.

#### Scenario: Shipping update requires confirmation

- **Given** attaching shipping data mutates provider-side orderForm state
- **When** a caller executes `checkout.update-shipping-data` without explicit confirmation
- **Then** the shared Pipeline MUST return `CONFIRMATION_REQUIRED`
- **And** VTEX Checkout MUST NOT be called.

#### Scenario: Reject incomplete shipping data

- **Given** a caller executes `checkout.update-shipping-data` without an orderForm ID, selected address, or logistics selection
- **When** Exito Tools validates the input
- **Then** it MUST fail with stable structured error code `CHECKOUT_INVALID_INPUT`
- **And** VTEX Checkout MUST NOT be called.

#### Scenario: Preserve VTEX geoCoordinates order

- **Given** a shipping address includes `geoCoordinates`
- **When** Exito Tools sends the shipping data attachment request
- **Then** the request MUST preserve the input array in VTEX order: longitude before latitude.

#### Scenario: Redact shipping address values

- **Given** VTEX Checkout returns shipping data with receiver, street, postal code, or coordinate values
- **When** Exito Tools emits the CLI JSON envelope
- **Then** stdout MUST expose only safe presence, total, and SLA diagnostics
- **And** stdout MUST NOT echo raw shipping address values.
