# Configuration Resolver Delta

## ADDED Requirements

### Requirement: Saved Default Profile can be persisted to YAML

The Configuration Resolver MUST provide a narrow persistence path for updating the saved Default Profile in the non-sensitive YAML Configuration File without writing credentials.

#### Scenario: Existing YAML default profile is updated

- **GIVEN** the selected Configuration File contains `defaultProfile: staging`
- **WHEN** the saved Default Profile is persisted as `prod`
- **THEN** the same Configuration File contains `defaultProfile: prod`
- **AND** no credential keys are written

#### Scenario: Missing YAML default profile is appended

- **GIVEN** the selected Configuration File exists without `defaultProfile`
- **WHEN** the saved Default Profile is persisted as `qa`
- **THEN** the file contains a top-level `defaultProfile: qa` entry

#### Scenario: No configuration file creates local project config

- **GIVEN** no explicit path, no environment path, no local config, and no user config exist
- **WHEN** the saved Default Profile is persisted as `dev`
- **THEN** local `./exito.yaml` is created with `defaultProfile: dev`
