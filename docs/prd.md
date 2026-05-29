# PRD: Exito Tools CLI/TUI Architecture

## Problem Statement

Exito operational work needs a single tool that can serve two very different consumers without duplicating domain logic: agents need a stable, scriptable, machine-readable CLI, while everyday users need a guided terminal experience that helps them discover and run useful actions. The project needs clear architecture and roadmap guidance so future domains can be added consistently, Checkout purchase assembly can be built safely, and agents can rely on stable contracts for discovery, execution, output, errors, configuration, and observability.

## Solution

Build Exito Tools as one Go application with two interaction surfaces: a machine-first CLI and a task-first TUI. Operational Domains expose neutral Capabilities backed by shared Use Cases. The CLI adapts Capabilities into Commands and emits JSON by default. The TUI adapts suitable Capabilities into human-friendly Actions using Bubble Tea. Capabilities are registered explicitly during boot into an immutable registry, can be discovered through a machine-readable capabilities command, and can be executed through both explicit domain commands and a generic capability run command.

The first implementation established the core architecture and initial read-only capabilities. The next roadmap expands Exito Tools into a guided VTEX Checkout purchase-assembly tool: create or load an orderForm, add products selected from search results, update client profile data, update shipping data and logistics selections, and inspect the prepared cart. Final order placement and payment processing remain explicitly separate high-risk roadmap steps that require dedicated approval and confirmation gates. Configuration should use non-sensitive YAML plus environment/dotenv credentials. Provider tokens and any customer-sensitive values must remain outside committed files.

## User Stories

