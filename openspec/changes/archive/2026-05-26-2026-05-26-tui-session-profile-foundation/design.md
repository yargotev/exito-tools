# Design: TUI session profile foundation

## Approach

Extend the Bubble Tea model with a small Session Profile state that is entered from the idle shell. The mode stores typed text only; on submit, the trimmed value replaces `Model.profile`, which is already the value used by TUI capability execution requests.

## Decisions

- `p` opens the Session Profile form from the idle shell.
- The form renders the current active profile and an empty new-profile input.
- Empty submissions stay in the form and do not change the Session Profile.
- `esc` closes the form without changing the Session Profile.
- The change is model-local and does not write YAML, dotenv, environment variables, or saved defaults.

## Risks

- Profile names are not validated against persisted config because profile listing/parsing is not implemented yet.
- Existing provider clients are created at application boot, so this slice updates execution metadata/profile context but does not rebuild domain provider configuration for the new profile.
