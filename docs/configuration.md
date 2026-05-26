# Configuration

## YAML configuration

Exito Tools reads non-sensitive application settings from the selected YAML Configuration File. The supported foundation key is:

```yaml
defaultProfile: staging
```

`defaultProfile` stores the saved Default Profile used when neither `--profile` nor `EXITO_PROFILE` is set. Secrets and provider tokens must stay in environment variables or non-committed dotenv files, not YAML.

Use `exito config set-default-profile <profile>` to persist a new saved Default Profile explicitly. The command updates the selected YAML Configuration File, or creates local `./exito.yaml` when no configuration file exists.

## Environment variables

### Geo

```env
EXITO_GEO_BASE_URL=https://sitidataws.sitimapa.co
EXITO_GEO_TOKEN=...
```

`EXITO_GEO_TOKEN` is sensitive and must only live in the real process environment or non-committed dotenv files.

Geo provider values are resolved in this order:

1. Real process environment variables.
2. Profile-specific `.env.<profile>` file for the Effective Profile.
3. General `.env` file.

The resolver marks Geo as configured only when both `EXITO_GEO_BASE_URL` and `EXITO_GEO_TOKEN` are present. The token value is kept out of JSON serialization; only token presence/source metadata may be exposed.

### Orders

```env
EXITO_ORDERS_BASE_URL=https://orders.example.test
EXITO_ORDERS_TOKEN=...
```

`EXITO_ORDERS_TOKEN` is sensitive and must only live in the real process environment or non-committed dotenv files.

Orders provider values are resolved in this order:

1. Real process environment variables.
2. Profile-specific `.env.<profile>` file for the Effective Profile.
3. General `.env` file.

The resolver marks Orders as configured only when both `EXITO_ORDERS_BASE_URL` and `EXITO_ORDERS_TOKEN` are present. The token value is kept out of JSON serialization; only token presence/source metadata may be exposed.