1. As an agent, I want CLI commands to return JSON by default, so that I can parse results without extra flags.
2. As an agent, I want a standard JSON envelope, so that every command has predictable success and failure shapes.
3. As an agent, I want structured errors with stable error codes, so that I can react programmatically to failures.
4. As an agent, I want generic exit codes plus specific error codes, so that I can make coarse and detailed decisions separately.
5. As an agent, I want command logs separated from stdout, so that JSON output is never corrupted by debug messages.
6. As an agent, I want every command result to include a request ID, so that I can correlate output with logs and upstream API calls.
7. As an agent, I want to provide an optional correlation ID, so that I can group multiple command executions under one higher-level automation run.
8. As an agent, I want to discover available capabilities through a machine-readable command, so that I can inspect what Exito Tools can do.
9. As an agent, I want capability discovery to include input schemas, output contracts, domains, versions, risk levels, audiences, and surface mappings, so that I can decide how to invoke capabilities safely.
10. As an agent, I want stable capability IDs, so that my automations do not break due to casual renames.
11. As an agent, I want incompatible capability changes to use new versioned IDs, so that breaking changes do not silently alter existing contracts.
12. As an agent, I want a generic capability execution path by ID, so that I can call capabilities without relying only on the explicit command tree.
13. As an agent, I want the generic run command to accept complete input objects, so that I can provide schema-shaped input through JSON, files, or stdin.
14. As an agent, I want explicit domain commands with typed flags, so that common workflows remain convenient and discoverable.
15. As an agent, I want command names to follow a standard verb vocabulary, so that commands are predictable across domains.
16. As an agent, I want risky commands to fail with a structured confirmation-required error instead of prompting, so that scripts never hang on interactive confirmation.
17. As an agent, I want list commands to paginate explicitly, so that commands do not perform unbounded work by default.
18. As an agent, I want warnings in structured metadata, so that I can decide whether partial data, fallbacks, or deprecations matter.
19. As an agent, I want credentials to come from environment variables or dotenv files, so that local and CI execution can inject secrets consistently.
20. As an agent, I want process environment variables to override dotenv files, so that CI and automation can control configuration reliably.
21. As an agent, I want project-local configuration to override user configuration, so that repository context is respected.
22. As an agent, I want `exito` without arguments to show brief help instead of launching an interactive UI, so that accidental invocations do not block automation.
23. As an agent, I want an explicit `exito tui` command for the TUI, so that interactive behavior is opt-in.
24. As an agent, I want Geo credentials to use standard environment variable names, so that I can configure the Geo provider predictably.
25. As an agent, I want the Geo token omitted from committed files, so that secrets are not leaked.
26. As an agent, I want `orders.get`, so that I can fetch an order by identifier through a stable contract.
27. As an agent, I want `geo.geocode-address`, so that I can geocode a city/address pair and receive only the selected normalized fields.
28. As an agent, I want external API DTOs mapped to domain-owned results, so that upstream API shape changes do not leak into my contracts.
29. As an agent, I want errors translated at the layer that understands them, so that error codes preserve useful technical or domain meaning.
30. As an agent, I want request metadata propagated to outbound HTTP where applicable, so that provider calls can be correlated with command results.
31. As a person, I want a task-first TUI, so that I can find what I need to do without understanding the internal domain architecture.
32. As a person, I want curated primary TUI actions, so that the first screen is useful rather than overwhelming.
33. As a person, I want a command palette, so that I can search for available actions across domains.
34. As a person, I want result filters, so that I can refine data already loaded inside a task.
35. As a person, I want agent-only capabilities hidden from primary TUI navigation, so that the UI does not promote technical actions that are not meant for everyday use.
36. As a person, I want the TUI to show the active profile, so that I know which environment and credentials are being used.
37. As a person, I want to change the TUI session profile temporarily, so that I can inspect another environment without changing future defaults.
38. As a person, I want to set a default profile explicitly, so that future CLI and TUI sessions use my preferred environment.
39. As a person, I want long-running TUI actions to show loading, success, failure, and cancellation states, so that I understand what the tool is doing.
40. As a person, I want long-running actions to be cancellable when possible, so that I can stop work I no longer need.
41. As a person, I want risky TUI actions to show impact-aware confirmation prompts, so that I do not accidentally perform destructive operations.
42. As a person, I want user-facing copy in English, so that the interface is consistent for now.
43. As a developer, I want one shared Use Case implementation per capability, so that CLI and TUI do not duplicate business logic.
44. As a developer, I want Operational Domains to avoid importing CLI or TUI frameworks, so that domain logic remains reusable and testable.
45. As a developer, I want a neutral Capability model, so that surfaces can adapt the same contract into Commands and Actions.
46. As a developer, I want schema-first inputs, so that CLI flags, TUI forms, and agent discovery come from one source of truth.
47. As a developer, I want structured, format-neutral Use Case Results, so that JSON, future output formats, and TUI views do not shape domain logic.
48. As a developer, I want presenters to own output formatting, so that domains do not return pretty strings or UI models.
49. As a developer, I want explicit domain wiring, so that registered capabilities are visible and testable.
50. As a developer, I want the Capability Registry to be immutable after boot, so that discovery is predictable at runtime.
51. As a developer, I want internal modular domains rather than dynamic plugins initially, so that the first implementation avoids unnecessary complexity.
52. As a developer, I want Workflow Capabilities for cross-domain orchestration, so that future cross-domain use cases do not get forced into unrelated domains.
53. As a developer, I want workflow capability IDs to use business names, so that visible contracts do not expose technical implementation terms.
54. As a developer, I want shared HTTP infrastructure plus domain API clients, so that cross-cutting HTTP concerns are reused while domain semantics stay near the domain.
55. As a developer, I want Viper only behind an explicit configuration resolver, so that helper libraries do not define product precedence rules.
56. As a developer, I want slog-based structured logging, so that logs are structured without adding unnecessary dependencies.
57. As a developer, I want a committed `.env.example`, so that setup is discoverable without committing secrets.
58. As a developer, I want local dotenv files ignored, so that profile-specific credentials are protected by default.
59. As a maintainer, I want decisions captured as ADRs and glossary terms, so that future agents understand why the architecture exists.
60. As a maintainer, I want a PRD that synthesizes the architecture, so that implementation can start in future sessions without repeating the interview.

## Implementation Decisions

