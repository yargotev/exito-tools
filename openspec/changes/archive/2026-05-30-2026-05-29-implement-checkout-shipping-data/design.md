# Design: Implement Checkout Shipping Data

## Technical Approach

Extend the existing Checkout Domain pattern used by add-items and client-profile. `checkout.update-shipping-data` remains domain-owned, confirmation-gated by capability metadata, and surfaced through CLI as JSON input adapted into capability input.

## Architecture Decisions

| Decision | Choice | Alternatives considered | Rationale |
|----------|--------|-------------------------|-----------|
| Input shape | Domain structs for `selectedAddresses` and `logisticsInfo` | Raw `map[string]any` passthrough | Keeps validation/redaction domain-owned while matching VTEX attachment shape |
| Coordinate handling | Preserve supplied `geoCoordinates []float64` order | Rename to latitude/longitude fields | VTEX examples use array order; spec explicitly requires longitude before latitude preservation |
| Output | Reuse `OrderFormSummary`, add shipping total and SLA diagnostics | Echo shipping DTO | Avoids PII leakage and follows existing redacted-summary contract |
| CLI | `--input-json` object plus `--confirm` | Many address flags | Shipping payload is nested; JSON keeps CLI narrow and machine-first |

## Data Flow

    CLI JSON --decode--> capability.Input
      -> Pipeline confirmation gate
      -> Checkout use case validation
      -> HTTPClient POST /orderForm/{id}/attachments/shippingData
      -> VTEX orderForm DTO --redact/map--> OrderFormSummary

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `internal/domain/checkout/checkout.go` | Modify | Add constants, models, use case, capability, input mapping, validation |
| `internal/domain/checkout/http_client.go` | Modify | Add `UpdateShippingData` POST and map shipping diagnostics |
| `internal/domain/checkout/brand_client.go` | Modify | Route brand-specific shipping updater |
| `internal/app/app.go` | Modify | Register new capability and extend provider interfaces |
| `internal/surface/cli/root.go` | Modify | Add `checkout update-shipping-data` command |
| tests under `internal/...` | Modify | Domain, HTTP, app, and CLI coverage |

## Interfaces / Contracts

```go
const CapabilityUpdateShippingDataID = "checkout.update-shipping-data"

type UpdateShippingDataInput struct {
    Brand string
    OrderFormID string
    ShippingData ShippingDataInput
}
```

The JSON shape accepted by CLI/capability uses `selectedAddresses` and `logisticsInfo`. Each logistics item requires `itemIndex`, `selectedSla`, and `selectedDeliveryChannel`.

## Testing Strategy

| Layer | What to Test | Approach |
|-------|--------------|----------|
| Unit | Domain normalization/validation and redacted result | Go tests with recording client |
| Unit | HTTP path/body/diagnostics | `httptest.Server` |
| Unit | CLI confirmation and JSON adapter | Cobra root tests with fake registry |
| Unit | App registration | Registry discovery test |

## Migration / Rollout

No migration required. This adds a new capability/command without changing existing Checkout IDs.

## Open Questions

None.
