package test

import (
	"path"

	"github.com/casbin/casbin/v3"
	"github.com/chan-jui-huang/go-backend-package/pkg/booter"
	"github.com/chan-jui-huang/go-backend-package/pkg/booter/config"
	"github.com/chan-jui-huang/go-backend-package/pkg/booter/service"
	"github.com/chan-jui-huang/go-backend-package/pkg/database"
	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

type rdbmsMigration struct {
	dir string
}

var RdbmsMigration *rdbmsMigration

func NewRdbmsMigration() *rdbmsMigration {
	booterConfig := config.Registry.Get("booter").(booter.Config)

	return &rdbmsMigration{
		dir: path.Join(booterConfig.RootDir, "internal/migration/rdbms/test"),
	}
}

func (rm *rdbmsMigration) Run(callbacks ...func()) {
	databaseConfig := config.Registry.Get("database").(database.Config)
	database := service.Registry.Get("database").(*gorm.DB)
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
	databaseConfig := config.Registry.Get("database").(database.Config)
	database := service.Registry.Get("database").(*gorm.DB)
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

	enforcer := service.Registry.Get("casbinEnforcer").(*casbin.SyncedCachedEnforcer)
	if err := enforcer.LoadPolicy(); err != nil {
		panic(err)
	}
}
