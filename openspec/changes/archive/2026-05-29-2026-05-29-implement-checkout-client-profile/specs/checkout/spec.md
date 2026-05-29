# Checkout Specification Delta

## ADDED Requirements

### Requirement: Checkout client profile implementation

Exito Tools MUST implement `checkout.update-client-profile` as a confirmation-required safe-write capability for attaching customer profile data to an existing VTEX orderForm.

#### Scenario: CLI attaches client profile data

- **Given** VTEX Checkout is configured for the requested profile and brand
- **And** a caller has an orderForm ID and customer profile data
- **When** the caller runs `exito checkout update-client-profile --brand exito --order-form-id <id> --input-json <profile-json> --confirm`
- **Then** Exito Tools MUST execute `checkout.update-client-profile` through the shared Pipeline
- **And** the VTEX Checkout client MUST call the orderForm client profile attachment endpoint
- **And** the result MUST include the updated redacted orderForm summary
- **And** the result MUST NOT echo raw customer profile values.

#### Scenario: Client-profile update requires confirmation

- **Given** attaching client profile data mutates provider-side orderForm state
- **When** a caller executes `checkout.update-client-profile` without explicit confirmation
- **Then** the shared Pipeline MUST return `CONFIRMATION_REQUIRED`
- **And** VTEX Checkout MUST NOT be called.

#### Scenario: Reject incomplete client profile input

- **Given** a caller executes `checkout.update-client-profile` without an orderForm ID or required profile fields
- **When** Exito Tools validates the input
- **Then** it MUST fail with stable structured error code `CHECKOUT_INVALID_INPUT`
- **And** VTEX Checkout MUST NOT be called.
