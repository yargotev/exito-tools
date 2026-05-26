# CLI Root Orders Get Delta Specification

## ADDED Requirements

### Requirement: Orders get command is exposed explicitly

The CLI Surface MUST expose `exito orders get --id <order-id>` as an explicit domain command for the `orders.get` Capability.

#### Scenario: Orders get runs through the shared pipeline

- GIVEN the Application has registered `orders.get`
- WHEN a user runs `exito orders get --id A123`
- THEN the command executes the `orders.get` Capability through the shared execution Pipeline
- AND stdout contains a standard JSON Envelope
- AND `meta.capabilityId` is `orders.get`

#### Scenario: Orders get requires an ID flag

- GIVEN the Orders get command is available
- WHEN a user runs `exito orders get` without `--id`
- THEN the command fails before execution
- AND it does not emit a JSON success envelope
