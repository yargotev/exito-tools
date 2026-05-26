# Design: TUI Default Profile foundation

## Approach

Extend the TUI model with a small `defaultProfileForm` state that is separate from the existing temporary Session Profile form. Pressing `d` opens the Default Profile form. Submitting a non-empty profile calls the existing `config.SetDefaultProfile` persistence helper with the boot options captured by application wiring.

The save operation is modeled as a Bubble Tea command returning a message so the TUI can render success or failure after the command completes. On success, the running Session Profile is updated to the saved profile and a short status message identifies the configuration file that was updated. On failure, the form closes and the TUI renders the error message while leaving the current Session Profile unchanged.

## Decisions

- `p` remains a temporary Session Profile change; `d` is the explicit persistent Default Profile flow.
- The application wiring stores normalized config boot options on `app.Application` so surfaces can reuse the same selected config path/workdir/home/env behavior.
- The TUI does not validate profile existence yet; it relies on the same non-blank/single-line validation as `config.SetDefaultProfile`.
- The TUI updates its current Session Profile after a successful save because the user explicitly chose that profile as the new default.

## Risks

- Provider clients are not rebuilt after changing the running profile; this matches the existing temporary Session Profile foundation and can be addressed in a future rebootstrap slice.
- The save command performs file I/O; errors are surfaced in the TUI as a non-fatal status instead of terminating the program.
