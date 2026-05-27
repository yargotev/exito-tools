## MODIFIED Requirements

### Requirement: Orders get command

The CLI Surface MUST expose `exito orders get --id <order-id>` as an explicit domain command for the `orders.get` Capability and MAY accept an optional GEOMS order-type filter.

#### Scenario: Orders get supports Carulla order type

- **GIVEN** the Application has registered `orders.get`
- **WHEN** a user runs `exito orders get --id A123 --order-type CarullaEcomm`
- **THEN** the command executes the `orders.get` Capability through the shared execution Pipeline
- **AND** the Capability input contains `id` equal to `A123`
- **AND** the Capability input contains `orderType` equal to `CarullaEcomm`
