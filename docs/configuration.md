# Configuration

## YAML configuration

Exito Tools reads non-sensitive application settings from the selected YAML Configuration File. The supported foundation key is:

```yaml
defaultProfile: staging
profiles:
  staging:
    geo:
      baseUrl: https://sitidataws.sitimapa.co
    orders:
      baseUrl: https://orders.example.test
```

`defaultProfile` stores the saved Default Profile used when neither `--profile` nor `EXITO_PROFILE` is set. Secrets and provider tokens must stay in environment variables or non-committed dotenv files, not YAML.

`profiles.<profile>.<provider>.baseUrl` stores non-sensitive provider endpoints for the Effective Profile. Environment and dotenv values still override YAML base URLs; tokens are never read from YAML.

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
4. Selected YAML `profiles.<profile>.geo.baseUrl` for non-sensitive base URL only.

The resolver marks Geo as configured only when a Geo base URL and `EXITO_GEO_TOKEN` are present. The token value is kept out of JSON serialization; only token presence/source metadata may be exposed.

### Orders

```env
# Non-sensitive endpoint can also live in YAML profiles.<profile>.orders.baseUrl.
EXITO_ORDERS_BASE_URL=https://bromoqa.grupo-exito.com/apioms/api/v1/TSYHCSNIDTKZJWM/geoms

# Preferred GEOMS client-credentials variables for non-committed env/dotenv files.
EXITO_ORDERS_CLIENT_ID=...
EXITO_ORDERS_CLIENT_SECRET=...
EXITO_ORDERS_SCOPE=...

# Optional override; defaults to the GEOMS Azure AD tenant token URL.
EXITO_ORDERS_TOKEN_URL=https://login.microsoftonline.com/40f94963-1b34-45ce-a5fb-6f1fde2f1a27/oauth2/v2.0/token
```

`EXITO_ORDERS_CLIENT_SECRET` and any GEOMS credential bundle are sensitive and must only live in the real process environment or non-committed dotenv files. The resolver also accepts the legacy GEOMS bundle variables `GEOMS_CREDENTIALS_QA` for non-prod profiles and `GEOMS_CREDENTIALS_PDN` for `prod`/`production`/`pdn` profiles, extracting `client_id`, `client_secret`, and `scope`.

Orders provider values are resolved in this order:

1. Real process environment variables.
2. Profile-specific `.env.<profile>` file for the Effective Profile.
3. General `.env` file.
4. Selected YAML `profiles.<profile>.orders.baseUrl` for non-sensitive base URL only.

The resolver marks Orders as configured when an Orders base URL is present and either a pre-fetched `EXITO_ORDERS_TOKEN` or complete GEOMS client credentials are present. Secret values are kept out of JSON serialization; only presence/source metadata may be exposed.
