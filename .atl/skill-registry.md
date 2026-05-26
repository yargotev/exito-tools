# Skill Registry

**Delegator use only.** Any agent that launches sub-agents reads this registry to resolve compact rules, then injects them directly into sub-agent prompts. Sub-agents do **not** need to read this registry or individual `SKILL.md` files unless the launch prompt explicitly asks them to.

Last refined: 2026-05-26  
Project: `exito-tools`

## Selection stance for this project

Exito Tools is a Go CLI/TUI application using Cobra, Bubble Tea, OpenSpec/SDD artifacts, strict domain/surface separation, JSON-first CLI contracts, and retrieval-led local documentation. Prefer skills that reinforce those constraints.

### Core skills to use often

| Trigger | Skill | Path |
| --- | --- | --- |
| Go tests, coverage, Bubble Tea model tests, teatest, golden files | `go-testing` | `/home/yerson-alexander-argote-vasquez/.codex/skills/go-testing/SKILL.md` |
| Library/framework/API docs or code examples | `documentation-lookup` | `/home/yerson-alexander-argote-vasquez/.agents/skills/documentation-lookup/SKILL.md` |
| Cobra command tree work in `internal/surface/cli` | `golang-spf13-cobra` | `.agents/skills/golang-spf13-cobra/SKILL.md` |
| Stress-test a plan against project terminology and docs | `grill-with-docs` | `.agents/skills/grill-with-docs/SKILL.md` |
| Significant code/design review, "judgment day", adversarial/dual review | `judgment-day` | `/home/yerson-alexander-argote-vasquez/.codex/skills/judgment-day/SKILL.md` |
| Architecture, ADRs, guides, onboarding, review-facing docs | `cognitive-doc-design` | `/home/yerson-alexander-argote-vasquez/.config/opencode/skills/cognitive-doc-design/SKILL.md` |
| Creating issues, bug reports, feature requests | `issue-creation` | `/home/yerson-alexander-argote-vasquez/.codex/skills/issue-creation/SKILL.md` |
| Creating/opening/preparing PRs | `branch-pr` | `/home/yerson-alexander-argote-vasquez/.codex/skills/branch-pr/SKILL.md` |
| Implementation work that may need clean commits or slices | `work-unit-commits` | `/home/yerson-alexander-argote-vasquez/.config/opencode/skills/work-unit-commits/SKILL.md` |
| PRs over 400 changed lines or stacked/chained review slices | `chained-pr` | `/home/yerson-alexander-argote-vasquez/.config/opencode/skills/chained-pr/SKILL.md` |

### Useful situational skills

| Trigger | Skill | Path |
| --- | --- | --- |
| Need to find/install more skills | `find-skills` | `/home/yerson-alexander-argote-vasquez/.agents/skills/find-skills/SKILL.md` |
| Create or update an AI agent skill | `skill-creator` | `/home/yerson-alexander-argote-vasquez/.codex/skills/skill-creator/SKILL.md` |
| Audit or refactor existing skills | `skill-improver` | `/home/yerson-alexander-argote-vasquez/.config/opencode/skills/skill-improver/SKILL.md` |
| Human-facing PR/issue/review/Slack comments | `comment-writer` | `/home/yerson-alexander-argote-vasquez/.config/opencode/skills/comment-writer/SKILL.md` |
| Turn current conversation context into a PRD | `to-prd` | `.agents/skills/to-prd/SKILL.md` |
| OpenAI API/product docs | `openai-docs` | `/home/yerson-alexander-argote-vasquez/.codex/skills/.system/openai-docs/SKILL.md` |
| Azure DevOps CLI automation | `azure-devops-cli` | `/home/yerson-alexander-argote-vasquez/.agents/skills/azure-devops-cli/SKILL.md` |
| Raster image generation/editing | `imagegen` | `/home/yerson-alexander-argote-vasquez/.codex/skills/.system/imagegen/SKILL.md` |
| Codex plugin scaffolding | `plugin-creator` | `/home/yerson-alexander-argote-vasquez/.codex/skills/.system/plugin-creator/SKILL.md` |
| Install Codex skills | `skill-installer` | `/home/yerson-alexander-argote-vasquez/.codex/skills/.system/skill-installer/SKILL.md` |
| Project-local Cobra CLI guidance | `golang-spf13-cobra` | `.agents/skills/golang-spf13-cobra/SKILL.md` |

### Deliberately not selected by default

