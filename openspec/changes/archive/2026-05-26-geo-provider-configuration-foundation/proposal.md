# Proposal: Geo provider configuration foundation

## Summary

Resolve the Geo provider base URL and token into shared Application Configuration so later Geo HTTP client slices can use a single deterministic configuration seam.

## Scope

- Resolve `EXITO_GEO_BASE_URL` and `EXITO_GEO_TOKEN` from process environment, `.env.<profile>`, then `.env`.
- Expose Geo provider readiness (`configured`) and source metadata in `config.Effective`.
- Keep the actual token omitted from JSON serialization while exposing only token presence.
- Add resolver tests for precedence, dotenv fallback, missing-token behavior, and secret redaction.

## Out of Scope

- Real Geo provider HTTP client.
- Wiring a real geocoder into `geo.geocode-address`.
- YAML configuration parsing.
- Logging or CLI commands that display configuration values.
