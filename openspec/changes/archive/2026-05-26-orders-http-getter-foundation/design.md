# Design

## Domain-owned provider client

`internal/domain/orders.HTTPGetter` owns the Orders provider DTOs and implements the existing `orders.Getter` interface. It uses `internal/platform/httpclient.Client` for base URL joining, bearer auth, JSON requests, bounded response decoding, and outbound request metadata headers.

## Provider request and response

The first narrow provider contract posts `{"id":"..."}` to `/orders/get`, matching the shared HTTP infrastructure foundation scenario. The response is decoded from the documented initial shape:

```json
{"order":{"id":"...","status":"...","createdAt":"..."}}
```

The DTO is immediately mapped to the domain-owned `orders.Order` so surfaces and use cases do not depend on provider JSON details.

## Application wiring

`app.New` keeps explicit dependency wiring. When `Effective.OrdersProvider.Configured` is false, it preserves `orders.UnavailableGetter{}` behavior. When configured, it builds `orders.NewHTTPGetter` from the resolved base URL and token.

## Error translation

The Orders domain translates provider failures into stable structured errors:

- Missing or invalid provider configuration: `ORDERS_NOT_CONFIGURED`.
- HTTP 404: `ORDER_NOT_FOUND`.
- Transport or non-2xx failures: `ORDERS_PROVIDER_UNAVAILABLE`.
- Invalid JSON: `ORDERS_PROVIDER_INVALID_RESPONSE`.
