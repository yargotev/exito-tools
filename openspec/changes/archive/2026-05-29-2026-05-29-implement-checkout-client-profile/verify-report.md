# Verification Report: Checkout client profile attachment

## Summary

Implemented `checkout.update-client-profile` as the next Checkout purchase-assembly slice after add-items. The capability is registered, confirmation-gated, exposed through CLI, and backed by VTEX Checkout HTTP attachment mapping.

## Checks

- ✅ Checkout Domain owns client-profile use case, validation, result model, and capability metadata.
- ✅ `checkout.update-client-profile` is marked safe-write and requires confirmation.
- ✅ Missing orderForm ID or required client profile fields fail with `CHECKOUT_INVALID_INPUT` before provider execution.
- ✅ VTEX HTTP adapter posts client profile JSON to `/api/checkout/pub/orderForm/{orderFormId}/attachments/clientProfileData` and maps the redacted orderForm summary.
- ✅ CLI exposes `--input-json` for profile fields while keeping brand/orderForm as explicit flags.
- ✅ Output does not echo raw customer profile values; result shape returns only the redacted orderForm summary.
- ✅ No shipping, place-order, process-order, or payment capability was introduced.

## Verification Commands

- `go test ./internal/domain/checkout ./internal/surface/cli ./internal/app`
- `make test`
