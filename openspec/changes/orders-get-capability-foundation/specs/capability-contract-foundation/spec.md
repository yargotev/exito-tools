# Capability Contract Orders Delta Specification

## ADDED Requirements

### Requirement: Orders get capability exposes a neutral contract

The system MUST expose `orders.get` as a read-only Orders Domain Capability with schema-first input and domain-owned result models.

#### Scenario: Orders get definition is discoverable

- GIVEN `orders.get` is registered
- WHEN a surface inspects its Capability definition
- THEN the definition includes domain `orders`, version metadata, read-only risk, agents and people audiences, CLI/TUI/command-palette visibility, and a required string `id` input field

#### Scenario: Orders get default dependency is not configured

- GIVEN the Application has no real Orders API client yet
- WHEN `orders.get` is executed with valid input through the generic pipeline
- THEN execution returns a structured `ORDERS_NOT_CONFIGURED` failure envelope
