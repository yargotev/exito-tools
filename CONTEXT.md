# Exito Tools

A single application for interacting with Exito operational domains through interfaces suited to both automation and everyday users.

## Language

**Application**:
The single product that exposes Exito operational capabilities across multiple interaction surfaces. It is not two separate products.
_Avoid_: Separate CLI app, separate TUI app

**User-facing Language**:
The language used for human-readable labels, help, and messages in Exito Tools. The current User-facing Language is English only.
_Avoid_: Mixed Spanish/English UI copy before i18n exists

**Application Configuration**:
Shared YAML configuration used by all **Interaction Surfaces** for non-sensitive concerns such as profiles, base endpoints, and domain-specific overrides. Secrets belong in environment variables or non-committed dotenv files.
_Avoid_: Separate CLI config, separate TUI config, committed secrets

**YAML Configuration**:
The non-sensitive configuration file format for Exito Tools. YAML stores profiles, endpoints, and domain overrides, but not credentials.
_Avoid_: JSON config for human-edited settings, secrets in YAML

**Configuration File**:
The YAML file that stores non-sensitive **Application Configuration**. Discovery precedence is explicit `--config`, then `EXITO_CONFIG`, then local `./exito.yaml`, then `~/.config/exito-tools/config.yaml`, then defaults.
_Avoid_: Hidden config precedence

**Local Project Configuration**:
A repository-local `./exito.yaml` file. When present, it takes precedence over user-level configuration because it represents the current project context.
_Avoid_: Ignoring repo context

**Configuration Resolver**:
The component that computes effective application configuration by applying Exito Tools precedence rules for flags, environment variables, dotenv files, profiles, and defaults. Libraries may help read sources, but the resolver owns the rules.
_Avoid_: Library-defined precedence as product behavior

**Credentials Source**:
The place where Exito Tools obtains sensitive authentication values. Credentials come from environment variables or non-committed `.env` files, not OS-specific keychains.
_Avoid_: OS keychain, committed credentials

**Dotenv File**:
A non-committed file containing environment variables for local execution of Exito Tools. Exito Tools supports a general `.env` and profile-specific `.env.{profile}` files, with real process environment variables taking precedence.
_Avoid_: Committed config file with secrets

**Dotenv Layering**:
The order for resolving environment-backed configuration: real process environment variables, then `.env.{profile}` for the **Effective Profile**, then general `.env`, then non-sensitive configuration defaults.
_Avoid_: Unspecified secret precedence

**Profile**:
A named configuration context made of an environment and credentials. When no Profile is specified, Exito Tools uses the configured default Profile; the initial default is `staging`.
_Avoid_: Store, region, domain filter

**Default Profile**:
The **Profile** used when a Command or TUI flow does not explicitly select one. The Default Profile can be changed by the user.
_Avoid_: Hard-coded environment

**Effective Profile**:
The Profile selected for a specific command or TUI flow after applying precedence: explicit `--profile`, then `EXITO_PROFILE`, then saved **Default Profile**, then `staging` fallback.
_Avoid_: Implicit environment when a Profile has been resolved

**Session Profile**:
The Profile currently active within a running **TUI Surface** session. Changing the Session Profile does not change the **Default Profile** unless the user explicitly chooses to set it as default.
_Avoid_: Silent default profile change

**Interaction Surface**:
A way of using the **Application** tailored to a type of user and workflow. The known surfaces are **CLI Surface** and **TUI Surface**.
_Avoid_: Separate app, frontend

**CLI Surface**:
The machine-first interaction surface intended for agents and automation. It emits stable output suitable for programmatic consumption, with JSON as the default **Output Format**.
_Avoid_: Bot app, agent app, human-first terminal UI

**Cobra**:
The Go framework used to route and execute **Commands** in the **CLI Surface**. Cobra is a surface concern and should not leak into Operational Domains.
_Avoid_: Domain dependency on Cobra

**Viper**:
The Go configuration library used as an implementation helper behind the **Configuration Resolver**. Viper is not the source of truth for Exito Tools configuration precedence.
_Avoid_: Viper-driven product behavior

**Bubble Tea**:
The Go framework used to build the **TUI Surface**. Bubble Tea is a surface concern and should not leak into Operational Domains.
_Avoid_: Domain dependency on Bubble Tea

**Output Format**:
The representation emitted by the **CLI Surface** for command results. JSON is the default Output Format; other formats may be selected when supported.
_Avoid_: View, rendering mode

