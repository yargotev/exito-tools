# Verification Report: Checkout add-items

## Summary

Implemented `checkout.add-items` as the next Checkout purchase-assembly slice after the orderForm base. The capability is registered, confirmation-gated, exposed through CLI, and backed by VTEX Checkout HTTP request mapping.

## Checks

- ✅ Checkout Domain owns add-items use case, validation, result model, and capability metadata.
- ✅ `checkout.add-items` is marked safe-write and requires confirmation.
- ✅ Empty item lists, blank SKU, non-positive quantity, blank seller, blank orderForm ID, and unsupported brands fail with `CHECKOUT_INVALID_INPUT`.
- ✅ VTEX HTTP adapter posts `orderItems` to `/api/checkout/pub/orderForm/{orderFormId}/items?allowedOutdatedData=false` and maps the redacted orderForm summary.
- ✅ CLI exposes repeated `--item sku=<sku>,quantity=<qty>[,seller=<seller>]` flags and defaults seller to `1` when omitted.
- ✅ No client-profile, shipping, place-order, process-order, or payment capability was introduced.

## Verification Commands

- `go test ./internal/domain/checkout ./internal/surface/cli ./internal/app`
- `make test`
