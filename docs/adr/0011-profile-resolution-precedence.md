# Profile resolution precedence

Exito Tools resolves the effective Profile using the closest explicit choice first: `--profile` flag, then `EXITO_PROFILE`, then the saved Default Profile, then the initial fallback `staging`. This gives automation a predictable override path while keeping everyday use simple.
