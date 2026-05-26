# TUI Surface Delta Specification

## ADDED Requirements

### Requirement: Command Palette discovers people-facing Actions

The TUI Surface MUST provide a Command Palette discovery mode that lists people-facing Actions across domains based on Capability metadata.

#### Scenario: Palette includes command-palette-visible people Actions

- GIVEN a Capability has people audience and command-palette visibility
- WHEN the Command Palette is opened
- THEN that Capability appears as a discoverable Action

#### Scenario: Palette filters Actions by query

- GIVEN the Command Palette is open with multiple discoverable Actions
- WHEN a user types a query matching an Action title or Capability ID
- THEN only matching Actions are shown

#### Scenario: Palette is distinct from primary navigation

- GIVEN a Capability has command-palette visibility but not TUI primary visibility
- WHEN the initial shell is rendered
- THEN that Capability is not shown as a primary Action
- WHEN the Command Palette is opened
- THEN that Capability is shown as a discoverable Action
