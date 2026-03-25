---
applyTo: '**/*.go'
---

## Core Principles

### Integration-First Testing (Test-Driven Development)

- Write contract or integration tests before implementation for critical flows
- Every critical user or business flow MUST have at least one integration or end-to-end test without mocking core infrastructure
- Mock or stub only external systems and side-effect boundaries such as network calls, messaging systems, file systems, time, randomness, OS services, or repository/gateway-backed datastores
- Do NOT mock internal domain logic such as entities, value objects, or pure domain services
