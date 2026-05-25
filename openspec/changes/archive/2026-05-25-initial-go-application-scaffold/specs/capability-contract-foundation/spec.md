# Capability Contract Foundation Specification

## Purpose

Define shared foundation contracts.

## Requirements

### Requirement: Shared contract skeletons exist

The system MUST provide shared contract types for Capability metadata, Structured Errors, and JSON Envelope-shaped results so later slices can extend them without changing boundaries.

#### Scenario: Future capability can depend on shared types

- GIVEN a later slice adds a real capability
- WHEN it needs common metadata or failure shapes
- THEN it can use those shared contract types

#### Scenario: Runtime envelope emission remains deferred

- GIVEN the initial scaffold only implements root help
- WHEN no explicit machine-readable command runs
- THEN shared envelope types may exist without requiring JSON output

### Requirement: Visible contracts remain English-only

The system MUST keep user-facing labels and messages English-only.

#### Scenario: Shared CLI-facing text is English-only

- GIVEN the scaffold exposes root help or contract-facing messages
- WHEN a user reads them
- THEN the visible text is in English only
