package test

import (
	"path"

	"github.com/chan-jui-huang/go-backend-framework/v2/internal/deps"
	"github.com/pressly/goose/v3"
)

type rdbmsMigration struct {
	dir string
}

var RdbmsMigration *rdbmsMigration

func NewRdbmsMigration() *rdbmsMigration {
	booterConfig := deps.BooterConfig()

	return &rdbmsMigration{
		dir: path.Join(booterConfig.RootDir, "internal/migration/rdbms/test"),
	}
}

func (rm *rdbmsMigration) Run(callbacks ...func()) {
	databaseConfig := deps.DatabaseConfig()
	database := deps.Database()
	db, err := database.DB()
	if err != nil {
		panic(err)
	}

	if err := goose.SetDialect(string(databaseConfig.Driver)); err != nil {
		panic(err)
	}
	if err := goose.Up(db, rm.dir); err != nil {
		panic(err)
	}

	for _, callback := range callbacks {
		callback()
	}
}

func (rm *rdbmsMigration) Reset() {
	databaseConfig := deps.DatabaseConfig()
	database := deps.Database()
	db, err := database.DB()
	if err != nil {
		panic(err)
	}

	if err := goose.SetDialect(string(databaseConfig.Driver)); err != nil {
		panic(err)
	}

	if err := goose.Reset(db, rm.dir); err != nil {
		panic(err)
	}

	err = database.Exec("DELETE FROM casbin_rules").Error
	if err != nil {
		panic(err)
	}

	enforcer := deps.CasbinEnforcer()
	if err := enforcer.LoadPolicy(); err != nil {
		panic(err)
	}
}
