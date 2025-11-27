<!--
Sync Impact Report
==================
Version Change: 1.0.1 → 2.0.0
Rationale: MAJOR - Remove all path dependencies, make constitution fully portable

Modified Principles:
- Principle I: Removed ALL path references (`internal/pkg/` → "business logic layer")
- Principle III: Removed ALL path references (`internal/test/`, `test.HttpHandler` → abstract terms)
- Principle V: Removed ALL path references, uses pure architectural layer terminology

Added Sections:
- Technical Requirements (replaces Technology Stack with abstract requirements)
- Layered Architecture (replaces any directory-specific constraints)
- Layer dependency rules (prohibited/allowed relationships)

Removed Sections:
- Directory Structure Constraints (BREAKING - path mappings should be maintained in CLAUDE.md)

Modified Sections:
- Security Requirements: Removed path reference (`storage/` → "environment-specific configuration files")
- Technology Stack → Technical Requirements: Replaced concrete technology choices (Go, Gin, GORM) with abstract requirements (quantum-resistant auth, GPU-resistant hashing, ACID transactions)
- Compliance Review: Removed tool-specific reference ("pull requests" → "proposed changes")
- Guidance Integration → Implementation Guidance: Removed project-specific references (CLAUDE.md, AI agent) with abstract documentation structure

Templates Requiring Updates:
✅ All templates remain compatible (already use abstract terminology)

Migration Guide for Projects:
- Projects using this constitution should maintain their own path mappings in project-specific docs (e.g., CLAUDE.md)
- Projects should maintain their own technology choices mapping abstract requirements to concrete implementations
- For go-backend-framework: Path mappings and technology choices already exist in CLAUDE.md

Previous Changes:
- v1.0.1 (2025-11-27): Clarification patch - Expanded testing scenarios
- v1.0.0 (2025-11-27): Initial constitution ratification
-->

# go-backend-framework Constitution

## Core Principles

### I. Single Responsibility Principle (SOLID Foundation)

Every module, class, or function MUST have one and only one reason to change. This is the PRIMARY architectural constraint.

**Non-negotiable rules:**
- Each module in the business logic layer MUST serve a single domain or capability
- Interface controllers MUST NOT contain business logic; they MUST delegate to the business logic layer
- Business logic MUST reside in a dedicated layer, separate from all interface implementations
- Data models MUST only define schema and relationships, never business operations
- Cross-cutting concerns (authentication, authorization, logging) MUST be handled by dedicated, single-purpose components

**Rationale:** Single responsibility enables independent testing, parallel development, and confident refactoring. Violating SRP creates cascading changes, hidden dependencies, and fragile tests.

### II. Keep It Simple, Stupid (KISS)

Start with the simplest solution that solves the current problem. Avoid premature abstraction and over-engineering.

**Non-negotiable rules:**
- Three similar lines of code is BETTER than a premature abstraction
- Do NOT create helpers, utilities, or abstractions for one-time operations
- Do NOT design for hypothetical future requirements (YAGNI principle applies)
- Do NOT add configuration flags or feature toggles when direct code changes suffice
- Only extract shared code when it appears in THREE or more places with identical behavior

**Rationale:** Complexity is the enemy of maintainability. Simple code is easier to understand, debug, and change. Premature abstraction creates indirection without proven value.

### III. Integration-First Testing (Test-Driven Development)

Tests MUST be written BEFORE implementation (Red-Green-Refactor). Integration tests take priority over isolated unit tests.

**Non-negotiable rules:**
- Write contract/integration tests FIRST before any implementation
- Tests MUST fail initially (Red phase) to prove they test real behavior
- Tests MUST use realistic environments:
  - Real databases (in-memory for tests, production-grade for staging/production)
  - Actual service instances via test fixture infrastructure
  - Real interface handlers through test harness framework
- Mock/stub ONLY external third-party services (payment providers, external APIs)
- Every controller MUST test ALL possible return scenarios including but not limited to:
  - **Success cases**: 200 OK, 201 Created, 204 No Content (with expected data/empty body)
  - **Client errors (4xx)**:
    - 400 Bad Request: Invalid input format, malformed JSON, type mismatches
    - 401 Unauthorized: Missing token, expired token, invalid token signature
    - 403 Forbidden: Valid authentication but insufficient permissions, role-based access denial
    - 404 Not Found: Resource does not exist, invalid resource ID
    - 409 Conflict: Duplicate entries (e.g., email already exists), concurrent modification conflicts
    - 422 Unprocessable Entity: Business rule violations (e.g., insufficient balance, invalid state transitions)
    - 429 Too Many Requests: Rate limit exceeded
  - **Server errors (5xx)**:
    - 500 Internal Server Error: Unexpected exceptions, database connection failures
    - 503 Service Unavailable: External service timeouts, circuit breaker open
  - **Edge cases and boundary conditions**:
    - Empty request bodies, null values, missing required fields
    - Maximum/minimum value boundaries (e.g., string length limits, numeric ranges)
    - Special characters and Unicode in text fields
    - Large payloads approaching size limits
    - Concurrent requests with race conditions
    - Database constraint violations (foreign key, unique, not null)
- Tests MUST be colocated with implementation following language conventions

**Rationale:** Integration tests validate real-world behavior and catch configuration errors, connection issues, and integration bugs that unit tests miss. Writing tests first forces clear requirements and prevents untested code.

