# Application Bootstrap Specification

## Purpose

Define scaffold boot behavior.

## Requirements

### Requirement: Runnable application bootstrap

The system MUST build a runnable Go Application through explicit Application Wiring before real Operational Domains are required.

#### Scenario: Root application boots the CLI surface

- GIVEN the scaffold has been built
- WHEN a user runs `exito`
- THEN the Application boots successfully and reaches the CLI Surface entrypoint

#### Scenario: No real business domains are required yet

- GIVEN the foundation slice is installed
- WHEN the Application boots
- THEN `orders.get` and `geo.geocode-address` remain deferred and do not block startup

### Requirement: Foundation boundaries stay explicit

The system MUST preserve documented package boundaries so surfaces depend on shared contracts while Operational Domains stay surface-independent.

#### Scenario: Surface dependency stays one-way

- GIVEN the scaffold packages exist
- WHEN future capabilities are added
- THEN Cobra concerns remain in surface packages rather than domain packages
