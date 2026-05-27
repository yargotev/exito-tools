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


### Requirement: Primary TUI Actions are keyboard navigable

The TUI Surface MUST let users navigate primary Actions with arrow keys and Vim-style `j`/`k` keys, and MUST run the selected primary Action through the same input, confirmation, and Pipeline path used by Command Palette Actions.

#### Scenario: Primary Action selection moves with keyboard

- GIVEN multiple people-facing TUI-visible Capabilities are shown as primary Actions
- WHEN the user presses `down` or `j`
- THEN the selected primary Action moves to the next Action
- WHEN the user presses `up` or `k`
- THEN the selected primary Action moves to the previous Action

#### Scenario: Selected primary Action runs through shared execution path

- GIVEN a primary Action is selected
- WHEN the user presses `enter`
- THEN the TUI either collects that Action's string inputs or asks for required confirmation when applicable
- AND the TUI executes the Action through the shared Capability execution Pipeline after required inputs and confirmation are satisfied

### Requirement: TUI uses Catppuccin Mocha visual language

The TUI Surface MUST render the shell with a Catppuccin Mocha-inspired visual theme while keeping user-facing labels and keyboard hints readable in plain terminal output.

#### Scenario: Keyboard hints are visible

- GIVEN the TUI shell is rendered
- WHEN a user reads the initial view
- THEN the view shows keyboard hints for arrows, Vim-style navigation, Command Palette, profile actions, cancellation, and quitting

### Requirement: Command Palette discovers people-facing Actions

The TUI Surface MUST provide a Command Palette discovery mode that lists people-facing Actions across domains based on Capability metadata, and MUST let users move the selected Palette Action with arrow keys and Vim-style `j`/`k` keys.

#### Scenario: Palette includes command-palette-visible people Actions

- GIVEN a Capability has people audience and command-palette visibility
- WHEN the Command Palette is opened
- THEN that Capability appears as a discoverable Action

#### Scenario: Palette filters Actions by query

- GIVEN the Command Palette is open with multiple discoverable Actions
- WHEN a user types a query matching an Action title or Capability ID
- THEN only matching Actions are shown

#### Scenario: Palette selection moves with keyboard

- GIVEN the Command Palette is open with multiple discoverable Actions
- WHEN the user presses `down` or `j`
- THEN the selected Palette Action moves to the next Action
- WHEN the user presses `up` or `k`
- THEN the selected Palette Action moves to the previous Action

#### Scenario: Palette is distinct from primary navigation

- GIVEN a Capability has command-palette visibility but not TUI primary visibility
- WHEN the initial shell is rendered
- THEN that Capability is not shown as a primary Action
- WHEN the Command Palette is opened
- THEN that Capability is shown as a discoverable Action

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

### Requirement: Command Palette Actions collect required string input

The TUI Surface MUST collect string inputs declared by a selected Action's Capability Input Schema before executing that Action, and MUST prevent submission only for required string fields left empty.

#### Scenario: Selected Action opens an input form

- GIVEN a people-facing command-palette-visible Capability declares string input fields
- WHEN the user selects it from the Command Palette
- THEN the TUI shows an input form for those fields instead of executing immediately

#### Scenario: Submitted form executes with collected input

- GIVEN a user has filled all required string fields in the TUI form
- WHEN the user submits the final field
- THEN the TUI executes the selected Capability through the shared Pipeline with those input values

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

### Requirement: Session Profile can be changed temporarily

The TUI Surface MUST let a user change the active Session Profile for the running TUI without changing the saved Default Profile.

#### Scenario: Session Profile form changes active profile

- GIVEN the TUI shell is showing an active Session Profile
- WHEN a user opens the Session Profile form and submits a non-empty profile name
- THEN the TUI shows that profile as the active Profile
- AND the change is limited to the running TUI model

#### Scenario: Session Profile form can be cancelled

- GIVEN the Session Profile form is open
- WHEN a user presses `esc`
- THEN the TUI closes the form
- AND the active Profile remains unchanged

#### Scenario: Subsequent Actions use changed Session Profile

- GIVEN the Session Profile has been changed in the running TUI
- WHEN a user executes a Command Palette Action
- THEN the shared Pipeline receives the changed profile in the execution request

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

### Requirement: Default Profile can be persisted explicitly from the TUI

The TUI Surface MUST provide an explicit user action for saving a new Default Profile through the shared non-sensitive configuration persistence path, without treating temporary Session Profile changes as saved defaults.

#### Scenario: Default Profile form saves profile

- GIVEN the TUI shell is showing an active Session Profile
- WHEN a user opens the Default Profile form and submits a non-empty profile name
- THEN the TUI persists that profile as the saved Default Profile
- AND the running TUI shows that profile as the active Profile
- AND the TUI renders a success message identifying the updated Configuration File

#### Scenario: Default Profile form can be cancelled

- GIVEN the Default Profile form is open
- WHEN a user presses `esc`
- THEN the TUI closes the form
- AND the active Profile remains unchanged
- AND no Default Profile is persisted

#### Scenario: Default Profile save failure is shown

- GIVEN the Default Profile form is open
- WHEN saving the submitted profile fails
- THEN the TUI renders a failure message
- AND the active Profile remains unchanged
