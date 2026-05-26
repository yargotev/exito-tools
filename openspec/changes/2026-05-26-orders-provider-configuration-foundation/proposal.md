# Proposal: Orders provider configuration foundation

## Problem

The `orders.get` capability exists but still uses an unavailable dependency because Orders provider connection settings are not resolved by application configuration. Before adding an HTTP Orders getter, the application needs a safe, profile-aware way to resolve Orders base URL and token values.

## Scope

- Add Orders provider base URL and token environment variable names.
- Resolve Orders provider settings through the existing credential layer order: process environment, `.env.<profile>`, then `.env`.
- Expose non-sensitive Orders provider metadata in `config.Effective` while omitting token values from JSON serialization.
- Document Orders environment variables and update the committed `.env.example` template.

## Out of Scope

- Implementing the Orders HTTP getter.
- Wiring `orders.get` to a real provider client.
- Parsing non-sensitive YAML configuration.
