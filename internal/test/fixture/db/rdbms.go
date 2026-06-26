package db

import (
	"path"

	"github.com/casbin/casbin/v3"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/booter"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/database"
	"github.com/pressly/goose/v3"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type RdbmsMigrationDependencies struct {
	fx.In

	BooterConfig   *booter.Config
	DatabaseConfig *database.Config
	Database       *gorm.DB
	CasbinEnforcer *casbin.SyncedCachedEnforcer
}

type RdbmsMigration struct {
	dir            string
	databaseConfig *database.Config
	database       *gorm.DB
	casbinEnforcer *casbin.SyncedCachedEnforcer
}

func NewRdbmsMigration(dependencies RdbmsMigrationDependencies) *RdbmsMigration {
	return &RdbmsMigration{
		dir:            path.Join(dependencies.BooterConfig.RootDir, "internal/migration/rdbms/test"),
		databaseConfig: dependencies.DatabaseConfig,
		database:       dependencies.Database,
		casbinEnforcer: dependencies.CasbinEnforcer,
	}
}

func (rm *RdbmsMigration) Database() *gorm.DB {
	return rm.database
}

func (rm *RdbmsMigration) Enforcer() *casbin.SyncedCachedEnforcer {
	return rm.casbinEnforcer
}

func (rm *RdbmsMigration) Run(callbacks ...func()) {
	db, err := rm.database.DB()
	if err != nil {
		panic(err)
	}

	if err := goose.SetDialect(string(rm.databaseConfig.Driver)); err != nil {
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

	db, err := rm.database.DB()
	if err != nil {
		panic(err)
	}

	if err := goose.SetDialect(string(rm.databaseConfig.Driver)); err != nil {
		panic(err)
	}

	if err := goose.Reset(db, rm.dir); err != nil {
		panic(err)
	}

	err = rm.database.Exec("DELETE FROM casbin_rules").Error
	if err != nil {
		panic(err)
	}

	if err := rm.casbinEnforcer.LoadPolicy(); err != nil {
		panic(err)
	}
}