**JSON Envelope**:
The standard JSON response shape emitted by the **CLI Surface**. Successful responses contain `ok: true`, `data`, and `meta`; failed responses contain `ok: false`, `error`, and `meta`. Standard metadata includes Request ID, optional Correlation ID, Profile, Capability ID, and duration.
_Avoid_: Raw command-specific JSON as the top-level response, secrets in metadata

**Envelope Metadata**:
The standard `meta` object in a **JSON Envelope**, containing `requestId`, optional `correlationId`, `profile`, `capabilityId`, and `durationMs`. It may include warnings, pagination, or deprecation metadata when relevant.
_Avoid_: Secrets or sensitive headers in metadata

**Pagination Metadata**:
Optional **Envelope Metadata** for list results, including cursor-based fields such as `nextCursor` and `hasMore`.
_Avoid_: Hidden auto-pagination metadata

**Structured Warning**:
A non-fatal machine-readable warning included in **Envelope Metadata**. Structured Warnings have stable codes, messages, and optional details, and do not change `ok: true` by themselves.
_Avoid_: Warning text mixed into stdout outside JSON

**Cursor**:
An opaque pagination marker used to request the next page of a list result. Cursors should be passed through, not interpreted by Exito Tools consumers.
_Avoid_: Offset when backend exposes cursor semantics

**Command Output Stream**:
The stdout stream used by CLI Commands for their primary output, including the JSON Envelope. Logs must not contaminate this stream.
_Avoid_: Debug logs mixed into JSON stdout

**Log Stream**:
The stderr stream or log file used for logs, diagnostics, and debug information. TUI logs should go to a file or non-disruptive diagnostic path rather than the main viewport.
_Avoid_: Logs as command output

**Structured Logging**:
Logging with levels and structured fields, implemented with Go `log/slog`. Structured Logging writes to the **Log Stream**, not the **Command Output Stream**.
_Avoid_: printf debugging in command output

**Operational Domain**:
A business capability area with its own language, use cases, and rules. Known operational domains include **Orders Domain** and **Geo Domain**, and more may be added later.
_Avoid_: Command category, technical module

**Domain Package**:
A Go package that implements an **Operational Domain**. Domain Packages expose Use Cases and Capabilities but do not depend on interaction surface frameworks.
_Avoid_: Cobra-aware domain package, Bubble Tea-aware domain package

**HTTP Infrastructure**:
Shared low-level networking support for base URLs, authentication headers, timeouts, retries, and request metadata. It does not own domain-specific API semantics.
_Avoid_: Generic layer containing domain business meaning

**Domain API Client**:
A domain-owned client for external APIs used by an **Operational Domain**. It uses shared **HTTP Infrastructure** but keeps endpoint DTOs and API semantics close to the domain.
_Avoid_: One giant shared API client for all domains

**External DTO**:
A data shape received from or sent to an external API. External DTOs belong near the relevant **Domain API Client** and should not be exposed as **Use Case Results**.
_Avoid_: API response as domain result

**Domain Model**:
A domain-owned representation used by Use Cases and Use Case Results. Domain Models express Exito Tools language rather than external API shapes.
_Avoid_: External DTO as model

**Surface Package**:
A Go package that adapts Capabilities into an **Interaction Surface**, such as CLI Commands or TUI Actions.
_Avoid_: Business logic owner

**Checkout Domain**:
The operational domain concerned with VTEX Checkout cart and orderForm assembly before an order is placed. It owns shopping-cart state, orderForm attachments, and checkout-step mutations.
_Avoid_: Treating cart assembly as Orders lookup, burying Checkout writes inside Catalog search

**VTEX Order Form**:
The VTEX Checkout cart state object that carries items, client profile data, shipping data, logistics selections, payment data, totals, and related checkout context before order placement. In Exito Tools, orderForm interactions belong to the **Checkout Domain**.
_Avoid_: OMS order, GEOMS order, generic form

**Checkout Session State**:
The provider-side and cookie-backed VTEX Checkout state associated with a VTEX Order Form. Exito Tools may create or update this state only through explicit confirmation-gated Checkout capabilities.
_Avoid_: Hidden browser cookie mutation, implicit shared cart state

**Checkout Attachment**:
A VTEX Checkout section attached to a VTEX Order Form, such as client profile data, shipping data, client preferences, marketing data, merchant context data, or payment data. Attachments are updated as explicit Checkout capabilities, not as side effects of product search.
_Avoid_: Random payload blob, surface-owned form state

**Purchase Assembly Flow**:
A guided workflow that prepares a VTEX Order Form for purchase by creating or loading the orderForm, adding selected items, updating customer and fulfillment attachments, and validating totals/options. Placing and paying the final order are separate high-risk steps unless explicitly approved.
_Avoid_: One hidden checkout macro, silent order placement

