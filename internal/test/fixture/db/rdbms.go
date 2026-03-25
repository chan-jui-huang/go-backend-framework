package db

import (
	"path"

	"github.com/chan-jui-huang/go-backend-framework/v2/internal/deps"
	"github.com/pressly/goose/v3"
)

type RdbmsMigration struct {
	dir string
}

func NewRdbmsMigration() *RdbmsMigration {
	booterConfig := deps.BooterConfig()

	return &RdbmsMigration{
		dir: path.Join(booterConfig.RootDir, "internal/migration/rdbms/test"),
	}
}

func (rm *RdbmsMigration) Run(callbacks ...func()) {
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

func (rm *RdbmsMigration) Reset() {
	if rm == nil {
		return
	}

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
