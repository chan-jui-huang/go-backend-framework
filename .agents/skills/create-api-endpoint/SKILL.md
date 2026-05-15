---
name: create-api-endpoint
description: Build or modify API endpoints in this go-backend-framework repository. Use when adding controller files, request/response DTOs, Swagger annotations, canonical error messages, logging, internal/pkg business logic, routes, permissions, seed data, and controller tests for a new or changed API.
---

# Create API Endpoint

## Goal

Create repository-standard API endpoints with thin HTTP controllers, reusable business logic, synchronized authorization presets, accurate Swagger contracts, and tests for every reachable response path.

## Before Editing

1. Read the root `AGENTS.md`, `.github/copilot-instructions.md` if present, and every `.github/instructions/**/*.instructions.md` file whose `applyTo` matches the files you will touch.
2. Inspect nearby files before choosing names or patterns:
   - Controllers: `internal/http/controller/<area>/**`
   - Routes: `internal/http/route/**`
   - Business logic: `internal/pkg/<domain>/**`
   - Error responses: `internal/http/response/error_message.go`
   - Permissions: `cmd/kit/permission_seeder/permission_seeder.go` and `internal/test/fake/permission.go`
   - Tests: matching `internal/http/controller/**/*_test.go`
3. Prefer extending the closest existing area. Create a new area only when the API does not fit an existing controller/route/domain package.
4. Keep paths in kebab case, Go package names lowercase without underscores, struct names singular PascalCase, and `Id` casing for identifiers.

## Endpoint Workflow

1. Define the API contract:
   - HTTP method, path, auth requirements, request body/query/path params, response data, and all reachable errors.
   - Decide whether the route is public, authenticated only, or protected by admin Casbin authorization.
   - Use existing canonical errors where possible before adding new messages.

2. Implement or extend `internal/pkg/<domain>` first when the behavior is reusable or shared:
   - Keep functions framework-agnostic; do not import Gin or HTTP response packages.
   - Accept `*gorm.DB` when database work is needed, matching existing helpers.
   - Match existing query shape in that package: preload choices, ordering, pagination, and error wrapping.
   - Wrap multiple writes in a transaction in the business layer when they must succeed together.
   - Add colocated unit tests where the logic has meaningful branches.

3. Create or update controller files under `internal/http/controller/<area>/`:
   - Use `<feature>_<action>.go`, for example `user_register.go` or `permission_create.go`.
   - Put request structs near the handler that binds them unless the local package already centralizes DTOs.
   - Put response data structs in `shared.go` when reused; otherwise keep them near the handler.
   - Use `binding` tags for Gin input validation and `json`, `mapstructure`, `validate`, and `format` tags on response DTOs following nearby examples.
   - Controllers should validate input, call `internal/pkg`, map errors to response messages, log, and format output.

4. Handle requests and responses consistently:
   - Body input: `c.ShouldBindBodyWithJSON(reqBody)` so `response.ErrorResponse.MakeLogFields(c)` can still include the request body on failures.
   - Query input: `c.ShouldBindQuery(reqQuery)`.
   - Path params: bind or parse explicitly, and return `response.RequestValidationFailed` or `response.BadRequest` according to nearby controller behavior.
   - Validation failure pattern:

```go
errResp := response.NewErrorResponse(response.RequestValidationFailed, errors.WithStack(err), response.MakeValidationErrorContext(err))
logger.Warn(errResp.Message, errResp.MakeLogFields(c)...)
c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
return
```

- Success pattern:

```go
c.JSON(http.StatusOK, response.NewResponse(data))
```

5. Log every error response:
   - Start handlers with `logger := deps.Logger()` when the handler can return errors.
   - Use `logger.Warn(errResp.Message, errResp.MakeLogFields(c)...)` for expected request/business failures.
   - Use `logger.Error(...)` only for truly unexpected server-side failures if nearby code distinguishes them.
   - Do not log sensitive values or expose stack traces through user-facing messages.

6. Add canonical error messages only when needed:
   - Define message constants exclusively in `internal/http/response/error_message.go`.
   - Add a unique `<status>-<sequence>` entry to `MessageToCode`.
   - Keep the map grouped by status and update `internal/http/response/error_message_test.go`.
   - Swagger must list only the messages and status codes the handler can actually return.