| Skill | Reason |
| --- | --- |
| `commit` | Its `TYPE(PAGE): MESSAGE` convention targets Exito storefront pages, not this Go CLI/TUI repo. Use `AGENTS.md` and repo history instead. |
| `vercel-react-best-practices` | Useful for React/Next.js repositories, not for this Go CLI/TUI codebase. |
| `url-via-markdown-new` | Useful in Claude/WebFetch flows, but Codex already has web browsing and source-citation rules. |
| External generic code-review skills | `judgment-day` already provides a project-aware dual-review protocol and resolves compact rules from this registry. |

## Compact Rules

Delegators copy matching blocks into sub-agent prompts as `## Project Standards (auto-resolved)`.

### go-testing
- Prefer table-driven Go tests for multiple inputs, edge cases, and error paths.
- Test external behavior and stable contracts, not implementation details.
- For Bubble Tea, test `Model.Update` state transitions directly before using terminal snapshots.
- Use `teatest` only for full interactive flows; keep model tests fast and focused.
- Use `t.TempDir()` for filesystem tests and fake servers/clients for network boundaries.
- Golden files are appropriate for stable view output; avoid brittle snapshots for behavior.
- Run the smallest relevant `go test` first, then `make test` before declaring implementation complete.

### documentation-lookup
- Use current official/library docs for framework, API, or code-example questions instead of relying on memory.
- Resolve the exact library/package and prefer version-specific docs when the user mentions a version.
- Prefer primary/official docs over community forks.
- Keep fetched examples narrow and adapt them to Exito Tools architecture rather than copying blindly.
- For OpenAI product/API questions, use `openai-docs` instead.

### golang-spf13-cobra
- Apply only to Cobra surface code, mainly `internal/surface/cli`; domain packages must stay Cobra-free.
- Use `RunE`/`*E` hooks so errors return through the CLI error/JSON envelope path; avoid `Run`, `panic`, or `os.Exit` in handlers.
- Set `SilenceUsage` and `SilenceErrors` where appropriate so Exito Tools controls structured errors and stdout/stderr discipline.
- Validate positional args with Cobra `Args` validators, not ad hoc `len(args)` checks inside `RunE`.
- Use `cmd.OutOrStdout()` and `cmd.ErrOrStderr()` for testable output; never contaminate stdout JSON with diagnostics.
- Build fresh command trees per test because Cobra accumulates state across `Execute()` calls.
- Treat Viper guidance from the upstream skill as advisory only: Exito Tools keeps Viper behind the Configuration Resolver and product-defined precedence.

### grill-with-docs
- Before architecture or terminology changes, read `CONTEXT.md`, relevant ADRs, specs, and code.
- Challenge terms that conflict with the glossary; propose precise canonical language.
- Explore code when a question can be answered locally instead of asking the user.
- Update `CONTEXT.md` immediately when domain language is resolved; keep it glossary-only.
- Create ADRs only for hard-to-reverse, surprising, trade-off decisions.
- Ask one design question at a time and provide a recommended answer.

### judgment-day
- Use only for explicit adversarial/dual review requests or high-confidence review gates.
- Resolve this registry first and inject relevant compact rules into both judges and the fix agent.
- Launch two independent blind judge agents in parallel with identical criteria.
- Classify findings as confirmed, suspect, contradiction, or clean; do not let judges collaborate.
- Fix confirmed criticals/real warnings narrowly, then re-judge as prescribed.
- Treat theoretical warnings as informational unless a normal user can trigger them.

### cognitive-doc-design
- Lead with the answer, decision, or action; move context after the quick path.
- Use progressive disclosure: happy path first, then details, edge cases, references.
- Prefer tables, checklists, examples, and templates over dense prose.
- Make reviewer path explicit: what to review first, out of scope, verification, next step.
- Keep sections small and signposted; optimize for fast scanning and low recall burden.
- Avoid documenting implementation trivia that will rot quickly.

### issue-creation
- Search for duplicates before creating an issue.
- Use the repo issue template; fill all required fields and pre-flight checks.
- Issues start as `status:needs-review`; PRs require `status:approved` first when that workflow applies.
- Bug reports need reproducible steps, expected behavior, actual behavior, environment, and logs when relevant.
- Feature requests need problem, user-facing solution, affected area, alternatives, and context.
- Questions should go to discussion/chat channels, not issues.

