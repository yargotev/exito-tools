# TUI Surface Delta

## ADDED Requirements

### Requirement: Command Palette Actions collect required string input

The TUI Surface MUST collect required string inputs declared by a selected Action's Capability Input Schema before executing that Action.

#### Scenario: Selected Action opens an input form

- GIVEN a people-facing command-palette-visible Capability declares required string input fields
- WHEN the user selects it from the Command Palette
- THEN the TUI shows an input form for those fields instead of executing immediately

#### Scenario: Submitted form executes with collected input

- GIVEN a user has filled all required string fields in the TUI form
- WHEN the user submits the final field
- THEN the TUI executes the selected Capability through the shared Pipeline with those input values
