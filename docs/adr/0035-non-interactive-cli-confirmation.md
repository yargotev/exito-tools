# Non-interactive CLI confirmation

The Exito Tools CLI Surface will not prompt interactively for risky operation confirmation. If confirmation is required and missing, the command fails with a structured `CONFIRMATION_REQUIRED` error; safe writes may use `--confirm`, while destructive operations should require a stronger explicit confirmation such as confirming the target resource identifier.
