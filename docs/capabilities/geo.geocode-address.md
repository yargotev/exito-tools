# geo.geocode-address

Geocodes a city/address pair using the configured Geo provider and returns the normalized address and selected location fields.

## Contract

- **Capability ID**: `geo.geocode-address`
- **Domain**: Geo
- **CLI Command**: `exito geo geocode-address --city <city> --address <address>`
- **TUI Action**: `Geocode address`
- **Risk Level**: read-only
- **Audience**: agents and people

## Input Schema

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `city` | string | yes | City name accepted by the Geo provider. |
| `address` | string | yes | Address to geocode. |

## Initial Result Shape

```json
{
  "message": "Geocoding successful.",
  "success": true,
  "location": {
    "latitude": "4.598090587",
    "longitude": "-74.160580822"
  },
  "status": "M",
  "normalizedAddress": "CL 57 H SUR # 68 D - 75",
  "neighborhood": "VILLA DEL RIO",
  "daneCode": "110010001"
}
```


## Provider Request

When `EXITO_GEO_BASE_URL` and `EXITO_GEO_TOKEN` are configured, Exito Tools sends a `POST` request to `<EXITO_GEO_BASE_URL>/geocode-address` with a bearer token and JSON body:

```json
{
  "city": "Bogota",
  "address": "CL 57 H SUR # 68 D - 75"
}
```

## External Response Mapping

| External field | Use Case Result field |
| --- | --- |
| `message` | `message` |
| `success` | `success` |
| `data.latitude` | `location.latitude` |
| `data.longitude` | `location.longitude` |
| `data.estado` | `status` |
| `data.dirtrad` | `normalizedAddress` |
| `data.barrio` | `neighborhood` |
| `data.coddane` | `daneCode` |
