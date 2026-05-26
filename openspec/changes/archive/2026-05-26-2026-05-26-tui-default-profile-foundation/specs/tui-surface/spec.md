# TUI Surface Delta Specification

## ADDED Requirements

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
