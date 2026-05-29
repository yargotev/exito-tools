# Change: Implement Checkout orderForm base

## Motivation

The Checkout roadmap needs the first concrete VTEX orderForm slice so agents can create or load a cart/orderForm before later purchase assembly steps. The source-of-truth Checkout spec already defines `checkout.get-order-form` and `checkout.create-order-form`; this change records the implementation slice that delivered those contracts.

This record is retrospective because implementation commit `d0ccf4c` was created before the OpenSpec change artifact. The purpose of this change is to restore SDD traceability, document verification, and ensure the implementation remains bounded to the approved first slice.

## Scope

- Add `vtexCheckout` non-sensitive provider configuration by profile and brand.
- Add `internal/domain/checkout` with domain-owned inputs, use cases, result summaries, validation, and capability metadata.
- Register `checkout.get-order-form` as read-only.
- Register `checkout.create-order-form` as confirmation-required `safe-write`.
- Add a VTEX Checkout HTTP client for public orderForm create/load endpoints with fake-server tests.
- Add minimal CLI commands:
  - `exito checkout get-order-form --brand <brand> --order-form-id <id>`
  - `exito checkout create-order-form --brand <brand> --sales-channel <sc> --confirm`
- Keep stdout output as standard JSON envelopes and map provider DTOs into redacted domain-owned summaries.

## Out of Scope

- `checkout.add-items`.
- `checkout.update-client-profile`.
- `checkout.update-shipping-data`.
- Place-order, process-order, or payment flows.
- Browser cookie persistence, raw cookie output, or raw customer PII diagnostics.
- TUI primary navigation for Checkout.

## Safety and rollback

Checkout creation mutates provider-side cart/session state and therefore remains confirmation-gated by the shared Pipeline. Missing CLI confirmation returns `CONFIRMATION_REQUIRED` before the HTTP client is called.

Rollback is a normal revert of implementation commit `d0ccf4c` plus this OpenSpec record. No schema migrations or persisted local state are introduced.
