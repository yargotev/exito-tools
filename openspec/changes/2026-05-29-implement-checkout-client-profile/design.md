# Design: Checkout client profile attachment

## Domain

The Checkout Domain gains an `UpdateClientProfileUseCase` and a `ClientProfileUpdater` port. Inputs are normalized inside the domain: blank brand defaults to `exito`, the orderForm ID and profile string fields are trimmed, and required profile fields must be present before any provider call.

The result returns only the existing redacted `OrderFormSummary`; it does not echo customer email, document, phone, or names.

## Provider

The VTEX Checkout HTTP adapter calls:

`POST /api/checkout/pub/orderForm/{orderFormId}/attachments/clientProfileData`

with a JSON body containing `email`, `firstName`, `lastName`, `documentType`, `document`, and `phone`. The adapter maps the returned provider orderForm into the existing redacted `OrderFormSummary`.

## Surfaces

The CLI accepts `--input-json` as an inline JSON object containing the client profile fields. Brand and orderForm ID remain explicit flags so the command shape matches the roadmap and keeps routing fields separate from sensitive profile data.

## Risks

Client profile attachment mutates provider-side Checkout session/orderForm state and carries PII. The capability is marked safe-write, requires explicit confirmation through the shared execution pipeline, and returns redacted output only.