- Build one Go application named Exito Tools with two interaction surfaces: CLI and TUI.
- Use Cobra for CLI command routing.
- Use Bubble Tea for the TUI.
- Use Viper only as an implementation helper behind an explicit Configuration Resolver.
- Use Go `log/slog` for structured logging.
- Keep visible contracts, help text, labels, messages, commands, and capability IDs English-only for now.
- Treat Orders and Geo as Operational Domains, not command categories.
- Model Capabilities as neutral contracts backed by Use Cases.
- Keep Use Cases shared across CLI and TUI.
- Keep Operational Domains independent from Cobra, Bubble Tea, and surface-specific packages.
- Use explicit application wiring to register domains and capabilities.
- Treat the Capability Registry as immutable after boot.
- Support internal modular domains now; defer dynamic external plugins.
- Allow Workflow Capabilities for cross-domain orchestration, but use business-oriented visible IDs rather than a technical workflow prefix.
- Use stable Capability IDs in the form `<domain>.<action>`, with lower-case kebab-case for multi-word actions.
- Use capability version metadata for compatible changes and new versioned IDs for incompatible changes.
- Use schema-first inputs for capabilities.
- Use structured, format-neutral Use Case Results.
- Use Presenters for JSON and future output formats.
- Emit JSON by default for CLI command output.
- Use a standard JSON Envelope with `ok`, `data` or `error`, and `meta`.
- Include `requestId`, optional `correlationId`, `profile`, `capabilityId`, and `durationMs` in standard envelope metadata.
- Keep logs and diagnostics out of stdout JSON; use stderr or log files.
- Use generic exit codes and specific structured error codes.
- Represent non-fatal issues as structured warnings in envelope metadata.
- Use explicit cursor-based pagination for list-style capabilities; do not auto-paginate by default.
- Use risk and confirmation metadata for capabilities.
- Keep CLI confirmation non-interactive; missing confirmation returns `CONFIRMATION_REQUIRED`.
- Make Capability Execution context-aware and cancellable where possible.
- Generate a unique request ID per capability execution.
- Accept an optional correlation ID for grouping multiple executions.
- Use `exito tui` as the explicit TUI entrypoint.
- Make `exito` without arguments show brief textual help rather than opening TUI.
- Expose `exito capabilities` as the machine-readable capability inventory.
- Support a conceptual generic run path with `exito run <capability-id>`.
- Make generic run accept complete input objects through inline JSON, input files, or piped stdin.
- Use explicit domain commands in the form `exito <domain> <action> [flags]`.
- Use standard action verbs: `get`, `list`, `search`, `validate`, `diagnose`, and mutating verbs only when state changes.
- Use task-first TUI navigation with curated primary Actions.
- Use a Command Palette for broader action discovery.
- Keep Result Filters distinct from Command Palette discovery.
- Use capability visibility and audience metadata to control CLI, TUI, and palette exposure.
- Support agent-only capabilities without promoting them in human-facing TUI flows by default.
- Use shared Application Configuration across CLI and TUI.
- Use Profiles to represent environment plus credentials.
- Use `staging` as the initial default Profile.
- Resolve Effective Profile using explicit profile flag, then `EXITO_PROFILE`, then saved default, then `staging` fallback.
- Keep TUI Session Profile changes separate from setting the Default Profile.
- Store non-sensitive configuration as YAML.
- Discover configuration by explicit config path, then `EXITO_CONFIG`, then local project config, then user config, then internal defaults.
- Use environment variables and non-committed dotenv files for credentials; do not use OS keychains.
- Layer dotenv values by real process environment, then profile-specific dotenv, then general dotenv, then non-sensitive defaults.
- Use a committed dotenv example template with no real secrets.
- Ignore local dotenv files while allowing the dotenv example template.
- Use shared low-level HTTP infrastructure for base URLs, auth headers, timeouts, retries, and request metadata.
- Keep external API semantics and DTOs in domain-owned API clients.
- Map external DTOs to domain-owned results before returning from Use Cases.
- Translate errors at the layer that understands their meaning.
- Define `orders.get` as the first Orders Capability.
- Define `geo.geocode-address` as the first Geo Capability.
- Define Checkout as its own Operational Domain for VTEX orderForm and cart assembly behavior.
- Add a checkout roadmap that starts with `checkout.get-order-form`, `checkout.create-order-form`, `checkout.add-items`, `checkout.update-client-profile`, and `checkout.update-shipping-data` before any final place-order or payment capabilities.
- Keep product discovery in Catalog; Checkout may consume selected SKU IDs from search results but must not hide Catalog search side effects inside cart mutations.
- Treat VTEX Checkout write operations as confirmation-required safe-write capabilities, with no parallel orderForm mutations in a single flow.
- Keep final order placement and payment processing out of the first Checkout slice until risk, credentials, and non-production validation are explicitly approved.
- Map the Geo provider response to only message, success, latitude, longitude, status, normalized address, neighborhood, and DANE code.
- Configure Geo with `EXITO_GEO_BASE_URL` and `EXITO_GEO_TOKEN`.
- Keep the Geo endpoint path owned by the Geo Domain API Client.
- Never commit or document real provider tokens.