**Orders Domain**:
The operational domain concerned with orders.
_Avoid_: Orders commands

**Geo Domain**:
The operational domain concerned with geolocation.
_Avoid_: Geolocation commands, location utilities

**Geo Credentials**:
The sensitive provider credentials used by the **Geo Domain**. The standard token variable is `EXITO_GEO_TOKEN`; it belongs in environment variables or non-committed dotenv files.
_Avoid_: Geo token in YAML or committed docs

**Geo Base URL**:
The provider base URL used by the **Geo Domain**. The standard environment variable is `EXITO_GEO_BASE_URL`; endpoint paths remain in the Geo Domain API Client.
_Avoid_: Full endpoint URL as the only configuration point

**Use Case**:
A domain action exposed by an **Operational Domain** and reused by all relevant **Interaction Surfaces**. If a capability is available in both CLI and TUI, both surfaces invoke the same Use Case rather than implementing it separately.
_Avoid_: Screen action, command handler business logic

**Capability Execution**:
The act of running a **Capability** through its **Use Case**. Capability Execution is context-aware so callers can apply cancellation, timeouts, and request-scoped metadata.
_Avoid_: Uncancelable long-running operation

**Request ID**:
A unique identifier generated by Exito Tools for each **Capability Execution**. Request ID appears in JSON Envelope metadata, Structured Logging, and outbound HTTP metadata when applicable.
_Avoid_: User-supplied grouping identifier

**Correlation ID**:
An optional identifier supplied by an agent or user to group multiple **Capability Executions**. Correlation ID complements Request ID; it does not replace per-execution Request ID.
_Avoid_: Reusing one request ID for many executions

**Task Runner**:
A TUI-side execution helper that runs **Capability Execution** while representing loading, success, structured error, and cancelled states.
_Avoid_: Blocking the TUI event loop

**Capability**:
A neutral description of something an **Operational Domain** can do, backed by a **Use Case** and adaptable into surface-specific forms. The CLI Surface presents a Capability as a **Command**; the TUI Surface presents it as an **Action**.
_Avoid_: CLI command definition, TUI screen definition, tool

**Workflow Capability**:
A **Capability** that orchestrates multiple **Operational Domains** without belonging cleanly to one of them. Workflow Capabilities are registered like domain Capabilities but live outside Domain Packages; their visible IDs should use business names rather than a technical `workflow.*` prefix.
_Avoid_: Forcing cross-domain orchestration into an unrelated domain, `workflow.*` as user-facing ID

**Application Wiring**:
The explicit composition layer that creates dependencies, loads Operational Domains, and registers their Capabilities.
_Avoid_: Side-effect imports, hidden auto-registration

**Capability Registry**:
The runtime catalog of registered **Capabilities** used by interaction surfaces and discovery commands. Capabilities enter the registry through **Application Wiring** during boot, then the registry is treated as immutable.
_Avoid_: Global implicit registry populated by imports, runtime capability mutation

**Capability ID**:
A stable machine-readable identifier for a **Capability**, using `<domain>.<action>` with lower-case kebab-case action names when needed. Incompatible changes use a new versioned ID such as `.v2` rather than silently changing the old contract.
_Avoid_: Display title as identifier, casually renamed ID

**Capability Version**:
Metadata describing the version of a **Capability** contract. Compatible changes may increment version metadata while keeping the same **Capability ID**; incompatible changes require a new versioned Capability ID.
_Avoid_: Silent breaking change

**Capability Visibility**:
Metadata that declares where a **Capability** should appear, such as CLI, TUI, or Command Palette. Visibility prevents technical or agent-only capabilities from leaking into human-facing flows.
_Avoid_: Showing every capability everywhere by default

**Capability Audience**:
Metadata that declares the intended consumers of a **Capability**, such as agents, people, or both. Agent-only capabilities remain machine-accessible but are not promoted in human-facing TUI flows unless explicitly visible there.
_Avoid_: Assuming all capabilities are for all users

**Agent-only Capability**:
A **Capability** intended for agents or automation rather than everyday people. Agent-only does not mean unsafe; it means the TUI should not promote it as a human-facing Action by default.
_Avoid_: Hidden business capability, unsafe command

**Risk Level**:
Metadata that classifies a **Capability** by operational risk, such as read-only, safe-write, or destructive. Interaction surfaces use Risk Level to decide when confirmation is required.
_Avoid_: Treating all capabilities as equally safe

