# Design: Generic Run Command Foundation

## Approach

The CLI surface owns only argument/input adaptation. `exito run <capability-id>` parses a complete JSON object into `capability.Input`, boots the Application with the shared root flags, then invokes `execution.Pipeline.Execute` with the resolved profile and optional correlation ID.

## Input handling

The command supports exactly one input source:

1. `--input-json` inline JSON object.
2. `--input-file` JSON object file.
3. stdin when data is piped or tests set a non-default input reader.

If no input is supplied, the command passes an empty object. This keeps zero-input capabilities invocable while still aligning the generic command with schema-shaped input objects.

## Output handling

The command writes the Pipeline envelope through the shared JSON presenter. Unknown Capability IDs are represented as structured `CAPABILITY_NOT_FOUND` envelopes from the Pipeline rather than as Cobra usage errors.
