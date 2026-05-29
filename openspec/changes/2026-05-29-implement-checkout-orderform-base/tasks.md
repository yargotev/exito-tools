# Tasks

- [x] Add OpenSpec change record for the first Checkout orderForm implementation slice.
- [x] Add `vtexCheckout` provider configuration resolution by profile and brand.
- [x] Add Checkout Domain package with use cases, validation, result summaries, and capability metadata.
- [x] Implement `checkout.get-order-form` as read-only.
- [x] Implement `checkout.create-order-form` as confirmation-required safe-write.
- [x] Implement VTEX Checkout HTTP client with fake-server tests.
- [x] Add minimal CLI commands for get/create orderForm.
- [x] Wire capability registration/discovery through explicit app wiring.
- [x] Keep add-items, profile, shipping, place-order, and payment out of scope.
- [x] Validate PII/cookies are not exposed in domain summaries or config serialization.
- [x] Run `make precommit` for implementation commit `d0ccf4c`.
- [x] Re-run validation after adding this OpenSpec record.
