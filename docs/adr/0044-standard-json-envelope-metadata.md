# Standard JSON envelope metadata

The Exito Tools JSON Envelope includes standard metadata for every command result: `requestId`, optional `correlationId`, `profile`, `capabilityId`, and `durationMs`. Additional metadata such as warnings, pagination, or deprecation notices may appear when relevant, but secrets and sensitive headers must never be included.
