# Capability Registry Orders Delta Specification

## MODIFIED Requirements

### Requirement: Registry accepts boot-time registration

The system MUST support explicit capability registration during application boot so Operational Domains can add capabilities through visible wiring.

#### Scenario: Orders get is available after finalization

- GIVEN application boot wires the Orders Domain
- WHEN the registry is finalized
- THEN `orders.get` is discoverable by stable Capability ID
