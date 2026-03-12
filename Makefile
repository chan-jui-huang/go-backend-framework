include .env

.PHONY: all clean test
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

test:
	$(eval args?=./...)
	ENV=test go test ${args}

benchmark:
	$(eval args?=./...)
	ENV=test go test -bench=. -run=none -benchmem ${args}