### branch-pr
- Check `git status --short` before PR work and do not overwrite user changes.
- Every PR must link the relevant approved issue if the repo enforces issue-first PRs.
- Use a focused branch name like `feat/<description>` or `fix/<description>` when creating branches.
- PR body should include linked issue, type, summary, changes table, and test plan.
- Add exactly one PR type label when required by the repo.
- Run relevant verification before PR handoff and keep commit messages free of AI attribution.

### work-unit-commits
- Split implementation into reviewable work units with tests/docs beside the code they verify.
- Keep commits focused: one behavior or artifact slice per commit.
- Avoid `git add -A`; stage explicit files after reviewing diffs.
- Do not mix unrelated formatting, refactors, and behavior changes in one commit.
- Use conventional, concise commit messages without AI/co-author attribution.
- If a slice grows too large, switch to `chained-pr` planning before PR handoff.

### chained-pr
- Split PRs over roughly 400 changed lines unless a maintainer accepts a size exception.
- Keep each child PR independently reviewable and focused on one deliverable work unit.
- Preserve clean diffs by retargeting or rebasing when parent changes pollute a child PR.
- State PR order, dependency diagram, current boundary, verification, and out-of-scope work.
- Do not mix stacked-main and feature-branch-chain strategies after choosing one.

### find-skills
- Identify the domain/task first, then search with specific keywords using `npx skills find <query>`.
- Present found skills with purpose, install command, and skills.sh link.
- Do not install skills silently; install only after user approval or an explicit install request.
- If no skill fits, say so and proceed with general capabilities or suggest creating a project-specific skill.
- After installing/removing skills, refresh this registry.

### skill-creator
- Use when creating durable AI instructions for recurring workflows, not one-off notes.
- Create a valid `SKILL.md` with clear frontmatter name/description and trigger language.
- Keep instructions action-oriented: rules, workflows, examples, gotchas.
- Prefer bundling scripts/templates/assets over retyping large repeatable content.
- Test that the skill trigger is specific enough and does not collide with existing skills.

### skill-improver
- Audit skills for trigger clarity, compact actionable rules, missing gotchas, and stale assumptions.
- Remove fluff, motivation, and redundant examples that do not affect execution.
- Preserve author intent and compatibility with existing tool/agent conventions.
- Prefer small, reviewable edits and validate frontmatter after changes.
- Update the registry after changing skill behavior or paths.

### comment-writer
- Start with the actionable point; avoid long recaps.
- Be warm and direct, like a teammate, not a corporate bot.
- Keep comments short: 1-3 paragraphs or a tight bullet list.
- Explain why when requesting a change.
- Avoid pile-ons; comment on the highest-value issues.
- Match thread/user language; avoid em dashes.

### to-prd
- Synthesize from current conversation and codebase context; do not over-interview.
- Use project glossary vocabulary and respect ADRs/specs.
- Describe modules and decisions at product/design level, not volatile file-by-file implementation detail.
- Include user stories, implementation decisions, testing decisions, out of scope, and further notes.
- Check with the user on major module boundaries and test focus before publishing if uncertainty remains.

### openai-docs
- Use for OpenAI API/product/model questions; prefer official OpenAI docs and citations.
- Check local/bundled references first where available, browse official OpenAI domains as fallback.
- Verify latest model/API details because they change frequently.
- Keep recommendations tied to the user's use case and include upgrade implications when relevant.

### azure-devops-cli
- Use for Azure DevOps projects, repos, pipelines, builds, PRs, work items, artifacts, and service endpoints.
- Prefer `az devops`/`az pipelines` CLI workflows and inspect current organization/project context first.
- Avoid destructive Azure changes without explicit user approval.
- Report exact commands and IDs used for reproducibility.

### imagegen
- Use for raster images/photos/illustrations/textures/sprites/mockups, not SVG/code-native UI assets.
- Use the built-in image tool by default; CLI fallback only when explicitly requested or confirmed.
- For project-bound assets, move/copy final files into the workspace and do not overwrite existing assets without permission.
- For transparent outputs, use chroma-key plus local alpha removal unless true native transparency is explicitly confirmed.

### plugin-creator
- Create Codex plugins with a required `.codex-plugin/plugin.json` and valid manifest defaults.
- Prefer personal-marketplace entries by default when scaffolding local plugins.
- Reuse templates/scripts from the skill rather than hand-writing large boilerplate.
- Use cache-buster/reinstall flow when updating a plugin during development.

