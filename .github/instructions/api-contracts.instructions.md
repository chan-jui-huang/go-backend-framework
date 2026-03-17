---
applyTo: 'internal/http/**/*.go,docs/**/*.go,docs/**/*.json,docs/**/*.yaml'
---

## Business Logic Layer Architecture

### API Development Workflow

When developing a new API endpoint, follow this mandatory pattern:

#### Controller Implementation

- Create or extend handler files under `internal/http/controller/<area>/`
- Controllers should ONLY:
  - Validate HTTP requests
  - Call business logic from `internal/pkg/`
  - Format HTTP responses
- Never put business logic directly in controllers
