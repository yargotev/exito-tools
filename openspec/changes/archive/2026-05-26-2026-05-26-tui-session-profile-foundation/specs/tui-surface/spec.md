# TUI Surface Delta: Session Profile foundation

## ADDED Requirements

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
