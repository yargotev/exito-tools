# TUI Surface Specification

## Purpose

Define the interactive TUI Surface behavior for task-first, people-facing workflows.

## Requirements

### Requirement: Initial TUI shell shows profile and primary actions

The TUI Surface MUST start from the shared Application and render a task-first shell that shows the effective Profile and curated primary Actions.

#### Scenario: People-facing TUI-visible capability is shown

- GIVEN a Capability has people audience and TUI visibility
- WHEN the initial TUI model is rendered
- THEN that Capability appears as a primary Action

#### Scenario: Agent-only capability is not promoted as primary

- GIVEN a Capability has only agent audience
- WHEN the initial TUI model is rendered
- THEN that Capability is not shown as a primary Action

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

### Requirement: Command Palette Action execution uses shared Pipeline

The TUI Surface MUST execute selected Command Palette Actions through the shared Capability execution Pipeline rather than invoking domain use cases directly.

#### Scenario: Selected palette Action succeeds

- GIVEN a people-facing command-palette-visible Capability can execute with the provided TUI input
- WHEN the user selects it from the Command Palette
- THEN the TUI shows loading while execution is in progress
- AND after completion shows a success state for that Capability

#### Scenario: Selected palette Action returns structured failure

- GIVEN a people-facing command-palette-visible Capability cannot execute with the provided TUI input
- WHEN the user selects it from the Command Palette
- THEN the TUI shows a failure state containing the structured error code and message
