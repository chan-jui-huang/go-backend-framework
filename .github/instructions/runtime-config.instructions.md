---
applyTo: '**/*.yml,**/*.yaml,.env*,Makefile,Dockerfile,deployment/**'
---

## Dependency Wiring

- **Config Snapshot**: YAML configuration with environment variable expansion (`${VAR_NAME}`) is loaded through registrar providers and stored in `internal/deps`.
- **Service Snapshot**: Initialized services are built by `fx` providers and exposed to existing code through typed `internal/deps` accessors.

## Bootstrap Sequence

```
load env/autoload → build booter config → fx.Supply + fx.Provide
                  → fx.Invoke(RegisterConfigDependencies, RegisterServiceDependencies)
                  → fx lifecycle hooks start validator, scheduler, and HTTP server
```

## Application Lifecycle Hooks (`fx`)

| Hook        | Registration               | Purpose                                                                  |
| ----------- | -------------------------- | ------------------------------------------------------------------------ |
| Validator   | `fx.OnStart`               | Register request-validation tag name handling before requests are served |
| Scheduler   | `fx.OnStart` / `fx.OnStop` | Start and stop background jobs                                           |
| HTTP Server | `fx.OnStart` / `fx.OnStop` | Run the Gin server and perform graceful shutdown                         |

See `cmd/app/main.go` for the canonical `fx` wiring example.

## Key Architectural Patterns

| Pattern                      | Purpose                                                                                     |
| ---------------------------- | ------------------------------------------------------------------------------------------- |
| **Fx Provider Pattern**      | Modular initialization is composed through registrar constructors and invokes               |
| **Dependency Injection**     | `fx` builds dependencies, and `internal/deps` exposes typed accessors for runtime/test code |
| **Configuration Management** | YAML + `.env` file support with variable expansion                                          |
| **Database Layer**           | Goose migrations + GORM ORM wrapper for type-safe operations                                |
| **Authentication**           | Ed25519 JWT (quantum-resistant) generated via `./bin/jwt`                                   |
| **Authorization**            | Casbin RBAC policies seeded via `cmd/kit/permission_seeder`                                 |
| **Logging**                  | Zap structured logging with rotation to `storage/log/`                                      |

## Architecture Constraints

### Technical Requirements

#### Data Persistence

- MUST support ACID transactions for business-critical operations
- MUST use connection pooling for database access
- MUST provide schema versioning and migration capabilities
- MUST support both relational and time-series data storage patterns when needed

#### Observability

- MUST provide structured logging with log rotation
- MUST support configuration via environment variables and configuration files
- MUST generate machine-readable API documentation
