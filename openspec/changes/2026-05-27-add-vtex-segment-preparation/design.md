# Design: Explicit VTEX segment preparation

## Capability

`catalog.create-vtex-segment` is a Catalog Domain capability because the resulting VTEX segment is primarily used to regionalize Catalog/Intelligent Search calls.

Metadata:

- Domain: `catalog`
- Risk: `safe-write`
- Requires confirmation: `true`
- Audience: agents first, people discoverable later through command palette if needed
- Visibility: CLI

## Inputs

- `brand`: `exito` or `carulla`, default `exito`.
- `regionId`: required VTEX region ID from a trusted upstream step.
- `salesChannel`: required VTEX sales channel / trade policy.
- `includeCookie`: optional boolean. When true, return `vtex_segment=<token>` in the result for immediate CLI copy/paste.

## Provider call

Use configured public VTEX brand base URL and call:

```http
POST /io/api/sessions
Content-Type: application/json

{
  "public": {
    "regionId": {"value": "REGION_ID"},
    "sc": {"value": "1"}
  }
}
```

The HTTP adapter maps provider response fields into domain-owned metadata. Expected token field is `segmentToken`; fallback parsing may inspect common token names if VTEX shape differs.

## Output and redaction

Stable result includes:

- `brand`
- `regionId`
- `salesChannel`
- `tokenSet`
- `tokenLength`
- optional `cookie` only when `includeCookie` is true
- `diagnostics.requestPath`
- `diagnostics.requestPayload` with only non-secret input values
- `diagnostics.providerPayload` with token-bearing fields redacted

Do not log, persist, or expose the token in diagnostics. The only unredacted token output is the explicit `cookie` field when `includeCookie` is true.

## CLI

`exito catalog create-vtex-segment --brand exito --region-id REGION --sales-channel 1 --confirm`

The command passes `Confirmed` to the shared Pipeline only when `--confirm` is present. Without confirmation, the command emits a standard `CONFIRMATION_REQUIRED` JSON envelope.
