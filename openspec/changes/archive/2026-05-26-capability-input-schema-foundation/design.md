# Design: Capability Input Schema Foundation

## Approach

`capability.Definition` remains the neutral discovery contract. This slice adds an optional `InputSchema` pointer with ordered `InputField` entries. Each field has a stable name, a shared `InputType`, a required marker, and English description text.

The schema is intentionally smaller than full JSON Schema. It is enough for inventory discovery and future CLI/TUI adapters while avoiding premature validation complexity.

## JSON shape

Capability inventory serializes input metadata under `inputSchema.fields` only when a definition provides a schema. Field names use lower camel-case JSON keys to match the existing inventory contract.

## Registry immutability

Input schemas contain a slice of fields. The registry now deep-copies `InputSchema.Fields` during registration, finalization, `All`, and `Find`, preserving immutable runtime discovery semantics.
