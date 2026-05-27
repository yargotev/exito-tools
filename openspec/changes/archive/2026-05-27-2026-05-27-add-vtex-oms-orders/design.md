# Design

`orders.get-vtex` lives in the Orders domain because OMS order detail is order business behavior; VTEX is treated as an external provider, not as a first-class operational domain.

The Orders domain owns the VTEX OMS HTTP getter and maps the provider response into a domain-owned `VTEXOMSOrder` result. The getter calls `GET /api/oms/pvt/orders/<id>` and sends `X-VTEX-API-AppKey` and `X-VTEX-API-AppToken` only from server-side configuration.

Configuration adds `VTEXOMSProvider` alongside the existing GEOMS Orders provider. YAML may hold non-sensitive `profiles.<profile>.vtexOms.<brand>.baseUrl`; credentials are resolved from process environment, `.env.<profile>`, then `.env`. Non-production profiles use the QA variable names, while `prod`/`production`/`pdn` use production variable names.
