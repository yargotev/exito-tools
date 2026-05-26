# Proposal: TUI session profile foundation

## Summary

Add the first temporary TUI Session Profile switcher so people can change the active profile for the running TUI without changing the saved Default Profile.

## Motivation

The TUI already shows the Effective Profile and executes Actions through the shared Pipeline with that profile. ADR 0014 distinguishes temporary Session Profile changes from explicit Default Profile changes, but the TUI has no way to update the session profile yet.

## Scope

- Add a minimal Session Profile mode opened from the TUI shell.
- Collect a non-empty profile name and apply it to the running TUI model only.
- Ensure subsequent TUI Action executions use the updated Session Profile.
- Keep the shared configuration resolver and saved Default Profile unchanged.

## Non-goals

- Persisting a new Default Profile.
- Listing configured profiles from YAML.
- Validating profile names against a configuration file.
- Changing the profile while another input mode or task execution is active.
