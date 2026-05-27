# orders.get

Gets an order by its identifier through GEOMS `findOrders`, enriched with `getOrder` details and `findItemsByOrder` line items.

## Contract

- **Capability ID**: `orders.get`
- **Domain**: Orders
- **CLI Command**: `exito orders get --id <order-id> [--order-type ExitoEcomm|CarullaEcomm]`
- **TUI Action**: `Get order by ID`
- **Risk Level**: read-only
- **Audience**: agents and people

## Input Schema

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `id` | string | yes | Order identifier mapped to the GEOMS `orderNumber` filter. |
| `orderType` | string | no | GEOMS order type filter. Defaults to `ExitoEcomm`; use `CarullaEcomm` for Carulla. |

## Provider behavior

The Orders HTTP client obtains a GEOMS bearer token with Azure AD client credentials, using `expires_in` from the token response for cache expiry. It then posts to `<orders.baseUrl>/findOrders` with the selected order type, defaulting to `ExitoEcomm`, calls `<orders.baseUrl>/getOrder` for details, and calls `<orders.baseUrl>/findItemsByOrder` twice: `notFood: false` for food items and `notFood: true` for non-food items.

## Initial Result Shape

```json
{
  "order": {
    "id": "1557551083896",
    "status": "7500",
    "createdAt": "2026-05-19T18:16:01.651192",
    "customerName": "pepitos apellido",
    "email": "jcartagena@grupo-exito.com",
    "orderTotal": 980584,
    "statusOrderMax": "7500",
    "statusOrderMin": "7500",
    "items": {
      "food": [],
      "notFood": []
    },
    "details": {
      "infClient": {},
      "statusOrder": {},
      "paymentInformation": {},
      "paymentMethodInf": []
    }
  }
}
```
