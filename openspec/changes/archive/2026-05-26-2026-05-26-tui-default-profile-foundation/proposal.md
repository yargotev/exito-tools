# Proposal: TUI Default Profile foundation

## Summary

Add a narrow TUI flow for explicitly saving a new Default Profile while keeping temporary Session Profile changes separate.

## Motivation

The TUI already lets users change the active Session Profile temporarily, and the CLI can persist the saved Default Profile to YAML. PRD story 38 requires a person-facing way to set the Default Profile explicitly so future CLI and TUI sessions use the preferred environment. The next smallest slice is to expose that persistence path in the TUI without changing domain execution or configuration precedence.

## Scope

- Add a TUI form for setting the saved Default Profile explicitly.
- Persist the submitted profile through the existing Configuration Resolver persistence path.
- Update the running TUI Session Profile to the saved value after a successful save.
- Render a success or failure message without writing secrets.
- Cover save, cancel, and failure behavior with focused TUI model tests.

## Out of Scope

- Listing or validating known profile names.
- Editing full YAML configuration or provider endpoints.
- Persisting credentials or secrets.
- Rebootstrapping provider clients after changing the running profile.
