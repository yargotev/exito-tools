# Immutable Capability Registry after boot

Exito Tools builds the Capability Registry during application boot through explicit wiring and treats it as immutable afterward. Interaction surfaces read from the finalized registry, which keeps capability discovery predictable for agents and avoids runtime surprises.
