# orders.get-vtex

Gets an order by its VTEX OMS identifier through the VTEX OMS private order detail API.

## Contract

- **Capability ID**: `orders.get-vtex`
- **Domain**: Orders
- **CLI Command**: `exito orders get-vtex --id <order-id> [--brand exito|carulla]`
- **Risk Level**: read-only
- **Audience**: agents and people

## Input Schema

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `id` | string | yes | VTEX OMS order identifier, such as `1611511090420-01`. |
| `brand` | string | no | VTEX brand account to query. Defaults to `exito`; use `carulla` for Carulla. |

## Provider behavior

The Orders domain keeps VTEX OMS independent from the GEOMS-backed `orders.get` capability. The VTEX OMS getter sends server-side credentials only with these headers:

- `X-VTEX-API-AppKey`
- `X-VTEX-API-AppToken`

It calls `<vtexOms.<brand>.baseUrl>/api/oms/pvt/orders/<id>` and maps selected fields into a domain-owned `VTEXOMSOrder` while preserving the provider payload under `details` for diagnostics.

## Initial Result Shape

```json
{
  "order": {
    "id": "1611511090420-01",
    "sequence": "12345",
    "status": "ready-for-handling",
    "statusDescription": "Ready for handling",
    "creationDate": "2026-05-27T15:00:00Z",
    "clientName": "Ada Lovelace",
    "email": "ada@example.test",
    "totalValue": 123456,
    "brand": "exito",
    "details": {}
  }
}
```
