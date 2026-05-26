# capability-registry Delta

## ADDED Requirements

### Requirement: Application boot registers Geo geocode-address

The application boot process SHALL explicitly register `geo.geocode-address` in the immutable Capability Registry.

#### Scenario: Generic run can route to Geo geocode-address

- **WHEN** the CLI runs `exito run geo.geocode-address` with complete `city` and `address` input
- **THEN** the shared execution Pipeline routes to the registered Capability and emits a standard JSON envelope with `meta.capabilityId` set to `geo.geocode-address`
