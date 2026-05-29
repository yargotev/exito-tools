# VTEX Checkout orderForm research

## Sources

- VTEX Checkout API overview: https://developers.vtex.com/docs/guides/checkout-api-overview
- VTEX orderForm fields: https://developers.vtex.com/docs/guides/orderform-fields
- Get current or create a new cart: https://developers.vtex.com/vtex-rest-api/docs/create-a-new-cart/
- Add cart items: https://developers.vtex.com/vtex-rest-api/docs/add-cart-items/
- Headless cart and checkout: https://developers.vtex.com/docs/guides/headless-cart-and-checkout
- Create a regular order from an existing cart: https://developers.vtex.com/docs/guides/create-a-regular-order-from-an-existing-cart

## Findings

- VTEX Checkout centers cart state on the `orderForm`, which carries items, client profile data, shipping data, payment data, totals, and other checkout context.
- A current cart can be retrieved through Checkout, and a fresh cart/orderForm can be created by requesting a new cart.
- Adding items, client profile data, shipping address/logistics selections, marketing data, payment data, and related attachments are separate Checkout operations that usually return an updated orderForm.
- VTEX documents that Checkout data modification operations must not be performed in parallel; clients should enqueue them to avoid overwrites or race/competition errors.
- Headless flows may rely on Checkout cookies such as orderForm ownership to access unmasked personal data. Exito Tools should treat cookie values and PII as sensitive execution state.
- Creating a regular order from an existing cart continues beyond orderForm assembly into place-order, payment, and processing APIs. Those steps are higher risk and should be separated from the first purchase assembly slice.