### IV. Semantic Versioning with Conventional Commits

All changes MUST follow Git-based version control with semantic versioning and Conventional Commits 1.0.0 specification.

**Non-negotiable rules:**
- Commit messages MUST follow format: `<type>[optional scope]: <description>`
- Commit types: `feat`, `fix`, `chore`, `docs`, `refactor`, `perf`, `test`, `ci`, `build`
- Breaking changes MUST be marked with `!` or `BREAKING CHANGE:` footer
- Version numbers MUST follow MAJOR.MINOR.PATCH:
  - MAJOR: Breaking changes (incompatible API changes, removed features)
  - MINOR: New features (backward-compatible functionality additions)
  - PATCH: Bug fixes (backward-compatible corrections)
- All commits MUST be atomic (single logical change)
- Commit messages MUST be in English with correct grammar

**Rationale:** Conventional commits enable automated changelog generation, semantic versioning enforcement, and clear change history. Semantic versioning communicates impact to consumers.

### V. Separation of Concerns (Layered Architecture)

The framework MUST maintain strict separation between business logic and interface layers.

**Non-negotiable rules:**
- Business logic MUST be interface-agnostic and isolated in a dedicated layer
- Interface layers (HTTP, CLI, scheduler, etc.) MUST only handle:
  - Input validation
  - Calling business logic layer functions
  - Output formatting
- All interface implementations MUST delegate to the same business logic layer
- No business logic duplication across interface layers
- Database transactions MUST be managed in the business logic layer, not interface controllers

**Rationale:** Separation enables business logic reuse across multiple interfaces, independent testing of business rules, and flexibility to add new interfaces (gRPC, WebSocket) without duplicating logic.

## Architecture Constraints

### Technical Requirements

**Security & Cryptography:**
- Authentication MUST use quantum-resistant algorithms
- Password hashing MUST use GPU-resistant algorithms (e.g., Argon2, scrypt, bcrypt with high cost)
- Cryptographic keys MUST use industry-standard key lengths (minimum 256-bit for symmetric, 2048-bit for RSA)

**Data Persistence:**
- MUST support ACID transactions for business-critical operations
- MUST use connection pooling for database access
- MUST provide schema versioning and migration capabilities
- MUST support both relational and time-series data storage patterns when needed

**Code Quality:**
- All code MUST pass static analysis security checks
- All code MUST follow language-standard formatting conventions
- All production builds MUST be optimized for performance

**Observability:**
- MUST provide structured logging with log rotation
- MUST support configuration via environment variables and configuration files
- MUST generate machine-readable API documentation

**Authorization:**
- MUST implement attribute-based or role-based access control
- MUST support policy-based authorization rules

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
- Each interface layer delegates to the business logic layer

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

### Security Requirements

- NEVER commit secrets (use environment-specific configuration files)
- NEVER update git config via automation
- NEVER run destructive git commands (`push --force`, `hard reset`) without explicit user request
- NEVER skip git hooks (`--no-verify`, `--no-gpg-sign`) unless explicitly requested
- ALL user input MUST be validated before processing
- Passwords MUST use Argon2id hashing (never bcrypt or plain SHA)
- JWT tokens MUST use Ed25519 signatures
- Error messages MUST NOT leak sensitive information (stack traces, SQL queries)
- Database credentials MUST come from environment variables, never hardcoded

### Commit Guidelines

**Required format:**
```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**Best practices:**
- Combine conversation context with `git diff --staged` when drafting messages
- Keep subject line under 72 characters
- Use imperative mood ("add feature" not "added feature")
- Reference issues in footer: `Fixes #123`, `Relates to #456`
- Include Co-Authored-By when appropriate

**Example:**
```
feat(user): add password reset endpoint

Implements user story US-042 for self-service password reset.
Adds new POST /api/user/reset-password endpoint with email validation.

BREAKING CHANGE: /api/user/password endpoint removed; use /api/user/reset-password

Fixes #123
```

## Governance

### Constitution Authority

This constitution supersedes all other development practices and documentation. When conflicts arise, this document takes precedence.

### Amendment Process

1. Proposed changes MUST be documented with rationale
2. Impact analysis MUST identify affected templates and code
3. Version number MUST increment according to semantic versioning:
   - MAJOR: Backward-incompatible principle removals or redefinitions
   - MINOR: New principles or materially expanded guidance
   - PATCH: Clarifications, wording fixes, non-semantic refinements
4. All dependent templates (`plan-template.md`, `spec-template.md`, `tasks-template.md`) MUST be updated
5. Sync Impact Report MUST be prepended to constitution as HTML comment
6. `LAST_AMENDED_DATE` MUST be updated to current date

### Compliance Review

- All proposed changes MUST be reviewed for compliance with these principles before integration
- Code reviewers MUST challenge complexity and demand justification
- Violations require documented rationale in implementation plans
- Architecture decisions MUST reference applicable constitutional principles

### Implementation Guidance

Projects SHOULD maintain separate operational documentation that complements this constitution:
- **Constitution**: Defines non-negotiable architectural principles and requirements (this document)
- **Implementation Guide**: Provides project-specific operational instructions, concrete technology choices, directory mappings, and development workflows

The constitution remains stable and principle-focused, while implementation guides evolve with project-specific decisions and tooling.

**Version**: 2.0.0 | **Ratified**: 2025-11-27 | **Last Amended**: 2025-11-27
