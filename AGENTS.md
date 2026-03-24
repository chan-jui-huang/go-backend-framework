# AGENTS.md

This document guides AI agents and contributors working in this repository. It defines project structure, conventions, workflows, and safety rules. Its scope covers the entire repository. If a more deeply nested AGENTS.md exists, it overrides conflicting guidance for its subtree.

## Repository Structure & Assets
- `bin/`: Build artifacts emitted by the Makefile targets (created on demand).
- `cmd/app`: Main HTTP entrypoint and wiring; builds to `bin/app`.
- `cmd/kit`: Helper CLIs (`jwt`, `http_route`, `rdbms_seeder`, `permission_seeder`).
- `cmd/template`: Minimal `fx` bootstrap executable that wires registrar providers and invokes; useful for scaffolding.
- `internal`: Domain logic and supporting modules with colocated tests (`*_test.go`).
  - `deps/`: Typed accessors for config and service dependencies populated during `fx` bootstrap.
  - `http/`: Gin HTTP stack — `controller/`, `middleware/`, `response/`, `route/`, `server.go`.
  - `scheduler/`: Background jobs (cron-like tasks) and example jobs under `job/`.
  - `registrar/`: `fx` providers, lifecycle hooks, and dependency registration used during bootstrap.
  - `pkg/`: Project-specific packages (database, models, permission logic, user domain helpers, etc.).
  - `migration/`: Database migrations split into `rdbms/` and `clickhouse/` trees.
  - `test/`: `fx`-based test runtime plus fixtures under `runtime/`, `fixture/{db,domain,http,scenario}/`, and `fake/`.
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
- Tests & benchmarks: `make test [args=./...]`, `make benchmark`, or `go test ./...` (ensure SQLite and `.env.test` are available).
- When using `make test`, do not set the `GOCACHE` environment variable.

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
5. For admin/protected capabilities, synchronise seeds between runtime and tests: update `cmd/kit/permission_seeder/permission_seeder.go` and mirror the same permission preset in `internal/test/fake/permission.go`; adjust `internal/test/fixture/domain/permission.go` or `internal/test/fixture/scenario/admin_api.go` when the admin fixture flow changes.
6. Whenever you introduce a new mock service (gomock or hand-written), also register an empty instance in `internal/test/runtime/runtime.go:emptyMockedServices()` so the test harness has placeholders for injected services.

General guidance:
- Keep changes minimal and localized; avoid unrelated refactors.
- Favor composition over duplication; reuse helpers under `internal/pkg/` and existing shared utilities.
- Update Swagger comments whenever API shapes change, then run `make swagger`.
- Wrap work in a transaction only when multiple insert/update/delete statements need to succeed together.

## Testing Guidelines
- Test framework: standard `testing`; `testify` is available for assertions.
- Location: place `*_test.go` files alongside the code they cover; prefer table-driven tests or testify suites.
- Bootstrapping: leverage `internal/test/runtime` and `internal/test/fixture/*` helpers for environment loading, migrations, HTTP handlers, seeded users, and scenario setup instead of reimplementing test wiring.
- Mock services: when introducing a new mock service, also register it in `internal/test/runtime/runtime.go` within `emptyMockedServices` so the test harness initializes it.
- Migrations: run required migrations first (e.g., `make sqlite-migration args=up`). Tests expect `.env.test` to provide DSNs and secrets.
- Coverage: exercise both success and failure paths. **REQUIRED**: Test all possible return values from controllers (success, validation failures, auth failures, business errors). Reference existing patterns in `internal/http/controller/**/*_test.go`. Document skipped integration tests in PR descriptions. Add controller/service tests whenever adding new endpoints.

## Migrations
- Use Make targets with environment loaded from `.env`: `make mysql-migration args="up"`, `make pgsql-migration args="up"`, `make sqlite-migration args="up"`, or `make clickhouse-migration args="up"`.
- Ensure relevant `DB_*` variables are set for the target database.
- Keep migrations idempotent and reversible; document non-trivial data movements.
- Migration filenames MUST use a timestamp reflecting the actual creation time; generate a fresh `YYYYMMDDhhmmss` prefix when the file is added, and ensure each migration file has a unique timestamp prefix (do not reuse the same `YYYYMMDDhhmmss` across multiple migration files). When creating paired migrations (e.g., prod/test), keep the timestamp identical across the pair while still unique per table/change.

