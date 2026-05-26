# Proposal: YAML Profile Provider Base URLs

## Summary

Allow the selected non-sensitive YAML Configuration File to provide profile-scoped provider base URLs for the initial Orders and Geo providers.

## Motivation

Profiles represent environment plus credentials. The resolver already selects an Effective Profile and loads provider values from environment/dotenv layers, but YAML only stores the saved Default Profile. Adding profile-scoped base URLs lets users keep non-sensitive environment endpoints in shared configuration while preserving tokens in environment or non-committed dotenv files.

## Scope

- Parse simple profile-scoped YAML base URLs from the selected Configuration File.
- Use YAML base URLs only as a non-sensitive fallback after process environment, profile-specific dotenv, and general dotenv values.
- Keep provider tokens resolved only from environment/dotenv layers and omitted from JSON serialization.
- Document the supported YAML shape.

## Non-goals

- Full YAML support or arbitrary nested configuration parsing.
- Storing provider tokens or secrets in YAML.
- Changing profile precedence or Configuration File discovery precedence.
- Rebooting TUI provider clients after a Session Profile change.
