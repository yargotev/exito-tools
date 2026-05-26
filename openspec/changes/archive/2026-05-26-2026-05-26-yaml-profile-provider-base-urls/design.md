# Design: YAML Profile Provider Base URLs

## Approach

Extend the Configuration Resolver with a narrow dependency-free parser for the supported YAML shape:

```yaml
defaultProfile: staging
profiles:
  staging:
    geo:
      baseUrl: https://geo.example.test
    orders:
      baseUrl: https://orders.example.test
```

The parser remains inside `internal/config` and extracts only the known non-sensitive fields. Provider token values continue to come from process environment, `.env.<profile>`, and `.env` only.

## Precedence

For provider base URLs, the resolver keeps existing environment/dotenv precedence and adds YAML as the last non-sensitive source:

1. Process environment variable.
2. Profile-specific dotenv file.
3. General dotenv file.
4. Selected YAML `profiles.<effective-profile>.<provider>.baseUrl`.

Provider tokens keep the existing credential-only precedence and are never read from YAML.

## Contracts

`Effective.GeoProvider.BaseURLSource` and `Effective.OrdersProvider.BaseURLSource` report `config-file` when the URL came from YAML. Token serialization remains unchanged (`json:"-"`).
