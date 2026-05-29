# Sequential confirmation-gated Checkout writes

VTEX Checkout write capabilities in Exito Tools will execute sequentially and require explicit confirmation before provider mutation.

Checkout API writes update shared orderForm/session state. Exito Tools must not issue parallel orderForm mutations inside one purchase assembly flow because competing writes can overwrite older values or fail due to provider-side race conditions. Each write capability should return the updated VTEX orderForm summary so the next step starts from the latest state.

The first Checkout roadmap slice stops at purchase assembly: create/load orderForm, add items, update client profile data, update shipping data/logistics selections, and inspect the prepared cart. Final order placement and payment processing are separate higher-risk capabilities and are not included until the user explicitly approves that slice.
