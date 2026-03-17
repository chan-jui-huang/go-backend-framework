---
description: 'Project-specific Go architecture and layering constraints extracted from AGENTS.md'
applyTo: '**/*.go'
---

# Go Backend Architecture Instructions

### Dependency Wiring

- **Config Snapshot**: YAML configuration with environment variable expansion (`${VAR_NAME}`) is loaded through registrar providers and stored in `internal/deps`.
- **Service Snapshot**: Initialized services are built by `fx` providers and exposed to existing code through typed `internal/deps` accessors.

### Bootstrap Sequence

```
load env/autoload → build booter config → fx.Supply + fx.Provide
                  → fx.Invoke(RegisterConfigDependencies, RegisterServiceDependencies)
                  → fx lifecycle hooks start validator, scheduler, and HTTP server
```

### Application Lifecycle Hooks (`fx`)

| Hook        | Registration               | Purpose                                                                  |
| ----------- | -------------------------- | ------------------------------------------------------------------------ |
| Validator   | `fx.OnStart`               | Register request-validation tag name handling before requests are served |
| Scheduler   | `fx.OnStart` / `fx.OnStop` | Start and stop background jobs                                           |
| HTTP Server | `fx.OnStart` / `fx.OnStop` | Run the Gin server and perform graceful shutdown                         |

See `cmd/app/main.go` for the canonical `fx` wiring example.

### Key Architectural Patterns

| Pattern                      | Purpose                                                                                     |
| ---------------------------- | ------------------------------------------------------------------------------------------- |
| **Fx Provider Pattern**      | Modular initialization is composed through registrar constructors and invokes               |
| **Dependency Injection**     | `fx` builds dependencies, and `internal/deps` exposes typed accessors for runtime/test code |
| **Configuration Management** | YAML + `.env` file support with variable expansion                                          |
| **Database Layer**           | Goose migrations + GORM ORM wrapper for type-safe operations                                |
| **Authentication**           | Ed25519 JWT (quantum-resistant) generated via `./bin/jwt`                                   |
| **Authorization**            | Casbin RBAC policies seeded via `cmd/kit/permission_seeder`                                 |
| **Logging**                  | Zap structured logging with rotation to `storage/log/`                                      |

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
