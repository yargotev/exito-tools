# TUI Surface Delta Specification

## ADDED Requirements

### Requirement: Successful task results can be filtered locally

The TUI Surface MUST provide a local Result Filter mode for successful task results that refines already loaded result rows without executing the Capability again.

#### Scenario: Result filter refines loaded task data

- GIVEN a Command Palette Action has completed successfully with result data
- WHEN the user opens the Result Filter and types a query matching one result row
- THEN the TUI shows the active filter query
- AND only matching loaded result rows are displayed
- AND the Capability is not executed again

#### Scenario: Result filter is distinct from Command Palette

- GIVEN a successful task result is displayed
- WHEN the user opens the Result Filter
- THEN the Command Palette remains closed
- AND the displayed filter is labelled as a Result Filter
