# Design: Geo geocode-address CLI command

## Approach

Follow the existing explicit Orders command pattern:

- Add `newGeoCommand` as the `geo` command group.
- Add `newGeoGeocodeAddressCommand` for `geocode-address`.
- Map `--city` and `--address` flags into a complete `capability.Input` object.
- Execute `geo.CapabilityGeocodeAddressID` through `execution.Pipeline`.
- Write the standard JSON envelope via `presenter.WriteJSON`.

The command will currently return `GEO_NOT_CONFIGURED` with standard metadata because the real provider client is deferred, but this verifies command routing and shared execution behavior.
