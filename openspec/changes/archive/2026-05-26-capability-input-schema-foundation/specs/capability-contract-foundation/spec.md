# Capability Contract Input Schema Delta Specification

## ADDED Requirements

### Requirement: Capability definitions expose neutral input schemas

The system MUST model schema-first Capability input metadata as neutral contracts so interaction surfaces and agents can discover complete input object requirements without redefining them per surface.

#### Scenario: Capability input schema is serialized in inventory output

- GIVEN a Capability definition includes an input schema with fields
- WHEN the CLI Surface emits the machine-readable capabilities inventory
- THEN the JSON definition includes `inputSchema.fields` with field name, type, required marker, and description

### Requirement: Capability input schemas use shared primitive categories

The system MUST provide stable input type categories for common JSON-shaped values so future Operational Domains avoid ad-hoc type strings.

#### Scenario: Future domains use shared input type constants

- GIVEN a future domain registers a Capability with an input schema
- WHEN it declares a string, number, boolean, object, or array input field
- THEN it can use shared capability input type constants
