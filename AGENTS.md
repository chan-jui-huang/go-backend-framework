# AGENTS.md

This document guides AI agents and contributors working in this repository. It defines project structure, conventions, workflows, and safety rules. Its scope covers the entire repository. If a more deeply nested AGENTS.md exists, it overrides conflicting guidance for its subtree.

## Repository Structure & Assets
- `bin/`: Build artifacts emitted by the Makefile targets (created on demand).
- `cmd/app`: Main HTTP entrypoint and wiring; builds to `bin/app`.
- `cmd/kit`: Helper CLIs (`jwt`, `http_route`, `rdbms_seeder`, `permission_seeder`).
- `cmd/template`: Minimal bootstrap executable that wires registrars; useful for scaffolding.
- `internal`: Domain logic and supporting modules with colocated tests (`*_test.go`).
  - `http/`: Gin HTTP stack — `controller/`, `middleware/`, `response/`, `route/`, `server.go`.
  - `scheduler/`: Background jobs (cron-like tasks) and example jobs under `job/`.
  - `registrar/`: Service and dependency registration for bootstrapping.
  - `pkg/`: Project-specific packages (database, models, permission logic, user domain helpers, etc.).
  - `migration/`: Database migrations split into `rdbms/` and `clickhouse/` trees.
  - `test/`: Fixtures and helpers for tests (admin, permission, migration, HTTP utilities).
- `docs`: Swagger outputs (`swagger.json|yaml`, `docs.go`).
- `deployment/`: Docker assets (e.g., `deployment/docker/`).
- `storage`: Runtime artifacts such as `storage/log/access.log` and `storage/log/app.log`.
- Configuration & env templates live at the root: `.env*`, `.air.toml`, `config*.yml`, `.golangci.yml`.

## Key Technologies
- Language: Go (Gin web framework).
- Data: GORM (MySQL/PostgreSQL/SQLite), ClickHouse, Redis.
- Config: Viper (file-based + environment variables).
- AuthN/Z: JWT with Casbin-based authorization policies.
- Docs: Swagger annotations generate specs in `docs/`.
- Scheduling: Cron-style background tasks.

## Core Architecture Concepts