## Security & Configuration
- Never commit secrets. Place credentials in `storage/...` as described in the README.
- Provide `.env.test` and `.env.staging` locally when needed; avoid committing real values.
- Validate configuration via `make run` and integration tests before deployment.

## Commit & PR Guidelines
- Use `.github/skills/commit-message-writer` for commit message drafting guidance.
- Before opening a PR, run `make linter` and `make test`, summarize behavioral changes, link relevant issues, and attach API diffs or screenshots when endpoints or docs change.

## Helper CLIs (under `cmd/kit`)
- `jwt`: Issue and inspect JWTs for environment-specific usage.
- `http_route`: Generate or lint HTTP route scaffolding.
- `rdbms_seeder`: Seed relational databases with base data.
- `permission_seeder`: Seed Casbin roles, permissions, and grouping policies.

## Agent Notes & Precedence
- Scope: This file applies to the entire repository.
- Precedence: More deeply nested AGENTS.md files override this one for their subtree.
- Behavior: Be precise, safe, and helpful. Do not fix unrelated bugs. Ask for clarification when requirements are ambiguous.
- Formatting: Follow existing code style. Avoid intrusive refactors not required by the task.

## GitHub Copilot Instructions Compatibility

To mirror GitHub Copilot coding agent behavior for repository custom instructions, agents working in this repository MUST also load applicable instruction files under `.github/`.

### Supported Instruction Sources
- Repository-wide instructions: `.github/copilot-instructions.md`
- Path-specific instructions: `.github/instructions/**/*.instructions.md`
- Agent instructions: this `AGENTS.md` file and any more deeply nested `AGENTS.md` files

### When To Load `.github/instructions`
- Before planning, editing, reviewing, or generating code for specific files, determine the in-scope target files first.
- Load every `.github/instructions/**/*.instructions.md` file whose front matter `applyTo` glob matches at least one in-scope file.
- Re-evaluate applicable path-specific instructions whenever the scope changes to additional files or directories.
- If the task is repository-wide and no concrete file is known yet, first inspect likely target files, then load matching path-specific instructions before making code changes.

### Path Matching Rules
- Path-specific instruction files MUST end with `.instructions.md`.
- The file MUST begin with front matter containing `applyTo`.
- `applyTo` uses glob syntax and may contain multiple comma-separated patterns such as `**/*.go,**/go.mod,**/go.sum`.
- Match `applyTo` patterns against repository-relative paths.
- If multiple path-specific instruction files match, combine all of them.

### `excludeAgent` Rules
- If a path-specific instruction file contains `excludeAgent: "coding-agent"`, agents MUST ignore that file.
- If `excludeAgent` is omitted, treat the file as applicable to coding agents.
- `excludeAgent: "code-review"` does not exclude coding agents and should still be applied for implementation tasks.

### Effective Precedence
- For repository custom instruction compatibility, use this Copilot-aligned order inside the repository:
  1. Matching path-specific instructions from `.github/instructions/**/*.instructions.md`
  2. Repository-wide instructions from `.github/copilot-instructions.md`
  3. Agent instructions from `AGENTS.md`
- When instructions conflict, prefer the higher item in the list above.
- Direct user instructions still override repository instruction files.
- A more deeply nested `AGENTS.md` still overrides this file for its subtree, but it does not cancel applicable `.github/instructions` matches unless it states so explicitly.

### Operational Notes
- Treat `.github` instruction files as additive guidance, similar to how Copilot combines multiple relevant instruction sources.
- Keep `.github` instruction files short, concrete, and broadly reusable for the paths they target.
- Do not assume `.github/copilot-instructions.md` exists; load it only if present.
- For reviews, analysis, and code generation involving Go files in this repository, the existing `.github/instructions/go.instructions.md` file applies to `**/*.go`, `**/go.mod`, and `**/go.sum` unless a future change narrows or excludes it.
