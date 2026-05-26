# Design: Capability Metadata Foundation

## Approach

The `capability.Definition` type remains the neutral source of truth for discovery metadata. This slice adds typed string aliases and constants for fields that are already documented in ADRs:

- `Domain` and `Version` support stable capability contracts and compatible version metadata.
- `Risk` and `RequiresConfirmation` give surfaces the information needed to handle safe-write or destructive flows later.
- `Audiences` and `Visibility` let surfaces decide where to promote capabilities without hiding machine-accessible capabilities from generic CLI discovery.

## Registry immutability

Adding slice-backed metadata means the previous shallow-copy registry behavior is no longer sufficient. Registry registration, finalization, `All`, and `Find` now clone definitions so outside callers cannot mutate audience or visibility slices after boot.

## JSON shape

Metadata fields use lower camel-case JSON names and `omitempty` for optional fields so existing empty foundation capabilities do not produce noisy zero values.
