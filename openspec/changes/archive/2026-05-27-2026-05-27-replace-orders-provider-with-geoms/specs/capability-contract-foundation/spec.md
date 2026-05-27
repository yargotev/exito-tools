## MODIFIED Requirements

### Requirement: orders.get domain execution

The system MUST execute `orders.get` through a domain-owned GEOMS client and return the stable `orders.GetResult` result shape.

#### Scenario: Configured GEOMS provider returns an order

- **Given** the Orders provider is configured with a GEOMS base URL and client credentials
- **And** the GEOMS token endpoint returns an access token with `expires_in`
- **And** GEOMS `findOrders` returns at least one order for order number `A123`
- **When** `orders.get` is executed with input `{"id":"A123"}`
- **Then** the result contains the mapped domain order
- **And** the provider request uses `filters.orderNumber` equal to `A123`
- **And** the provider request uses `filters.orderType` equal to `ExitoEcomm`
- **And** the system calls GEOMS `getOrder` with `order` equal to `A123`
- **And** the system calls GEOMS `findItemsByOrder` with `notFood` equal to `false` and `true`
- **And** the result includes `details`, `items.food`, and `items.notFood`

#### Scenario: Carulla order type is selected

- **Given** the Orders provider is configured
- **When** `orders.get` is executed with input `{"id":"A123","orderType":"CarullaEcomm"}`
- **Then** the GEOMS `findOrders` request uses `filters.orderType` equal to `CarullaEcomm`
