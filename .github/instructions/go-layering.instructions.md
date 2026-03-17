---
applyTo: '**/*.go'
---

## Core Principles

### Single Responsibility Principle (SOLID Foundation)

- Every module, class, or function MUST have one and only one reason to change.
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

### Keep It Simple, Stupid (KISS)

- Three similar lines of code is BETTER than a premature abstraction
- Do NOT create helpers, utilities, or abstractions for one-time operations
- Do NOT design for hypothetical future requirements (YAGNI principle applies)
- Do NOT add configuration flags or feature toggles when direct code changes suffice
- Only extract shared code when it appears in THREE or more places with identical behavior

### Separation of Concerns (Layered Architecture)

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

## Business Logic Layer Architecture

### Separation of Concerns

- Reusable business logic belongs in `internal/pkg/`
- `internal/http/` is limited to HTTP-specific request handling, middleware, routing, and response formatting
- Other interface layers such as schedulers and CLI commands SHOULD reuse business logic instead of duplicating it

### API Development Workflow

#### Controller Implementation

- Create or extend handler files under `internal/http/controller/<area>/`
- Controllers should ONLY:
  - Validate HTTP requests
  - Call business logic from `internal/pkg/`
  - Format HTTP responses
- Never put business logic directly in controllers

#### Business Logic First

- Check if the required business logic exists in `internal/pkg/<domain>/`
- If not, create or extend functions in `internal/pkg/<domain>/` first
- Example: For user operations, add functions to `internal/pkg/user/user.go`
- Keep these functions framework-agnostic and reusable
