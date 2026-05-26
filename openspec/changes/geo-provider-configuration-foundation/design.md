# Design: Geo provider configuration foundation

## Approach

Extend the existing `internal/config` resolver rather than adding domain-specific environment reads in Geo code:

- Add `config.GeoProvider` to `config.Effective`.
- Reuse the existing credential layer order: process environment, `.env.<profile>`, then `.env`.
- Parse only simple dotenv `KEY=value` lines needed by the resolver, ignoring missing files and comments.
- Treat blank values as absent so a lower layer can still supply a usable value.
- Mark the provider configured only when both `EXITO_GEO_BASE_URL` and `EXITO_GEO_TOKEN` are present.
- Keep `GeoProvider.Token` available in memory for future client wiring, but tag it `json:"-"`; JSON can expose `tokenSet` and source metadata, never the token value.

This slice intentionally leaves application wiring on `geo.UnavailableGeocoder{}`. The next slice can introduce the Geo HTTP client and decide whether to wire it when `Config.GeoProvider.Configured` is true.