### skill-installer
- Use to list curated installable Codex skills or install a skill from a GitHub repo path.
- Confirm target location and scope before installing if the user did not specify.
- After install/remove/update, refresh this registry.
- Do not treat installation as endorsement; evaluate relevance before making it core.

## External candidates found with `$find-skills`

Reviewed on 2026-05-26 and installed only the candidate that was useful without conflicting with Exito Tools' architecture.

| Candidate | Decision | Reason |
| --- | --- | --- |
| `samber/cc-skills-golang@golang-spf13-cobra` | **Installed locally** with `DISABLE_TELEMETRY=1 npx -y skills@latest add 'samber/cc-skills-golang@golang-spf13-cobra' --agent codex --copy -y` | Useful focused Cobra command-tree guidance. Registry compact rules constrain it to `internal/surface/cli` and override upstream Viper assumptions with Exito Tools' Configuration Resolver rule. |
| `gentleman-programming/engram@gentleman-bubbletea` | Not installed | Too specific to Gentleman.Dots installer paths/screens; current `go-testing` already covers Bubble Tea test patterns and Exito Tools has its own TUI language. |
| `samber/cc-skills-golang@golang-cli` | Not installed | Useful generic CLI advice, but it repeatedly assumes direct Cobra+Viper binding and generic project layout that can conflict with Exito Tools' JSON-first CLI and Configuration Resolver contracts. |
| `eduardo-sl/go-agent-skills@go-architecture-review` | Not installed | Good dependency-direction checklist, but its generic service/handler/store layout and env-only config advice conflict with local package layout and YAML Configuration Resolver decisions. |
| `addyosmani/agent-skills@documentation-and-adrs` | Not installed | Good generic docs advice, but it recommends `docs/decisions/` and examples that conflict with this repo's `docs/adr/`, `CONTEXT.md`, OpenSpec, and existing doc skills. |

Searches run on 2026-05-26: `go testing`, `cli`, `tui`, `documentation`, `code review`, `bubbletea`, `cobra cli go`, `golang architecture`, `adr architecture`, `github issue pr`.

## Project Conventions

| File | Path | Notes |
| --- | --- | --- |
| `AGENTS.md` | `AGENTS.md` | Project-level agent instructions and local docs index. |
| `CONTEXT.md` | `CONTEXT.md` | Domain language, product terms, architectural vocabulary. |
| `docs/prd.md` | `docs/prd.md` | Product roadmap and intended capability sequence. |
| `docs/package-layout.md` | `docs/package-layout.md` | Package boundaries and dependency rule. |
| `docs/adr/0024-go-package-layout.md` | `docs/adr/0024-go-package-layout.md` | Current package-layout ADR explicitly referenced by AGENTS.md. |
| `docs/adr/*.md` | `docs/adr/*.md` | Specific architecture decisions. |
| `openspec/specs/*/spec.md` | `openspec/specs/*/spec.md` | Current source-of-truth requirements. |
| `openspec/changes/<change-id>/` | `openspec/changes/<change-id>/` | Active change proposals, specs, designs, tasks. |
| `docs/configuration.md` | `docs/configuration.md` | Configuration behavior. |
| `openspec/specs/configuration-resolver/spec.md` | `openspec/specs/configuration-resolver/spec.md` | Configuration Resolver requirements. |
| `docs/capabilities/*.md` | `docs/capabilities/*.md` | Capability contracts. |
| `openspec/specs/capability-*/spec.md` | `openspec/specs/capability-*/spec.md` | Capability requirements. |

## Project-specific rules to inject even without a matching skill

- Start with `git status --short`; stop and ask if unexpected user changes appear.
- Keep changes narrow and reviewable; avoid repo-wide rewrites unless explicitly requested.
- Domain packages under `internal/domain/*` must not import Cobra, Bubble Tea, or `internal/surface/*`.
- Surface packages adapt capabilities into CLI/TUI behavior; business behavior belongs in domain use cases or shared execution/application layers.
- Use explicit wiring in `internal/app`; do not introduce hidden side-effect registration.
- Treat the capability registry as immutable after boot and preserve stable Capability IDs.
- CLI stdout is for JSON envelopes; logs and diagnostics go to stderr or non-disruptive TUI paths.
- Long-running execution must be context-aware and cancellable.
- Risk/destructive actions require explicit confirmation; non-interactive CLI flows must not silently proceed.
- User-facing labels, help, errors, and visible contracts remain English until i18n exists.
