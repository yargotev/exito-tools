# Proposal: TUI confirmation prompt foundation

## Summary

Add the narrow foundation for impact-aware TUI confirmation prompts before running confirmation-required Capabilities. The slice keeps confirmation local to the TUI Surface, reuses Capability risk/confirmation metadata, and passes explicit confirmation into the shared execution Pipeline only after the user confirms.

## Problem

The shared Pipeline already rejects confirmation-required Capabilities unless a surface passes explicit confirmation. The CLI has a non-interactive `--confirm` flow, but the TUI currently sends selected Actions directly to the Pipeline. A people-facing risky Action would therefore fail with `CONFIRMATION_REQUIRED` instead of presenting the interactive confirmation UX required by the PRD.

## Scope

- Detect command-palette Actions whose Capability definition requires confirmation.
- Render an English confirmation prompt with action title, Capability ID, risk level, and description when available.
- Let `y`/`enter` confirm and execute through the shared Pipeline with explicit confirmation.
- Let `n`/`esc` cancel the prompt without executing the Capability.
- Cover prompt, confirm, and cancel behavior with focused TUI model tests.

## Out of Scope

- Destructive target-identifier re-entry or domain-specific impact summaries.
- CLI confirmation behavior changes.
- Persisting confirmation decisions across Actions or sessions.
- Adding new mutating production Capabilities.

## Rollback

Revert this change to restore the previous TUI behavior where confirmation-required Actions reach the Pipeline without surface-level confirmation and fail with `CONFIRMATION_REQUIRED`.
