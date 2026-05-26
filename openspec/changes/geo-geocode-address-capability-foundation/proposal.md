# Proposal: Geo geocode-address Capability foundation

## Summary

Add the first Geo Domain capability foundation for `geo.geocode-address`, matching ADR 0051 and the documented capability contract while deferring the real provider client and explicit `exito geo geocode-address` command to later slices.

## Scope

- Add an internal Geo Domain package with domain-owned input/result models.
- Add a `Geocoder` seam and `GeocodeAddressUseCase` with no surface dependencies.
- Register an executable `geo.geocode-address` Capability with metadata and input schema.
- Wire the capability explicitly during application boot with a default unavailable dependency.
- Cover domain execution, application wiring, and generic `exito run` behavior.

## Out of Scope

- Real HTTP provider integration.
- Dotenv/YAML parsing of Geo credentials.
- Explicit `exito geo geocode-address` CLI command.
- TUI action implementation.
