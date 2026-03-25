# go-backend-framework

A **production-ready Go backend framework** built on the modular [go-backend-package](https://github.com/chan-jui-huang/go-backend-package) library with best practices in mind.

This framework provides a solid foundation for building scalable REST APIs with:
- **Security-first design**: Ed25519 JWT (quantum-resistant), Argon2id password hashing (GPU-resistant)
- **Modular infrastructure**: 12 reusable packages with zero cross-dependencies
- **Complete stack**: Authentication, authorization, database ORM, caching, logging, job scheduling
- **Production features**: Graceful shutdown, signal handling, connection pooling, log rotation
- **Developer experience**: Hot-reload, migrations, testing utilities, API documentation

## Table of Contents

- [Quick Start](#quick-start)
- [Architecture Overview](#architecture-overview)
- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [Core Concepts](#core-concepts)
- [Common Tasks](#common-tasks)
- [Testing](#testing)
- [Dependencies](#dependencies)
- [Additional Resources](#additional-resources)

---

## Quick Start

### 1. Create a New Project from the Framework

Visit the [release page](https://github.com/chan-jui-huang/go-backend-framework/releases) to find available versions.

```bash
# Set your desired version
VERSION=2.X.X  # Replace with actual version, e.g., 2.5.0

# Download and extract the installer
curl -L https://github.com/chan-jui-huang/go-backend-framework/archive/refs/tags/v$VERSION.tar.gz | \
  tar -zxv --strip-components=1 go-backend-framework-$VERSION/gbf-installer.sh

# Run the interactive installer
./gbf-installer.sh

# Compile the framework tools and utilities
cd <your-project-name> && make
```

### 2. Prerequisites

Before running the application, install the following:

| Requirement | Purpose |
|------------|---------|
| **MySQL** (v8.0+) or **PostgreSQL** (v17+) | Production database |
| **SQLite** (v3.0+) | Running tests locally |
| **Go** (v1.25.4+) | Language runtime |
| **Make** | Build automation |

### 3. Configure Environment

1. **Set up environment variables:**
   ```bash
   # Copy and edit the configuration
   cp .env.example .env
   # Edit .env with your database credentials and secrets
   ```

2. **Run database migrations:**
   ```bash
   make mysql-migration args=up  # For MySQL
   make pgsql-migration args=up  # For PostgreSQL
   ```

3. **Seed initial data:**
   ```bash
   ./bin/rdbms_seeder         # Populate base relational data
   ./bin/permission_seeder    # Set up access control policies (Casbin)
   ```

4. **Generate JWT secrets:**
   ```bash
   ./bin/jwt                  # Generate JWT for development
   ./bin/jwt -env=.env.test   # Generate JWT for test runs
   ./bin/jwt -env=.env.staging  # Generate JWT for staging
   ```

### 4. Run Tests

```bash
make test  # Run full test suite
# Or test specific packages:
make test args=./internal/http/controller/user
```

### 5. Start the Application

```bash
# Option A: With hot reloading (requires 'air' tool)
make air

# Option B: Direct execution
go run cmd/app/*

# Option C: Production build
make all && ./bin/app
```

**Verify the server is running:**
```bash
curl http://127.0.0.1:8080/api/ping
```

---

## Architecture Overview

The framework follows the standard [Go Project Layout](https://github.com/golang-standards/project-layout) with a clean separation of concerns:

### Layered Design

```
HTTP Request
    ↓
[Middleware Layer] - Authentication, Authorization, Logging, Rate Limiting
    ↓
[Controller Layer] - HTTP Handlers, Request Validation
    ↓
[Service Layer] - Business Logic (internal/pkg)
    ↓
[Data Layer] - GORM Models, Database Access
    ↓
Database
```

### Key Technologies

| Component | Technology | Purpose |
|-----------|-----------|---------|
| **Core Framework** | [go-backend-package](https://github.com/chan-jui-huang/go-backend-package) | 12 modular infrastructure packages (zero cross-dependencies) |
| **Lifecycle & Bootstrap** | `fx`, `internal/registrar`, `internal/deps` | Application lifecycle, dependency injection, and typed dependency access |
| **HTTP Server** | [Gin](https://github.com/gin-gonic/gin) | Fast, lightweight web framework |
| **Database ORM** | [GORM](https://github.com/go-gorm/gorm) + `pkg/database` | Type-safe database operations with connection pooling |
| **Authentication** | `pkg/authentication` (Ed25519) | Quantum-resistant JWT token implementation |
| **Password Hashing** | `pkg/argon2` | GPU-resistant Argon2id password hashing |
| **Authorization** | [Casbin](https://github.com/casbin/casbin) | Flexible role-based access control (RBAC, ABAC) |
| **Logging** | `pkg/logger` + [Zap](https://github.com/uber-go/zap) | Structured logging with rotation |
| **Configuration** | `pkg/booter/config` + [Viper](https://github.com/spf13/viper) | YAML and environment-based configuration with variable expansion |
| **Caching** | `pkg/redis` + [Redis](https://github.com/redis/go-redis) | In-memory caching with pooling |
| **Job Scheduling** | `pkg/scheduler` | Cron-style jobs with second-precision timing |
| **Validation** | [Validator](https://github.com/go-playground/validator) | Struct field validation |
| **Migrations** | [Goose](https://github.com/pressly/goose) | Database versioning & migrations |
| **API Docs** | [Swagger/Swag](https://github.com/swaggo/swag) | Auto-generated REST API documentation |

### Dependency Wiring

- **Config Snapshot**: YAML configuration with environment variable expansion (`${VAR_NAME}`) is loaded through registrar providers and stored in `internal/deps`.
- **Service Snapshot**: Initialized services are built by `fx` providers and exposed to existing code through typed `internal/deps` accessors.

### Bootstrap Sequence

```
load env/autoload → build booter config → fx.Supply + fx.Provide
                  → fx.Invoke(RegisterConfigDependencies, RegisterServiceDependencies)
                  → fx lifecycle hooks start validator, scheduler, and HTTP server
```

### Application Lifecycle Hooks (`fx`)

| Hook        | Registration               | Purpose                                                                  |
| ----------- | -------------------------- | ------------------------------------------------------------------------ |
| Validator   | `fx.OnStart`               | Register request-validation tag name handling before requests are served |
| Scheduler   | `fx.OnStart` / `fx.OnStop` | Start and stop background jobs                                           |
| HTTP Server | `fx.OnStart` / `fx.OnStop` | Run the Gin server and perform graceful shutdown                         |

See `cmd/app/main.go` for the canonical `fx` wiring example.

### Key Architectural Patterns

| Pattern                      | Purpose                                                                                     |
| ---------------------------- | ------------------------------------------------------------------------------------------- |
| **Fx Provider Pattern**      | Modular initialization is composed through registrar constructors and invokes               |
| **Dependency Injection**     | `fx` builds dependencies, and `internal/deps` exposes typed accessors for runtime/test code |
| **Configuration Management** | YAML + `.env` file support with variable expansion                                          |
| **Database Layer**           | Goose migrations + GORM ORM wrapper for type-safe operations                                |
| **Authentication**           | Ed25519 JWT (quantum-resistant) generated via `./bin/jwt`                                   |
| **Authorization**            | Casbin RBAC policies seeded via `cmd/kit/permission_seeder`                                 |
| **Logging**                  | Zap structured logging with rotation to `storage/log/`                                      |

---

## Development Setup

### Local Development Server

Use [air](https://github.com/cosmtrek/air) for hot-reloading during development:

```bash
# Review the configuration
cat .air.toml.example

# Copy or customize for your environment
cp .air.toml.example .air.toml

# Start with hot reloading
make air
```

This automatically restarts the server when you modify Go source files.

### Makefile Commands

```bash
# Building & Running
make all              # Build app and all helper tools (jwt, seeders, etc.)
make run              # Run app directly with race detector
make air              # Run with hot reloading (requires 'air' tool)
make debug-app        # Build debug version with race detector
make clean            # Clean bin directory

# Testing & Quality
make test [args=...]  # Run tests (default: all packages)
make benchmark        # Run benchmark tests
make linter           # Run code linter (golangci-lint)

# Database Migrations
make mysql-migration args=up     # MySQL: forward
make mysql-migration args=down   # MySQL: rollback
make pgsql-migration args=up     # PostgreSQL: forward
make sqlite-migration args=up    # SQLite: forward (for testing)
make clickhouse-migration args=up # ClickHouse: forward

# Code Generation
make swagger   # Generate Swagger/OpenAPI documentation

# Helper Tools (build only)
make jwt                    # Build jwt token generator
make rdbms_seeder          # Build database seeder
make http_route            # Build route lister
make permission_seeder     # Build policy seeder
```

---

## Project Structure

```
go-backend-framework/
├── cmd/                              # CLI commands and entry points
│   ├── app/                          # Main HTTP server
│   │   └── main.go                   # Application bootstrap
│   └── kit/                          # Helper utilities
│       ├── rdbms_seeder/             # Populate base relational data
│       ├── permission_seeder/        # Set up Casbin policies
│       ├── jwt/                      # JWT token generation
│       └── http_route/               # List and validate routes
│
├── internal/                         # Private application code
│   ├── http/                         # HTTP server & routing layer
│   │   ├── controller/               # HTTP handlers by domain
│   │   │   ├── user/                 # User endpoints
│   │   │   └── admin/                # Admin endpoints
│   │   ├── middleware/               # Request/response middleware
│   │   │   ├── authentication_middleware.go
│   │   │   ├── authorization_middleware.go
│   │   │   ├── rate_limit_middleware.go
│   │   │   └── ...
│   │   ├── response/                 # HTTP response helpers
│   │   │   ├── response.go           # Response builders
│   │   │   └── error_message.go      # Error codes & messages
│   │   ├── route/                    # Route definitions
│   │   │   └── api_route.go          # Route registration
│   │   └── server.go                 # Server configuration
│   │
│   ├── pkg/                          # Reusable business logic (internal use only)
│   │   ├── user/                     # User domain logic
│   │   ├── permission/               # Permission domain logic
│   │   ├── database/                 # Database helper logic
│   │   ├── model/                    # GORM model definitions
│   │   └── ...
│   │
│   ├── migration/                    # Database versioning
│   │   ├── rdbms/                    # SQL migrations (MySQL, PostgreSQL)
│   │   │   ├── 202308xx_create_users_table.sql
│   │   │   └── ...
│   │   ├── rdbms/test/               # RDBMS migrations for testing
│   │   ├── rdbms/seeder/             # Relational seeders
│   │   └── clickhouse/               # ClickHouse migrations
│   │
│   ├── registrar/                    # Dependency injection & wiring
│   │   ├── app.go                    # fx lifecycle registration
│   │   ├── http.go                   # HTTP server providers
│   │   ├── scheduler.go              # Scheduler providers
│   │   └── ...
│   │
│   ├── scheduler/                    # Background jobs & cron tasks
│   │   ├── scheduler.go
│   │   └── job/
│   │       └── example_job.go
│   │
│   └── test/                         # Test utilities & fixtures
│       ├── test.go                   # Public test facade
│       ├── runtime/
│       │   └── runtime.go            # Runtime bootstrap & lifecycle
│       ├── fixture/
│       │   ├── http/                 # HTTP request helpers and handler fixture
│       │   ├── db/                   # RDBMS / ClickHouse migration fixtures
│       │   ├── domain/               # User and permission domain fixtures
│       │   └── scenario/             # End-to-end scenario fixtures
│       └── fake/                     # Fixture data presets
│
├── pkg/                              # Public, reusable packages
│   └── ...                           # (Future: extracted as standalone modules)
│
├── docs/                             # Generated API documentation
│   ├── swagger.json
│   ├── swagger.yaml
│   └── docs.go
│
├── storage/                          # Runtime data storage
│   └── log/                          # Application logs
│       ├── app.log
│       └── access.log
│
├── deployment/                       # Container & deployment configs
│   └── docker/
│       └── Dockerfile
│
├── config.yml                        # Default/local configuration
├── config.test.yml                   # Test configuration
├── config.staging.yml                # Staging configuration
├── config.production.yml             # Production configuration
├── .env.example                      # Environment template
├── go.mod & go.sum                   # Dependency management
├── Makefile                          # Build automation
└── README.md                         # This file
```

---

## Core Concepts

The framework's architecture is built on **go-backend-package** patterns. For comprehensive details, refer to:

📚 **[go-backend-package AGENTS.md](https://github.com/chan-jui-huang/go-backend-package/blob/main/AGENTS.md)**

### Business Logic Placement

- Put reusable business logic in `internal/pkg/`
- Keep controllers, schedulers, and CLI commands focused on input handling, orchestration, and output formatting
- When a flow is used by a single interface only, keep it in that interface layer until reuse justifies extraction

### Dependency Access

- Register configuration and services through `fx` providers in `internal/registrar/`
- Access runtime dependencies through typed helpers in `internal/deps/`
- Keep bootstrap and lifecycle wiring in `cmd/app/main.go` and registrar packages, not inside business logic

### Test Runtime

- Use `internal/test/runtime` to bootstrap test environments
- Use fixtures from `internal/test/fixture/` and preset data from `internal/test/fake/`
- Prefer integration-style controller tests that exercise the real HTTP stack and database setup

---

## Common Tasks

### Adding a New API Endpoint

1. **Add or extend business logic** in `internal/pkg/<domain>/` when the flow is reusable across interfaces.
2. **Create the handler** in `internal/http/controller/<domain>/` and keep it focused on validation, orchestration, and response formatting.
3. **Register routes** in `internal/http/route/<domain>/api_route.go`, then add new routers to `internal/http/route/api_route.go` when needed.
4. **Add controller tests** beside the handler and use `internal/test/runtime`, `internal/test/fixture`, and `internal/test/fake` to exercise the real HTTP stack.
5. **Update permissions** in `cmd/kit/permission_seeder/permission_seeder.go` and mirror any required test fixtures under `internal/test/`.
6. **Regenerate Swagger docs** after API contract changes:
   ```bash
   make swagger
   ```

### Running Database Migrations

```bash
# Apply all pending migrations
make mysql-migration args=up

# Rollback one migration
make mysql-migration args=down

# Create a new migration (SQL format)
make mysql-migration args="create add_new_column sql"

# Migrate to a specific version
make mysql-migration args="up-to 202311251234"

# Check migration status
make mysql-migration args=status

# Reset all migrations (rollback everything)
make mysql-migration args=reset
```

### Viewing API Documentation

After running `make swagger`, visit:
```
http://localhost:8080/swagger/index.html
```

The Swagger UI displays all endpoints, request/response schemas, and allows testing.

### Logging

Structured logs are written to `storage/log/app.log` and `storage/log/access.log` through the registrar-wired logging setup.

For details, refer to: [go-backend-package/pkg/logger](https://github.com/chan-jui-huang/go-backend-package/tree/main/pkg/logger)

### Caching with Redis

Redis is initialized during bootstrap and is available to application packages through typed dependencies and registrar-wired services.

For details, refer to: [go-backend-package/pkg/redis](https://github.com/chan-jui-huang/go-backend-package/tree/main/pkg/redis)

### Password Hashing

Use **Argon2id** (GPU-resistant) for password hashing. See example in `internal/http/controller/user/user_register.go`.

For details, refer to: [go-backend-package/pkg/argon2](https://github.com/chan-jui-huang/go-backend-package/tree/main/pkg/argon2)

### Background Job Scheduling

The scheduler with **second-precision timing** is started and stopped through `fx` lifecycle hooks. Add jobs in `internal/scheduler/job/`.

For details, refer to: [go-backend-package/pkg/scheduler](https://github.com/chan-jui-huang/go-backend-package/tree/main/pkg/scheduler)

### Pagination & Error Handling

For data pagination and stacktrace extraction, use utilities from go-backend-package:

- **Pagination**: [pkg/pagination](https://github.com/chan-jui-huang/go-backend-package/tree/main/pkg/pagination)
- **Stacktrace**: [pkg/stacktrace](https://github.com/chan-jui-huang/go-backend-package/tree/main/pkg/stacktrace)

---

## Testing

### Test Execution

```bash
# Run all tests
make test

# Run tests for a specific package
make test args=./internal/http/controller/user

# Run a specific test function
go test -run TestUserLogin ./internal/http/controller/user

# Run with coverage report
go test -cover ./...

# Generate HTML coverage report
go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out
```

### Test Setup

Tests use helpers from `internal/test/` with testify's suite pattern:

```go
package user_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/chan-jui-huang/go-backend-framework/v2/internal/http/controller/user"
    "github.com/chan-jui-huang/go-backend-framework/v2/internal/test"
    "github.com/chan-jui-huang/go-backend-framework/v2/internal/test/fake"
    "github.com/stretchr/testify/suite"
)

type UserLoginTestSuite struct {
    suite.Suite
    runtime *test.Runtime
}

// SetupSuite runs once before all tests
func (suite *UserLoginTestSuite) SetupSuite() {
    suite.runtime = test.NewRdbmsRuntime(suite.T())
    suite.runtime.Users.Register(fake.User())
}

// Test the login endpoint
func (suite *UserLoginTestSuite) Test() {
    reqBody := user.UserLoginRequest{
        Email:    fake.User().Email,
        Password: fake.User().Password,
    }

    reqBodyBytes, _ := json.Marshal(reqBody)
    req := httptest.NewRequest("POST", "/api/user/login", bytes.NewReader(reqBodyBytes))
    suite.runtime.HTTP.AddCsrfToken(req)
    resp := httptest.NewRecorder()

    suite.runtime.HTTP.ServeHTTP(resp, req)

    suite.Equal(http.StatusOK, resp.Code)
}

func (suite *UserLoginTestSuite) TearDownSuite() {
    suite.runtime.Close()
}

// Run all tests in the suite
func TestUserLoginTestSuite(t *testing.T) {
    suite.Run(t, new(UserLoginTestSuite))
}
```

### Testing Best Practices

- **Test both success and failure paths**
- **Use table-driven tests** for parametric testing with multiple similar scenarios
- **Keep tests independent** (no shared state)
- **Use runtime and fixtures** from `internal/test/` to avoid duplication
- **Mock external services** (APIs, payment providers) — use [uber-go/mock](https://github.com/uber-go/mock)
- **Run `make linter`** before committing to catch style issues

---

## Dependencies

### Core Infrastructure

Built on **[go-backend-package](https://github.com/chan-jui-huang/go-backend-package)** with 12 modular packages (zero cross-dependencies):

**Bootstrap & Configuration**: [pkg/booter](https://github.com/chan-jui-huang/go-backend-package/tree/main/pkg/booter)

**Data & Storage**: [pkg/database](https://github.com/chan-jui-huang/go-backend-package/tree/main/pkg/database), [pkg/redis](https://github.com/chan-jui-huang/go-backend-package/tree/main/pkg/redis), [pkg/clickhouse](https://github.com/chan-jui-huang/go-backend-package/tree/main/pkg/clickhouse)

**Security**: [pkg/authentication](https://github.com/chan-jui-huang/go-backend-package/tree/main/pkg/authentication), [pkg/argon2](https://github.com/chan-jui-huang/go-backend-package/tree/main/pkg/argon2), [pkg/random](https://github.com/chan-jui-huang/go-backend-package/tree/main/pkg/random)

**Operations**: [pkg/logger](https://github.com/chan-jui-huang/go-backend-package/tree/main/pkg/logger), [pkg/scheduler](https://github.com/chan-jui-huang/go-backend-package/tree/main/pkg/scheduler), [pkg/pagination](https://github.com/chan-jui-huang/go-backend-package/tree/main/pkg/pagination), [pkg/stacktrace](https://github.com/chan-jui-huang/go-backend-package/tree/main/pkg/stacktrace)

### Framework Stack

HTTP: [Gin](https://github.com/gin-gonic/gin) | ORM: [GORM](https://github.com/go-gorm/gorm) | Auth: [Casbin](https://github.com/casbin/casbin)

Validation: [Validator](https://github.com/go-playground/validator), [Form](https://github.com/go-playground/form), [Mapstructure](https://github.com/mitchellh/mapstructure)

Migrations: [Goose](https://github.com/pressly/goose) | Docs: [Swag](https://github.com/swaggo/swag)

Dev Tools: [Air](https://github.com/cosmtrek/air), [GolangCI-Lint](https://github.com/golangci/golangci-lint), [Testify](https://github.com/stretchr/testify)

---

## Additional Resources

- **[Go Project Layout](https://github.com/golang-standards/project-layout)** - Recommended project structure
- **[Gin Documentation](https://gin-gonic.com/)** - HTTP framework guide
- **[GORM Documentation](https://gorm.io/)** - ORM reference
- **[Casbin Documentation](https://casbin.org/)** - Authorization guide
- **[JWT.io](https://jwt.io/)** - JWT specification and tools
- **[Swagger/OpenAPI](https://swagger.io/)** - API documentation standard
