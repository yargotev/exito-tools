# Application Bootstrap Orders Delta Specification

## MODIFIED Requirements

### Requirement: Runnable application bootstrap

The system MUST build a runnable Go Application through explicit Application Wiring and MAY register the first real Operational Domain capability once its neutral contract exists.

#### Scenario: Orders get capability is wired during boot

- GIVEN the Application boots successfully
- WHEN application wiring finalizes the Capability Registry
- THEN the registry contains `orders.get`
- AND the registered entry has an executable handler