### Major Modules To Build

- Configuration Resolver: deep module that computes effective config from flags, environment variables, dotenv layers, YAML config, profiles, and defaults.
- Capability Core: deep module defining Capability metadata, Input Schema, Use Case Result, Structured Error, warnings, risk, audience, visibility, and execution contracts.
- Capability Registry: deep module for explicit boot-time registration and immutable runtime discovery.
- Execution Pipeline: deep module that validates inputs, applies profile/config context, creates request metadata, executes Use Cases, measures duration, and returns structured results/errors.
- Presenter Layer: deep module for JSON Envelope output and future output formats.
- CLI Surface: Cobra adapter that builds root/domain/run/capabilities/tui commands from registry and execution primitives.
- TUI Surface: Bubble Tea adapter for task-first navigation, Command Palette, result filters, profile display, and task execution states.
- HTTP Infrastructure: deep module for authenticated provider requests, timeouts, retries, request IDs, and technical error translation.
- Orders Domain: domain module exposing order lookup capabilities such as `orders.get` and `orders.get-vtex`.
- Geo Domain: domain module exposing `geo.geocode-address`, VTEX region diagnostics, and mapping provider DTOs to domain results.
- Catalog Domain: domain module exposing product discovery capabilities used before Checkout cart mutation.
- Checkout Domain: domain module owning VTEX orderForm creation/loading, cart item updates, client profile attachments, shipping/logistics attachments, and later payment/place-order steps.
- Workflow Layer: module for business-named cross-domain capabilities such as guided purchase assembly from Catalog search into Checkout orderForm updates.

## Testing Decisions

- Tests should verify external behavior and stable contracts, not internal implementation details.
- Configuration Resolver should have table-driven tests for profile resolution, config file discovery, dotenv layering, and environment override precedence.
- Capability Core should test schema validation, structured error representation, risk metadata, audience/visibility metadata, and versioned ID behavior.
- Capability Registry should test duplicate ID rejection, successful registration, immutable finalized registry behavior, and discovery output ordering if ordering becomes part of the contract.
- Execution Pipeline should test request ID generation, correlation ID propagation, duration metadata, context cancellation, structured success, structured failure, and warning propagation.
- Presenter Layer should test exact JSON Envelope shapes for success, error, warnings, pagination, and metadata.
- CLI Surface should test root help behavior, JSON default output, stdout/stderr separation, exit code mapping, `exito capabilities`, `exito run`, and explicit command mappings.
- CLI confirmation behavior should test that risky commands fail with `CONFIRMATION_REQUIRED` instead of prompting.
- TUI Surface should be tested around model update behavior, navigation state, command palette filtering, result filtering, profile display, loading/success/error/cancelled task states, and confirmation flows.
- HTTP Infrastructure should use fake servers to test base URL handling, auth header injection, request ID propagation, timeouts, retries if implemented, and technical error translation.
- Geo Domain should test provider DTO mapping to the selected domain result fields and semantic provider error translation.
- Orders Domain should test the `orders.get` contract through a fake client.
- End-to-end CLI tests should run the compiled command or command root with fake dependencies and assert stdout JSON, stderr logging behavior, and exit codes.
- Bubble Tea tests should prefer model-level tests and command/message simulation rather than terminal snapshot fragility.
- No tests should require real provider tokens or real network calls.

## Out of Scope

- Dynamic external plugins, runtime plugin loading, version negotiation, or marketplace behavior.
- Full internationalization; English-only is the current product language.
- OS keychain, Windows Credential Manager, macOS Keychain, or Linux Secret Service integration.
- A graphical web UI.
- Implicit TUI launch from `exito` without arguments.
- Storing secrets in YAML configuration or committed documentation.
- Auto-pagination by default for list commands.
- Domain-specific exit codes for every possible failure.
- Letting external API DTOs define CLI/TUI contracts directly.
- Building every future domain in the first slice; Checkout is now on the roadmap but must be delivered incrementally.
- Implementing destructive domain operations unless a future capability requires them.

## Further Notes

The repo now contains the Go scaffold and several implemented capabilities. The next product roadmap is VTEX Checkout purchase assembly. Start with non-production, confirmation-gated orderForm operations and documentation/tests for each public contract before adding broader TUI workflows. The existing ADRs, glossary, OpenSpec specs, and capability docs are the source of truth for terminology and decisions.
