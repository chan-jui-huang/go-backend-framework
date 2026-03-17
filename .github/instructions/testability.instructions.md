---
applyTo: '**/*.go'
---

## Core Principles

### Integration-First Testing (Test-Driven Development)

Tests MUST be written BEFORE implementation (Red-Green-Refactor). Integration tests take priority over isolated unit tests.

#### Test-Driven Development Process

**Non-negotiable rules:**

- Write contract/integration tests FIRST before any implementation

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

#### Rationale

Integration tests validate real-world behavior and catch configuration errors, connection issues, and integration bugs that unit tests miss. Writing tests first forces clear requirements and prevents untested code. Comprehensive mock verification helps ensure that, along the tested execution paths, the system not only returns the expected results but also calls the right collaborators with the right arguments at the right time.
