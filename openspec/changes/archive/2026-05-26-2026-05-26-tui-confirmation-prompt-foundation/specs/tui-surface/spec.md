# TUI Surface Delta Specification

## ADDED Requirements

### Requirement: TUI Actions require impact-aware confirmation before risky execution

The TUI Surface MUST show an interactive confirmation prompt before executing a people-facing Action whose Capability metadata requires confirmation, and MUST pass explicit confirmation to the shared Pipeline only after the user confirms.

#### Scenario: Confirmation-required Action shows prompt before execution

- GIVEN a people-facing command-palette-visible Capability requires confirmation
- WHEN a user selects that Action from the Command Palette
- THEN the TUI shows a confirmation prompt containing the Action title, Capability ID, and risk level
- AND the Capability is not executed yet

#### Scenario: Confirmed Action executes with explicit confirmation

- GIVEN a confirmation prompt is open for a confirmation-required Action
- WHEN the user confirms the prompt
- THEN the TUI executes the selected Capability through the shared Pipeline with explicit confirmation
- AND the TUI shows loading while execution is in progress

#### Scenario: Confirmation prompt can be cancelled

- GIVEN a confirmation prompt is open for a confirmation-required Action
- WHEN the user cancels the prompt
- THEN the TUI closes the prompt
- AND the Capability is not executed
