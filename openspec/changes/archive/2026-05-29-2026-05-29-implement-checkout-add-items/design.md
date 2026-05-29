# Design: Checkout add-items

## Domain

The Checkout Domain gains an `AddItemsUseCase` and an `Adder` port. Inputs are normalized inside the domain: blank brand defaults to `exito`, item SKU IDs and sellers are trimmed, and quantities must be positive integers. Empty item lists fail with `CHECKOUT_INVALID_INPUT`.

## Provider

The VTEX Checkout HTTP adapter calls:

`POST /api/checkout/pub/orderForm/{orderFormId}/items?allowedOutdatedData=false`

with an `orderItems` body containing item `id`, `quantity`, `seller`, and zero-based `index`. The adapter maps the returned provider orderForm into the existing redacted `OrderFormSummary`.

## Surfaces

The CLI accepts repeated `--item` flags in the form `sku=<sku>,quantity=<qty>[,seller=<seller>]`. Seller defaults to `1` when omitted so common SKU selections can be added without a seller lookup.

## Risks

Adding items mutates provider-side Checkout session/orderForm state. The capability is marked safe-write and requires explicit confirmation through the shared execution pipeline.
