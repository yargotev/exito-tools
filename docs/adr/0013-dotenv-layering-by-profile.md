# Dotenv layering by Profile

Exito Tools will support layered dotenv loading: real process environment variables win, then `.env.{profile}` for the Effective Profile, then the general `.env`, then non-sensitive configuration defaults. This lets local development keep shared defaults while allowing profile-specific credentials and endpoints without committing secrets.