**Foundation:** The framework is built on **[go-backend-package](https://github.com/chan-jui-huang/go-backend-package)** with 12 modular packages and zero cross-dependencies. For comprehensive architectural details, refer to: **[go-backend-package AGENTS.md](https://github.com/chan-jui-huang/go-backend-package/blob/main/AGENTS.md)**

### Dual-Registry System
- **Config Registry**: YAML configuration with environment variable expansion (`${VAR_NAME}`)
- **Service Registry**: Initialized service instances accessed via `service.Registry.Get(key)`

### Bootstrap Sequence

```
loadEnv() → bootConfig() → BeforeExecute() → Execute() → AfterExecute()
                           [Registrar Boot() + Register()]
```

### Application Lifecycle Phases (pkg/app)

| Phase | Type | Blocking | Purpose |
|-------|------|----------|---------|
| STARTING | Sequential | Yes | Pre-startup setup and initialization |
| EXECUTION | Goroutine | No | Main application logic runs here |
| STARTED | Sequential | Yes | Post-startup hooks and callbacks |
| SIGNALS | Goroutines | Yes | OS signal handling for graceful shutdown |
| ASYNC | Goroutines | No | Background tasks and scheduled jobs |
| TERMINATED | Sequential | Yes | Cleanup and resource release before exit |

See `cmd/app/main.go` for callback registration examples.

### Key Architectural Patterns

| Pattern | Purpose |
|---------|---------|
| **Registrar Pattern** | Modular initialization via `Boot()` + `Register()` methods |
| **Dependency Injection** | Services registered and retrieved via Service Registry |
| **Configuration Management** | YAML + `.env` file support with variable expansion |
| **Database Layer** | Goose migrations + GORM ORM wrapper for type-safe operations |
| **Authentication** | Ed25519 JWT (quantum-resistant) generated via `./bin/jwt` |
| **Authorization** | Casbin RBAC policies seeded via `cmd/kit/permission_seeder` |
| **Logging** | Zap structured logging with rotation to `storage/log/` |

## Business Logic Layer Architecture

### Separation of Concerns
This framework enforces a strict separation between business logic and presentation/interface layers:

- **`internal/pkg/`**: Contains all reusable business logic that can be invoked by any interface layer (HTTP APIs, schedulers, CLI commands, etc.). Functions here are pure business operations that handle data manipulation, database transactions, and domain-specific logic.

- **`internal/http/`**: Contains HTTP-specific handlers, middleware, routing, and request/response structures. Controllers in this layer orchestrate HTTP interactions but delegate actual business operations to `internal/pkg/`.

- **Other Interface Layers**: Schedulers (`internal/scheduler/`), CLI commands (`cmd/kit/`), and any future interfaces (gRPC, WebSocket, etc.) should similarly call functions from `internal/pkg/` rather than duplicating logic.

### API Development Workflow

When developing a new API endpoint, follow this mandatory pattern:

1. **Business Logic First**:
   - Check if the required business logic exists in `internal/pkg/<domain>/`
   - If not, create or extend functions in `internal/pkg/<domain>/` first
   - Example: For user operations, add functions to `internal/pkg/user/user.go`
   - Keep these functions framework-agnostic and reusable

2. **Controller Implementation**:
   - Create or extend handler files under `internal/http/controller/<area>/`
   - Controllers should ONLY:
     - Validate HTTP requests
     - Call business logic from `internal/pkg/`
     - Format HTTP responses
   - Never put business logic directly in controllers

3. **Testing Requirements**:
   - Write comprehensive tests for all controller return scenarios
   - Test coverage must include:
     - Success case (200 OK with expected data)
     - All error cases (validation failures, authentication/authorization failures, business logic errors)
     - Edge cases and boundary conditions
   - Reference existing tests in `internal/http/controller/**/*_test.go` for patterns
   - Example: `internal/http/controller/user/user_register_test.go` demonstrates testing multiple return paths

### Example Pattern

**Bad** - Business logic in controller:
```go
func CreateUser(c *gin.Context) {
    // DON'T: Database operations directly in controller
    db.Table("users").Create(&user)
}
```

**Good** - Calls business logic from internal/pkg:
```go
func CreateUser(c *gin.Context) {
    // DO: Delegate to business logic layer
    if err := user.Create(database.NewTx(), u); err != nil {
        // Handle error response
    }
    // Return success response
}
```

### Benefits of This Architecture

- **Reusability**: Business logic in `internal/pkg/` can be called from APIs, schedulers, CLI tools, or future interfaces
- **Testability**: Business logic can be unit tested independently from HTTP concerns
- **Maintainability**: Changes to business rules are centralized in `internal/pkg/`
- **Consistency**: All interfaces use the same validated business logic

## Build, Run, and Tooling
- `make`: Compiles the main service and helper binaries with the `jsoniter` build tag enabled.
- `make all`: Builds `bin/app` alongside helper CLIs.
- `make run`: Runs the service locally with the race detector enabled.
- `go run cmd/app/*`: Spins up the API without the Makefile; use `make air` for hot reloading (requires the `air` tool).
- CLI builds: `make jwt`, `make http_route`, `make rdbms_seeder`, `make permission_seeder` build individual helpers into `bin/`.
- Docs: `make swagger` regenerates Swagger artifacts after updating annotations.
- Quality: `make linter` (golangci-lint with `errcheck` and `gosec`) and `golangci-lint run ./...` should pass before commits.
- Tests & benchmarks: `make test [args=./...]`, `make benchmark`, or `go test ./...` (ensure SQLite and `.env.testing` are available).

## Coding Style & Naming Conventions
- Run `gofmt`, `goimports`, and `go vet ./...` prior to review. Indentation must use tabs (gofmt default, width 4 spaces).
- Keep package names lowercase with no underscores; prefer descriptive names (`internal/http/controller/user`).
- Use natural exported identifiers (e.g., `UserLoginRequest`). For camelCase names containing "id", spell it as `Id` (e.g., `userId`).
- Avoid import aliases unless required for conflict resolution or clarity; rely on the default import name whenever possible.
- Reuse established helper patterns (e.g., response builders in `internal/http/response`) rather than handcrafting JSON payloads.
- Respect the `jsoniter` build tag where the Makefile enables it.
- Follow directory-local conventions and keep reusable code in `internal/pkg/`; leave HTTP wiring and other application layers under `internal/`.
- API route paths must use kebab case (e.g., `pet-store`).
- Struct names: singular, PascalCase (e.g., `User`, `SchedulerJob`); use another casing only when a schema or user directive explicitly requires it.
- Table tags: let GORM infer plural table names; only add a `TableName()` override when migrations already created a non-standard name (e.g., `AuditLog` mapping to `"audit_logs"`).
- Slice fields for relations: use plural names for collections (e.g., `Roles []Role`, `Permissions []Permission`).
- Single relations: keep singular names (e.g., `Role *Role`, `Profile *Profile`).
- Many-to-many join structs: keep the join struct singular (e.g., `UserRole` for explicit join rows). Parent structs still use plural slices for the many-to-many relation (e.g., `User` has `Roles []Role`, `Role` has `Users []User`, both pointing to the same join table or struct).
- Foreign key fields: camelCase with an `Id` suffix so GORM can map to `snake_case` columns (e.g., `UserId uint`, `RoleId uint`); add a `gorm:"column:..."` tag only when the column name diverges.
- Packages: lower-case with no underscores. Files use `snake_case.go` only when it clarifies multiword names already used in the codebase.
- If a directory name includes underscores for readability (e.g., `user_activity`), keep the Go package name lower-case without underscores and import with an explicit alias only when needed.
- Follow the naming/style conventions of files in the same directory before introducing new patterns.
- Helper/query patterns: when adding or modifying helpers, match the existing query shape in that file/package (ordering, preload use, error wrapping). Keep helpers minimal—do not add joins/preloads or extra queries unless the file already does or the requirement demands it.
- Swagger `// @param` comments must keep the description placeholder as `" "` (do not insert real tokens like `"id"`).
- Keep Swagger annotations aligned with actual runtime responses; if a handler logs an error instead of returning it, remove the unused error entry from Swagger.
- When a foreign key guarantees a required parent (e.g., `User` on `UserPermission`), do not add nil-guards before using the preloaded relation.

### Example: indentation must use tabs
```golang
// use tab is Good example
package main

func main() {
	fmt.Println("example")
}

// use space is Bad example
package main

func main() {
  fmt.Println("example")
}
```

## Error Responses
- Define canonical error messages and codes exclusively in `internal/http/response/error_message.go`.
- Message constants describe user-facing text. Update `MessageToCode` with a unique `<status>-<sequence>` string (e.g., `400-001`).
- Organize the map by HTTP status headers (`// 400`, `// 401`, etc.) so related errors stay grouped.
- Ensure handlers return responses whose status aligns with the declared message and add corresponding tests when introducing new errors.

## Development Workflow
When adding a new API endpoint:
1. Extend or create handler files under `internal/http/controller/<area>` (for example, `internal/http/controller/user/user_register.go`). Follow the existing `<feature>_<action>.go` naming and add matching `*_test.go` coverage.
2. **REQUIRED**: Implement or extend business logic in `internal/pkg/<domain>` (e.g., `internal/pkg/user`, `internal/pkg/permission`) BEFORE writing controller code. Controllers MUST delegate to `internal/pkg/` and never contain business logic directly. Keep helpers reusable and include unit tests where practical.
3. Wire routes in the corresponding router under `internal/http/route/<area>/api_route.go`, ensuring it implements `route.Router` and guards handlers with middleware as needed.
4. If you introduce a brand-new router, register it in `internal/http/route/api_route.go` by appending it to the `routers` slice so it participates in `AttachRoutes`.
5. For admin/protected capabilities, synchronise seeds between runtime and tests: update `cmd/kit/permission_seeder/permission_seeder.go` and mirror the same permissions/roles in `internal/test/permission_service.go` (and adjust `internal/test/admin_service.go` when admin fixtures change).
6. Whenever you introduce a new mock service (gomock or hand-written), also register an empty instance in `internal/test/test.go:emptyMockedServices()` so the test harness has placeholders for injected services.

General guidance:
- Keep changes minimal and localized; avoid unrelated refactors.
- Favor composition over duplication; reuse helpers under `internal/pkg/` and existing shared utilities.
- Update Swagger comments whenever API shapes change, then run `make swagger`.
- Wrap work in a transaction only when multiple insert/update/delete statements need to succeed together.

## Testing Guidelines
- Test framework: standard `testing`; `testify` is available for assertions.
- Location: place `*_test.go` files alongside the code they cover; prefer table-driven tests or testify suites.
- Bootstrapping: leverage `internal/test` utilities for environment loading, seeded users, migrations, and CSRF helpers instead of reimplementing setup logic.
- Mock services: when introducing a new mock service, also register it in `internal/test/test.go` within `emptyMockedServices` so the test harness initializes it.
- Migrations: run required migrations first (e.g., `make sqlite-migration args=up`). Tests expect `.env.testing` to provide DSNs and secrets.
- Coverage: exercise both success and failure paths. **REQUIRED**: Test all possible return values from controllers (success, validation failures, auth failures, business errors). Reference existing patterns in `internal/http/controller/**/*_test.go`. Document skipped integration tests in PR descriptions. Add controller/service tests whenever adding new endpoints.

## Migrations
- Use Make targets with environment loaded from `.env`: `make mysql-migration args="up"`, `make pgsql-migration args="up"`, `make sqlite-migration args="up"`, or `make clickhouse-migration args="up"`.
- Ensure relevant `DB_*` variables are set for the target database.
- Keep migrations idempotent and reversible; document non-trivial data movements.
- Migration filenames MUST use a timestamp reflecting the actual creation time; generate a fresh `YYYYMMDDhhmmss` prefix when the file is added, and ensure each migration file has a unique timestamp prefix (do not reuse the same `YYYYMMDDhhmmss` across multiple migration files). When creating paired migrations (e.g., prod/test), keep the timestamp identical across the pair while still unique per table/change.

## Security & Configuration
- Never commit secrets. Place credentials in `storage/...` as described in the README.
- Provide `.env.dev` and `.env.testing` locally; avoid committing real values.
- Validate configuration via `make run` and integration tests before deployment.

## Commit & PR Guidelines
- Follow Conventional Commits (`feat(scope): ...`, `fix(scope): ...`, `chore: ...`, `docs: ...`, `refactor: ...`, etc.).
- Keep commits focused; include related migrations or Swagger updates in the same commit when relevant.
- Before opening a PR, run `make linter` and `make test`, summarize behavioral changes, link relevant issues, and attach API diffs or screenshots when endpoints or docs change.
- When drafting commit messages, combine the diff of staged files with any relevant AI agent conversation tied to those changes; if no conversation applies, rely on the staged diff alone.

## Helper CLIs (under `cmd/kit`)
- `jwt`: Issue and inspect JWTs for local/dev.
- `http_route`: Generate or lint HTTP route scaffolding.
- `rdbms_seeder`: Seed relational databases with base data.
- `permission_seeder`: Seed Casbin roles, permissions, and grouping policies.

## Agent Notes & Precedence
- Scope: This file applies to the entire repository.
- Precedence: More deeply nested AGENTS.md files override this one for their subtree.
- Behavior: Be precise, safe, and helpful. Do not fix unrelated bugs. Ask for clarification when requirements are ambiguous.
- Formatting: Follow existing code style. Avoid intrusive refactors not required by the task.

## Specification-Driven Development with SpecKit

This project uses [SpecKit](https://github.com/github/spec-kit) for specification-driven development. SpecKit provides a structured workflow from feature specification to implementation.

### Constitutional Authority

**MANDATORY**: Before starting ANY development work, agents MUST read `.specify/memory/constitution.md`. This document contains the project's non-negotiable architectural principles and requirements that supersede all other guidance.

- Location: `.specify/memory/constitution.md`
- Authority: Constitutional principles override any conflicting guidance in other documentation
- Scope: All development activities (features, bug fixes, refactoring, testing)
- Compliance: All code, architecture decisions, and implementation plans MUST adhere to constitutional principles

### SpecKit Workflow Stages

SpecKit implements a five-stage specification-driven development workflow. Each stage is initiated through slash commands (custom commands integrated with your AI coding assistant).

**Core Philosophy**: Specifications don't serve code—code serves specifications. The specification is the source of truth, with implementation plans and code continuously regenerated from it.

#### Stage 1: Constitution (`/speckit.constitution`)
- **Purpose**: Establish foundational governance principles that guide all development decisions
- **When to use**: At project initialization or when amending architectural principles
- **Output**: `.specify/memory/constitution.md` containing non-negotiable principles covering:
  - Code quality standards
  - Testing requirements
  - Architectural constraints
  - User experience guidelines
  - Performance expectations
- **Key principle**: The constitution is referenced throughout all subsequent phases to ensure consistency

#### Stage 2: Specification (`/speckit.specify`)
- **Purpose**: Define **what** to build and **why**, without specifying **how**
- **When to use**: At the start of any new feature or significant change
- **Focus**: Be explicit about functional requirements and user needs; avoid tech stack details at this stage
- **Output**: `specs/[number]-[short-name]/spec.md` containing:
  - Detailed feature requirements
  - User stories
  - Acceptance criteria
  - Constraints and assumptions
- **Important**: Marks ambiguities with `[NEEDS CLARIFICATION]` to prevent unvalidated assumptions
- **Next steps**: `/speckit.clarify` (optional) or `/speckit.plan`

#### Stage 3: Planning (`/speckit.plan`)
- **Purpose**: Translate business requirements into technical architecture and technology choices
- **When to use**: After specification is complete and validated
- **Focus**: Make technology decisions (frameworks, databases, infrastructure) based on solidified requirements
- **Output**:
  - `plan.md` - High-level technical approach and architecture
  - `research.md` - Technology investigation and decisions
  - `data-model.md` - Entity schemas and relationships
  - `contracts/` - API specifications
  - `quickstart.md` - Validation scenarios
- **Constitutional check**: Ensures alignment with principles before proceeding
- **Next steps**: `/speckit.tasks`

#### Stage 4: Task Generation (`/speckit.tasks`)
- **Purpose**: Break down the implementation plan into actionable, dependency-ordered work items
- **When to use**: After planning is approved
- **Process**:
  - Analyzes contracts, data models, and test scenarios
  - Converts abstract requirements into concrete tasks
  - Marks independent tasks as parallelizable `[P]`
  - Organizes by user story and dependency order
- **Output**: `tasks.md` containing executable task breakdown ready for implementation
- **Task format**: `- [ ] [TaskID] [P?] [Story?] Description with file path`
- **Next steps**: `/speckit.analyze` (optional) or `/speckit.implement`

#### Stage 5: Implementation (`/speckit.implement`)
- **Purpose**: Execute the task plan systematically to build the feature
- **When to use**: After tasks are generated and optionally analyzed
- **Process**:
  - Executes tasks phase-by-phase following dependency order
  - Respects parallel execution opportunities marked with `[P]`
  - Maintains test-driven approach where applicable
  - Marks completed tasks with `[X]` in `tasks.md`
- **Constitutional compliance**: All implementation MUST adhere to constitutional principles
- **Output**: Working implementation of the specified feature

### Optional Quality Commands

These commands enhance quality but are not required in the core workflow:

#### `/speckit.clarify`
- **Purpose**: Structured questioning to address specification gaps before planning
- **When to use**: After `/speckit.specify` if ambiguities remain
- **Process**: Asks up to 5 targeted questions, updates spec incrementally
- **Output**: Updated `spec.md` with resolved ambiguities

#### `/speckit.analyze`
- **Purpose**: Cross-artifact consistency and coverage analysis
- **When to use**: After `/speckit.tasks`, before implementation
- **Process**: Non-destructive validation of spec/plan/tasks alignment
- **Output**: Analysis report identifying gaps, conflicts, or constitution violations

#### `/speckit.checklist`
- **Purpose**: Generate custom validation checklists for requirement quality
- **When to use**: At any stage to validate completeness, clarity, consistency
- **Key concept**: "Unit tests for requirements" - tests requirement quality, NOT implementation
- **Output**: Domain-specific checklists (e.g., `ux.md`, `security.md`, `api.md`)

### Integration with Development Workflow

SpecKit complements the existing workflow described in "Development Workflow" section:
- Use `/speckit.specify` and `/speckit.plan` BEFORE implementing new features or endpoints
- Task breakdown in `tasks.md` provides structured guidance for implementation
- Constitutional principles in `.specify/memory/constitution.md` enforce architecture constraints
- All SpecKit stages MUST respect the constitutional authority and "Business Logic Layer Architecture" patterns

### Constitution and SpecKit Templates

The constitution (`.specify/memory/constitution.md`) defines principles that inform SpecKit templates:
- `spec-template.md`: Feature specification structure
- `plan-template.md`: Implementation plan structure
- `tasks-template.md`: Task breakdown structure

When the constitution is amended, dependent templates are updated automatically to maintain consistency.

## Project Specifications (`docs/specs`)
The `docs/specs` directory is the designated location for project function specifications.

### Conventions
1. **Directory Structure**: Each specification must be contained within a directory following the `00X-NAME` pattern.
   - `00X`: A three-digit sequential number (e.g., `001`, `002`).
   - `NAME`: A concise description of the specification (e.g., `user-registration`).
   - Inside this directory, there must be at least two files: `spec.md` and `plan.md`.
2. **SpecKit Compliance**: The content and structure of `spec.md` and `plan.md` must strictly adhere to the guidelines defined in the [Specification-Driven Development with SpecKit](#specification-driven-development-with-speckit) section.
3. **Sequential Relationship**: The numbering represents a strict order or dependency sequence between documents. Ensure new specifications respect this existing order.
4. **Language**: All documentation within `docs/specs` must be written in **Traditional Chinese (zh-TW)**. However, software engineering terminology must use **English (en-US)**.
5. **Security**: No sensitive data (passwords, API keys, private keys, credentials, etc.) is allowed in `docs/specs`.
6. **Accuracy**: Do not write documentation that contradicts the actual code logic. Ensure the specification matches the implementation.
7. **Content Separation**:
   - `spec.md`: Must focus strictly on **Why** and **What**. It is forbidden to include **How** (implementation details).
   - `plan.md`: Must focus strictly on **How** (implementation details, technology stack, flow control). It is forbidden to include **Why** and **What**.
8. **Spec Structure**: When describing exposed interfaces (e.g., APIs), strict separation of concerns is required:
   - **In `spec.md` (Functional Requirements)**:
     - **Why**: The business or system reason for this feature.
     - **What**: The functional description of the feature.
     - **User Story**: In the format "Given [context], When [action], Then [outcome]".
   - **In `plan.md` (Technical Implementation)**:
     - **How**: The overall technical approach to realize the requirements.
     - **Technical Details**: Specifics on flow control, technology stack, and data models (including but not limited to these aspects).
     - **Implementation Strategy**: High-level phases or approach for delivery. **Do not** include detailed task lists or checkboxes (these belong in `tasks.md`).
9. **Verification**: Whenever modifying `spec.md` or `plan.md`, explicitly verify that technical implementation details (How) have not leaked into `spec.md`, and that functional requirements (Why/What) are not duplicated in `plan.md`.

## Conventional Commits

This project follows [Conventional Commits 1.0.0](https://www.conventionalcommits.org/en/v1.0.0/).

### Format

The commit message should be structured as follows:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Commit Types

- **feat:** New features (correlates with `MINOR` in semantic versioning)
- **fix:** Bug fixes (correlates with `PATCH` in semantic versioning)
- **chore:** Maintenance, dependencies, tooling (e.g., `chore(deps): update dependencies`)
- **docs:** Documentation changes (e.g., `docs: update README`)
- **refactor:** Code refactoring without feature changes or bug fixes
- **perf:** Performance improvements
- **test:** Adding or updating tests
- **ci:** CI/CD configuration changes
- **build:** Build system changes

### Breaking Changes

Indicate breaking changes by appending `!` before the colon or using a `BREAKING CHANGE:` footer:

```
feat(api)!: remove deprecated user endpoint

BREAKING CHANGE: /api/v1/user endpoint is no longer available; use /api/v2/user
```

Breaking changes correlate with `MAJOR` in semantic versioning.

For complete specification details, see the [official documentation](https://www.conventionalcommits.org/en/v1.0.0/).

## Commit Message Writing Guidelines

The recommended way to write a commit message is to **combine conversation history with code changes**. This allows the AI assistant to fully understand the **intent** and **implementation** of the changes. If no conversation history is available, use the alternate method.

### Preferred Approach: Combine Conversation and `git diff`

This is the most recommended and simplest approach. The AI assistant will automatically analyze the entire interaction and final code changes to generate the most accurate commit message.

**Workflow:**
1. After completing a task, stage your changes (`git add .`).
2. Ask the AI assistant to generate a commit message.

**Example Prompt:**
> "We're done. Please generate a commit message for me."

### Alternate Approach: Use `git diff` Only

If there is no relevant conversation history (e.g., you completed the work offline and now need to commit), the AI assistant can analyze the `git diff` output alone to generate a commit message.

**Workflow:**
1. Stage your changes (`git add .`).
2. Run `git diff --staged`.
3. Provide the complete `diff` output to the AI assistant.

**Example Prompt:**
> "Please generate a commit message for the following `git diff`:"
> ```diff
> diff --git a/src/utils/math.js b/src/utils/math.js
> index 6e9b2f7..8b4e6ad 100644
> --- a/src/utils/math.js
> +++ b/src/utils/math.js
> @@ -1,5 +1,9 @@
>  function add(a, b) {
>    return a + b;
>  }
> +
> +function subtract(a, b) {
> +  return a - b;
> +}
>
> -module.exports = { add };
> +module.exports = { add, subtract };
> ```

### Language and Grammar

- **Use English:** Commit messages must be written entirely in English to facilitate international collaboration and tool processing.
- **Grammatical Correctness:** Before committing, review your commit message to ensure correct grammar, proper spelling, and appropriate punctuation. A clear and professional message helps others understand your changes.
