## ADDED Requirements

### Requirement: Geo provider HTTP mapping

The Geo Domain SHALL map configured provider geocode responses to the stable `geo.geocode-address` result contract without exposing provider DTOs to callers.

#### Scenario: Provider response is mapped to domain result

- **GIVEN** `EXITO_GEO_BASE_URL` and `EXITO_GEO_TOKEN` are configured
- **WHEN** `geo.geocode-address` receives a successful provider response containing `data.latitude`, `data.longitude`, `data.estado`, `data.dirtrad`, `data.barrio`, and `data.coddane`
- **THEN** the capability result contains `location.latitude`, `location.longitude`, `status`, `normalizedAddress`, `neighborhood`, and `daneCode`
