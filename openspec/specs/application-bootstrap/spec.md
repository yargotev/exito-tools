# Application Bootstrap Specification

## Purpose

Define Go Application boot behavior and explicit wiring boundaries.

## Requirements

### Requirement: Runnable application bootstrap

The system MUST build a runnable Go Application through explicit Application Wiring and register implemented Operational Domain capabilities once their neutral contracts exist.

#### Scenario: Root application boots the CLI surface

- GIVEN the application has been built
- WHEN a user runs `exito`
- THEN the Application boots successfully and reaches the CLI Surface entrypoint

#### Scenario: Orders and Geo capabilities are wired during boot

- GIVEN the Application boots successfully
- WHEN application wiring finalizes the Capability Registry
- THEN the registry contains `orders.get`
- AND the registry contains `geo.geocode-address`
- AND both registered entries have executable handlers

### Requirement: Foundation boundaries stay explicit

The system MUST preserve documented package boundaries so surfaces depend on shared contracts while Operational Domains stay surface-independent.

#### Scenario: Surface dependency stays one-way

- GIVEN the scaffold packages exist
- WHEN future capabilities are added
- THEN Cobra concerns remain in surface packages rather than domain packages
