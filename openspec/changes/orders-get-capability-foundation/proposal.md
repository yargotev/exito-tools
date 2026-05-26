# Proposal: Orders Get Capability Foundation

## Summary

Add the first Orders Domain capability, `orders.get`, as a neutral executable Capability wired explicitly during application boot.

## Scope

- Add `internal/domain/orders` with domain-owned input, result, order model, use case, and getter interface.
- Expose `orders.get` metadata with read-only risk, agents/people audiences, CLI/TUI/palette visibility, and required string `id` input schema.
- Wire `orders.get` into `app.New` through explicit application wiring.
- Provide a default unavailable Orders dependency that returns a structured configuration error until real API infrastructure exists.
- Verify generic `exito run orders.get` reaches the bootstrapped capability and returns a standard not-configured envelope.

## Non-goals

- Call a real Orders API.
- Add explicit `exito orders get --id` domain command wiring.
- Add HTTP infrastructure, authentication, retries, or external DTO mapping.
- Implement Geo capability.
