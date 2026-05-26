# Proposal: Orders Get CLI Command

## Summary

Expose the first explicit domain command, `exito orders get --id <order-id>`, as a convenient CLI adapter over the existing neutral `orders.get` Capability.

## Scope

- Add an `orders` command group to the CLI Surface.
- Add `orders get --id <order-id>` with required ID flag.
- Execute the registered `orders.get` Capability through the shared execution Pipeline.
- Emit the same standard JSON Envelope metadata as generic `exito run`.
- Keep the command wired to the current unavailable Orders dependency until a real Orders client exists.

## Non-goals

- Add a real Orders API client.
- Change the `orders.get` domain use case or result contract.
- Add non-JSON output formats.
- Add more Orders commands.
