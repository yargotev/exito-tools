# Checkout Specification Delta

## ADDED Requirements

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
