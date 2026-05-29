# Checkout domain for VTEX orderForm assembly

Exito Tools will model VTEX Checkout orderForm and cart-assembly behavior as a dedicated Checkout Domain, separate from Orders and Catalog.

Catalog remains responsible for product discovery and SKU selection. Orders remains responsible for already-created order lookups through GEOMS or VTEX OMS. Checkout owns creating or loading VTEX orderForms, adding/updating cart items, updating checkout attachments such as client profile data and shipping data, and later place-order/payment steps if explicitly approved.

This keeps stateful cart mutation out of read-only product search, avoids overloading Orders with pre-order cart state, and gives interaction surfaces a clear place to apply safe-write confirmation, PII redaction, cookie handling, and sequential execution rules for VTEX Checkout API writes.
