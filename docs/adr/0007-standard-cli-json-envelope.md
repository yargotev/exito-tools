# Standard CLI JSON envelope

The Exito Tools CLI Surface will wrap JSON output in a standard envelope: successful commands return `ok: true` with `data` and optional `meta`, while failed commands return `ok: false` with a structured `error` and optional `meta`. This gives agents a consistent contract across domains and leaves room for metadata such as warnings, pagination, timestamps, and request identifiers.
