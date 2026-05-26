# Design: Orders Get Capability Foundation

## Approach

The Orders Domain lives under `internal/domain/orders` to preserve the documented package boundary between domains and surfaces. It exposes a `Getter` interface, a `GetUseCase`, domain-owned `Order`/`GetResult` models, and `NewGetCapability` for application wiring.

`app.New` explicitly registers `orders.NewGetCapability(orders.UnavailableGetter{})`. This makes `orders.get` discoverable and runnable through the generic pipeline without introducing fake data or real network dependencies. Until the real Orders client exists, valid executions fail with a structured `ORDERS_NOT_CONFIGURED` error.

## Capability contract

`orders.get` uses the documented stable ID, read-only risk, both agents and people audiences, CLI/TUI/command-palette visibility, and one required string input field: `id`.

## Deferred work

A later slice can replace `UnavailableGetter` with a real domain API client and add the explicit `exito orders get --id <order-id>` command. This slice keeps domain logic surface-independent and avoids committing to external DTOs prematurely.
