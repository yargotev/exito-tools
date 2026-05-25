# Viper behind explicit configuration resolver

Exito Tools will use Viper only behind an explicit configuration resolver, not as the architectural source of truth for precedence. Viper is useful with Cobra for reading config files, flags, and environment variables, but Exito Tools must still enforce its own documented resolution rules for Profiles and dotenv layering.
