---
name: make-migration
description: Operate this repository's Makefile-backed goose migration targets safely, including inspecting status, creating migrations, validating files, and applying or rolling back schema changes.
compatibility: Requires the go-backend-framework workspace, root Makefile, Go toolchain, goose tool dependency, and valid database environment variables from `.env` or the active shell.
metadata:
  version: '0.1.0'
---

# Make Migration

## Use When

Use this skill when the user asks to run, inspect, create, validate, apply, or roll back migrations through this repository's `make *-migration` targets.

## Repository Targets

The root `Makefile` wraps `go tool goose` with `-allow-missing` and fixed migration directories:

- `make mysql-migration args='<command>'`: RDBMS migrations in `internal/migration/rdbms`, MySQL DSN from `DB_USERNAME`, `DB_PASSWORD`, `DB_HOST`, `DB_PORT`, and `DB_DATABASE`.
- `make pgsql-migration args='<command>'`: RDBMS migrations in `internal/migration/rdbms`, PostgreSQL DSN from the same `DB_*` variables.
- `make sqlite-migration args='<command>'`: RDBMS migrations in `internal/migration/rdbms`, SQLite database path from `DB_DATABASE`.
- `make clickhouse-migration args='<command>'`: ClickHouse migrations in `internal/migration/clickhouse`, DSN from `CLICKHOUSE_ADDR_01`, `CLICKHOUSE_DATABASE`, `CLICKHOUSE_USERNAME`, and `CLICKHOUSE_PASSWORD`.

## Workflow

1. Inspect `Makefile` before operating, because target wiring or directories may change.
2. Choose the target from the requested database; default to `sqlite-migration` only when the user asks for local/test SQLite or the repository context clearly indicates SQLite.
3. For discovery, run `make <target> args='-h'` or `go tool goose -h`; these only print goose help and target command expansion.
4. Before destructive operations, run `make <target> args='status'` unless the user explicitly asks to skip it.
5. Treat `down`, `down-to`, `redo`, and `reset` as destructive. Confirm intent unless the user already gave an explicit rollback/reset instruction.
6. Create new SQL migrations with `make <target> args='create <name> sql'`; keep generated timestamps unique and use the same timestamp only for intentionally paired prod/test migrations.
7. Validate migration files with `make <target> args='validate'` before applying.
8. Apply pending migrations with `make <target> args='up'`, or use `up-by-one` / `up-to VERSION` when limiting scope.
9. After apply or rollback, run `make <target> args='status'` again and summarize the resulting version/status.

## Goose Commands

Supported commands observed from `go tool goose -h`:

- `up`: migrate to the most recent version.
- `up-by-one`: migrate up by one.
- `up-to VERSION`: migrate to a specific version.
- `down`: roll back one version.
- `down-to VERSION`: roll back to a specific version.
- `redo`: re-run the latest migration.
- `reset`: roll back all migrations.
- `status`: show migration status.
- `version`: print current database version.
- `create NAME [sql|go]`: create a timestamped migration file.
- `fix`: apply sequential ordering to migration filenames.
- `validate`: check migration files without running them.

## Useful Options

Pass options inside `args` after the Makefile's fixed driver and DSN:

- `-h`: print goose help.
- `-v`: enable verbose output.
- `-no-color`: disable colored output.
- `-timeout <duration>`: limit query duration, such as `-timeout 30s`.
- `-table <name>`: use a non-default migration version table.
- `-no-versioning`: apply migrations without version tracking, in file order.
- `-s`: use sequential numbering for new migrations.
- MySQL-only TLS options: `-certfile`, `-ssl-cert`, `-ssl-key`.

The Makefile already supplies `-dir` and `-allow-missing`; do not duplicate or override them unless the task explicitly requires changing the Makefile.

## Test Migration Rule

- Repository test migrations under `internal/migration/rdbms/test/` are replayed from an empty database in the test harness.
- Prefer the simplest valid schema change for the test file.
- If the change can be expressed with `ALTER TABLE ... ADD COLUMN ...`, use that directly in the test migration.
- Do not rebuild the whole table or write extra copy/drop/rename SQL unless the database engine or schema change genuinely requires it.
- Mirror the production migration intent, but keep the test migration minimal when the empty-database startup makes additional compatibility SQL unnecessary.

## Example Prompts

```text
$make-migration Check the current SQLite migration status, validate files, then apply pending migrations if validation passes.
```

```text
$make-migration Create a new MySQL SQL migration named add_user_indexes, then show me the generated file path.
```

```text
$make-migration Roll back ClickHouse migrations down to version 20260513090000 after showing current status.
```
