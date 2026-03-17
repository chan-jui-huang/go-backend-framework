---
applyTo: 'cmd/**/*.go,internal/deps/**/*.go,internal/registrar/**/*.go,internal/http/server.go,internal/scheduler/**/*.go,**/*.yml,**/*.yaml,.env*,Makefile,Dockerfile,deployment/**'
---

## Runtime Constraints

#### Data Persistence

- MUST support ACID transactions for business-critical operations
- MUST use connection pooling for database access
- MUST provide schema versioning and migration capabilities
- MUST support both relational and time-series data storage patterns when needed

#### Observability

- MUST provide structured logging with log rotation
- MUST support configuration via environment variables and configuration files
- MUST generate machine-readable API documentation

#### Configuration

- Runtime configuration MUST load from YAML and environment variables
- Environment-backed configuration values SHOULD use variable expansion instead of hardcoded secrets

#### Runtime Integrations

- Authentication changes MUST remain compatible with the project's Ed25519 JWT flow
- Authorization changes MUST integrate with the existing Casbin RBAC policy model
- Runtime logging MUST use the project's structured logging pipeline and rotation behavior
