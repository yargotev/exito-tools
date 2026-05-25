# orders.get

Gets an order by its identifier.

## Contract

- **Capability ID**: `orders.get`
- **Domain**: Orders
- **CLI Command**: `exito orders get --id <order-id>`
- **TUI Action**: `Get order by ID`
- **Risk Level**: read-only
- **Audience**: agents and people

## Input Schema

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `id` | string | yes | Order identifier. |

## Initial Result Shape

```json
{
  "order": {
    "id": "...",
    "status": "...",
    "createdAt": "..."
  }
}
```
