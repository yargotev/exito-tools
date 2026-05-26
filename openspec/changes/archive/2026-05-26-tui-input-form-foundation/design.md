# Design: TUI input form foundation

## Approach

Extend the TUI model with a small form state that is entered when a selected Command Palette Action has required string fields in its neutral Input Schema. The form stores the target Capability ID, ordered fields, current cursor, and typed values.

Pressing enter in the Command Palette starts the form when input is required; otherwise it keeps the existing immediate execution behavior. Pressing enter in a field advances to the next field, and pressing enter on the final field submits the collected `capability.Input` through the same Pipeline command used by direct action execution.

## Decisions

- Only required string fields are supported in this slice.
- Field order follows Input Schema order.
- Empty field submission stays in place instead of calling the Pipeline, avoiding a known required-input failure for simple forms.
- The form is Action-scoped and replaces the palette view while active.

## Risks

- Capabilities with required non-string fields still cannot be filled from the TUI and will use existing failure behavior until richer form controls are added.
- There is no persistent edit history or advanced cursor movement yet.
