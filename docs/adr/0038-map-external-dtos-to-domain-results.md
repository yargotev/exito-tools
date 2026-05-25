# Map external DTOs to domain results

Domain API clients in Exito Tools may use external DTOs internally, but Use Cases return domain-owned models or Use Case Results rather than leaking external API response shapes. This keeps the domain language stable, protects interaction surfaces from API changes, and lets presenters format results consistently.
