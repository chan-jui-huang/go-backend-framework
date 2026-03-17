---
applyTo: '**/*.go'
---

## Core Principles

### Single Responsibility Principle (SOLID Foundation)

Every module, class, or function MUST have one and only one reason to change. This is the PRIMARY architectural constraint.

**Non-negotiable rules:**

- Each module in the business logic layer MUST serve a single domain or capability
- Primary Rule (Single-Interface Scope/KISS Override): If a data access or business flow
  is used by a single interface only, it SHOULD remain in that interface layer
  (controller/scheduler/CLI). Do NOT extract to a shared business-logic layer
  unless reuse reaches three or more call sites or an existing helper already
  matches.
- Default Rule (Separation of Concerns): Controllers SHOULD delegate business
  logic to the business layer to keep interface code thin, except when the
  Single-Interface Scope/KISS Override applies.
- Business logic MUST reside in a dedicated layer, separate from all interface implementations
- Data models MUST only define schema and relationships, never business operations
- Cross-cutting concerns (authentication, authorization, logging) MUST be handled by dedicated, single-purpose components

**Rationale:** Single responsibility enables independent testing, parallel development, and confident refactoring. Violating SRP creates cascading changes, hidden dependencies, and fragile tests.

### Keep It Simple, Stupid (KISS)

Start with the simplest solution that solves the current problem. Avoid premature abstraction and over-engineering.

**Non-negotiable rules:**
- Three similar lines of code is BETTER than a premature abstraction
- Do NOT create helpers, utilities, or abstractions for one-time operations
- Do NOT design for hypothetical future requirements (YAGNI principle applies)
- Do NOT add configuration flags or feature toggles when direct code changes suffice
- Only extract shared code when it appears in THREE or more places with identical behavior

**Rationale:** Complexity is the enemy of maintainability. Simple code is easier to understand, debug, and change. Premature abstraction creates indirection without proven value.

### Separation of Concerns (Layered Architecture)

The framework MUST maintain strict separation between business logic and interface layers.

**Non-negotiable rules:**

- Except when the Single-Interface Scope/KISS Override applies, business logic
  MUST be interface-agnostic and isolated in a dedicated layer.
- Interface layers (HTTP, CLI, scheduler, etc.) SHOULD only handle:
  - Input validation
  - Calling business logic layer functions
  - Output formatting
  unless the Single-Interface Scope/KISS Override applies to a single-interface flow.
- All interface implementations SHOULD delegate to the same business logic
  layer when shared business logic exists.
- No business logic duplication across interface layers unless required by the
  Single-Interface Scope/KISS Override for single-interface use cases.
- Database transactions SHOULD be managed in the business logic layer, except
  when the Single-Interface Scope/KISS Override requires them in the interface layer
  for single-interface flows.

**Rationale:** Separation enables business logic reuse across multiple interfaces, independent testing of business rules, and flexibility to add new interfaces (gRPC, WebSocket) without duplicating logic.

### Layered Architecture

The system MUST be organized into distinct architectural layers with clear responsibilities:

**Business Logic Layer:**
- Contains all domain logic, business rules, and data operations
- Interface-agnostic: no dependencies on HTTP, CLI, or other interface concerns
- Manages database transactions and data integrity
- Organized by business domains/capabilities

**Interface Layers:**
- HTTP/REST API: Request handling, routing, response formatting
- CLI: Command-line argument parsing and output formatting
- Scheduler: Background job triggering and scheduling logic
- Each interface layer delegates to the business logic layer by default, except
  where the Single-Interface Scope/KISS Override applies

**Infrastructure Layer:**
- Service registration and dependency injection
- Configuration management
- Logging, monitoring, and observability
- Database migrations and schema management

**Test Infrastructure:**
- Test fixtures and utilities
- Test harness for integration testing
- Colocated with implementation following language conventions

**Layer Dependencies:**
- Interface layers → Business logic layer (allowed)
- Business logic layer → Infrastructure layer (allowed)
- Interface layers → Interface layers (prohibited)
- Business logic layer → Interface layers (prohibited)

## Business Logic Layer Architecture

### Separation of Concerns

This framework enforces a strict separation between business logic and presentation/interface layers:

- **`internal/pkg/`**: Contains all reusable business logic that can be invoked by any interface layer (HTTP APIs, schedulers, CLI commands, etc.). Functions here are pure business operations that handle data manipulation, database transactions, and domain-specific logic.

- **`internal/http/`**: Contains HTTP-specific handlers, middleware, routing, and request/response structures. Controllers in this layer orchestrate HTTP interactions but delegate actual business operations to `internal/pkg/`.

- **Other Interface Layers**: Schedulers (`internal/scheduler/`), CLI commands (`cmd/kit/`), and any future interfaces (gRPC, WebSocket, etc.) should similarly call functions from `internal/pkg/` rather than duplicating logic.

### API Development Workflow

When developing a new API endpoint, follow this mandatory pattern:

#### Business Logic First

- Check if the required business logic exists in `internal/pkg/<domain>/`
- If not, create or extend functions in `internal/pkg/<domain>/` first
- Example: For user operations, add functions to `internal/pkg/user/user.go`
- Keep these functions framework-agnostic and reusable

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
