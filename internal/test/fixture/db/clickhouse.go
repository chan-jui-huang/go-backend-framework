package db

import (
	"database/sql"
	"fmt"
	"path"

	ch "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/booter"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/clickhouse"
	"github.com/pressly/goose/v3"
	"go.uber.org/fx"
)

type ClickhouseMigrationDependencies struct {
	fx.In

	BooterConfig     *booter.Config
	ClickhouseConfig *clickhouse.Config
	Clickhouse       driver.Conn
}

type ClickhouseMigration struct {
	dir              string
	clickhouseConfig *clickhouse.Config
	clickhouse       driver.Conn
}

func NewClickhouseMigration(dependencies ClickhouseMigrationDependencies) *ClickhouseMigration {
	return &ClickhouseMigration{
		dir:              path.Join(dependencies.BooterConfig.RootDir, "internal/migration/clickhouse/test"),
		clickhouseConfig: dependencies.ClickhouseConfig,
		clickhouse:       dependencies.Clickhouse,
	}
}

func (cm *ClickhouseMigration) Conn() driver.Conn {
	return cm.clickhouse
}

func (cm *ClickhouseMigration) Run(callbacks ...func()) {
	clickhouseConfig := cm.clickhouseConfig
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

	if err := goose.Up(conn, cm.dir, goose.WithAllowMissing()); err != nil {
		panic(err)
	}

	for _, callback := range callbacks {
		callback()
	}
}

func (cm *ClickhouseMigration) Reset() {
	if cm == nil {
		return
	}

	clickhouseConfig := cm.clickhouseConfig
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
