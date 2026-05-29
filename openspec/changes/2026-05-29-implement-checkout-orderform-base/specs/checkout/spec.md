# Checkout Specification Delta

## ADDED Requirements

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
