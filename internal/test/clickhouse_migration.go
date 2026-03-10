package test

import (
	"database/sql"
	"fmt"
	"path"

	ch "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/chan-jui-huang/go-backend-package/pkg/booter"
	"github.com/chan-jui-huang/go-backend-package/pkg/booter/config"
	"github.com/chan-jui-huang/go-backend-package/pkg/clickhouse"
	"github.com/pressly/goose/v3"
)

type clickhouseMigration struct {
	dir string
}

var ClickhouseMigration *clickhouseMigration

func NewClickhouseMigration() *clickhouseMigration {
	booterConfig := config.Registry.Get("booter").(booter.Config)

	return &clickhouseMigration{
		dir: path.Join(booterConfig.RootDir, "internal/migration/clickhouse/test"),
	}
}

func (cm *clickhouseMigration) Run(callbacks ...func()) {
	clickhouseConfig := config.Registry.Get("clickhouse").(clickhouse.Config)
	conn, err := sql.Open("clickhouse", fmt.Sprintf("tcp://%s?username=%s&password=%s", clickhouseConfig.Addr[0], clickhouseConfig.Username, clickhouseConfig.Password))
	if err != nil {
		panic(err)
	}

	if _, err := conn.Exec(
		"CREATE DATABASE IF NOT EXISTS {database:Identifier}",
		ch.Named("database", clickhouseConfig.Database),
	); err != nil {
		panic(err)
	}

	if err := conn.Close(); err != nil {
		panic(err)
	}

	conn, err = sql.Open("clickhouse", fmt.Sprintf("tcp://%s?username=%s&password=%s&database=%s", clickhouseConfig.Addr[0], clickhouseConfig.Username, clickhouseConfig.Password, clickhouseConfig.Database))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	if err := goose.SetDialect(string(goose.DialectClickHouse)); err != nil {
		panic(err)
	}

	if err := goose.Up(conn, cm.dir); err != nil {
		panic(err)
	}

	for _, callback := range callbacks {
		callback()
	}
}

func (cm *clickhouseMigration) Reset() {
	clickhouseConfig := config.Registry.Get("clickhouse").(clickhouse.Config)
	conn, err := sql.Open("clickhouse", fmt.Sprintf("tcp://%s?username=%s&password=%s&database=%s", clickhouseConfig.Addr[0], clickhouseConfig.Username, clickhouseConfig.Password, clickhouseConfig.Database))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	if err := goose.SetDialect(string(goose.DialectClickHouse)); err != nil {
		panic(err)
	}

	if err := goose.Reset(conn, cm.dir); err != nil {
		panic(err)
	}
}
