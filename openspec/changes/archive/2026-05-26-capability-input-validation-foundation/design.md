# Design: Capability Input Validation Foundation

## Approach

The execution pipeline is the first shared enforcement point for schema-shaped inputs because both generic CLI execution and future surfaces pass through it. After registry lookup and before handler invocation, the pipeline validates the request input against the registered Capability definition's optional `InputSchema`.

Validation is deliberately narrow:

- Missing or `nil` required fields fail.
- Present fields must match the declared primitive JSON-shaped `InputType`.
- Optional missing fields are allowed.
- Unknown extra input fields are ignored for now.
- Capabilities without schemas preserve previous behavior.

## Error contract

Invalid inputs return a standard failed JSON-envelope-shaped result with `error.code = INVALID_INPUT`. The handler is not invoked when validation fails, keeping domain use cases focused on business behavior instead of surface/input-shape checks.
