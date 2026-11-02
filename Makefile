include .env

.PHONY: all clean test security security-go-version security-govulncheck security-osv
SHELL:=/bin/bash

BIN_DIR:=${PWD}/bin
CMD_DIR:=${PWD}/cmd
KIT_DIR:=${CMD_DIR}/kit

APP_DIR:=${CMD_DIR}/app
APP_BIN:=app


GOFILES:=$(shell find . -type f -name "*.go")
TAGS:="jsoniter"

OBJECTS:=jwt rdbms_seeder http_route permission_seeder

all:${OBJECTS}
	go build -o ${BIN_DIR}/${APP_BIN} -v -tags ${TAGS} -ldflags "-s -w" ${APP_DIR}

run:
	go run -race -tags ${TAGS} ${APP_DIR}

debug-app:
	go build -o ${BIN_DIR}/$@ -v -race -tags ${TAGS} ${APP_DIR}

jwt:
	go build -o ${BIN_DIR}/$@ -v -race -ldflags "-s -w" ${KIT_DIR}/$@

rdbms_seeder:
	go build -o ${BIN_DIR}/$@ -v -race -ldflags "-s -w" ${KIT_DIR}/$@

http_route:
	go build -o ${BIN_DIR}/$@ -v -race -ldflags "-s -w" ${KIT_DIR}/$@

permission_seeder:
	go build -o ${BIN_DIR}/$@ -v -race -ldflags "-s -w" ${KIT_DIR}/$@

clean:
	rm -rf ${BIN_DIR}

mysql-migration:
	go tool goose -dir internal/migration/rdbms -allow-missing mysql "${DB_USERNAME}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_DATABASE}?parseTime=true&loc=UTC" ${args}

pgsql-migration:
	go tool goose -dir internal/migration/rdbms -allow-missing postgres "user=${DB_USERNAME} password=${DB_PASSWORD} host=${DB_HOST} port=${DB_PORT} dbname=${DB_DATABASE} sslmode=disable" ${args}

sqlite-migration:
	go tool goose -dir internal/migration/rdbms -allow-missing sqlite3 ${DB_DATABASE} ${args}

clickhouse-migration:
	go tool goose -dir internal/migration/clickhouse -allow-missing clickhouse "tcp://${CLICKHOUSE_ADDR_01}/${CLICKHOUSE_DATABASE}?username=${CLICKHOUSE_USERNAME}&password=${CLICKHOUSE_PASSWORD}" ${args}

linter:
	go tool golangci-lint run ./...

swagger:
	go tool swag init --pd -g cmd/app/main.go

air:
	go tool air

security: security-go-version security-govulncheck security-osv

# Check all known standard library and toolchain advisories for the active Go release.
security-go-version:
	@set -e; \
	go_version=$$(go env GOVERSION); \
	go_version=$${go_version#go}; \
	if [[ ! "$$go_version" =~ ^[0-9]+\.[0-9]+\.[0-9]+$$ ]]; then \
		printf 'Cannot check unsupported Go release version: %s\n' "$$go_version" >&2; \
		exit 1; \
	fi; \
	scan_dir=$$(mktemp -d); \
	trap 'rm -rf "$$scan_dir"' EXIT; \
	printf 'Checking Go %s: standard library and toolchain\n' "$$go_version"; \
	printf '{"results":[{"packages":[{"package":{"ecosystem":"Go","name":"stdlib","version":"%s"}},{"package":{"ecosystem":"Go","name":"toolchain","version":"%s"}}]}]}\n' "$$go_version" "$$go_version" > "$$scan_dir/osv-scanner.json"; \
	go run github.com/google/osv-scanner/v2/cmd/osv-scanner@v2.7.0 scan source --no-call-analysis=go --lockfile="osv-scanner:$$scan_dir/osv-scanner.json"

security-govulncheck:
	go tool govulncheck -scan=package -test -tags ${TAGS} ./...
	go tool govulncheck -scan=package -tags ${TAGS} tool

# Scan the complete selected module graph, including tool dependencies.
security-osv:
	@set -e; \
	scan_dir=$$(mktemp -d); \
	trap 'rm -rf "$$scan_dir"' EXIT; \
	{ \
		printf 'module security-scan\n\nrequire (\n'; \
		go list -m -f '{{if not .Main}}{{if .Replace}}{{if .Replace.Version}}{{.Replace.Path}} {{.Replace.Version}}{{end}}{{else}}{{.Path}} {{.Version}}{{end}}{{end}}' all; \
		printf ')\n'; \
	} > "$$scan_dir/go.mod"; \
	go run github.com/google/osv-scanner/v2/cmd/osv-scanner@v2.7.0 scan source --no-call-analysis=go --lockfile="$$scan_dir/go.mod"

test:
	$(eval args?=./...)
	ENV=test go test ${args}

benchmark:
	$(eval args?=./...)
	ENV=test go test -bench=. -run=none -benchmem ${args}
