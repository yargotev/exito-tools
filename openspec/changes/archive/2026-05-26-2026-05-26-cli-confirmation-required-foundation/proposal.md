# Proposal: CLI confirmation required foundation

## Summary

Add the narrow shared execution and CLI foundation for risky Capability confirmation so non-interactive CLI commands fail with a structured `CONFIRMATION_REQUIRED` error when confirmation-required Capabilities are invoked without explicit confirmation.

## Motivation

ADR 0034 and ADR 0035 require risk and confirmation metadata to prevent silent risky execution. Capability definitions already expose `requiresConfirmation`, but the execution path does not enforce it and the generic `run` command has no explicit confirmation input.

## Scope

- Add a confirmation flag to the generic `exito run <capability-id>` command.
- Pass confirmation intent through the shared execution Pipeline request.
- Make the Pipeline reject confirmation-required Capabilities with `CONFIRMATION_REQUIRED` when confirmation is missing.
- Preserve existing read-only behavior and successful execution when confirmation is provided.

## Non-goals

- Add destructive domain capabilities.
- Add TUI impact-aware confirmation prompts.
- Add target-specific destructive confirmation text.
- Add confirmation flags to domain commands that do not currently require confirmation.