**Confirmation Requirement**:
Metadata indicating whether a **Capability** requires explicit confirmation before execution. Risky CLI Commands fail with `CONFIRMATION_REQUIRED` instead of prompting when confirmation is missing; TUI Actions may require an impact-aware confirmation prompt.
_Avoid_: Silent destructive execution, interactive CLI confirmation prompt

**Confirmation Error**:
A **Structured Error** emitted when a risky Command is invoked without required confirmation. The standard Error Code is `CONFIRMATION_REQUIRED`.
_Avoid_: Hanging for interactive confirmation in CLI

**Input Schema**:
A neutral description of the inputs required by a **Capability**. Interaction surfaces use the Input Schema to collect and validate values without redefining the same inputs separately.
_Avoid_: CLI flags as source of truth, TUI form as source of truth

**Use Case Result**:
The structured, format-neutral output returned by a **Use Case**. It is presented by interaction surfaces and output presenters, not pre-rendered by the domain.
_Avoid_: Pretty string, table output, TUI view model as domain result

**Presenter**:
A surface or format-specific transformation from a **Use Case Result** into something a consumer can read, such as JSON, Markdown, HTML, Toon, or a TUI view.
_Avoid_: Domain formatter

**Structured Error**:
A stable error representation with a machine-readable code, human-readable message, structured details, and retryability information when applicable.
_Avoid_: Plain error string, panic output

**Error Translation**:
The conversion of technical, external, or domain failures into **Structured Errors**. The layer that understands the failure's meaning owns the translation.
_Avoid_: Translating all errors at the CLI/TUI boundary

**Error Code**:
A stable machine-readable identifier for a **Structured Error**. Error Codes are part of the automation contract and should not change casually.
_Avoid_: Error message as identifier

**Exit Code**:
The process-level result emitted by a **Command** for coarse automation decisions. Exit Codes stay generic; detailed semantics belong in **Error Codes**.
_Avoid_: Domain-specific exit code for every failure

**Action**:
A user-facing task people can discover and start from the **TUI Surface**, usually backed by a **Use Case**. Primary TUI navigation uses curated Actions; the Command Palette may expose additional suitable Actions.
_Avoid_: Tool, exposing every capability as a primary action

**Command**:
A machine-facing invocation exposed by the **CLI Surface**, usually backed by a **Use Case**. Explicit domain Commands follow `exito <domain> <action> [flags]` and map to a **Capability ID**.
_Avoid_: Tool, command names disconnected from Capability IDs

**Action Verb**:
The verb used in a **Command** or **Capability ID** to describe the operation. Standard verbs are `get`, `list`, `search`, `validate`, `diagnose`, and explicit mutating verbs such as `create`, `update`, or `delete` when state changes.
_Avoid_: check, consult, show, fetch unless a domain reason exists

**Root Command**:
The bare `exito` invocation. It shows brief textual help and does not open the TUI or emit domain command output.
_Avoid_: Implicit TUI launch, hidden command execution

**Capabilities Command**:
A machine-readable Command for discovering available **Capabilities**, including their domain, descriptions, schemas, and mappings to CLI Commands and TUI Actions. It is the explicit discovery path for agents.
_Avoid_: Scraping root help for automation, listing only CLI commands

**Run Command**:
A generic agent-oriented Command for executing a registered **Capability** by **Capability ID**, expected to look like `exito run <capability-id>`. It receives a complete input object through `--input-json`, `--input-file`, or piped stdin, and reuses the same execution path as explicit domain Commands.
_Avoid_: Duplicate execution path, separate business logic, dynamic per-field flags in run

**TUI Surface**:
The task-first interaction surface intended for everyday people using a guided terminal experience. It may use domains internally, but its navigation should prioritize what the person wants to accomplish.
_Avoid_: User app, visual CLI, domain-first menu

**TUI Entrypoint**:
The explicit Command that starts the **TUI Surface**, expected to be `exito tui`. Exito Tools should not silently open the TUI when `exito` is run without arguments.
_Avoid_: Implicit interactive startup

**Command Palette**:
A TUI discovery mechanism that helps people find available **Actions** across domains. It is distinct from filtering already-loaded results.
_Avoid_: Filter, menu search, tool finder

**Result Filter**:
A refinement mechanism applied to a set of results already shown within a task. It is distinct from discovering actions in the **Command Palette**.
_Avoid_: Command search, tool finder

## Example Dialogue

**Dev**: Are we building a CLI and a TUI as separate tools?

**Domain Expert**: No. We are building Exito Tools as one application with two interaction surfaces: agents use the CLI Surface, and people use the TUI Surface. Both surfaces expose Operational Domains such as Orders and Geo, and both call the same Use Cases when they offer the same capability.
