---
applyTo: '**/*_test.go'
---

## Core Principles

### Integration-First Testing (Test-Driven Development)

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

- Every controller test suite MUST cover success, validation failure, authentication failure, authorization failure, not found, business error, and server error scenarios when those outcomes are reachable
- Tests SHOULD include malformed input, relevant boundary conditions, and other scenario-specific cases for the handler as needed

## API Development Workflow

#### Testing Requirements

- Write comprehensive tests for all controller return scenarios
- Test coverage must include:
  - Success case (200 OK with expected data)
  - All error cases (validation failures, authentication/authorization failures, business logic errors)
  - Edge cases, boundary conditions, and other relevant scenarios
- Reference existing tests in `internal/http/controller/**/*_test.go` for patterns
- Example: `internal/http/controller/user/user_register_test.go` demonstrates testing multiple return paths
