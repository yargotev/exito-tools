# Configuration Resolver Delta: YAML Default Profile Foundation

## ADDED Requirements

### Requirement: YAML configuration can provide saved Default Profile

The Configuration Resolver MUST load a saved Default Profile from the selected non-sensitive YAML Configuration File when no saved Default Profile is supplied directly to the resolver.

#### Scenario: Local YAML default profile is used

- GIVEN no explicit profile and no `EXITO_PROFILE` are set
- AND local `./exito.yaml` contains `defaultProfile: prod`
- WHEN configuration is resolved
- THEN the Effective Profile is `prod`
- AND the Effective Profile source is saved default

#### Scenario: Explicit profile overrides YAML default profile

- GIVEN explicit profile input is `qa`
- AND local `./exito.yaml` contains `defaultProfile: prod`
- WHEN configuration is resolved
- THEN the Effective Profile is `qa`
- AND the Effective Profile source is explicit

#### Scenario: Environment profile overrides YAML default profile

- GIVEN no explicit profile is set
- AND `EXITO_PROFILE` is `dev`
- AND local `./exito.yaml` contains `defaultProfile: prod`
- WHEN configuration is resolved
- THEN the Effective Profile is `dev`
- AND the Effective Profile source is environment