7. Add Swagger annotations above the handler:
   - Include `// @tags`, `// @accept json` when JSON input is accepted, and `// @produce json`.
   - Use placeholder descriptions as `" "` in every `// @param`.
   - Include `X-XSRF-TOKEN` for state-changing routes that require CSRF.
   - Include `Authorization` for authenticated/protected routes.
   - Use actual DTO names in request and response annotations.
   - Keep route annotations aligned with runtime path and method, for example `// @router /api/admin/user-role [put]`.
   - After changing annotations, run `make swagger` unless the user explicitly asks to skip it.

8. Wire routes:
   - For an existing router, update `internal/http/route/<area>/api_route.go`.
   - For a new router, create it so it implements `route.Router`, then append it in `internal/http/route/api_route.go`.
   - Public routes usually live under a specific router group such as `api/user`.
   - Admin/protected routes should use `middleware.Authenticate()` and `middleware.Authorize()` consistently with `internal/http/route/admin/api_route.go`.

9. Add permissions for protected admin capabilities:
   - Add permission names and `p` Casbin rules to `cmd/kit/permission_seeder/permission_seeder.go`.
   - Mirror the same permissions and rules in `internal/test/fake/permission.go`.
   - Match route paths exactly as Casbin sees them, including `:id` path params and uppercase HTTP methods.
   - If the admin fixture flow changes, update `internal/test/fixture/domain/permission.go` or `internal/test/fixture/scenario/admin_api.go`.
   - If a new mock service is introduced, register an empty instance in `internal/test/runtime/runtime.go:emptyMockedServices()`.

10. Write controller tests alongside the handler:

- Use the high-level `internal/test` runtime and fixture abstractions instead of rebuilding app wiring by hand.
- Inspect `internal/test/test.go` and `internal/test/runtime/runtime.go`, then choose the smallest available runtime constructor or `RuntimeOptions` setup that matches the endpoint dependencies such as HTTP routing, authentication, RDBMS migrations, ClickHouse migrations, Redis, or mocked services.
- Reuse existing runtime fields and fixture/scenario helpers from `internal/test/fixture/**` before adding new test wiring.
- Cover success, validation failure, auth failure, authorization failure, business errors, and server errors when reachable.
- Public mutating routes should include CSRF mismatch tests when middleware applies.
- Protected routes should test missing/wrong token returns `Unauthorized` and insufficient permission returns `Forbidden`.
- Decode `response.Response` data into the controller response DTO and assert important fields.
- Decode failures into `response.ErrorResponse` and assert HTTP status, `Message`, `Code`, and validation `Context` when present.

11. Format and validate:

- Run `gofmt`/`goimports` on changed Go files.
- Run targeted `go test` for touched packages; broaden to `make test args=./...` or `go test ./...` when risk is wider.
- Run `make swagger` after annotation changes.
- Run `make linter` when the change touches shared behavior or before final handoff if time allows.

## File Checklist

Use this as a scoped checklist; not every endpoint needs every file.

- `internal/pkg/<domain>/*.go`: business logic and data access helpers.
- `internal/pkg/<domain>/*_test.go`: unit tests for reusable logic.
- `internal/http/controller/<area>/<feature>_<action>.go`: handler, request DTO, local response DTO if not shared.
- `internal/http/controller/<area>/shared.go`: reusable response data and `Fill` methods.
- `internal/http/controller/<area>/<feature>_<action>_test.go`: controller/integration tests.
- `internal/http/route/<area>/api_route.go`: route registration.
- `internal/http/route/api_route.go`: only when adding a new router.
- `internal/http/response/error_message.go` and `error_message_test.go`: only when adding canonical errors.
- `cmd/kit/permission_seeder/permission_seeder.go`: runtime admin permission seeds.
- `internal/test/fake/permission.go`: mirrored test permission preset.
- `internal/test/fixture/**`: only when fixture/scenario setup must change.
- `docs/*`: regenerated Swagger output after annotation changes.

## Example Prompts

```text
$create-api-endpoint Add an admin GET /api/admin/audit-log endpoint with pagination, permission audit-log-read, Swagger docs, and controller tests.
```

```text
$create-api-endpoint Add POST /api/user/profile that validates JSON input, delegates profile updates to internal/pkg/user, returns a ProfileData response, and covers validation/auth/business errors in tests.
```
