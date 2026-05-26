# TUI Surface Delta Specification

## MODIFIED Requirements

### Requirement: Command Palette Action execution uses shared Pipeline

The TUI Surface MUST execute selected Command Palette Actions through the shared Capability execution Pipeline rather than invoking domain use cases directly, and in-flight executions MUST be cancellable through the TUI Task Runner context.

#### Scenario: Selected palette Action succeeds

- GIVEN a people-facing command-palette-visible Capability can execute with the provided TUI input
- WHEN the user selects it from the Command Palette
- THEN the TUI shows loading while execution is in progress
- AND after completion shows a success state for that Capability

#### Scenario: Selected palette Action returns structured failure

- GIVEN a people-facing command-palette-visible Capability cannot execute with the provided TUI input
- WHEN the user selects it from the Command Palette
- THEN the TUI shows a failure state containing the structured error code and message

#### Scenario: In-flight Action can be cancelled

- GIVEN a selected Command Palette Action is running through the shared Pipeline
- WHEN the user presses `esc`
- THEN the TUI cancels the task execution context
- AND shows a cancelled state for that Capability
- AND a late completion message for the cancelled task does not replace the cancelled state
