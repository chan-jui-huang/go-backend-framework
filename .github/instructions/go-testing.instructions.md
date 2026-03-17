---
applyTo: '**/*_test.go'
---

## Core Principles

### Integration-First Testing (Test-Driven Development)

Tests MUST be written BEFORE implementation (Red-Green-Refactor). Integration tests take priority over isolated unit tests.

#### Test-Driven Development Process

**Non-negotiable rules:**

- Write contract/integration tests FIRST before any implementation
- Tests MUST fail initially (Red phase) to prove they test real behavior
- Tests MUST be colocated with implementation following language conventions

#### Integration Testing Requirements

**Testing Environment:**
Tests MUST use realistic environments:
- Real databases (in-memory for tests, production-grade for staging/production)
- Actual service instances via test fixture infrastructure
- Real interface handlers through test harness framework

**Critical Flow Coverage:**
For every critical user or business flow, there MUST be at least one integration or end-to-end test that exercises the flow without mocking core infrastructure (database, message bus, HTTP stack) to validate real behavior across boundaries.

#### Test Doubles: When and How to Mock

**Scope of Mocking:**
Mock/stub MUST be limited to external systems and side-effect boundaries:
- Network calls (third-party APIs, internal microservices, payment providers)
- Messaging systems and background job queues
- File systems, system clock, randomness sources, OS-level services
- Datastores accessed via well-defined repository or gateway interfaces

**Prohibited Mocking:**
Internal domain logic (entities, value objects, pure domain services) MUST NOT be mocked.

#### Mock Verification Standards

When using test doubles:

**1. Choosing the Right Test Double:**
- Use **stubs** to provide predetermined responses without verifying interactions
- Use **mocks** only when interactions (who calls whom, with what arguments) are part of the observable contract / behavior specification

**2. What to Verify:**
Tests SHOULD verify arguments of mocked operations that are part of business rules or observable behavior, not incidental implementation details.

**3. Invocation Sequence:**
Tests SHOULD verify invocation sequence ONLY when the order itself is a business or protocol requirement.

**4. Comprehensive Coverage:**
Mock verification MUST cover all interactions that define the behavior under test, but MUST NOT assert on internal implementation details (such as private helper calls, logging, or framework internals) that are free to change.

#### Controller Testing Requirements

Every controller MUST test ALL possible return scenarios including but not limited to:

**Success Cases:**
- 200 OK, 201 Created, 204 No Content (with expected data/empty body)

**Client Errors (4xx):**
- **400 Bad Request**: Invalid input format, malformed JSON, type mismatches
- **401 Unauthorized**: Missing token, expired token, invalid token signature
- **403 Forbidden**: Valid authentication but insufficient permissions, role-based access denial
- **404 Not Found**: Resource does not exist, invalid resource ID
- **409 Conflict**: Duplicate entries (e.g., email already exists), concurrent modification conflicts
- **422 Unprocessable Entity**: Business rule violations (e.g., insufficient balance, invalid state transitions)
- **429 Too Many Requests**: Rate limit exceeded

**Server Errors (5xx):**
- **500 Internal Server Error**: Unexpected exceptions, database connection failures
- **503 Service Unavailable**: External service timeouts, circuit breaker open

**Edge Cases and Boundary Conditions:**
- Empty request bodies, null values, missing required fields
- Maximum/minimum value boundaries (e.g., string length limits, numeric ranges)
- Special characters and Unicode in text fields
- Large payloads approaching size limits
- Concurrent requests with race conditions
- Database constraint violations (foreign key, unique, not null)

## API Development Workflow

When developing a new API endpoint, follow this mandatory pattern:

#### Testing Requirements

- Write comprehensive tests for all controller return scenarios
- Test coverage must include:
  - Success case (200 OK with expected data)
  - All error cases (validation failures, authentication/authorization failures, business logic errors)
  - Edge cases and boundary conditions
- Reference existing tests in `internal/http/controller/**/*_test.go` for patterns
- Example: `internal/http/controller/user/user_register_test.go` demonstrates testing multiple return paths

#### Rationale

Integration tests validate real-world behavior and catch configuration errors, connection issues, and integration bugs that unit tests miss. Writing tests first forces clear requirements and prevents untested code. Comprehensive mock verification helps ensure that, along the tested execution paths, the system not only returns the expected results but also calls the right collaborators with the right arguments at the right time.
